package docgen

import (
	"context"
	"testing"

	"github.com/getkin/kin-openapi/openapi3"
	reflectopenapi "github.com/podhmo/reflect-openapi"
	"github.com/podhmo/reflect-openapi/info"
)

type Person struct {
	Name string `json:"name" description:"name of person"`
	Age  string `json:"age,omitempty"`
	// 	Father *Person `json:"father"`
	Group Group `json:"group,omitempty"`
}

type Group struct {
	Name string
}

func TestTypeString(t *testing.T) {
	c := &reflectopenapi.Config{SkipExtractComments: true, Info: info.New()}
	doc, err := c.BuildDoc(context.Background(), func(m *reflectopenapi.Manager) {
		m.RegisterType(Person{}).After(func(s *openapi3.Schema) {
			s.Description = "Person object\n- foo\n- bar\n- boo"
		})
	})
	if err != nil {
		t.Fatalf("unexpected setup failure: %+v", err)
	}

	name := "Person"
	schema := doc.Components.Schemas[name]
	got := TypeString(doc, c.Info, schema)
	t.Logf("name:%s\n```go\n%s```", name, got)
}
