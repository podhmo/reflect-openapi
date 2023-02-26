package main

import (
	"github.com/getkin/kin-openapi/openapi3"
	reflectopenapi "github.com/podhmo/reflect-openapi"
)

func main() {
	c := &reflectopenapi.Config{}
	c.EmitDoc(func(m *reflectopenapi.Manager) {
		m.RegisterFunc(ListUser, func(op *openapi3.Operation) {
			m.Doc.AddOperation("/users", "GET", op)
		})
	})
}

type PaginatedInput[T any] struct {
	Cursor   string `in:"query" query:"cursor"`
	PageSize int    `in:"query" query:"pageSize"`

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
