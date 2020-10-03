package shape_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/podhmo/reflect-openapi/pkg/shape"
)

type Person struct {
	Name string `json:"name"`
	Age  int
}

func TestPrimitive(t *testing.T) {
	t.Run("int", func(t *testing.T) {
		got := shape.Extract(1)
		if _, ok := got.(shape.Primitive); !ok {
			t.Errorf("expected Primitive, but %T", got)
		}

		// format
		if got := fmt.Sprintf("%v", got); got != "int" {
			t.Errorf("expected string expression is %q but %q", "int", got)
		}
	})

	t.Run("new type", func(t *testing.T) {
		type MyInt int
		got := shape.Extract(MyInt(1))
		if _, ok := got.(shape.Primitive); !ok {
			t.Errorf("expected Primitive, but %T", got)
		}

		// format
		if got, want := fmt.Sprintf("%v", got), "github.com/podhmo/reflect-openapi/pkg/shape_test.MyInt"; want != got {
			t.Errorf("expected string expression is %q but %q", want, got)
		}
	})

	t.Run("type alias", func(t *testing.T) {
		type MyInt = int
		got := shape.Extract(MyInt(1))
		if _, ok := got.(shape.Primitive); !ok {
			t.Errorf("expected Primitive, but %T", got)
		}

		// format
		if got, want := fmt.Sprintf("%v", got), "int"; want != got {
			t.Errorf("expected string expression is %q but %q", want, got)
		}
	})
}

func TestStruct(t *testing.T) {
	t.Run("user defined", func(t *testing.T) {
		got := shape.Extract(Person{})
		v, ok := got.(shape.Struct)
		if !ok {
			t.Errorf("expected Struct, but %T", got)
		}

		if len(v.Fields.Values) != 2 {
			t.Errorf("expected the number of Person's fields is 1, but %v", len(v.Fields.Values))
		}

		if got := v.FieldName(0); got != "name" {
			t.Errorf("expected field name with json tag is %q, but %q", "name", got)
		}
		if got := v.FieldName(1); got != "Age" {
			t.Errorf("expected field name without json tag is %q, but %q", "name", got)
		}

		// format
		if got, want := fmt.Sprintf("%v", got), "github.com/podhmo/reflect-openapi/pkg/shape_test.Person"; want != got {
			t.Errorf("expected string expression is %q but %q", want, got)
		}
	})

	t.Run("time.Time", func(t *testing.T) {
		var z time.Time
		got := shape.Extract(z)
		if _, ok := got.(shape.Struct); !ok {
			t.Errorf("expected Struct, but %T", got)
		}

		// format
		if got := fmt.Sprintf("%v", got); got != "time.Time" {
			t.Errorf("expected string expression is %q but %q", "int", got)
		}
	})
}

func TestContainer(t *testing.T) {
	t.Run("slice", func(t *testing.T) {
		t.Run("primitive", func(t *testing.T) {
			got := shape.Extract([]int{})
			v, ok := got.(shape.Container)
			if !ok {
				t.Errorf("expected Container, but %T", got)
			}
			if got := len(v.Args); got != 1 {
				t.Errorf("expected the length of slices's args is %v, but %v", 1, got)
			}

			if got, want := fmt.Sprintf("%v", got), "slice[int]"; want != got {
				t.Errorf("expected string expression is %q but %q", want, got)
			}
		})
		t.Run("struct", func(t *testing.T) {
			got := shape.Extract([]Person{})
			v, ok := got.(shape.Container)
			if !ok {
				t.Errorf("expected Container, but %T", got)
			}
			if got := len(v.Args); got != 1 {
				t.Errorf("expected the length of slices's args is %v, but %v", 1, got)
			}

			// format
			if got, want := fmt.Sprintf("%v", got), "slice[github.com/podhmo/reflect-openapi/pkg/shape_test.Person]"; want != got {
				t.Errorf("expected string expression is %q but %q", want, got)
			}
		})
	})

	t.Run("map", func(t *testing.T) {
		t.Run("primitive", func(t *testing.T) {
			got := shape.Extract(map[string]int{})
			v, ok := got.(shape.Container)
			if !ok {
				t.Errorf("expected Container, but %T", got)
			}
			if got := len(v.Args); got != 2 {
				t.Errorf("expected the length of slices's args is %v, but %v", 1, got)
			}

			// format
			if got, want := fmt.Sprintf("%v", got), "map[string, int]"; want != got {
				t.Errorf("expected string expression is %q but %q", want, got)
			}
		})
		t.Run("struct", func(t *testing.T) {
			got := shape.Extract(map[string][]Person{})
			v, ok := got.(shape.Container)
			if !ok {
				t.Errorf("expected Container, but %T", got)
			}
			if got := len(v.Args); got != 2 {
				t.Errorf("expected the length of slices's args is %v, but %v", 1, got)
			}

			// format
			if got, want := fmt.Sprintf("%v", got), "map[string, slice[github.com/podhmo/reflect-openapi/pkg/shape_test.Person]]"; want != got {
				t.Errorf("expected string expression is %q but %q", want, got)
			}
		})
	})
}

type ListUserInput struct {
	Query string
	Limit int
}

func ListUser(ctx context.Context, input ListUserInput) ([]Person, error) {
	return nil, nil
}

func TestFunction(t *testing.T) {
	got := shape.Extract(ListUser)
	_, ok := got.(shape.Function)
	if !ok {
		t.Errorf("expected Container, but %T", got)
	}
	if got, want := fmt.Sprintf("%v", got), "github.com/podhmo/reflect-openapi/pkg/shape_test.ListUser(context.Context, github.com/podhmo/reflect-openapi/pkg/shape_test.ListUserInput) (slice[github.com/podhmo/reflect-openapi/pkg/shape_test.Person], error)"; want != got {
		t.Errorf("expected string expression is %q but %q", want, got)
	}
}
