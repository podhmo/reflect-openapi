package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"go/token"
	"log"
	"os"

	"github.com/getkin/kin-openapi/openapi3"
	reflectopenapi "github.com/podhmo/reflect-openapi"
	"github.com/podhmo/reflect-openapi/docgen"
	"github.com/podhmo/reflect-openapi/info"
	reflectshape "github.com/podhmo/reflect-shape"
)

var options struct {
	DocFile string
}

func main() {
	flag.StringVar(&options.DocFile, "docfile", "", "write openapi doc file")
	flag.Parse()

	if err := run(); err != nil {
		log.Fatalf("!! %+v", err)
	}
}

type Error struct {
	Code    int32  `json:"code"`    // Error code
	Message string `json:"message"` // Error message
}

type User struct {
	// Name of the user
	Name string `json:"name"`
	// Age of the user
	Age int `json:"age,omitempty"`
}

type GetUserInput struct {
	Pretty bool   `in:"query" query:"pretty"`
	ID     string `in:"path" path:"id"`
}
type GetUserOutput struct {
	User User `json:"user"`
}

// get user
func GetUser(ctx context.Context, input GetUserInput) (*GetUserOutput, error) { return nil, nil }

func run() error {
	c := &reflectopenapi.Config{
		Info:                info.New(), // need!
		DefaultError:        Error{},
		DefaultErrorExample: Error{Code: 444, Message: "unexpected error!"},
		EnableAutoTag:       true,
		DisableInputRef:     true,
		DisableOutputRef:    true,
		GoPositionFunc: func(fset *token.FileSet, fn *reflectshape.Func) string {
			filepos := fset.Position(fn.Pos())
			return fmt.Sprintf("https://github.com/podhmo/reflect-openapi/blob/main/_examples/d02markdown-with-wrapper/main.go#L%d", filepos.Line)
		},
	}
	ctx := context.Background()
	tree, err := c.BuildDoc(ctx, func(m *reflectopenapi.Manager) {
		m.Doc.Info.Title = "Swagger Petstore"
		m.Doc.Info.Version = "1.0.0"
		m.Doc.Info.Description = "A sample API that uses a petstore as an example to demonstrate features in the OpenAPI 3.0 specification"

		mount(m)
	})
	if err != nil {
		return fmt.Errorf("build: %w", err)
	}

	if options.DocFile != "" {
		f, err := os.Create(options.DocFile)
		if err != nil {
			return fmt.Errorf("open openapi doc: %w", err)
		}
		enc := json.NewEncoder(f)
		enc.SetIndent("", "  ")
		if err := enc.Encode(tree); err != nil {
			return fmt.Errorf("write openapi doc: %w", err)
		}
	}

	doc := docgen.Generate(tree, c.Info)
	if err := docgen.WriteDoc(os.Stdout, doc); err != nil {
		return fmt.Errorf("generate: %w", err)
	}
	return nil
}

func mount(m *reflectopenapi.Manager) {
	m.RegisterFunc(GetUser).After(func(op *openapi3.Operation) {
		m.Doc.AddOperation("/users/{id}", "GET", op)
	})
}
