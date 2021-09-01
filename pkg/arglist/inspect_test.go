package arglist

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"reflect"
	"testing"
)

func inspectFuncFromFile(f *ast.File, name string) (NameSet, error) {
	ob := f.Scope.Lookup(name)
	if ob == nil {
		return NameSet{}, fmt.Errorf("not found %q", name)
	}
	decl, ok := ob.Decl.(*ast.FuncDecl)
	if !ok {
		return NameSet{}, fmt.Errorf("unexpected decl %T", ob)
	}
	return InspectFunc(decl)
}

func TestInspectFunc(t *testing.T) {
	const code = `package foo
func Sum(x int, y,z int) int {
	return x + y + z
}
func Sum2(xs ...int) int {
	return 0
}
func Sprintf(ctx context.Context, fmt string, vs ...interface{}) (string, error) {
	return fmt.Sprintf(fmt, vs...), nil
}
func Sprintf2(ctx context.Context, fmt string, vs ...interface{}) (s string, err error) {
	return fmt.Sprintf(fmt, vs...), nil
}
`
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, "foo.go", code, parser.ParseComments)
	if err != nil {
		t.Fatalf("something wrong in parse-file %v", err)
	}

	cases := []struct {
		name string
		want NameSet
	}{
		{
			name: "Sum",
			want: NameSet{
				Name:    "Sum",
				Args:    []string{"x", "y", "z"},
				Returns: []string{"ret0"},
			},
		},
		{
			name: "Sum2",
			want: NameSet{
				Name:    "Sum2",
				Args:    []string{"*xs"},
				Returns: []string{"ret0"},
			},
		},
		{
			name: "Sprintf",
			want: NameSet{
				Name:    "Sprintf",
				Args:    []string{"ctx", "fmt", "*vs"},
				Returns: []string{"ret0", "ret1"},
			},
		},
		{
			name: "Sprintf2",
			want: NameSet{
				Name:    "Sprintf2",
				Args:    []string{"ctx", "fmt", "*vs"},
				Returns: []string{"s", "err"},
			},
		},
	}
	for _, c := range cases {
		c := c
		t.Run(c.name, func(t *testing.T) {
			got, err := inspectFuncFromFile(f, c.name)
			if err != nil {
				t.Fatalf("unexpected error %+v", err)
			}
			if !reflect.DeepEqual(c.want, got) {
				t.Errorf("want:\n\t%q\nbut got:\n\t%q\n", c.want, got)
			}
		})
	}
}

// Anonymous function is not supported, so.

func TestInspectAnonymousFunc(t *testing.T) {
	l := NewLookup()
	fn := func(x string) (string, error) { return "", nil }
	ns, err := l.LookupNameSetFromFunc(fn)
	if err == nil {
		t.Errorf("error is expected, but not error is occured")
	}
	if len(ns.Args) != 0 {
		t.Errorf("len(ns.Args) == 0 is expected, but got %d", len(ns.Args))
	}
	if len(ns.Returns) != 0 {
		t.Errorf("len(ns.Returns) == 0 is expected, but got %d", len(ns.Returns))
	}
}
