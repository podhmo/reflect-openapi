package reflectopenapi

import (
	"reflect"
	"strconv"

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

	extractor *shape.Extractor
}

func isRequiredDefault(tag reflect.StructTag) bool {
	s, ok := tag.Lookup("required")
	if !ok {
		return false
	}
	v, _ := strconv.ParseBool(s)
	return v
}

func NewVisitor(resolver Resolver, selector Selector, extractor *shape.Extractor) *Visitor {
	return &Visitor{
		Transformer: (&Transformer{
			cache:            map[int]interface{}{},
			interceptFuncMap: map[reflect.Type]func(*shape.Shape) *openapi3.Schema{},
			Resolver:         resolver,
			IsRequired:       isRequiredDefault,
			Selector:         selector,
			extractor:        extractor,
		}).Builtin(),
		Schemas:    map[int]*openapi3.Schema{},
		Operations: map[int]*openapi3.Operation{},
		extractor:  extractor,
	}
}

func (v *Visitor) VisitType(ob interface{}, modifiers ...func(*openapi3.Schema)) *openapi3.SchemaRef {
	in := v.extractor.Extract(ob)
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
				"x-new-type": in.FullName(),
			}
		}
		v.Transformer.cache[id] = out
	}
	return v.ResolveSchema(out, in)
}
func (v *Visitor) VisitFunc(ob interface{}, modifiers ...func(*openapi3.Operation)) *openapi3.Operation {
	in := v.extractor.Extract(ob)
	out := v.Transform(in).(*openapi3.Operation)
	// FIXME: comment lookup
	// if v.CommentLookup != nil {
	// 	description, err := v.CommentLookup.LookupCommentTextFromFunc(ob)
	// 	if err != nil {
	// 		log.Printf("comment lookup failed, %v", ob)
	// 	} else {
	// 		parts := strings.Split(out.OperationID, ".")
	// 		description := strings.TrimSpace(strings.TrimPrefix(description, parts[len(parts)-1]))
	// 		out.Description = description
	// 		out.Summary = strings.SplitN(description, "\n", 2)[0]
	// 	}
	// }

	for _, m := range modifiers {
		m(out)
	}

	v.Operations[in.Number] = out
	return out
}
