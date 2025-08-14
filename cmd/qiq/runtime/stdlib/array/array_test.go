package array

import (
	"QIQ/cmd/qiq/runtime"
	"QIQ/cmd/qiq/runtime/values"
	"testing"
)

// -------------------------------------- array_first -------------------------------------- MARK: array_first

func TestArrayFirst(t *testing.T) {
	context := runtime.NewContext(nil, nil, nil)

	array := values.NewArray()
	actual, err := nativeFn_array_first([]values.RuntimeValue{array}, context)
	if err != nil {
		t.Errorf("Unexpected error: %s", err)
	} else if actual.GetType() != values.NullValue {
		t.Errorf("Expected: \"null\", Got \"%s\"", actual)
	}

	array.SetElement(nil, values.NewInt(42))
	actual, err = nativeFn_array_first([]values.RuntimeValue{array}, context)
	if err != nil {
		t.Errorf("Unexpected error: %s", err)
	} else if actual.(*values.Int).Value != 42 {
		t.Errorf("Expected: 42, Got %d", actual.(*values.Int).Value)
	}

	array = values.NewArray()
	array.SetElement(values.NewInt(10), values.NewInt(20))
	array.SetElement(values.NewInt(11), values.NewInt(21))
	actual, err = nativeFn_array_first([]values.RuntimeValue{array}, context)
	if err != nil {
		t.Errorf("Unexpected error: %s", err)
	} else if actual.(*values.Int).Value != 20 {
		t.Errorf("Expected: 20, Got %d", actual.(*values.Int).Value)
	}
}

// -------------------------------------- array_key_exists -------------------------------------- MARK: array_key_exists

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

// -------------------------------------- array_key_first -------------------------------------- MARK: array_key_first

func TestArrayKeyFirst(t *testing.T) {
	array := values.NewArray()
	if actual := FirstKey(array); actual.GetType() != values.NullValue {
		t.Errorf("Expected: \"null\", Got \"%s\"", actual)
	}

	array.SetElement(nil, values.NewInt(42))
	if actual := FirstKey(array); actual.(*values.Int).Value != 0 {
		t.Errorf("Expected: 0, Got %d", actual.(*values.Int).Value)
	}

	array = values.NewArray()
	array.SetElement(values.NewInt(42), values.NewInt(43))
	if actual := FirstKey(array); actual.(*values.Int).Value != 42 {
		t.Errorf("Expected: 42, Got %d", actual.(*values.Int).Value)
	}

	array = values.NewArray()
	array.SetElement(values.NewStr("str"), values.NewInt(43))
	if actual := FirstKey(array); actual.(*values.Str).Value != "str" {
		t.Errorf("Expected: \"str\", Got \"%s\"", actual.(*values.Str).Value)
	}
}

// -------------------------------------- array_key_last -------------------------------------- MARK: array_key_last

func TestArrayKeyLast(t *testing.T) {
	array := values.NewArray()
	if actual := LastKey(array); actual.GetType() != values.NullValue {
		t.Errorf("Expected: \"null\", Got \"%s\"", actual)
	}

	array.SetElement(nil, values.NewInt(42))
	array.SetElement(nil, values.NewInt(43))
	if actual := LastKey(array); actual.(*values.Int).Value != 1 {
		t.Errorf("Expected: 0, Got %d", actual.(*values.Int).Value)
	}

	array = values.NewArray()
	array.SetElement(values.NewInt(42), values.NewInt(43))
	if actual := LastKey(array); actual.(*values.Int).Value != 42 {
		t.Errorf("Expected: 42, Got %d", actual.(*values.Int).Value)
	}

	array = values.NewArray()
	array.SetElement(nil, values.NewInt(42))
	array.SetElement(values.NewStr("str"), values.NewInt(43))
	if actual := LastKey(array); actual.(*values.Str).Value != "str" {
		t.Errorf("Expected: \"str\", Got \"%s\"", actual.(*values.Str).Value)
	}
}

// -------------------------------------- array_last -------------------------------------- MARK: array_last

