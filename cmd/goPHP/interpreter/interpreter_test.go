package interpreter

import (
	"GoPHP/cmd/goPHP/ast"
	"testing"
)

func TestText(t *testing.T) {
	actual, err := NewInterpreter().Process("<html>...</html>")
	if err != nil {
		t.Errorf("Unexpected error: \"%s\"", err)
		return
	}
	expected := "<html>...</html>"
	if actual != expected {
		t.Errorf("Expected: \"%s\", Got \"%s\"", expected, actual)
	}
}

func TestEchoShortTag(t *testing.T) {
	actual, err := NewInterpreter().Process(`<html><?= "abc" ?><?= 42; ?></html>`)
	if err != nil {
		t.Errorf("Unexpected error: \"%s\"", err)
		return
	}
	expected := "<html>abc42</html>"
	if actual != expected {
		t.Errorf("Expected: \"%s\", Got \"%s\"", expected, actual)
	}
}

func TestEchoExpression(t *testing.T) {
	actual, err := NewInterpreter().Process(`<html><?php echo "abc", 42 ?><?php echo "def", 24; ?></html>`)
	if err != nil {
		t.Errorf("Unexpected error: \"%s\"", err)
		return
	}
	expected := "<html>abc42def24</html>"
	if actual != expected {
		t.Errorf("Expected: \"%s\", Got \"%s\"", expected, actual)
	}
}

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

func TestRuntimeValueToBool(t *testing.T) {
	// boolean to boolean
	actual, err := runtimeValueToBool(NewBooleanRuntimeValue(true))
	if err != nil {
		t.Errorf("Unexpected error: \"%s\"", err)
		return
	}
	expected := true
	if actual != expected {
		t.Errorf("Expected: \"%t\", Got \"%t\"", expected, actual)
	}

	actual, err = runtimeValueToBool(NewBooleanRuntimeValue(false))
	if err != nil {
		t.Errorf("Unexpected error: \"%s\"", err)
		return
	}
	expected = false
	if actual != expected {
		t.Errorf("Expected: \"%t\", Got \"%t\"", expected, actual)
	}

	// integer to boolean
	actual, err = runtimeValueToBool(NewIntegerRuntimeValue(0))
	if err != nil {
		t.Errorf("Unexpected error: \"%s\"", err)
		return
	}
	expected = false
	if actual != expected {
		t.Errorf("Expected: \"%t\", Got \"%t\"", expected, actual)
	}

	actual, err = runtimeValueToBool(NewIntegerRuntimeValue(-0))
	if err != nil {
		t.Errorf("Unexpected error: \"%s\"", err)
		return
	}
	expected = false
	if actual != expected {
		t.Errorf("Expected: \"%t\", Got \"%t\"", expected, actual)
	}

	actual, err = runtimeValueToBool(NewIntegerRuntimeValue(1))
	if err != nil {
		t.Errorf("Unexpected error: \"%s\"", err)
		return
	}
	expected = true
	if actual != expected {
		t.Errorf("Expected: \"%t\", Got \"%t\"", expected, actual)
	}

	actual, err = runtimeValueToBool(NewIntegerRuntimeValue(42))
	if err != nil {
		t.Errorf("Unexpected error: \"%s\"", err)
		return
	}
	expected = true
	if actual != expected {
		t.Errorf("Expected: \"%t\", Got \"%t\"", expected, actual)
	}

	actual, err = runtimeValueToBool(NewIntegerRuntimeValue(-2))
	if err != nil {
		t.Errorf("Unexpected error: \"%s\"", err)
		return
	}
	expected = true
	if actual != expected {
		t.Errorf("Expected: \"%t\", Got \"%t\"", expected, actual)
	}

	// floating to boolean
	actual, err = runtimeValueToBool(NewFloatingRuntimeValue(0.0))
	if err != nil {
		t.Errorf("Unexpected error: \"%s\"", err)
		return
	}
	expected = false
	if actual != expected {
		t.Errorf("Expected: \"%t\", Got \"%t\"", expected, actual)
	}

	actual, err = runtimeValueToBool(NewFloatingRuntimeValue(-0.0))
	if err != nil {
		t.Errorf("Unexpected error: \"%s\"", err)
		return
	}
	expected = false
	if actual != expected {
		t.Errorf("Expected: \"%t\", Got \"%t\"", expected, actual)
	}

	actual, err = runtimeValueToBool(NewFloatingRuntimeValue(1.5))
	if err != nil {
		t.Errorf("Unexpected error: \"%s\"", err)
		return
	}
	expected = true
	if actual != expected {
		t.Errorf("Expected: \"%t\", Got \"%t\"", expected, actual)
	}

	actual, err = runtimeValueToBool(NewFloatingRuntimeValue(42.0))
	if err != nil {
		t.Errorf("Unexpected error: \"%s\"", err)
		return
	}
	expected = true
	if actual != expected {
		t.Errorf("Expected: \"%t\", Got \"%t\"", expected, actual)
	}

	actual, err = runtimeValueToBool(NewFloatingRuntimeValue(-2.0))
	if err != nil {
		t.Errorf("Unexpected error: \"%s\"", err)
		return
	}
	expected = true
	if actual != expected {
		t.Errorf("Expected: \"%t\", Got \"%t\"", expected, actual)
	}

	// string to boolean
	actual, err = runtimeValueToBool(NewStringRuntimeValue(""))
	if err != nil {
		t.Errorf("Unexpected error: \"%s\"", err)
		return
	}
	expected = false
	if actual != expected {
		t.Errorf("Expected: \"%t\", Got \"%t\"", expected, actual)
	}

	actual, err = runtimeValueToBool(NewStringRuntimeValue("0"))
	if err != nil {
		t.Errorf("Unexpected error: \"%s\"", err)
		return
	}
	expected = false
	if actual != expected {
		t.Errorf("Expected: \"%t\", Got \"%t\"", expected, actual)
	}

	actual, err = runtimeValueToBool(NewStringRuntimeValue("Hi"))
	if err != nil {
		t.Errorf("Unexpected error: \"%s\"", err)
		return
	}
	expected = true
	if actual != expected {
		t.Errorf("Expected: \"%t\", Got \"%t\"", expected, actual)
	}

	// null to boolean
	actual, err = runtimeValueToBool(NewNullRuntimeValue())
	if err != nil {
		t.Errorf("Unexpected error: \"%s\"", err)
		return
	}
	expected = false
	if actual != expected {
		t.Errorf("Expected: \"%t\", Got \"%t\"", expected, actual)
	}
}

