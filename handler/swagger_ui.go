package handler

import (
	"fmt"
	"net/http"

	"github.com/getkin/kin-openapi/openapi3"
)

func SwaggerUIHandler(doc *openapi3.Swagger, basePath string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, TEMPLATE, basePath+"/doc")
	}
}

const TEMPLATE = `<!DOCTYPE html>
<html>

<head>
    <link type="text/css" rel="stylesheet" href="https://cdn.jsdelivr.net/npm/swagger-ui-dist@3.30.0/swagger-ui.css">
    <title>OpenAPI Docs</title>
</head>

<body>
    <div id="swagger-ui"></div>
    <script src="https://cdn.jsdelivr.net/npm/swagger-ui-dist@3.30.0/swagger-ui-bundle.js"></script>
    <script>
        const ui = SwaggerUIBundle({
            url: '%s', // the endpoint returns openAPI doc
            dom_id: '#swagger-ui',
            presets: [
                SwaggerUIBundle.presets.apis,
                SwaggerUIBundle.SwaggerUIStandalonePreset
            ],
            layout: "BaseLayout",
            deepLinking: true,
            showExtensions: true,
            showCommonExtensions: true
        })

    </script>
</body>

</html>
`
