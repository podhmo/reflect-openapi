package reflectopenapi_test

import (
	"context"
	"encoding/json"
	"os"

	"github.com/getkin/kin-openapi/openapi3"
	reflectopenapi "github.com/podhmo/reflect-openapi"
)

type Owner struct {
	Name string `json:"name" required:"true" openapi-override:"{'pattern': '^[A-Z][-A-Za-z]+$'}"`
	Age  int    `json:"age"`
}

type ListOwnerInput struct {
	Sort string `json:"sort" in:"query" openapi-override:"{'enum': ['desc', 'asc'], 'default': 'asc'}"`
}

func ListOwner(ctx context.Context, input ListOwnerInput) ([]*Owner, error) {
	return nil, nil
}

func ExampleConfig() {
	c := reflectopenapi.Config{
		TagNameOption: &reflectopenapi.TagNameOption{
			NameTag:        "json",
			ParamTypeTag:   "in",
			DescriptionTag: "description",
			OverrideTag:    "openapi-override",
		},
		SkipValidation: true,
		Extractor:      shapeCfg,
	}
	doc, _ := c.BuildDoc(context.Background(), func(m *reflectopenapi.Manager) {
		m.RegisterFunc(ListOwner).After(func(op *openapi3.Operation) {
			m.Doc.AddOperation("/owners", "GET", op)
		})
	})
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "\t")
	enc.Encode(doc)
	// Output:
	// {
	// 	"components": {
	// 		"schemas": {
	// 			"ListOwnerInput": {
	// 				"default": {
	// 					"sort": ""
	// 				},
	// 				"properties": {
	// 					"sort": {
	// 						"default": "asc",
	// 						"enum": [
	// 							"desc",
	// 							"asc"
	// 						],
	// 						"type": "string"
	// 					}
	// 				},
	// 				"type": "object"
	// 			},
	// 			"Owner": {
	// 				"properties": {
	// 					"age": {
	// 						"type": "integer"
	// 					},
	// 					"name": {
	// 						"pattern": "^[A-Z][-A-Za-z]+$",
	// 						"type": "string"
	// 					}
	// 				},
	// 				"required": [
	// 					"name"
	// 				],
	// 				"type": "object"
	// 			}
	// 		}
	// 	},
	// 	"info": {
	// 		"description": "-",
	// 		"title": "Sample API",
	// 		"version": "0.0.0"
	// 	},
	// 	"openapi": "3.0.0",
	// 	"paths": {
	// 		"/owners": {
	// 			"get": {
	// 				"operationId": "github.com/podhmo/reflect-openapi_test.ListOwner",
	// 				"requestBody": {
	// 					"content": {
	// 						"application/json": {
	// 							"schema": {
	// 								"$ref": "#/components/schemas/ListOwnerInput"
	// 							}
	// 						}
	// 					}
	// 				},
	// 				"responses": {
	// 					"200": {
	// 						"content": {
	// 							"application/json": {
	// 								"schema": {
	// 									"items": {
	// 										"$ref": "#/components/schemas/Owner"
	// 									},
	// 									"type": "array"
	// 								}
	// 							}
	// 						},
	// 						"description": ""
	// 					},
	// 					"default": {
	// 						"description": ""
	// 					}
	// 				}
	// 			}
	// 		}
	// 	},
	// 	"servers": [
	// 		{
	// 			"description": "local development server",
	// 			"url": "http://localhost:8888"
	// 		}
	// 	]
	// }
}
