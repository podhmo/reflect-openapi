package reflectopenapi

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"reflect"
	"sort"
	"strconv"
	"strings"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/podhmo/reflect-openapi/info"
	shape "github.com/podhmo/reflect-shape"
)

func NewDoc() (*openapi3.T, error) {
	skeleton := []byte(`{
  "openapi": "3.0.0",
  "info": {
    "title": "Sample API",
    "description": "",
    "version": "0.0.0"
  },
  "servers": [
    {
      "url": "http://localhost:8888",
      "description": "local development server"
    },
  ],
  "paths": {}
}`)
	return NewDocFromSkeleton(skeleton)
}

// TODO: add api function
func NewDocFromSkeleton(spec []byte) (*openapi3.T, error) {
	l := openapi3.NewLoader()
	return l.LoadFromData(spec)
}

type TagNameOption struct {
	NameTag        string
	RequiredTag    string
	ParamTypeTag   string
	OverrideTag    string
	DescriptionTag string

	XNewTypeTag string
}

func DefaultTagNameOption() *TagNameOption {
	return &TagNameOption{
		NameTag:        "json",
		RequiredTag:    "required",
		ParamTypeTag:   "in",
		DescriptionTag: "description",
		OverrideTag:    "openapi-override",
		XNewTypeTag:    "x-go-type",
	}
}

type Config struct {
	*TagNameOption
	Info *info.Info // go/types.Info like object (tracking metadata)

	Doc    *openapi3.T
	Loaded bool // if true, skip registerType() and registerFunc() actions

	Resolver  Resolver
	Selector  Selector
	Extractor Extractor

	StrictSchema        bool // if true, use `{additionalProperties: false}` as default
	SkipValidation      bool // if true, skip validation for api doc definition
	SkipExtractComments bool // if true, skip extracting comments as a description

	EnableAutoTag bool // if true, adding package name as tag

	DisableInputRef  bool
	DisableOutputRef bool

	DefaultError            interface{}
	DefaultErrorExample     interface{}
	IsRequiredCheckFunction func(reflect.StructTag) bool // handling required, default is always false
}

func (c *Config) DefaultResolver() Resolver {
	if c.TagNameOption == nil {
		c.TagNameOption = DefaultTagNameOption()
	}
	if c.Resolver != nil {
		return c.Resolver
	}
	resolver := &UseRefResolver{
		NameStore:        NewNameStore(),
		DisableInputRef:  c.DisableInputRef,
		DisableOutputRef: c.DisableOutputRef,
	}
	if c.Info != nil {
		resolver.NameStore.info = c.Info
	}
	if c.StrictSchema {
		ng := false
		resolver.AdditionalPropertiesAllowed = &ng
	}
	c.Resolver = resolver
	return c.Resolver
}

func (c *Config) DefaultExtractor() Extractor {
	if c.Extractor != nil {
		return c.Extractor
	}
	cfg := &shape.Config{
		FillArgNames:    true,
		FillReturnNames: true,
		SkipComments:    c.SkipExtractComments,
	}
	c.Extractor = cfg
	return c.Extractor
}

func (c *Config) DefaultSelector() Selector {
	if c.Selector != nil {
		return c.Selector
	}
	c.Selector = &DefaultSelector{}
	return c.Selector
}

