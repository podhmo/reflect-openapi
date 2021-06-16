package main

import (
	"os"
	"strconv"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/guregu/null"
	reflectopenapi "github.com/podhmo/reflect-openapi"
)

type Person struct {
	ID        string      `json:"id" required:"true"`
	Name      string      `json:"name"`
	NickName  null.String `json:"nickname"`
	NickName2 null.String `json:"nickname2"`
}

func main() {
	c := reflectopenapi.Config{
		SkipValidation: true,
	}
	if ok, _ := strconv.ParseBool(os.Getenv("WITHOUT_REF")); ok {
		c.Resolver = &reflectopenapi.NoRefResolver{}
	}
	c.EmitDoc(func(m *reflectopenapi.Manager) {
		installNullable(m)
		{
			m.Visitor.VisitType(Person{})
		}
	})
}

func installNullable(m *reflectopenapi.Manager) {
	{
		var v null.Bool
		m.RegisterType(v, func(schema *openapi3.Schema) {
			schema.Title = "Null" + schema.Title
			schema.Nullable = true
			schema.Properties = nil
		})
	}
	{
		var v null.Float
		m.RegisterType(v, func(schema *openapi3.Schema) {
			schema.Title = "Null" + schema.Title
			schema.Nullable = true
			schema.Properties = nil
		})
	}
	{
		var v null.Int
		m.RegisterType(v, func(schema *openapi3.Schema) {
			schema.Title = "Null" + schema.Title
			schema.Type = "integer"
			schema.Nullable = true
			schema.Properties = nil

		})
	}
	{
		var v null.String
		m.RegisterType(v, func(schema *openapi3.Schema) {
			schema.Title = "Null" + schema.Title
			schema.Type = "string"
			schema.Nullable = true
			schema.Properties = nil
		})
	}
	// {
	// 	var v null.Time
	// 	m.RegisterType(v, func(schema *openapi3.Schema) {
	// 		schema.Nullable = true
	// 	})
	// }
}
