package main

import (
	"context"

	"github.com/getkin/kin-openapi/openapi3"
	reflectopenapi "github.com/podhmo/reflect-openapi"
)

type HasTitle struct {
	Title string `json:"title"`
}

// Content object
type Content struct {
	HasTitle
	Author string `json:"author"` // default john
	Pretty bool   `in:"query" query:"pretty"`
	Flags
}

type Flags struct {
	BoolOK       bool    `in:"query" query:"boolOK"`
	BoolNG       bool    `in:"query" query:"boolNG"`
	BoolPtrOK    *bool   `in:"query" query:"boolPtrOK"`
	BoolPtrNG    *bool   `in:"query" query:"boolPtrNG"`
	BoolPtrNil   *bool   `in:"query" query:"boolPtrNil"`
	StringPtrNil *string `in:"query" query:"stringPtrNil"`
}

type PostContentInput Content

type PostContentInput2 struct {
	Content
	*Ref
}

type Ref struct {
	XXX string `json:"xxx" required:"false"`
}

// PostContent action
func PostContent(ctx context.Context, input PostContentInput) (Content, error) {
	return Content{}, nil
}

// PostContent2 action
func PostContent2(ctx context.Context, input PostContentInput2) (Content, error) {
	return Content{}, nil
}

func main() {
	c := &reflectopenapi.Config{}
	c.EmitDoc(func(m *reflectopenapi.Manager) {
		ok := true
		ng := false
		flags := Flags{
			BoolOK:     ok,
			BoolNG:     ng,
			BoolPtrOK:  &ok,
			BoolPtrNG:  &ng,
			BoolPtrNil: nil,
		}
		m.RegisterFunc(PostContent).After(func(op *openapi3.Operation) {
			m.Doc.AddOperation("/Contents", "POST", op)
		}).DefaultInput(PostContentInput{Author: "john", Pretty: true, Flags: flags})
		m.RegisterFunc(PostContent2).After(func(op *openapi3.Operation) {
			m.Doc.AddOperation("/Contents2", "POST", op)
		}).DefaultInput(PostContentInput2{Content: Content{Author: "john", Pretty: true, Flags: flags}})
	})
}
