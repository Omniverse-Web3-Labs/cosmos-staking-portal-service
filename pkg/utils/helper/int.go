package helper

import (
	"reflect"
	"strconv"
)

func IntVal(i interface{}) int {
	var k int
	if v, ok := i.(string); ok {
		k, _ = strconv.Atoi(v)
	} else if v, ok := (i.(bool)); ok {
		if v {
			k = 1
		}
	} else if v, ok := (i.(int)); ok {
		k = v
	} else if v, ok := (i.(int8)); ok {
		k = int(v)
	} else if v, ok := (i.(int16)); ok {
		k = int(v)
	} else if v, ok := (i.(int32)); ok {
		k = int(v)
	} else if v, ok := (i.(int64)); ok {
		k = int(v)
	} else if v, ok := (i.(uint)); ok {
		k = int(v)
	} else if v, ok := (i.(uint8)); ok {
		k = int(v)
	} else if v, ok := (i.(uint16)); ok {
		k = int(v)
	} else if v, ok := (i.(uint32)); ok {
		k = int(v)
	} else if v, ok := (i.(uint64)); ok {
		k = int(v)
	} else if v, ok := (i.(float32)); ok {
		k = int(v)
	} else if v, ok := (i.(float64)); ok {
		k = int(v)
	} else {
		//TODO
	}
	return k
}

func ABS(i interface{}) int {
	k := IntVal(i)
	if k < 0 {
		k = -k
	}
	return k
}

func IsIntegerType(v any) bool {
	t := reflect.TypeOf(v)
	if t != nil {
		kind := t.Kind()
		return kind == reflect.Int || kind == reflect.Int8 || kind == reflect.Int16 ||
			kind == reflect.Int32 || kind == reflect.Int64 || kind == reflect.Uint ||
			kind == reflect.Uint8 || kind == reflect.Uint16 || kind == reflect.Uint32 ||
			kind == reflect.Uint64
	}
	return false
}
