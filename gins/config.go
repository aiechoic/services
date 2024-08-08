package gins

import (
	"github.com/aiechoic/services/openapi"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"os"
	"path"
)

const DefaultConfigSection ConfigSection = "gin-service"

var defaultConfigData = []byte(
	`# Gin Service Config

# title for openapi
api_title: "API Documentation"
# version of the api
api_version: "1.0.0"
# api servers
api_servers:
  - url: "http://localhost:8080/api/v1"
    description: "Local Server"
# http port
http_port: 8080
# gin mode, can be "debug", "release", "test"
gin_mode: "debug"
# enable log
log: true
# enable cross-origin resource sharing
enable_cors: true
# api root
api_root: "/api/v1"
# configure this to serve static files
static_routes: []
#  - route: "/static"
#    dir: "./static"
#    not_found: "404.html" # "index.html" for vuejs
`)

type OpenAPIServer struct {
	// on production, this should be the actual url, e.g. https://api.example.com/api/v1
	Url         string `mapstructure:"url"`
	Description string `mapstructure:"description"`
}

type StaticRoute struct {
	Route    string `mapstructure:"route"`
	Dir      string `mapstructure:"dir"`
	NotFound string `mapstructure:"not_found"`
}

type Config struct {
	ApiTitle     string          `mapstructure:"api_title"`   // for openapi
	ApiVersion   string          `mapstructure:"api_version"` // for openapi
	ApiServers   []OpenAPIServer `mapstructure:"api_servers"` // for openapi
	HttpPort     int             `mapstructure:"http_port"`
	GinMode      string          `mapstructure:"gin_mode"`
	Log          bool            `mapstructure:"log"`
	EnableCORS   bool            `mapstructure:"enable_cors"`
	APIRoot      string          `mapstructure:"api_root"`
	StaticRoutes []StaticRoute   `mapstructure:"static_routes"`
}

func (c *Config) NewServer() *Server {
	engine, router := c.NewEngine()
	api := c.NewOpenAPI()
	return &Server{
		API:       api,
		Port:      c.HttpPort,
		Engine:    engine,
		APIRouter: router,
	}
}

func (g *Config) NewEngine() (*gin.Engine, gin.IRouter) {
	if g.GinMode == "" {
		gin.SetMode(g.GinMode)
	}
	r := gin.New()
	if g.Log {
		r.Use(gin.Logger())
	}
	r.Use(gin.Recovery())
	if g.EnableCORS {
		r.Use(cors.Default())
	}
	for _, sr := range g.StaticRoutes {
		if sr.NotFound != "" {
			r.GET(sr.Route+"/*filepath", g.TryServeFiles(sr))
		} else {
			r.Static(sr.Route, sr.Dir)
		}
	}
	var i gin.IRouter = r
	if g.APIRoot != "" {
		i = r.Group(g.APIRoot)
	}
	return r, i
}

// TryServeFiles attempts to serve static files from the specified directory.
// If the requested file does not exist, it serves a custom "not found" page.
func (g *Config) TryServeFiles(sr StaticRoute) gin.HandlerFunc {
	return func(c *gin.Context) {
		filePath := path.Join(sr.Dir, c.Param("filepath"))
		if _, err := os.Stat(filePath); os.IsNotExist(err) {
			notFoundPage := path.Join(sr.Dir, sr.NotFound)
			c.File(notFoundPage)
		} else {
			c.File(filePath)
		}
	}
}

func (g *Config) NewOpenAPI() *openapi.Openapi {
	var servers []*openapi.Server
	for _, s := range g.ApiServers {
		servers = append(servers, &openapi.Server{
			Url:         s.Url,
			Description: s.Description,
		})
	}
	info := &openapi.Info{
		Title:   g.ApiTitle,
		Version: g.ApiVersion,
	}
	return &openapi.Openapi{
		Openapi: "3.1.0",
		Info:    info,
		Servers: servers,
		Components: &openapi.Components{
			SecuritySchemes: map[string]*openapi.SecurityScheme{},
			Schemas:         map[string]*openapi.Schema{},
		},
		Paths: map[string]openapi.PathItem{},
	}
}
