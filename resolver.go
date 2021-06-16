package reflectopenapi

import (
	"fmt"
	"strings"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/podhmo/reflect-openapi/pkg/shape"
)

// without ref

type NoRefResolver struct {
	AdditionalPropertiesAllowed *bool // set as Config.StrictSchema
}

func (r *NoRefResolver) ResolveSchema(v *openapi3.Schema, s shape.Shape) *openapi3.SchemaRef {
	switch s.(type) {
	case shape.Primitive, shape.Container:
		return &openapi3.SchemaRef{Value: v}
	default:
		if r.AdditionalPropertiesAllowed != nil && v.Type == "object" {
			v.AdditionalPropertiesAllowed = r.AdditionalPropertiesAllowed
		}
		return &openapi3.SchemaRef{Value: v}
	}
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

// with ref

type UseRefResolver struct {
	Schemas                     []*openapi3.SchemaRef
	AdditionalPropertiesAllowed *bool // set as Config.StrictSchema
}

func (r *UseRefResolver) ResolveSchema(v *openapi3.Schema, s shape.Shape) *openapi3.SchemaRef {
	useOriginalDef := false
	switch s.(type) {
	case shape.Primitive, shape.Container:
		if len(v.Extensions) == 0 {
			useOriginalDef = true
		}
	}
	if useOriginalDef {
		return &openapi3.SchemaRef{Value: v}
	}

	if r.AdditionalPropertiesAllowed != nil && v.Type == "object" {
		v.AdditionalPropertiesAllowed = r.AdditionalPropertiesAllowed
	}
	if s.GetName() == "" {
		return &openapi3.SchemaRef{Value: v}
	}

	name := v.Title
	if name == "" {
		name = s.GetName()
	}
	ref := fmt.Sprintf("#/components/schemas/%s", name)
	r.Schemas = append(r.Schemas, &openapi3.SchemaRef{Ref: ref, Value: v})
	return &openapi3.SchemaRef{Ref: ref, Value: v}
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

func (r *UseRefResolver) Bind(doc *openapi3.T) {
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
