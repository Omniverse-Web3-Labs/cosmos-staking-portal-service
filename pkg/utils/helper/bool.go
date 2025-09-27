package helper

import "reflect"

func BoolVal(i interface{}) bool {
	if v, ok := i.(int); ok {
		return v != 0
	} else if v, ok := i.(int8); ok {
		return v != 0
	} else if v, ok := i.(int16); ok {
		return v != 0
	} else if v, ok := i.(int32); ok {
		return v != 0
	} else if v, ok := i.(int64); ok {
		return v != 0
	} else if v, ok := i.(uint); ok {
		return v != 0
	} else if v, ok := i.(uint8); ok {
		return v != 0
	} else if v, ok := i.(uint16); ok {
		return v != 0
	} else if v, ok := i.(uint32); ok {
		return v != 0
	} else if v, ok := i.(uint64); ok {
		return v != 0
	} else if v, ok := i.(string); ok {
		return v != ""
	} else if v, ok := i.(bool); ok {
		return v
	} else if v, ok := i.(float32); ok {
		return v != float32(0)
	} else if v, ok := i.(float64); ok {
		return v != float64(0)
	} else {
		v := reflect.ValueOf(i)
		switch v.Kind() {
		case reflect.Array, reflect.Map, reflect.Slice:
			return v.Len() == 0
		case reflect.Pointer, reflect.Interface:
			return v.IsNil()
		case reflect.Struct:
			return v.Interface() == struct{}{}
		}
	}
	return false
}
