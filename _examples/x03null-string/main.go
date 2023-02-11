package main

import (
	"os"
	"reflect"
	"strconv"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/guregu/null"
	"github.com/podhmo/nullable"
	reflectopenapi "github.com/podhmo/reflect-openapi"
	reflectshape "github.com/podhmo/reflect-shape"
)

type Person struct {
	ID   string `json:"id" required:"true"`
	Name string `json:"name"`

	NickName  null.String           `json:"nickname,omitempty"`
	NickName2 null.String           `json:"nickname2,omitempty"`
	NickName3 nullable.Type[string] `json:"nickname3,omitempty"`
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
		m.RegisterType(v).After(func(schema *openapi3.Schema) {
			schema.Title = "Null" + schema.Title
			schema.Nullable = true
			schema.Properties = nil
		})
	}
	{
		var v null.Float
		m.RegisterType(v).After(func(schema *openapi3.Schema) {
			schema.Title = "Null" + schema.Title
			schema.Nullable = true
			schema.Properties = nil
		})
	}
	{
		var v null.Int
		m.RegisterType(v).After(func(schema *openapi3.Schema) {
			schema.Title = "Null" + schema.Title
			schema.Type = "integer"
			schema.Nullable = true
			schema.Properties = nil
		})
	}
	{
		var v null.String
		m.RegisterType(v).After(func(schema *openapi3.Schema) {
			schema.Title = "Null" + schema.Title
			schema.Type = "string"
			schema.Nullable = true
			schema.Properties = nil
		})
	}
	{
		var v null.Time
		m.RegisterType(v).After(func(schema *openapi3.Schema) {
			schema.Title = "Null" + schema.Title
			schema.Type = "string"
			schema.Format = "date-time"
			schema.Nullable = true
			schema.Properties = nil
		})
	}

	// or use RegisterInterception (x-go-type is not set)
	{
		var v nullable.Type[string]
		m.RegisterInterception(reflect.TypeOf(v), func(*reflectshape.Shape) *openapi3.Schema {
			schema := openapi3.NewStringSchema().WithNullable()
			schema.Title = "NullableString"
			return schema
		})
	}
}
