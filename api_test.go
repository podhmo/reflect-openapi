package reflectopenapi_test

import (
	"context"
	"testing"

	"github.com/getkin/kin-openapi/openapi3"
	reflectopenapi "github.com/podhmo/reflect-openapi"
	"github.com/podhmo/reflect-openapi/pkg/jsonequal"
)

func TestEmpty(t *testing.T) {
	cases := []struct {
		Msg    string
		GenDoc func() (*openapi3.Swagger, error)
		Output string
	}{
		{
			Msg: "empty",
			GenDoc: func() (*openapi3.Swagger, error) {
				c := reflectopenapi.Config{
					SkipValidation: true,
				}
				return c.BuildDoc(context.Background(), func(m *reflectopenapi.Manager) {})
			},
			Output: "{}",
		},
		{
			Msg: "operation only",
			GenDoc: func() (*openapi3.Swagger, error) {
				c := reflectopenapi.Config{
					SkipValidation: true,
				}
				return c.BuildDoc(context.Background(), func(m *reflectopenapi.Manager) {
					op := m.Visitor.VisitFunc(func() string { return "" })
					m.Doc.AddOperation("/ping", "GET", op)
				})
			},
			Output: `
{
  "/ping": {
    "get": {
      "operationId": "github.com/podhmo/reflect-openapi_test.TestEmpty.func2.1.1",
      "responses": {
        "200": {
          "content": {
            "application/json": {
              "schema": {
                "type": "string"
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
  }
}
`,
		},
		{
			Msg: "operation with default error",
			GenDoc: func() (*openapi3.Swagger, error) {
				type Error struct {
					Message string `json:"message"`
				}
				c := reflectopenapi.Config{
					SkipValidation: true,
					DefaultError:   Error{},
					Resolver:       &reflectopenapi.NoRefResolver{},
				}
				return c.BuildDoc(context.Background(), func(m *reflectopenapi.Manager) {
					op := m.Visitor.VisitFunc(func() string { return "" })
					m.Doc.AddOperation("/ping", "GET", op)
				})
			},
			Output: `
{
  "/ping": {
    "get": {
      "operationId": "github.com/podhmo/reflect-openapi_test.TestEmpty.func3.1.1",
      "responses": {
        "200": {
          "content": {
            "application/json": {
              "schema": {
                "type": "string"
              }
            }
          },
          "description": ""
        },
        "default": {
          "content": {
            "application/json": {
              "schema": {
                "properties": {
                  "message": {
                    "type": "string"
                  }
                },
                "type": "object"
              }
            }
          },
          "description": "default error"
        }
      }
    }
  }
}
`,
		},
	}

	for _, c := range cases {
		t.Run(c.Msg, func(t *testing.T) {
			doc, err := c.GenDoc()
			if err != nil {
				t.Fatalf("unexpected error: %+v", err)
			}

			want := c.Output
			got := doc.Paths
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
}
