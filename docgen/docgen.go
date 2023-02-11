package docgen

import (
	"embed"
	"fmt"
	"io"
	"text/template"

	"github.com/getkin/kin-openapi/openapi3"
)

//go:embed templates
var fs embed.FS

type Doc struct {
	Title       string
	Version     string
	Description string

	Endpoints []Endpoint
}

type Endpoint struct {
	Method      string
	Path        string
	OperationID string
	DocumentInfo

	HtmlID string
}

type DocumentInfo struct {
	Summary     string
	Description string
}

func Generate(doc *openapi3.T) *Doc {
	endpoints := make([]Endpoint, 0, len(doc.Paths))
	for path, pathItem := range doc.Paths {
		if op := pathItem.Connect; op != nil {
			method := "CONNECT"
			endpoints = append(endpoints, Endpoint{
				OperationID:  op.OperationID,
				Method:       method,
				Path:         path,
				DocumentInfo: toDocumentInfo(op.Summary, op.Description),
				HtmlID:       toHtmlID(op.OperationID, method, path),
			})
		}
		if op := pathItem.Delete; op != nil {
			method := "DELETE"
			endpoints = append(endpoints, Endpoint{
				OperationID:  op.OperationID,
				Method:       method,
				Path:         path,
				DocumentInfo: toDocumentInfo(op.Summary, op.Description),
				HtmlID:       toHtmlID(op.OperationID, method, path),
			})
		}
		if op := pathItem.Get; op != nil {
			method := "GET"
			endpoints = append(endpoints, Endpoint{
				OperationID:  op.OperationID,
				Method:       method,
				Path:         path,
				DocumentInfo: toDocumentInfo(op.Summary, op.Description),
				HtmlID:       toHtmlID(op.OperationID, method, path),
			})
		}
		if op := pathItem.Head; op != nil {
			method := "HEAD"
			endpoints = append(endpoints, Endpoint{
				OperationID:  op.OperationID,
				Method:       method,
				Path:         path,
				DocumentInfo: toDocumentInfo(op.Summary, op.Description),
				HtmlID:       toHtmlID(op.OperationID, method, path),
			})
		}
		if op := pathItem.Options; op != nil {
			method := "OPTIONS"
			endpoints = append(endpoints, Endpoint{
				OperationID:  op.OperationID,
				Method:       method,
				Path:         path,
				DocumentInfo: toDocumentInfo(op.Summary, op.Description),
				HtmlID:       toHtmlID(op.OperationID, method, path),
			})
		}
		if op := pathItem.Patch; op != nil {
			method := "PATCH"
			endpoints = append(endpoints, Endpoint{
				OperationID:  op.OperationID,
				Method:       method,
				Path:         path,
				DocumentInfo: toDocumentInfo(op.Summary, op.Description),
				HtmlID:       toHtmlID(op.OperationID, method, path),
			})
		}
		if op := pathItem.Post; op != nil {
			method := "POST"
			endpoints = append(endpoints, Endpoint{
				OperationID:  op.OperationID,
				Method:       method,
				Path:         path,
				DocumentInfo: toDocumentInfo(op.Summary, op.Description),
				HtmlID:       toHtmlID(op.OperationID, method, path),
			})
		}
		if op := pathItem.Put; op != nil {
			method := "PUT"
			endpoints = append(endpoints, Endpoint{
				OperationID:  op.OperationID,
				Method:       method,
				Path:         path,
				DocumentInfo: toDocumentInfo(op.Summary, op.Description),
				HtmlID:       toHtmlID(op.OperationID, method, path),
			})
		}
		if op := pathItem.Trace; op != nil {
			method := "TRACE"
			endpoints = append(endpoints, Endpoint{
				OperationID:  op.OperationID,
				Method:       method,
				Path:         path,
				DocumentInfo: toDocumentInfo(op.Summary, op.Description),
				HtmlID:       toHtmlID(op.OperationID, method, path),
			})
		}
	}
	return &Doc{
		Title:       doc.Info.Title,
		Description: doc.Info.Description,
		Version:     doc.Info.Version,
		Endpoints:   endpoints,
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
