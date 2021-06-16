package handler

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/getkin/kin-openapi/openapi3"
)

func OpenAPIDocHandler(doc *openapi3.T) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		enc := json.NewEncoder(w)
		if err := enc.Encode(doc); err != nil {
			fmt.Fprintf(w, `{"error": %q}`, err.Error())
			return
		}
	}
}
