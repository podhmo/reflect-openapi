package shape

import (
	"log"

	"github.com/podhmo/reflect-openapi/pkg/arglist"
)

func fixupArglist(lookup *arglist.Lookup, fn *Function, ob interface{}, fullname string) {
	params := fn.Params.Keys
	returns := fn.Returns.Keys

	// fixup names
	nameset, err := lookup.LookupNameSetFromFunc(ob)
	if err != nil {
		log.Printf("function %q, arglist lookup is failed %v", fullname, err)
	}
	if len(nameset.Args) != len(params) {
		log.Printf("the length of arguments is mismatch, %d != %d", len(nameset.Args), len(params))
	} else {
		fn.Params.Keys = nameset.Args
	}
	if len(nameset.Returns) != len(returns) {
		log.Printf("the length of returns is mismatch, %d != %d", len(nameset.Returns), len(returns))
	} else {
		fn.Returns.Keys = nameset.Returns
	}
}
