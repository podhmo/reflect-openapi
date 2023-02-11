package reflectopenapi

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/podhmo/reflect-openapi/info"
)

func TestInfoOrderedProperties(t *testing.T) {
	newVisitor := func(info *info.Info) *Visitor {
		c := Config{SkipExtractComments: true}
		visitor := NewVisitor(*c.TagNameOption, c.DefaultResolver(), c.DefaultSelector(), c.DefaultExtractor())
		visitor.Transformer.info = info
		return visitor
	}

	type w struct{ VS []string }
	type Input struct {
		X, Y, Z string
		Value   int
	}

	type InputWithJSONTag struct {
		X     string `json:"x,omitempty"`
		Y     string `json:"y,omitempty"`
		Z     string `json:"z,omitempty"`
		Value int    `json:"value,omitempty"`
	}

	type InputWithIgnore struct {
		X     string `json:"x,omitempty"`
		Y     string `json:"-,omitempty"`
		Z     string `json:"z,omitempty"`
		Value int    `json:"value,omitempty"`
	}

	type InputWithEmbedded_Embedded_Embedded struct {
		Z string `json:"z,omitempty"`
	}
	type InputWithEmbedded_Embedded struct {
		Y string `json:"y,omitempty"`
		*InputWithEmbedded_Embedded_Embedded
	}
	type InputWithEmbedded struct {
		X string `json:"x,omitempty"`
		InputWithEmbedded_Embedded
		Value int `json:"value,omitempty"`
	}

	type InputWithNested struct {
		X      string `json:"x,omitempty"`
		Nested struct {
			Y      string `json:"y,omitempty"`
			Nested struct {
				Z string `json:"z,omitempty"`
			} `json:"nested"`
		} `json:"nested"`
		Value int `json:"value,omitempty"`
	}

	cases := []struct {
		msg            string
		input          interface{}
		wantProperties []string
	}{
		{"simple", Input{}, []string{"X", "Y", "Z", "Value"}},
		{"with-json-tag", InputWithJSONTag{}, []string{"x", "y", "z", "value"}},
		{"with-ignore", InputWithIgnore{}, []string{"x", "z", "value"}},
		{"with-embedded", InputWithEmbedded{}, []string{"x", "y", "z", "value"}},
		{"with-nested", InputWithNested{}, []string{"x", "nested", "value"}},
	}

	info := info.New()
	visitor := newVisitor(info)

	for _, c := range cases {
		c := c
		t.Run(c.msg, func(t *testing.T) {
			ref := visitor.VisitType(visitor.Extractor.Extract(c.input))
			schema := info.LookupSchema(ref)
			if schema == nil {
				t.Fatalf("Transformer.Transform() with info: schema is not found")
			}
			sinfo, ok := info.Schemas[schema]
			if !ok {
				t.Errorf("Transformer.Transform() with info: SchemaInfo is not found")
			}

			if diff := cmp.Diff(w{c.wantProperties}, w{sinfo.OrderedProperties}); diff != "" {
				t.Errorf("Info.OrderedPropreties mismatch (-want +got):\n%s", diff)
			}
		})
	}

}
