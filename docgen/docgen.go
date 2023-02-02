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
}

func Generate(doc *openapi3.T) *Doc {
	return &Doc{
		Title:       doc.Info.Title,
		Description: doc.Info.Description,
		Version:     doc.Info.Version,
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
