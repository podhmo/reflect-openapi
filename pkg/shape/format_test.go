package shape_test

import (
	"context"
	"fmt"
	"io"
	"reflect"
	"testing"
	"time"

	"github.com/podhmo/reflect-openapi/pkg/shape"
)

func TestFormat(t *testing.T) {
	e := &shape.Extractor{
		Seen: map[reflect.Type]shape.Shape{},
	}

	cases := []struct {
		msg                string
		shape              shape.Shape
		stringFormat       string
		valueFormat        string
		verboseValueFormat string
	}{
		{
			msg: "primitive",
			shape: func() shape.Shape {
				var v int = 100
				return e.Extract(v)
			}(),
			stringFormat:       "int",
			valueFormat:        "int",
			verboseValueFormat: "int",
		},
		{
			msg: "struct",
			shape: func() shape.Shape {
				var v struct {
					v int
					X int
					Y int
					T time.Time
				}
				return e.Extract(v)
			}(),
			stringFormat:       "<anonymous>",
			valueFormat:        "{v, X, Y, T}",
			verboseValueFormat: "{v int, X int, Y int, T time.Time}",
		},
		{
			msg: "interface",
			shape: func() shape.Shape {
				var v = func() io.WriteCloser { return nil }
				return e.Extract(v).(shape.Function).Returns.Values[0]
			}(),
			stringFormat:       "io.WriteCloser",
			valueFormat:        "io.WriteCloser{Close(), Write()}",
			verboseValueFormat: "io.WriteCloser{Close(), Write()}", // TODO: type
		},
		{
			msg: "container, map",
			shape: func() shape.Shape {
				var v = map[string]time.Time{}
				return e.Extract(v)
			}(),
			stringFormat:       "map[string, time.Time]",
			valueFormat:        "map[string, time.Time{wall, ext, loc}]",
			verboseValueFormat: "map[string, time.Time{wall uint64, ext int64, loc *time.Location}]",
		},
		{
			msg: "container, slice",
			shape: func() shape.Shape {
				var v = []*string{}
				return e.Extract(v)
			}(),
			stringFormat:       "slice[*string]",
			valueFormat:        "slice[*string]",
			verboseValueFormat: "slice[*string]",
		},
		{
			msg: "function",
			shape: func() shape.Shape {
				var v = func(ctx context.Context, s string, v interface{}) (shape.Shape, error) { return nil, nil }
				return e.Extract(v)
			}(),
			stringFormat:       "github.com/podhmo/reflect-openapi/pkg/shape_test.TestFormat.func6.1(context.Context, string, ) (github.com/podhmo/reflect-openapi/pkg/shape.Shape, error)",
			valueFormat:        "github.com/podhmo/reflect-openapi/pkg/shape_test.TestFormat.func6.1(context.Context{Deadline(), Done(), Err(), Value}, string, {}) (github.com/podhmo/reflect-openapi/pkg/shape.Shape{Clone(), GetFullName(), GetIdentity(), GetLv(), GetName(), GetPackage(), GetReflectKind(), GetReflectType(), GetReflectValue(), ResetName(), ResetPackage(), ResetReflectType(), Shape(), deref(), info}, error{Error})",
			verboseValueFormat: "github.com/podhmo/reflect-openapi/pkg/shape_test.TestFormat.func6.1(context.Context{Deadline(), Done(), Err(), Value}, string, {}) (github.com/podhmo/reflect-openapi/pkg/shape.Shape{Clone(), GetFullName(), GetIdentity(), GetLv(), GetName(), GetPackage(), GetReflectKind(), GetReflectType(), GetReflectValue(), ResetName(), ResetPackage(), ResetReflectType(), Shape(), deref(), info}, error{Error})",
		},
	}

	for _, c := range cases {
		c := c
		t.Run(c.msg, func(t *testing.T) {
			{
				got := fmt.Sprintf("%+s", c.shape)
				want := c.stringFormat
				if want != got {
					t.Errorf("%%+s, want\n\t%s\nbut got\n\t%s", want, got)
				}
			}
			{
				got := fmt.Sprintf("%+v", c.shape)
				want := c.valueFormat
				if want != got {
					t.Errorf("%%+v, want\n\t%s\nbut got\n\t%s", want, got)
				}
			}
			{
				got := fmt.Sprintf("%+#v", c.shape)
				want := c.verboseValueFormat
				if want != got {
					t.Errorf("%%+v, want\n\t%s\nbut got\n\t%s", want, got)
				}
			}
		})
	}
}
