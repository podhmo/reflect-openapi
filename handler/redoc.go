package handler

import (
	"fmt"
	"net/http"

	"github.com/getkin/kin-openapi/openapi3"
)

func RedocHandler(doc *openapi3.T, basePath string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		title := fmt.Sprintf("%s (%s)", doc.Info.Title, doc.Info.Version)
		fmt.Fprintf(w, REDOC_TEMPLATE, title, basePath+"/doc")
	}
}

const REDOC_TEMPLATE = `<!DOCTYPE html>
<html>
<head>
<title>%s</title>
<!-- needed for adaptive design -->
<meta charset="utf-8"/>
<meta name="viewport" content="width=device-width, initial-scale=1"><!DOCTYPE html>

<link href="https://fonts.googleapis.com/css?family=Montserrat:300,400,700|Roboto:300,400,700" rel="stylesheet">

<!--
ReDoc doesn't change outer page styles
-->
<style>
  body {{
	margin: 0;
	padding: 0;
  }}
</style>
</head>
<body>
<redoc spec-url="%s"></redoc>
<script src="https://cdn.jsdelivr.net/npm/redoc@next/bundles/redoc.standalone.js"> </script>
</body>
</html>
`
