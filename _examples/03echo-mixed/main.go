package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"reflect"
	"regexp"
	"strconv"
	"strings"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	reflectopenapi "github.com/podhmo/reflect-openapi"
	"github.com/podhmo/reflect-openapi/docgen"
	"github.com/podhmo/reflect-openapi/dochandler"
	"github.com/podhmo/reflect-openapi/info"
)

// simplified version of this.
// https://swagger.io/docs/specification/basic-structure/

// required check with github.com/go-playground 's manner

type User struct {
	ID   int    `json:"id"`
	Name string `json:"name" validate:"required"` // for go-playground/validator
}

type CustomValidator struct {
	validator *validator.Validate
}

func (cv *CustomValidator) Validate(i interface{}) error {
	return cv.validator.Struct(i)
}

var (
	users = []User{
		{ID: 1, Name: "foo"},
		{ID: 2, Name: "bar"},
	}
)

// ListUsers returns a list of users.
func ListUsers() []User {
	return users
}

// InsertUser inserts user.
func InsertUser(user User) User {
	if user.ID == 0 {
		user.ID = len(users)
	}
	users = append(users, user)
	return user
}

type GetUserInput struct {
	UserID int `json:"userId" in:"path"`
}

// GetUser returns user
func GetUser(input GetUserInput) (User, error) {
	userID := input.UserID
	for _, u := range users {
		if u.ID == userID {
			return u, nil
		}
	}
	return User{}, fmt.Errorf("not found")
}

// ----------------------------------------
type Setup struct {
	Echo *echo.Echo
	*reflectopenapi.Manager
	Info *info.Info
}

var (
	rx = regexp.MustCompile(`:([^/:]+)`)
)

func (s *Setup) AddEndpoint(
	method, path string,
	interactor interface{},
	handler echo.HandlerFunc,
) {
	// for web api
	s.Echo.Add(
		method, path, handler,
	)

	// for doc
	openapiPath := rx.ReplaceAllString(path, `{$1}`)
	// log.Println("replace path: ", path, "->", openapiPath)

	s.RegisterFunc(interactor).After(func(op *openapi3.Operation) {
		s.Doc.AddOperation(openapiPath, method, op)
	})
}

func (s *Setup) SetupEndpoints() {
	s.AddEndpoint(
		"GET", "/users", ListUsers,
		func(c echo.Context) error {
			users := ListUsers()
			return c.JSON(200, users)
		},
	)
	s.AddEndpoint(
		"POST", "/users", InsertUser,
		func(c echo.Context) error {
			var u User
			if err := c.Bind(&u); err != nil {
				return c.JSON(400, map[string]string{"message": err.Error()})
			}
			if err := c.Validate(u); err != nil {
				return c.JSON(400, map[string]string{"message": err.Error()})
			}
			InsertUser(u)
			return c.JSON(201, users) // not supported in openapi doc
		},
	)
	s.AddEndpoint(
		"GET", "/users/:userId", GetUser,
		func(c echo.Context) error {
			var input GetUserInput
			{
				v := c.Param("userId")
				userID, err := strconv.ParseInt(v, 10, 0)
				if err != nil {
					return c.JSON(400, map[string]string{"message": err.Error()})
				}
				input.UserID = int(userID)
			}

			user, err := GetUser(input)
			if err != nil {
				return c.JSON(404, map[string]string{"message": err.Error()})
			}
			return c.JSON(200, user)
		},
	)
}

func (s *Setup) SetupSwaggerUI(addr string) {
	doc := s.Doc
	doc.Servers = append([]*openapi3.Server{{
		URL:         fmt.Sprintf("http://localhost%s", addr),
		Description: "local development server",
	}}, doc.Servers...)

	h := dochandler.New(doc, "/_doc", s.Info)
	s.Echo.Any("/_doc/*", echo.WrapHandler(h))
}

// ----------------------------------------

type FieldError struct {
	Path    string `json:"path"`
	Message string `json:"message"`
}
type APIError struct {
	Message string                `json:"message"`
	Details map[string]FieldError `json:"details"`
}

var options struct {
	useDoc bool

	docFile string
	mdFile  string
}

func main() {
	flag.BoolVar(&options.useDoc, "doc", false, "generate dod")
	flag.StringVar(&options.docFile, "docfile", "", "write file name (openapi.json)")
	flag.StringVar(&options.mdFile, "mdfile", "", "write file name (README.md)")
	flag.Parse()

	if err := run(); err != nil {
		log.Fatalf("%+v", err)
	}
}

func run() error {
	addr := os.Getenv("ADDR")
	if addr == "" {
		addr = ":44444"
	}

	e := echo.New()
	e.Validator = &CustomValidator{validator: validator.New()}
	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "method=${method}, uri=${uri}, status=${status}\n",
	}))

	ctx := context.Background()
	c := reflectopenapi.Config{
		SkipValidation: false,
		StrictSchema:   true,
		DefaultError:   APIError{},
		Info:           info.New(),
		EnableAutoTag:  true,
		IsRequiredCheckFunction: func(f reflect.StructTag) bool {
			v, ok := f.Lookup("validate")
			if !ok {
				return false
			}
			return strings.Contains(v, "required")
		},
	}
	doc, err := c.BuildDoc(ctx, func(m *reflectopenapi.Manager) {
		s := &Setup{Manager: m, Echo: e, Info: c.Info}
		s.SetupEndpoints()
		s.SetupSwaggerUI(addr)
	})
	if err != nil {
		return err
	}

	if options.useDoc {
		{
			log.Println("generate openapi doc")
			var w io.Writer = os.Stdout
			if options.docFile != "" {
				f, err := os.Create(options.docFile)
				if err != nil {
					return fmt.Errorf("open docfile: %w", err)
				}
				defer f.Close()
				w = f
			}

			enc := json.NewEncoder(w)
			enc.SetIndent("", "  ")
			if err := enc.Encode(doc); err != nil {
				return fmt.Errorf("write docfile: %w", err)
			}
		}

		if options.mdFile != "" {
			log.Println("generate README")
			f, err := os.Create(options.mdFile)
			if err != nil {
				return fmt.Errorf("open mdfile: %w", err)
			}
			defer f.Close()

			d := docgen.Generate(doc, c.Info)
			if err := docgen.WriteDoc(f, d); err != nil {
				return fmt.Errorf("write mdfile: %w", err)
			}
		}
		return nil
	}

	log.Println("listening ...", addr)
	return http.ListenAndServe(addr, e)
}
