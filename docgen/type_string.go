package docgen

import (
	"bytes"
	"fmt"
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

	// Array
	// String
	// Number

	// Object
	schema := info.LookupSchema(ref)
	fmt.Fprintf(w, "type %s struct {\n", schema.Title)
	for _, name := range info.Schemas[schema].OrderedProperties {
		prop := schema.Properties[name]

		suffix := "?"
		for _, x := range schema.Required {
			if name == x {
				suffix = ""
				break
			}
		} // or nullable?

		fmt.Fprintf(w, "\t%s%s %s\n", name, suffix, prop.Value.Type)
	}
	fmt.Fprintln(w, "}")
	return w.String()
}
