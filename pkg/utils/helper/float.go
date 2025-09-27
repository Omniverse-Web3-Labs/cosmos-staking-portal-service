package helper

import "reflect"

func IsFloatType(v any) bool {
	t := reflect.TypeOf(v)
	if t != nil {
		kind := t.Kind()
		return kind == reflect.Float32 || kind == reflect.Float64
	}
	return false
}
