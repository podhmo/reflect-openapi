package shape

import "reflect"

func Extract(ob interface{}) Shape {
	e := &Extractor{
		Seen: map[reflect.Type]Shape{},
	}
	return e.Extract(ob)
}

func NewExtractor() *Extractor {
	return &Extractor{
		Seen: map[reflect.Type]Shape{},
	}
}
