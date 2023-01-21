package main

import (
	"reflect"

	"github.com/getkin/kin-openapi/openapi3"
	reflectopenapi "github.com/podhmo/reflect-openapi"
	shape "github.com/podhmo/reflect-shape"
)

type Todo struct {
	ID    string `json:"id" required:"true"`
	Title string `json:"title"`
	Done  bool   `json:"done"`
}

func ListTodo() []Todo {
	return nil
}
func GetTodo(params struct {
	ID string `json:"id" in:"path"`
}) *Todo {
	return nil
}

type CustomSelector struct {
	reflectopenapi.Selector
	Extractor reflectopenapi.Extractor
}

// wrap with {"items": <>}
func (s *CustomSelector) SelectOutput(fn *shape.Func) *shape.Shape {
	out := s.Selector.SelectOutput(fn)
	if out.Kind == reflect.Slice {
		rt := reflect.StructOf([]reflect.StructField{
			{
				Name: "Items",
				Type: out.Type,
				Tag:  `json:"items"`,
			},
			{
				Name: "HasNext",
				Type: reflect.TypeOf(false),
				Tag:  `json:"hasNext"`,
			},
		})
		return s.Extractor.Extract(reflect.New(rt).Interface())
	}
	return out
}

func main() {
	c := reflectopenapi.Config{
		SkipValidation: true,
	}
	c.Selector = &CustomSelector{
		Selector:  c.DefaultSelector(),
		Extractor: c.DefaultExtractor(),
	}
	c.EmitDoc(func(m *reflectopenapi.Manager) {
		m.RegisterFunc(ListTodo).After(func(op *openapi3.Operation) {
			m.Doc.AddOperation("/todo", "GET", op)
		})
		m.RegisterFunc(GetTodo).After(func(op *openapi3.Operation) {
			m.Doc.AddOperation("/todo/{id}", "GET", op)
		})
	})
}
