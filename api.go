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
	shape "github.com/podhmo/reflect-shape"
)

func NewDoc() (*openapi3.T, error) {
	skeleton := []byte(`{
  "openapi": "3.0.0",
  "info": {
    "title": "Sample API",
    "description": "-",
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
	ParamTypeTag   string
	OverrideTag    string
	DescriptionTag string

	XNewTypeTag string
}

func DefaultTagNameOption() *TagNameOption {
	return &TagNameOption{
		NameTag:        "json",
		ParamTypeTag:   "in",
		DescriptionTag: "description",
		OverrideTag:    "openapi-override",
		XNewTypeTag:    "x-go-type",
	}
}

type Config struct {
	*TagNameOption

	Doc    *openapi3.T
	Loaded bool // if true, skip registerType() and registerFunc() actions

	Resolver  Resolver
	Selector  Selector
	Extractor Extractor

	StrictSchema        bool // if true, use `{additionalProperties: false}` as default
	SkipValidation      bool // if true, skip validation for api doc definition
	SkipExtractComments bool // if true, skip extracting comments as a description

	DefaultError            interface{}
	IsRequiredCheckFunction func(reflect.StructTag) bool // handling required, default is always false
}

func (c *Config) DefaultResolver() Resolver {
	if c.TagNameOption == nil {
		c.TagNameOption = DefaultTagNameOption()
	}
	if c.Resolver != nil {
		return c.Resolver
	}
	resolver := &UseRefResolver{NameStore: NewNameStore()}
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
			log.Printf("[INFO]  Skips execution because openapi doc is loaded from file")
			return doValidation()
		}

		// perform execution
		sort.Slice(m.Actions, func(i, j int) bool { return m.Actions[i].Phase < m.Actions[j].Phase })
		for _, ac := range m.Actions {
			ac.Action()
		}

		if c.DefaultError != nil {
			errSchema := v.VisitType(c.DefaultError)
			responseRef := &openapi3.ResponseRef{
				Value: openapi3.NewResponse().
					WithDescription("default error").
					WithJSONSchemaRef(errSchema),
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
	before func(*openapi3.Schema)
	after  func(*openapi3.SchemaRef)
}

func (a *RegisterTypeAction) After(f func(*openapi3.SchemaRef)) *RegisterTypeAction {
	if a.after == nil {
		a.after = f
		return a
	}

	prevFn := a.after
	a.after = func(ref *openapi3.SchemaRef) {
		prevFn(ref)
		f(ref)
	}
	return a
}
func (a *RegisterTypeAction) Before(f func(*openapi3.Schema)) *RegisterTypeAction {
	if a.before == nil {
		a.before = f
		return a
	}

	prevFn := a.before
	a.before = func(s *openapi3.Schema) {
		prevFn(s)
		f(s)
	}
	return a
}
func (a *RegisterTypeAction) Description(description string) *RegisterTypeAction {
	return a.After(func(ref *openapi3.SchemaRef) {
		ref.Value.Description = description
	})
}
func (a *RegisterTypeAction) Enum(values ...interface{}) *RegisterTypeAction {
	return a.Before(func(s *openapi3.Schema) {
		s.Enum = values
	})
}
func (a *RegisterTypeAction) Default(value interface{}) *RegisterTypeAction {
	return a.Before(func(s *openapi3.Schema) {
		s.Default = value
	})
}
func (a *RegisterTypeAction) Example(value interface{}) *RegisterTypeAction {
	return a.Before(func(s *openapi3.Schema) {
		s.Example = value
	})
}

func (m *Manager) RegisterType(ob interface{}, modifiers ...func(*openapi3.Schema)) *RegisterTypeAction {
	var ac *RegisterTypeAction
	ac = &RegisterTypeAction{
		registerAction: &registerAction{
			Manager: m,
			Phase:   phase1Action,
			Action: func() {
				if ac.before != nil {
					modifiers = append(modifiers, ac.before)
				}
				s := m.Visitor.VisitType(ob, modifiers...)
				if ac.after != nil {
					ac.after(s)
				}
			},
		},
	}
	m.Actions = append(m.Actions, ac.registerAction)
	return ac
}

type RegisterFuncAction struct {
	*registerAction
	after func(*openapi3.Operation)
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
func (a *RegisterFuncAction) Example(code int, mime string, title string, value interface{}) *RegisterFuncAction {
	return a.After(func(op *openapi3.Operation) {
		ref := op.Responses[strconv.Itoa(code)]
		if ref == nil {
			schemaRef := a.Manager.Visitor.VisitType(value)
			ref = &openapi3.ResponseRef{Value: openapi3.NewResponse().WithJSONSchemaRef(schemaRef).WithDescription("-")}
			op.Responses[strconv.Itoa(code)] = ref
		}
		if ref.Value != nil {
			if ref.Value.Content == nil {
				ref.Value.Content = openapi3.NewContentWithJSONSchema(openapi3.NewSchema())
			}
			mediatype := ref.Value.Content.Get(mime)
			if mediatype == nil {
				mediatype = openapi3.NewMediaType()
				ref.Value.Content[mime] = mediatype
			}
			if mediatype.Example == nil && mediatype.Examples == nil {
				mediatype.Example = value
			} else {
				if mediatype.Examples == nil {
					mediatype.Examples = openapi3.Examples{}
					if mediatype.Example != nil {
						mediatype.Examples["default"] = &openapi3.ExampleRef{Value: &openapi3.Example{
							Value: mediatype.Example,
						}}
						mediatype.Example = nil
					}
				}
				if title == "default" {
					title = "default" + strconv.Itoa(len(mediatype.Examples))
				}
				mediatype.Examples[title] = &openapi3.ExampleRef{Value: &openapi3.Example{
					Value: value,
				}}
			}
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
		registerAction: &registerAction{
			Phase:   phase2Action,
			Manager: m,
			Action: func() {
				op := m.Visitor.VisitFunc(fn, modifiers...)
				if ac.after != nil {
					ac.after(op)
				}
			},
		},
	}
	m.Actions = append(m.Actions, ac.registerAction)
	return ac
}

func (m *Manager) RegisterInterception(rt reflect.Type, intercept func(*shape.Shape) *openapi3.Schema) {
	m.Visitor.Transformer.RegisterInterception(rt, intercept)
}
