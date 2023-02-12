package docgen

import (
	"context"
	"testing"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/google/go-cmp/cmp"
	reflectopenapi "github.com/podhmo/reflect-openapi"
	"github.com/podhmo/reflect-openapi/info"
)

type Person struct {
	Name   string  `json:"name" description:"name of person"`
	Age    string  `json:"age,omitempty"`
	Father *Person `json:"father"`

	Group    Group            `json:"group,omitempty"`
	Children []Person         `json:"children,omitempty"`
	Skills   map[string]Skill `json:"skills,omitempty"`
}

type Group struct {
	Name string
}

type Skill struct {
	Name        string
	Description string
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

	want := `// Person object
// - foo
// - bar
// - boo
type Person struct {
	// name of person
	name string
	age? string
	father? Person // :recursive:
	group? struct { // Group
		Name string
	}
	children? []Person // :recursive:
	skills? map[string]struct {     // Skill
		Name string
		Description string
	}
}
`
	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("TypeString() mismatch (-want +got):\n%s", diff)
	}
}
