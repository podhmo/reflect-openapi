package shape_test

import (
	"fmt"
	"reflect"
	"testing"
	"time"

	"github.com/podhmo/reflect-openapi/pkg/shape"
)

func TestFormat(t *testing.T) {
	e := &shape.Extractor{
		Seen: map[reflect.Type]shape.Shape{},
	}

	t.Run("primitive", func(t *testing.T) {
		var v int = 100
		got := e.Extract(v)
		{
			got := fmt.Sprintf("%+v", got)
			want := "int"
			if want != got {
				t.Errorf("%%+v, want\n\t%s\nbut got\n\t%s", want, got)
			}
		}
		{
			got := fmt.Sprintf("%+#v", got)
			want := "int"
			if want != got {
				t.Errorf("%%+v, want\n\t%s\nbut got\n\t%s", want, got)
			}
		}
	})
	t.Run("struct", func(t *testing.T) {
		var v struct {
			X int
			Y int
			T time.Time
		}
		got := e.Extract(v)
		{
			got := fmt.Sprintf("%+v", got)
			want := "{X, Y, T}"
			if want != got {
				t.Errorf("%%+v, want\n\t%s\nbut got\n\t%s", want, got)
			}
		}
		{
			got := fmt.Sprintf("%+#v", got)
			want := "{X int, Y int, T time.Time}"
			if want != got {
				t.Errorf("%%+v, want\n\t%s\nbut got\n\t%s", want, got)
			}
		}
	})
}
