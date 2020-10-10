package reflectopenapi

import (
	"fmt"
	"log"
	"reflect"
	"strings"
	"time"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/podhmo/reflect-openapi/pkg/comment"
	"github.com/podhmo/reflect-openapi/pkg/shape"
)

// TODO: validation for schema
// TODO: support function input
// TODO: extra information
// TODO: integration with net/http.Handler
// TODO: integration with fasthttp.Handler
// TODO: json tag inline,omitempty support
// TODO: schema required, unrequired support
// TODO: schema nullable support (?)

// not visitor pattern
type Visitor struct {
	*Transformer
	Doc        *openapi3.Swagger
	Schemas    map[reflect.Type]*openapi3.Schema
	Operations map[reflect.Type]*openapi3.Operation

	CommentLookup *comment.Lookup
}

func NewVisitor(resolver Resolver) *Visitor {
	return &Visitor{
		Transformer: (&Transformer{
			cache:            map[reflect.Type]interface{}{},
			interceptFuncMap: map[reflect.Type]func(shape.Shape) *openapi3.Schema{},
			Resolver:         resolver,
			IsRequired:       func(tag reflect.StructTag) bool { return false },
		}).Builtin(),
		Schemas:    map[reflect.Type]*openapi3.Schema{},
		Operations: map[reflect.Type]*openapi3.Operation{},
	}
}

func (v *Visitor) VisitType(ob interface{}, modifiers ...func(*openapi3.Schema)) *openapi3.SchemaRef {
	in := shape.Extract(ob)
	out := v.Transform(in).(*openapi3.Schema)
	for _, m := range modifiers {
		m(out)
	}
	v.Schemas[in.GetReflectType()] = out
	return v.ResolveSchema(out, in)
}
func (v *Visitor) VisitFunc(ob interface{}, modifiers ...func(*openapi3.Operation)) *openapi3.Operation {
	in := shape.Extract(ob)
	out := v.Transform(in).(*openapi3.Operation)
	for _, m := range modifiers {
		m(out)
	}
	if v.CommentLookup != nil {
		description, err := v.CommentLookup.LookupCommentTextFromFunc(ob)
		if err != nil {
			log.Printf("comment lookup failed, %v", ob)
		} else {
			parts := strings.Split(out.OperationID, ".")
			out.Description = strings.TrimSpace(strings.TrimPrefix(description, parts[len(parts)-1]))
		}
	}
	v.Operations[in.GetReflectType()] = out
	return out
}

type Transformer struct {
	Resolver

	cache    map[reflect.Type]interface{}
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
	rt := s.GetReflectType()
	if retval, ok := t.cache[rt]; ok {
		t.CacheHit++
		return retval
	}

	// e.g. for time.Time as {"type": "string", "format": "date-time"}
	if intercept, ok := t.interceptFuncMap[rt]; ok {
		retval := intercept(s)
		t.cache[rt] = retval
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
			notImplementedYet(s)
		}
	case shape.Struct:
		schema := openapi3.NewObjectSchema()
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

			switch v.GetReflectKind() {
			case reflect.Struct:
				f := t.Transform(v).(*openapi3.Schema) // xxx

				if !s.Metadata[i].Anonymous {
					schema.Properties[name] = t.ResolveSchema(f, v)
					if t.IsRequired(s.Tags[i]) {
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
				f := t.Transform(v).(*openapi3.Schema) // xxx
				schema.Properties[name] = t.ResolveSchema(f, v)
				if t.IsRequired(s.Tags[i]) {
					schema.Required = append(schema.Required, name)
				}
			}
		}
		t.cache[rt] = schema
		return schema
	case shape.Function:
		// return *openapi.Operation
		// as interactor (TODO: meta tag? for specific usecase)

		op := openapi3.NewOperation()
		op.OperationID = s.GetFullName()
		op.Responses = openapi3.NewResponses()

		// parameters
		if len(s.Params.Values) > 0 {
			// todo: support (ctx, ob)

			// scan body
			inob := s.Params.Values[0]
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
						schema := t.Transform(v).(*openapi3.Schema)
						p := openapi3.NewPathParameter(inob.FieldName(i)).
							WithSchema(schema)
						params = append(params, t.ResolveParameter(p, v))
					case "query":
						schema := t.Transform(v).(*openapi3.Schema)
						p := openapi3.NewQueryParameter(inob.FieldName(i)).
							WithSchema(schema).
							WithRequired(t.IsRequired(inob.Tags[i]))
						params = append(params, t.ResolveParameter(p, v))
					case "header":
						schema := t.Transform(v).(*openapi3.Schema)
						p := openapi3.NewHeaderParameter(inob.FieldName(i)).
							WithSchema(schema).
							WithRequired(t.IsRequired(inob.Tags[i]))
						params = append(params, t.ResolveParameter(p, v))
					case "cookie":
						schema := t.Transform(v).(*openapi3.Schema)
						p := openapi3.NewCookieParameter(inob.FieldName(i)).
							WithSchema(schema).
							WithRequired(t.IsRequired(inob.Tags[i]))
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
		if len(s.Returns.Values) > 0 {
			// todo: support (ob, error)
			outob := s.Returns.Values[0]
			schema := t.Transform(outob).(*openapi3.Schema) // xxx
			op.Responses["200"] = t.ResolveResponse(
				openapi3.NewResponse().WithDescription("").WithJSONSchemaRef(
					t.ResolveSchema(schema, outob),
				),
				s.Returns.Values[0],
			)
		}
		t.cache[rt] = op
		return op
	case shape.Container:
		// container is map,slice,array
		switch s.GetReflectKind() {
		case reflect.Slice:
			schema := openapi3.NewArraySchema()
			inner := t.Transform(s.Args[0]).(*openapi3.Schema)
			schema.Items = t.ResolveSchema(inner, s.Args[0])
			t.cache[rt] = schema
			return schema
		default:
			notImplementedYet(s)
		}
		notImplementedYet(s)
	default:
		notImplementedYet(s)
	}
	panic("never")
}

func notImplementedYet(ob interface{}) {
	panic(ob)
}
