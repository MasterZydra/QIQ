package interpreter

import (
	"GoPHP/cmd/goPHP/ast"
	"GoPHP/cmd/goPHP/ini"
	"GoPHP/cmd/goPHP/phpError"
	"fmt"
	"testing"
)

// ------------------- MARK: function tests -------------------

func TestVariableExprToVariableName(t *testing.T) {
	// simple-variable-expression

	// $var
	interpreter := NewInterpreter(ini.NewDevIni(), &Request{}, "test.php")
	actual, err := interpreter.varExprToVarName(ast.NewSimpleVariableExpr(0, ast.NewVariableNameExpr(0, nil, "$var")), interpreter.env)
	if err != nil {
		t.Errorf("Unexpected error: \"%s\"", err)
		return
	}
	expected := "$var"
	if actual != expected {
		t.Errorf("Expected: \"%s\", Got \"%s\"", expected, actual)
	}

	// $$var
	interpreter = NewInterpreter(ini.NewDevIni(), &Request{}, "test.php")
	interpreter.env.declareVariable("$var", NewStringRuntimeValue("hi"))
	actual, err = interpreter.varExprToVarName(
		ast.NewSimpleVariableExpr(0,
			ast.NewSimpleVariableExpr(0, ast.NewVariableNameExpr(0, nil, "$var"))), interpreter.env)
	if err != nil {
		t.Errorf("Unexpected error: \"%s\"", err)
		return
	}
	expected = "$hi"
	if actual != expected {
		t.Errorf("Expected: \"%s\", Got \"%s\"", expected, actual)
	}

	// $$$var
	interpreter = NewInterpreter(ini.NewDevIni(), &Request{}, "test.php")
	interpreter.env.declareVariable("$var1", NewStringRuntimeValue("hi"))
	interpreter.env.declareVariable("$var", NewStringRuntimeValue("var1"))
	actual, err = interpreter.varExprToVarName(
		ast.NewSimpleVariableExpr(0,
			ast.NewSimpleVariableExpr(0,
				ast.NewSimpleVariableExpr(0, ast.NewVariableNameExpr(0, nil, "$var")))), interpreter.env)
	if err != nil {
		t.Errorf("Unexpected error: \"%s\"", err)
		return
	}
	expected = "$hi"
	if actual != expected {
		t.Errorf("Expected: \"%s\", Got \"%s\"", expected, actual)
	}
}

// ------------------- MARK: input output tests -------------------

func testForError(t *testing.T, php string, expected phpError.Error) {
	_, err := NewInterpreter(ini.NewDevIni(), &Request{}, "/home/admin/test.php").Process(php)
	if err.GetErrorType() != expected.GetErrorType() || err.GetMessage() != expected.GetMessage() {
		t.Errorf("\nCode: \"%s\"\nExpected: %s\nGot:      %s", php, expected, err)
	}
}

func testInputOutput(t *testing.T, php string, output string) *Interpreter {
	// Always use "\n" for tests so that they also pass on Windows
	PHP_EOL = "\n"
	interpreter := NewInterpreter(ini.NewDevIni(), &Request{}, "test.php")
	actual, err := interpreter.Process(php)
	if err != nil {
		t.Errorf("\nCode: \"%s\"\nUnexpected error: \"%s\"", php, err)
		return interpreter
	}
	if actual != output {
		t.Errorf("\nCode: \"%s\"\nExpected: \"%s\",\nGot \"%s\"", php, output, actual)
	}
	return interpreter
}

func TestOutput(t *testing.T) {
	// Without PHP
	testInputOutput(t, "<html>...</html>", "<html>...</html>")

	// Echo short tag
	testInputOutput(t, `<html><?= "abc\n" ?><?= 42; ?></html>`, "<html>abc\n42</html>")
	// Echo
	testInputOutput(t,
		`<html><?php echo "abc\t", 42 ?><?php echo "def", 24; ?></html>`, "<html>abc\t42def24</html>",
	)

	// Simple variable substitution
	testInputOutput(t, `<?php $a = 42; echo "a{$a}b";`, "a42b")
}

func TestConstants(t *testing.T) {
	// Predefined constants
	testInputOutput(t, `<?php echo E_USER_NOTICE;`, fmt.Sprintf("%d", phpError.E_USER_NOTICE))
	testInputOutput(t, `<?php echo E_ALL;`, fmt.Sprintf("%d", phpError.E_ALL))

	// Userdefined constants
	testInputOutput(t, `<?php const TRUTH = 42; const PI = "3.141";echo TRUTH, PI;`, "423.141")
}

func TestFileIncludes(t *testing.T) {
	testForError(t, `<?php require "include.php"; ?>`,
		phpError.NewError("Uncaught Error: Failed opening required 'include.php' (include_path='/home/admin') in /home/admin/test.php:1:15"),
	)
	testForError(t, `<?php require_once "include.php"; ?>`,
		phpError.NewError("Uncaught Error: Failed opening required 'include.php' (include_path='/home/admin') in /home/admin/test.php:1:20"),
	)
	testForError(t, `<?php include "include.php"; ?>`,
		phpError.NewWarning("include(): Failed opening 'include.php' for inclusion (include_path='/home/admin') in /home/admin/test.php:1:15"),
	)
	testForError(t, `<?php include_once "include.php"; ?>`,
		phpError.NewWarning("include(): Failed opening 'include.php' for inclusion (include_path='/home/admin') in /home/admin/test.php:1:20"),
	)
}

func TestVariable(t *testing.T) {
	// Undefined variable
	testInputOutput(t, `<?php echo is_null($a) ? "a" : "b";`, "Warning: Undefined variable $a\na")
	testInputOutput(t, `<?php echo intval($a);`, "Warning: Undefined variable $a\n0")
	testInputOutput(t, `<?php echo intval($$a);`, "Warning: Undefined variable $a\nWarning: Undefined variable $\n0")

	// Simple variable
	testInputOutput(t, `<?php $var = "hi"; $var = "hello"; echo $var, " world";`, "hello world")

	// Variable variable name
	testInputOutput(t, `<?php $var = "hi"; $$var = "hello"; echo $hi, " world";`, "hello world")

	// Chained variable declarations
	testInputOutput(t, `<?php $a = $b = $c = 42; echo $a, $b, $c;`, "424242")

	// Compound assignment
	testInputOutput(t, `<?php $a = 42; echo $a; $a += 2; echo $a; $a += $a; echo $a;`, "424488")

	// Parenthesized LHS
	testForError(t, `<?php ($a) = 42;`,
		phpError.NewParseError(`Statement must end with a semicolon. Got: "=" at /home/admin/test.php:1:12`),
	)
}

