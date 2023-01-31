package reflectopenapi_test

import (
	"context"
	"encoding/json"
	"reflect"
	"strconv"
	"testing"

	"github.com/getkin/kin-openapi/openapi3"
	reflectopenapi "github.com/podhmo/reflect-openapi"
	"github.com/podhmo/reflect-openapi/pkg/jsonequal"
)

func newVisitor(
	resolver reflectopenapi.Resolver,
	selector reflectopenapi.Selector,
	extractor reflectopenapi.Extractor,
) *reflectopenapi.Visitor {
	if selector == nil {
		selector = &reflectopenapi.DefaultSelector{}
	}
	if extractor == nil {
		extractor = shapeCfg
	}
	return reflectopenapi.NewVisitor(*reflectopenapi.DefaultTagNameOption(), resolver, selector, extractor)
}
func newVisitorDefault(
	resolver reflectopenapi.Resolver,
) *reflectopenapi.Visitor {
	return newVisitor(resolver, nil, nil)
}

func TestVisitType(t *testing.T) {
	intN := 1
	cases := []struct {
		Msg    string
		Input  interface{}
		Output string
	}{
		{
			Msg:    "primitive, integer",
			Input:  intN,
			Output: `{"type": "integer", "title": "int"}`,
		},
		{
			Msg:    "primitive, string",
			Input:  "foo",
			Output: `{"type": "string", "title": "string"}`,
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
			Msg: "struct, without json tag, unexported",
			Input: struct {
				Name       string
				unexported string
			}{},
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
			Msg: "struct, with json tag omitempty",
			Input: struct {
				Name string `json:"name,omitempty"`
			}{},
			Output: `{"type": "object", "properties": {"name": {"type": "string"}}}`,
		},
		{
			Msg: "struct, with in-tag=query, ignored",
			Input: struct {
				Name  string `json:"name"`
				Query string `json:"query" in:"query"`
			}{},
			Output: `{"type": "object", "properties": {"name": {"type": "string"}}}`,
		},
		{
			Msg: "struct, with -, ignored",
			Input: struct {
				Name string `json:"name"`
				Code int    `json:"-"`
			}{},
			Output: `{"type": "object", "properties": {"name": {"type": "string"}}}`,
		},
		{
			Msg: "struct, with openapi-override",
			Input: struct {
				Name string `json:"name" openapi-override:"{'pattern': '^[A-Z][a-zA-Z]+$'}"`
				Age  int    `json:"age" openapi-override:"{'minimum': 0, \"maximum\": 100}"`
			}{},
			Output: `{"type": "object", "properties": {"name": {"type": "string", "pattern": "^[A-Z][a-zA-Z]+$"}, "age": {"type": "integer", "maximum": 100, "minimum": 0}}}`,
		},
		{
			Msg: "struct, with openapi-override2",
			Input: struct {
				Name string `json:"name" openapi-override:"{'pattern': '^Test\\d+$'}"`
			}{},
			Output: `{"type": "object", "properties": {"name": {"type": "string", "pattern": "^Test\\d+$"}}}`,
		},
		// pointer
		{
			Msg:    "pointer, *integer",
			Input:  &intN,
			Output: `{"type": "integer", "title": "int"}`,
		},
		// slice
		{
			Msg:    "slice",
			Input:  []int{},
			Output: `{"type": "array", "items": {"type": "integer"}}`,
		},
		// map
		{
			Msg: "struct, for map[string, primitive] field",
			Input: struct {
				Points map[string]int `json:"points"`
			}{},
			Output: `{"type": "object", "properties": {"points": {"additionalProperties": {"type": "integer"}}}}`,
		},
		{
			Msg: "struct, for map[string, struct] field",
			Input: struct {
				Metadata map[string]struct {
					Field string `json:"field"`
					Type  string `json:"type"`
					Value string `json:"value"`
				} `json:"metadata"`
			}{},
			Output: `{"type": "object", "properties": {"metadata": {"additionalProperties": {"type": "object", "properties": {"type": {"type": "string"}, "value": {"type": "string"}, "field": {"type": "string"}}}}}}`,
		},
		// interface
		{
			Msg: "struct, for empty interface field",
			Input: struct {
				Metadata interface{} `json:"metadata"`
			}{},
			Output: `{"type": "object", "properties": {"metadata": {"type": "object", "additionalProperties": true, "description": "<Any type>"}}}`,
		},
		{
			Msg: "struct, for interface field",
			Input: struct {
				Name     string                          `json:"name"`
				Metadata interface{ Get(string) string } `json:"metadata"`
			}{},
			Output: `{"type": "object", "properties": {"name": {"type": "string"}}}`,
		},
		// unclear definition
		{
			Msg: "struct, for unclear definition",
			Input: struct {
				i interface{}
			}{},
			Output: `{"type": "object", "additionalProperties": true, "description": "<unclear definition>"}`,
		},
		{
			Msg: "struct, zero",
			Input: struct {
			}{},
			Output: `{"type": "object"}`,
		},
	}

	v := newVisitorDefault(&reflectopenapi.NoRefResolver{})

	for _, c := range cases {
		t.Run(c.Msg, func(t *testing.T) {
			got := v.VisitType(v.Transformer.Extractor.Extract(c.Input))

			if err := jsonequal.NoDiff(
				jsonequal.FromString(c.Output).Named("want"),
				jsonequal.From(got).Named("got"),
			); err != nil {
				t.Errorf("%+v", err)
			}
		})
	}
}

