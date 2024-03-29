package reflectopenapi

import (
	"fmt"
	"log"
	"reflect"
	"strings"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/podhmo/reflect-openapi/info"
	shape "github.com/podhmo/reflect-shape"
)

type Direction string

const (
	DirectionInput     Direction = "input"
	DirectionOutput    Direction = "output"
	DirectionInternal  Direction = "internal"
	DirectionParameter Direction = "parameter"
)

// without ref

type NoRefResolver struct {
	AdditionalPropertiesAllowed *bool // set as Config.StrictSchema
}

var _ Resolver = (*NoRefResolver)(nil)

func (r *NoRefResolver) ResolveSchema(v *openapi3.Schema, s *shape.Shape, typ Direction) *openapi3.SchemaRef {
	if r.AdditionalPropertiesAllowed != nil && v.Type == "object" && s.Kind == reflect.Struct && s.Type.NumField() > 0 {
		v.AdditionalProperties.Has = r.AdditionalPropertiesAllowed
	}
	return &openapi3.SchemaRef{Value: v}
}
func (r *NoRefResolver) ResolveParameter(v *openapi3.Parameter, s *shape.Shape) *openapi3.ParameterRef {
	return &openapi3.ParameterRef{Value: v}
}
func (r *NoRefResolver) ResolveRequestBody(v *openapi3.RequestBody, s *shape.Shape) *openapi3.RequestBodyRef {
	return &openapi3.RequestBodyRef{Value: v}
}
func (r *NoRefResolver) ResolveResponse(v *openapi3.Response, s *shape.Shape) *openapi3.ResponseRef {
	return &openapi3.ResponseRef{Value: v}
}

// with ref

type UseRefResolver struct {
	*NameStore // for Binder

	DisableInputRef             bool
	DisableOutputRef            bool
	AdditionalPropertiesAllowed *bool // set as Config.StrictSchema
}

var _ Resolver = (*UseRefResolver)(nil)
var _ Binder = (*UseRefResolver)(nil)

func (r *UseRefResolver) ResolveSchema(v *openapi3.Schema, s *shape.Shape, direction Direction) (ref *openapi3.SchemaRef) {
	useOriginalDef := false
	switch s.Kind {
	case reflect.Struct, reflect.Interface:
	default:
		if len(v.Extensions) == 0 {
			useOriginalDef = true
		}
	}
	if useOriginalDef {
		return &openapi3.SchemaRef{Value: v}
	}
	if r.AdditionalPropertiesAllowed != nil && v.Type == "object" && s.Kind == reflect.Struct && s.Type.NumField() > 0 {
		v.AdditionalProperties.Has = r.AdditionalPropertiesAllowed
	}
	if s.Name == "" {
		return &openapi3.SchemaRef{Value: v}
	}

	name := v.Title // after VisitType()
	if name == "" {
		name = s.Name
	}

	if (r.DisableInputRef && direction == DirectionInput) || r.DisableOutputRef && direction == DirectionOutput {
		if v.Extensions == nil {
			v.Extensions = map[string]interface{}{}
		}
		v.Extensions["x-go-id"] = s.FullName()
		return &openapi3.SchemaRef{Value: v}
	}
	return r.NameStore.GetOrCreatePair(v, name, s).Ref
}

func (r *UseRefResolver) ResolveParameter(v *openapi3.Parameter, s *shape.Shape) *openapi3.ParameterRef {
	return &openapi3.ParameterRef{Value: v}
}
func (r *UseRefResolver) ResolveRequestBody(v *openapi3.RequestBody, s *shape.Shape) *openapi3.RequestBodyRef {
	return &openapi3.RequestBodyRef{Value: v}
}
func (r *UseRefResolver) ResolveResponse(v *openapi3.Response, s *shape.Shape) *openapi3.ResponseRef {
	return &openapi3.ResponseRef{Value: v}
}

type RefPair struct {
	Name  string
	Shape *shape.Shape

	Def *openapi3.SchemaRef
	Ref *openapi3.SchemaRef
}

type NameStore struct {
	Prefix     string
	OnConflict func(*RefPair, int)

	pairMap map[string][]*RefPair
	info    *info.Info
}

func NewNameStore() *NameStore {
	ns := &NameStore{
		Prefix:  "#/components/schemas/",
		pairMap: map[string][]*RefPair{},
	}
	ns.OnConflict = ns.fixPairAsAddingSuffix
	return ns
}

func (ns *NameStore) fixPairAsAddingSuffix(pair *RefPair, i int) {
	if i > 0 {
		name := fmt.Sprintf("%s%02d", pair.Name, i)
		pair.Name = name
		pair.Ref.Ref = ns.Prefix + name
	}

	if pair.Def.Value.Extensions == nil {
		pair.Def.Value.Extensions = map[string]interface{}{}
	}
	pair.Def.Value.Extensions["x-go-id"] = pair.Shape.FullName() // FIXME: what is x-go-id?
}

func (ns *NameStore) GetOrCreatePair(v *openapi3.Schema, name string, shape *shape.Shape) *RefPair {
	// normalize name
	if strings.Contains(name, "[") {
		name = strings.ReplaceAll(name, "[", "_")
	}
	if strings.Contains(name, "]") {
		name = strings.ReplaceAll(name, "]", "")
	}

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
	if ns.info != nil {
		ns.info.SchemaValue[pair.Ref] = pair.Def.Value
	}

	ns.pairMap[name] = append(ns.pairMap[name], pair)
	return pair
}

func (ns *NameStore) BindSchemas(doc *openapi3.T) {
	if len(ns.pairMap) == 0 {
		return
	}
	if doc.Components == nil {
		doc.Components = &openapi3.Components{}
	}
	if doc.Components.Schemas == nil {
		doc.Components.Schemas = map[string]*openapi3.SchemaRef{}
	}
	schemas := doc.Components.Schemas

	for name, pairs := range ns.pairMap {
		if len(pairs) > 1 {
			for i, pair := range pairs {
				ns.OnConflict(pair, i)
				log.Printf("name conflict is occured, fix %s -> %s (%s)", name, pair.Name, pair.Shape.FullName())
			}
		}

		for _, pair := range pairs {
			schemas[pair.Name] = pair.Def
			if ns.info != nil {
				if sinfo, ok := ns.info.SchemaInfo[pair.Def.Value]; ok {
					sinfo.Name = pair.Name
				}
			}
		}
	}
}