func TestConditionals(t *testing.T) {
	// Conditional
	testInputOutput(t, `<?php echo 1 ? "y" : "n"; echo 0 ? "n" : "y"; echo false ?: "y";`, "yyy")

	// Coalesce
	testInputOutput(t,
		`<?php $a = null; echo $a ?? "a"; $a = "b"; echo $a ?? "a"; echo "c" ?? "d";`, "abc",
	)
	testInputOutput(t, `<?php echo $a ?? "a";`, "a")

	// If statment
	testInputOutput(t, `<?php $a = 42; if ($a === 42) { echo "42"; } elseif ($a === 41) { echo "41"; } else { echo "??"; }`, "42")
	testInputOutput(t, `<?php $a = 41; if ($a === 42) { echo "42"; } elseif ($a === 41) { echo "41"; } else { echo "??"; }`, "41")
	testInputOutput(t, `<?php $a = 40; if ($a === 42) { echo "42"; } elseif ($a === 41) { echo "41"; } else { echo "??"; }`, "??")
	// Alternative syntax
	testInputOutput(t, `<?php $a = 42; if ($a === 42): echo "42"; elseif ($a === 41): echo "41"; else: echo "??"; endif;`, "42")
	testInputOutput(t, `<?php $a = 41; if ($a === 42): echo "42"; elseif ($a === 41): echo "41"; else: echo "??"; endif;`, "41")
	testInputOutput(t, `<?php $a = 40; if ($a === 42): echo "42"; elseif ($a === 41): echo "41"; else: echo "??"; endif;`, "??")

	// While statment
	testInputOutput(t, `<?php $a = 40; while ($a < 42) { echo "1"; $a++; }`, "11")
	testInputOutput(t, `<?php $a = 42; while ($a < 42) { echo "1"; $a++; }`, "")
	testInputOutput(t, `<?php $a = 0; while (true) { echo "1"; $a++; if ($a == 5) { break; }}`, "11111")
	testInputOutput(t, `<?php $a = 0; while (true) { echo "1"; while (true) { echo "2"; break 2; }}`, "12")
	// Alternative syntax
	testInputOutput(t, `<?php $a = 40; while ($a < 42): echo "1";  $a++; endwhile;`, "11")
	testInputOutput(t, `<?php $a = 42; while ($a < 42): echo "1";  $a++; endwhile;`, "")

	// Do statement
	testInputOutput(t, `<?php $a = 40; do { echo "1"; $a++; } while ($a < 42);`, "11")
	testInputOutput(t, `<?php $a = 42; do { echo "1"; $a++; } while ($a < 42);`, "1")
	testInputOutput(t, `<?php $a = 0; do { echo "1"; $a++; if ($a == 5) { break; }} while (true);`, "11111")
	testInputOutput(t, `<?php $a = 0; do { echo "1"; while (true) { echo "2"; break 2; }} while (true);`, "12")
}

func TestIntrinsic(t *testing.T) {
	// Exit
	interpreter := testInputOutput(t, `Hello <?php exit("world");`, "Hello world")
	if interpreter.GetExitCode() != 0 {
		t.Errorf("Expected: \"%d\", Got \"%d\"", 0, interpreter.GetExitCode())
	}
	interpreter = testInputOutput(t, `Hello<?php exit;`, "Hello")
	if interpreter.GetExitCode() != 0 {
		t.Errorf("Expected: \"%d\", Got \"%d\"", 0, interpreter.GetExitCode())
	}
	interpreter = testInputOutput(t, `Hello<?php exit(42);`, "Hello")
	if interpreter.GetExitCode() != 42 {
		t.Errorf("Expected: \"%d\", Got \"%d\"", 42, interpreter.GetExitCode())
	}

	// Empty
	testInputOutput(t, `<?php echo empty(false) ? "y" : "n";`, "y")
	testInputOutput(t, `<?php echo empty(true) ? "y" : "n";`, "n")
	testInputOutput(t, `<?php echo empty(0) ? "y" : "n";`, "y")
	testInputOutput(t, `<?php echo empty(1) ? "y" : "n";`, "n")
	testInputOutput(t, `<?php echo empty(0.0) ? "y" : "n";`, "y")
	testInputOutput(t, `<?php echo empty(2.0) ? "y" : "n";`, "n")
	testInputOutput(t, `<?php echo empty("") ? "y" : "n";`, "y")
	testInputOutput(t, `<?php echo empty("0") ? "y" : "n";`, "y")
	testInputOutput(t, `<?php echo empty("1") ? "y" : "n";`, "n")
	testInputOutput(t, `<?php echo empty("00") ? "y" : "n";`, "n")
	testInputOutput(t, `<?php echo empty(null) ? "y" : "n";`, "y")
	testInputOutput(t, `<?php echo empty($a) ? "y" : "n";`, "y")
	testInputOutput(t, `<?php $a = 1; echo empty($a) ? "y" : "n";`, "n")

	// Isset
	testInputOutput(t, `<?php $a = 1; echo isset($a) ? "y" : "n";`, "y")
	testInputOutput(t, `<?php echo isset($a) ? "y" : "n";`, "n")
	testInputOutput(t, `<?php $a = 1; echo isset($a, $b) ? "y" : "n";`, "n")
	testInputOutput(t, `<?php echo isset($a, $b) ? "y" : "n";`, "n")

	// Unset
	testInputOutput(t, `<?php $a = 1; echo isset($a) ? "y" : "n"; unset($a); echo isset($a) ? "y" : "n";`, "yn")
	testInputOutput(t, `<?php echo isset($a) ? "y" : "n"; unset($a); echo isset($a) ? "y" : "n";`, "nn")
}

func TestUserFunctions(t *testing.T) {
	// Simple user defined function without types, params, ...
	// Check if function definition can be after the function call
	testInputOutput(t,
		`<?php
			echo "a";
			helloWorld();
			echo "b";
			function helloWorld() { echo "Hello World"; }
			echo "c";
		`,
		"aHello Worldbc",
	)
	testInputOutput(t,
		`<?php
			echo "a";
			helloWorld();
			echo "b";
			{{{ function helloWorld() { echo "Hello World"; } }}}
			echo "c";
		`,
		"aHello Worldbc",
	)

	// Test correct scoping: $a is available in the function body
	testInputOutput(t,
		`<?php
			$a = 1;
			func();
			function func() {
				echo isset($a) ? "y" : "n";
				$b = 2;
			}
			echo isset($b) ? "y" : "n";
		`,
		"nn",
	)

	// Parameters
	testInputOutput(t,
		`<?php
			$a = 1;
			func($a);
			func(2);
			function func($param) {
				echo $param;
			}
		`,
		"12",
	)
	testInputOutput(t,
		`<?php
			$a = 1;
			func($a, 2);
			func(1, 1+1);
			function func($param1, $param2) {
				echo $param1 + $param2;
			}
		`,
		"33",
	)
	// Typed parameters
	testInputOutput(t,
		`<?php
			$a = 1;
			func($a, 2);
			func(1, 1+1);
			function func(int|float $param1, int|float $param2) {
				echo $param1 + $param2;
			}
		`,
		"33",
	)

	// Return type
	testInputOutput(t,
		`<?php
			$a = 1;
			echo func($a, 2);
			echo func(1, 1+1);
			function func(int|float $param1, int|float $param2,): int|float {
				return $param1 + $param2;
			}
		`,
		"33",
	)
}

func TestArray(t *testing.T) {
	testInputOutput(t, `<?php $a = [0, 1, 2]; echo $a[0] === null ? "y" : "n";`, "n")
	testInputOutput(t, `<?php $a = [0, 1, 2]; echo $a[3] === null ? "y" : "n";`, "y")
	testInputOutput(t, `<?php $a = [0, 1]; echo $a[2] = 2; echo $a[2];`, "22")
	// TODO add test with nested: $b["a"]["b"]["c"]=1;

	// Pass by value not reference
	testInputOutput(t,
		`<?php $a = $b = [42]; var_dump($a[0], $b[0]); $b[0] = 43; var_dump($a[0], $b[0]);`,
		"int(42)\nint(42)\nint(42)\nint(43)\n",
	)
	testInputOutput(t,
		`<?php $b = [42]; $a = $b; var_dump($a[0], $b[0]); $b[0] = 43; var_dump($a[0], $b[0]);`,
		"int(42)\nint(42)\nint(42)\nint(43)\n",
	)
}

func TestCastExpression(t *testing.T) {
	testInputOutput(t, `<?php var_dump((array)42);`, "array(1) {\n  [0]=>\n  int(42)\n}\n")
	testInputOutput(t, `<?php var_dump((binary)42);`, `string(2) "42"`+"\n")
	testInputOutput(t, `<?php var_dump((bool)42);`, "bool(true)\n")
	testInputOutput(t, `<?php var_dump((boolean)42);`, "bool(true)\n")
	testInputOutput(t, `<?php var_dump((double)42);`, "float(42)\n")
	testInputOutput(t, `<?php var_dump((int)42);`, "int(42)\n")
	testInputOutput(t, `<?php var_dump((integer)42);`, "int(42)\n")
	testInputOutput(t, `<?php var_dump((float)42);`, "float(42)\n")
	// TODO testInputOutput(t, `<?php var_dump((object)42);`, "a")
	testInputOutput(t, `<?php var_dump((real)42);`, "float(42)\n")
	testInputOutput(t, `<?php var_dump((string)42);`, `string(2) "42"`+"\n")
}

