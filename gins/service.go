package gins

import (
	"github.com/aiechoic/services/openapi"
	"github.com/gin-gonic/gin"
)

type Request struct {
	Description string
	Query       any // from url query
	Json        any // from body
	Form        any // from body
	Xml         any // from body
}

func (r *Request) getContents(service string, api *openapi.Openapi) map[openapi.ContentType]*openapi.MediaType {
	var contents = map[openapi.ContentType]*openapi.MediaType{}
	if r.Json != nil {
		contents[openapi.ContentTypeJson] = &openapi.MediaType{
			Schema: api.NewSchema(service, r.Json, openapi.ContentTypeJson),
		}
	}
	if r.Form != nil {
		contents[openapi.ContentTypeForm] = &openapi.MediaType{
			Schema: api.NewSchema(service, r.Form, openapi.ContentTypeForm),
		}
	}
	if r.Xml != nil {
		contents[openapi.ContentTypeXml] = &openapi.MediaType{
			Schema: api.NewSchema(service, r.Xml, openapi.ContentTypeXml),
		}
	}
	return contents
}

type Response struct {
	Description string
	Json        any
	Xml         any
}

func (r *Response) getContents(service string, api *openapi.Openapi) map[openapi.ContentType]*openapi.MediaType {
	var contents = map[openapi.ContentType]*openapi.MediaType{}
	if r.Json != nil {
		contents[openapi.ContentTypeJson] = &openapi.MediaType{
			Schema: api.NewSchema(service, r.Json, openapi.ContentTypeJson),
		}
	}
	if r.Xml != nil {
		contents[openapi.ContentTypeXml] = &openapi.MediaType{
			Schema: api.NewSchema(service, r.Xml, openapi.ContentTypeXml),
		}
	}
	return contents
}

type Route struct {
	Method      string
	Path        string
	Summary     string
	Description string
	Security    Security
	Handler     Handler
}

type Handler struct {
	Request  Request
	Response Response
	Handler  func(c *gin.Context)
}

type Security interface {
	Auth(c *gin.Context)
	SecurityScheme() []map[string][]string
}

func (r Route) Use(h Security) Route {
	r.Security = h
	return r
}

type Service struct {
	Tag         string
	Description string
	Path        string
	Routes      []Route
}