// Function as API endpoint
func TestVisitFunc(t *testing.T) {
	cases := []struct {
		Msg       string
		Input     interface{}
		Output    string
		Selector  reflectopenapi.Selector
		Extractor reflectopenapi.Extractor
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
			Msg: "in=query,path as parameters",
			Input: func(data struct {
				Name   string
				Age    int
				ID     string `path:"id" in:"path"`
				Pretty bool   `query:"pretty" in:"query"`
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
		{
			Msg:   "embedded parameters",
			Input: func4,
			Output: `
{
  "operationId": "github.com/podhmo/reflect-openapi_test.func4",
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
            "age": {
              "type": "integer"
            },
            "name": {
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
			Msg:   "use merge-params selector",
			Input: func5,
			Output: `
{
  "operationId": "github.com/podhmo/reflect-openapi_test.func5",
  "parameters": [
    {
      "in": "query",
      "name": "pretty",
      "schema": {"type": "boolean"},
	  "description": "pretty output or not"
    }
  ],
  "requestBody": {
	  "content": {
		  "application/json": {
			  "schema": {
				  "properties": {
					  "x": {
						  "type": "integer"
					  },
					  "y": {
						  "type": "integer"
					  }
				  },
                  "required": [
                      "x",
                      "y"
                  ],
				  "type": "object"
			  }
		  }
	  }
  },
  "responses": {
	  "200": {
		  "content": {
			  "application/json": {
				  "schema": {
					  "items": {
						  "type": "integer"
					  },
					  "type": "array"
				  }
			  }
		  },
		  "description": "ans"
	  },
	  "default": {
		  "description": ""
	  }
  }
}
`,
			Selector: &struct {
				reflectopenapi.MergeParamsInputSelector
				reflectopenapi.FirstParamOutputSelector
			}{},
			Extractor: shapeCfg,
		},
		{
			Msg:   "use method",
			Input: new(S).M,
			Output: `{
				"description":"This is Method sample",
				"summary":"This is Method sample",
				"operationId":"github.com/podhmo/reflect-openapi_test.S.M",
				"requestBody":{
					"content":{
						"application/json":{
							"schema":{"properties":{"X":{"type":"integer"},"Y":{"type":"integer"}},"type":"object"}}}},
							"responses":{"200":{"content":{"application/json":{"schema":{"items":{"type":"integer"},"type":"array"}}},
							"description":""},
							"default":{"description":""}}}`,
		},
	}

	for _, c := range cases {
		t.Run(c.Msg, func(t *testing.T) {
			v := newVisitor(&reflectopenapi.NoRefResolver{}, c.Selector, c.Extractor)
			got := v.VisitFunc(v.Extractor.Extract(c.Input))

			if err := jsonequal.NoDiff(
				jsonequal.FromString(c.Output).Named("want"),
				jsonequal.From(got).Named("got"),
			); err != nil {
				t.Errorf("%+v", err)
			}
		})
	}
}

type EmbeddedParametersInput struct {
	Name string `json:"name"`
	EmbeddedParametersInputInner
}
type EmbeddedParametersInputInner struct {
	ID string `path:"id" in:"path"`
	EmbeddedParametersInputInnerInner
}
type EmbeddedParametersInputInnerInner struct {
	Age    int  `json:"age"`
	Pretty bool `query:"pretty" in:"query"`
}

func func4(
	ctx context.Context,
	input EmbeddedParametersInput,
) {
}

func func5(
	ctx context.Context,
	x, y int,
	pretty *bool, // pretty output or not
) []int /* ans */ {
	return nil
}

type S struct{}

// This is Method sample
func (s *S) M(ctx context.Context, input struct{ X, Y int }) []int { return nil }

type User struct {
	Name string `json:"string"`
}

type Group struct {
	Members []User `json:"members"`
}

func TestWithRef(t *testing.T) {
	r := &reflectopenapi.UseRefResolver{NameStore: reflectopenapi.NewNameStore()}
	v := newVisitorDefault(r)

	got := v.VisitType(v.Extractor.Extract(Group{}))

	t.Run("return value is ref", func(t *testing.T) {
		want := `{"$ref": "#/components/schemas/Group"}`

		if err := jsonequal.NoDiff(
			jsonequal.FromString(want).Named("want"),
			jsonequal.From(got).Named("got"),
		); err != nil {
			t.Errorf("%+v", err)
		}
	})

	t.Run("there are original definition in schemas", func(t *testing.T) {
		doc := &openapi3.T{}
		r.BindSchemas(doc)

		b, _ := json.Marshal(doc.Components.Schemas)
		want := `{
  "Group": {
    "properties": {
      "members": {
        "items": {
          "$ref": "#/components/schemas/User"
        },
        "type": "array"
      }
    },
    "title": "Group",
    "type": "object"
  },
  "User": {
    "properties": {
      "string": {
        "type": "string"
      }
    },
    "type": "object"
  }
}
`
		if err := jsonequal.NoDiff(
			jsonequal.FromString(want).Named("want"),
			jsonequal.FromBytes(b).Named("got"),
		); err != nil {
			t.Errorf("%+v", err)
		}
	})
}

type Person struct {
	ID   string `json:"id"`                   // required
	Name string `json:"name" required:"true"` // required
	Age  int    `json:"age" required:"false"` // unrequired
}

type WrapPerson struct {
	Person

	Father *Person `json:"father" required:"false"` // unrequired
	Mother *Person `json:"mother" required:"false"` // unrequired

	FamilyName string `json:"familyName"`
}

func TestIsRequiredFunction(t *testing.T) {
	r := &reflectopenapi.NoRefResolver{}
	v := newVisitorDefault(r)
	v.IsRequired = func(tag reflect.StructTag) bool {
		v, exists := tag.Lookup("required")
		if !exists {
			return true // required
		}
		required, err := strconv.ParseBool(v)
		if err != nil {
			return false // unrequired
		}
		return required
	}

	t.Run("plain", func(t *testing.T) {
		got := v.VisitType(v.Extractor.Extract(Person{}))
		want := `
{
  "properties": {
    "id": {
		"description": "required",
		"type": "string"
    },
    "name": {
		"description": "required",
		"type": "string"
    },
    "age": {
		"description": "unrequired",
		"type": "integer"
    }
  },
  "required": [
    "id",
    "name"
  ],
  "title": "Person",
  "type": "object"
}
`
		if err := jsonequal.NoDiff(
			jsonequal.FromString(want).Named("want"),
			jsonequal.From(got).Named("got"),
		); err != nil {
			t.Errorf("%+v", err)
		}
	})

	t.Run("embedded", func(t *testing.T) {
		got := v.VisitType(v.Extractor.Extract(WrapPerson{}))
		want := `
{
  "properties": {
    "age": {
		"description": "unrequired",
		"type": "integer"
    },
    "familyName": {
      "type": "string"
    },
    "father": {
      "properties": {
        "age": {
			"description": "unrequired",
			"type": "integer"
        },
        "id": {
			"description": "required",
			"type": "string"
        },
        "name": {
			"description": "required",
			"type": "string"
        }
      },
      "required": [
        "id",
        "name"
      ],
      "title": "Person",
      "type": "object"
    },
    "id": {
		"description": "required",
		"type": "string"
    },
    "mother": {
      "properties": {
        "age": {
			"description": "unrequired",
			"type": "integer"
        },
        "id": {
			"description": "required",
			"type": "string"
        },
        "name": {
			"description": "required",
			"type": "string"
        }
      },
      "required": [
        "id",
        "name"
      ],
      "title": "Person",
      "type": "object"
    },
    "name": {
		"description": "required",
		"type": "string"
    }
  },
  "required": [
    "id",
    "name",
    "familyName"
  ],
  "title": "WrapPerson",
  "type": "object"
}
`
		if err := jsonequal.NoDiff(
			jsonequal.FromString(want).Named("want"),
			jsonequal.From(got).Named("got"),
		); err != nil {
			t.Errorf("%+v", err)
		}
	})
}
