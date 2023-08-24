package main

import (
	"fmt"
	"go/parser"
	"go/token"
	"log"
)

func Hello() {}

func main() {
	fset := token.NewFileSet()
	t, err := parser.ParseFile(fset, "main.go", nil, parser.ParseComments)
	if err != nil {
		log.Fatalf("parse file: %+v", err)
	}

	pos := t.Scope.Lookup("Hello").Pos()
	fmt.Println(fset.Position(pos).String())
}
