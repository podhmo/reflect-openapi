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

var _ Selector = (*DefaultSelector)(nil)

type FirstParamInputSelector struct{}

func (s *FirstParamInputSelector) SelectInput(fn *shape.Func) (*shape.Shape, string) {
	args := fn.Args()
	if args.Len() == 0 {
		return nil, ""
	}
	for _, x := range args {
		if x.Shape.Type == rcontextType {
			continue
		}
		// TODO: set description
		return x.Shape, x.Doc
	}
	return nil, ""
}

type MergeParamsInputSelector struct {
	transformer *Transformer
}

func (s *MergeParamsInputSelector) NeedTransformer(t *Transformer) {
	s.transformer = t
}

var _ needTransformer = (*MergeParamsInputSelector)(nil)

func (s *MergeParamsInputSelector) SelectInput(fn *shape.Func) (*shape.Shape, string) {
	args := fn.Args()
	if args.Len() == 0 {
		return nil, ""
	}

	var fields []reflect.StructField
	for _, p := range args {
		if p.Shape.Type == rcontextType {
			continue
		}

		// todo: handling customization
		required := p.Shape.Lv == 0
		tag := fmt.Sprintf(`%s:%q`, s.transformer.TagNameOption.NameTag, p.Name)
		if !required {
			switch p.Shape.Kind {
			case reflect.Chan, reflect.Interface, reflect.Slice, reflect.Array, reflect.Struct:
			default:
				tag += fmt.Sprintf(` %s:"query"`, s.transformer.TagNameOption.ParamTypeTag)
			}
		} else {
			tag += ` required:"true"`
		}
		if p.Doc != "" {
			tag += fmt.Sprintf(` description:%q`, p.Doc)
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
	return s.transformer.Extractor.Extract(rval.Interface()), ""
}

type FirstParamOutputSelector struct{}

func (s *FirstParamOutputSelector) SelectOutput(fn *shape.Func) (*shape.Shape, string) {
	returns := fn.Returns()
	if returns.Len() == 0 {
		return nil, ""
	}
	return returns[0].Shape, returns[0].Doc
}
