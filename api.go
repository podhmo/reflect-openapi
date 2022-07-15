package reflectopenapi

import (
	"context"
	"encoding/json"
	"os"
	"reflect"
	"sort"
	"strconv"
	"strings"

	"github.com/getkin/kin-openapi/openapi3"
	shape "github.com/podhmo/reflect-shape"
	"github.com/podhmo/reflect-shape/arglist"
	"github.com/podhmo/reflect-shape/comment"
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

type Config struct {
	Doc *openapi3.T

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
	c.Extractor = &shape.Extractor{
		Seen: map[reflect.Type]shape.Shape{},
	}
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
		c.DefaultResolver(),
		c.DefaultSelector(),
		c.DefaultExtractor(),
	)
	if c.IsRequiredCheckFunction != nil {
		v.Transformer.IsRequired = c.IsRequiredCheckFunction
	}
	if !c.SkipExtractComments {
		v.CommentLookup = comment.NewLookup()
	}
	if _, ok := c.Selector.(useArglist); ok {
		if e, ok := v.extractor.(*shape.Extractor); ok {
			e.ArglistLookup = arglist.NewLookup()
		}
	}
	m := &Manager{
		Doc:      c.Doc,
		Resolver: c.Resolver,
		Visitor:  v,
	}

	return m, func(ctx context.Context) error {
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

		if !c.SkipValidation {
			if err := m.Doc.Validate(ctx); err != nil {
				return err
			}
		}
		return nil
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
	Phase  int // lowerst is first
	Action func()
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

func (m *Manager) RegisterType(ob interface{}, modifiers ...func(*openapi3.Schema)) *RegisterTypeAction {
	var ac *RegisterTypeAction
	ac = &RegisterTypeAction{
		registerAction: &registerAction{
			Phase: phase1Action,
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

// func (a *RegisterFuncAction) AnotherError(code int, typ interface{}) *RegisterFuncAction {
// 	return a.After(func(op *openapi3.Operation) {
// 		op.Responses[strconv.Itoa(code)] = // TODO: implement
// 	})
// }

func (m *Manager) RegisterFunc(fn interface{}, modifiers ...func(*openapi3.Operation)) *RegisterFuncAction {
	var ac *RegisterFuncAction
	ac = &RegisterFuncAction{
		registerAction: &registerAction{
			Phase: phase2Action,
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
