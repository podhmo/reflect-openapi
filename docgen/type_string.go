package docgen

import (
	"bytes"
	"fmt"
	"io"
	"reflect"
	"strconv"
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

	indent := ""
	schema := info.LookupSchema(ref)
	if description := schema.Description; description != "" {
		fmt.Fprintf(w, "%s// %s\n", indent, strings.Join(strings.Split(description, "\n"), fmt.Sprintf("\n%s// ", indent)))
	}
	fmt.Fprintf(w, "type %s ", schema.Title)
	writeType(w, doc, info, schema, nil)
	w.WriteRune('\n')
	return w.String()
}

// TODO: openapi-override

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
	case openapi3.TypeObject, "":
		// TODO: handling additionalProperties: true

		if ref := schema.AdditionalProperties.Schema; ref != nil { // map?
			writeMap(w, doc, info, schema, history)
			return
		}

		isRecursive := false
		meta := info.Schemas[schema]
		for _, id := range history {
			if meta.ID == id {
				isRecursive = true
				break
			}
		}

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
	// TODO: MinItems,MaxItems,UniqueItems
	if len(history) > 0 {
		if _, ok := schema.Extensions["x-go-type"]; ok {
			if len(history) > 0 {
				fmt.Fprintf(w, "[]%s", schema.Title)
				return
			}
		}
	}
	io.WriteString(w, "[]")
	subschema := info.LookupSchema(schema.Items)
	writeType(w, doc, info, subschema, history)
}
func writeBoolean(w *bytes.Buffer, doc *openapi3.T, info *info.Info, schema *openapi3.Schema, history []int) {
	if len(history) > 0 {
		if _, ok := schema.Extensions["x-go-type"]; ok {
			if len(history) > 0 {
				fmt.Fprintf(w, "%s[boolean]", schema.Title)
				return
			}
		}
	}

	io.WriteString(w, "boolean")
}
func writeInteger(w *bytes.Buffer, doc *openapi3.T, info *info.Info, schema *openapi3.Schema, history []int) {
	if len(history) > 0 {
		if _, ok := schema.Extensions["x-go-type"]; ok {
			if len(history) > 0 {
				fmt.Fprintf(w, "%s[integer]", schema.Title)
				return
			}
		}
	}

	io.WriteString(w, "integer")
}
func writeNumber(w *bytes.Buffer, doc *openapi3.T, info *info.Info, schema *openapi3.Schema, history []int) {
	// TODO: Min,Max,MultipleOf,exclusiveMinimum,exclusiveMaximum
	if len(history) > 0 {
		if _, ok := schema.Extensions["x-go-type"]; ok {
			if len(history) > 0 {
				fmt.Fprintf(w, "%s[number]", schema.Title)
				return
			}
		}
	}

	io.WriteString(w, "number")
}
func writeMap(w *bytes.Buffer, doc *openapi3.T, info *info.Info, schema *openapi3.Schema, history []int) {
	io.WriteString(w, "map[string]")
	subschema := info.LookupSchema(schema.AdditionalProperties.Schema)
	writeType(w, doc, info, subschema, history)
}
func writeObject(w *bytes.Buffer, doc *openapi3.T, info *info.Info, schema *openapi3.Schema, history []int) {
	// TODO: MinProps,MaxProps,(Discriminator)

	io.WriteString(w, "struct {")
	if len(history) > 0 {
		fmt.Fprintf(w, "\t// %s", schema.Title)
	}
	w.WriteRune('\n')
	meta := info.Schemas[schema]
	for _, name := range meta.OrderedProperties { // TODO: Nullable,Readonly,WriteOnly,AllowEmptyValue,Deprecated
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
	if len(history) > 0 {
		if _, ok := schema.Extensions["x-go-type"]; ok {
			if len(history) > 0 {
				fmt.Fprintf(w, "%s[string]", schema.Title)
				return
			}
		}
	}

	if len(schema.Enum) > 0 {
		values := make([]string, len(schema.Enum))
		if rt := reflect.TypeOf(schema.Enum[0]); rt.Kind() == reflect.String {
			for i, x := range schema.Enum {
				values[i] = strconv.Quote(reflect.ValueOf(x).String())
			}
		} else {
			for i, x := range schema.Enum {
				values[i] = fmt.Sprintf("%v", x)
			}
		}
		io.WriteString(w, strings.Join(values, " | "))
	} else {
		io.WriteString(w, "string")
	}

	{
		tags := make([]string, 0, 4)
		if schema.MinLength > 0 { // todo: openapi-override
			tags = append(tags, fmt.Sprintf(`minLength:"%d"`, schema.MinLength))
		}
		if schema.MaxLength != nil {
			tags = append(tags, fmt.Sprintf(`maxLength:"%d"`, *schema.MaxLength))
		}
		if schema.Pattern != "" {
			tags = append(tags, fmt.Sprintf(`pattern:"%s"`, schema.Pattern))
		}
		if len(tags) > 0 {
			fmt.Fprintf(w, " `%s`", strings.Join(tags, " "))
		}
	}

}
