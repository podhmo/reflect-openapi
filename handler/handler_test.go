package handler

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/getkin/kin-openapi/openapi3"
	reflectopenapi "github.com/podhmo/reflect-openapi"
	"github.com/podhmo/tenuki"
)

func Hello(input struct {
	Name string `json:"name"`
}) string {
	return ""
}

func TestEndpoints(t *testing.T) {
	c := &reflectopenapi.Config{}
	var handler http.Handler
	c.BuildDoc(context.Background(), func(m *reflectopenapi.Manager) {
		{
			op := m.Visitor.VisitFunc(
				Hello,
				func(op *openapi3.Operation) {
					op.Summary = "Hello world"
				})
			m.Doc.AddOperation("/hello", "POST", op)
		}
		{
			op := m.Visitor.VisitFunc(
				func(input struct {
					Name string `json:"name"`
				}) string {
					return ""
				},
				func(op *openapi3.Operation) {
					op.Summary = "Byebye world"
				})
			m.Doc.AddOperation("/byebye", "POST", op)
		}
		handler = NewHandler(m.Doc)
	})

	ts := httptest.NewServer(handler)
	defer ts.Close()

	f := tenuki.New(t)
	res := f.Do(
		f.NewRequest("GET", fmt.Sprintf("%s/", ts.URL), nil),
		tenuki.AssertStatus(http.StatusOK),
	)

	// assertion
	want := []Endpoint{
		{Method: "POST", Path: "/byebye", OperationID: "github.com/podhmo/reflect-openapi/handler.Byebye", Summary: "Byebye world"},
		{Method: "POST", Path: "/hello", OperationID: "github.com/podhmo/reflect-openapi/handler.Hello", Summary: "Hello world"},
	}
	var got []Endpoint
	f.Extract().JSON(res, &got)
	if !reflect.DeepEqual(want, got) {
		t.Errorf("response body\nwant\n\t%+v\nbut\n\t%+v", want, got)
	}
}
