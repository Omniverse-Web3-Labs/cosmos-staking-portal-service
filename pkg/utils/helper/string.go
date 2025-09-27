package helper

import (
	"strconv"

	jsoniter "github.com/json-iterator/go"
)

func StringVal(i interface{}) string {
	var k string
	if v, ok := i.(string); ok {
		k = v
	} else if v, ok := (i.(bool)); ok {
		if v {
			k = "true"
		} else {
			k = "false"
		}
	} else if v, ok := (i.(int)); ok {
		k = strconv.Itoa(v)
	} else if v, ok := (i.(int8)); ok {
		k = strconv.Itoa(int(v))
	} else if v, ok := (i.(int16)); ok {
		k = strconv.Itoa(int(v))
	} else if v, ok := (i.(int32)); ok {
		k = strconv.Itoa(int(v))
	} else if v, ok := (i.(int64)); ok {
		k = strconv.FormatInt(v, 10)
	} else if v, ok := (i.(uint)); ok {
		k = strconv.FormatUint(uint64(v), 10)
	} else if v, ok := (i.(uint8)); ok {
		k = strconv.FormatUint(uint64(v), 10)
	} else if v, ok := (i.(uint16)); ok {
		k = strconv.FormatUint(uint64(v), 10)
	} else if v, ok := (i.(uint32)); ok {
		k = strconv.FormatUint(uint64(v), 10)
	} else if v, ok := (i.(uint64)); ok {
		k = strconv.FormatUint(uint64(v), 10)
	} else if v, ok := (i.(float32)); ok {
		k = strconv.FormatFloat(float64(v), 'f', 2, 64)
	} else if v, ok := (i.(float64)); ok {
		k = strconv.FormatFloat(v, 'f', 2, 64)
	} else {
		k, _ = jsoniter.MarshalToString(i)
	}
	return k
}
