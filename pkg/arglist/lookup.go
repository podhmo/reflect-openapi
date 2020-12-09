package arglist

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io/ioutil"
	"reflect"
	"runtime"
)

// TODO: merge with ../comment

type Lookup struct {
	fset *token.FileSet

	fileCache map[string]*ast.File
	declCache map[*ast.File]map[string][]*ast.FuncDecl
}

// NewLookup is the factory function creating Lookup
func NewLookup() *Lookup {
	return &Lookup{
		fset:      token.NewFileSet(),
		fileCache: map[string]*ast.File{},
		declCache: map[*ast.File]map[string][]*ast.FuncDecl{},
	}
}

func (l *Lookup) LookupAST(filename string) (*ast.File, error) {
	if f, ok := l.fileCache[filename]; ok {
		return f, nil
	}
	mode := parser.ParseComments
	code, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	f, err := parser.ParseFile(l.fset, filename, code, mode)
	if err != nil {
		return nil, err
	}
	l.fileCache[filename] = f
	return f, nil
}

func (l *Lookup) LookupFuncDecl(filename string, targetName string) (*ast.FuncDecl, error) {
	f, err := l.LookupAST(filename)
	if err != nil {
		return nil, err
	}
	ob := f.Scope.Lookup(targetName)
	if ob == nil {
		return nil, fmt.Errorf("not found %q in %q", targetName, filename)
	}
	decl, ok := ob.Decl.(*ast.FuncDecl)
	if !ok {
		return nil, fmt.Errorf("%q is unexpected type %T", targetName, ob)
	}
	return decl, nil
}

func (l *Lookup) LookupArglistFromFunc(fn interface{}) ([]string, error) {
	rfunc := runtime.FuncForPC(reflect.ValueOf(fn).Pointer())
	filename, _ := rfunc.FileLine(rfunc.Entry())
	decl, err := l.LookupFuncDecl(filename, rfunc.Name())
	if err != nil {
		return nil, err
	}
	r, err := InspectFunc(decl)
	if err != nil {
		return nil, err
	}
	return r.Args, nil
}
