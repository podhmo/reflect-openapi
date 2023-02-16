package docgen

import (
	"embed"
	"fmt"
	"io"
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
}

type Object struct {
	Name       string
	TypeString string
	DocumentInfo

	HtmlID string
	Links  []Link
}

type DocumentInfo struct {
	Summary     string
	Description string
}

type Link = info.Link

func Generate(doc *openapi3.T, info *info.Info) *Doc {
	endpoints := make([]Endpoint, 0, len(doc.Paths))
	walknode.PathItem(doc, func(pathItem *openapi3.PathItem, path string) {
		walknode.Operation(pathItem, func(op *openapi3.Operation, method string) {
			htmlID := toHtmlID(op.OperationID, method, path) // url

			input := Object{Name: "input", TypeString: ActionInputString(doc, info, op)}
			if body := op.RequestBody; body != nil {
				if body.Value != nil { // not support request component
					media := body.Value.Content.Get("application/json")
					if media != nil {
						schema := info.LookupSchema(media.Schema)
						if sinfo, ok := info.SchemaInfo[schema]; ok {
							// log.Printf("[DEBUG] schema link: %q link input of %q", schema.Title, op.OperationID)
							sinfo.Links = append(sinfo.Links, Link{Title: fmt.Sprintf("input of %s", op.OperationID), URL: "#" + htmlID})
						}
					}
				}
			}

			outputList := make([]Object, 0, 2)
			walknode.Response(op, func(ref *openapi3.ResponseRef, name string) {
				outputList = append(outputList, Object{Name: name, TypeString: ActionOutputString(doc, info, ref, name)})
				if ref.Value != nil { // not support response component
					media := ref.Value.Content.Get("application/json")
					if media != nil {
						schema := info.LookupSchema(media.Schema)
						if sinfo, ok := info.SchemaInfo[schema]; ok {
							// log.Printf("[DEBUG] schema link: %q link output of %q (%s)", schema.Title, op.OperationID, name)
							sinfo.Links = append(sinfo.Links, Link{Title: fmt.Sprintf("output of %s (%s)", op.OperationID, name), URL: "#" + htmlID})
						}
					}
				}
			})

			endpoints = append(endpoints, Endpoint{
				OperationID:  op.OperationID,
				Method:       method,
				Path:         path,
				DocumentInfo: toDocumentInfo("", op.Summary, op.Description),
				HtmlID:       htmlID,
				Tags:         strings.Join(op.Tags, " "),

				Input:      input,
				OutputList: outputList,
			})
		})
	})

	var objects []Object
	if doc.Components != nil {
		objects = make([]Object, 0, len(doc.Components.Schemas))
		walknode.Schema(doc, func(ref *openapi3.SchemaRef, k string) {
			schema := info.LookupSchema(ref)
			links := info.SchemaInfo[schema].Links
			// log.Printf("[DEBUG] schema: %s\tlinks=%d", ref.Value.Title, len(links))
			objects = append(objects, Object{
				Name:         k,
				TypeString:   TypeString(doc, info, ref),
				DocumentInfo: toDocumentInfo(ref.Value.Title, "", ref.Value.Description),
				HtmlID:       toHtmlID(k),
				Links:        links,
			})
		})
	}
	return &Doc{
		Title:       doc.Info.Title,
		Description: doc.Info.Description,
		Version:     doc.Info.Version,
		Endpoints:   endpoints,
		Objects:     objects,
	}
}

func Docgen(w io.Writer, doc *Doc) error {
	tmpl, err := template.ParseFS(fs, "templates/doc.tmpl")
	if err != nil {
		return fmt.Errorf("lookup template: %w", err)
	}

	if err := tmpl.Execute(w, doc); err != nil {
		return fmt.Errorf("generate doc: %w", err)
	}
	return nil
}