func TestOperators(t *testing.T) {
	// Logical "not"
	testInputOutput(t, `<?php echo !true ? "y" : "y";`, "y")
	testInputOutput(t, `<?php echo !false ? "y" : "y";`, "y")
	testInputOutput(t, `<?php echo !42 ? "y" : "y";`, "y")

	// Logical "and" and "or"
	testInputOutput(t, `<?php echo 4 && 0 ? "t" : "f";`, "f")
	testInputOutput(t, `<?php echo 4 && 1 ? "t" : "f";`, "t")
	testInputOutput(t, `<?php echo 4 && false ? "t" : "f";`, "f")
	testInputOutput(t, `<?php echo 4 && true ? "t" : "f";`, "t")
	testInputOutput(t, `<?php echo 0 || 0 ? "t" : "f";`, "f")
	testInputOutput(t, `<?php echo 4 || 1 ? "t" : "f";`, "t")
	testInputOutput(t, `<?php echo false || false ? "t" : "f";`, "f")
	testInputOutput(t, `<?php echo 4 || true ? "t" : "f";`, "t")

	// Unary expression
	// Boolean
	testInputOutput(t, `<?php var_dump(+true);`, "int(1)\n")
	testInputOutput(t, `<?php var_dump(-true);`, "int(-1)\n")
	testInputOutput(t, `<?php var_dump(+false);`, "int(0)\n")
	testInputOutput(t, `<?php var_dump(-false);`, "int(0)\n")
	// Floating
	testInputOutput(t, `<?php var_dump(+2.5);`, "float(2.5)\n")
	testInputOutput(t, `<?php var_dump(+(-2.5));`, "float(-2.5)\n")
	testInputOutput(t, `<?php var_dump(-3.0);`, "float(-3)\n")
	testInputOutput(t, `<?php var_dump(-(-3.0));`, "float(3)\n")
	testInputOutput(t, `<?php var_dump(~(-3.0));`, "int(2)\n")
	testInputOutput(t, `<?php var_dump(~3.0);`, "int(-4)\n")
	// Integer
	testInputOutput(t, `<?php var_dump(+2);`, "int(2)\n")
	testInputOutput(t, `<?php var_dump(+(-2));`, "int(-2)\n")
	testInputOutput(t, `<?php var_dump(-3);`, "int(-3)\n")
	testInputOutput(t, `<?php var_dump(-(-3));`, "int(3)\n")
	testInputOutput(t, `<?php var_dump(~(-3));`, "int(2)\n")
	testInputOutput(t, `<?php var_dump(~3);`, "int(-4)\n")

	// Post-/Prefix Inc-/Decrement
	// Integer
	testInputOutput(t, `<?php $a = 42; var_dump($a++); var_dump($a);`, "int(42)\nint(43)\n")
	testInputOutput(t, `<?php $a = 42; var_dump($a--); var_dump($a);`, "int(42)\nint(41)\n")
	testInputOutput(t, `<?php $a = 42; var_dump(++$a); var_dump($a);`, "int(43)\nint(43)\n")
	testInputOutput(t, `<?php $a = 42; var_dump(--$a); var_dump($a);`, "int(41)\nint(41)\n")
	// Boolean
	testInputOutput(t, `<?php $a = true; var_dump($a++); var_dump($a);`, "bool(true)\nbool(true)\n")
	testInputOutput(t, `<?php $a = true; var_dump($a--); var_dump($a);`, "bool(true)\nbool(true)\n")
	testInputOutput(t, `<?php $a = false; var_dump($a++); var_dump($a);`, "bool(false)\nbool(false)\n")
	testInputOutput(t, `<?php $a = false; var_dump($a--); var_dump($a);`, "bool(false)\nbool(false)\n")
	// Floating
	testInputOutput(t, `<?php $a = 42.0; var_dump($a++); var_dump($a);`, "float(42)\nfloat(43)\n")
	testInputOutput(t, `<?php $a = 42.0; var_dump($a--); var_dump($a);`, "float(42)\nfloat(41)\n")
	// Null
	testInputOutput(t, `<?php $a = null; var_dump($a++); var_dump($a);`, "NULL\nint(1)\n")
	testInputOutput(t, `<?php $a = null; var_dump($a--); var_dump($a);`, "NULL\nNULL\n")
	// String
	testInputOutput(t, `<?php $a = ""; var_dump($a++); var_dump($a);`, `string(0) ""`+"\n"+`string(1) "1"`+"\n")
	testInputOutput(t, `<?php $a = ""; var_dump($a--); var_dump($a);`, `string(0) ""`+"\nint(-1)\n")

	// Base operators
	// Integer
	testInputOutput(t, `<?php echo 4 >> 2;`, "1")
	testInputOutput(t, `<?php echo 8 << 2;`, "32")
	testInputOutput(t, `<?php $a = 13; echo $a <<= 1;`, "26")
	testInputOutput(t, `<?php echo 4 ^ 4;`, "0")
	testInputOutput(t, `<?php echo 8 ^ 4;`, "12")
	testInputOutput(t, `<?php $a = 13; echo $a ^= 4;`, "9")
	testInputOutput(t, `<?php echo 4 | 4;`, "4")
	testInputOutput(t, `<?php echo 8 | 4;`, "12")
	testInputOutput(t, `<?php $a = 8; echo $a |= 4;`, "12")
	testInputOutput(t, `<?php echo 8 & 4;`, "0")
	testInputOutput(t, `<?php echo 12 & 8;`, "8")
	testInputOutput(t, `<?php $a = 12; echo $a &= 4;`, "4")
	testInputOutput(t, `<?php echo 42 + 1;`, "43")
	testInputOutput(t, `<?php $a = 42; echo $a += 1;`, "43")
	testInputOutput(t, `<?php echo 42 - 1;`, "41")
	testInputOutput(t, `<?php $a = 42; echo $a -= 1;`, "41")
	testInputOutput(t, `<?php echo 42 * 2;`, "84")
	testInputOutput(t, `<?php $a = 42; echo $a *= 2;`, "84")
	testInputOutput(t, `<?php echo 42 / 2;`, "21")
	testInputOutput(t, `<?php $a = 42; echo $a /= 2;`, "21")
	testInputOutput(t, `<?php echo 42 % 5;`, "2")
	testInputOutput(t, `<?php $a = 42; echo $a %= 5;`, "2")
	testInputOutput(t, `<?php echo 2 ** 4;`, "16")
	testInputOutput(t, `<?php echo 2 ** 2 ** 2;`, "16")
	testInputOutput(t, `<?php $a = 2; echo $a **= 4;`, "16")
	// Floating
	testInputOutput(t, `<?php echo 42.0 + 1.5;`, "43.5")
	testInputOutput(t, `<?php $a = 42.0; echo $a += 1.5;`, "43.5")
	testInputOutput(t, `<?php echo 42 - 1.5;`, "40.5")
	testInputOutput(t, `<?php $a = 42.0; echo $a -= 1.5;`, "40.5")
	testInputOutput(t, `<?php echo 42.1 * 2;`, "84.2")
	testInputOutput(t, `<?php $a = 42.1; echo $a *= 2;`, "84.2")
	testInputOutput(t, `<?php echo 43.0 / 2;`, "21.5")
	testInputOutput(t, `<?php $a = 43.0; echo $a /= 2;`, "21.5")
	testInputOutput(t, `<?php echo 2.0 ** 4;`, "16")
	testInputOutput(t, `<?php $a = 2.0; echo $a **= 4;`, "16")
	// String
	testInputOutput(t, `<?php echo "a" . "bc";`, "abc")
	testInputOutput(t, `<?php $a = "a"; echo $a .= "bc";`, "abc")
	// Combined additions and multiplications
	testInputOutput(t, `<?php echo 31 + 21 + 11;`, "63")
	testInputOutput(t, `<?php echo 4 * 3 * 2;`, "24")
	testInputOutput(t, `<?php echo 2 + 3 * 4;`, "14")
	testInputOutput(t, `<?php echo (2 + 3) * 4;`, "20")
	testInputOutput(t, `<?php echo 2 * 3 + 4 * 5 + 6;`, "32")
	testInputOutput(t, `<?php echo 2 * (3 + 4) * 5 + 6;`, "76")
	testInputOutput(t, `<?php echo 2 + 3 * 4 + 5 * 6;`, "44")
}

