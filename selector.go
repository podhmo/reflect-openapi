package reflectopenapi

import (
	"fmt"
	"reflect"

	"github.com/podhmo/reflect-openapi/pkg/shape"
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
		if inob.GetFullName() == "context.Background" {
			continue
		}
		return inob
	}
	return nil
}

type MergeParamsInputSelector struct{}

func (s *MergeParamsInputSelector) SelectInput(fn shape.Function) shape.Shape {
	if len(fn.Params.Values) == 0 {
		return nil
	}

	fields := shape.ShapeMap{}
	tags := make([]reflect.StructTag, 0, fn.Params.Len())
	metadata := make([]shape.FieldMetadata, 0, fn.Params.Len())
	for i, p := range fn.Params.Values {
		if p.GetFullName() == "context.Background" {
			continue
		}

		name := fn.Params.Keys[i]
		fields.Keys = append(fields.Keys, name)
		fields.Values = append(fields.Values, p)

		// todo: handling customization
		required := p.GetLv() == 0
		tags = append(tags, reflect.StructTag(
			fmt.Sprintf(`json:"%s" required:"%t"`, name, required),
		))
		metadata = append(metadata, shape.FieldMetadata{Required: required})
	}
	return shape.Struct{
		Info: &shape.Info{
			Name:    "body", // rename?
			Kind:    shape.Kind(reflect.Struct),
			Package: fn.Info.Package,
		},
		Fields:   fields,
		Tags:     tags,
		Metadata: metadata,
	}
}

type FirstParamOutputSelector struct{}

func (s *FirstParamOutputSelector) SelectOutput(fn shape.Function) shape.Shape {
	if len(fn.Returns.Values) == 0 {
		return nil
	}
	return fn.Returns.Values[0]
}
