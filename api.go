package reflectopenapi

import (
	"context"
	"encoding/json"
	"os"
	"reflect"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/podhmo/reflect-openapi/pkg/comment"
)

func NewDoc() (*openapi3.Swagger, error) {
	skeleton := []byte(`{
  "openapi": "3.0.0",
  "info": {
    "title": "Sample API",
    "description": "-",
    "version": "0.0.0"
  },
  "servers": [
    {
      "url": "http://localhost:8888",
      "description": "local development server"
    },
  ],
  "paths": {}
}`)
	return NewDocFromSkeleton(skeleton)
}

// TODO: add api function
func NewDocFromSkeleton(skeleton []byte) (*openapi3.Swagger, error) {
	l := openapi3.NewSwaggerLoader()
	return l.LoadSwaggerFromData(skeleton)
}

type Config struct {
	Doc      *openapi3.Swagger
	Resolver Resolver

	StrictSchema        bool // if true, use `{additionalProperties: false}` as default
	SkipValidation      bool // if true, skip validation for api doc definition
	SkipExtractComments bool // if true, skip extracting comments as a description

	DefaultError            interface{}
	IsRequiredCheckFunction func(reflect.StructTag) bool // handling required, default is always false
}

func (c *Config) BuildDoc(ctx context.Context, use func(m *Manager)) (*openapi3.Swagger, error) {
	if c.Resolver == nil {
		resolver := &UseRefResolver{}
		if c.StrictSchema {
			ng := false
			resolver.AdditionalPropertiesAllowed = &ng
		}
		c.Resolver = resolver
	}

	if c.Doc == nil {
		doc, err := NewDoc()
		if err != nil {
			return nil, err
		}
		c.Doc = doc
	}

	v := NewVisitor(c.Resolver)
	if c.IsRequiredCheckFunction != nil {
		v.Transformer.IsRequired = c.IsRequiredCheckFunction
	}
	if !c.SkipExtractComments {
		v.CommentLookup = comment.NewLookup()
	}

	m := &Manager{
		Doc:      c.Doc,
		Resolver: c.Resolver,
		Visitor:  v,
	}

	use(m)

	if c.DefaultError != nil {
		errSchema := v.VisitType(c.DefaultError)
		responseRef := &openapi3.ResponseRef{
			Value: openapi3.NewResponse().
				WithDescription("default error").
				WithJSONSchemaRef(errSchema),
		}
		for _, op := range v.Operations {
			if val, ok := op.Responses["default"]; !ok || val.Value == nil || val.Value.Description == nil || *val.Value.Description != "" {
				continue
			}
			op.Responses["default"] = responseRef
		}
	}

	if b, ok := c.Resolver.(Binder); ok {
		b.Bind(m.Doc)
	}

	if !c.SkipValidation {
		if err := m.Doc.Validate(ctx); err != nil {
			return m.Doc, err
		}
	}
	return m.Doc, nil
}

func (c *Config) EmitDoc(use func(m *Manager)) {
	ctx := context.Background()
	doc, err := c.BuildDoc(ctx, use)
	if err != nil {
		panic(err)
	}
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	if err := enc.Encode(doc); err != nil {
		panic(err)
	}
}

type Manager struct {
	Visitor  *Visitor
	Resolver Resolver

	Doc *openapi3.Swagger
}
