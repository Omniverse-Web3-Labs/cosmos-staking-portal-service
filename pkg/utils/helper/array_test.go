package helper_test

import (
	"app/pkg/utils/helper"
	"testing"
)

func TestInArray(t *testing.T) {
	ok := helper.InArray([]string{"1", "2", "3"}, "2")
	if ok != true {
		t.Errorf("expect %v, but is %v", true, ok)
	}
	ok = helper.InArray([]string{"1", "2", "3"}, "4")
	if ok != false {
		t.Errorf("expect %v, but is %v", false, ok)
	}
	ok = helper.InArray([]int{1, 2, 3}, 2)
	if ok != true {
		t.Errorf("expect %v, but is %v", true, ok)
	}
	ok = helper.InArray([]int{1, 2, 3}, 4)
	if ok != false {
		t.Errorf("expect %v, but is %v", false, ok)
	}
}
