package info

import "github.com/getkin/kin-openapi/openapi3"

// Info is the go/types.Info like object that handling metadata.
type Info struct {
	SchemaInfo  map[*openapi3.Schema]*SchemaInfo
	SchemaValue map[*openapi3.SchemaRef]*openapi3.Schema
}

func New() *Info {
	return &Info{
		SchemaInfo:  map[*openapi3.Schema]*SchemaInfo{},
		SchemaValue: map[*openapi3.SchemaRef]*openapi3.Schema{},
	}
}

func (i *Info) LookupSchema(ref *openapi3.SchemaRef) *openapi3.Schema {
	if ref.Value != nil {
		return ref.Value
	}
	if v, ok := i.SchemaValue[ref]; ok {
		return v
	}
	return nil
}

type SchemaInfo struct {
	ID int // reflectshape.Schema.Number

	OrderedProperties []string
	Links             []Link
}

type Link struct {
	Title string
	URL   string
}
