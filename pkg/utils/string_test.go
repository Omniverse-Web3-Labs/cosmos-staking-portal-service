package utils_test

import (
	"app/pkg/utils"
	"testing"
)

func TestString(t *testing.T) {
	id := utils.StringUtil.UniqueID()
	if len(id) != 16 {
		t.Errorf("UniqueID expect length to 16 but is %d", len(id))
	}
}
