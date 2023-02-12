package docgen

import (
	"bytes"
	"fmt"
	"io"
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

	indent := "" // TODO: nested
	schema := info.LookupSchema(ref)
	if description := schema.Description; description != "" {
		fmt.Fprintf(w, "%s// %s\n", indent, strings.Join(strings.Split(description, "\n"), fmt.Sprintf("\n%s// ", indent)))
	}
	fmt.Fprintf(w, "type %s ", schema.Title)
	writeType(w, doc, info, schema, nil)
	w.WriteRune('\n')
	return w.String()
}

func writeType(w *bytes.Buffer, doc *openapi3.T, info *info.Info, schema *openapi3.Schema, history []int) {
	switch schema.Type {
	case openapi3.TypeArray:
		writeArray(w, doc, info, schema, history)
	case openapi3.TypeBoolean:
		writeBoolean(w, doc, info, schema, history)
	case openapi3.TypeInteger:
		writeInteger(w, doc, info, schema, history)
	case openapi3.TypeNumber:
		writeNumber(w, doc, info, schema, history)
	case openapi3.TypeObject:
		isRecursive := false
		meta := info.Schemas[schema]
		for _, id := range history {
			if meta.ID == id {
				isRecursive = true
				break
			}
		}

		// TODO: map (additionalProperties)
		if isRecursive {
			fmt.Fprintf(w, "%s // :recursive:", schema.Title)
		} else {
			writeObject(w, doc, info, schema, history)
		}
	case openapi3.TypeString:
		writeString(w, doc, info, schema, history)
	default:
		panic(fmt.Sprintf("TypeString() unexpected schema type: %q", schema.Type))
	}
}
func writeArray(w *bytes.Buffer, doc *openapi3.T, info *info.Info, schema *openapi3.Schema, history []int) {
	io.WriteString(w, "[]")
	subschema := info.LookupSchema(schema.Items)
	writeType(w, doc, info, subschema, history)
}
func writeBoolean(w *bytes.Buffer, doc *openapi3.T, info *info.Info, schema *openapi3.Schema, history []int) {
	io.WriteString(w, "boolean")
}
func writeInteger(w *bytes.Buffer, doc *openapi3.T, info *info.Info, schema *openapi3.Schema, history []int) {
	io.WriteString(w, "integer")
}
func writeNumber(w *bytes.Buffer, doc *openapi3.T, info *info.Info, schema *openapi3.Schema, history []int) {
	io.WriteString(w, "number")
}
func writeObject(w *bytes.Buffer, doc *openapi3.T, info *info.Info, schema *openapi3.Schema, history []int) {
	io.WriteString(w, "struct {")
	if len(history) > 0 {
		fmt.Fprintf(w, "\t// %s", schema.Title)
	}
	w.WriteRune('\n')
	meta := info.Schemas[schema]
	for _, name := range meta.OrderedProperties {
		indent := strings.Repeat("\t", len(history)+1)

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

		fmt.Fprintf(w, "%s%s%s ", indent, name, suffix)

		subschema := info.LookupSchema(prop)
		writeType(w, doc, info, subschema, append(history, meta.ID))
		w.WriteRune('\n')
	}
	fmt.Fprintf(w, "%s}", strings.Repeat("\t", len(history)))
}

func writeString(w *bytes.Buffer, doc *openapi3.T, info *info.Info, schema *openapi3.Schema, history []int) {
	io.WriteString(w, "string")
}
