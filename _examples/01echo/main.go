package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"regexp"
	"strconv"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
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
		handler echo.HandlerFunc,
	)
}

type APIRouter struct {
	Echo *echo.Echo
}

func (r *APIRouter) AddEndpoint(
	method, path string,
	interactor interface{},
	handler echo.HandlerFunc,
) {
	r.Echo.Add(
		method, path, handler,
	)
}

type DocRouter struct {
	Doc      *openapi3.Swagger
	Resolver reflectopenapi.Resolver
	Visitor  *reflectopenapi.Visitor
}

var (
	rx = regexp.MustCompile(`:([^/:]+)`)
)

func (r *DocRouter) AddEndpoint(
	method, path string,
	interactor interface{},
	handler echo.HandlerFunc,
) {
	oaPath := rx.ReplaceAllString(path, `{$1}`)
	// log.Println("replace path: ", path, "->", oaPath)
	op := r.Visitor.VisitFunc(interactor)
	r.Doc.AddOperation(oaPath, method, op)
}

func Mount(r Router) {
	r.AddEndpoint(
		"GET", "/users", ListUsers,
		func(c echo.Context) error {
			users := ListUsers()
			return c.JSON(200, users)
		},
	)
	r.AddEndpoint(
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

	e := echo.New()
	router := &APIRouter{Echo: e}
	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "method=${method}, uri=${uri}, status=${status}\n",
	}))
	Mount(router)

	log.Println("listening ...", addr)
	return http.ListenAndServe(addr, e)
}
