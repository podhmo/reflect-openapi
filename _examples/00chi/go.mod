module m

go 1.15

replace github.com/podhmo/reflect-openapi => ../../

require (
	github.com/getkin/kin-openapi v0.83.0
	github.com/go-chi/chi v4.1.2+incompatible
	github.com/go-chi/render v1.0.1
	github.com/podhmo/reflect-openapi v0.0.12
	github.com/podhmo/reflect-shape v0.3.5 // indirect
	golang.org/x/net v0.0.0-20210917221730-978cfadd31cf // indirect
)
