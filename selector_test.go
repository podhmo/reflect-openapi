package reflectopenapi

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/podhmo/reflect-openapi/pkg/shape"
)

func dummyTestTargetFunc(struct {
	X int
	Y int
}, struct{ Z int }, int) {
}

func TestFirstParamInputSelector(t *testing.T) {
	selector := &FirstParamInputSelector{}
	fn := dummyTestTargetFunc
	got := selector.SelectInput(shape.Extract(fn).(shape.Function))

	{
		want := reflect.Struct
		got := got.GetReflectKind()
		if want != got {
			t.Errorf("kind:\nwant\n\t%v\nbut got\n\t%v", want, got)
		}
	}
	{
		want := "{X, Y}"
		got := fmt.Sprintf("%+v", got)
		if want != got {
			t.Errorf("format:\nwant\n\t%v\nbut got\n\t%v", want, got)
		}
	}
}

func TestMergeParamsInputSelector(t *testing.T) {
	fn := dummyTestTargetFunc
	selector := &MergeParamsInputSelector{}
	got := selector.SelectInput(shape.Extract(fn).(shape.Function))
	{
		want := reflect.Struct
		got := got.GetReflectKind()
		if want != got {
			t.Errorf("kind:\nwant\n\t%v\nbut got\n\t%v", want, got)
		}
	}
	{
		want := "github.com/podhmo/reflect-openapi.{args0, args1, args2}"
		got := fmt.Sprintf("%+v", got)
		if want != got {
			t.Errorf("format:\nwant\n\t%v\nbut got\n\t%v", want, got)
		}
	}
}