func (c *Config) NewManager() (*Manager, func(ctx context.Context) error, error) {
	if c.Doc == nil {
		doc, err := NewDoc()
		if err != nil {
			return nil, nil, err
		}
		c.Doc = doc
	}

	v := NewVisitor(
		*c.TagNameOption,
		c.DefaultResolver(),
		c.DefaultSelector(),
		c.DefaultExtractor(),
	)
	v.EnableAutoTag = c.EnableAutoTag
	v.info = c.Info
	if c.IsRequiredCheckFunction != nil {
		v.Transformer.IsRequired = c.IsRequiredCheckFunction
	}

	m := &Manager{
		Doc:      c.Doc,
		Resolver: c.Resolver,
		Visitor:  v,
	}

	return m, func(ctx context.Context) error {
		doValidation := func() error {
			if !c.SkipValidation {
				// preventing the error like `invalid components: schema <name>: invalid default: unhandled value of type <Type>``
				b, err := json.Marshal(m.Doc)
				if err != nil {
					return fmt.Errorf("marshal doc before validation: %w", err)
				}
				doc, err := openapi3.NewLoader().LoadFromData(b)
				if err != nil {
					return fmt.Errorf("load doc before validation: %w", err)
				}
				// m.Doc = doc // need?
				if err := doc.Validate(ctx); err != nil {
					return err
				}
			}
			return nil
		}
		if c.Loaded {
			log.Printf("[INFO]  openapi-doc building process is skipped, because the doc is loaded from a file.")
			return doValidation()
		}

		// perform execution
		sort.Slice(m.Actions, func(i, j int) bool { return m.Actions[i].Phase < m.Actions[j].Phase })
		for _, ac := range m.Actions {
			ac.Action()
		}

		if c.DefaultError != nil {
			errSchema := v.VisitType(v.Transformer.Extractor.Extract(c.DefaultError))
			responseRef := &openapi3.ResponseRef{
				Value: openapi3.NewResponse().
					WithDescription("default error").
					WithJSONSchemaRef(errSchema),
			}
			if c.DefaultErrorExample != nil {
				content := responseRef.Value.Content["application/json"]
				if content.Examples == nil {
					content.Examples = openapi3.Examples{}
				}
				content.Examples["default"] = &openapi3.ExampleRef{Value: openapi3.NewExample(c.DefaultErrorExample)}

				if c.Info != nil {
					c.Info.SchemaValue[errSchema].Example = c.DefaultErrorExample
				}
			}

			for _, op := range v.Operations {
				if val, ok := op.Responses["default"]; !ok || val.Value == nil || val.Value.Description == nil || *val.Value.Description != "" {
					continue
				}
				op.Responses["default"] = responseRef
			}
		}

		if b, ok := c.Resolver.(Binder); ok {
			b.BindSchemas(m.Doc)
		}

		return doValidation()
	}, nil
}

func (c *Config) BuildDoc(ctx context.Context, use func(m *Manager)) (*openapi3.T, error) {
	m, commit, err := c.NewManager()
	if err != nil {
		return m.Doc, err
	}
	use(m)
	if err := commit(ctx); err != nil {
		return m.Doc, err
	}
	return m.Doc, nil
}

func (c *Config) EmitDoc(use func(m *Manager)) {
	ctx := context.Background()
	doc, err := c.BuildDoc(ctx, use)
	if err != nil {
		panic(err)
	}
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	if err := enc.Encode(doc); err != nil {
		panic(err)
	}
}

type Manager struct {
	Visitor  *Visitor // TODO: changes to unexported field
	Resolver Resolver
	Actions  []*registerAction

	Doc *openapi3.T
}

const (
	phase1Action = 0
	phase2Action = 1
)

type registerAction struct {
	Phase   int // lowerst is first
	Manager *Manager
	Action  func()
}

type RegisterTypeAction struct {
	*registerAction
	before func(*shape.Shape)
	after  func(*openapi3.Schema)
}

