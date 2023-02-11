package reflectopenapi

import (
	"reflect"
	"testing"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/google/go-cmp/cmp"
	"github.com/podhmo/reflect-openapi/internal"
	shape "github.com/podhmo/reflect-shape"
)

func TestInfoOrderedProperties(t *testing.T) {
	newTransformer := func(info *internal.Info) *Transformer {
		c := Config{SkipExtractComments: true}
		transformer := (&Transformer{
			TagNameOption:    *DefaultTagNameOption(),
			cache:            map[int]interface{}{},
			defaultValues:    map[int]reflect.Value{},
			interceptFuncMap: map[reflect.Type]func(*shape.Shape) *openapi3.Schema{},
			Resolver:         c.DefaultResolver(),
			Selector:         c.DefaultSelector(),
			Extractor:        c.DefaultExtractor(),
		}).Builtin()
		transformer.IsRequired = transformer.isRequired
		if t, ok := transformer.Selector.(needTransformer); ok {
			t.NeedTransformer(transformer)
		}
		transformer.info = info
		return transformer
	}

	type ref struct{ VS []string }
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

	info := internal.NewInfo()
	transformer := newTransformer(info)

	for _, c := range cases {
		c := c
		t.Run(c.msg, func(t *testing.T) {

			got := transformer.Transform(transformer.Extractor.Extract(c.input))
			schema, ok := got.(*openapi3.Schema)
			if !ok {
				t.Errorf("Transformer.Transform() with info: unexpected type %T", schema)
			}
			sinfo, ok := info.Schemas[schema]
			if !ok {
				t.Errorf("Transformer.Transform() with info: SchemaInfo is not found")
			}

			if diff := cmp.Diff(ref{c.wantProperties}, ref{sinfo.OrderedProperties}); diff != "" {
				t.Errorf("Info.OrderedPropreties mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
