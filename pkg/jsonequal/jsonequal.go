package jsonequal

import (
	"encoding/json"
	"fmt"
	"io"
	"reflect"

	"github.com/google/go-cmp/cmp"
)

// Caller :
type Caller struct {
	LeftName  string
	RightName string

	EqualFunc func(left interface{}, right interface{}) bool
	WrapfFunc func(error, string) error
}

// From :
func From(iface interface{}) func() (interface{}, []byte, error) {
	return func() (interface{}, []byte, error) {
		b, err := json.Marshal(iface)
		if err != nil {
			return nil, nil, err
		}
		var v interface{}
		if err := json.Unmarshal(b, &v); err != nil {
			return nil, nil, err
		}
		return v, b, nil
	}
}

// FromRaw :
func FromRaw(iface interface{}) func() (interface{}, []byte, error) {
	return func() (interface{}, []byte, error) {
		b, err := json.Marshal(iface)
		if err != nil {
			return nil, nil, err
		}
		return iface, b, nil
	}
}

// FromRawWithBytes :
func FromRawWithBytes(iface interface{}, b []byte) func() (interface{}, []byte, error) {
	return func() (interface{}, []byte, error) {
		return iface, b, nil
	}
}

// FromReader :
func FromReader(reader io.Reader) func() (interface{}, []byte, error) {
	return func() (interface{}, []byte, error) {
		decoder := json.NewDecoder(reader)
		var v interface{}
		if err := decoder.Decode(&v); err != nil {
			return nil, nil, err
		}
		b, err := json.Marshal(&v)
		if err != nil {
			return nil, nil, err
		}
		return v, b, nil
	}
}

// FromReadCloser :
func FromReadCloser(reader io.ReadCloser) func() (interface{}, []byte, error) {
	return func() (interface{}, []byte, error) {
		defer reader.Close()
		return FromReader(reader)()
	}
}

// FromBytes :
func FromBytes(b []byte) func() (interface{}, []byte, error) {
	return func() (interface{}, []byte, error) {
		var v interface{}
		if err := json.Unmarshal(b, &v); err != nil {
			return nil, nil, err
		}
		return v, b, nil
	}
}

// FromString :
func FromString(s string) func() (interface{}, []byte, error) {
	return FromBytes([]byte(s))
}

// WithLeftName :
func WithLeftName(s string) func(*Caller) {
	return func(c *Caller) {
		c.LeftName = s
	}
}

// WithRightName :
func WithRightName(s string) func(*Caller) {
	return func(c *Caller) {
		c.RightName = s
	}
}

// ShouldBeSame :
func ShouldBeSame(
	lsrc func() (interface{}, []byte, error),
	rsrc func() (interface{}, []byte, error),
	options ...func(*Caller),
) error {
	caller := Caller{
		LeftName:  "left",
		RightName: "right",
	}
	for _, opt := range options {
		opt(&caller)
	}
	if caller.EqualFunc == nil {
		caller.EqualFunc = reflect.DeepEqual
	}
	if caller.WrapfFunc == nil {
		caller.WrapfFunc = func(err error, message string) error {
			return fmt.Errorf("%s: %w", message, err)
		}
	}

	lv, lb, err := lsrc()
	if err != nil {
		return caller.WrapfFunc(err, "on load left data")
	}
	rv, rb, err := rsrc()
	if err != nil {
		return caller.WrapfFunc(err, "on load right data")
	}

	_ = lb
	_ = rb
	diff := cmp.Diff(lv, rv)
	if diff == "" {
		return nil
	}
	return fmt.Errorf("JSON mismatch (-%s +%s):\n%s", caller.LeftName, caller.RightName, diff)
}

// Equal :
func Equal(
	lsrc func() (interface{}, []byte, error),
	rsrc func() (interface{}, []byte, error),
	options ...func(*Caller),
) bool {
	return ShouldBeSame(lsrc, rsrc, options...) == nil
}
