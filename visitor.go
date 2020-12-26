package reflectopenapi

import (
	"fmt"
	"log"
	"os"
	"reflect"
	"strconv"
	"strings"
	"time"

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

	Doc        *openapi3.Swagger
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

type Transformer struct {
	Resolver
	Selector Selector

	cache    map[shape.Identity]interface{}
	CacheHit int

	interceptFuncMap map[reflect.Type]func(shape.Shape) *openapi3.Schema
	IsRequired       func(reflect.StructTag) bool
}

func (t *Transformer) RegisterInterception(rt reflect.Type, intercept func(shape.Shape) *openapi3.Schema) {
	t.interceptFuncMap[rt] = intercept
}

func (t *Transformer) Builtin() *Transformer {
	// todo: handling required?
	{
		var z []byte
		t.RegisterInterception(reflect.ValueOf(z).Type(), func(s shape.Shape) *openapi3.Schema {
			v := openapi3.NewStringSchema()
			v.Format = "binary"
			return v
		})
	}
	{
		var z time.Time
		t.RegisterInterception(reflect.ValueOf(z).Type(), func(s shape.Shape) *openapi3.Schema {
			return openapi3.NewDateTimeSchema()
		})
	}
	return t
}

func (t *Transformer) Transform(s shape.Shape) interface{} { // *Operation | *Schema | *Response
	id := s.GetIdentity()
	if retval, ok := t.cache[id]; ok {
		t.CacheHit++
		return retval
	}

	// e.g. for time.Time as {"type": "string", "format": "date-time"}
	if intercept, ok := t.interceptFuncMap[s.GetReflectType()]; ok {
		retval := intercept(s)
		t.cache[id] = retval
		return retval
	}

	switch s := s.(type) {
	case shape.Primitive:
		switch s.GetReflectKind() {
		case reflect.Bool:
			return openapi3.NewBoolSchema()
		case reflect.String:
			return openapi3.NewStringSchema()
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			// Todo: use NewInt64Schema?
			return openapi3.NewIntegerSchema()
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64: // Uintptr
			return openapi3.NewIntegerSchema()
		case reflect.Float32, reflect.Float64:
			return openapi3.NewFloat64Schema()
		default:
			return notImplementedYet(s)
		}
	case shape.Struct:
		schema := openapi3.NewObjectSchema()
		t.cache[id] = schema
		for i, v := range s.Fields.Values {
			oaType, ok := s.Tags[i].Lookup("openapi")
			if ok {
				switch strings.ToLower(oaType) {
				case "cookie", "header", "path", "query":
					// log.debug: skip this is not body's field
					continue
				}
			}

			name := s.FieldName(i)
			if name == "-" {
				continue
			}

			// skip if json tag is not found and unexported field
			if fname := s.Fields.Keys[i]; s.Metadata[i].FieldName == "" && fname[0] == strings.ToLower(fname)[0] {
				continue
			}

			switch v.GetReflectKind() {
			case reflect.Struct:
				f, ok := t.Transform(v).(*openapi3.Schema) // xxx
				if !ok {
					continue
				}

				if !s.Metadata[i].Anonymous {
					schema.Properties[name] = t.ResolveSchema(f, v)
					if s.Metadata[i].Required || t.IsRequired(s.Tags[i]) {
						schema.Required = append(schema.Required, name)
					}
				} else { // embedded
					for subname, subf := range f.Properties {
						schema.Properties[subname] = subf
					}
					if len(f.Required) > 0 {
						schema.Required = append(schema.Required, f.Required...)
					}
				}
			case reflect.Func, reflect.Chan:
				continue
			default:
				f, ok := t.Transform(v).(*openapi3.Schema) // xxx
				if !ok {
					continue
				}
				schema.Properties[name] = t.ResolveSchema(f, v)
				if s.Metadata[i].Required || t.IsRequired(s.Tags[i]) {
					schema.Required = append(schema.Required, name)
				}
			}
		}
		return schema
	case shape.Function:
		// return *openapi.Operation
		// as interactor (TODO: meta tag? for specific usecase)

		op := openapi3.NewOperation()
		t.cache[id] = op
		{
			fullname := s.GetFullName()
			fullname = strings.ReplaceAll(strings.ReplaceAll(strings.ReplaceAll(fullname, "(", ""), ")", ""), "*", "")
			op.OperationID = strings.TrimSuffix(strings.TrimPrefix(fullname, "."), "-fm")
		}
		op.Responses = openapi3.NewResponses()

		// parameters
		if inob := t.Selector.SelectInput(s); inob != nil {
			schema := t.Transform(inob).(*openapi3.Schema) // xxx
			if len(schema.Properties) > 0 {
				// todo: required,content,description
				body := openapi3.NewRequestBody().
					WithJSONSchemaRef(t.ResolveSchema(schema, inob))
				op.RequestBody = t.ResolveRequestBody(body, inob)
			}

			// scan other
			switch inob := inob.(type) {
			case shape.Struct:
				params := openapi3.NewParameters()
				for i, v := range inob.Fields.Values {
					paramType, ok := inob.Tags[i].Lookup("openapi")
					if !ok {
						continue
					}

					// todo: required, type
					switch strings.ToLower(paramType) {
					case "json":
						continue
					case "path":
						p := openapi3.NewPathParameter(inob.FieldName(i))
						schema := t.Transform(v).(*openapi3.Schema)
						p.Schema = t.ResolveSchema(schema, v)
						params = append(params, t.ResolveParameter(p, v))
					case "query":
						p := openapi3.NewQueryParameter(inob.FieldName(i)).
							WithRequired(t.IsRequired(inob.Tags[i]))
						schema := t.Transform(v).(*openapi3.Schema)
						p.Schema = t.ResolveSchema(schema, v)
						params = append(params, t.ResolveParameter(p, v))
					case "header":
						p := openapi3.NewHeaderParameter(inob.FieldName(i)).
							WithRequired(t.IsRequired(inob.Tags[i]))
						schema := t.Transform(v).(*openapi3.Schema)
						p.Schema = t.ResolveSchema(schema, v)
						params = append(params, t.ResolveParameter(p, v))
					case "cookie":
						p := openapi3.NewCookieParameter(inob.FieldName(i)).
							WithRequired(t.IsRequired(inob.Tags[i]))
						schema := t.Transform(v).(*openapi3.Schema)
						p.Schema = t.ResolveSchema(schema, v)
						params = append(params, t.ResolveParameter(p, v))
					default:
						panic(paramType)
					}
				}
				if len(params) > 0 {
					op.Parameters = params
				}
			default:
				fmt.Println("only struct")
				panic(inob)
			}
		}

		// responses
		if outob := t.Selector.SelectOutput(s); outob != nil {
			// todo: support (ob, error)
			schema := t.Transform(outob).(*openapi3.Schema) // xxx
			op.Responses["200"] = t.ResolveResponse(
				openapi3.NewResponse().WithDescription("").WithJSONSchemaRef(
					t.ResolveSchema(schema, outob),
				),
				outob,
			)
		}
		return op
	case shape.Container:
		// container is map,slice,array
		switch s.GetReflectKind() {
		case reflect.Slice, reflect.Array:
			schema := openapi3.NewArraySchema()
			t.cache[id] = schema
			inner, ok := t.Transform(s.Args[0]).(*openapi3.Schema)
			if !ok {
				inner = openapi3.NewSchema()
			}
			schema.Items = t.ResolveSchema(inner, s.Args[0])
			return schema
		case reflect.Map:
			if s.Args[0].GetReflectKind() != reflect.String {
				panic(fmt.Sprintf("not supported type %v, support only map[string, <V>]", s))
			}
			schema := openapi3.NewSchema()
			t.cache[id] = schema
			inner := t.Transform(s.Args[1]).(*openapi3.Schema)
			schema.AdditionalProperties = t.ResolveSchema(inner, s.Args[1])
			return schema
		default:
			return notImplementedYet(s)
		}
	case shape.Interface:
		log.Printf("interface is not supported, ignored. %v", s)
		return nil
	default:
		return notImplementedYet(s)
	}
}

func notImplementedYet(s shape.Shape) interface{} {
	if ok, _ := strconv.ParseBool(os.Getenv("FORCE")); ok {
		log.Printf("not implemented yet for %+v", s)
		return nil
	}
	panic(fmt.Sprintf("not implemented yet for %v\nIf you want to run forcibly, execute with FORCE=1", s))
}
