package reflectopenapi

import (
	"context"

	"github.com/getkin/kin-openapi/openapi3"
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

	StrictSchema   bool // if true, use `{additionalProperties: false}` as default
	SkipValidation bool // if true, skip validation for api doc definition
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

	m := &Manager{
		Doc:      c.Doc,
		Resolver: c.Resolver,
		Visitor:  NewVisitor(c.Resolver),
	}
	use(m)

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

type Manager struct {
	Visitor  *Visitor
	Resolver Resolver

	Doc *openapi3.Swagger
}
