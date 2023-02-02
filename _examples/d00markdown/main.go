package main

import (
	"log"
	"os"

	"github.com/podhmo/reflect-openapi/docgen"
)

func main() {
	doc := &docgen.Doc{
		Title:       "Swagger Petstore",
		Version:     "1.0.0",
		Description: "A sample API that uses a petstore as an example to demonstrate features in the OpenAPI 3.0 specification",
	}
	if err := docgen.Docgen(os.Stdout, doc); err != nil {
		log.Fatalf("!! %+v", err)
	}
}
