package interpreter

import (
	"GoPHP/cmd/goPHP/ast"
	"fmt"
	"testing"
)

// ------------------- MARK: function tests -------------------

func TestVariableExprToVariableName(t *testing.T) {
	// simple-variable-expression

	// $var
	interpreter := NewInterpreter(NewDevConfig(), &Request{})
	actual, err := interpreter.varExprToVarName(ast.NewSimpleVariableExpression(ast.NewVariableNameExpression("$var")), interpreter.env)
	if err != nil {
		t.Errorf("Unexpected error: \"%s\"", err)
		return
	}
	expected := "$var"
	if actual != expected {
		t.Errorf("Expected: \"%s\", Got \"%s\"", expected, actual)
	}

	// $$var
	interpreter = NewInterpreter(NewDevConfig(), &Request{})
	interpreter.env.declareVariable("$var", NewStringRuntimeValue("hi"))
	actual, err = interpreter.varExprToVarName(
		ast.NewSimpleVariableExpression(
			ast.NewSimpleVariableExpression(ast.NewVariableNameExpression("$var"))), interpreter.env)
	if err != nil {
		t.Errorf("Unexpected error: \"%s\"", err)
		return
	}
	expected = "$hi"
	if actual != expected {
		t.Errorf("Expected: \"%s\", Got \"%s\"", expected, actual)
	}

	// $$$var
	interpreter = NewInterpreter(NewDevConfig(), &Request{})
	interpreter.env.declareVariable("$var1", NewStringRuntimeValue("hi"))
	interpreter.env.declareVariable("$var", NewStringRuntimeValue("var1"))
	actual, err = interpreter.varExprToVarName(
		ast.NewSimpleVariableExpression(
			ast.NewSimpleVariableExpression(
				ast.NewSimpleVariableExpression(ast.NewVariableNameExpression("$var")))), interpreter.env)
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

func testInputOutput(t *testing.T, php string, output string) *Interpreter {
	interpreter := NewInterpreter(NewDevConfig(), &Request{})
	actual, err := interpreter.Process(php)
	if err != nil {
		fmt.Println("    Code:", php)
		t.Errorf("Unexpected error: \"%s\"", err)
		return interpreter
	}
	if actual != output {
		fmt.Println("    Code:", php)
		t.Errorf("Expected: \"%s\", Got \"%s\"", output, actual)
	}
	return interpreter
}

func TestText(t *testing.T) {
	testInputOutput(t, "<html>...</html>", "<html>...</html>")
}

func TestEchoShortTag(t *testing.T) {
	testInputOutput(t, `<html><?= "abc" ?><?= 42; ?></html>`, "<html>abc42</html>")
}

func TestEchoExpression(t *testing.T) {
	testInputOutput(t,
		`<html><?php echo "abc", 42 ?><?php echo "def", 24; ?></html>`,
		"<html>abc42def24</html>",
	)
}

func TestStringVariableSubstitution(t *testing.T) {
	testInputOutput(t, `<?php $a = 42; echo "a{$a}b";`, "a42b")
}

func TestVariableDeclaration(t *testing.T) {
	// Simple variable
	testInputOutput(t,
		`<?php $var = "hi"; $var = "hello"; echo $var, " world";`,
		"hello world",
	)

	// Variable variable name
	testInputOutput(t,
		`<?php $var = "hi"; $$var = "hello"; echo $hi, " world";`,
		"hello world",
	)

	// Chained variable declarations
	testInputOutput(t,
		`<?php $a = $b = $c = 42; echo $a, $b, $c;`,
		"424242",
	)

	// Compound assignment
	testInputOutput(t,
		`<?php $a = 42; echo $a; $a += 2; echo $a; $a += $a; echo $a;`,
		"424488",
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

func TestPredefinedConstants(t *testing.T) {
	testInputOutput(t, `<?php echo E_USER_NOTICE;`, fmt.Sprintf("%d", E_USER_NOTICE))
	testInputOutput(t, `<?php echo E_ALL;`, fmt.Sprintf("%d", E_ALL))
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

func TestConstantDeclaration(t *testing.T) {
	testInputOutput(t,
		`<?php const TRUTH = 42; const PI = "3.141";echo TRUTH, PI;`,
		"423.141",
	)
}

func TestConditional(t *testing.T) {
	testInputOutput(t,
		`<?php echo 1 ? "a" : "b"; echo 0 ? "b" : "a"; echo false ?: "a";`,
		"aaa",
	)
}

func TestCoalesce(t *testing.T) {
	testInputOutput(t,
		`<?php $a = null; echo $a ?? "a"; $a = "b"; echo $a ?? "a"; echo "c" ?? "d";`,
		"abc",
	)
	testInputOutput(t, `<?php echo $a ?? "a";`, "a")
}

func TestCalculation(t *testing.T) {
	// Boolean
	testInputOutput(t, `<?php echo 4 && 0 ? "t" : "f";`, "f")
	testInputOutput(t, `<?php echo 4 && 1 ? "t" : "f";`, "t")
	testInputOutput(t, `<?php echo 4 && false ? "t" : "f";`, "f")
	testInputOutput(t, `<?php echo 4 && true ? "t" : "f";`, "t")
	testInputOutput(t, `<?php echo 0 || 0 ? "t" : "f";`, "f")
	testInputOutput(t, `<?php echo 4 || 1 ? "t" : "f";`, "t")
	testInputOutput(t, `<?php echo false || false ? "t" : "f";`, "f")
	testInputOutput(t, `<?php echo 4 || true ? "t" : "f";`, "t")

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

	// Combined additions and multiplications
	testInputOutput(t, `<?php echo 31 + 21 + 11;`, "63")
	testInputOutput(t, `<?php echo 4 * 3 * 2;`, "24")
	testInputOutput(t, `<?php echo 2 + 3 * 4;`, "14")
	testInputOutput(t, `<?php echo (2 + 3) * 4;`, "20")
	testInputOutput(t, `<?php echo 2 * 3 + 4 * 5 + 6;`, "32")
	testInputOutput(t, `<?php echo 2 * (3 + 4) * 5 + 6;`, "76")
	testInputOutput(t, `<?php echo 2 + 3 * 4 + 5 * 6;`, "44")

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
}

func TestComparison(t *testing.T) {
	// ===
	testInputOutput(t, `<?php echo "abc" === "abc" ? "a" : "b";`, "a")
	testInputOutput(t, `<?php echo "abc" !== "abc" ? "a" : "b";`, "b")
	testInputOutput(t, `<?php echo "abc" !== "abcd" ? "a" : "b";`, "a")
	testInputOutput(t, `<?php echo "abc" === "abcd" ? "a" : "b";`, "b")
	testInputOutput(t, `<?php echo "123" !== 123 ? "a" : "b";`, "a")
	testInputOutput(t, `<?php echo "123" === 123 ? "a" : "b";`, "b")
}

func TestIncDec(t *testing.T) {
	// Boolean
	testInputOutput(t, `<?php $a = true; var_dump($a++); var_dump($a);`, "bool(true)\nbool(true)\n")
	testInputOutput(t, `<?php $a = true; var_dump($a--); var_dump($a);`, "bool(true)\nbool(true)\n")
	testInputOutput(t, `<?php $a = true; var_dump(++$a); var_dump($a);`, "bool(true)\nbool(true)\n")
	testInputOutput(t, `<?php $a = true; var_dump(--$a); var_dump($a);`, "bool(true)\nbool(true)\n")
	testInputOutput(t, `<?php $a = false; var_dump($a++); var_dump($a);`, "bool(false)\nbool(false)\n")
	testInputOutput(t, `<?php $a = false; var_dump($a--); var_dump($a);`, "bool(false)\nbool(false)\n")
	testInputOutput(t, `<?php $a = false; var_dump(++$a); var_dump($a);`, "bool(false)\nbool(false)\n")
	testInputOutput(t, `<?php $a = false; var_dump(--$a); var_dump($a);`, "bool(false)\nbool(false)\n")

	// Floating
	testInputOutput(t, `<?php $a = 42.0; var_dump($a++); var_dump($a);`, "float(42)\nfloat(43)\n")
	testInputOutput(t, `<?php $a = 42.0; var_dump($a--); var_dump($a);`, "float(42)\nfloat(41)\n")
	testInputOutput(t, `<?php $a = 42.0; var_dump(++$a); var_dump($a);`, "float(43)\nfloat(43)\n")
	testInputOutput(t, `<?php $a = 42.0; var_dump(--$a); var_dump($a);`, "float(41)\nfloat(41)\n")

	// Integer
	testInputOutput(t, `<?php $a = 42; var_dump($a++); var_dump($a);`, "int(42)\nint(43)\n")
	testInputOutput(t, `<?php $a = 42; var_dump($a--); var_dump($a);`, "int(42)\nint(41)\n")
	testInputOutput(t, `<?php $a = 42; var_dump(++$a); var_dump($a);`, "int(43)\nint(43)\n")
	testInputOutput(t, `<?php $a = 42; var_dump(--$a); var_dump($a);`, "int(41)\nint(41)\n")

	// Null
	testInputOutput(t, `<?php $a = null; var_dump($a++); var_dump($a);`, "NULL\nint(1)\n")
	testInputOutput(t, `<?php $a = null; var_dump($a--); var_dump($a);`, "NULL\nNULL\n")
	testInputOutput(t, `<?php $a = null; var_dump(++$a); var_dump($a);`, "int(1)\nint(1)\n")
	testInputOutput(t, `<?php $a = null; var_dump(--$a); var_dump($a);`, "NULL\nNULL\n")

	// String
	testInputOutput(t, `<?php $a = ""; var_dump($a++); var_dump($a);`, `string(0) ""`+"\n"+`string(1) "1"`+"\n")
	testInputOutput(t, `<?php $a = ""; var_dump($a--); var_dump($a);`, `string(0) ""`+"\nint(-1)\n")
	testInputOutput(t, `<?php $a = ""; var_dump(++$a); var_dump($a);`, `string(1) "1"`+"\n"+`string(1) "1"`+"\n")
	testInputOutput(t, `<?php $a = ""; var_dump(--$a); var_dump($a);`, "int(-1)\nint(-1)\n")
}

func TestUnaryExpression(t *testing.T) {
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
}

func TestLogicalExpression(t *testing.T) {
	// Not
	testInputOutput(t, `<?php echo !true ? "a" : "b";`, "b")
	testInputOutput(t, `<?php echo !false ? "a" : "b";`, "a")
	testInputOutput(t, `<?php echo !42 ? "a" : "b";`, "b")
}

func TestWarningUndefinedVariable(t *testing.T) {
	testInputOutput(t,
		`<?php echo is_null($a) ? "a" : "b";`,
		"Warning: Undefined variable $a\na",
	)

	testInputOutput(t,
		`<?php echo intval($a);`,
		"Warning: Undefined variable $a\n0",
	)

	testInputOutput(t,
		`<?php echo intval($$a);`,
		"Warning: Undefined variable $a\nWarning: Undefined variable $\n0",
	)
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
	testInputOutput(t, `<?php echo empty(false) ? "a" : "b";`, "a")
	testInputOutput(t, `<?php echo empty(true) ? "a" : "b";`, "b")
	testInputOutput(t, `<?php echo empty(0) ? "a" : "b";`, "a")
	testInputOutput(t, `<?php echo empty(1) ? "a" : "b";`, "b")
	testInputOutput(t, `<?php echo empty(0.0) ? "a" : "b";`, "a")
	testInputOutput(t, `<?php echo empty(2.0) ? "a" : "b";`, "b")
	testInputOutput(t, `<?php echo empty("") ? "a" : "b";`, "a")
	testInputOutput(t, `<?php echo empty("0") ? "a" : "b";`, "a")
	testInputOutput(t, `<?php echo empty("1") ? "a" : "b";`, "b")
	testInputOutput(t, `<?php echo empty("00") ? "a" : "b";`, "b")
	testInputOutput(t, `<?php echo empty(null) ? "a" : "b";`, "a")
	testInputOutput(t, `<?php echo empty($a) ? "a" : "b";`, "a")
	testInputOutput(t, `<?php $a = 1; echo empty($a) ? "a" : "b";`, "b")

	// Isset
	testInputOutput(t, `<?php $a = 1; echo isset($a) ? "a" : "b";`, "a")
	testInputOutput(t, `<?php echo isset($a) ? "a" : "b";`, "b")
	testInputOutput(t, `<?php $a = 1; echo isset($a, $b) ? "a" : "b";`, "b")
	testInputOutput(t, `<?php echo isset($a, $b) ? "a" : "b";`, "b")

	// Unset
	testInputOutput(t, `<?php $a = 1; echo isset($a) ? "y" : "n"; unset($a); echo isset($a) ? "y" : "n";`, "yn")
	testInputOutput(t, `<?php echo isset($a) ? "y" : "n"; unset($a); echo isset($a) ? "y" : "n";`, "nn")
}

func TestCompoundStmt(t *testing.T) {
	testInputOutput(t, `<?php { echo "1"; echo "2";} {}`, "12")
}

func TestIfStmt(t *testing.T) {
	testInputOutput(t, `<?php $a = 42; if ($a === 42) { echo "42"; } elseif ($a === 41) { echo "41"; } else { echo "??"; }`, "42")
	testInputOutput(t, `<?php $a = 41; if ($a === 42) { echo "42"; } elseif ($a === 41) { echo "41"; } else { echo "??"; }`, "41")
	testInputOutput(t, `<?php $a = 40; if ($a === 42) { echo "42"; } elseif ($a === 41) { echo "41"; } else { echo "??"; }`, "??")
	// Alternative syntax
	testInputOutput(t, `<?php $a = 42; if ($a === 42): echo "42"; elseif ($a === 41): echo "41"; else: echo "??"; endif;`, "42")
	testInputOutput(t, `<?php $a = 41; if ($a === 42): echo "42"; elseif ($a === 41): echo "41"; else: echo "??"; endif;`, "41")
	testInputOutput(t, `<?php $a = 40; if ($a === 42): echo "42"; elseif ($a === 41): echo "41"; else: echo "??"; endif;`, "??")
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
	// TODO if "compareRelationString" - string is implemented
	// testInputOutput(t, `<?php var_dump(NULL < "");`, "bool(false)\n")
	// testInputOutput(t, `<?php var_dump(NULL <= "");`, "bool(true)\n")
	// testInputOutput(t, `<?php var_dump(NULL <=> "");`, "int(0)\n")
	// testInputOutput(t, `<?php var_dump(NULL < "abc");`, "bool(true)\n")
	// testInputOutput(t, `<?php var_dump(NULL <= "abc");`, "bool(true)\n")
	// testInputOutput(t, `<?php var_dump(NULL <=> "abc");`, "int(-1)\n")

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
	// TODO if "compareRelationString" - string is implemented
	// testInputOutput(t, `<?php var_dump("" < NULL);`, "bool(false)\n")
	// testInputOutput(t, `<?php var_dump("" <= NULL);`, "bool(true)\n")
	// testInputOutput(t, `<?php var_dump("" <=> NULL);`, "int(0)\n")
	// testInputOutput(t, `<?php var_dump("22" < NULL);`, "bool(false)\n")
	// testInputOutput(t, `<?php var_dump("22" <= NULL);`, "bool(false)\n")
	// testInputOutput(t, `<?php var_dump("22" <=> NULL);`, "int(1)\n")
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
}
