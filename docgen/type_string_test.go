package docgen

import (
	"context"
	"testing"

	reflectopenapi "github.com/podhmo/reflect-openapi"
	"github.com/podhmo/reflect-openapi/info"
)

type Person struct {
	Name   string  `json:"name"`
	Age    string  `json:"age,omitempty"`
	Father *Person `json:"father"`
}

func TestTypeString(t *testing.T) {
	c := &reflectopenapi.Config{SkipExtractComments: true, Info: info.New()}
	doc, err := c.BuildDoc(context.Background(), func(m *reflectopenapi.Manager) {
		m.RegisterType(Person{})
	})
	if err != nil {
		t.Fatalf("unexpected setup failure: %+v", err)
	}

	name := "Person"
	schema := doc.Components.Schemas[name]
	got := TypeString(doc, c.Info, schema)
	t.Logf("name:%s\n```go\n%s```", name, got)
}
