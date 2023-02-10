package docgen

import (
	"embed"
	"fmt"
	"io"
	"regexp"
	"strings"
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

var (
	TRUNCATE_SIZE = 88
)

func toDocumentInfo(summary, description string) (di DocumentInfo) {
	// fmt.Fprintf(os.Stderr, "summary: %q\n", summary)
	// fmt.Fprintf(os.Stderr, "description: %q\n", description)
	// fmt.Fprintf(os.Stderr, "--\n")

	defer func() {
		if len(di.Summary) > TRUNCATE_SIZE {
			di.Summary = di.Summary[:TRUNCATE_SIZE]
		}
	}()

	parts := strings.Split(description, "\n")
	if len(parts) > 3 && strings.TrimSpace(parts[1]) == "" {
		di.Summary = strings.TrimSpace(parts[0])
		di.Description = strings.TrimSpace(strings.Join(parts[2:], "\n"))
		return
	} else if summary != "" {
		di.Summary = strings.TrimSpace(summary)
		di.Description = strings.TrimSpace(description)
		return
	} else {
		di.Summary = strings.TrimSpace(parts[0])
		di.Description = strings.TrimSpace(strings.Join(parts[1:], "\n"))
		return
	}
}

var (
	toDashRegex  = regexp.MustCompile(`[ \t]+`)
	toEmptyRegex = regexp.MustCompile(`[{/\.}]+`)
)

func toHtmlID(operationID, method, path string) string {
	s := fmt.Sprintf("%s %s %s", operationID, method, path)
	s = strings.ToLower(s)
	s = toEmptyRegex.ReplaceAllString(s, "")
	s = toDashRegex.ReplaceAllString(s, "-")
	return s
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
