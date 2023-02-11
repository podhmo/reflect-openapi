package internal

import "github.com/getkin/kin-openapi/openapi3"

// Info is the go/types.Info like object that handling metadata.
type Info struct {
	Schemas map[*openapi3.Schema]SchemaInfo
	Refs    map[*openapi3.SchemaRef]*openapi3.Schema
}

func NewInfo() *Info {
	return &Info{
		Schemas: map[*openapi3.Schema]SchemaInfo{},
		Refs:    map[*openapi3.SchemaRef]*openapi3.Schema{},
	}
}

func (i *Info) RegisterSchemaInfo(schema *openapi3.Schema, id int, props []string) {
	i.Schemas[schema] = SchemaInfo{ID: id, OrderedProperties: props}
}

func (i *Info) RegisterRef(ref *openapi3.SchemaRef, schema *openapi3.Schema) {
	i.Refs[ref] = schema
}

func (i *Info) LookupSchema(ref *openapi3.SchemaRef) *openapi3.Schema {
	if ref.Value != nil {
		return ref.Value
	}
	if v, ok := i.Refs[ref]; ok {
		return v
	}
	return nil
}

type SchemaInfo struct {
	ID int // reflectshape.Schema.Number

	OrderedProperties []string
}
