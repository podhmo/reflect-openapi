package main

import (
	"context"
	"fmt"
	"log"
	"os"

	reflectopenapi "github.com/podhmo/reflect-openapi"
	"github.com/podhmo/reflect-openapi/docgen"
)

func main() {
	if err := run(); err != nil {
		log.Fatalf("!! %+v", err)
	}
}

func run() error {
	c := &reflectopenapi.Config{}
	ctx := context.Background()
	tree, err := c.BuildDoc(ctx, func(m *reflectopenapi.Manager) {
		m.Doc.Info.Title = "Swagger Petstore"
		m.Doc.Info.Version = "1.0.0"
		m.Doc.Info.Description = "A sample API that uses a petstore as an example to demonstrate features in the OpenAPI 3.0 specification"
	})
	if err != nil {
		return fmt.Errorf("build: %w", err)
	}

	doc := docgen.Generate(tree)
	if err := docgen.Docgen(os.Stdout, doc); err != nil {
		return fmt.Errorf("generate: %w", err)
	}
	return nil
}
