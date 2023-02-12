package docgen

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"reflect"
	"strconv"
	"strings"
	"sync"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/podhmo/reflect-openapi/info"
)

var PADDING = `	`

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

	//  top level tags
	writeTags(w, info, schema, "// tags: ")
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
		log.Printf("[WARN]  TypeString() unexpected schema type: %q", schema.Type)
		fmt.Fprintf(w, "// TypeString() unexpected schema type: %q", schema.Type)
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
		fmt.Fprintf(w, "%s// %s", PADDING, schema.Title)
	}
	w.WriteRune('\n')
	meta := info.Schemas[schema]
	for _, name := range meta.OrderedProperties { // TODO: Nullable,Readonly,WriteOnly,AllowEmptyValue,Deprecated
		indent := strings.Repeat(PADDING, len(history)+1)

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
		writeTags(w, info, subschema, " ")
		w.WriteRune('\n')
	}
	fmt.Fprintf(w, "%s}", strings.Repeat(PADDING, len(history)))
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
}

func writeTags(w *bytes.Buffer, info *info.Info, schema *openapi3.Schema, prefix string) {
	tags := make([]string, 0, 20)

	switch schema.Type {
	case openapi3.TypeArray:
		tags = putTags(w, schema, tags)
		subschema := info.LookupSchema(schema.Items)
		tags = putTags(w, subschema, tags)
	case openapi3.TypeBoolean, openapi3.TypeInteger, openapi3.TypeNumber, openapi3.TypeString:
		tags = putTags(w, schema, tags)
	case openapi3.TypeObject, "":
		tags = putTags(w, schema, tags)
	default:
		log.Printf("[WARN]  writeTags() unexpected schema type: %q", schema.Type)
		return
	}

	if len(tags) > 0 {
		fmt.Fprintf(w, "%s`%s`", prefix, strings.Join(tags, " "))
	}
}

func putTags(w *bytes.Buffer, schema *openapi3.Schema, tags []string) []string {
	if schema.Format != "" {
		tags = append(tags, fmt.Sprintf(`format:"%s"`, schema.Format))
	}

	// properties
	if schema.Nullable {
		tags = append(tags, `nullable:"true"`)
	}
	if schema.ReadOnly {
		tags = append(tags, `readonly:"true"`)
	}
	if schema.WriteOnly {
		tags = append(tags, `writeonly:"true"`)
	}
	if schema.AllowEmptyValue {
		tags = append(tags, `allowEmptyValue:"true"`)
	}
	if schema.Deprecated {
		tags = append(tags, `deprecated:"true"`)
	}

	// array
	if schema.MinItems > 0 {
		tags = append(tags, fmt.Sprintf(`minItems:"%d"`, schema.MinItems))
	}
	if schema.MaxItems != nil {
		tags = append(tags, fmt.Sprintf(`maxItems:"%d"`, *schema.MaxItems))
	}
	if schema.UniqueItems {
		tags = append(tags, `uniqueItems:"true"`)
	}

	switch schema.Type {
	case openapi3.TypeInteger:
		if schema.Min != nil {
			tags = append(tags, fmt.Sprintf(`minimum:"%d"`, int64(*schema.Min)))
		}
		if schema.Max != nil {
			tags = append(tags, fmt.Sprintf(`maximum:"%d"`, int64(*schema.Max)))
		}
		if schema.MultipleOf != nil {
			tags = append(tags, fmt.Sprintf(`multipleOf:"%d"`, int64(*schema.MultipleOf)))
		}
	case openapi3.TypeNumber:
		if schema.Min != nil {
			tags = append(tags, fmt.Sprintf(`minimum:"%f"`, *schema.Min))
		}
		if schema.Max != nil {
			tags = append(tags, fmt.Sprintf(`maximum:"%f"`, *schema.Max))
		}
		if schema.MultipleOf != nil {
			tags = append(tags, fmt.Sprintf(`multipleOf:"%f"`, *schema.MultipleOf))
		}
	}
	if schema.ExclusiveMin {
		tags = append(tags, `exclusiveMinimum:"true"`)
	}
	if schema.ExclusiveMax {
		tags = append(tags, `exclusiveMaximum:"true"`)
	}

	// object
	if schema.MinProps > 0 {
		tags = append(tags, fmt.Sprintf(`minProps:"%d"`, schema.MinProps))
	}
	if schema.MaxProps != nil {
		tags = append(tags, fmt.Sprintf(`maxProps:"%d"`, *schema.MaxProps))
	}

	if schema.Discriminator != nil {
		tags = append(tags, fmt.Sprintf(`discriminator:"%s"`, schema.Discriminator.PropertyName))
	}

	// string
	if schema.MinLength > 0 {
		tags = append(tags, fmt.Sprintf(`minLength:"%d"`, schema.MinLength))
	}
	if schema.MaxLength != nil {
		tags = append(tags, fmt.Sprintf(`maxLength:"%d"`, *schema.MaxLength))
	}
	if schema.Pattern != "" {
		tags = append(tags, fmt.Sprintf(`pattern:"%s"`, schema.Pattern))
	}

	return tags
}
