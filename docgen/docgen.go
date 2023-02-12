package docgen

import (
	"embed"
	"fmt"
	"io"
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
}

type Endpoint struct {
	Method      string
	Path        string
	OperationID string
	DocumentInfo

	HtmlID string
}

type Object struct {
	Name       string
	TypeString string
	DocumentInfo

	HtmlID string
}

type DocumentInfo struct {
	Summary     string
	Description string
}

func Generate(doc *openapi3.T, info *info.Info) *Doc {
	endpoints := make([]Endpoint, 0, len(doc.Paths))
	walknode.PathItem(doc, func(pathItem *openapi3.PathItem, path string) {
		walknode.Operation(pathItem, func(op *openapi3.Operation, method string) {
			endpoints = append(endpoints, Endpoint{
				OperationID:  op.OperationID,
				Method:       method,
				Path:         path,
				DocumentInfo: toDocumentInfo(op.Summary, op.Description),
				HtmlID:       toHtmlID(op.OperationID, method, path),
			})
		})
	})

	var objects []Object
	if doc.Components != nil {
		objects = make([]Object, 0, len(doc.Components.Schemas))
		walknode.Schema(doc, func(ref *openapi3.SchemaRef, k string) {
			objects = append(objects, Object{
				Name:         k,
				TypeString:   TypeString(doc, info, ref),
				DocumentInfo: toDocumentInfo("", ref.Value.Description),
				HtmlID:       toHtmlID(k),
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
