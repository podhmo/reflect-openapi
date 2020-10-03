package jsonequal_test

import (
	"bytes"
	"testing"

	"github.com/podhmo/reflect-openapi/pkg/jsonequal"
)

func TestIt(t *testing.T) {
	v := map[string]int{"foo": 1}
	b := []byte(`{"foo": 1}`)
	r := bytes.NewBufferString(`{"foo": 1}`)

	if err := jsonequal.ShouldBeSame(jsonequal.From(v), jsonequal.FromBytes(b)); err != nil {
		t.Errorf("mismatch: %s", err)
	}

	if err := jsonequal.ShouldBeSame(jsonequal.From(v), jsonequal.FromReader(r)); err != nil {
		t.Errorf("mismatch: %s", err)
	}
}
