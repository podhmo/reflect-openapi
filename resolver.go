package reflectopenapi

import (
	"fmt"

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
	return r.NameStore.GetOrCreatePair(v, name, s).Ref
}

type NameStore struct {
	Prefix  string
	pairMap map[string][]*RefPair
}

func NewNameStore() *NameStore {
	return &NameStore{
		Prefix:  "#/components/schemas/",
		pairMap: map[string][]*RefPair{},
	}
}

type RefPair struct {
	Name  string
	Shape shape.Shape

	Def *openapi3.SchemaRef
	Ref *openapi3.SchemaRef
}

func (ns *NameStore) GetOrCreatePair(v *openapi3.Schema, name string, shape shape.Shape) *RefPair {
	pairs, existed := ns.pairMap[name]
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

	pair := &RefPair{
		Name:  name,
		Shape: shape,
		Def:   &openapi3.SchemaRef{Value: v},
		Ref:   &openapi3.SchemaRef{Ref: ns.Prefix + name, Value: v},
	}

	ns.pairMap[name] = append(ns.pairMap[name], pair)
	return pair
}

func (ns *NameStore) FixRefs() { // TODO: more variation to fix name-conflict
	for _, pairs := range ns.pairMap {
		if len(pairs) == 1 {
			continue
		}

		pair := pairs[0]
		if pair.Def.Value.Extensions == nil {
			pair.Def.Value.Extensions = map[string]interface{}{}
		}
		pair.Def.Value.Extensions["x-go-id"] = pair.Shape.GetIdentity()

		for i, pair := range pairs[1:] {
			name := fmt.Sprintf("%s%02d", pair.Name, i+1)
			pair.Name = name
			pair.Ref.Ref = ns.Prefix + name

			if pair.Def.Value.Extensions == nil {
				pair.Def.Value.Extensions = map[string]interface{}{}
			}
			pair.Def.Value.Extensions["x-go-id"] = pair.Shape.GetIdentity()
		}
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

func (r *UseRefResolver) Bind(doc *openapi3.T) {
	if len(r.NameStore.pairMap) == 0 {
		return
	}

	if doc.Components.Schemas == nil {
		doc.Components.Schemas = map[string]*openapi3.SchemaRef{}
	}

	r.NameStore.FixRefs()

	schemas := doc.Components.Schemas
	for _, pairs := range r.NameStore.pairMap {
		for _, pair := range pairs {
			schemas[pair.Name] = pair.Def
		}
	}
}
