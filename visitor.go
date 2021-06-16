package reflectopenapi

import (
	"log"
	"reflect"
	"strconv"
	"strings"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/podhmo/reflect-openapi/pkg/comment"
	"github.com/podhmo/reflect-openapi/pkg/shape"
)

// TODO: extra information
// TODO: json tag inline,omitempty support
// TODO: schema nullable support (?)

// not visitor pattern
type Visitor struct {
	*Transformer
	CommentLookup *comment.Lookup

	Doc        *openapi3.T
	Schemas    map[shape.Identity]*openapi3.Schema
	Operations map[shape.Identity]*openapi3.Operation

	extractor Extractor
}

func isRequiredDefault(tag reflect.StructTag) bool {
	s, ok := tag.Lookup("required")
	if !ok {
		return false
	}
	v, _ := strconv.ParseBool(s)
	return v
}

func NewVisitor(resolver Resolver, selector Selector, extractor Extractor) *Visitor {
	return &Visitor{
		Transformer: (&Transformer{
			cache:            map[shape.Identity]interface{}{},
			interceptFuncMap: map[reflect.Type]func(shape.Shape) *openapi3.Schema{},
			Resolver:         resolver,
			IsRequired:       isRequiredDefault,
			Selector:         selector,
		}).Builtin(),
		Schemas:    map[shape.Identity]*openapi3.Schema{},
		Operations: map[shape.Identity]*openapi3.Operation{},
		extractor:  extractor,
	}
}

func (v *Visitor) VisitType(ob interface{}, modifiers ...func(*openapi3.Schema)) *openapi3.SchemaRef {
	in := v.extractor.Extract(ob)
	out := v.Transform(in).(*openapi3.Schema)
	out.Title = in.GetName()
	for _, m := range modifiers {
		m(out)
	}

	id := in.GetIdentity()
	v.Schemas[id] = out
	if len(modifiers) > 0 {
		if out.Extensions == nil {
			out.Extensions = map[string]interface{}{
				"x-new-type": in.GetFullName(),
			}
		}
		v.Transformer.cache[id] = out
	}
	return v.ResolveSchema(out, in)
}
func (v *Visitor) VisitFunc(ob interface{}, modifiers ...func(*openapi3.Operation)) *openapi3.Operation {
	in := v.extractor.Extract(ob)
	out := v.Transform(in).(*openapi3.Operation)
	if v.CommentLookup != nil {
		description, err := v.CommentLookup.LookupCommentTextFromFunc(ob)
		if err != nil {
			log.Printf("comment lookup failed, %v", ob)
		} else {
			parts := strings.Split(out.OperationID, ".")
			description := strings.TrimSpace(strings.TrimPrefix(description, parts[len(parts)-1]))
			out.Description = description
			out.Summary = strings.SplitN(description, "\n", 2)[0]
		}
	}

	for _, m := range modifiers {
		m(out)
	}

	v.Operations[in.GetIdentity()] = out
	return out
}
