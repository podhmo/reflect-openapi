package reflectopenapi

import (
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
	NameStore *NameStore

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

	name := v.Title // after VisitType()
	if name == "" {
		name = s.GetName()
	}
	return r.NameStore.GetOrCreatePair(v, name).Ref
}

type NameStore struct {
	PairMap map[string][]*refPair
}

func NewNameStore() *NameStore {
	return &NameStore{
		PairMap: map[string][]*refPair{},
	}
}

type refPair struct {
	Def *openapi3.SchemaRef
	Ref *openapi3.SchemaRef
}

func (ns *NameStore) GetOrCreatePair(v *openapi3.Schema, name string) *refPair {
	pairs, existed := ns.PairMap[name]
	if existed {
		if len(pairs) == 1 && pairs[0].Def.Value == v {
			return pairs[0]
		}
		for _, pair := range pairs {
			if pair.Def.Value == v {
				return pair
			}
		}
	}

	pair := &refPair{
		Def: &openapi3.SchemaRef{Value: v},
		Ref: &openapi3.SchemaRef{Ref: "#/components/schemas/" + name, Value: v},
	}

	ns.PairMap[name] = append(ns.PairMap[name], pair)
	return pair
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
	if len(r.NameStore.PairMap) == 0 {
		return
	}

	if doc.Components.Schemas == nil {
		doc.Components.Schemas = map[string]*openapi3.SchemaRef{}
	}

	schemas := doc.Components.Schemas
	for name, pairs := range r.NameStore.PairMap {
		for _, pair := range pairs {
			schemas[name] = pair.Def
		}
	}
}
