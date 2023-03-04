package docgen

import (
	"embed"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"strings"
	"text/template"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/podhmo/reflect-openapi/info"
	"github.com/podhmo/reflect-openapi/walknode"
)

//go:embed templates
var fs embed.FS

type Doc struct {
	Title       string
	Version     string
	Description string

	Endpoints []Endpoint
	HTMLs     []HTMLEndpoint
	Objects   []Object

	SkipMetadata bool // skip header metadata
}

type Endpoint struct {
	Method      string
	Path        string
	OperationID string
	DocumentInfo

	HtmlID string
	Tags   string

	Input      Object
	OutputList []Object
	HasExample bool

	GoPositionURL string
}
type HTMLEndpoint Endpoint

type Object struct {
	Name        string
	TypeExpr    string
	TypeString  string
	ContentType string
	DocumentInfo

	HtmlID   string
	Links    []Link
	Examples []Example
}

type DocumentInfo struct {
	Summary     string
	Description string
}

type Link = info.Link

type Example struct {
	Title       string
	Description string
	Value       string
}

func Generate(doc *openapi3.T, info *info.Info) *Doc {
	endpoints := make([]Endpoint, 0, len(doc.Paths))
	htmls := make([]HTMLEndpoint, 0, len(doc.Paths))

	walknode.PathItem(doc, func(pathItem *openapi3.PathItem, path string) {
		walknode.Operation(pathItem, func(op *openapi3.Operation, method string) {
			htmlID := toHtmlID(op.OperationID, method, path) // url
			numOfExamples := 0

			input := Object{Name: "input", TypeString: ActionInputString(doc, info, op)}
			if body := op.RequestBody; body != nil {
				if body.Value != nil { // not support request component
					media := body.Value.Content.Get("application/json")
					if media != nil {
						schema, typ := toInnerSchemaAndTypeExpr(info, media.Schema)
						if sinfo, ok := info.SchemaInfo[schema]; ok {
							// log.Printf("[DEBUG] schema link: %q link input of %q", typ, op.OperationID)
							sinfo.Links = append(sinfo.Links, Link{Title: fmt.Sprintf("input of %s as `%s`", op.OperationID, typ), URL: "#" + htmlID})
						}
						input.TypeExpr = typ
						input.HtmlID = toHtmlID(schema.Title) // TODO: name conflict
					}
				}
			}

			outputList := make([]Object, 0, 2)
			walknode.Response(op, func(ref *openapi3.ResponseRef, name string) {
				if ref.Value != nil { // not support response component
					for contentType, media := range ref.Value.Content {
						output := Object{Name: name, TypeString: ActionOutputString(doc, info, ref, name), ContentType: contentType}
						switch contentType {
						case "application/json":
							schema, typ := toInnerSchemaAndTypeExpr(info, media.Schema)
							if sinfo, ok := info.SchemaInfo[schema]; ok {
								// log.Printf("[DEBUG] schema link: %q link output of %q (%s)", typ, op.OperationID, name)
								sinfo.Links = append(sinfo.Links, Link{Title: fmt.Sprintf("output of %s (%s) as `%s`", op.OperationID, name, typ), URL: "#" + htmlID})
							}
							output.TypeExpr = typ
							output.HtmlID = toHtmlID(schema.Title) // TODO: name conflict

							walknode.Example(media.Examples, func(ref *openapi3.ExampleRef, title string) {
								b, err := json.MarshalIndent(ref.Value.Value, "", "  ")
								if err != nil {
									log.Printf("[INFO ] docgen.Generate() operationID=%q -- %+v", op.OperationID, err)
									b = []byte(fmt.Sprintf(`<! %s>`, err.Error()))
								}
								output.Examples = append(output.Examples, Example{Title: title, Description: ref.Value.Description, Value: string(b)})
							})
							numOfExamples += len(output.Examples)
						default:
							if ref.Value.Description != nil {
								output.DocumentInfo = toDocumentInfo("", "", *ref.Value.Description)
							}
							// log.Printf("[DEBUG] htmls: %q content-type=%q summary=%q", op.OperationID, contentType, output.Summary)
						}
						outputList = append(outputList, output)
					}
				}
			})

			ep := Endpoint{
				OperationID:  op.OperationID,
				Method:       method,
				Path:         path,
				DocumentInfo: toDocumentInfo("", op.Summary, op.Description),
				HtmlID:       htmlID,
				Tags:         strings.Join(op.Tags, " "),

				Input:      input,
				OutputList: outputList,
				HasExample: numOfExamples > 0,
			}
			if op.Extensions != nil {
				if v, ok := op.Extensions["x-go-position"]; ok {
					if v := v.(string); v != "" {
						ep.GoPositionURL = v
					}
				}
			}

			if len(ep.OutputList) >= 1 && ep.OutputList[0].ContentType != "application/json" { // maybe: [200, default] or [200]
				htmls = append(htmls, HTMLEndpoint(ep))
			} else {
				endpoints = append(endpoints, ep)
			}
		})
	})

	var objects []Object
	if doc.Components != nil {
		objects = make([]Object, 0, len(doc.Components.Schemas))
		walknode.Schema(doc, func(ref *openapi3.SchemaRef, k string) {
			schema := info.LookupSchema(ref)

			// log.Printf("[DEBUG] schema: %s\tlinks=%d", ref.Value.Title, len(links))
			object := Object{
				Name:         k,
				TypeString:   TypeString(doc, info, ref),
				DocumentInfo: toDocumentInfo(ref.Value.Title, "", ref.Value.Description),
				HtmlID:       toHtmlID(k),
			}
			if schemaInfo := info.SchemaInfo[schema]; schemaInfo != nil {
				object.Links = info.SchemaInfo[schema].Links
			}

			if schema.Example != nil {
				b, err := json.MarshalIndent(schema.Example, "", "  ")
				if err != nil {
					log.Printf("[INFO ] docgen.Generate() operationID=%q -- %+v", schema.Title, err)
					b = []byte(fmt.Sprintf(`<! %s>`, err.Error()))
				}
				object.Examples = []Example{{Value: string(b)}}
			}

			objects = append(objects, object)
		})
	}
	return &Doc{
		Title:       doc.Info.Title,
		Description: doc.Info.Description,
		Version:     doc.Info.Version,
		Endpoints:   endpoints,
		HTMLs:       htmls,
		Objects:     objects,
	}
}

func WriteDoc(w io.Writer, doc *Doc) error {
	tmpl, err := template.ParseFS(fs, "templates/doc.tmpl")
	if err != nil {
		return fmt.Errorf("lookup template: %w", err)
	}

	if err := tmpl.Execute(w, doc); err != nil {
		return fmt.Errorf("write doc: %w", err)
	}
	return nil
}

func toInnerSchemaAndTypeExpr(info *info.Info, ref *openapi3.SchemaRef) (*openapi3.Schema, string) {
	schema := info.LookupSchema(ref)
	typ := schema.Title
	switch schema.Type {
	case openapi3.TypeArray:
		schema = info.LookupSchema(schema.Items)
		typ = "[]" + schema.Title
	case openapi3.TypeObject:
		if schema.AdditionalProperties.Schema != nil {
			schema = info.LookupSchema(schema.Items)
			typ = "map[string]" + schema.Title
		}
	}
	return schema, typ
}
