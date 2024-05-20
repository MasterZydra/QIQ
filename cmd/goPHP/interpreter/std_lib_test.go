package interpreter

import "testing"

func TestArrayKeyExists(t *testing.T) {
	array := NewArrayRuntimeValue(map[IRuntimeValue]IRuntimeValue{
		NewIntegerRuntimeValue(0): NewIntegerRuntimeValue(42),
	})
	if actual, _ := lib_array_key_exists(NewIntegerRuntimeValue(0), array); !actual {
		t.Errorf("Expected: \"%t\", Got \"%t\"", true, actual)
	}
	if actual, _ := lib_array_key_exists(NewIntegerRuntimeValue(1), array); actual {
		t.Errorf("Expected: \"%t\", Got \"%t\"", false, actual)
	}
}

func TestLibBoolval(t *testing.T) {
	doTest := func(runtimeValue IRuntimeValue, expected bool) {
		actual, err := lib_boolval(runtimeValue)
		if err != nil {
			t.Errorf("Unexpected error: \"%s\"", err)
			return
		}
		if actual != expected {
			t.Errorf("Expected: \"%t\", Got \"%t\"", expected, actual)
		}
	}

	// array to boolean
	doTest(NewArrayRuntimeValue(map[IRuntimeValue]IRuntimeValue{}), false)
	doTest(NewArrayRuntimeValue(map[IRuntimeValue]IRuntimeValue{NewIntegerRuntimeValue(0): NewIntegerRuntimeValue(42)}), true)

	// boolean to boolean
	doTest(NewBooleanRuntimeValue(true), true)
	doTest(NewBooleanRuntimeValue(false), false)

	// integer to boolean
	doTest(NewIntegerRuntimeValue(0), false)
	doTest(NewIntegerRuntimeValue(-0), false)
	doTest(NewIntegerRuntimeValue(1), true)
	doTest(NewIntegerRuntimeValue(42), true)
	doTest(NewIntegerRuntimeValue(-2), true)

	// floating to boolean
	doTest(NewFloatingRuntimeValue(0.0), false)
	doTest(NewFloatingRuntimeValue(1.5), true)
	doTest(NewFloatingRuntimeValue(42.0), true)
	doTest(NewFloatingRuntimeValue(-2.0), true)

	// string to boolean
	doTest(NewStringRuntimeValue(""), false)
	doTest(NewStringRuntimeValue("0"), false)
	doTest(NewStringRuntimeValue("Hi"), true)

	// null to boolean
	doTest(NewNullRuntimeValue(), false)
}

func TestLibFloatval(t *testing.T) {
	doTest := func(runtimeValue IRuntimeValue, expected float64) {
		actual, err := lib_floatval(runtimeValue)
		if err != nil {
			t.Errorf("Unexpected error: \"%s\"", err)
			return
		}
		if actual != expected {
			t.Errorf("Expected: \"%g\", Got \"%g\"", expected, actual)
		}
	}

	// boolean to floating
	doTest(NewBooleanRuntimeValue(true), 1)
	doTest(NewBooleanRuntimeValue(false), 0)

	// integer to floating
	doTest(NewIntegerRuntimeValue(0), 0)
	doTest(NewIntegerRuntimeValue(-0), 0)
	doTest(NewIntegerRuntimeValue(1), 1)
	doTest(NewIntegerRuntimeValue(42), 42)
	doTest(NewIntegerRuntimeValue(-2), -2)

	// floating to floating
	doTest(NewFloatingRuntimeValue(0.0), 0)
	doTest(NewFloatingRuntimeValue(1.5), 1.5)
	doTest(NewFloatingRuntimeValue(42.0), 42)
	doTest(NewFloatingRuntimeValue(-2.0), -2)

	// string to floating
	// doTest(NewStringRuntimeValue(""), false)
	// doTest(NewStringRuntimeValue("0"), false)
	// doTest(NewStringRuntimeValue("Hi"), true)

	// null to floating
	doTest(NewNullRuntimeValue(), 0)
}

