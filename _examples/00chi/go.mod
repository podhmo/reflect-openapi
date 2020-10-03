module m

go 1.15

replace github.com/podhmo/reflect-openapi => ../../

require (
	github.com/getkin/kin-openapi v0.22.1
	github.com/go-chi/chi v4.1.2+incompatible
	github.com/go-chi/render v1.0.1
	github.com/podhmo/reflect-openapi v0.0.0-00010101000000-000000000000
	golang.org/x/net v0.0.0-20201002202402-0a1ea396d57c // indirect
)
