package docgen

import (
	"bytes"
	"fmt"
	"strings"
	"sync"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/podhmo/reflect-openapi/info"
)

var pool = &sync.Pool{
	New: func() interface{} { return new(bytes.Buffer) },
}

func TypeString(doc *openapi3.T, info *info.Info, ref *openapi3.SchemaRef) string {
	w := pool.Get().(*bytes.Buffer)
	defer pool.Put(w)
	w.Reset()

	schema := info.LookupSchema(ref)
	switch schema.Type {
	case openapi3.TypeArray:
		writeArray(w, doc, info, schema)
	case openapi3.TypeBoolean:
		writeBoolean(w, doc, info, schema)
	case openapi3.TypeInteger:
		writeInteger(w, doc, info, schema)
	case openapi3.TypeNumber:
		writeNumber(w, doc, info, schema)
	case openapi3.TypeObject:
		writeObject(w, doc, info, schema)
	case openapi3.TypeString:
		writeString(w, doc, info, schema)
	default:
		panic(fmt.Sprintf("TypeString() unexpected schema type: %q", schema.Type))
	}
	return w.String()
}

func writeArray(w *bytes.Buffer, doc *openapi3.T, info *info.Info, schema *openapi3.Schema) {
}
func writeBoolean(w *bytes.Buffer, doc *openapi3.T, info *info.Info, schema *openapi3.Schema) {
}
func writeInteger(w *bytes.Buffer, doc *openapi3.T, info *info.Info, schema *openapi3.Schema) {
}
func writeNumber(w *bytes.Buffer, doc *openapi3.T, info *info.Info, schema *openapi3.Schema) {
}
func writeObject(w *bytes.Buffer, doc *openapi3.T, info *info.Info, schema *openapi3.Schema) {
	indent := "" // TODO: nested
	if description := schema.Description; description != "" {
		fmt.Fprintf(w, "%s// %s\n", indent, strings.Join(strings.Split(description, "\n"), fmt.Sprintf("\n%s// ", indent)))
	}
	fmt.Fprintf(w, "type %s struct {\n", schema.Title)
	for _, name := range info.Schemas[schema].OrderedProperties {
		indent := "\t" // TODO: nested

		prop := schema.Properties[name]

		suffix := "?"
		for _, x := range schema.Required {
			if name == x {
				suffix = ""
				break
			}
		} // or nullable?

		if description := prop.Value.Description; description != "" && prop.Value.Type != "object" {
			fmt.Fprintf(w, "%s// %s\n", indent, strings.Join(strings.Split(description, "\n"), fmt.Sprintf("\n%s// ", indent)))
		}

		fmt.Fprintf(w, "%s%s%s %s", indent, name, suffix, prop.Value.Type)
		w.WriteRune('\n')
	}
	fmt.Fprintln(w, "}")
}

func writeString(w *bytes.Buffer, doc *openapi3.T, info *info.Info, schema *openapi3.Schema) {
}
