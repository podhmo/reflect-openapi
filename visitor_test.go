package reflectopenapi_test

import (
	"context"
	"encoding/json"
	"reflect"
	"strconv"
	"testing"

	"github.com/getkin/kin-openapi/openapi3"
	reflectopenapi "github.com/podhmo/reflect-openapi"
	"github.com/podhmo/reflect-openapi/pkg/arglist"
	"github.com/podhmo/reflect-openapi/pkg/jsonequal"
	"github.com/podhmo/reflect-openapi/pkg/shape"
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
		extractor = &shape.Extractor{
			Seen: map[reflect.Type]shape.Shape{},
		}
	}
	return reflectopenapi.NewVisitor(resolver, selector, extractor)
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
			Msg: "struct, with openapi-tag=query, ignored",
			Input: struct {
				Name  string `json:"name"`
				Query string `json:"query" openapi:"query"`
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
		// pointer
		{
			Msg:    "pointer, *integer",
			Input:  &intN,
			Output: `{"type": "integer"}`,
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
	}

	v := newVisitorDefault(&reflectopenapi.NoRefResolver{})

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
		{
			Msg:   "use merge-params selector",
			Input: func4,
			Output: `
{
  "operationId": "github.com/podhmo/reflect-openapi_test.func4",
  "parameters": [
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
		  "description": ""
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
			Extractor: &shape.Extractor{
				Seen:          map[reflect.Type]shape.Shape{},
				ArglistLookup: arglist.NewLookup(),
			},
		},
	}

	for _, c := range cases {
		t.Run(c.Msg, func(t *testing.T) {
			v := newVisitor(&reflectopenapi.NoRefResolver{}, c.Selector, c.Extractor)
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
	v := newVisitorDefault(r)

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
  "Group": {
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

func TestIsRequiredFunction(t *testing.T) {
	type Person struct {
		ID   string `json:"id"`                   // required
		Name string `json:"name" required:"true"` // required
		Age  int    `json:"age" required:"false"` // unrequired
	}

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
		got := v.VisitType(Person{})
		want := `
{
  "properties": {
    "id": {
      "type": "string"
    },
    "name": {
      "type": "string"
    },
    "age": {
      "type": "integer"
    }
  },
  "required": [
    "id",
    "name"
  ],
  "type": "object"
}
`
		if err := jsonequal.ShouldBeSame(
			jsonequal.FromString(want),
			jsonequal.From(got),
			jsonequal.WithLeftName("want"),
			jsonequal.WithRightName("got"),
		); err != nil {
			t.Errorf("%+v", err)
		}
	})

	t.Run("embedded", func(t *testing.T) {
		type WrapPerson struct {
			Person

			Father *Person `json:"father" required:"false"` // unrequired
			Mother *Person `json:"mother" required:"false"` // unrequired

			FamilyName string `json:"familyName"`
		}

		got := v.VisitType(WrapPerson{})
		want := `
{
  "properties": {
    "age": {
      "type": "integer"
    },
    "familyName": {
      "type": "string"
    },
    "father": {
      "properties": {
        "age": {
          "type": "integer"
        },
        "id": {
          "type": "string"
        },
        "name": {
          "type": "string"
        }
      },
      "required": [
        "id",
        "name"
      ],
      "type": "object"
    },
    "id": {
      "type": "string"
    },
    "mother": {
      "properties": {
        "age": {
          "type": "integer"
        },
        "id": {
          "type": "string"
        },
        "name": {
          "type": "string"
        }
      },
      "required": [
        "id",
        "name"
      ],
      "type": "object"
    },
    "name": {
      "type": "string"
    }
  },
  "required": [
    "id",
    "name",
    "familyName"
  ],
  "type": "object"
}
`
		if err := jsonequal.ShouldBeSame(
			jsonequal.FromString(want),
			jsonequal.From(got),
			jsonequal.WithLeftName("want"),
			jsonequal.WithRightName("got"),
		); err != nil {
			t.Errorf("%+v", err)
		}
	})
}

func func4(ctx context.Context, x, y int, pretty *bool) []int {
	return nil
}
