package dochandler

import (
	"log"
	"net/http"
	"strings"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/podhmo/reflect-openapi/info"
)

func New(doc *openapi3.T, basePath string, info *info.Info, mdtext string) http.Handler {
	mux := &http.ServeMux{}
	basePath = strings.TrimSuffix(basePath, "/")

	for _, server := range doc.Servers {
		log.Printf("[INFO]  openapi-doc handler is registerd in GET %s%s/", server.URL, basePath)
	}

	mux.Handle(basePath+"/", ListEndpointHandler(
		doc,
		Endpoint{Method: "GET", Path: basePath + "/doc", OperationID: "OpenAPIDocHandler", Summary: "(added by github.com/podhmo/reflect-openapi/dochandler)"},
		Endpoint{Method: "GET", Path: basePath + "/ui", OperationID: "SwaggerUIHandler", Summary: "(added by github.com/podhmo/reflect-openapi/dochandler)"},
		Endpoint{Method: "GET", Path: basePath + "/redoc", OperationID: "RedocHandler", Summary: "(added by github.com/podhmo/reflect-openapi/dochandler)"},
		Endpoint{Method: "GET", Path: basePath + "/mddoc", OperationID: "MdDocHandler", Summary: "(added by github.com/podhmo/reflect-openapi/dochandler)"},
	))
	mux.Handle(basePath+"/doc", OpenAPIDocHandler(doc))
	mux.Handle(basePath+"/ui", SwaggerUIHandler(doc, basePath))
	mux.Handle(basePath+"/redoc", RedocHandler(doc, basePath))
	if info != nil {
		h := NewMdDocHandler(doc, info)
		if mdtext != "" {
			log.Printf("[INFO]  mddoc is loaded from a file.")
			h.text = mdtext
		}
		mux.HandleFunc(basePath+"/mddoc", h.HTML)
		mux.HandleFunc(basePath+"/mddoc.md", h.Text)
	}
	return mux
}
