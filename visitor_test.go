package reflectopenapi_test

import (
	"encoding/json"
	"testing"

	"github.com/getkin/kin-openapi/openapi3"
	reflectopenapi "github.com/podhmo/reflect-openapi"
	"github.com/podhmo/reflect-openapi/pkg/jsonequal"
)

func TestVisitType(t *testing.T) {
	cases := []struct {
		Msg    string
		Input  interface{}
		Output string
	}{
		{
			Msg:    "primitive, integer",
			Input:  1,
			Output: `{"type": "integer"}`,
		},
		{
			Msg:    "primitive, string",
			Input:  "foo",
			Output: `{"type": "string"}`,
		},
		{
			Msg:    "primitive, []byte",
			Input:  []byte("foo"),
			Output: `{"type": "string", "format": "binary"}`,
		},
		{
			Msg:    "struct, without json tag",
			Input:  struct{ Name string }{},
			Output: `{"type": "object", "properties": {"Name": {"type": "string"}}}`,
		},
		{
			Msg: "struct, with json tag",
			Input: struct {
				Name string `json:"name"`
			}{},
			Output: `{"type": "object", "properties": {"name": {"type": "string"}}}`,
		},
		{
			Msg: "struct, with openapitag=query, ignored",
			Input: struct {
				Name string `json:"name" openapi:"query"`
			}{},
			Output: `{"type": "object"}`,
		},
	}

	v := reflectopenapi.NewVisitor(&reflectopenapi.NoRefResolver{})

	for _, c := range cases {
		t.Run(c.Msg, func(t *testing.T) {
			got := v.VisitType(c.Input)

			if err := jsonequal.ShouldBeSame(
				jsonequal.FromString(c.Output),
				jsonequal.From(got),
				jsonequal.WithLeftName("want"),
				jsonequal.WithRightName("got"),
			); err != nil {
				t.Errorf("%+v", err)
			}
		})
	}
}

// Function as API endpoint
func TestVisitFunc(t *testing.T) {
	cases := []struct {
		Msg    string
		Input  interface{}
		Output string
	}{
		{
			Msg:   "return value as response",
			Input: func() int { return 1 },
			Output: `
{
  "operationId": "github.com/podhmo/reflect-openapi_test.TestVisitFunc.func1",
  "responses": {
    "200": {
      "content": {
        "application/json": {
          "schema": {
            "type": "integer"
          }
        }
      },
      "description": ""
    },
    "default": {
      "description": ""
    }
  }
}`,
		},
		{
			Msg: "arguments.0 as request body",
			Input: func(data struct {
				Name string
				Age  int
			}) {
			},
			Output: `
{
  "operationId": "github.com/podhmo/reflect-openapi_test.TestVisitFunc.func2",
  "requestBody": {
    "content": {
      "application/json": {
        "schema": {
          "type": "object",
          "properties": {
            "Age": {
              "type": "integer"
            },
            "Name": {
              "type": "string"
            }
          }
        }
      }
    }
  },
  "responses": {
    "default": {
      "description": ""
    }
  }
}
`,
		},
		{
			Msg: "openapi=query,path as parameters",
			Input: func(data struct {
				Name   string
				Age    int
				ID     string `json:"id" openapi:"path"`
				Pretty bool   `json:"pretty" openapi:"query"`
			}) {
			},
			Output: `
{
  "operationId": "github.com/podhmo/reflect-openapi_test.TestVisitFunc.func3",
  "parameters": [
    {
      "in": "path",
      "name": "id",
      "required": true,
      "schema": {"type": "string"}
    },
    {
      "in": "query",
      "name": "pretty",
      "schema": {"type": "boolean"}
    }
  ],
  "requestBody": {
    "content": {
      "application/json": {
        "schema": {
          "type": "object",
          "properties": {
            "Age": {
              "type": "integer"
            },
            "Name": {
              "type": "string"
            }
          }
        }
      }
    }
  },
  "responses": {
    "default": {
      "description": ""
    }
  }
}
`,
		},
	}

	v := reflectopenapi.NewVisitor(&reflectopenapi.NoRefResolver{})

	for _, c := range cases {
		t.Run(c.Msg, func(t *testing.T) {
			got := v.VisitFunc(c.Input)

			if err := jsonequal.ShouldBeSame(
				jsonequal.FromString(c.Output),
				jsonequal.From(got),
				jsonequal.WithLeftName("want"),
				jsonequal.WithRightName("got"),
			); err != nil {
				t.Errorf("%+v", err)
			}
		})
	}
}

func TestWithRef(t *testing.T) {
	type User struct {
		Name string `json:"string"`
	}

	type Group struct {
		Members []User `json:"members"`
	}

	r := &reflectopenapi.UseRefResolver{}
	v := reflectopenapi.NewVisitor(r)

	got := v.VisitType(Group{})

	t.Run("return value is ref", func(t *testing.T) {
		want := `{"$ref": "#/components/schemas/Group"}`

		if err := jsonequal.ShouldBeSame(
			jsonequal.FromString(want),
			jsonequal.From(got),
			jsonequal.WithLeftName("want"),
			jsonequal.WithRightName("got"),
		); err != nil {
			t.Errorf("%+v", err)
		}
	})

	t.Run("there are original definition in schemas", func(t *testing.T) {
		doc := &openapi3.Swagger{}
		r.Bind(doc)

		b, _ := json.Marshal(doc.Components.Schemas)
		want := `{
  "#/components/schemas/Group": {
    "properties": {
      "members": {
        "items": {
          "$ref": "#/components/schemas/User"
        },
        "type": "array"
      }
    },
    "type": "object"
  },
  "#/components/schemas/User": {
    "properties": {
      "string": {
        "type": "string"
      }
    },
    "type": "object"
  }
}
`
		if err := jsonequal.ShouldBeSame(
			jsonequal.FromString(want),
			jsonequal.FromBytes(b),
			jsonequal.WithLeftName("want"),
			jsonequal.WithRightName("got"),
		); err != nil {
			t.Errorf("%+v", err)
		}
	})
}
