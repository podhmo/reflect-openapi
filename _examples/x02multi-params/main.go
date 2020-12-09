package main

import (
	reflectopenapi "github.com/podhmo/reflect-openapi"
)

func Add(x int, y int) int {
	return x + y
}
func Add2(x int, y, z int) int {
	return x + y + z
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
		SkipValidation:     true,
		SkipExtractArglist: false, // default is false, but explicitly
		Resolver:           &reflectopenapi.NoRefResolver{},
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
			op := m.Visitor.VisitFunc(Sum)
			m.Doc.AddOperation("/Sum", "POST", op)
		}
		{
			op := m.Visitor.VisitFunc(Sum2)
			m.Doc.AddOperation("/Sum2", "POST", op)
		}
	})
}
