package handler

import (
	"net/http"
	"strings"

	"github.com/getkin/kin-openapi/openapi3"
)

func NewHandler(doc *openapi3.Swagger, basePath string) http.Handler {
	mux := &http.ServeMux{}
	basePath = strings.TrimSuffix(basePath, "/")

	mux.Handle(basePath+"/", ListEndpointHandler(
		doc,
		Endpoint{Method: "GET", Path: basePath + "/doc", OperationID: "OpenAPIDocHandler", Summary: "(added by github.com/podhmo/reflect-openapi/handler)"},
		Endpoint{Method: "GET", Path: basePath + "/ui", OperationID: "SwaggerUIHandler", Summary: "(added by github.com/podhmo/reflect-openapi/handler)"},
	))
	mux.Handle(basePath+"/doc", OpenAPIDocHandler(doc))
	mux.Handle(basePath+"/ui", SwaggerUIHandler(doc, basePath))
	return mux
}
