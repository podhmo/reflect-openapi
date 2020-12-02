package reflectopenapi

import "github.com/podhmo/reflect-openapi/pkg/shape"

type DefaultSelector struct {
	UseFirstInputSelector
	UseFirstOutputSelector
}

type UseFirstInputSelector struct{}

func (s *UseFirstInputSelector) SelectInput(fn shape.Function) shape.Shape {
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

type UseFirstOutputSelector struct{}

func (s *UseFirstOutputSelector) SelectOutput(fn shape.Function) shape.Shape {
	if len(fn.Returns.Values) == 0 {
		return nil
	}
	return fn.Returns.Values[0]
}
