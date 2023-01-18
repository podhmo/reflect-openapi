package reflectopenapi

import (
	"context"
	"fmt"
	"reflect"
	"strings"

	shape "github.com/podhmo/reflect-shape"
)

var rcontextType = reflect.TypeOf(func(context.Context) {}).In(0)

type DefaultSelector struct {
	FirstParamInputSelector
	FirstParamOutputSelector
}

type FirstParamInputSelector struct{}

func (s *FirstParamInputSelector) SelectInput(fn *shape.Func) *shape.Shape {
	args := fn.Args()
	if args.Len() == 0 {
		return nil
	}
	for _, x := range args {
		if x.Shape.Type == rcontextType {
			continue
		}
		return x.Shape
	}
	return nil
}

type MergeParamsInputSelector struct {
	Extractor Extractor
}

func (s *MergeParamsInputSelector) useArglist() {
}
func (s *MergeParamsInputSelector) SelectInput(fn *shape.Func) *shape.Shape {
	args := fn.Args()
	if args.Len() == 0 {
		return nil
	}

	var fields []reflect.StructField
	for _, p := range args {
		if p.Shape.Type == rcontextType {
			continue
		}

		// todo: handling customization
		required := p.Shape.Lv == 0
		tag := fmt.Sprintf(`json:%q`, p.Name)
		if !required {
			switch p.Shape.Kind {
			case reflect.Chan, reflect.Interface, reflect.Slice, reflect.Array, reflect.Struct:
			default:
				tag += ` openapi:"query"`
			}
		} else {
			tag += ` required:"true"`
		}
		fields = append(fields, reflect.StructField{
			Name: strings.ToTitle(p.Name),
			Type: p.Shape.Type,
			Tag:  reflect.StructTag(tag),
		})
	}

	// create new struct with reflect
	rtype := reflect.StructOf(fields)
	rval := reflect.New(rtype)
	return s.Extractor.Extract(rval.Interface())
}

type FirstParamOutputSelector struct{}

func (s *FirstParamOutputSelector) SelectOutput(fn *shape.Func) *shape.Shape {
	returns := fn.Returns()
	if returns.Len() == 0 {
		return nil
	}
	return returns[0].Shape
}
