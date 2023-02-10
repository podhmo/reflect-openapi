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
	Summary     string
	Description string
}

func Generate(doc *openapi3.T) *Doc {
	endpoints := make([]Endpoint, 0, len(doc.Paths))
	for path, pathItem := range doc.Paths {
		if op := pathItem.Connect; op != nil {
			endpoints = append(endpoints, Endpoint{
				OperationID: op.OperationID,
				Method:      "CONNECT",
				Path:        path,
				Summary:     op.Summary,
				Description: op.Description,
			})
		}
		if op := pathItem.Delete; op != nil {
			endpoints = append(endpoints, Endpoint{
				OperationID: op.OperationID,
				Method:      "DELETE",
				Path:        path,
				Summary:     op.Summary,
				Description: op.Description,
			})
		}
		if op := pathItem.Get; op != nil {
			endpoints = append(endpoints, Endpoint{
				OperationID: op.OperationID,
				Method:      "GET",
				Path:        path,
				Summary:     op.Summary,
				Description: op.Description,
			})
		}
		if op := pathItem.Head; op != nil {
			endpoints = append(endpoints, Endpoint{
				OperationID: op.OperationID,
				Method:      "HEAD",
				Path:        path,
				Summary:     op.Summary,
				Description: op.Description,
			})
		}
		if op := pathItem.Options; op != nil {
			endpoints = append(endpoints, Endpoint{
				OperationID: op.OperationID,
				Method:      "OPTIONS",
				Path:        path,
				Summary:     op.Summary,
				Description: op.Description,
			})
		}
		if op := pathItem.Patch; op != nil {
			endpoints = append(endpoints, Endpoint{
				OperationID: op.OperationID,
				Method:      "PATCH",
				Path:        path,
				Summary:     op.Summary,
				Description: op.Description,
			})
		}
		if op := pathItem.Post; op != nil {
			endpoints = append(endpoints, Endpoint{
				OperationID: op.OperationID,
				Method:      "POST",
				Path:        path,
				Summary:     op.Summary,
				Description: op.Description,
			})
		}
		if op := pathItem.Put; op != nil {
			endpoints = append(endpoints, Endpoint{
				OperationID: op.OperationID,
				Method:      "PUT",
				Path:        path,
				Summary:     op.Summary,
				Description: op.Description,
			})
		}
		if op := pathItem.Trace; op != nil {
			endpoints = append(endpoints, Endpoint{
				OperationID: op.OperationID,
				Method:      "TRACE",
				Path:        path,
				Summary:     op.Summary,
				Description: op.Description,
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