func TestStrictComparison(t *testing.T) {
	// Table from https://www.php.net/manual/en/types.comparisons.php
	// true
	testInputOutput(t, `<?php var_dump(true === true);`, "bool(true)\n")
	testInputOutput(t, `<?php var_dump(true === false);`, "bool(false)\n")
	testInputOutput(t, `<?php var_dump(true === 1);`, "bool(false)\n")
	testInputOutput(t, `<?php var_dump(true === 0);`, "bool(false)\n")
	testInputOutput(t, `<?php var_dump(true === -1);`, "bool(false)\n")
	testInputOutput(t, `<?php var_dump(true === "1");`, "bool(false)\n")
	testInputOutput(t, `<?php var_dump(true === "0");`, "bool(false)\n")
	testInputOutput(t, `<?php var_dump(true === "-1");`, "bool(false)\n")
	testInputOutput(t, `<?php var_dump(true === null);`, "bool(false)\n")
	testInputOutput(t, `<?php var_dump(true === []);`, "bool(false)\n")
	testInputOutput(t, `<?php var_dump(true === "php");`, "bool(false)\n")
	testInputOutput(t, `<?php var_dump(true === "");`, "bool(false)\n")
	// false
	testInputOutput(t, `<?php var_dump(false === false);`, "bool(true)\n")
	testInputOutput(t, `<?php var_dump(false === 1);`, "bool(false)\n")
	testInputOutput(t, `<?php var_dump(false === 0);`, "bool(false)\n")
	testInputOutput(t, `<?php var_dump(false === -1);`, "bool(false)\n")
	testInputOutput(t, `<?php var_dump(false === "1");`, "bool(false)\n")
	testInputOutput(t, `<?php var_dump(false === "0");`, "bool(false)\n")
	testInputOutput(t, `<?php var_dump(false === "-1");`, "bool(false)\n")
	testInputOutput(t, `<?php var_dump(false === null);`, "bool(false)\n")
	testInputOutput(t, `<?php var_dump(false === []);`, "bool(false)\n")
	testInputOutput(t, `<?php var_dump(false === "php");`, "bool(false)\n")
	testInputOutput(t, `<?php var_dump(false === "");`, "bool(false)\n")
	// 1
	testInputOutput(t, `<?php var_dump(1 === 1);`, "bool(true)\n")
	testInputOutput(t, `<?php var_dump(1 === 0);`, "bool(false)\n")
	testInputOutput(t, `<?php var_dump(1 === -1);`, "bool(false)\n")
	testInputOutput(t, `<?php var_dump(1 === "1");`, "bool(false)\n")
	testInputOutput(t, `<?php var_dump(1 === "0");`, "bool(false)\n")
	testInputOutput(t, `<?php var_dump(1 === "-1");`, "bool(false)\n")
	testInputOutput(t, `<?php var_dump(1 === null);`, "bool(false)\n")
	testInputOutput(t, `<?php var_dump(1 === []);`, "bool(false)\n")
	testInputOutput(t, `<?php var_dump(1 === "php");`, "bool(false)\n")
	testInputOutput(t, `<?php var_dump(1 === "");`, "bool(false)\n")
	// 0
	testInputOutput(t, `<?php var_dump(0 === 0);`, "bool(true)\n")
	testInputOutput(t, `<?php var_dump(0 === -1);`, "bool(false)\n")
	testInputOutput(t, `<?php var_dump(0 === "1");`, "bool(false)\n")
	testInputOutput(t, `<?php var_dump(0 === "0");`, "bool(false)\n")
	testInputOutput(t, `<?php var_dump(0 === "-1");`, "bool(false)\n")
	testInputOutput(t, `<?php var_dump(0 === null);`, "bool(false)\n")
	testInputOutput(t, `<?php var_dump(0 === []);`, "bool(false)\n")
	testInputOutput(t, `<?php var_dump(0 === "php");`, "bool(false)\n")
	testInputOutput(t, `<?php var_dump(0 === "");`, "bool(false)\n")
	// -1
	testInputOutput(t, `<?php var_dump(-1 === -1);`, "bool(true)\n")
	testInputOutput(t, `<?php var_dump(-1 === "1");`, "bool(false)\n")
	testInputOutput(t, `<?php var_dump(-1 === "0");`, "bool(false)\n")
	testInputOutput(t, `<?php var_dump(-1 === "-1");`, "bool(false)\n")
	testInputOutput(t, `<?php var_dump(-1 === null);`, "bool(false)\n")
	testInputOutput(t, `<?php var_dump(-1 === []);`, "bool(false)\n")
	testInputOutput(t, `<?php var_dump(-1 === "php");`, "bool(false)\n")
	testInputOutput(t, `<?php var_dump(-1 === "");`, "bool(false)\n")
	// "1"
	testInputOutput(t, `<?php var_dump("1" === "1");`, "bool(true)\n")
	testInputOutput(t, `<?php var_dump("1" === "0");`, "bool(false)\n")
	testInputOutput(t, `<?php var_dump("1" === "-1");`, "bool(false)\n")
	testInputOutput(t, `<?php var_dump("1" === null);`, "bool(false)\n")
	testInputOutput(t, `<?php var_dump("1" === []);`, "bool(false)\n")
	testInputOutput(t, `<?php var_dump("1" === "php");`, "bool(false)\n")
	testInputOutput(t, `<?php var_dump("1" === "");`, "bool(false)\n")
	// "0"
	testInputOutput(t, `<?php var_dump("0" === "0");`, "bool(true)\n")
	testInputOutput(t, `<?php var_dump("0" === "-1");`, "bool(false)\n")
	testInputOutput(t, `<?php var_dump("0" === null);`, "bool(false)\n")
	testInputOutput(t, `<?php var_dump("0" === []);`, "bool(false)\n")
	testInputOutput(t, `<?php var_dump("0" === "php");`, "bool(false)\n")
	testInputOutput(t, `<?php var_dump("0" === "");`, "bool(false)\n")
	// "-1"
	testInputOutput(t, `<?php var_dump("-1" === "-1");`, "bool(true)\n")
	testInputOutput(t, `<?php var_dump("-1" === null);`, "bool(false)\n")
	testInputOutput(t, `<?php var_dump("-1" === []);`, "bool(false)\n")
	testInputOutput(t, `<?php var_dump("-1" === "php");`, "bool(false)\n")
	testInputOutput(t, `<?php var_dump("-1" === "");`, "bool(false)\n")
	// null
	testInputOutput(t, `<?php var_dump(null === null);`, "bool(true)\n")
	testInputOutput(t, `<?php var_dump(null === []);`, "bool(false)\n")
	testInputOutput(t, `<?php var_dump(null === "php");`, "bool(false)\n")
	testInputOutput(t, `<?php var_dump(null === "");`, "bool(false)\n")
	// []
	testInputOutput(t, `<?php var_dump([] === []);`, "bool(true)\n")
	testInputOutput(t, `<?php var_dump([] === "php");`, "bool(false)\n")
	testInputOutput(t, `<?php var_dump([] === "");`, "bool(false)\n")
	// "php"
	testInputOutput(t, `<?php var_dump("php" === "php");`, "bool(true)\n")
	testInputOutput(t, `<?php var_dump("php" === "");`, "bool(false)\n")
	// ""
	testInputOutput(t, `<?php var_dump("" === "");`, "bool(true)\n")
}

