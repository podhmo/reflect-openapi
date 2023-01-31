package reflectopenapi

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/getkin/kin-openapi/openapi3"
	shape "github.com/podhmo/reflect-shape"
)

type Transformer struct {
	Resolver
	Selector Selector

	Extractor     Extractor
	TagNameOption TagNameOption

	cache    map[int]interface{}
	CacheHit int

	interceptFuncMap map[reflect.Type]func(*shape.Shape) *openapi3.Schema
	IsRequired       func(reflect.StructTag) bool
}

func (t *Transformer) RegisterInterception(rt reflect.Type, intercept func(*shape.Shape) *openapi3.Schema) {
	t.interceptFuncMap[rt] = intercept
}

func (t *Transformer) Builtin() *Transformer {
	// todo: handling required?
	{
		var z []byte
		t.RegisterInterception(reflect.ValueOf(z).Type(), func(s *shape.Shape) *openapi3.Schema {
			v := openapi3.NewStringSchema()
			v.Format = "binary"
			return v
		})
	}
	{
		var z time.Time
		t.RegisterInterception(reflect.ValueOf(z).Type(), func(s *shape.Shape) *openapi3.Schema {
			return openapi3.NewDateTimeSchema()
		})
	}
	return t
}

func (t *Transformer) Transform(s *shape.Shape) interface{} { // *Operation | *Schema | *Response
	id := s.Number
	if retval, ok := t.cache[id]; ok {
		t.CacheHit++
		return retval
	}

	// e.g. for time.Time as {"type": "string", "format": "date-time"}
	if intercept, ok := t.interceptFuncMap[s.Type]; ok {
		retval := intercept(s)
		t.cache[id] = retval
		return retval
	}

	switch s.Kind {
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
	case reflect.Struct:
		schema := openapi3.NewObjectSchema()
		t.cache[id] = schema
		ob := s.Struct()

		// add default value
		if rv := s.DefaultValue; rv.IsValid() && !rv.IsZero() && s.Name != "" {
			if !shape.IsZeroRecursive(s.Type, s.DefaultValue) {
				schema.Default = s.DefaultValue.Interface()
			}
		}

		// description
		if doc := ob.Doc(); doc != "" {
			schema.Description = doc
		}

		for _, f := range ob.Fields() {
			oaType, ok := f.Tag.Lookup(t.TagNameOption.ParamTypeTag)
			if ok {
				switch strings.ToLower(oaType) {
				case "cookie", "header", "path", "query":
					// log.debug: skip this is not body's field
					continue
				}
			}

			name, hasJsonTag := f.Tag.Lookup(t.TagNameOption.NameTag)
			hasOmitEmpty := false
			if !hasJsonTag {
				name = f.Name
			} else if strings.Contains(name, ",") {
				parts := strings.SplitN(name, ",", 2)
				name = parts[0]
				hasOmitEmpty = len(parts) > 1 && strings.Contains(parts[1], "omitempty")
			}
			if name == "-" {
				continue
			}

			if !hasJsonTag && !f.IsExported() {
				// skip if json tag is not found and unexported field
				continue
			}

			switch f.Shape.Kind {
			case reflect.Struct:
				subschema, ok := t.Transform(f.Shape).(*openapi3.Schema) // xxx
				if !ok {
					continue
				}

				if !f.Anonymous {
					schema.Properties[name] = t.ResolveSchema(subschema, f.Shape)
					if t.IsRequired(f.Tag) { // TODO: s.Metadata[i].Required
						schema.Required = append(schema.Required, name)
					}
				} else { // embedded
					for subname, subf := range subschema.Properties {
						schema.Properties[subname] = subf
					}
					if len(subschema.Required) > 0 {
						schema.Required = append(schema.Required, subschema.Required...)
					}
				}
			case reflect.Func, reflect.Chan:
				continue
			default:
				subschema, ok := t.Transform(f.Shape).(*openapi3.Schema) // xxx
				if !ok {
					continue
				}
				ref := t.ResolveSchema(subschema, f.Shape)
				schema.Properties[name] = ref
				if t.IsRequired(f.Tag) && !hasOmitEmpty { // TODO: s.Metadata[i].Required
					schema.Required = append(schema.Required, name)
				}

				// description
				if ref.Value != nil {
					doc := f.Doc
					if v, ok := f.Tag.Lookup(t.TagNameOption.DescriptionTag); ok {
						doc = v
					}
					if doc != "" {
						ref.Value.Description = doc
					}
				}

				// override: e.g. `openapi-override:"{'minimum': 0}"`
				if ref.Value != nil {
					if v, ok := f.Tag.Lookup(t.TagNameOption.OverrideTag); ok {
						if ref.Value.Extensions == nil {
							var overrideValues map[string]interface{}
							if err := json.Unmarshal([]byte(strings.ReplaceAll(strings.ReplaceAll(v, `\`, `\\`), "'", "\"")), &overrideValues); err != nil {
								log.Printf("[WARN]  openapi-override: unmarshal json is failed: %q", v)
							}
							ref.Value.Extensions = overrideValues
						}
					}
				}
			}
		}

		// too add-hoc?
		if len(schema.Properties) == 0 && ob.Fields().Len() > 0 {
			ok := true
			schema.AdditionalProperties.Has = &ok
			schema.Description = "<unclear definition>"
		}
		return schema
	case reflect.Func:
		// return *openapi.Operation
		op := openapi3.NewOperation()
		t.cache[id] = op
		{
			fullname := s.FullName()
			op.OperationID = fullname
		}
		op.Responses = openapi3.NewResponses()

		fn := s.Func()

		// description
		if doc := fn.Doc(); doc != "" {
			op.Description = doc
		}

		// parameters
		if inob := t.Selector.SelectInput(fn); inob != nil {
			schema := t.Transform(inob).(*openapi3.Schema) // xxx
			if len(schema.Properties) > 0 {
				// todo: required,content,description
				body := openapi3.NewRequestBody().
					WithJSONSchemaRef(t.ResolveSchema(schema, inob))
				op.RequestBody = t.ResolveRequestBody(body, inob)
			}

			// scan other
			if inob.Kind != reflect.Struct {
				log.Printf("[WARN]  only struct: but %s", inob.Kind)
				panic(inob)
			} else {
				params := openapi3.NewParameters()
				inob := inob.Struct()
				for _, f := range flattenFields(inob.Fields()) {
					paramType, ok := f.Tag.Lookup(t.TagNameOption.ParamTypeTag)
					if !ok {
						continue
					}

					name := f.Name
					if v, ok := f.Tag.Lookup(t.TagNameOption.NameTag); ok {
						name = v
					}

					// todo: required, type
					switch strings.ToLower(paramType) {
					case "json":
						continue
					case "path":
						if v, ok := f.Tag.Lookup("path"); ok {
							name = v
						}
						p := openapi3.NewPathParameter(name)
						schema := t.Transform(f.Shape).(*openapi3.Schema)
						p.Schema = t.ResolveSchema(schema, f.Shape)
						p.Description = f.Doc
						if v, ok := f.Tag.Lookup(t.TagNameOption.DescriptionTag); ok {
							p.Description = v
						}
						params = append(params, t.ResolveParameter(p, f.Shape))
					case "query":
						if v, ok := f.Tag.Lookup("query"); ok {
							name = v
						}
						p := openapi3.NewQueryParameter(name).
							WithRequired(t.IsRequired(f.Tag))
						schema := t.Transform(f.Shape).(*openapi3.Schema)
						p.Schema = t.ResolveSchema(schema, f.Shape)
						p.Description = f.Doc
						if v, ok := f.Tag.Lookup(t.TagNameOption.DescriptionTag); ok {
							p.Description = v
						}
						params = append(params, t.ResolveParameter(p, f.Shape))
					case "header":
						if v, ok := f.Tag.Lookup("header"); ok {
							name = v
						}
						p := openapi3.NewHeaderParameter(name).
							WithRequired(t.IsRequired(f.Tag))
						schema := t.Transform(f.Shape).(*openapi3.Schema)
						p.Schema = t.ResolveSchema(schema, f.Shape)
						p.Description = f.Doc
						if v, ok := f.Tag.Lookup(t.TagNameOption.DescriptionTag); ok {
							p.Description = v
						}
						params = append(params, t.ResolveParameter(p, f.Shape))
					case "cookie":
						if v, ok := f.Tag.Lookup("cookie"); ok {
							name = v
						}
						p := openapi3.NewCookieParameter(name).
							WithRequired(t.IsRequired(f.Tag))
						schema := t.Transform(f.Shape).(*openapi3.Schema)
						p.Schema = t.ResolveSchema(schema, f.Shape)
						p.Description = f.Doc
						if v, ok := f.Tag.Lookup(t.TagNameOption.DescriptionTag); ok {
							p.Description = v
						}
						params = append(params, t.ResolveParameter(p, f.Shape))
					default:
						log.Printf("[WARN]  invalid openapiTag: %q in %s.%s, suppored values are [path, query, header, cookie]", inob.Shape.Type, f.Name, f.Tag.Get(t.TagNameOption.ParamTypeTag))
					}
				}
				if len(params) > 0 {
					op.Parameters = params
				}
			}
		}

		// responses
		if outob := t.Selector.SelectOutput(fn); outob != nil {
			// todo: support (ob, error)
			schema := t.Transform(outob).(*openapi3.Schema) // xxx
			ref := t.ResolveSchema(schema, outob)
			doc := ""
			for _, p := range fn.Returns() {
				if p.Shape.Number == outob.Number && p.Doc != "" {
					doc = p.Doc
				}
			}
			op.Responses["200"] = t.ResolveResponse(
				openapi3.NewResponse().WithDescription(doc).WithJSONSchemaRef(ref),
				outob,
			)
		}
		return op
	case reflect.Slice, reflect.Array:
		schema := openapi3.NewArraySchema()
		t.cache[id] = schema

		var rob reflect.Value
		if s.DefaultValue.Len() > 0 {
			rob = s.DefaultValue.Index(0)
			if rob.Type().Kind() == reflect.Ptr && rob.IsNil() {
				rob = newInnerValue(s.Type)
			}
		} else {
			rob = newInnerValue(s.Type)
		}
		innerShape := t.Extractor.Extract(rob.Interface())

		inner, ok := t.Transform(innerShape).(*openapi3.Schema)
		if !ok {
			inner = openapi3.NewSchema()
		}
		schema.Items = t.ResolveSchema(inner, innerShape)
		return schema
	case reflect.Map:
		if s.Type.Key().Kind() != reflect.String {
			panic(fmt.Sprintf("not supported type %v, support only map[string, <V>]", s))
		}
		schema := openapi3.NewSchema()
		t.cache[id] = schema

		rob := newInnerValue(s.Type)
		innerShape := t.Extractor.Extract(rob.Interface())

		inner := t.Transform(innerShape).(*openapi3.Schema)
		schema.AdditionalProperties.Schema = t.ResolveSchema(inner, innerShape)
		return schema
	case reflect.Interface:
		iface := s.Interface()

		schema := openapi3.NewObjectSchema()
		if iface.Methods().Len() == 0 {
			schema.Description = "<Any type>"
		} else {
			log.Printf("`[INFO]  %v` is not supported, ignored.", s.Type)
			// schema.Description = fmt.Sprintf("`%v` is not supported, ignored", s.Type)
			return nil
		}

		ok := true
		schema.AdditionalProperties.Has = &ok

		return schema
	default:
		return notImplementedYet(s)
	}
}

func notImplementedYet(s *shape.Shape) interface{} {
	if ok, _ := strconv.ParseBool(os.Getenv("FORCE")); ok {
		log.Printf("[INFO]  not implemented yet for %+v", s)
		return nil
	}
	panic(fmt.Sprintf("not implemented yet for %v\nIf you want to run forcibly, execute with FORCE=1", s))
}

func flattenFields(fields shape.FieldList) shape.FieldList {
	r := make([]*shape.Field, 0, fields.Len())
	for _, f := range fields {
		if !f.Anonymous {
			r = append(r, f)
			continue
		}
		if f.Shape.Kind == reflect.Struct {
			r = append(r, flattenFields(f.Shape.Struct().Fields())...)
		}
	}
	return r
}

// return not zero inner value from map or slice
func newInnerValue(rt reflect.Type) reflect.Value {
	innerType := rt.Elem()
	lv := 0
	if innerType.Kind() == reflect.Ptr {
		lv++
		innerType = innerType.Elem()
	}
	rob := reflect.New(innerType).Elem()
	for i := 0; i < lv; i++ {
		rob = rob.Addr()
	}
	return rob
}
