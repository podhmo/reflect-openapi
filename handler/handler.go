package handler

import (
	"net/http"

	"github.com/getkin/kin-openapi/openapi3"
)

func NewHandler(doc *openapi3.Swagger) http.Handler {
	mux := &http.ServeMux{}
	mux.Handle("/", EndpointsHandler(doc))
	mux.Handle("/swagger-ui", SwaggerUIHandler(doc))
	return mux
}
