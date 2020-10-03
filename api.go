package reflectopenapi

import "github.com/getkin/kin-openapi/openapi3"

func NewDoc() (*openapi3.Swagger, error) {
	skeleton := []byte(`{
  "openapi": "3.0.0",
  "info": {
    "title": "Sample API",
    "description": "-",
    "version": "0.0.0"
  },
  "servers": [
    {
      "url": "http://localhost:8888",
      "description": "local development server"
    },
  ],
}`)
	return NewDocFromSkeleton(skeleton)
}

// TODO: add api function
func NewDocFromSkeleton(skeleton []byte) (*openapi3.Swagger, error) {
	l := openapi3.NewSwaggerLoader()
	return l.LoadSwaggerFromData(skeleton)
}
