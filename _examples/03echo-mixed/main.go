package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"regexp"
	"strconv"

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
type Setup struct {
	Echo *echo.Echo
	*reflectopenapi.Manager
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
	oaPath := rx.ReplaceAllString(path, `{$1}`)
	// log.Println("replace path: ", path, "->", oaPath)
	op := s.Visitor.VisitFunc(interactor)
	s.Doc.AddOperation(oaPath, method, op)
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

	e := echo.New()
	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "method=${method}, uri=${uri}, status=${status}\n",
	}))

	ctx := context.Background()
	c := reflectopenapi.Config{
		SkipValidation: false,
		StrictSchema:   true,
	}
	doc, err := c.BuildDoc(ctx, func(m *reflectopenapi.Manager) {
		s := &Setup{Manager: m, Echo: e}
		s.SetupEndpoints()
	})
	if err != nil {
		return err
	}

	if useDoc {
		log.Println("generate openapi doc")
		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "  ")
		return enc.Encode(doc)
	}

	log.Println("listening ...", addr)
	return http.ListenAndServe(addr, e)
}
