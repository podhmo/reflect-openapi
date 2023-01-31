package reflectopenapi

import (
	"reflect"
	"strconv"
	"strings"

	"github.com/getkin/kin-openapi/openapi3"
	shape "github.com/podhmo/reflect-shape"
)

// TODO: extra information
// TODO: json tag inline,omitempty support
// TODO: schema nullable support (?)

// not visitor pattern
type Visitor struct {
	*Transformer

	Doc        *openapi3.T
	Schemas    map[int]*openapi3.Schema
	Operations map[int]*openapi3.Operation
}

func isRequiredDefault(tag reflect.StructTag) bool {
	s, ok := tag.Lookup("required")
	if !ok {
		return false
	}
	v, _ := strconv.ParseBool(s)
	return v
}

func NewVisitor(tagNameOption TagNameOption, resolver Resolver, selector Selector, extractor Extractor) *Visitor {
	transformer := (&Transformer{
		TagNameOption:    tagNameOption,
		cache:            map[int]interface{}{},
		interceptFuncMap: map[reflect.Type]func(*shape.Shape) *openapi3.Schema{},
		Resolver:         resolver,
		IsRequired:       isRequiredDefault,
		Selector:         selector,
		Extractor:        extractor,
	}).Builtin()
	if t, ok := selector.(needTransformer); ok {
		t.NeedTransformer(transformer)
	}
	return &Visitor{
		Transformer: transformer,
		Schemas:     map[int]*openapi3.Schema{},
		Operations:  map[int]*openapi3.Operation{},
	}
}

func (v *Visitor) VisitType(in *shape.Shape, modifiers ...func(*openapi3.Schema)) *openapi3.SchemaRef {
	out := v.Transform(in).(*openapi3.Schema)
	out.Title = in.Name
	for _, m := range modifiers {
		m(out)
	}

	id := in.Number
	v.Schemas[id] = out
	if len(modifiers) > 0 {
		if out.Extensions == nil {
			out.Extensions = map[string]interface{}{
				v.TagNameOption.XNewTypeTag: in.FullName(),
			}
		}
		if doc := in.Named().Doc(); doc != "" {
			out.Description = doc
		}
		v.Transformer.cache[id] = out
	}
	return v.ResolveSchema(out, in)
}
func (v *Visitor) VisitFunc(in *shape.Shape, modifiers ...func(*openapi3.Operation)) *openapi3.Operation {
	out := v.Transform(in).(*openapi3.Operation)

	if doc := in.Func().Doc(); doc != "" {
		out.Description = doc
		out.Summary = strings.SplitN(doc, "\n", 2)[0]
	}

	for _, m := range modifiers {
		m(out)
	}

	v.Operations[in.Number] = out
	return out
}
