package main

import (
	"os"
	"strconv"
	"time"

	"github.com/getkin/kin-openapi/openapi3"
	reflectopenapi "github.com/podhmo/reflect-openapi"
)

type Todo struct {
	ID        string    `json:"id" required:"true"`
	Title     string    `json:"title"`
	Done      bool      `json:"done"`
	CreatedAt time.Time `json:"createdAt"`
}

type SortOrder string

const (
	SortOrderDesc SortOrder = "desc"
	SortOrderAsc  SortOrder = "asc"
)

func GetTodo(params struct {
	ID string `json:"id" in:"path"`
}) *Todo {
	return nil
}

func ListTodo(params struct {
	Sort SortOrder `json:"sort" in:"query"`
}) []Todo {
	return nil
}

func main() {
	c := reflectopenapi.Config{
		SkipValidation: true,
	}
	if ok, _ := strconv.ParseBool(os.Getenv("WITHOUT_REF")); ok {
		c.Resolver = &reflectopenapi.NoRefResolver{}
	}

	c.EmitDoc(func(m *reflectopenapi.Manager) {
		m.RegisterFunc(ListTodo).After(func(op *openapi3.Operation) {
			m.Doc.AddOperation("/todo", "GET", op)
		})
		m.RegisterType(SortOrderAsc).Enum(SortOrderDesc, SortOrderAsc)
		m.RegisterFunc(GetTodo).After(func(op *openapi3.Operation) {
			m.Doc.AddOperation("/todo/{id}", "GET", op)
		})
	})
}
