package arglist

import (
	"fmt"
	"go/ast"
)

type NameSet struct {
	Name    string
	Args    []string
	Returns []string
}

func InspectFunc(decl *ast.FuncDecl) (NameSet, error) {
	var r NameSet
	r.Name = decl.Name.Name
	if decl.Type.Params != nil {
		var names []string
		i := 0
		for _, x := range decl.Type.Params.List {
			if len(x.Names) == 0 {
				names = append(names, fmt.Sprintf("arg%d", i))
				i++
				continue
			}
			if _, ok := x.Type.(*ast.Ellipsis); ok {
				names = append(names, fmt.Sprintf("*%s", x.Names[0].Name))
				continue
			}
			for _, ident := range x.Names {
				names = append(names, ident.Name)
			}
		}
		r.Args = names
	}
	if decl.Type.Results != nil {
		var names []string
		i := 0
		for _, x := range decl.Type.Results.List {
			if len(x.Names) == 0 {
				names = append(names, fmt.Sprintf("ret%d", i))
				i++
				continue
			}
			for _, ident := range x.Names {
				names = append(names, ident.Name)
			}
		}
		r.Returns = names
	}
	return r, nil
}
