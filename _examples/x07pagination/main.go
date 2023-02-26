package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/getkin/kin-openapi/openapi3"
	reflectopenapi "github.com/podhmo/reflect-openapi"
	"github.com/podhmo/reflect-openapi/docgen"
	"github.com/podhmo/reflect-openapi/info"
)

func main() {
	if err := run(); err != nil {
		log.Fatalf("!! %+v", err)
	}
}

func run() error {
	c := &reflectopenapi.Config{Info: info.New()}
	ctx := context.Background()
	doc, err := c.BuildDoc(ctx, func(m *reflectopenapi.Manager) {
		m.RegisterFunc(ListUser, func(op *openapi3.Operation) {
			m.Doc.AddOperation("/users", "GET", op)
		})
	})
	if err != nil {
		return fmt.Errorf("build doc: %w", err)
	}

	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	if err := enc.Encode(doc); err != nil {
		return fmt.Errorf("write doc: %w", err)
	}

	md := docgen.Generate(doc, c.Info)
	f, err := os.Create("README.md")
	if err != nil {
		return fmt.Errorf("new md: %w", err)
	}
	defer f.Close()
	if err := docgen.WriteDoc(f, md); err != nil {
		return fmt.Errorf("write md: %w", err)
	}
	return nil
}

type PaginatedInput[T any] struct {
	Cursor   string `in:"query" query:"cursor"`
	PageSize int    `in:"query" query:"pageSize" openapi-override:"{'default': 20, 'maximum': 100}"`

	Value T `embedded:"true"` // TODO: need embedded in generics
}

type PaginatedOutput[T any] struct {
	HasMore    bool   `json:"hasMore"`
	Cursor     string `json:"cursor"`
	NextCursor string `json:"nextCursor"`

	Items T `json:"items"`
}

type User struct {
	Name string `json:"name"`
}

type ListUserInput struct {
	Sort  string `in:"query" query:"sort" openapi-override:"{'enum': ['asc', 'desc'], 'default': 'desc'}"`
	Query string `in:"query" query:"query"`
}

func ListUser(input PaginatedInput[ListUserInput]) (*PaginatedOutput[[]User], error) {
	return nil, nil
}
