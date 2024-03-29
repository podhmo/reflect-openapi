package reflectopenapi

import (
	"github.com/getkin/kin-openapi/openapi3"
	shape "github.com/podhmo/reflect-shape"
)

type Extractor interface {
	Extract(interface{}) *shape.Shape
}

type Selector interface {
	SelectInput(*shape.Func) (*shape.Shape, string)
	SelectOutput(*shape.Func) (*shape.Shape, string)
}

type Resolver interface {
	ResolveSchema(v *openapi3.Schema, s *shape.Shape, typ Direction) *openapi3.SchemaRef
	ResolveParameter(v *openapi3.Parameter, s *shape.Shape) *openapi3.ParameterRef
	ResolveRequestBody(v *openapi3.RequestBody, s *shape.Shape) *openapi3.RequestBodyRef
	ResolveResponse(v *openapi3.Response, s *shape.Shape) *openapi3.ResponseRef
}

type Binder interface {
	BindSchemas(doc *openapi3.T)
}

type needTransformer interface {
	NeedTransformer(*Transformer) // this is work-around
}
