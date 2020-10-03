package reflectopenapi

import (
	"fmt"
	"strings"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/podhmo/reflect-openapi/pkg/shape"
)

type Resolver interface {
	ResolveSchema(v *openapi3.Schema, s shape.Shape) *openapi3.SchemaRef
	ResolveParameter(v *openapi3.Parameter, s shape.Shape) *openapi3.ParameterRef
	ResolveRequestBody(v *openapi3.RequestBody, s shape.Shape) *openapi3.RequestBodyRef
	ResolveResponse(v *openapi3.Response, s shape.Shape) *openapi3.ResponseRef
}

type Binder interface {
	Bind(doc *openapi3.Swagger)
}

// without ref

type NoRefResolver struct{}

func (r *NoRefResolver) ResolveSchema(v *openapi3.Schema, s shape.Shape) *openapi3.SchemaRef {
	return &openapi3.SchemaRef{Value: v}
}
func (r *NoRefResolver) ResolveParameter(v *openapi3.Parameter, s shape.Shape) *openapi3.ParameterRef {
	return &openapi3.ParameterRef{Value: v}
}
func (r *NoRefResolver) ResolveRequestBody(v *openapi3.RequestBody, s shape.Shape) *openapi3.RequestBodyRef {
	return &openapi3.RequestBodyRef{Value: v}
}
func (r *NoRefResolver) ResolveResponse(v *openapi3.Response, s shape.Shape) *openapi3.ResponseRef {
	return &openapi3.ResponseRef{Value: v}
}
func (r *NoRefResolver) Bind(doc *openapi3.Swagger) {
}

// with ref

type UseRefResolver struct {
	Schemas []*openapi3.SchemaRef
}

func (r *UseRefResolver) ResolveSchema(v *openapi3.Schema, s shape.Shape) *openapi3.SchemaRef {
	switch s := s.(type) {
	case shape.Primitive, shape.Container:
		return &openapi3.SchemaRef{Value: v}
	default:
		ref := fmt.Sprintf("#/components/schemas/%s", s.GetName())
		r.Schemas = append(r.Schemas, &openapi3.SchemaRef{Ref: ref, Value: v})
		return &openapi3.SchemaRef{Ref: ref, Value: v}
	}
}
func (r *UseRefResolver) ResolveParameter(v *openapi3.Parameter, s shape.Shape) *openapi3.ParameterRef {
	return &openapi3.ParameterRef{Value: v}
}
func (r *UseRefResolver) ResolveRequestBody(v *openapi3.RequestBody, s shape.Shape) *openapi3.RequestBodyRef {
	return &openapi3.RequestBodyRef{Value: v}
}
func (r *UseRefResolver) ResolveResponse(v *openapi3.Response, s shape.Shape) *openapi3.ResponseRef {
	return &openapi3.ResponseRef{Value: v}
}

func (r *UseRefResolver) Bind(doc *openapi3.Swagger) {
	if len(r.Schemas) == 0 {
		return
	}

	if doc.Components.Schemas == nil {
		doc.Components.Schemas = map[string]*openapi3.SchemaRef{}
	}
	for _, ref := range r.Schemas {
		ref := ref
		path := ref.Ref
		ref.Ref = ""
		doc.Components.Schemas[strings.TrimPrefix(path, "#/components/schemas/")] = ref
	}
}
