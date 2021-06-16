package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/gofiber/fiber/v2"
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
type Setup interface {
	AddEndpoint(
		method, path string,
		interactor interface{},
		handler fiber.Handler,
	)
}

type APISetup struct {
	App *fiber.App
}

func (s *APISetup) AddEndpoint(
	method, path string,
	interactor interface{},
	handler fiber.Handler,
) {
	s.App.Add(
		method, path, handler,
	)
}

type DocSetup struct {
	*reflectopenapi.Manager
}

var (
	rx = regexp.MustCompile(`:([^/:]+)`) // TODO: support optional parameter?
)

func (s *DocSetup) AddEndpoint(
	method, path string,
	interactor interface{},
	handler fiber.Handler,
) {
	oaPath := rx.ReplaceAllString(path, `{$1}`)
	// log.Println("replace path: ", path, "->", oaPath)
	s.RegisterFunc(interactor).After(func(op *openapi3.Operation) {
		s.Doc.AddOperation(path, method, op)
	})
}

func Mount(s Setup) {
	s.AddEndpoint(
		"GET", "/users", ListUsers,
		func(c *fiber.Ctx) error {
			users := ListUsers()
			return c.JSON(users)
		},
	)
	s.AddEndpoint(
		"GET", "/users/:userId", GetUser,
		func(c *fiber.Ctx) error {
			var input GetUserInput
			{
				v := c.Params("userId")
				userID, err := strconv.ParseInt(v, 10, 0)
				if err != nil {
					return c.Status(400).JSON(map[string]string{"message": err.Error()})
				}
				input.UserID = int(userID)
			}

			user, err := GetUser(input)
			if err != nil {
				return c.Status(404).JSON(map[string]string{"message": err.Error()})
			}
			return c.JSON(user)
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
		}
		doc, err := c.BuildDoc(context.Background(), func(m *reflectopenapi.Manager) {
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

	s := &APISetup{App: fiber.New()}
	Mount(s)

	log.Println("listening ...", addr)
	return s.App.Listen(addr)
}
