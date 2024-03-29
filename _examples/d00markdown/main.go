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

type NewPet struct {
	// Name of the pet
	Name string `json:"name"`
	// Type of the pet
	Tag string `json:"tag,omitempty"`
}

// Pet : pet object.
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
	NewPet
}

// Creates a new pet
//
// Creates a new pet in the store. Duplicates are allowed
func AddPet(
	input AddPetInput, // Pet to add to the store
) (*Pet /*pet response */, error /* unexpected error*/) {
	return nil, nil
}

type FindPetByIDInput struct {
	ID int64 `in:"path" path:"id"` // ID of pet to fetch
}

// Returns a pet by ID
//
// Returns a pet based on a single ID
func FindPetByID(input FindPetByIDInput) *Pet/* pet response */ { return nil }

type DeletePetInput struct {
	ID int64 `in:"path" path:"id"` // ID of pet to delete
}

// Deletes a pet by ID
//
// deletes a single pet based on the ID supplied
func DeletePet(input DeletePetInput) struct{}/* pet deleted */ { return struct{}{} }

func run() error {
	c := &reflectopenapi.Config{
		Info:                info.New(), // need!
		DefaultError:        Error{},
		DefaultErrorExample: Error{Code: 444, Message: "unexpected error!"},
		EnableAutoTag:       true,
		GoPositionFunc: func(fset *token.FileSet, fn *reflectshape.Func) string {
			filepos := fset.Position(fn.Pos())
			return fmt.Sprintf("https://github.com/podhmo/reflect-openapi/blob/main/_examples/d00markdown/main.go#L%d", filepos.Line)
		},
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
	if err := docgen.WriteDoc(os.Stdout, doc); err != nil {
		return fmt.Errorf("generate: %w", err)
	}
	return nil
}

func mount(m *reflectopenapi.Manager) {
	m.RegisterFunc(FindPets).After(func(op *openapi3.Operation) {
		m.Doc.AddOperation("/pets", "GET", op)
		op.Tags = []string{"pet", "read"}
	}).Example(200, "application/json", "", "sample output", []Pet{{ID: 1, NewPet: NewPet{Name: "foo", Tag: "A"}}, {ID: 2, NewPet: NewPet{Name: "bar", Tag: "A"}}, {ID: 3, NewPet: NewPet{Name: "boo", Tag: "B"}}})

	m.RegisterFunc(AddPet).After(func(op *openapi3.Operation) {
		m.Doc.AddOperation("/pets", "POST", op)
		op.Tags = []string{"pet", "write"}
	})
	m.RegisterFunc(FindPetByID).After(func(op *openapi3.Operation) {
		m.Doc.AddOperation("/pets/{id}", "GET", op)
		op.Tags = []string{"pet", "read"}
	})
	m.RegisterFunc(DeletePet).After(func(op *openapi3.Operation) {
		m.Doc.AddOperation("/pets/{id}", "DELETE", op)
		op.Tags = []string{"pet", "write"}
	}).Status(204)
}