func TestLooseComparison(t *testing.T) {
	// Table from https://www.php.net/manual/en/types.comparisons.php
	// true
	testInputOutput(t, `<?php var_dump(true == true);`, "bool(true)\n")
	testInputOutput(t, `<?php var_dump(true == false);`, "bool(false)\n")
	testInputOutput(t, `<?php var_dump(true == 1);`, "bool(true)\n")
	testInputOutput(t, `<?php var_dump(true == 0);`, "bool(false)\n")
	testInputOutput(t, `<?php var_dump(true == -1);`, "bool(true)\n")
	testInputOutput(t, `<?php var_dump(true == "1");`, "bool(true)\n")
	testInputOutput(t, `<?php var_dump(true == "0");`, "bool(false)\n")
	testInputOutput(t, `<?php var_dump(true == "-1");`, "bool(true)\n")
	testInputOutput(t, `<?php var_dump(true == null);`, "bool(false)\n")
	testInputOutput(t, `<?php var_dump(true == []);`, "bool(false)\n")
	testInputOutput(t, `<?php var_dump(true == "php");`, "bool(true)\n")
	testInputOutput(t, `<?php var_dump(true == "");`, "bool(false)\n")
	// false
	testInputOutput(t, `<?php var_dump(false == false);`, "bool(true)\n")
	testInputOutput(t, `<?php var_dump(false == 1);`, "bool(false)\n")
	testInputOutput(t, `<?php var_dump(false == 0);`, "bool(true)\n")
	testInputOutput(t, `<?php var_dump(false == -1);`, "bool(false)\n")
	testInputOutput(t, `<?php var_dump(false == "1");`, "bool(false)\n")
	testInputOutput(t, `<?php var_dump(false == "0");`, "bool(true)\n")
	testInputOutput(t, `<?php var_dump(false == "-1");`, "bool(false)\n")
	testInputOutput(t, `<?php var_dump(false == null);`, "bool(true)\n")
	testInputOutput(t, `<?php var_dump(false == []);`, "bool(true)\n")
	testInputOutput(t, `<?php var_dump(false == "php");`, "bool(false)\n")
	testInputOutput(t, `<?php var_dump(false == "");`, "bool(true)\n")
	// 1
	testInputOutput(t, `<?php var_dump(1 == 1);`, "bool(true)\n")
	testInputOutput(t, `<?php var_dump(1 == 0);`, "bool(false)\n")
	testInputOutput(t, `<?php var_dump(1 == -1);`, "bool(false)\n")
	testInputOutput(t, `<?php var_dump(1 == "1");`, "bool(true)\n")
	testInputOutput(t, `<?php var_dump(1 == "0");`, "bool(false)\n")
	testInputOutput(t, `<?php var_dump(1 == "-1");`, "bool(false)\n")
	testInputOutput(t, `<?php var_dump(1 == null);`, "bool(false)\n")
	testInputOutput(t, `<?php var_dump(1 == []);`, "bool(false)\n")
	testInputOutput(t, `<?php var_dump(1 == "php");`, "bool(false)\n")
	testInputOutput(t, `<?php var_dump(1 == "");`, "bool(false)\n")
	// 0
	testInputOutput(t, `<?php var_dump(0 == 0);`, "bool(true)\n")
	testInputOutput(t, `<?php var_dump(0 == -1);`, "bool(false)\n")
	testInputOutput(t, `<?php var_dump(0 == "1");`, "bool(false)\n")
	testInputOutput(t, `<?php var_dump(0 == "0");`, "bool(true)\n")
	testInputOutput(t, `<?php var_dump(0 == "-1");`, "bool(false)\n")
	testInputOutput(t, `<?php var_dump(0 == null);`, "bool(true)\n")
	testInputOutput(t, `<?php var_dump(0 == []);`, "bool(false)\n")
	testInputOutput(t, `<?php var_dump(0 == "php");`, "bool(false)\n")
	testInputOutput(t, `<?php var_dump(0 == "");`, "bool(false)\n")
	// -1
	testInputOutput(t, `<?php var_dump(-1 == -1);`, "bool(true)\n")
	testInputOutput(t, `<?php var_dump(-1 == "1");`, "bool(false)\n")
	testInputOutput(t, `<?php var_dump(-1 == "0");`, "bool(false)\n")
	testInputOutput(t, `<?php var_dump(-1 == "-1");`, "bool(true)\n")
	testInputOutput(t, `<?php var_dump(-1 == null);`, "bool(false)\n")
	testInputOutput(t, `<?php var_dump(-1 == []);`, "bool(false)\n")
	testInputOutput(t, `<?php var_dump(-1 == "php");`, "bool(false)\n")
	testInputOutput(t, `<?php var_dump(-1 == "");`, "bool(false)\n")
	// "1"
	testInputOutput(t, `<?php var_dump("1" == "1");`, "bool(true)\n")
	testInputOutput(t, `<?php var_dump("1" == "0");`, "bool(false)\n")
	testInputOutput(t, `<?php var_dump("1" == "-1");`, "bool(false)\n")
	testInputOutput(t, `<?php var_dump("1" == null);`, "bool(false)\n")
	testInputOutput(t, `<?php var_dump("1" == []);`, "bool(false)\n")
	testInputOutput(t, `<?php var_dump("1" == "php");`, "bool(false)\n")
	testInputOutput(t, `<?php var_dump("1" == "");`, "bool(false)\n")
	// "0"
	testInputOutput(t, `<?php var_dump("0" == "0");`, "bool(true)\n")
	testInputOutput(t, `<?php var_dump("0" == "-1");`, "bool(false)\n")
	testInputOutput(t, `<?php var_dump("0" == null);`, "bool(false)\n")
	testInputOutput(t, `<?php var_dump("0" == []);`, "bool(false)\n")
	testInputOutput(t, `<?php var_dump("0" == "php");`, "bool(false)\n")
	testInputOutput(t, `<?php var_dump("0" == "");`, "bool(false)\n")
	// "-1"
	testInputOutput(t, `<?php var_dump("-1" == "-1");`, "bool(true)\n")
	testInputOutput(t, `<?php var_dump("-1" == null);`, "bool(false)\n")
	testInputOutput(t, `<?php var_dump("-1" == []);`, "bool(false)\n")
	testInputOutput(t, `<?php var_dump("-1" == "php");`, "bool(false)\n")
	testInputOutput(t, `<?php var_dump("-1" == "");`, "bool(false)\n")
	// null
	testInputOutput(t, `<?php var_dump(null == null);`, "bool(true)\n")
	testInputOutput(t, `<?php var_dump(null == []);`, "bool(true)\n")
	testInputOutput(t, `<?php var_dump(null == "php");`, "bool(false)\n")
	testInputOutput(t, `<?php var_dump(null == "");`, "bool(true)\n")
	// []
	testInputOutput(t, `<?php var_dump([] == []);`, "bool(true)\n")
	testInputOutput(t, `<?php var_dump([] == "php");`, "bool(false)\n")
	testInputOutput(t, `<?php var_dump([] == "");`, "bool(false)\n")
	// "php"
	testInputOutput(t, `<?php var_dump("php" == "php");`, "bool(true)\n")
	testInputOutput(t, `<?php var_dump("php" == "");`, "bool(false)\n")
	// ""
	testInputOutput(t, `<?php var_dump("" == "");`, "bool(true)\n")

	// Warning-Box from https://www.php.net/manual/en/language.operators.comparison.php
	testInputOutput(t, `<?php var_dump(0 == "a");`, "bool(false)\n")
	testInputOutput(t, `<?php var_dump("a" == 0);`, "bool(false)\n")
	testInputOutput(t, `<?php var_dump("1" == "01a");`, "bool(false)\n")
	testInputOutput(t, `<?php var_dump("01a" == "1");`, "bool(false)\n")
	testInputOutput(t, `<?php var_dump(1 == "1a");`, "bool(false)\n")
	testInputOutput(t, `<?php var_dump("1a" == 1);`, "bool(false)\n")
	testInputOutput(t, `<?php var_dump("1" == "01");`, "bool(true)\n")
	testInputOutput(t, `<?php var_dump("01" == "1");`, "bool(true)\n")
	testInputOutput(t, `<?php var_dump("10" == "1e1");`, "bool(true)\n")
	testInputOutput(t, `<?php var_dump("1e1" == "10");`, "bool(true)\n")
	testInputOutput(t, `<?php var_dump(100 == "1e2");`, "bool(true)\n")
	testInputOutput(t, `<?php var_dump(1e2 == "100");`, "bool(true)\n")
}

