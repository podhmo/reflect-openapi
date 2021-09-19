package reflectopenapi

import (
	"reflect"

	shape "github.com/podhmo/reflect-shape"
)

type DefaultSelector struct {
	FirstParamInputSelector
	FirstParamOutputSelector
}

type FirstParamInputSelector struct{}

func (s *FirstParamInputSelector) SelectInput(fn shape.Function) shape.Shape {
	if len(fn.Params.Values) == 0 {
		return nil
	}
	for _, inob := range fn.Params.Values {
		if inob.GetFullName() == "context.Context" {
			continue
		}
		return inob
	}
	return nil
}

type MergeParamsInputSelector struct{}

func (s *MergeParamsInputSelector) useArglist() {
}
func (s *MergeParamsInputSelector) SelectInput(fn shape.Function) shape.Shape {
	if len(fn.Params.Values) == 0 {
		return nil
	}
	fields := shape.ShapeMap{}
	tags := make([]reflect.StructTag, 0, fn.Params.Len())
	metadata := make([]shape.FieldMetadata, 0, fn.Params.Len())
	for i, p := range fn.Params.Values {
		if p.GetFullName() == "context.Context" {
			continue
		}

		name := fn.Params.Keys[i]
		fields.Keys = append(fields.Keys, name)
		fields.Values = append(fields.Values, p)

		// todo: handling customization
		required := p.GetLv() == 0
		var tag reflect.StructTag
		if !required {
			if _, ok := p.(shape.Primitive); ok {
				tag = reflect.StructTag(`openapi:"query"`)
			}
		}
		metadata = append(metadata, shape.FieldMetadata{
			FieldName: name,
			Required:  required,
		})
		tags = append(tags, tag)
	}

	retval := shape.Struct{
		Info: &shape.Info{
			Name:    "", // not ref
			Kind:    shape.Kind(reflect.Struct),
			Package: fn.Info.Package,
		},
		Fields:   fields,
		Tags:     tags,
		Metadata: metadata,
	}
	retval.ResetReflectType(reflect.PtrTo(fn.GetReflectType()))
	return retval
}

type FirstParamOutputSelector struct{}

func (s *FirstParamOutputSelector) SelectOutput(fn shape.Function) shape.Shape {
	if len(fn.Returns.Values) == 0 {
		return nil
	}
	return fn.Returns.Values[0]
}