func TestLibIntval(t *testing.T) {
	doTest := func(runtimeValue IRuntimeValue, expected int64) {
		actual, err := lib_intval(runtimeValue)
		if err != nil {
			t.Errorf("Unexpected error: \"%s\"", err)
			return
		}
		if actual != expected {
			t.Errorf("Expected: \"%d\", Got \"%d\"", expected, actual)
		}
	}

	// array to integer
	doTest(NewArrayRuntimeValue(map[IRuntimeValue]IRuntimeValue{}), 0)
	doTest(NewArrayRuntimeValue(map[IRuntimeValue]IRuntimeValue{NewIntegerRuntimeValue(0): NewIntegerRuntimeValue(42)}), 1)

	// boolean to integer
	doTest(NewBooleanRuntimeValue(true), 1)
	doTest(NewBooleanRuntimeValue(false), 0)

	// integer to integer
	doTest(NewIntegerRuntimeValue(0), 0)
	doTest(NewIntegerRuntimeValue(-0), 0)
	doTest(NewIntegerRuntimeValue(1), 1)
	doTest(NewIntegerRuntimeValue(42), 42)
	doTest(NewIntegerRuntimeValue(-2), -2)

	// floating to integer
	// doTest(NewFloatingRuntimeValue(0.0), 0)
	// doTest(NewFloatingRuntimeValue(1.5), 1)
	// doTest(NewFloatingRuntimeValue(42.0), 42)
	// doTest(NewFloatingRuntimeValue(-2.0), -2)

	// string to integer
	// doTest(NewStringRuntimeValue(""), false)
	// doTest(NewStringRuntimeValue("0"), false)
	// doTest(NewStringRuntimeValue("Hi"), true)

	// null to integer
	doTest(NewNullRuntimeValue(), 0)
}

func TestLibIsNull(t *testing.T) {
	actual := lib_is_null(NewNullRuntimeValue())
	if expected := true; actual != expected {
		t.Errorf("Expected: \"%t\", Got \"%t\"", expected, actual)
	}
	actual = lib_is_null(NewIntegerRuntimeValue(42))
	if expected := false; actual != expected {
		t.Errorf("Expected: \"%t\", Got \"%t\"", expected, actual)
	}
}

func TestLibStrval(t *testing.T) {
	doTest := func(runtimeValue IRuntimeValue, expected string) {
		actual, err := lib_strval(runtimeValue)
		if err != nil {
			t.Errorf("Unexpected error: \"%s\"", err)
			return
		}
		if actual != expected {
			t.Errorf("Expected: \"%s\", Got \"%s\"", expected, actual)
		}
	}

	// array to string
	doTest(NewArrayRuntimeValue(map[IRuntimeValue]IRuntimeValue{}), "Array")

	// boolean to string
	doTest(NewBooleanRuntimeValue(true), "1")
	doTest(NewBooleanRuntimeValue(false), "")

	// integer to string
	doTest(NewIntegerRuntimeValue(0), "0")
	doTest(NewIntegerRuntimeValue(-0), "0")
	doTest(NewIntegerRuntimeValue(1), "1")
	doTest(NewIntegerRuntimeValue(42), "42")
	doTest(NewIntegerRuntimeValue(-2), "-2")

	// floating to string
	doTest(NewFloatingRuntimeValue(0.0), "0")
	doTest(NewFloatingRuntimeValue(1.5), "1.5")
	doTest(NewFloatingRuntimeValue(42.0), "42")
	doTest(NewFloatingRuntimeValue(-2.0), "-2")

	// string to string
	doTest(NewStringRuntimeValue(""), "")
	doTest(NewStringRuntimeValue("0"), "0")
	doTest(NewStringRuntimeValue("Hi"), "Hi")

	// null to string
	doTest(NewNullRuntimeValue(), "")
}
