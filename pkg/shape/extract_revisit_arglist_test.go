package shape_test

import (
	"testing"

	"github.com/podhmo/reflect-openapi/pkg/arglist"
	"github.com/podhmo/reflect-openapi/pkg/shape"
)

type DB struct {
}

func Foo(db *DB)        {}
func Bar(anotherDB *DB) {}

func TestRevisitArglist(t *testing.T) {
	lookup := arglist.NewLookup()

	t.Run("without-lookup", func(t *testing.T) {
		e := shape.NewExtractor()
		e.ArglistLookup = nil

		{
			s := e.Extract(Foo).(shape.Function)
			want := "args0"
			got := s.Params.Keys[0]
			if s.Params.Len() != 1 {
				t.Errorf("%s: invalid arg list, len(args) == %d", s.GetName(), s.Params.Len())
			}
			if want != got {
				t.Errorf("%s: args[0] name, want %q but got %q", s.GetName(), want, got)
			}
		}
		{
			s := e.Extract(Bar).(shape.Function)
			want := "args0"
			got := s.Params.Keys[0]
			if s.Params.Len() != 1 {
				t.Errorf("%s: invalid arg list, len(args) == %d", s.GetName(), s.Params.Len())
			}
			if want != got {
				t.Errorf("%s: args[0] name, want %q but got %q", s.GetName(), want, got)
			}
		}
	})

	t.Run("disable", func(t *testing.T) {
		e := shape.NewExtractor()
		e.ArglistLookup = lookup

		{
			s := e.Extract(Foo).(shape.Function)
			want := "db"
			got := s.Params.Keys[0]
			if s.Params.Len() != 1 {
				t.Errorf("%s: invalid arg list, len(args) == %d", s.GetName(), s.Params.Len())
			}
			if want != got {
				t.Errorf("%s: args[0] name, want %q but got %q", s.GetName(), want, got)
			}
		}
		{
			s := e.Extract(Bar).(shape.Function)
			want := "db"
			got := s.Params.Keys[0]
			if s.Params.Len() != 1 {
				t.Errorf("%s: invalid arg list, len(args) == %d", s.GetName(), s.Params.Len())
			}
			if want != got {
				t.Errorf("%s: args[0] name, want %q but got %q", s.GetName(), want, got)
			}
		}
	})

	t.Run("enable", func(t *testing.T) {
		e := shape.NewExtractor()
		e.ArglistLookup = lookup
		e.RevisitArglist = true

		{
			s := e.Extract(Foo).(shape.Function)
			want := "db"
			got := s.Params.Keys[0]
			if s.Params.Len() != 1 {
				t.Errorf("%s: invalid arg list, len(args) == %d", s.GetName(), s.Params.Len())
			}
			if want != got {
				t.Errorf("%s: args[0] name, want %q but got %q", s.GetName(), want, got)
			}
		}
		{
			s := e.Extract(Bar).(shape.Function)
			want := "anotherDB"
			got := s.Params.Keys[0]
			if s.Params.Len() != 1 {
				t.Errorf("%s: invalid arg list, len(args) == %d", s.GetName(), s.Params.Len())
			}
			if want != got {
				t.Errorf("%s: args[0] name, want %q but got %q", s.GetName(), want, got)
			}
		}
	})
}
