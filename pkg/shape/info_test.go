package shape_test

import (
	"reflect"
	"testing"

	"github.com/podhmo/reflect-openapi/pkg/shape"
)

type Article struct {
	Title string
}

func TestInfo(t *testing.T) {
	got := shape.Extract(Article{})

	t.Run("Shape", func(t *testing.T) {
		got := got.Shape()
		want := "struct"
		if want != got {
			t.Errorf("expected %q, but %q", want, got)
		}
	})
	t.Run("GetName", func(t *testing.T) {
		got := got.GetName()
		want := "Article"
		if want != got {
			t.Errorf("expected %q, but %q", want, got)
		}
	})
	t.Run("GetPackage", func(t *testing.T) {
		got := got.GetPackage()
		want := "github.com/podhmo/reflect-openapi/pkg/shape_test"
		if want != got {
			t.Errorf("expected %q, but %q", want, got)
		}
	})
	t.Run("GetReflectKind", func(t *testing.T) {
		got := got.GetReflectKind()
		want := reflect.Struct
		if want != got {
			t.Errorf("expected %v, but %v", want, got)
		}
	})
	t.Run("GetReflectType", func(t *testing.T) {
		got := got.GetReflectType()
		want := reflect.TypeOf(Article{})
		if want != got {
			t.Errorf("expected %v, but %v", want, got)
		}
	})
	t.Run("GetReflectValue", func(t *testing.T) {
		got := got.GetReflectValue()
		want := reflect.ValueOf(Article{})
		if want != got {
			t.Errorf("expected %v, but %v", want, got)
		}
	})
}
