package main

import (
	"context"

	"github.com/getkin/kin-openapi/openapi3"
	reflectopenapi "github.com/podhmo/reflect-openapi"
)

// Action is wrapped type for wrapping function something like DefineGet()
type Action[I any, O any] func(ctx context.Context, input I) (O, error)

// Foo obejct
type Foo struct {
	ID   int
	Name string // name of foo
}

// GetFoo's input parameters (ignored)
type GetFooInput struct {
	ID int `in:"path" path:"id"` // id of foo
}

// GetFoo returns matched foo object.
func GetFoo(ctx context.Context, input GetFooInput) (*Foo, error) { return nil, nil }

func main() {
	c := reflectopenapi.Config{}
	c.EmitDoc(func(m *reflectopenapi.Manager) {
		DefineGet(m, "/foo/{id}", GetFoo)
	})
}

func DefineGet[I any, O any](m *reflectopenapi.Manager, path string, action Action[I, O]) {
	m.RegisterFunc(action).After(func(op *openapi3.Operation) {
		m.Doc.AddOperation(path, "GET", op)
	})
}
