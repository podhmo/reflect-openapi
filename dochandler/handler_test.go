package dochandler

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/google/go-cmp/cmp"
	reflectopenapi "github.com/podhmo/reflect-openapi"
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
			m.RegisterFunc(Hello).After(func(op *openapi3.Operation) {
				m.Doc.AddOperation("/hello", "POST", op)
			})
		}
		{
			m.RegisterFunc(
				func(input struct {
					Name string `json:"name"`
				}) string {
					return ""
				}).After(
				func(op *openapi3.Operation) {
					op.Summary = "Byebye world"
					op.OperationID = "github.com/podhmo/reflect-openapi/dochandler.Byebye"
					m.Doc.AddOperation("/byebye", "POST", op)
				})

		}
		handler = New(m.Doc, "")
	})

	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, httptest.NewRequest("GET", "/", nil))
	res := rec.Result()

	if want, got := 200, res.StatusCode; want != got {
		t.Errorf("unexpoected status code: want:%d, but got:%d", want, got)
	}

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
	if err := json.NewDecoder(res.Body).Decode(&got); err != nil {
		t.Fatalf("unexpected decode error: %+v", err)
	}
	defer res.Body.Close()
	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("handler response mismatch (-want +got):\n%s", diff)
	}
}
