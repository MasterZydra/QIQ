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
	// TODO additive and multiplicative expressions
	// Integer
	testInputOutput(t, `<?php $a = 42; echo $a += 1;`, "43")
	testInputOutput(t, `<?php $a = 42; echo $a -= 1;`, "41")
	testInputOutput(t, `<?php $a = 42; echo $a *= 2;`, "84")
	testInputOutput(t, `<?php $a = 42; echo $a /= 2;`, "21")
	testInputOutput(t, `<?php $a = 42; echo $a %= 5;`, "2")
	testInputOutput(t, `<?php $a = 2; echo $a **= 4;`, "16")

	// Floating
	testInputOutput(t, `<?php $a = 42.0; echo $a += 1;`, "43")
	testInputOutput(t, `<?php $a = 42.0; echo $a -= 1;`, "41")
	testInputOutput(t, `<?php $a = 42.0; echo $a *= 2;`, "84")
	testInputOutput(t, `<?php $a = 42.0; echo $a /= 2;`, "21")
	testInputOutput(t, `<?php $a = 2.0; echo $a **= 4;`, "16")
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
