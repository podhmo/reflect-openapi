package main

import (
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

type GetUserInput struct {
	UserID int `json:"userId" openapi:"path"`
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
type Router interface {
	AddEndpoint(
		method, path string,
		interactor interface{},
		handler http.HandlerFunc,
	)
}

type APIRouter struct {
	Router chi.Router
}

func (r *APIRouter) AddEndpoint(
	method, path string,
	interactor interface{},
	handler http.HandlerFunc,
) {
	r.Router.Method(
		method, path, handler,
	)
}

type DocRouter struct {
	Doc      *openapi3.Swagger
	Resolver reflectopenapi.Resolver
	Visitor  *reflectopenapi.Visitor
}

func (r *DocRouter) AddEndpoint(
	method, path string,
	interactor interface{},
	handler http.HandlerFunc,
) {
	op := r.Visitor.VisitFunc(interactor)
	r.Doc.AddOperation(path, method, op)
}

func Mount(r Router) {
	r.AddEndpoint(
		"GET", "/users", ListUsers,
		func(w http.ResponseWriter, req *http.Request) {
			users := ListUsers()
			render.JSON(w, req, users)
		},
	)
	r.AddEndpoint(
		"GET", "/users/{userId}", GetUser,
		func(w http.ResponseWriter, req *http.Request) {
			var input GetUserInput
			if err := input.Bind(req); err != nil {
				render.Status(req, 400)
				render.JSON(w, req, map[string]string{"message": err.Error()})
				return
			}

			user, err := GetUser(input)
			if err != nil {
				render.Status(req, 404)
				render.JSON(w, req, map[string]string{"message": err.Error()})
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

		doc, err := reflectopenapi.NewDoc()
		if err != nil {
			return err
		}

		r := &reflectopenapi.UseRefResolver{}
		v := reflectopenapi.NewVisitor(r)

		router := &DocRouter{
			Resolver: r,
			Visitor:  v,
			Doc:      doc,
		}
		Mount(router)
		r.Bind(doc)

		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "  ")
		return enc.Encode(doc)
	}

	router := &APIRouter{Router: chi.NewRouter()}
	router.Router.Use(middleware.RealIP)
	router.Router.Use(middleware.Logger)
	Mount(router)

	log.Println("listening ...", addr)
	return http.ListenAndServe(addr, router.Router)
}
