package parser

import (
	"GoPHP/cmd/goPHP/ast"
	"testing"
)

func TestText(t *testing.T) {
	program, _ := NewParser().ProduceAST("<html>")
	expected := ast.NewTextExpression("<html>")
	actual := ast.ExprToTextExpr(ast.StmtToExprStatement(program.GetStatements()[0]).GetExpression())
	if expected.String() != actual.String() || expected.GetValue() != actual.GetValue() {
		t.Errorf("Expected: \"%s\", Got \"%s\"", expected, actual)
	}
}

func TestVariableName(t *testing.T) {
	program, _ := NewParser().ProduceAST("<?php $myVar;")
	expected := ast.NewSimpleVariableExpression(ast.NewVariableNameExpression("$myVar"))
	actual := ast.ExprToSimpleVarExpr(ast.StmtToExprStatement(program.GetStatements()[0]).GetExpression())
	if expected.String() != actual.String() {
		t.Errorf("Expected: \"%s\", Got \"%s\"", expected, actual)
	}

	program, _ = NewParser().ProduceAST("<?php $$myVar;")
	expected = ast.NewSimpleVariableExpression(ast.NewSimpleVariableExpression(ast.NewVariableNameExpression("$myVar")))
	actual = ast.ExprToSimpleVarExpr(ast.StmtToExprStatement(program.GetStatements()[0]).GetExpression())
	if expected.String() != actual.String() {
		t.Errorf("Expected: \"%s\", Got \"%s\"", expected, actual)
	}
}

func TestIntegerLiteral(t *testing.T) {
	// decimal-literal
	program, _ := NewParser().ProduceAST("<?php 42;")
	expected := ast.NewIntegerLiteralExpression(42)
	actual := ast.ExprToIntLitExpr(ast.StmtToExprStatement(program.GetStatements()[0]).GetExpression())
	if expected.String() != actual.String() || expected.GetValue() != actual.GetValue() {
		t.Errorf("Expected: \"%s\", Got \"%s\"", expected, actual)
	}

	// octal-literal
	program, _ = NewParser().ProduceAST("<?php 042;")
	expected = ast.NewIntegerLiteralExpression(34)
	actual = ast.ExprToIntLitExpr(ast.StmtToExprStatement(program.GetStatements()[0]).GetExpression())
	if expected.String() != actual.String() || expected.GetValue() != actual.GetValue() {
		t.Errorf("Expected: \"%s\", Got \"%s\"", expected, actual)
	}

	// hexadecimal-literal
	program, _ = NewParser().ProduceAST("<?php 0x42;")
	expected = ast.NewIntegerLiteralExpression(66)
	actual = ast.ExprToIntLitExpr(ast.StmtToExprStatement(program.GetStatements()[0]).GetExpression())
	if expected.String() != actual.String() || expected.GetValue() != actual.GetValue() {
		t.Errorf("Expected: \"%s\", Got \"%s\"", expected, actual)
	}

	// binary-literal
	program, _ = NewParser().ProduceAST("<?php 0b110110101;")
	expected = ast.NewIntegerLiteralExpression(437)
	actual = ast.ExprToIntLitExpr(ast.StmtToExprStatement(program.GetStatements()[0]).GetExpression())
	if expected.String() != actual.String() || expected.GetValue() != actual.GetValue() {
		t.Errorf("Expected: \"%s\", Got \"%s\"", expected, actual)
	}
}

func TestFloatingLiteral(t *testing.T) {
	program, _ := NewParser().ProduceAST("<?php .5;")
	expected := ast.NewFloatingLiteralExpression(0.5)
	actual := ast.ExprToFloatLitExpr(ast.StmtToExprStatement(program.GetStatements()[0]).GetExpression())
	if expected.String() != actual.String() || expected.GetValue() != actual.GetValue() {
		t.Errorf("Expected: \"%s\", Got \"%s\"", expected, actual)
	}

	program, _ = NewParser().ProduceAST("<?php 1.2;")
	expected = ast.NewFloatingLiteralExpression(1.2)
	actual = ast.ExprToFloatLitExpr(ast.StmtToExprStatement(program.GetStatements()[0]).GetExpression())
	if expected.String() != actual.String() || expected.GetValue() != actual.GetValue() {
		t.Errorf("Expected: \"%s\", Got \"%s\"", expected, actual)
	}

	program, _ = NewParser().ProduceAST("<?php .5e-4;")
	expected = ast.NewFloatingLiteralExpression(0.5e-4)
	actual = ast.ExprToFloatLitExpr(ast.StmtToExprStatement(program.GetStatements()[0]).GetExpression())
	if expected.String() != actual.String() || expected.GetValue() != actual.GetValue() {
		t.Errorf("Expected: \"%s\", Got \"%s\"", expected, actual)
	}

	program, _ = NewParser().ProduceAST("<?php 2.5e+4;")
	expected = ast.NewFloatingLiteralExpression(2.5e+4)
	actual = ast.ExprToFloatLitExpr(ast.StmtToExprStatement(program.GetStatements()[0]).GetExpression())
	if expected.String() != actual.String() || expected.GetValue() != actual.GetValue() {
		t.Errorf("Expected: \"%s\", Got \"%s\"", expected, actual)
	}

	program, _ = NewParser().ProduceAST("<?php 2e4;")
	expected = ast.NewFloatingLiteralExpression(2e4)
	actual = ast.ExprToFloatLitExpr(ast.StmtToExprStatement(program.GetStatements()[0]).GetExpression())
	if expected.String() != actual.String() || expected.GetValue() != actual.GetValue() {
		t.Errorf("Expected: \"%s\", Got \"%s\"", expected, actual)
	}
}

