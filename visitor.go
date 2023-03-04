package reflectopenapi

import (
	"reflect"
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

	EnableAutoTag bool
}

func NewVisitor(tagNameOption TagNameOption, resolver Resolver, selector Selector, extractor Extractor) *Visitor {
	transformer := (&Transformer{
		TagNameOption:    tagNameOption,
		cache:            map[int]interface{}{},
		defaultValues:    map[int]reflect.Value{},
		interceptFuncMap: map[reflect.Type]func(*shape.Shape) *openapi3.Schema{},
		Resolver:         resolver,
		Selector:         selector,
		Extractor:        extractor,
	}).Builtin()
	transformer.IsRequired = transformer.isRequired
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

	if out.Default == nil {
		if in.Package.Path != "" && !shape.IsZeroRecursive(in.Type, in.DefaultValue) {
			out.Default = in.DefaultValue.Interface()
		}
	}
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
	return v.ResolveSchema(out, in, DirectionInternal)
}
func (v *Visitor) VisitFunc(in *shape.Shape, modifiers ...func(*openapi3.Operation)) *openapi3.Operation {
	out := v.Transform(in).(*openapi3.Operation)

	fn := in.Func()
	if doc := fn.Doc(); doc != "" {
		out.Description = doc
		out.Summary = strings.SplitN(doc, "\n", 2)[0]
	}

	for _, m := range modifiers {
		m(out)
	}

	if v.EnableAutoTag {
		out.Tags = append(out.Tags, fn.Shape.Package.Name)
	}

	v.Operations[in.Number] = out
	return out
}
