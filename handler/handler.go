package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sort"

	"github.com/getkin/kin-openapi/openapi3"
)

func NewHandler(doc *openapi3.Swagger) http.Handler {
	mux := &http.ServeMux{}
	mux.Handle("/", EndpointsHandler(doc))
	return mux
}

func EndpointsHandler(doc *openapi3.Swagger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		list := endpoints(doc)
		w.Header().Add("Content-type", "application/json")
		enc := json.NewEncoder(w)
		if err := enc.Encode(list); err != nil {
			fmt.Fprintf(w, `{"error": %q}`, err.Error())
			return
		}
	}
}

type Endpoint struct {
	Method      string `json:"method"`
	Path        string `json:"path"`
	OperationID string `json:"operationId"`
	Summary     string `json:"summary"`
}

func (r *Endpoint) String() string {
	return fmt.Sprintf("%s %s	%s	%s\n", r.Method, r.Path, r.OperationID, r.Summary)
}

func endpoints(doc *openapi3.Swagger) []Endpoint {
	var r []Endpoint
	for path, pathItem := range doc.Paths {
		if pathItem.Get != nil {
			x := Endpoint{Path: path, Method: "GET", OperationID: pathItem.Get.OperationID, Summary: pathItem.Get.Summary}
			r = append(r, x)
		}
		if pathItem.Post != nil {
			x := Endpoint{Path: path, Method: "POST", OperationID: pathItem.Post.OperationID, Summary: pathItem.Post.Summary}
			r = append(r, x)
		}
		if pathItem.Patch != nil {
			x := Endpoint{Path: path, Method: "PATCH", OperationID: pathItem.Patch.OperationID, Summary: pathItem.Patch.Summary}
			r = append(r, x)
		}
		if pathItem.Put != nil {
			x := Endpoint{Path: path, Method: "PUT", OperationID: pathItem.Put.OperationID, Summary: pathItem.Put.Summary}
			r = append(r, x)
		}
		if pathItem.Head != nil {
			x := Endpoint{Path: path, Method: "HEAD", OperationID: pathItem.Head.OperationID, Summary: pathItem.Head.Summary}
			r = append(r, x)
		}
		if pathItem.Options != nil {
			x := Endpoint{Path: path, Method: "OPTIONS", OperationID: pathItem.Options.OperationID, Summary: pathItem.Options.Summary}
			r = append(r, x)
		}
		if pathItem.Trace != nil {
			x := Endpoint{Path: path, Method: "TRACE", OperationID: pathItem.Trace.OperationID, Summary: pathItem.Trace.Summary}
			r = append(r, x)
		}
	}
	sort.SliceStable(r, func(i, j int) bool { return r[i].Path < r[j].Path })
	return r
}
