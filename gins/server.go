package gins

import (
	"context"
	"errors"
	"fmt"
	"github.com/aiechoic/services/openapi"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"os"
	"os/signal"
	"slices"
	"strings"
	"syscall"
	"time"
)

type Server struct {
	API       *openapi.Openapi
	Port      int
	Engine    *gin.Engine
	APIRouter gin.IRouter
}

func (s *Server) Run(ctx context.Context) {
	address := fmt.Sprintf(":%d", s.Port)
	srv := &http.Server{
		Addr:    address,
		Handler: s.Engine,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	select {
	case <-quit:
		log.Println("received system signal, shutting down gracefully")
	case <-ctx.Done():
		log.Println("context cancelled, shutting down gracefully")
	}

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Fatal("server forced to shutdown:", err)
	}

	log.Println("server exiting")
}

func (s *Server) Register(services ...*Service) {
	for _, service := range services {
		s.register(service)
	}
}

func (s *Server) SetSecuritySchemes(schemes openapi.SecuritySchemes) {
	if s.API.Components.SecuritySchemes == nil {
		s.API.Components.SecuritySchemes = make(openapi.SecuritySchemes)
	}
	for name, scheme := range schemes {
		if _, ok := s.API.Components.SecuritySchemes[name]; ok {
			panic(fmt.Sprintf("security scheme %s already exists", name))
		}
		s.API.Components.SecuritySchemes[name] = scheme
	}
}

func (s *Server) register(service *Service) {
	o := s.API
	r := s.APIRouter
	o.Tags = append(o.Tags, &openapi.Tag{
		Name:        service.Tag,
		Description: service.Description,
	})
	for _, route := range service.Routes {
		// openapi spec requires method to be lowercase
		route.Method = strings.ToLower(route.Method)
		if service.Path == "/" {
			service.Path = ""
		}
		if !strings.HasPrefix(route.Path, "/") {
			route.Path = "/" + route.Path
		}
		path := service.Path + route.Path
		pathItem, ok := o.Paths[path]
		if !ok {
			pathItem = make(openapi.PathItem)
			o.Paths[path] = pathItem
		}
		op := &openapi.Operation{
			Tags:        []string{service.Tag},
			Summary:     route.Summary,
			Description: route.Description,
			Responses: map[openapi.ResponseCode]*openapi.ResponseBody{
				"200": {
					Content:     route.Handler.Response.getContents(service.Tag, o),
					Description: route.Handler.Response.Description,
				},
			},
		}
		if route.Security != nil {
			op.Security = route.Security.SecurityScheme()
		}
		if route.Handler.Request.Json != nil || route.Handler.Request.Form != nil {
			if route.Handler.Request.Json != nil && route.Handler.Request.Form != nil {
				panic(fmt.Sprintf(
					"service %s route %s: cannot have both json and form Body parameters",
					service.Tag, route.Path,
				))
			}
			if route.Method != "post" && route.Method != "put" {
				panic(fmt.Sprintf(
					"service %s route %s: request json/form Body parameter only allowed for POST and PUT methods, got %s",
					service.Tag, route.Path, route.Method,
				))
			}
			op.RequestBody = &openapi.RequestBody{
				Content:     route.Handler.Request.getContents(service.Tag, o),
				Description: route.Handler.Request.Description,
			}
		}
		if route.Handler.Request.Query != nil {
			schema := o.NewSchema(service.Tag, route.Handler.Request.Query, openapi.ContentTypeForm)
			if schema.Ref != "" {
				schema = o.GetRefSchema(schema.Ref)
			}
			for name, prop := range schema.Properties {
				op.Parameters = append(op.Parameters, &openapi.Parameter{
					Name:        name,
					In:          "query",
					Schema:      prop,
					Description: prop.Description,
					Required:    slices.Contains(schema.Required, name),
				})
			}
		}
		pathItem[route.Method] = op
		var handlers []gin.HandlerFunc
		if route.Security != nil {
			handlers = append(handlers, route.Security.Auth)
		}
		handlers = append(handlers, route.Handler.Handler)
		r.Handle(strings.ToUpper(route.Method), path, handlers...)
	}
}
