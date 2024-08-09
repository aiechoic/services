package openapi

import (
	"fmt"
	"reflect"
	"strings"
)

type Contact struct {
	Name  string `json:"name,omitempty"`
	Url   string `json:"url,omitempty"`
	Email string `json:"email,omitempty"`
}

type License struct {
	Name       string `json:"name,omitempty"`
	Url        string `json:"url,omitempty"`
	Identifier string `json:"identifier,omitempty"`
}

type Info struct {
	Title          string   `json:"title,omitempty"`
	Summary        string   `json:"summary,omitempty"`
	Description    string   `json:"description,omitempty"`
	TermsOfService string   `json:"termsOfService,omitempty"`
	Contact        *Contact `json:"contact,omitempty"`
	License        *License `json:"license,omitempty"`
	Version        string   `json:"version,omitempty"`
}

type Server struct {
	Url         string `json:"url,omitempty"`
	Description string `json:"description,omitempty"`
}

type SecuritySchemes map[string]*SecurityScheme

type SecuritySchemeType string

const (
	SecuritySchemeTypeApiKey    SecuritySchemeType = "apiKey"
	SecuritySchemeTypeHttp      SecuritySchemeType = "http"
	SecuritySchemeTypeOauth     SecuritySchemeType = "oauth2"
	SecuritySchemeTypeOpenId    SecuritySchemeType = "openIdConnect"
	SecuritySchemeTypeMutualTLS SecuritySchemeType = "mutualTLS"
)

type OAuthFlow struct {
	AuthorizationUrl string            `json:"authorizationUrl,omitempty"`
	TokenUrl         string            `json:"tokenUrl,omitempty"`
	RefreshUrl       string            `json:"refreshUrl,omitempty"`
	Scopes           map[string]string `json:"scopes,omitempty"`
}

type OAuthFlows struct {
	Implicit          *OAuthFlow `json:"implicit,omitempty"`
	Password          *OAuthFlow `json:"password,omitempty"`
	ClientCredentials *OAuthFlow `json:"clientCredentials,omitempty"`
	AuthorizationCode *OAuthFlow `json:"authorizationCode,omitempty"`
}

type SecurityScheme struct {
	Type             SecuritySchemeType `json:"type,omitempty"`
	Description      string             `json:"description,omitempty"`
	Name             string             `json:"name,omitempty"`
	In               string             `json:"in,omitempty"` // header, query, cookie
	Scheme           string             `json:"scheme,omitempty"`
	BearerFormat     string             `json:"bearerFormat,omitempty"`
	OpenIdConnectUrl string             `json:"openIdConnectUrl,omitempty"`
	Flows            *OAuthFlows        `json:"flows,omitempty"`
}

type Components struct {
	SecuritySchemes SecuritySchemes    `json:"securitySchemes,omitempty"`
	Schemas         map[string]*Schema `json:"schemas,omitempty"`
}

type ContentType string

const (
	ContentTypeJson ContentType = "application/json"
	ContentTypeXml  ContentType = "application/xml"
	ContentTypeForm ContentType = "application/x-www-form-urlencoded"
)

var contentStructTags = map[ContentType]string{
	ContentTypeJson: "json",
	ContentTypeXml:  "xml",
	ContentTypeForm: "form",
}

func (c ContentType) GetStructTag() string {
	return contentStructTags[c]
}

type MediaType struct {
	Schema  *Schema `json:"schema,omitempty"`
	Example any     `json:"example,omitempty"`
}

type RequestBody struct {
	Description string                     `json:"description,omitempty"`
	Content     map[ContentType]*MediaType `json:"content,omitempty"`
	Required    bool                       `json:"required,omitempty"`
}

type ResponseCode string

type Parameter struct {
	Name            string  `json:"name,omitempty"`
	In              string  `json:"in,omitempty"` // query, path, header, cookie
	Schema          *Schema `json:"schema,omitempty"`
	Description     string  `json:"description,omitempty"`
	Required        bool    `json:"required,omitempty"`
	Deprecated      bool    `json:"deprecated,omitempty"`
	AllowEmptyValue bool    `json:"allowEmptyValue,omitempty"`
}

type ResponseBody struct {
	Description string                     `json:"description,omitempty"`
	Content     map[ContentType]*MediaType `json:"content,omitempty"`
}

type Operation struct {
	Tags        []string                       `json:"tags,omitempty"`
	Summary     string                         `json:"summary,omitempty"`
	Description string                         `json:"description,omitempty"`
	RequestBody *RequestBody                   `json:"requestBody,omitempty"`
	Parameters  []*Parameter                   `json:"parameters,omitempty"`
	Responses   map[ResponseCode]*ResponseBody `json:"responses,omitempty"`
	Deprecated  bool                           `json:"deprecated,omitempty"`
	Security    []map[string][]string          `json:"security,omitempty"`
}

type PathItem map[string]*Operation

type Tag struct {
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
}

type Schema struct {
	Ref         string             `json:"$ref,omitempty"`
	Description string             `json:"description,omitempty"`
	Type        string             `json:"type,omitempty"`
	Enum        []string           `json:"enum,omitempty"`
	Format      string             `json:"format,omitempty"`
	Required    []string           `json:"required,omitempty"`
	Properties  map[string]*Schema `json:"properties,omitempty"`
	Items       *Schema            `json:"items,omitempty"`
}

