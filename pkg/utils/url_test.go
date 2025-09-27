package utils_test

import (
	"app/pkg/utils"
	"testing"
)

func TestSuffix(t *testing.T) {
	if suffix := utils.UrlUtil.Suffix("https://www.example.com/abc/test.jpg?version=1"); suffix != "jpg" {
		t.Errorf("suffix https://www.example.com/abc/test.jpg?version=1 is not jpg but is %s", suffix)
	}
	if suffix := utils.UrlUtil.Suffix("https://www.example.com/abc/test.jpg:large"); suffix != "jpg" {
		t.Errorf("suffix https://www.example.com/abc/test.jpg?version=1 is not jpg but is %s", suffix)
	}
}
