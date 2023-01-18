package jsonequal

// MustNormalize :
func MustNormalize(src interface{}) interface{} {
	dst, err := Normalize(src)
	if err != nil {
		panic(err)
	}
	return dst
}

// Normalize :
func Normalize(src interface{}) (interface{}, error) {
	dst, _, err := From(src)()
	if err != nil {
		return nil, err
	}
	return dst, nil
}
