package main

import (
	"github.com/getkin/kin-openapi/openapi3"
	reflectopenapi "github.com/podhmo/reflect-openapi"
)

type Name string

type Article struct {
	Title string
	Owner Name
}

func main() {
	c := &reflectopenapi.Config{}
	c.EmitDoc(func(m *reflectopenapi.Manager) {
		m.RegisterType(Name("")).After(func(schema *openapi3.Schema) {
			schema.Pattern = `^[A-Z][A-Za-z\-]+$`
			schema.Description = "Name of something"
		})
		m.RegisterType(Article{})
	})
}
