package dochandler

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/google/go-cmp/cmp"
	reflectopenapi "github.com/podhmo/reflect-openapi"
	"github.com/podhmo/tenuki"
)

// Hello world
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
			)
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
					op.OperationID = "github.com/podhmo/reflect-openapi/dochandler.Byebye"
				})
			m.Doc.AddOperation("/byebye", "POST", op)
		}
		handler = New(m.Doc, "")
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
		{Method: "POST", Path: "/byebye", OperationID: "github.com/podhmo/reflect-openapi/dochandler.Byebye", Summary: "Byebye world"},
		{Method: "POST", Path: "/hello", OperationID: "github.com/podhmo/reflect-openapi/dochandler.Hello", Summary: "Hello world"},
		// added by dochandler package
		{Method: "GET", Path: "/doc", OperationID: "OpenAPIDocHandler", Summary: "(added by github.com/podhmo/reflect-openapi/dochandler)"},
		{Method: "GET", Path: "/ui", OperationID: "SwaggerUIHandler", Summary: "(added by github.com/podhmo/reflect-openapi/dochandler)"},
		{Method: "GET", Path: "/redoc", OperationID: "RedocHandler", Summary: "(added by github.com/podhmo/reflect-openapi/dochandler)"},
	}
	var got []Endpoint
	f.Extract().BindJSON(res, &got)
	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("handler response mismatch (-want +got):\n%s", diff)
	}
}
