package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/getkin/kin-openapi/openapi3"
	reflectopenapi "github.com/podhmo/reflect-openapi"
	"github.com/podhmo/reflect-openapi/docgen"
	"github.com/podhmo/reflect-openapi/info"
)

var options struct {
	DocFile string
}

func main() {
	flag.StringVar(&options.DocFile, "docfile", "", "write openapi doc file")
	flag.Parse()

	if err := run(); err != nil {
		log.Fatalf("!! %+v", err)
	}
}

func run() error {
	c := &reflectopenapi.Config{
		Info:          info.New(), // need!
		EnableAutoTag: true,
	}
	ctx := context.Background()
	tree, err := c.BuildDoc(ctx, func(m *reflectopenapi.Manager) {
		m.Doc.Info.Title = "hello"
		m.Doc.Info.Version = "1.0.0"
		m.Doc.Info.Description = `This is the example has text/html output`

		mount(m)

		// login
		m.RegisterFuncText(Login, "text/html", func(op *openapi3.Operation) {
			m.Doc.AddOperation("/login", "POST", op)
		})
	})
	if err != nil {
		return fmt.Errorf("build: %w", err)
	}

	if options.DocFile != "" {
		f, err := os.Create(options.DocFile)
		if err != nil {
			return fmt.Errorf("open openapi doc: %w", err)
		}
		enc := json.NewEncoder(f)
		enc.SetIndent("", "  ")
		if err := enc.Encode(tree); err != nil {
			return fmt.Errorf("write openapi doc: %w", err)
		}
	}

	doc := docgen.Generate(tree, c.Info)
	if err := docgen.WriteDoc(os.Stdout, doc); err != nil {
		return fmt.Errorf("generate: %w", err)
	}
	return nil
}

func mount(m *reflectopenapi.Manager) {
	m.RegisterFunc(Hello, func(op *openapi3.Operation) {
		m.Doc.AddOperation("/api/hello", "POST", op)
	})

	m.RegisterFuncText(HelloHTML, "text/html", func(op *openapi3.Operation) {
		m.Doc.AddOperation("/hello/{name}", "GET", op)
	})

	// with custom error response (responses['default'])
	m.RegisterFuncText(HelloHTML2, "text/html", func(op *openapi3.Operation) {
		m.Doc.AddOperation("/hello2/{name}", "GET", op)
	}).Error(Error{}, "default error response")

	// with response header
	m.RegisterFuncText(HelloHTML3, "text/html", func(op *openapi3.Operation) {
		m.Doc.AddOperation("/hello3/{name}", "GET", op)
	}).Error(Error{}, "default error response")
}

func Hello(input struct {
	Name string `json:"name"`
}) (output struct {
	Message string `json:"message"`
}) {
	output.Message = fmt.Sprintf("hello %s", input.Name)
	return
}

func HelloHTML(input struct {
	Name string `path:"name" in:"path"`
}) string /* html with greeting message */ {
	return fmt.Sprintf("<p>hello %s</p>", input.Name)
}

func HelloHTML2(input struct {
	Name string `path:"name" in:"path"`
}) string /* html with greeting message */ {
	return fmt.Sprintf("<p>hello %s</p>", input.Name)
}

func HelloHTML3(input struct {
	Name string `path:"name" in:"path"`
}) (output struct {
	Body       string //html with greeting message
	XSomething string `in:"header" header:"X-SOMETHING"`
}) {
	output.Body = fmt.Sprintf("<p>hello %s</p>", input.Name)
	output.XSomething = "something"
	return
}

// Error is custom error response
type Error struct {
	Message string `json:"message"`
}

// https://swagger.io/docs/specification/authentication/cookie-authentication/
type LoginInput struct {
	Name     string `json:"name"`
	Password string `json:"password"`
}
type LoginOutput struct {
	CookieHeader string `in:"header" header:"Set-Cookie" opneapi-override:"{'example': 'JSESSIONID=abcde12345; Path=/; HttpOnly'}"`
	Body         string
}

// Successfully authenticated.
// The session ID is returned in a cookie named `JSESSIONID`. You need to include this cookie in subsequent request
func Login(LoginInput) LoginOutput {
	return LoginOutput{}
}
