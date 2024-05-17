package interpreter

import (
	"GoPHP/cmd/goPHP/ast"
	"testing"
)

// ------------------- MARK: function tests -------------------

func TestVariableExprToVariableName(t *testing.T) {
	// simple-variable-expression

	// $var
	interpreter := NewInterpreter()
	actual, err := interpreter.varExprToVarName(ast.NewSimpleVariableExpression(ast.NewVariableNameExpression("$var")))
	if err != nil {
		t.Errorf("Unexpected error: \"%s\"", err)
		return
	}
	expected := "$var"
	if actual != expected {
		t.Errorf("Expected: \"%s\", Got \"%s\"", expected, actual)
	}

	// $$var
	interpreter = NewInterpreter()
	interpreter.env.declareVariable("$var", NewStringRuntimeValue("hi"))
	actual, err = interpreter.varExprToVarName(
		ast.NewSimpleVariableExpression(
			ast.NewSimpleVariableExpression(ast.NewVariableNameExpression("$var"))))
	if err != nil {
		t.Errorf("Unexpected error: \"%s\"", err)
		return
	}
	expected = "$hi"
	if actual != expected {
		t.Errorf("Expected: \"%s\", Got \"%s\"", expected, actual)
	}

	// $$$var
	interpreter = NewInterpreter()
	interpreter.env.declareVariable("$var1", NewStringRuntimeValue("hi"))
	interpreter.env.declareVariable("$var", NewStringRuntimeValue("var1"))
	actual, err = interpreter.varExprToVarName(
		ast.NewSimpleVariableExpression(
			ast.NewSimpleVariableExpression(
				ast.NewSimpleVariableExpression(ast.NewVariableNameExpression("$var")))))
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

func testInputOutput(t *testing.T, php string, output string) {
	actual, err := NewInterpreter().Process(php)
	if err != nil {
		t.Errorf("Unexpected error: \"%s\"", err)
		return
	}
	if actual != output {
		t.Errorf("Expected: \"%s\", Got \"%s\"", output, actual)
	}
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

func TestCalculation(t *testing.T) {
	// Integer
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
}
