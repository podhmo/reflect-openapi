package reflectopenapi

import "github.com/podhmo/reflect-openapi/pkg/shape"

type Selector interface {
	SelectInput(shape.Function) shape.Shape
	SelectOutput(shape.Function) shape.Shape
}

type DefaultSelector struct {
	Extractor *shape.Extractor
}

func (s *DefaultSelector) SelectInput(fn shape.Function) shape.Shape {
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

func (s *DefaultSelector) SelectOutput(fn shape.Function) shape.Shape {
	if len(fn.Returns.Values) == 0 {
		return nil
	}
	return fn.Returns.Values[0]
}
