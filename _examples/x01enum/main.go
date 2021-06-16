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
	ID string `json:"id" openapi:"path"`
}) *Todo {
	return nil
}

func ListTodo(params struct {
	Sort SortOrder `json:"sort" openapi:"query"`
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
		m.RegisterType(SortOrderAsc, func(schema *openapi3.Schema) {
			schema.Enum = []interface{}{
				SortOrderDesc,
				SortOrderAsc,
			}
		})
		m.RegisterFunc(GetTodo).After(func(op *openapi3.Operation) {
			m.Doc.AddOperation("/todo/{id}", "GET", op)
		})
	})
}
