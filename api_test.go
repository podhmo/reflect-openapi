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
					m.RegisterFunc(func() string { return "" }, func(op *openapi3.Operation) {
						m.Doc.AddOperation("/ping", "GET", op)
					})
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
					m.RegisterFunc(func() string { return "" }, func(op *openapi3.Operation) {
						m.Doc.AddOperation("/ping", "GET", op)
					})
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
			if err := jsonequal.NoDiff(
				jsonequal.FromString(want).Named("want"),
				jsonequal.From(got).Named("got"),
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
			m.RegisterType(Sin{})

			type A struct {
				Sin     *Sin
				Message string
			}
			m.RegisterType(A{})
		}
		{
			type Sin struct {
				Name string
				Text string
			}

			// name-conflict is occured
			m.RegisterType(Sin{})

			type B struct {
				Sin     *Sin
				Message string

				RelatedList []*Sin
			}
			m.RegisterType(B{})
		}
		{
			type Sin struct {
				Info string
			}

			// prevent name-conflict by hand
			m.RegisterType(Sin{}).After(func(s *openapi3.Schema) {
				s.Title = "SinForC"
			})

			type C struct {
				Sin *Sin
			}
			m.RegisterType(C{})
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
      "x-go-type": "github.com/podhmo/reflect-openapi_test.Sin"
	}
  }
}
`
	b, err := json.Marshal(doc.Components)
	if err != nil {
		t.Errorf("unexpected marshal error %+v", err)
	}
	if err := jsonequal.NoDiff(
		jsonequal.FromString(want).Named("want"),
		jsonequal.FromBytes(b).Named("got"),
	); err != nil {
		t.Errorf("%+v", err)
	}
}

func TestDefaultInput(t *testing.T) {
	type SimpleInput struct {
		Query string `in:"query" query:"query"`
		Value string `in:"query" query:"value"`
	}

	type EmbeddedInputInput struct {
		SimpleInput
	}
	type EmbeddedInput struct {
		EmbeddedInputInput
		Pretty bool `in:"query" query:"pretty"`
	}

	type WithBodyInput struct {
		Name     string `json:"name"`
		NickName string `json:"nickname"`

		Pretty  bool  `in:"query" query:"pretty"`
		Verbose *bool `in:"query" query:"verbose"`
	}

	type Flags struct {
		BoolOK     bool  `in:"query" query:"boolOK"`
		BoolNG     bool  `in:"query" query:"boolNG"`
		BoolPtrOK  *bool `in:"query" query:"boolPtrOK"`
		BoolPtrNG  *bool `in:"query" query:"boolPtrNG"`
		BoolPtrNil *bool `in:"query" query:"boolPtrNil"`
	}
	type WFlags struct {
		*Flags
	}

	ok := true
	ng := false
	cases := []struct {
		fn                        interface{}
		method, path, operationID string
		defaultValue              interface{}
		want                      string
	}{
		{
			fn:     func(input SimpleInput) {},
			method: "GET", path: "/something", operationID: "ListWithSimpleInput",
			defaultValue: SimpleInput{Query: "desc"},
			want: `{
				"operationId": "ListWithSimpleInput",
				"parameters": [
					{"in": "query", "name": "query", "schema": {"default": "desc", "type": "string"}},
					{"in": "query", "name": "value", "schema": {"type": "string"}}
				],
				"responses": {"default": {"description": ""}}
			}`,
		},
		{
			fn:     func(input *EmbeddedInput) {},
			method: "GET", path: "/something", operationID: "ListWithEmbeddedInput",
			defaultValue: EmbeddedInput{EmbeddedInputInput: EmbeddedInputInput{SimpleInput: SimpleInput{Query: "desc"}}},
			want: `{
				"operationId": "ListWithEmbeddedInput",
				"parameters": [
					{"in": "query", "name": "query", "schema": {"default": "desc", "type": "string"}},
					{"in": "query", "name": "value", "schema": {"type": "string"}},
					{"in": "query", "name": "pretty", "schema": {"default": false, "type": "boolean"}}
				],
				"responses": {"default": {"description": ""}}
			}`,
		},
		{
			fn:     func(ctx context.Context, input WithBodyInput) {},
			method: "POST", path: "/something", operationID: "PostWithBody",
			defaultValue: WithBodyInput{Pretty: true, Verbose: &ok, Name: "john"},
			want: `{
				"operationId": "PostWithBody",
				"parameters": [
					{"in": "query", "name": "pretty", "schema": {"default": true, "type": "boolean"}},
					{"in": "query", "name": "verbose", "schema": {"default": true, "type": "boolean"}}
				],
				"requestBody": {
					"content": {
						"application/json": {
							"schema": {
								"type": "object",
								"properties": {
									"name": {"type": "string", "default": "john"},
									"nickname": {"type": "string"}
								}
							}
						}
					}
				},
				"responses": {"default": {"description": ""}}
			}`,
		},
		{
			fn:     func(ctx context.Context, input Flags) {},
			method: "POST", path: "/flagsEmbedded", operationID: "flagsEmbedded",
			defaultValue: Flags{BoolOK: true, BoolPtrOK: &ok, BoolPtrNG: &ng},
			want: `{
				"operationId": "flagsEmbedded",
				"parameters": [
					{"in": "query", "name": "boolOK", "schema": {"default": true, "type": "boolean"}},
					{"in": "query", "name": "boolNG", "schema": {"default": false, "type": "boolean"}},
					{"in": "query", "name": "boolPtrOK", "schema": {"default": true, "type": "boolean"}},
					{"in": "query", "name": "boolPtrNG", "schema": {"default": false, "type": "boolean"}},
					{"in": "query", "name": "boolPtrNil", "schema": {"type": "boolean"}}
				],
				"responses": {"default": {"description": ""}}
			}`,
		},
		{
			fn:     func(ctx context.Context, input *WFlags) {},
			method: "POST", path: "/wflagsEmbedded", operationID: "wflagsEmbedded",
			defaultValue: &WFlags{Flags: &Flags{BoolOK: true, BoolPtrOK: &ok, BoolPtrNG: &ng}},
			want: `{
				"operationId": "wflagsEmbedded",
				"parameters": [
					{"in": "query", "name": "boolOK", "schema": {"default": true, "type": "boolean"}},
					{"in": "query", "name": "boolNG", "schema": {"default": false, "type": "boolean"}},
					{"in": "query", "name": "boolPtrOK", "schema": {"default": true, "type": "boolean"}},
					{"in": "query", "name": "boolPtrNG", "schema": {"default": false, "type": "boolean"}},
					{"in": "query", "name": "boolPtrNil", "schema": {"type": "boolean"}}
				],
				"responses": {"default": {"description": ""}}
			}`,
		},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.operationID, func(t *testing.T) {
			c := reflectopenapi.Config{
				SkipValidation: true,
				Extractor:      shapeCfg,
				Resolver:       &reflectopenapi.NoRefResolver{},
			}

			doc, err := c.BuildDoc(context.Background(), func(m *reflectopenapi.Manager) {
				m.RegisterFunc(tc.fn).After(func(op *openapi3.Operation) {
					op.OperationID = tc.operationID
					m.Doc.AddOperation(tc.path, tc.method, op)
				}).DefaultInput(tc.defaultValue)
			})
			if err != nil {
				t.Fatalf("c.BuildDoc(): unexpected error: %+v", err)
			}

			want := tc.want
			got := doc.Paths.Find(tc.path).GetOperation(tc.method)
			if err := jsonequal.NoDiff(
				jsonequal.FromString(want).Named("want"),
				jsonequal.From(got).Named("got"),
			); err != nil {
				t.Errorf("%+v", err)
			}
		})
	}
}
