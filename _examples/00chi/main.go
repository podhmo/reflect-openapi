package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
	reflectopenapi "github.com/podhmo/reflect-openapi"
)

// simplified version of this.
// https://swagger.io/docs/specification/basic-structure/

type User struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Age  int    `json:"age" openapi-override:"{'minimum': 0}"`

	WithNickname
}
type WithNickname struct {
	Nickname string `json:"nickname,omitempty" openapi-override:"{'minLength': 1}"`
}

var (
	users = []User{
		{ID: 1, Name: "foo", Age: 20},
		{ID: 2, Name: "bar"},
	}
)

// ListUsers returns a list of users.
func ListUsers() []User {
	return users
}

type GetUserInput struct {
	UserID int `json:"userId" in:"path"`
}

func (input *GetUserInput) Bind(req *http.Request) error {
	v := chi.URLParam(req, "userId")
	userID, err := strconv.ParseInt(v, 10, 0)
	if err != nil {
		return err
	}
	input.UserID = int(userID)

	return nil
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
type Setup interface {
	AddEndpoint(
		method, path string,
		interactor interface{},
		handler http.HandlerFunc,
	)
}

type APISetup struct {
	Router chi.Router
}

func (s *APISetup) AddEndpoint(
	method, path string,
	interactor interface{},
	handler http.HandlerFunc,
) {
	s.Router.Method(
		method, path, handler,
	)
}

type DocSetup struct {
	*reflectopenapi.Manager
}

func (s *DocSetup) AddEndpoint(
	method, path string,
	interactor interface{},
	handler http.HandlerFunc,
) {
	s.RegisterFunc(interactor).After(func(op *openapi3.Operation) {
		s.Doc.AddOperation(path, method, op)
	}).
		Example(404, "application/json", "default", "not found value", APIError{"not found"}).
		Example(400, "application/json", "default", "validation error", APIError{"bad request"})
}

type APIError struct {
	Message string `json:"message"`
}

func Mount(s Setup) {
	s.AddEndpoint(
		"GET", "/users", ListUsers,
		func(w http.ResponseWriter, req *http.Request) {
			users := ListUsers()
			render.JSON(w, req, users)
		},
	)
	s.AddEndpoint(
		"GET", "/users/{userId}", GetUser,
		func(w http.ResponseWriter, req *http.Request) {
			var input GetUserInput
			if err := input.Bind(req); err != nil {
				render.Status(req, 400)
				render.JSON(w, req, APIError{err.Error()})
				return
			}

			user, err := GetUser(input)
			if err != nil {
				render.Status(req, 404)
				render.JSON(w, req, APIError{err.Error()})
				return
			}
			render.JSON(w, req, user)
		},
	)
}

// ----------------------------------------

func main() {
	useDocF := flag.Bool("doc", false, "generate doc")
	flag.Parse()

	if err := run(*useDocF); err != nil {
		log.Fatalf("%+v", err)
	}
}

func run(useDoc bool) error {
	addr := os.Getenv("ADDR")
	if addr == "" {
		addr = ":44444"
	}

	if useDoc {
		log.Println("generate openapi doc")
		c := reflectopenapi.Config{
			SkipValidation: false,
			StrictSchema:   true,
			DefaultError:   APIError{},
		}
		doc, err := c.BuildDoc(context.Background(), func(m *reflectopenapi.Manager) {
			m.RegisterType(User{Name: "foo"})
			s := &DocSetup{Manager: m}
			Mount(s)
		})
		if err != nil {
			return err
		}

		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "  ")
		return enc.Encode(doc)
	}

	s := &APISetup{Router: chi.NewRouter()}
	s.Router.Use(middleware.RealIP)
	s.Router.Use(middleware.Logger)
	Mount(s)

	log.Println("listening ...", addr)
	return http.ListenAndServe(addr, s.Router)
}
