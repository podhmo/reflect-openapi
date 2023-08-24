package main

import (
	"fmt"

	reflectopenapi "github.com/podhmo/reflect-openapi"
	reflectshape "github.com/podhmo/reflect-shape"
)

type Wrap[T any] struct {
	Value T
}

func Example() {
	c := &reflectshape.Config{}

	fmt.Println(c.Extract(reflectopenapi.Header{}).Name)
	fmt.Println(c.Extract(Wrap[reflectopenapi.Header]{Value: reflectopenapi.Header{}}).Name)
	// Output:
	// Header
	// Wrap[github.com/podhmo/reflect-openapi.Header]
}
