package helper_test

import (
	"app/pkg/utils/helper"
	"testing"
)

func TestBoolVal(t *testing.T) {
	testCases := map[interface{}]bool{
		0:        false,
		1:        true,
		2:        true,
		uint8(1): true,
		uint8(0): false,
		int8(1):  true,
		int8(0):  false,
		"":       false,
		"1":      true,
	}
	for k, v := range testCases {
		if r := helper.BoolVal(k); v != r {
			t.Errorf("BoolVal(%v) expect to %v but is %v", k, v, r)
		}
	}
}
