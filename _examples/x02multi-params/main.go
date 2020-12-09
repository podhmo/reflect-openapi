package main

import (
	"context"
	"fmt"

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
		{
			op := m.Visitor.VisitFunc(Add)
			m.Doc.AddOperation("/Add", "POST", op)
		}
		{
			op := m.Visitor.VisitFunc(Add2)
			m.Doc.AddOperation("/Add2", "POST", op)
		}
		{
			op := m.Visitor.VisitFunc(Hello)
			m.Doc.AddOperation("/Hello", "POST", op)
		}
		{
			op := m.Visitor.VisitFunc(Hello2)
			m.Doc.AddOperation("/Hello2", "POST", op)
		}
		{
			op := m.Visitor.VisitFunc(Hello3)
			m.Doc.AddOperation("/Hello3", "POST", op)
		}
		{
			op := m.Visitor.VisitFunc(Sum)
			m.Doc.AddOperation("/Sum", "POST", op)
		}
		{
			op := m.Visitor.VisitFunc(Sum2)
			m.Doc.AddOperation("/Sum2", "POST", op)
		}
	})
}
