module m

go 1.15

replace github.com/podhmo/reflect-openapi => ../../

require (
	github.com/getkin/kin-openapi v0.75.0
	github.com/go-chi/chi v4.1.2+incompatible
	github.com/go-chi/render v1.0.1
	github.com/podhmo/reflect-openapi v0.0.11
	golang.org/x/net v0.0.0-20210917221730-978cfadd31cf // indirect
)
