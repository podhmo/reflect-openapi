package reflectopenapi

import (
	"github.com/getkin/kin-openapi/openapi3"
	shape "github.com/podhmo/reflect-shape"
)

type Selector interface {
	SelectInput(*shape.Func) *shape.Shape
	SelectOutput(*shape.Func) *shape.Shape
}

type Resolver interface {
	ResolveSchema(v *openapi3.Schema, s *shape.Shape) *openapi3.SchemaRef
	ResolveParameter(v *openapi3.Parameter, s *shape.Shape) *openapi3.ParameterRef
	ResolveRequestBody(v *openapi3.RequestBody, s *shape.Shape) *openapi3.RequestBodyRef
	ResolveResponse(v *openapi3.Response, s *shape.Shape) *openapi3.ResponseRef
}

type Binder interface {
	BindSchemas(doc *openapi3.T)
}

// xxx
type useArglist interface {
	useArglist()
}
