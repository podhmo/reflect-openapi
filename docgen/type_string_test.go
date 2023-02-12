package docgen

import (
	"context"
	"strings"
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

type HelloInput struct {
	Pretty bool   `in:"query" query:"pretty"`
	Name   string `json:"name"`
}
type HelloOutput struct {
	Message string `json:"message"`
}
type Error struct {
	Code    int64  `json:"code"`
	Message string `json:"message"`
}

func Hello(HelloInput) *HelloOutput { return nil }

func TestActionInputString(t *testing.T) {
	PADDING = "@@"
	defer func() { PADDING = "\t" }()

	c := &reflectopenapi.Config{SkipExtractComments: true, Info: info.New()}
	doc, err := c.BuildDoc(context.Background(), func(m *reflectopenapi.Manager) {
		m.RegisterFunc(Hello).After(func(op *openapi3.Operation) {
			m.Doc.AddOperation("/Hello", "POST", op)
			op.Parameters[0].Value.Description = "if true\npretty print is activate"
		})
	})
	if err != nil {
		t.Fatalf("unexpected setup failure: %+v", err)
	}

	op := doc.Paths.Find("/Hello").GetOperation("POST")
	got := ActionInputString(doc, c.Info, op)

	want := `
type Input struct {
@@// if true
@@// pretty print is activate
@@pretty? boolean ` + "`in:\"query\"`" + `

@@JSONBody struct {@@// HelloInput
@@@@name string
@@}
}
`

	if diff := cmp.Diff(strings.TrimSpace(want), strings.TrimSpace(got)); diff != "" {
		t.Errorf("ActionOutputString() mismatch (-want +got):\n%s", diff)
	}
	// t.Logf("%s", got)
}

func TestActionOutputString(t *testing.T) {
	PADDING = "@@"
	defer func() { PADDING = "\t" }()

	info := info.New()
	c := &reflectopenapi.Config{
		SkipExtractComments: true,
		Info:                info,
		DefaultError:        Error{},
	}
	doc, err := c.BuildDoc(context.Background(), func(m *reflectopenapi.Manager) {
		m.RegisterFunc(Hello).After(func(op *openapi3.Operation) {
			m.Doc.AddOperation("/Hello", "POST", op)
		})
	})
	if err != nil {
		t.Fatalf("unexpected setup failure: %+v", err)
	}

	op := doc.Paths.Find("/Hello").GetOperation("POST")

	cases := []struct {
		name string
		want string
	}{
		{name: "default", want: `
// default error
type OutputDefault struct {
@@code integer ` + "`format:\"int64\"`" + `
@@message string
}`},
		{name: "200", want: `
type Output200 struct {
@@message string
}`},
	}

	for _, c := range cases {
		c := c
		name := c.name
		t.Run(name, func(t *testing.T) {
			ref := op.Responses[name]
			got := ActionOutputString(doc, info, ref, name)
			if diff := cmp.Diff(strings.TrimSpace(c.want), strings.TrimSpace(got)); diff != "" {
				t.Errorf("ActionOutputString() mismatch (-want +got):\n%s", diff)
			}
			// t.Logf("%s", got)
		})
	}
}
