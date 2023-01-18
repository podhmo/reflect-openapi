package reflectopenapi_test

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/getkin/kin-openapi/openapi3"
	reflectopenapi "github.com/podhmo/reflect-openapi"
	"github.com/podhmo/reflect-openapi/pkg/jsonequal"
	shape "github.com/podhmo/reflect-shape"
)

var shapeCfg = &shape.Config{
	IncludeGoTestFiles: true,
	FillArgNames:       true,
	FillReturnNames:    true,
}

func TestEmpty(t *testing.T) {
	cases := []struct {
		Msg    string
		GenDoc func() (*openapi3.T, error)
		Output string
	}{
		{
			Msg: "empty",
			GenDoc: func() (*openapi3.T, error) {
				c := reflectopenapi.Config{
					SkipValidation: true,
					Extractor:      shapeCfg,
				}
				return c.BuildDoc(context.Background(), func(m *reflectopenapi.Manager) {})
			},
			Output: "{}",
		},
		{
			Msg: "operation only",
			GenDoc: func() (*openapi3.T, error) {
				c := reflectopenapi.Config{
					SkipValidation: true,
					Extractor:      shapeCfg,
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
			GenDoc: func() (*openapi3.T, error) {
				type Error struct {
					Message string `json:"message"`
				}
				c := reflectopenapi.Config{
					SkipValidation: true,
					DefaultError:   Error{},
					Resolver:       &reflectopenapi.NoRefResolver{},
					Extractor:      shapeCfg,
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
                "title": "Error",
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

func TestNameConflict(t *testing.T) {
	c := reflectopenapi.Config{
		SkipValidation: true,
		Extractor:      shapeCfg,
	}

	doc, err := c.BuildDoc(context.Background(), func(m *reflectopenapi.Manager) {
		{
			type Sin struct {
				Value float64
			}
			m.Visitor.VisitType(Sin{})

			type A struct {
				Sin     *Sin
				Message string
			}
			m.Visitor.VisitType(A{})
		}
		{
			type Sin struct {
				Name string
				Text string
			}

			// name-conflict is occured
			m.Visitor.VisitType(Sin{})

			type B struct {
				Sin     *Sin
				Message string

				RelatedList []*Sin
			}
			m.Visitor.VisitType(B{})
		}
		{
			type Sin struct {
				Info string
			}

			// prevent name-conflict by hand
			m.Visitor.VisitType(Sin{}, func(s *openapi3.Schema) {
				s.Title = "SinForC"
			})

			type C struct {
				Sin *Sin
			}
			m.Visitor.VisitType(C{})
		}
	})

	if err != nil {
		t.Fatalf("unexpected error %+v", err)
	}

	want := `
{
  "schemas": {
	"A": {
	  "properties": {
		"Message": {
		  "type": "string"
		},
		"Sin": {
		  "$ref": "#/components/schemas/Sin"
		}
	  },
	  "title": "A",
	  "type": "object"
	},
	"B": {
	  "properties": {
		"Message": {
		  "type": "string"
		},
		"RelatedList": {
		  "items": {
			"$ref": "#/components/schemas/Sin01"
		  },
		  "type": "array"
		},
		"Sin": {
		  "$ref": "#/components/schemas/Sin01"
		}
	  },
	  "title": "B",
	  "type": "object"
	},
	"C": {
		"type": "object",
		"properties": {
			"Sin": {
				"$ref": "#/components/schemas/SinForC"
			}
		},
		"title": "C"
	},
	"Sin": {
	  "properties": {
		"Value": {
		  "type": "number"
		}
	  },
	  "title": "Sin",
	  "type": "object",
	  "x-go-id": "github.com/podhmo/reflect-openapi_test.Sin"
	},
	"Sin01": {
	  "properties": {
		"Name": {
		  "type": "string"
		},
		"Text": {
		  "type": "string"
		}
	  },
	  "title": "Sin",
	  "type": "object",
	  "x-go-id": "github.com/podhmo/reflect-openapi_test.Sin"
	},
	"SinForC": {
	  "properties": {
		  "Info": {
			  "type": "string"
		  }
	  },
	  "title": "SinForC",
	  "type": "object",
      "x-new-type": "github.com/podhmo/reflect-openapi_test.Sin"
	}
  }
}
`
	b, err := json.Marshal(doc.Components)
	if err != nil {
		t.Errorf("unexpected marshal error %+v", err)
	}
	if err := jsonequal.ShouldBeSame(
		jsonequal.FromString(want),
		jsonequal.FromBytes(b),
		jsonequal.WithLeftName("want"),
		jsonequal.WithRightName("got"),
	); err != nil {
		t.Errorf("%+v", err)
	}
}
