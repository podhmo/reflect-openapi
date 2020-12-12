package handler

import (
	"net/http"

	"github.com/getkin/kin-openapi/openapi3"
)

func SwaggerUIHandler(doc *openapi3.Swagger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
	}
}
