package main

import (
	"context"
	"fmt"

	"github.com/getkin/kin-openapi/openapi3"
	reflectopenapi "github.com/podhmo/reflect-openapi"
)

type Greeter interface {
	Greet(string)
}

func Add(x int, y int) int {
	return x + y
}
func Add2(x int, y, z int) int {
	return x + y + z
}
func Hello(name string, pretty *bool) (string, error) {
	if *pretty {
		return fmt.Sprintf("** Hello %s **", name), nil
	}
	return fmt.Sprintf("Hello %s", name), nil
}
func Hello2(ctx context.Context, g Greeter, name string, pretty *bool) (string, error) {
	if *pretty {
		return fmt.Sprintf("** Hello %s **", name), nil
	}
	return fmt.Sprintf("Hello %s", name), nil
}

type Person struct {
	Name string `json:"name"`
}

func Hello3(subject *Person, object string, pretty *bool) (string, error) {
	var prefix string
	if subject != nil {
		prefix = fmt.Sprintf("%s: ", subject.Name)
	}
	if *pretty {
		return fmt.Sprintf("** %sHello %s **", prefix, object), nil
	}
	return fmt.Sprintf("%sHello %s", prefix, object), nil
}
func Sum(xs []int) int {
	n := 0
	for _, x := range xs {
		n += x
	}
	return n
}
func Sum2(xs ...int) int {
	n := 0
	for _, x := range xs {
		n += x
	}
	return n
}

func main() {
	c := reflectopenapi.Config{
		SkipValidation: true,
		Selector: &struct {
			reflectopenapi.MergeParamsInputSelector
			reflectopenapi.FirstParamOutputSelector
		}{},
	}
	c.EmitDoc(func(m *reflectopenapi.Manager) {
		m.RegisterFunc(Add).After(func(op *openapi3.Operation) {
			m.Doc.AddOperation("/Add", "POST", op)
		})
		m.RegisterFunc(Add2).After(func(op *openapi3.Operation) {
			m.Doc.AddOperation("/Add2", "POST", op)
		})
		m.RegisterFunc(Hello).After(func(op *openapi3.Operation) {
			m.Doc.AddOperation("/Hello", "POST", op)
		})
		m.RegisterFunc(Hello2).After(func(op *openapi3.Operation) {
			m.Doc.AddOperation("/Hello2", "POST", op)
		})
		m.RegisterFunc(Hello3).After(func(op *openapi3.Operation) {
			m.Doc.AddOperation("/Hello3", "POST", op)
		})
		m.RegisterFunc(Sum).After(func(op *openapi3.Operation) {
			m.Doc.AddOperation("/Sum", "POST", op)
		})
		m.RegisterFunc(Sum2).After(func(op *openapi3.Operation) {
			m.Doc.AddOperation("/Sum2", "POST", op)
		})
	})
}
