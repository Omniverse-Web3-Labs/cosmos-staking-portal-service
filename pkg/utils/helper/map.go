package helper

import jsoniter "github.com/json-iterator/go"

func StructToMap(a interface{}) map[string]interface{} {
	data := make(map[string]interface{}, 0)
	bs, _ := jsoniter.Marshal(a)
	jsoniter.Unmarshal(bs, &data)
	return data
}

func MapMerge(a ...map[string]interface{}) map[string]interface{} {
	data := make(map[string]interface{}, 0)
	for i := range a {
		for k, v := range a[i] {
			data[k] = v
		}
	}
	return data
}