func TestCompareRelation(t *testing.T) {
	// Array
	// Array - Array
	testInputOutput(t, `<?php var_dump(["a"] < []);`, "bool(false)\n")
	testInputOutput(t, `<?php var_dump(["a"] <= []);`, "bool(false)\n")
	testInputOutput(t, `<?php var_dump(["a"] <=> []);`, "int(1)\n")
	testInputOutput(t, `<?php var_dump([] < []);`, "bool(false)\n")
	testInputOutput(t, `<?php var_dump([] <= []);`, "bool(true)\n")
	testInputOutput(t, `<?php var_dump([] <=> []);`, "int(0)\n")
	testInputOutput(t, `<?php var_dump(["a"] <=> ["a"]);`, "int(0)\n")
	testInputOutput(t, `<?php var_dump([] < ["a"]);`, "bool(true)\n")
	testInputOutput(t, `<?php var_dump([] <= ["a"]);`, "bool(true)\n")
	testInputOutput(t, `<?php var_dump([] <=> ["a"]);`, "int(-1)\n")
	testInputOutput(t, `<?php var_dump([21] <=> [42]);`, "int(-1)\n")
	testInputOutput(t, `<?php var_dump([21] <=> [21, 42]);`, "int(-1)\n")
	// Array - Boolean
	testInputOutput(t, `<?php var_dump([] < true);`, "bool(true)\n")
	testInputOutput(t, `<?php var_dump([] <= true);`, "bool(true)\n")
	testInputOutput(t, `<?php var_dump([] <=> true);`, "int(-1)\n")
	testInputOutput(t, `<?php var_dump([] < false);`, "bool(false)\n")
	testInputOutput(t, `<?php var_dump([] <= false);`, "bool(true)\n")
	testInputOutput(t, `<?php var_dump([] <=> false);`, "int(0)\n")
	testInputOutput(t, `<?php var_dump([42] < true);`, "bool(false)\n")
	testInputOutput(t, `<?php var_dump([42] <= true);`, "bool(true)\n")
	testInputOutput(t, `<?php var_dump([42] <=> true);`, "int(0)\n")
	testInputOutput(t, `<?php var_dump([42] < false);`, "bool(false)\n")
	testInputOutput(t, `<?php var_dump([42] <= false);`, "bool(false)\n")
	testInputOutput(t, `<?php var_dump([42] <=> false);`, "int(1)\n")
	// Array - Floating
	testInputOutput(t, `<?php var_dump([] < 2.5);`, "bool(false)\n")
	testInputOutput(t, `<?php var_dump([] <= 2.5);`, "bool(false)\n")
	testInputOutput(t, `<?php var_dump([] <=> 2.5);`, "int(1)\n")
	// Array - Integer
	testInputOutput(t, `<?php var_dump([] < 2);`, "bool(false)\n")
	testInputOutput(t, `<?php var_dump([] <= 2);`, "bool(false)\n")
	testInputOutput(t, `<?php var_dump([] <=> 2);`, "int(1)\n")
	// Array - Null
	testInputOutput(t, `<?php var_dump([] < NULL);`, "bool(false)\n")
	testInputOutput(t, `<?php var_dump([] <= NULL);`, "bool(true)\n")
	testInputOutput(t, `<?php var_dump([] <=> NULL);`, "int(0)\n")
	testInputOutput(t, `<?php var_dump(["abc"] < NULL);`, "bool(false)\n")
	testInputOutput(t, `<?php var_dump(["abc"] <= NULL);`, "bool(false)\n")
	testInputOutput(t, `<?php var_dump(["abc"] <=> NULL);`, "int(1)\n")
	// Array - String
	testInputOutput(t, `<?php var_dump([] < "");`, "bool(false)\n")
	testInputOutput(t, `<?php var_dump([] <= "");`, "bool(false)\n")
	testInputOutput(t, `<?php var_dump([] <=> "");`, "int(1)\n")

	// Boolean
	// Boolean - Array
	testInputOutput(t, `<?php var_dump(true < []);`, "bool(false)\n")
	testInputOutput(t, `<?php var_dump(true <= []);`, "bool(false)\n")
	testInputOutput(t, `<?php var_dump(true <=> []);`, "int(1)\n")
	testInputOutput(t, `<?php var_dump(false < []);`, "bool(false)\n")
	testInputOutput(t, `<?php var_dump(false <= []);`, "bool(true)\n")
	testInputOutput(t, `<?php var_dump(false <=> []);`, "int(0)\n")
	testInputOutput(t, `<?php var_dump(true < [42]);`, "bool(false)\n")
	testInputOutput(t, `<?php var_dump(true <= [42]);`, "bool(true)\n")
	testInputOutput(t, `<?php var_dump(true <=> [42]);`, "int(0)\n")
	testInputOutput(t, `<?php var_dump(false < [42]);`, "bool(true)\n")
	testInputOutput(t, `<?php var_dump(false <= [42]);`, "bool(true)\n")
	testInputOutput(t, `<?php var_dump(false <=> [42]);`, "int(-1)\n")
	// Boolean - Boolean
	testInputOutput(t, `<?php var_dump(true < true);`, "bool(false)\n")
	testInputOutput(t, `<?php var_dump(true <= true);`, "bool(true)\n")
	testInputOutput(t, `<?php var_dump(true <=> true);`, "int(0)\n")
	testInputOutput(t, `<?php var_dump(true < false);`, "bool(false)\n")
	testInputOutput(t, `<?php var_dump(true <= false);`, "bool(false)\n")
	testInputOutput(t, `<?php var_dump(true <=> false);`, "int(1)\n")
	testInputOutput(t, `<?php var_dump(false < true);`, "bool(true)\n")
	testInputOutput(t, `<?php var_dump(false <= true);`, "bool(true)\n")
	testInputOutput(t, `<?php var_dump(false <=> true);`, "int(-1)\n")
	testInputOutput(t, `<?php var_dump(false < false);`, "bool(false)\n")
	testInputOutput(t, `<?php var_dump(false <= false);`, "bool(true)\n")
	testInputOutput(t, `<?php var_dump(false <=> false);`, "int(0)\n")
	// Boolean - Floating
	testInputOutput(t, `<?php var_dump(true < 2.5);`, "bool(false)\n")
	testInputOutput(t, `<?php var_dump(true <= 2.5);`, "bool(true)\n")
	testInputOutput(t, `<?php var_dump(true <=> 2.5);`, "int(0)\n")
	testInputOutput(t, `<?php var_dump(true < 0.5);`, "bool(false)\n")
	testInputOutput(t, `<?php var_dump(true <= 0.5);`, "bool(true)\n")
	testInputOutput(t, `<?php var_dump(true <=> 0.5);`, "int(0)\n")
	testInputOutput(t, `<?php var_dump(false < 2.5);`, "bool(true)\n")
	testInputOutput(t, `<?php var_dump(false <= 2.5);`, "bool(true)\n")
	testInputOutput(t, `<?php var_dump(false <=> 2.5);`, "int(-1)\n")
	testInputOutput(t, `<?php var_dump(false < 0.0);`, "bool(false)\n")
	testInputOutput(t, `<?php var_dump(false <= 0.0);`, "bool(true)\n")
	testInputOutput(t, `<?php var_dump(false <=> 0.0);`, "int(0)\n")
	// Boolean - Integer
	testInputOutput(t, `<?php var_dump(true < 2);`, "bool(false)\n")
	testInputOutput(t, `<?php var_dump(true <= 2);`, "bool(true)\n")
	testInputOutput(t, `<?php var_dump(true <=> 2);`, "int(0)\n")
	testInputOutput(t, `<?php var_dump(true < 0);`, "bool(false)\n")
	testInputOutput(t, `<?php var_dump(true <= 0);`, "bool(false)\n")
	testInputOutput(t, `<?php var_dump(true <=> 0);`, "int(1)\n")
	testInputOutput(t, `<?php var_dump(false < 2);`, "bool(true)\n")
	testInputOutput(t, `<?php var_dump(false <= 2);`, "bool(true)\n")
	testInputOutput(t, `<?php var_dump(false <=> 2);`, "int(-1)\n")
	testInputOutput(t, `<?php var_dump(false < 0);`, "bool(false)\n")
	testInputOutput(t, `<?php var_dump(false <= 0);`, "bool(true)\n")
	testInputOutput(t, `<?php var_dump(false <=> 0);`, "int(0)\n")
	// Boolean - Null
	testInputOutput(t, `<?php var_dump(true < NULL);`, "bool(false)\n")
	testInputOutput(t, `<?php var_dump(true <= NULL);`, "bool(false)\n")
	testInputOutput(t, `<?php var_dump(true <=> NULL);`, "int(1)\n")
	testInputOutput(t, `<?php var_dump(false < NULL);`, "bool(false)\n")
	testInputOutput(t, `<?php var_dump(false <= NULL);`, "bool(true)\n")
	testInputOutput(t, `<?php var_dump(false <=> NULL);`, "int(0)\n")

	// Floating
	// Floating - Array
	testInputOutput(t, `<?php var_dump(2.5 < []);`, "bool(true)\n")
	testInputOutput(t, `<?php var_dump(2.5 <= []);`, "bool(true)\n")
	testInputOutput(t, `<?php var_dump(2.5 <=> []);`, "int(-1)\n")
	// Floating - Boolean
	testInputOutput(t, `<?php var_dump(2.5 < true);`, "bool(false)\n")
	testInputOutput(t, `<?php var_dump(2.5 <= true);`, "bool(true)\n")
	testInputOutput(t, `<?php var_dump(2.5 <=> true);`, "int(0)\n")
	testInputOutput(t, `<?php var_dump(0.5 < true);`, "bool(false)\n")
	testInputOutput(t, `<?php var_dump(0.5 <= true);`, "bool(true)\n")
	testInputOutput(t, `<?php var_dump(0.5 <=> true);`, "int(0)\n")
	testInputOutput(t, `<?php var_dump(2.5 < false);`, "bool(false)\n")
	testInputOutput(t, `<?php var_dump(2.5 <= false);`, "bool(false)\n")
	testInputOutput(t, `<?php var_dump(2.5 <=> false);`, "int(1)\n")
	testInputOutput(t, `<?php var_dump(0.0 < false);`, "bool(false)\n")
	testInputOutput(t, `<?php var_dump(0.0 <= false);`, "bool(true)\n")
	testInputOutput(t, `<?php var_dump(0.0 <=> false);`, "int(0)\n")
	// Floating - Floating
	testInputOutput(t, `<?php var_dump(2.5 < 3.5);`, "bool(true)\n")
	testInputOutput(t, `<?php var_dump(2.5 <= 3.5);`, "bool(true)\n")
	testInputOutput(t, `<?php var_dump(2.5 <=> 3.5);`, "int(-1)\n")
	testInputOutput(t, `<?php var_dump(3.5 < 3.5);`, "bool(false)\n")
	testInputOutput(t, `<?php var_dump(3.5 <= 3.5);`, "bool(true)\n")
	testInputOutput(t, `<?php var_dump(3.5 <=> 3.5);`, "int(0)\n")
	testInputOutput(t, `<?php var_dump(12.5 < 3.5);`, "bool(false)\n")
	testInputOutput(t, `<?php var_dump(12.5 <= 3.5);`, "bool(false)\n")
	testInputOutput(t, `<?php var_dump(12.5 <=> 3.5);`, "int(1)\n")
	// Floating - Integer
	testInputOutput(t, `<?php var_dump(2.5 < 3);`, "bool(true)\n")
	testInputOutput(t, `<?php var_dump(2.5 <= 3);`, "bool(true)\n")
	testInputOutput(t, `<?php var_dump(2.5 <=> 3);`, "int(-1)\n")
	testInputOutput(t, `<?php var_dump(3.0 < 3);`, "bool(false)\n")
	testInputOutput(t, `<?php var_dump(3.0 <= 3);`, "bool(true)\n")
	testInputOutput(t, `<?php var_dump(3.0 <=> 3);`, "int(0)\n")
	testInputOutput(t, `<?php var_dump(12.5 < 3);`, "bool(false)\n")
	testInputOutput(t, `<?php var_dump(12.5 <= 3);`, "bool(false)\n")
	testInputOutput(t, `<?php var_dump(12.5 <=> 3);`, "int(1)\n")
	// Floating - Null
	testInputOutput(t, `<?php var_dump(2.5 < NULL);`, "bool(false)\n")
	testInputOutput(t, `<?php var_dump(2.5 <= NULL);`, "bool(false)\n")
	testInputOutput(t, `<?php var_dump(2.5 <=> NULL);`, "int(1)\n")
	testInputOutput(t, `<?php var_dump(-2.5 < NULL);`, "bool(true)\n")
	testInputOutput(t, `<?php var_dump(-2.5 <= NULL);`, "bool(true)\n")
	testInputOutput(t, `<?php var_dump(-2.5 <=> NULL);`, "int(-1)\n")

	// Integer
	// Integer - Array
	testInputOutput(t, `<?php var_dump(2 < []);`, "bool(true)\n")
	testInputOutput(t, `<?php var_dump(2 <= []);`, "bool(true)\n")
	testInputOutput(t, `<?php var_dump(2 <=> []);`, "int(-1)\n")
	// Integer - Boolean
	testInputOutput(t, `<?php var_dump(2 < true);`, "bool(false)\n")
	testInputOutput(t, `<?php var_dump(2 <= true);`, "bool(true)\n")
	testInputOutput(t, `<?php var_dump(2 <=> true);`, "int(0)\n")
	testInputOutput(t, `<?php var_dump(2 < false);`, "bool(false)\n")
	testInputOutput(t, `<?php var_dump(2 <= false);`, "bool(false)\n")
	testInputOutput(t, `<?php var_dump(2 <=> false);`, "int(1)\n")
	testInputOutput(t, `<?php var_dump(-2 < true);`, "bool(false)\n")
	testInputOutput(t, `<?php var_dump(-2 <= true);`, "bool(true)\n")
	testInputOutput(t, `<?php var_dump(-2 <=> true);`, "int(0)\n")
	testInputOutput(t, `<?php var_dump(-2 < false);`, "bool(false)\n")
	testInputOutput(t, `<?php var_dump(-2 <= false);`, "bool(false)\n")
	testInputOutput(t, `<?php var_dump(-2 <=> false);`, "int(1)\n")
	// Integer - Floating
	testInputOutput(t, `<?php var_dump(2 < 3.5);`, "bool(true)\n")
	testInputOutput(t, `<?php var_dump(2 <= 3.5);`, "bool(true)\n")
	testInputOutput(t, `<?php var_dump(2 <=> 3.5);`, "int(-1)\n")
	testInputOutput(t, `<?php var_dump(3 < 3.0);`, "bool(false)\n")
	testInputOutput(t, `<?php var_dump(3 <= 3.0);`, "bool(true)\n")
	testInputOutput(t, `<?php var_dump(3 <=> 3.0);`, "int(0)\n")
	testInputOutput(t, `<?php var_dump(12 < 3.5);`, "bool(false)\n")
	testInputOutput(t, `<?php var_dump(12 <= 3.5);`, "bool(false)\n")
	testInputOutput(t, `<?php var_dump(12 <=> 3.5);`, "int(1)\n")
	// Integer - Integer
	testInputOutput(t, `<?php var_dump(2 < 2);`, "bool(false)\n")
	testInputOutput(t, `<?php var_dump(2 <= 2);`, "bool(true)\n")
	testInputOutput(t, `<?php var_dump(2 <=> 2);`, "int(0)\n")
	testInputOutput(t, `<?php var_dump(2 < 0);`, "bool(false)\n")
	testInputOutput(t, `<?php var_dump(2 <= 0);`, "bool(false)\n")
	testInputOutput(t, `<?php var_dump(2 <=> 0);`, "int(1)\n")
	testInputOutput(t, `<?php var_dump(-2 < 2);`, "bool(true)\n")
	testInputOutput(t, `<?php var_dump(-2 <= 2);`, "bool(true)\n")
	testInputOutput(t, `<?php var_dump(-2 <=> 2);`, "int(-1)\n")
	// Integer - Null
	testInputOutput(t, `<?php var_dump(2 < NULL);`, "bool(false)\n")
	testInputOutput(t, `<?php var_dump(2 <= NULL);`, "bool(false)\n")
	testInputOutput(t, `<?php var_dump(2 <=> NULL);`, "int(1)\n")
	testInputOutput(t, `<?php var_dump(0 < NULL);`, "bool(false)\n")
	testInputOutput(t, `<?php var_dump(0 <= NULL);`, "bool(true)\n")
	testInputOutput(t, `<?php var_dump(0 <=> NULL);`, "int(0)\n")
	testInputOutput(t, `<?php var_dump(-2 < NULL);`, "bool(true)\n")
	testInputOutput(t, `<?php var_dump(-2 <= NULL);`, "bool(true)\n")
	testInputOutput(t, `<?php var_dump(-2 <=> NULL);`, "int(-1)\n")

	// Null
	// Null - Array
	testInputOutput(t, `<?php var_dump(NULL < []);`, "bool(false)\n")
	testInputOutput(t, `<?php var_dump(NULL <= []);`, "bool(true)\n")
	testInputOutput(t, `<?php var_dump(NULL <=> []);`, "int(0)\n")
	testInputOutput(t, `<?php var_dump(NULL < ["abc"]);`, "bool(true)\n")
	testInputOutput(t, `<?php var_dump(NULL <= ["abc"]);`, "bool(true)\n")
	testInputOutput(t, `<?php var_dump(NULL <=> ["abc"]);`, "int(-1)\n")
	// Null - Boolean
	testInputOutput(t, `<?php var_dump(NULL < true);`, "bool(true)\n")
	testInputOutput(t, `<?php var_dump(NULL <= true);`, "bool(true)\n")
	testInputOutput(t, `<?php var_dump(NULL <=> true);`, "int(-1)\n")
	testInputOutput(t, `<?php var_dump(NULL < false);`, "bool(false)\n")
	testInputOutput(t, `<?php var_dump(NULL <= false);`, "bool(true)\n")
	testInputOutput(t, `<?php var_dump(NULL <=> false);`, "int(0)\n")
	// Null - Floating
	testInputOutput(t, `<?php var_dump(NULL < 2.5);`, "bool(true)\n")
	testInputOutput(t, `<?php var_dump(NULL <= 2.5);`, "bool(true)\n")
	testInputOutput(t, `<?php var_dump(NULL <=> 2.5);`, "int(-1)\n")
	testInputOutput(t, `<?php var_dump(NULL < -2.5);`, "bool(false)\n")
	testInputOutput(t, `<?php var_dump(NULL <= -2.5);`, "bool(false)\n")
	testInputOutput(t, `<?php var_dump(NULL <=> -2.5);`, "int(1)\n")
	// Null - Integer
	testInputOutput(t, `<?php var_dump(NULL < 2);`, "bool(true)\n")
	testInputOutput(t, `<?php var_dump(NULL <= 2);`, "bool(true)\n")
	testInputOutput(t, `<?php var_dump(NULL <=> 2);`, "int(-1)\n")
	testInputOutput(t, `<?php var_dump(NULL < -2);`, "bool(false)\n")
	testInputOutput(t, `<?php var_dump(NULL <= -2);`, "bool(false)\n")
	testInputOutput(t, `<?php var_dump(NULL <=> -2);`, "int(1)\n")
	// Null - Null
	testInputOutput(t, `<?php var_dump(NULL < NULL);`, "bool(false)\n")
	testInputOutput(t, `<?php var_dump(NULL <= NULL);`, "bool(true)\n")
	testInputOutput(t, `<?php var_dump(NULL <=> NULL);`, "int(0)\n")
	// Null - String
	testInputOutput(t, `<?php var_dump(NULL < "");`, "bool(false)\n")
	testInputOutput(t, `<?php var_dump(NULL <= "");`, "bool(true)\n")
	testInputOutput(t, `<?php var_dump(NULL <=> "");`, "int(0)\n")
	testInputOutput(t, `<?php var_dump(NULL < "abc");`, "bool(true)\n")
	testInputOutput(t, `<?php var_dump(NULL <= "abc");`, "bool(true)\n")
	testInputOutput(t, `<?php var_dump(NULL <=> "abc");`, "int(-1)\n")

	// String
	// String - Array
	testInputOutput(t, `<?php var_dump("" < []);`, "bool(true)\n")
	testInputOutput(t, `<?php var_dump("" <= []);`, "bool(true)\n")
	testInputOutput(t, `<?php var_dump("" <=> []);`, "int(-1)\n")
	// String - Boolean
	testInputOutput(t, `<?php var_dump("" < true);`, "bool(true)\n")
	testInputOutput(t, `<?php var_dump("" <= true);`, "bool(true)\n")
	testInputOutput(t, `<?php var_dump("" <=> true);`, "int(-1)\n")
	testInputOutput(t, `<?php var_dump("a" < true);`, "bool(false)\n")
	testInputOutput(t, `<?php var_dump("a" <= true);`, "bool(true)\n")
	testInputOutput(t, `<?php var_dump("a" <=> true);`, "int(0)\n")
	testInputOutput(t, `<?php var_dump("" < false);`, "bool(false)\n")
	testInputOutput(t, `<?php var_dump("" <= false);`, "bool(true)\n")
	testInputOutput(t, `<?php var_dump("" <=> false);`, "int(0)\n")
	testInputOutput(t, `<?php var_dump("a" < false);`, "bool(false)\n")
	testInputOutput(t, `<?php var_dump("a" <= false);`, "bool(false)\n")
	testInputOutput(t, `<?php var_dump("a" <=> false);`, "int(1)\n")
	// String - Null
	testInputOutput(t, `<?php var_dump("" < NULL);`, "bool(false)\n")
	testInputOutput(t, `<?php var_dump("" <= NULL);`, "bool(true)\n")
	testInputOutput(t, `<?php var_dump("" <=> NULL);`, "int(0)\n")
	testInputOutput(t, `<?php var_dump("22" < NULL);`, "bool(false)\n")
	testInputOutput(t, `<?php var_dump("22" <= NULL);`, "bool(false)\n")
	testInputOutput(t, `<?php var_dump("22" <=> NULL);`, "int(1)\n")
}
