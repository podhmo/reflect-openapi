package docgen

import (
	"context"
	"testing"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/google/go-cmp/cmp"
	reflectopenapi "github.com/podhmo/reflect-openapi"
	"github.com/podhmo/reflect-openapi/info"
)

type Sort string

const (
	SortASC  Sort = "asc"
	SortDESC Sort = "desc"
)

type Person struct {
	Name   string      `json:"name" description:"name of person"`
	Age    PositiveInt `json:"age,omitempty" openapi-override:"{'exclusiveMinimum': true}"`
	Father *Person     `json:"father"`

	Group    Group            `json:"group,omitempty"`
	Children []Person         `json:"children,omitempty"`
	Skills   map[string]Skill `json:"skills,omitempty"`

	Sort Sort `json:"sort"`
}

type PositiveInt int

type Group struct {
	Name string
}

type Skill struct {
	Name        string
	Description string
}

func TestTypeString(t *testing.T) {
	PADDING = "@@"
	defer func() { PADDING = "\t" }()

	c := &reflectopenapi.Config{SkipExtractComments: true, Info: info.New()}
	doc, err := c.BuildDoc(context.Background(), func(m *reflectopenapi.Manager) {
		m.RegisterType(PositiveInt(0)).After(func(s *openapi3.Schema) { var z float64; s.Min = &z })
		m.RegisterType(SortASC).Enum(SortASC, SortDESC)
		m.RegisterType(Person{}).After(func(s *openapi3.Schema) {
			s.Description = `Person object
- foo
- bar
- boo`
			s.Properties["name"].Value.Pattern = `^[A-Z][a-zA-z\-_]+$`
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
@@// name of person
@@name string ` + "`pattern:\"^[A-Z][a-zA-z\\-_]+$\"`" + `
@@age? PositiveInt[integer] ` + "`minimum:\"0\" exclusiveMinimum:\"true\"`" + `
@@father? Person // :recursive:
@@group? struct {@@// Group
@@@@Name string
@@}
@@children? []Person // :recursive:
@@skills? map[string]struct {@@// Skill
@@@@Name string
@@@@Description string
@@}
@@sort Sort[string]
}
`
	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("TypeString() mismatch (-want +got):\n%s", diff)
	}

	// _ = want
	// t.Logf("%s", got)
}
