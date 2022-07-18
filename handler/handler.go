package handler

import (
	"log"
	"net/http"
	"strings"

	"github.com/getkin/kin-openapi/openapi3"
)

func NewHandler(doc *openapi3.T, basePath string) http.Handler {
	mux := &http.ServeMux{}
	basePath = strings.TrimSuffix(basePath, "/")

	for _, server := range doc.Servers {
		log.Printf("[INFO]  openapi-doc handler is registerd in GET %s%s/", server.URL, basePath)
	}

	mux.Handle(basePath+"/", ListEndpointHandler(
		doc,
		Endpoint{Method: "GET", Path: basePath + "/doc", OperationID: "OpenAPIDocHandler", Summary: "(added by github.com/podhmo/reflect-openapi/handler)"},
		Endpoint{Method: "GET", Path: basePath + "/ui", OperationID: "SwaggerUIHandler", Summary: "(added by github.com/podhmo/reflect-openapi/handler)"},
	))
	mux.Handle(basePath+"/doc", OpenAPIDocHandler(doc))
	mux.Handle(basePath+"/ui", SwaggerUIHandler(doc, basePath))
	return mux
}
