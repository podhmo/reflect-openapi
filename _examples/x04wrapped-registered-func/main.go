package main

import (
	"context"

	"github.com/getkin/kin-openapi/openapi3"
	reflectopenapi "github.com/podhmo/reflect-openapi"
)

type Action[I any, O any] func(ctx context.Context, input I) (O, error)

type Foo struct {
	ID   int
	Name string
}

type GetFooInput struct {
	ID int `in:"path" path:"id"`
}

func GetFoo(ctx context.Context, input GetFooInput) (*Foo, error) { return nil, nil }

func main() {
	c := reflectopenapi.Config{}
	c.EmitDoc(func(m *reflectopenapi.Manager) {
		m.RegisterFunc(GetFoo).After(func(op *openapi3.Operation) {
			m.Doc.AddOperation("/foo/{id}", "GET", op)
		})
	})
}
