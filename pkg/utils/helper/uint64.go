package helper

import "strconv"

func Uint64Val(i interface{}) uint64 {
	var k uint64
	if v, ok := i.(string); ok {
		k, _ = strconv.ParseUint(v, 10, 64)
	} else if v, ok := (i.(*string)); ok {
		k, _ = strconv.ParseUint(*v, 10, 64)
	} else if v, ok := (i.(bool)); ok {
		if v {
			k = 1
		}
	} else if v, ok := (i.(int)); ok {
		k = uint64(v)
	} else if v, ok := (i.(int8)); ok {
		k = uint64(v)
	} else if v, ok := (i.(int16)); ok {
		k = uint64(v)
	} else if v, ok := (i.(int32)); ok {
		k = uint64(v)
	} else if v, ok := (i.(int64)); ok {
		k = uint64(v)
	} else if v, ok := (i.(uint)); ok {
		k = uint64(v)
	} else if v, ok := (i.(uint8)); ok {
		k = uint64(v)
	} else if v, ok := (i.(uint16)); ok {
		k = uint64(v)
	} else if v, ok := (i.(uint32)); ok {
		k = uint64(v)
	} else if v, ok := (i.(uint64)); ok {
		k = uint64(v)
	} else if v, ok := (i.(float32)); ok {
		k = uint64(v)
	} else if v, ok := (i.(float64)); ok {
		k = uint64(v)
	} else if v, ok := (i.(*int)); ok {
		k = uint64(*v)
	} else if v, ok := (i.(*int8)); ok {
		k = uint64(*v)
	} else if v, ok := (i.(*int16)); ok {
		k = uint64(*v)
	} else if v, ok := (i.(*int32)); ok {
		k = uint64(*v)
	} else if v, ok := (i.(*int64)); ok {
		k = uint64(*v)
	} else if v, ok := (i.(*uint)); ok {
		k = uint64(*v)
	} else if v, ok := (i.(*uint8)); ok {
		k = uint64(*v)
	} else if v, ok := (i.(*uint16)); ok {
		k = uint64(*v)
	} else if v, ok := (i.(*uint32)); ok {
		k = uint64(*v)
	} else if v, ok := (i.(*uint64)); ok {
		k = uint64(*v)
	} else if v, ok := (i.(*float32)); ok {
		k = uint64(*v)
	} else if v, ok := (i.(*float64)); ok {
		k = uint64(*v)
	}
	return k
}
