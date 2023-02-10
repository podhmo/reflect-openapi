package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/getkin/kin-openapi/openapi3"
	reflectopenapi "github.com/podhmo/reflect-openapi"
	"github.com/podhmo/reflect-openapi/docgen"
)

func main() {
	if err := run(); err != nil {
		log.Fatalf("!! %+v", err)
	}
}

type FindPetsInput struct {
}

// Returns all pets
//
// Returns all pets from the system that the user has access to
// Nam sed condimentum est. Maecenas tempor sagittis sapien, nec rhoncus sem sagittis sit amet. Aenean at gravida augue, ac iaculis sem. Curabitur odio lorem, ornare eget elementum nec, cursus id lectus. Duis mi turpis, pulvinar ac eros ac, tincidunt varius justo. In hac habitasse platea dictumst. Integer at adipiscing ante, a sagittis ligula. Aenean pharetra tempor ante molestie imperdiet. Vivamus id aliquam diam. Cras quis velit non tortor eleifend sagittis. Praesent at enim pharetra urna volutpat venenatis eget eget mauris. In eleifend fermentum facilisis. Praesent enim enim, gravida ac sodales sed, placerat id erat. Suspendisse lacus dolor, consectetur non augue vel, vehicula interdum libero. Morbi euismod sagittis libero sed lacinia.
//
// Sed tempus felis lobortis leo pulvinar rutrum. Nam mattis velit nisl, eu condimentum ligula luctus nec. Phasellus semper velit eget aliquet faucibus. In a mattis elit. Phasellus vel urna viverra, condimentum lorem id, rhoncus nibh. Ut pellentesque posuere elementum. Sed a varius odio. Morbi rhoncus ligula libero, vel eleifend nunc tristique vitae. Fusce et sem dui. Aenean nec scelerisque tortor. Fusce malesuada accumsan magna vel tempus. Quisque mollis felis eu dolor tristique, sit amet auctor felis gravida. Sed libero lorem, molestie sed nisl in, accumsan tempor nisi. Fusce sollicitudin massa ut lacinia mattis. Sed vel eleifend lorem. Pellentesque vitae felis pretium, pulvinar elit eu, euismod sapien.
func FindPets(input FindPetsInput) {}

type AddPetInput struct {
}

// Creates a new pet
//
// Creates a new pet in the store. Duplicates are allowed
func AddPet(input AddPetInput) {}

type FindPetByIDInput struct {
	ID string `in:"path" path:"id"`
}

// Returns a pet by ID
//
// Returns a pet based on a single ID
func FindPetByID(input FindPetByIDInput) {}

type DeletePetInput struct {
	ID string `in:"path" path:"id"`
}

// Deletes a pet by ID
//
// deletes a single pet based on the ID supplied
func DeletePet(input DeletePetInput) {}

func run() error {
	c := &reflectopenapi.Config{}
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

	doc := docgen.Generate(tree)
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
