package jsonequal

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/google/go-cmp/cmp"
)

type Node struct {
	name string
	v    interface{}
	b    []byte
	err  error
}

func (n *Node) Named(v string) *Node {
	n.name = v
	return n
}

// From :
func From(iface interface{}) *Node {
	b, err := json.Marshal(iface)
	if err != nil {
		return &Node{"", nil, nil, err}
	}
	var v interface{}
	if err := json.Unmarshal(b, &v); err != nil {
		return &Node{"", nil, nil, err}
	}
	return &Node{"", v, b, nil}
}

// FromRaw :
func FromRaw(iface interface{}) *Node {
	b, err := json.Marshal(iface)
	if err != nil {
		return &Node{"", nil, nil, err}
	}
	return &Node{"", iface, b, nil}
}

// FromRawWithBytes :
func FromRawWithBytes(iface interface{}, b []byte) *Node {
	return &Node{"", iface, b, nil}
}

// FromReader :
func FromReader(reader io.Reader) *Node {

	decoder := json.NewDecoder(reader)
	var v interface{}
	if err := decoder.Decode(&v); err != nil {
		return &Node{"", nil, nil, err}
	}
	b, err := json.Marshal(&v)
	if err != nil {
		return &Node{"", nil, nil, err}
	}
	return &Node{"", v, b, nil}
}

// FromReadCloser :
func FromReadCloser(reader io.ReadCloser) *Node {
	defer reader.Close()
	return FromReader(reader)
}

// FromBytes :
func FromBytes(b []byte) *Node {
	var v interface{}
	if err := json.Unmarshal(b, &v); err != nil {
		return &Node{"", nil, nil, err}
	}
	return &Node{"", v, b, nil}
}

// FromString :
func FromString(s string) *Node {
	return FromBytes([]byte(s))
}

// NoDiff :
func NoDiff(
	l *Node,
	r *Node,
) error {
	if l.name == "" {
		l.name = "left"
	}
	if r.name == "" {
		r.name = "right"
	}
	wrapfFunc := func(err error, message string) error {
		return fmt.Errorf("%s: %w", message, err)
	}

	if l.err != nil {
		return wrapfFunc(l.err, "on load left data")
	}
	if r.err != nil {
		return wrapfFunc(r.err, "on load right data")
	}

	diff := cmp.Diff(l.v, r.v)
	if diff == "" {
		return nil
	}
	return fmt.Errorf("JSON mismatch (-%s +%s):\n%s", l.name, r.name, diff)
}

// Equal :
func Equal(
	lsrc *Node,
	rsrc *Node,
) bool {
	return NoDiff(lsrc, rsrc) == nil
}
