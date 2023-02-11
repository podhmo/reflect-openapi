package internal

import "github.com/getkin/kin-openapi/openapi3"

// Info is the go/types.Info like object that handling metadata.
type Info struct {
	Schemas map[*openapi3.Schema]SchemaInfo
}

func NewInfo() *Info {
	return &Info{
		Schemas: map[*openapi3.Schema]SchemaInfo{},
	}
}

func (i *Info) RegisterSchemaInfo(schema *openapi3.Schema, id int, props []string) {
	i.Schemas[schema] = SchemaInfo{ID: id, OrderedProperties: props}
}

type SchemaInfo struct {
	ID int // reflectshape.Schema.Number

	OrderedProperties []string
}