func (a *RegisterTypeAction) After(f func(*openapi3.Schema)) *RegisterTypeAction {
	if a.after == nil {
		a.after = f
		return a
	}

	prevFn := a.after
	a.after = func(ref *openapi3.Schema) {
		prevFn(ref)
		f(ref)
	}
	return a
}
func (a *RegisterTypeAction) Before(f func(*shape.Shape)) *RegisterTypeAction {
	if a.before == nil {
		a.before = f
		return a
	}

	prevFn := a.before
	a.before = func(s *shape.Shape) {
		prevFn(s)
		f(s)
	}
	return a
}
func (a *RegisterTypeAction) Description(description string) *RegisterTypeAction {
	return a.After(func(s *openapi3.Schema) {
		s.Description = description
	})
}
func (a *RegisterTypeAction) Enum(values ...interface{}) *RegisterTypeAction {
	return a.After(func(s *openapi3.Schema) {
		s.Enum = values
	})
}
func (a *RegisterTypeAction) Default(value interface{}) *RegisterTypeAction {
	return a.After(func(s *openapi3.Schema) {
		s.Default = value
	})
}
func (a *RegisterTypeAction) Example(value interface{}) *RegisterTypeAction {
	return a.After(func(s *openapi3.Schema) {
		s.Example = value
	})
}

func (m *Manager) RegisterType(ob interface{}, modifiers ...func(*openapi3.SchemaRef)) *RegisterTypeAction {
	var ac *RegisterTypeAction
	ac = &RegisterTypeAction{
		registerAction: &registerAction{
			Manager: m,
			Phase:   phase1Action,
			Action: func() {
				in := m.Visitor.Transformer.Extractor.Extract(ob)
				if ac.before != nil {
					ac.before(in)
				}
				var options []func(*openapi3.Schema)
				if ac.after != nil {
					options = append(options, ac.after)
				}
				ref := m.Visitor.VisitType(in, options...)
				for _, m := range modifiers {
					m(ref)
				}
			},
		},
	}
	m.Actions = append(m.Actions, ac.registerAction)
	return ac
}

type RegisterFuncAction struct {
	*registerAction
	before func(*shape.Func)
	after  func(*openapi3.Operation)
	fn     interface{}
}

func (a *RegisterFuncAction) After(f func(*openapi3.Operation)) *RegisterFuncAction {
	if a.after == nil {
		a.after = f
		return a
	}
	prevFn := a.after
	a.after = func(op *openapi3.Operation) {
		prevFn(op)
		f(op)
	}
	return a
}
func (a *RegisterFuncAction) Before(f func(*shape.Func)) *RegisterFuncAction {
	if a.before == nil {
		a.before = f
		return a
	}

	prevFn := a.before
	a.before = func(s *shape.Func) {
		prevFn(s)
		f(s)
	}
	return a
}
func (a *RegisterFuncAction) Description(description string) *RegisterFuncAction {
	return a.After(func(op *openapi3.Operation) {
		op.Description = strings.TrimSpace(description)
	})
}
func (a *RegisterFuncAction) Status(code int) *RegisterFuncAction {
	return a.After(func(op *openapi3.Operation) {
		def, ok := op.Responses["200"]
		if ok {
			delete(op.Responses, "200")
			op.Responses[strconv.Itoa(code)] = def
		}
	})
}
func (a *RegisterFuncAction) Error(value interface{}, description string) *RegisterFuncAction {
	var errSchema *openapi3.SchemaRef
	a.Manager.RegisterType(value, func(ref *openapi3.SchemaRef) {
		errSchema = ref
	})

	return a.After(func(op *openapi3.Operation) {
		op.AddResponse(0, openapi3.NewResponse().WithDescription(description).WithContent(openapi3.NewContentWithJSONSchemaRef(errSchema)))
	})
}
func (a *RegisterFuncAction) DefaultInput(value interface{}) *RegisterFuncAction {
	if value == nil {
		return a
	}

	return a.Before(func(fn *shape.Func) {
		t := a.Manager.Visitor.Transformer
		inob, _ := t.Selector.SelectInput(fn)

		rt := reflect.TypeOf(value)
		rv := reflect.ValueOf(value)
		for rt.Kind() == reflect.Pointer {
			rt = rt.Elem()
			rv = rv.Elem()
		}

		// FIXME: MergeParamsSInputSelector is not supported
		if inob.Type != rt {
			if FORCE {
				log.Printf("[WARN]  DefaultInput: expected type is %v but got %v, ignored...", inob.Type, rt)
				return
			}
			panic(fmt.Sprintf("DefaultInput: expected type is %v but got %v", inob.Type, rt))
		}
		t.defaultValues[inob.Number] = rv
	})
}
func (a *RegisterFuncAction) Example(code int, mime string, title, description string, value interface{}) *RegisterFuncAction {
	// does not use Example, Examples only.
	return a.After(func(op *openapi3.Operation) {
		ref := op.Responses[strconv.Itoa(code)]
		if ref == nil {
			schemaRef := a.Manager.Visitor.VisitType(a.Manager.Visitor.Extractor.Extract(value))
			ref = &openapi3.ResponseRef{Value: openapi3.NewResponse().WithJSONSchemaRef(schemaRef).WithDescription("-")}
			op.Responses[strconv.Itoa(code)] = ref
		}

		title := strings.TrimSpace(title)
		description := strings.TrimSpace(description)

		if ref.Value != nil {
			if ref.Value.Content == nil {
				ref.Value.Content = openapi3.NewContentWithJSONSchema(openapi3.NewSchema())
			}
			mediatype := ref.Value.Content.Get(mime)
			if mediatype == nil {
				mediatype = openapi3.NewMediaType()
				ref.Value.Content[mime] = mediatype
			}

			if mediatype.Examples == nil {
				mediatype.Examples = openapi3.Examples{}
			}
			if mediatype.Example != nil {
				mediatype.Examples["default"] = &openapi3.ExampleRef{Value: &openapi3.Example{
					Value: mediatype.Example,
				}}
				mediatype.Example = nil
			}
			if title == "" || title == "default" {
				title = "default"
				if n := len(mediatype.Examples); n > 0 {
					title += strconv.Itoa(n)
				}
			}
			mediatype.Examples[title] = &openapi3.ExampleRef{Value: &openapi3.Example{
				Value:       value,
				Description: description,
				Summary:     description,
			}}
		}
	})
}

