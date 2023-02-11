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