func TestVariableDeclaration(t *testing.T) {
	// Simple variable
	actual, err := NewInterpreter().Process(`<?php $var = "hi"; $var = "hello"; echo $var, " world";`)
	if err != nil {
		t.Errorf("Unexpected error: \"%s\"", err)
		return
	}
	expected := "hello world"
	if actual != expected {
		t.Errorf("Expected: \"%s\", Got \"%s\"", expected, actual)
	}

	// Variable variable name
	actual, err = NewInterpreter().Process(`<?php $var = "hi"; $$var = "hello"; echo $hi, " world";`)
	if err != nil {
		t.Errorf("Unexpected error: \"%s\"", err)
		return
	}
	expected = "hello world"
	if actual != expected {
		t.Errorf("Expected: \"%s\", Got \"%s\"", expected, actual)
	}

	// Chained variable declarations
	actual, err = NewInterpreter().Process(`<?php $a = $b = $c = 42; echo $a, $b, $c;`)
	if err != nil {
		t.Errorf("Unexpected error: \"%s\"", err)
		return
	}
	expected = "424242"
	if actual != expected {
		t.Errorf("Expected: \"%s\", Got \"%s\"", expected, actual)
	}

	// Compound assignment
	actual, err = NewInterpreter().Process(`<?php $a = 42; echo $a; $a += 2; echo $a; $a += $a; echo $a;`)
	if err != nil {
		t.Errorf("Unexpected error: \"%s\"", err)
		return
	}
	expected = "424488"
	if actual != expected {
		t.Errorf("Expected: \"%s\", Got \"%s\"", expected, actual)
	}
}

func TestConditional(t *testing.T) {
	actual, err := NewInterpreter().Process(`<?php echo 1 ? "a" : "b"; echo 0 ? "b" : "a"; echo false ?: "a";`)
	if err != nil {
		t.Errorf("Unexpected error: \"%s\"", err)
		return
	}
	expected := "aaa"
	if actual != expected {
		t.Errorf("Expected: \"%s\", Got \"%s\"", expected, actual)
	}
}