func TestArrayLast(t *testing.T) {
	context := runtime.NewContext(nil, nil, nil)

	array := values.NewArray()
	actual, err := nativeFn_array_last([]values.RuntimeValue{array}, context)
	if err != nil {
		t.Errorf("Unexpected error: %s", err)
	} else if actual.GetType() != values.NullValue {
		t.Errorf("Expected: \"null\", Got \"%s\"", actual)
	}

	array.SetElement(nil, values.NewInt(42))
	actual, err = nativeFn_array_last([]values.RuntimeValue{array}, context)
	if err != nil {
		t.Errorf("Unexpected error: %s", err)
	} else if actual.(*values.Int).Value != 42 {
		t.Errorf("Expected: 42, Got %d", actual.(*values.Int).Value)
	}

	array = values.NewArray()
	array.SetElement(values.NewInt(10), values.NewInt(20))
	array.SetElement(values.NewInt(11), values.NewInt(21))
	actual, err = nativeFn_array_last([]values.RuntimeValue{array}, context)
	if err != nil {
		t.Errorf("Unexpected error: %s", err)
	} else if actual.(*values.Int).Value != 21 {
		t.Errorf("Expected: 21, Got %d", actual.(*values.Int).Value)
	}
}

// -------------------------------------- array_pop -------------------------------------- MARK: array_pop

func TestArrayPop(t *testing.T) {
	context := runtime.NewContext(nil, nil, nil)

	array := values.NewArray()
	actual, err := nativeFn_array_pop([]values.RuntimeValue{array}, context)
	if err != nil {
		t.Errorf("Unexpected error: %s", err)
	} else if actual.GetType() != values.NullValue {
		t.Errorf("Expected: \"null\", Got \"%s\"", actual)
	}

	array.SetElement(nil, values.NewInt(42))
	actual, err = nativeFn_array_pop([]values.RuntimeValue{array}, context)
	if err != nil {
		t.Errorf("Unexpected error: %s", err)
	} else if actual.(*values.Int).Value != 42 {
		t.Errorf("Expected: 42, Got %d", actual.(*values.Int).Value)
	}
	if !array.IsEmpty() {
		t.Error("Expected array to be empty after pop")
	}

	array = values.NewArray()
	array.SetElement(nil, values.NewInt(42))
	array.SetElement(nil, values.NewInt(43))
	actual, err = nativeFn_array_pop([]values.RuntimeValue{array}, context)
	if err != nil {
		t.Errorf("Unexpected error: %s", err)
	} else if actual.(*values.Int).Value != 43 {
		t.Errorf("Expected: 43, Got %d", actual.(*values.Int).Value)
	}
	if len(array.Keys) != 1 || array.Keys[0].(*values.Int).Value != 0 {
		t.Error("Expected array to contain one element with key 0 after pop")
	}
}

// -------------------------------------- array_push -------------------------------------- MARK: array_push

func TestArrayPush(t *testing.T) {
	context := runtime.NewContext(nil, nil, nil)

	array := values.NewArray()
	actual, err := nativeFn_array_push([]values.RuntimeValue{array, values.NewInt(42)}, context)
	if err != nil {
		t.Errorf("Unexpected error: %s", err)
	} else if actual.(*values.Int).Value != 1 {
		t.Errorf("Expected: 1, Got %d", actual.(*values.Int).Value)
	}
	if array.IsEmpty() {
		t.Error("Expected array not to be empty after push")
	}
	if len(array.Keys) != 1 || array.Keys[0].(*values.Int).Value != 0 {
		t.Error("Expected array to contain one element with key 0 after push")
	}

	actual, err = nativeFn_array_push([]values.RuntimeValue{array, values.NewInt(43)}, context)
	if err != nil {
		t.Errorf("Unexpected error: %s", err)
	} else if actual.(*values.Int).Value != 2 {
		t.Errorf("Expected: 2, Got %d", actual.(*values.Int).Value)
	}
	if len(array.Keys) != 2 || array.Keys[1].(*values.Int).Value != 1 {
		t.Error("Expected array to contain one element with key 0 after push")
	}

	array = values.NewArray()
	actual, err = nativeFn_array_push([]values.RuntimeValue{array, values.NewInt(42), values.NewInt(43)}, context)
	if err != nil {
		t.Errorf("Unexpected error: %s", err)
	} else if actual.(*values.Int).Value != 2 {
		t.Errorf("Expected: 2, Got %d", actual.(*values.Int).Value)
	}
	if len(array.Keys) != 2 || array.Keys[0].(*values.Int).Value != 0 || array.Keys[1].(*values.Int).Value != 1 {
		t.Error("Expected array to contain two elements with keys 0 and 1 after push")
	}
}
