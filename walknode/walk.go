package walknode

import (
	"sort"

	"github.com/getkin/kin-openapi/openapi3"
)

func PathItem(doc *openapi3.T, fn func(item *openapi3.PathItem, path string)) {
	paths := make([]string, 0, len(doc.Paths))
	for k := range doc.Paths {
		paths = append(paths, k)
	}
	sort.Strings(paths)

	for _, path := range paths {
		fn(doc.Paths[path], path)
	}
}

func Operation(pathItem *openapi3.PathItem, fn func(op *openapi3.Operation, method string)) {
	if op := pathItem.Connect; op != nil {
		method := "CONNECT"
		fn(op, method)
	}
	if op := pathItem.Delete; op != nil {
		method := "DELETE"
		fn(op, method)
	}
	if op := pathItem.Get; op != nil {
		method := "GET"
		fn(op, method)
	}
	if op := pathItem.Head; op != nil {
		method := "HEAD"
		fn(op, method)
	}
	if op := pathItem.Options; op != nil {
		method := "OPTIONS"
		fn(op, method)
	}
	if op := pathItem.Patch; op != nil {
		method := "PATCH"
		fn(op, method)
	}
	if op := pathItem.Post; op != nil {
		method := "POST"
		fn(op, method)
	}
	if op := pathItem.Put; op != nil {
		method := "PUT"
		fn(op, method)
	}
	if op := pathItem.Trace; op != nil {
		method := "TRACE"
		fn(op, method)
	}
}

func Schema(doc *openapi3.T, fn func(*openapi3.SchemaRef, string)) {
	schemas := make([]string, 0, len(doc.Components.Schemas))
	for k := range doc.Components.Schemas {
		schemas = append(schemas, k)
	}
	sort.Strings(schemas)
	for _, k := range schemas {
		fn(doc.Components.Schemas[k], k)
	}
}

func Response(op *openapi3.Operation, fn func(*openapi3.ResponseRef, string)) {
	responses := make([]string, 0, len(op.Responses))
	for k := range op.Responses {
		responses = append(responses, k)
	}
	sort.Strings(responses)
	for _, k := range responses {
		fn(op.Responses[k], k)
	}
}

func Example(examples openapi3.Examples, fn func(*openapi3.ExampleRef, string)) {
	if ref, ok := examples["default"]; ok {
		fn(ref, "default")
	}
	keys := make([]string, 0, len(examples))
	for k := range examples {
		if k == "default" {
			continue
		}
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		fn(examples[k], k)
	}
}