func TestStringLiteral(t *testing.T) {
	program, _ := NewParser().ProduceAST(`<?php b'A "single quoted" \'string\'';`)
	expected := ast.NewStringLiteralExpression(`A "single quoted" \'string\'`, ast.SingleQuotedString)
	actual := ast.ExprToStrLitExpr(ast.StmtToExprStatement(program.GetStatements()[0]).GetExpression())
	if expected.String() != actual.String() ||
		expected.GetValue() != actual.GetValue() || expected.GetStringType() != actual.GetStringType() {
		t.Errorf("Expected: \"%s\", Got \"%s\"", expected, actual)
	}

	program, _ = NewParser().ProduceAST(`<?php 'A "single quoted" \'string\'';`)
	expected = ast.NewStringLiteralExpression(`A "single quoted" \'string\'`, ast.SingleQuotedString)
	actual = ast.ExprToStrLitExpr(ast.StmtToExprStatement(program.GetStatements()[0]).GetExpression())
	if expected.String() != actual.String() ||
		expected.GetValue() != actual.GetValue() || expected.GetStringType() != actual.GetStringType() {
		t.Errorf("Expected: \"%s\", Got \"%s\"", expected, actual)
	}

	program, _ = NewParser().ProduceAST(`<?php b"A \"single quoted\" 'string'";`)
	expected = ast.NewStringLiteralExpression(`A \"single quoted\" 'string'`, ast.DoubleQuotedString)
	actual = ast.ExprToStrLitExpr(ast.StmtToExprStatement(program.GetStatements()[0]).GetExpression())
	if expected.String() != actual.String() ||
		expected.GetValue() != actual.GetValue() || expected.GetStringType() != actual.GetStringType() {
		t.Errorf("Expected: \"%s\", Got \"%s\"", expected, actual)
	}

	program, _ = NewParser().ProduceAST(`<?php "A \"single quoted\" 'string'";`)
	expected = ast.NewStringLiteralExpression(`A \"single quoted\" 'string'`, ast.DoubleQuotedString)
	actual = ast.ExprToStrLitExpr(ast.StmtToExprStatement(program.GetStatements()[0]).GetExpression())
	if expected.String() != actual.String() ||
		expected.GetValue() != actual.GetValue() || expected.GetStringType() != actual.GetStringType() {
		t.Errorf("Expected: \"%s\", Got \"%s\"", expected, actual)
	}
}

func TestEchoStatement(t *testing.T) {
	program, _ := NewParser().ProduceAST(`<?php echo 12, "abc", $var;`)

	expected := ast.NewIntegerLiteralExpression(12)
	actual := ast.ExprToIntLitExpr(ast.StmtToEchoStatement(program.GetStatements()[0]).GetExpressions()[0])
	if expected.String() != actual.String() || expected.GetValue() != actual.GetValue() {
		t.Errorf("Expected: \"%s\", Got \"%s\"", expected, actual)
	}

	expected1 := ast.NewStringLiteralExpression("abc", ast.DoubleQuotedString)
	actual1 := ast.ExprToStrLitExpr(ast.StmtToEchoStatement(program.GetStatements()[0]).GetExpressions()[1])
	if expected1.String() != actual1.String() ||
		expected1.GetValue() != actual1.GetValue() || expected1.GetStringType() != actual1.GetStringType() {
		t.Errorf("Expected: \"%s\", Got \"%s\"", expected1, actual1)
	}

	expected2 := ast.NewSimpleVariableExpression(ast.NewVariableNameExpression("$var"))
	actual2 := ast.ExprToSimpleVarExpr(ast.StmtToEchoStatement(program.GetStatements()[0]).GetExpressions()[2])
	if expected2.String() != actual2.String() || expected2.GetVariableName().GetKind() != actual2.GetVariableName().GetKind() {
		t.Errorf("Expected: \"%s\", Got \"%s\"", expected2, actual2)
	}
}

// TODO continue here
// func TestVariableDeclaration(t *testing.T) {
// 	program, _ := NewParser().ProduceAST(`<?php $variable = "abc";`)
// 	expected := ast.NewIntegerLiteralExpression(12)
// 	actual := ast.ExprToIntLitExpr(ast.StmtToEchoStatement(program.GetStatements()[0]).GetExpressions()[0])
// 	if expected.String() != actual.String() || expected.GetValue() != actual.GetValue() {
// 		t.Errorf("Expected: \"%s\", Got \"%s\"", expected, actual)
// 	}
// }