// func (a *RegisterFuncAction) AnotherError(code int, typ interface{}) *RegisterFuncAction {
// 	return a.After(func(op *openapi3.Operation) {
// 		op.Responses[strconv.Itoa(code)] = // TODO: implement
// 	})
// }

func (m *Manager) RegisterFunc(fn interface{}, modifiers ...func(*openapi3.Operation)) *RegisterFuncAction {
	var ac *RegisterFuncAction
	ac = &RegisterFuncAction{
		fn: fn,
		registerAction: &registerAction{
			Phase:   phase2Action,
			Manager: m,
			Action: func() {
				in := m.Visitor.Transformer.Extractor.Extract(fn)
				if in.Kind != reflect.Func {
					panic(fmt.Sprintf("not function: %v", in))
				}
				sfn := in.Func()
				if ac.before != nil {
					ac.before(sfn)
				}
				if ac.after != nil {
					modifiers = append(modifiers, ac.after)
				}
				m.Visitor.VisitFunc(in, modifiers...)
			},
		},
	}
	m.Actions = append(m.Actions, ac.registerAction)
	return ac
}

func (m *Manager) RegisterFuncText(fn interface{}, contentType string, modifiers ...func(*openapi3.Operation)) *RegisterFuncAction {
	return m.RegisterFunc(fn, append([]func(*openapi3.Operation){
		func(op *openapi3.Operation) {
			res := op.Responses.Get(200).Value
			ref := res.Content.Get("application/json").Schema
			ref.Value = openapi3.NewStringSchema()
			res.Content = openapi3.NewContentWithSchemaRef(ref, []string{contentType})
		},
	}, modifiers...)...)
}

func (m *Manager) RegisterInterception(rt reflect.Type, intercept func(*shape.Shape) *openapi3.Schema) {
	m.Visitor.Transformer.RegisterInterception(rt, intercept)
}

var FORCE bool

func init() {
	if ok, _ := strconv.ParseBool(os.Getenv("FORCE")); ok {
		FORCE = ok
	}
}
