package interpreter

import (
	"GoPHP/cmd/goPHP/runtime/values"
	"testing"
)

// ------------------- MARK: boolval -------------------

func TestLibBoolval(t *testing.T) {
	doTest := func(runtimeValue values.RuntimeValue, expected bool) {
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
	doTest(values.NewArray(), false)
	array := values.NewArray()
	array.SetElement(nil, values.NewInt(42))
	doTest(array, true)

	// boolean to boolean
	doTest(values.NewBool(true), true)
	doTest(values.NewBool(false), false)

	// integer to boolean
	doTest(values.NewInt(0), false)
	doTest(values.NewInt(-0), false)
	doTest(values.NewInt(1), true)
	doTest(values.NewInt(42), true)
	doTest(values.NewInt(-2), true)

	// floating to boolean
	doTest(values.NewFloat(0.0), false)
	doTest(values.NewFloat(1.5), true)
	doTest(values.NewFloat(42.0), true)
	doTest(values.NewFloat(-2.0), true)

	// string to boolean
	doTest(values.NewStr(""), false)
	doTest(values.NewStr("0"), false)
	doTest(values.NewStr("Hi"), true)

	// null to boolean
	doTest(values.NewNull(), false)
}

// ------------------- MARK: floatval -------------------

func TestLibFloatval(t *testing.T) {
	doTest := func(runtimeValue values.RuntimeValue, expected float64) {
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
	doTest(values.NewBool(true), 1)
	doTest(values.NewBool(false), 0)

	// integer to floating
	doTest(values.NewInt(0), 0)
	doTest(values.NewInt(-0), 0)
	doTest(values.NewInt(1), 1)
	doTest(values.NewInt(42), 42)
	doTest(values.NewInt(-2), -2)

	// floating to floating
	doTest(values.NewFloat(0.0), 0)
	doTest(values.NewFloat(1.5), 1.5)
	doTest(values.NewFloat(42.0), 42)
	doTest(values.NewFloat(-2.0), -2)

	// string to floating
	doTest(values.NewStr(""), 0)
	doTest(values.NewStr("0"), 0)
	doTest(values.NewStr("42"), 42)
	doTest(values.NewStr("+42"), 42)
	doTest(values.NewStr("-42.4"), -42.4)
	doTest(values.NewStr("Hi"), 0)

	// null to floating
	doTest(values.NewNull(), 0)
}

// ------------------- MARK: get_debug_type -------------------

func TestLibGetDebugType(t *testing.T) {
	testInputOutput(t, `<?php echo get_debug_type(false);`, "bool")
	testInputOutput(t, `<?php echo get_debug_type(true);`, "bool")
	testInputOutput(t, `<?php echo get_debug_type(0);`, "int")
	testInputOutput(t, `<?php echo get_debug_type(-1);`, "int")
	testInputOutput(t, `<?php echo get_debug_type(42);`, "int")
	testInputOutput(t, `<?php echo get_debug_type(0.0);`, "float")
	testInputOutput(t, `<?php echo get_debug_type(-1.5);`, "float")
	testInputOutput(t, `<?php echo get_debug_type(42.5);`, "float")
	testInputOutput(t, `<?php echo get_debug_type("");`, "string")
	testInputOutput(t, `<?php echo get_debug_type("abc");`, "string")
	testInputOutput(t, `<?php echo get_debug_type([]);`, "array")
	testInputOutput(t, `<?php echo get_debug_type([42]);`, "array")
	testInputOutput(t, `<?php echo get_debug_type(null);`, "null")
}

// ------------------- MARK: gettype -------------------

func TestLibGettype(t *testing.T) {
	testInputOutput(t, `<?php echo gettype(false);`, "boolean")
	testInputOutput(t, `<?php echo gettype(true);`, "boolean")
	testInputOutput(t, `<?php echo gettype(0);`, "integer")
	testInputOutput(t, `<?php echo gettype(-1);`, "integer")
	testInputOutput(t, `<?php echo gettype(42);`, "integer")
	testInputOutput(t, `<?php echo gettype(0.0);`, "double")
	testInputOutput(t, `<?php echo gettype(-1.5);`, "double")
	testInputOutput(t, `<?php echo gettype(42.5);`, "double")
	testInputOutput(t, `<?php echo gettype("");`, "string")
	testInputOutput(t, `<?php echo gettype("abc");`, "string")
	testInputOutput(t, `<?php echo gettype([]);`, "array")
	testInputOutput(t, `<?php echo gettype([42]);`, "array")
	testInputOutput(t, `<?php echo gettype(null);`, "NULL")
}

// ------------------- MARK: intval -------------------

func TestLibIntval(t *testing.T) {
	doTest := func(runtimeValue values.RuntimeValue, expected int64) {
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
	doTest(values.NewArray(), 0)
	array := values.NewArray()
	array.SetElement(nil, values.NewInt(42))
	doTest(array, 1)

	// boolean to integer
	doTest(values.NewBool(true), 1)
	doTest(values.NewBool(false), 0)

	// integer to integer
	doTest(values.NewInt(0), 0)
	doTest(values.NewInt(-0), 0)
	doTest(values.NewInt(1), 1)
	doTest(values.NewInt(42), 42)
	doTest(values.NewInt(-2), -2)

	// floating to integer
	doTest(values.NewFloat(0.0), 0)
	doTest(values.NewFloat(1.5), 1)
	doTest(values.NewFloat(42.0), 42)
	doTest(values.NewFloat(-2.0), -2)

	// string to integer
	doTest(values.NewStr(""), 0)
	doTest(values.NewStr("0"), 0)
	doTest(values.NewStr("1"), 1)
	doTest(values.NewStr("42"), 42)
	doTest(values.NewStr("+42"), 42)
	doTest(values.NewStr("-42"), -42)
	doTest(values.NewStr("Hi"), 0)

	// null to integer
	doTest(values.NewNull(), 0)
}

// ------------------- MARK: is_XXX -------------------

func TestLibIsType(t *testing.T) {
	// is_array
	testInputOutput(t, `<?php $a = [true]; var_dump(is_array($a));`, "bool(true)\n")
	testInputOutput(t, `<?php $a = 0; var_dump(is_array($a));`, "bool(false)\n")

	// is_bool
	testInputOutput(t, `<?php $a = true; var_dump(is_bool($a));`, "bool(true)\n")
	testInputOutput(t, `<?php $a = 0; var_dump(is_bool($a));`, "bool(false)\n")

	// is_float
	testInputOutput(t, `<?php $a = 42.0; var_dump(is_float($a));`, "bool(true)\n")
	testInputOutput(t, `<?php $a = 0; var_dump(is_float($a));`, "bool(false)\n")

	// is_int
	testInputOutput(t, `<?php $a = 42; var_dump(is_int($a));`, "bool(true)\n")
	testInputOutput(t, `<?php $a = "42"; var_dump(is_int($a));`, "bool(false)\n")

	// is_null
	testInputOutput(t, `<?php $a = null; var_dump(is_null($a));`, "bool(true)\n")
	testInputOutput(t, `<?php $a = 42; var_dump(is_null($a));`, "bool(false)\n")

	// is_scalar
	testInputOutput(t, `<?php $a = true; var_dump(is_scalar($a));`, "bool(true)\n")
	testInputOutput(t, `<?php $a = false; var_dump(is_scalar($a));`, "bool(true)\n")
	testInputOutput(t, `<?php $a = 42; var_dump(is_scalar($a));`, "bool(true)\n")
	testInputOutput(t, `<?php $a = 3.5; var_dump(is_scalar($a));`, "bool(true)\n")
	testInputOutput(t, `<?php $a = "abc"; var_dump(is_scalar($a));`, "bool(true)\n")
	testInputOutput(t, `<?php $a = null; var_dump(is_scalar($a));`, "bool(false)\n")
	testInputOutput(t, `<?php $a = []; var_dump(is_scalar($a));`, "bool(false)\n")

	// is_string
	testInputOutput(t, `<?php $a = " "; var_dump(is_string($a));`, "bool(true)\n")
	testInputOutput(t, `<?php $a = 42; var_dump(is_string($a));`, "bool(false)\n")
}

// ------------------- MARK: print_r -------------------

func TestLibPrintR(t *testing.T) {
	testInputOutput(t, `<?php print_r(3.5);`, "3.5")
	testInputOutput(t, `<?php print_r(42);`, "42")
	testInputOutput(t, `<?php print_r("abc");`, "abc")
	testInputOutput(t, `<?php print_r(true);`, "1")
	testInputOutput(t, `<?php print_r(false);`, "")
	testInputOutput(t, `<?php print_r(null);`, "")
	testInputOutput(t, `<?php print_r([]);`, "Array\n(\n)\n")
	testInputOutput(t, `<?php print_r([1,2]);`, "Array\n(\n    [0] => 1\n    [1] => 2\n)\n")
	testInputOutput(t, `<?php print_r([1, [1]]);`, "Array\n(\n    [0] => 1\n    [1] => Array\n        (\n            [0] => 1\n        )\n\n)\n")
}

// ------------------- MARK: strval -------------------

func TestLibStrval(t *testing.T) {
	doTest := func(runtimeValue values.RuntimeValue, expected string) {
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
	doTest(values.NewArray(), "Array")

	// boolean to string
	doTest(values.NewBool(true), "1")
	doTest(values.NewBool(false), "")

	// integer to string
	doTest(values.NewInt(0), "0")
	doTest(values.NewInt(-0), "0")
	doTest(values.NewInt(1), "1")
	doTest(values.NewInt(42), "42")
	doTest(values.NewInt(-2), "-2")

	// floating to string
	doTest(values.NewFloat(0.0), "0")
	doTest(values.NewFloat(1.5), "1.5")
	doTest(values.NewFloat(42.0), "42")
	doTest(values.NewFloat(-2.0), "-2")

	// string to string
	doTest(values.NewStr(""), "")
	doTest(values.NewStr("0"), "0")
	doTest(values.NewStr("Hi"), "Hi")

	// null to string
	doTest(values.NewNull(), "")
}

// ------------------- MARK: var_dump -------------------

func TestLibVarDump(t *testing.T) {
	testInputOutput(t, `<?php var_dump(3.5);`, "float(3.5)\n")
	testInputOutput(t, `<?php var_dump(3.5, 42, true, false, null);`, "float(3.5)\nint(42)\nbool(true)\nbool(false)\nNULL\n")
	testInputOutput(t, `<?php var_dump([]);`, "array(0) {\n}\n")
	testInputOutput(t, `<?php var_dump([1,2]);`, "array(2) {\n  [0]=>\n  int(1)\n  [1]=>\n  int(2)\n}\n")
	testInputOutput(t, `<?php var_dump([1, [1]]);`, "array(2) {\n  [0]=>\n  int(1)\n  [1]=>\n  array(1) {\n    [0]=>\n    int(1)\n  }\n}\n")
}

// ------------------- MARK: var_export -------------------

func TestLibVarExport(t *testing.T) {
	testInputOutput(t, `<?php var_export(3.5);`, "3.5")
	testInputOutput(t, `<?php var_export(42);`, "42")
	testInputOutput(t, `<?php var_export("abc");`, "'abc'")
	testInputOutput(t, `<?php var_export(true);`, "true")
	testInputOutput(t, `<?php var_export(false);`, "false")
	testInputOutput(t, `<?php var_export(null);`, "NULL")
	testInputOutput(t, `<?php var_export([]);`, "array (\n)")
	testInputOutput(t, `<?php var_export([1,2]);`, "array (\n  0 => 1,\n  1 => 2,\n)")
	testInputOutput(t, `<?php var_export([1, [1]]);`, "array (\n  0 => 1,\n  1 => \n  array (\n    0 => 1,\n  ),\n)")
}
