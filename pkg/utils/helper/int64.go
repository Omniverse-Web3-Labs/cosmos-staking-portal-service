package helper

import "strconv"

func Int64Val(i interface{}) int64 {
	var k int64
	if v, ok := i.(string); ok {
		k, _ = strconv.ParseInt(v, 10, 64)
	} else if v, ok := (i.(*string)); ok {
		k, _ = strconv.ParseInt(*v, 10, 64)
	} else if v, ok := (i.(bool)); ok {
		if v {
			k = 1
		}
	} else if v, ok := (i.(int)); ok {
		k = int64(v)
	} else if v, ok := (i.(int8)); ok {
		k = int64(v)
	} else if v, ok := (i.(int16)); ok {
		k = int64(v)
	} else if v, ok := (i.(int32)); ok {
		k = int64(v)
	} else if v, ok := (i.(int64)); ok {
		k = int64(v)
	} else if v, ok := (i.(uint)); ok {
		k = int64(v)
	} else if v, ok := (i.(uint8)); ok {
		k = int64(v)
	} else if v, ok := (i.(uint16)); ok {
		k = int64(v)
	} else if v, ok := (i.(uint32)); ok {
		k = int64(v)
	} else if v, ok := (i.(uint64)); ok {
		k = int64(v)
	} else if v, ok := (i.(float32)); ok {
		k = int64(v)
	} else if v, ok := (i.(float64)); ok {
		k = int64(v)
	} else if v, ok := (i.(*int)); ok {
		k = int64(*v)
	} else if v, ok := (i.(*int8)); ok {
		k = int64(*v)
	} else if v, ok := (i.(*int16)); ok {
		k = int64(*v)
	} else if v, ok := (i.(*int32)); ok {
		k = int64(*v)
	} else if v, ok := (i.(*int64)); ok {
		k = int64(*v)
	} else if v, ok := (i.(*uint)); ok {
		k = int64(*v)
	} else if v, ok := (i.(*uint8)); ok {
		k = int64(*v)
	} else if v, ok := (i.(*uint16)); ok {
		k = int64(*v)
	} else if v, ok := (i.(*uint32)); ok {
		k = int64(*v)
	} else if v, ok := (i.(*uint64)); ok {
		k = int64(*v)
	} else if v, ok := (i.(*float32)); ok {
		k = int64(*v)
	} else if v, ok := (i.(*float64)); ok {
		k = int64(*v)
	}
	return k
}
