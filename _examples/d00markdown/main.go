package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/getkin/kin-openapi/openapi3"
	reflectopenapi "github.com/podhmo/reflect-openapi"
	"github.com/podhmo/reflect-openapi/docgen"
	"github.com/podhmo/reflect-openapi/info"
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
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type NewPet struct {
	// Name of the pet
	Name string `json:"name"`
	// Type of the pet
	Tag string `json:"tag,omitempty"`
}

type Pet struct {
	// Unique id of the pet
	ID int64 `json:"id"`

	NewPet
}

type FindPetsInput struct {
	// tags to filter by
	Tags []string `in:"query" query:"tags"`

	// maximum number of results to return
	Limit int32 `in:"query" query:"limit"`
}

// Returns all pets
//
// Returns all pets from the system that the user has access to
// Nam sed condimentum est. Maecenas tempor sagittis sapien, nec rhoncus sem sagittis sit amet. Aenean at gravida augue, ac iaculis sem. Curabitur odio lorem, ornare eget elementum nec, cursus id lectus. Duis mi turpis, pulvinar ac eros ac, tincidunt varius justo. In hac habitasse platea dictumst. Integer at adipiscing ante, a sagittis ligula. Aenean pharetra tempor ante molestie imperdiet. Vivamus id aliquam diam. Cras quis velit non tortor eleifend sagittis. Praesent at enim pharetra urna volutpat venenatis eget eget mauris. In eleifend fermentum facilisis. Praesent enim enim, gravida ac sodales sed, placerat id erat. Suspendisse lacus dolor, consectetur non augue vel, vehicula interdum libero. Morbi euismod sagittis libero sed lacinia.
//
// Sed tempus felis lobortis leo pulvinar rutrum. Nam mattis velit nisl, eu condimentum ligula luctus nec. Phasellus semper velit eget aliquet faucibus. In a mattis elit. Phasellus vel urna viverra, condimentum lorem id, rhoncus nibh. Ut pellentesque posuere elementum. Sed a varius odio. Morbi rhoncus ligula libero, vel eleifend nunc tristique vitae. Fusce et sem dui. Aenean nec scelerisque tortor. Fusce malesuada accumsan magna vel tempus. Quisque mollis felis eu dolor tristique, sit amet auctor felis gravida. Sed libero lorem, molestie sed nisl in, accumsan tempor nisi. Fusce sollicitudin massa ut lacinia mattis. Sed vel eleifend lorem. Pellentesque vitae felis pretium, pulvinar elit eu, euismod sapien.
func FindPets(
	input FindPetsInput,
) []Pet /* pet response */ {
	return nil
}

type AddPetInput struct {
	NewPet // Pet to add to the store /* TODO: */
}

// Creates a new pet
//
// Creates a new pet in the store. Duplicates are allowed
func AddPet(
	input AddPetInput,
) (*Pet /*pet response TODO: */, error) {
	return nil, nil
}

type FindPetByIDInput struct {
	ID string `in:"path" path:"id"`
}

// Returns a pet by ID
//
// Returns a pet based on a single ID
func FindPetByID(input FindPetByIDInput) *Pet { return nil }

type DeletePetInput struct {
	ID string `in:"path" path:"id"`
}

// Deletes a pet by ID
//
// deletes a single pet based on the ID supplied
func DeletePet(input DeletePetInput) {}

func run() error {
	c := &reflectopenapi.Config{
		Info:         info.New(), // need!
		DefaultError: Error{},
	}
	ctx := context.Background()
	tree, err := c.BuildDoc(ctx, func(m *reflectopenapi.Manager) {
		m.Doc.Info.Title = "Swagger Petstore"
		m.Doc.Info.Version = "1.0.0"
		m.Doc.Info.Description = "A sample API that uses a petstore as an example to demonstrate features in the OpenAPI 3.0 specification"

		// TODO:
		m.Doc.Info.TermsOfService = "http://swagger.io/terms/"
		m.Doc.Info.Contact = &openapi3.Contact{
			Name:  "Swagger API Team",
			Email: "apiteam@swagger.io",
			URL:   "http://swagger.io",
		}
		m.Doc.Info.License = &openapi3.License{
			Name: "Apache 2.0",
			URL:  "https://www.apache.org/licenses/LICENSE-2.0.html",
		}

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
	if err := docgen.Docgen(os.Stdout, doc); err != nil {
		return fmt.Errorf("generate: %w", err)
	}
	return nil
}

func mount(m *reflectopenapi.Manager) {
	m.RegisterFunc(FindPets).After(func(op *openapi3.Operation) {
		m.Doc.AddOperation("/pets", "GET", op)
	})
	m.RegisterFunc(AddPet).After(func(op *openapi3.Operation) {
		m.Doc.AddOperation("/pets", "POST", op)
	})
	m.RegisterFunc(FindPetByID).After(func(op *openapi3.Operation) {
		m.Doc.AddOperation("/pets/{id}", "GET", op)
	})
	m.RegisterFunc(DeletePet).After(func(op *openapi3.Operation) {
		m.Doc.AddOperation("/pets/{id}", "DELETE", op)
	})
}
