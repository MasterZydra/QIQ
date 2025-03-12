package array

import (
	"GoPHP/cmd/goPHP/runtime/values"
	"testing"
)

func TestArrayKeyExists(t *testing.T) {
	array := values.NewArray()
	array.SetElement(nil, values.NewInt(42))
	if actual, _ := lib_array_key_exists(values.NewInt(0), array); !actual {
		t.Errorf("Expected: \"%t\", Got \"%t\"", true, actual)
	}
	if actual, _ := lib_array_key_exists(values.NewInt(1), array); actual {
		t.Errorf("Expected: \"%t\", Got \"%t\"", false, actual)
	}
}
