package main

import (
	"os"
	"strconv"
	"strings"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/guregu/null"
	"github.com/podhmo/nullable"
	reflectopenapi "github.com/podhmo/reflect-openapi"
)

type Person struct {
	ID   string `json:"id" required:"true"`
	Name string `json:"name"`

	NickName  null.String           `json:"nickname"`
	NickName2 null.String           `json:"nickname2"`
	NickName3 nullable.Type[string] `json:"nickname3"`
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
			m.RegisterType(Person{})
		}
	})
}

func installNullable(m *reflectopenapi.Manager) {
	{
		var v null.Bool
		m.RegisterType(v).Before(func(schema *openapi3.Schema) {
			schema.Title = "Null" + schema.Title
			schema.Nullable = true
			schema.Properties = nil
		})
	}
	{
		var v null.Float
		m.RegisterType(v).Before(func(schema *openapi3.Schema) {
			schema.Title = "Null" + schema.Title
			schema.Nullable = true
			schema.Properties = nil
		})
	}
	{
		var v null.Int
		m.RegisterType(v).Before(func(schema *openapi3.Schema) {
			schema.Title = "Null" + schema.Title
			schema.Type = "integer"
			schema.Nullable = true
			schema.Properties = nil
		})
	}
	{
		var v null.String
		m.RegisterType(v).Before(func(schema *openapi3.Schema) {
			schema.Title = "Null" + schema.Title
			schema.Type = "string"
			schema.Nullable = true
			schema.Properties = nil
		})
	}
	{
		var v null.Time
		m.RegisterType(v, func(schema *openapi3.Schema) {
			schema.Title = "Null" + schema.Title
			schema.Type = "string"
			schema.Format = "date-time"
			schema.Nullable = true
			schema.Properties = nil
		})
	}

	{
		var v nullable.Type[string]
		m.RegisterType(v).Before(func(schema *openapi3.Schema) {
			schema.Title = strings.ReplaceAll(strings.ReplaceAll("Nullable"+schema.Title, "[", "_"), "]", "")
			schema.Type = "string"
			schema.Nullable = true
			schema.Properties = nil
		})
	}
}