func (s *Schema) RequiredFields(fields []string) {
	existRequired := map[string]bool{}
	for _, field := range s.Required {
		existRequired[field] = true
	}
	for _, field := range fields {
		if _, ok := existRequired[field]; !ok {
			s.Required = append(s.Required, field)
		}
	}
}

type Openapi struct {
	Openapi    string              `json:"openapi,omitempty"`
	Info       *Info               `json:"info,omitempty"`
	Servers    []*Server           `json:"servers,omitempty"`
	Components *Components         `json:"components,omitempty"`
	Paths      map[string]PathItem `json:"paths,omitempty"`
	Tags       []*Tag              `json:"tags,omitempty"`
}

func (o *Openapi) GetRefSchema(ref string) *Schema {
	name := strings.TrimPrefix(ref, "#/components/schemas/")
	s := o.Components.Schemas[name]
	if s == nil {
		panic(fmt.Sprintf("schema %s not found", name))
	}
	for pn, ps := range s.Properties {
		if ps.Ref != "" {
			s.Properties[pn] = o.GetRefSchema(ps.Ref)
		}
	}
	return s
}

func (o *Openapi) NewSchema(service string, obj any, ct ContentType) *Schema {
	tag := ct.GetStructTag()
	if tag == "" {
		panic(fmt.Sprintf("unsupported content type: %s", ct))
	}
	if obj == nil {
		return &Schema{
			Type: "null",
		}
	}
	rt := reflect.TypeOf(obj)
	return o.newSchema(service, rt, tag)
}

func (o *Openapi) newSchema(service string, rt reflect.Type, tag string) *Schema {
	if rt == nil {
		return &Schema{
			Type: "null",
		}
	}
	if rt.Kind() == reflect.Ptr {
		rt = rt.Elem()
	}
	rk := rt.Kind()
	switch rk {
	case reflect.String:
		return &Schema{
			Type: "string",
		}
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		s := &Schema{
			Type:   "integer",
			Format: "int32",
		}
		if rk == reflect.Int64 || rk == reflect.Uint64 {
			s.Format = "int64"
		}
		return s
	case reflect.Float32, reflect.Float64:
		s := &Schema{
			Type: "number",
		}
		if rk == reflect.Float64 {
			s.Format = "double"
		} else {
			s.Format = "float"
		}
		return s
	case reflect.Bool:
		return &Schema{
			Type: "boolean",
		}
	case reflect.Slice, reflect.Array:
		return &Schema{
			Type:  "array",
			Items: o.newSchema(service, rt.Elem(), tag),
		}
	case reflect.Struct:
		return o.newRefStructSchema(service, rt, tag)
	default:
		return &Schema{
			Type: "null",
		}
	}
}

func (o *Openapi) newRefStructSchema(service string, rt reflect.Type, tagName string) *Schema {
	name := rt.Name()
	if service != "" {
		name = fmt.Sprintf("%s-%s", service, name)
	}
	name = fmt.Sprintf("%s-%s", name, tagName)
	ref := fmt.Sprintf("#/components/schemas/%s", name)
	if o.Components == nil {
		o.Components = &Components{
			Schemas:         map[string]*Schema{},
			SecuritySchemes: SecuritySchemes{},
		}
	}
	if schema, ok := o.Components.Schemas[name]; ok {
		return &Schema{
			Ref: ref,
		}
	} else {
		schema = o.newStructSchema(service, rt, tagName)
		o.Components.Schemas[name] = schema
		return &Schema{
			Ref: ref,
		}
	}
}

func (o *Openapi) newStructSchema(service string, rt reflect.Type, tagName string) *Schema {

	schema := &Schema{
		Type:       "object",
		Properties: map[string]*Schema{},
	}
	for i := 0; i < rt.NumField(); i++ {
		field := rt.Field(i)
		if field.PkgPath != "" { // unexported
			continue
		}
		fieldName := field.Name
		// 如果 tag 中有 json:"-" 则忽略该字段
		if tag := field.Tag.Get(tagName); tag == "-" {
			continue
		} else if tag != "" {
			// 如果 tag 中有 json:"name" 则使用 name 作为字段名
			fieldName = strings.Split(tag, ",")[0]
		}
		prop := o.newSchema(service, field.Type, tagName)
		if field.Anonymous {
			if prop.Ref != "" {
				prop = o.GetRefSchema(prop.Ref)
			}
			if prop.Type == "object" && len(prop.Properties) > 0 {
				// Merge anonymous nested struct properties
				for nestedFieldName, nestedFieldSchema := range prop.Properties {
					if _, ok := schema.Properties[nestedFieldName]; !ok {
						schema.Properties[nestedFieldName] = nestedFieldSchema
					}
				}
				if len(prop.Required) > 0 {
					schema.RequiredFields(prop.Required)
				}
			}
		} else {
			if binding := field.Tag.Get("binding"); binding != "" {
				values := strings.Split(binding, ",")
				for _, value := range values {
					switch value {
					case "required":
						schema.Required = append(schema.Required, fieldName)
					case "email":
						prop.Format = "email"
					}
				}
			}
			if description := field.Tag.Get("description"); description != "" {
				prop.Description = description
			}
			if enum := field.Tag.Get("enum"); enum != "" {
				prop.Enum = strings.Split(enum, ",")
			}
			schema.Properties[fieldName] = prop
		}
	}
	return schema
}
