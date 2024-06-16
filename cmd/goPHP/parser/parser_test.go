package parser

import (
	"GoPHP/cmd/goPHP/ast"
	"GoPHP/cmd/goPHP/ini"
	"testing"
)

func testExpr(t *testing.T, php string, expected ast.IExpression) {
	testExprs(t, php, []ast.IExpression{expected})
}

func testExprs(t *testing.T, php string, expected []ast.IExpression) {
	program, err := NewParser(ini.NewDevIni()).ProduceAST(php, "test.php")
	if err != nil {
		t.Errorf("Unexpected error: \"%s\"", err)
		return
	}
	for index, expect := range expected {
		actual := ast.StmtToExprStmt(program.GetStatements()[index]).GetExpression()
		if expect.String() != actual.String() {
			t.Errorf("\nExpected: \"%s\"\nGot       \"%s\"", expect, actual)
			return
		}
	}
}

func testStmt(t *testing.T, php string, expected ast.IStatement) {
	testStmts(t, php, []ast.IStatement{expected})
}

func testStmts(t *testing.T, php string, expected []ast.IStatement) {
	program, err := NewParser(ini.NewDevIni()).ProduceAST(php, "test.php")
	if err != nil {
		t.Errorf("Unexpected error: \"%s\"", err)
		return
	}
	for index, expect := range expected {
		actual := program.GetStatements()[index]
		if expect.String() != actual.String() {
			t.Errorf("\nExpected: \"%s\"\nGot       \"%s\"", expect, actual)
			return
		}
	}
}

func TestText(t *testing.T) {
	testExpr(t, "<html>", ast.NewTextExpr("<html>"))
}

func TestVariable(t *testing.T) {
	// Lookup
	testExpr(t, "<?php $myVar;", ast.NewSimpleVariableExpr(ast.NewVariableNameExpr(nil, "$myVar")))
	testExpr(t, "<?php $$myVar;", ast.NewSimpleVariableExpr(ast.NewSimpleVariableExpr(ast.NewVariableNameExpr(nil, "$myVar"))))

	// Simple assignment
	testExpr(t, `<?php $variable = "abc";`, ast.NewSimpleAssignmentExpr(
		ast.NewSimpleVariableExpr(ast.NewVariableNameExpr(nil, "$variable")),
		ast.NewStringLiteralExpr(nil, "abc", ast.DoubleQuotedString),
	))

	// Compound assignment
	testExprs(t, `<?php $a = 42; $a += 2;`, []ast.IExpression{
		ast.NewSimpleAssignmentExpr(ast.NewSimpleVariableExpr(ast.NewVariableNameExpr(nil, "$a")), ast.NewIntegerLiteralExpr(nil, 42)),
		ast.NewCompoundAssignmentExpr(ast.NewSimpleVariableExpr(ast.NewVariableNameExpr(nil, "$a")), "+", ast.NewIntegerLiteralExpr(nil, 2)),
	})
}

func TestArray(t *testing.T) {
	// Subscript
	testExpr(t, "<?php $myVar[];", ast.NewSubscriptExpr(ast.NewSimpleVariableExpr(ast.NewVariableNameExpr(nil, "$myVar")), nil))
	testExpr(t, "<?php $myVar[0];",
		ast.NewSubscriptExpr(ast.NewSimpleVariableExpr(ast.NewVariableNameExpr(nil, "$myVar")), ast.NewIntegerLiteralExpr(nil, 0)),
	)
}

func TestFunctionCall(t *testing.T) {
	// Without argument
	testExpr(t, "<?php func();", ast.NewFunctionCallExpr(nil, "func", []ast.IExpression{}))
	// With argument
	testExpr(t, "<?php func(42);", ast.NewFunctionCallExpr(nil, "func", []ast.IExpression{ast.NewIntegerLiteralExpr(nil, 42)}))
}

func TestLiteral(t *testing.T) {
	// Array literal
	expected := ast.NewArrayLiteralExpr(nil)
	expected.AddElement(ast.NewIntegerLiteralExpr(nil, 0), ast.NewIntegerLiteralExpr(nil, 1))
	expected.AddElement(ast.NewIntegerLiteralExpr(nil, 1), ast.NewStringLiteralExpr(nil, "a", ast.DoubleQuotedString))
	expected.AddElement(ast.NewIntegerLiteralExpr(nil, 2), ast.NewBooleanLiteralExpr(nil, false))
	testExpr(t, `<?php [1, "a", false];`, expected)
	testExpr(t, `<?php array(1, "a", false);`, expected)

	// Boolean literal
	testExpr(t, "<?php true;", ast.NewBooleanLiteralExpr(nil, true))
	testExpr(t, "<?php false;", ast.NewBooleanLiteralExpr(nil, false))

	// Null literal
	testExpr(t, "<?php null;", ast.NewNullLiteralExpr(nil))

	// Integer literal
	// decimal-literal
	testExpr(t, "<?php 42;", ast.NewIntegerLiteralExpr(nil, 42))
	// octal-literal
	testExpr(t, "<?php 042;", ast.NewIntegerLiteralExpr(nil, 34))
	// hexadecimal-literal
	testExpr(t, "<?php 0x42;", ast.NewIntegerLiteralExpr(nil, 66))
	// binary-literal
	testExpr(t, "<?php 0b110110101;", ast.NewIntegerLiteralExpr(nil, 437))

	// Floating literal
	testExpr(t, "<?php .5;", ast.NewFloatingLiteralExpr(nil, 0.5))
	testExpr(t, "<?php 1.2;", ast.NewFloatingLiteralExpr(nil, 1.2))
	testExpr(t, "<?php .5e-4;", ast.NewFloatingLiteralExpr(nil, 0.5e-4))
	testExpr(t, "<?php 2.5e+4;", ast.NewFloatingLiteralExpr(nil, 2.5e+4))
	testExpr(t, "<?php 2e4;", ast.NewFloatingLiteralExpr(nil, 2e4))

	// String literal
	testExpr(t, `<?php b'A "single quoted" \'string\'';`, ast.NewStringLiteralExpr(nil, `A "single quoted" 'string'`, ast.SingleQuotedString))
	testExpr(t, `<?php 'A "single quoted" \'string\'';`, ast.NewStringLiteralExpr(nil, `A "single quoted" 'string'`, ast.SingleQuotedString))
	testExpr(t, `<?php b"A \"double quoted\" 'string'";`, ast.NewStringLiteralExpr(nil, `A "double quoted" 'string'`, ast.DoubleQuotedString))
	testExpr(t, `<?php "A \"double quoted\" 'string'";`, ast.NewStringLiteralExpr(nil, `A "double quoted" 'string'`, ast.DoubleQuotedString))
}

func TestEchoStatement(t *testing.T) {
	testStmt(t, `<?php echo 12, "abc", $var;`, ast.NewEchoStmt(nil, []ast.IExpression{
		ast.NewIntegerLiteralExpr(nil, 12), ast.NewStringLiteralExpr(nil, "abc", ast.DoubleQuotedString),
		ast.NewSimpleVariableExpr(ast.NewVariableNameExpr(nil, "$var")),
	}))
}

func TestConditional(t *testing.T) {
	// Conditional
	testExpr(t, `<?php 1 ? "a" : "b";`, ast.NewConditionalExpr(
		ast.NewIntegerLiteralExpr(nil, 1), ast.NewStringLiteralExpr(nil, "a", ast.DoubleQuotedString),
		ast.NewStringLiteralExpr(nil, "b", ast.DoubleQuotedString),
	))

	testExpr(t, `<?php 1 ?: "b";`,
		ast.NewConditionalExpr(ast.NewIntegerLiteralExpr(nil, 1), nil, ast.NewStringLiteralExpr(nil, "b", ast.DoubleQuotedString)),
	)

	testExpr(t, `<?php 1 ? "a" : 2 ? "b": "c";`, ast.NewConditionalExpr(
		ast.NewIntegerLiteralExpr(nil, 1), ast.NewStringLiteralExpr(nil, "a", ast.DoubleQuotedString),
		ast.NewConditionalExpr(
			ast.NewIntegerLiteralExpr(nil, 2), ast.NewStringLiteralExpr(nil, "b", ast.DoubleQuotedString),
			ast.NewStringLiteralExpr(nil, "c", ast.DoubleQuotedString),
		),
	))

	// Coalesce
	testExpr(t, `<?php "a" ?? "b";`, ast.NewCoalesceExpr(
		ast.NewStringLiteralExpr(nil, "a", ast.DoubleQuotedString), ast.NewStringLiteralExpr(nil, "b", ast.DoubleQuotedString),
	))

	testExpr(t, `<?php "a" ?? "b" ?? "c";`, ast.NewCoalesceExpr(
		ast.NewStringLiteralExpr(nil, "a", ast.DoubleQuotedString),
		ast.NewCoalesceExpr(
			ast.NewStringLiteralExpr(nil, "b", ast.DoubleQuotedString), ast.NewStringLiteralExpr(nil, "c", ast.DoubleQuotedString),
		),
	))
}

func TestCastExpression(t *testing.T) {
	testExpr(t, `<?php (string)42;`, ast.NewCastExpr(nil, "string", ast.NewIntegerLiteralExpr(nil, 42)))
}

func TestParenthesizedExpression(t *testing.T) {
	testExpr(t, `<?php (1+2);`, ast.NewBinaryOpExpr(ast.NewIntegerLiteralExpr(nil, 1), "+", ast.NewIntegerLiteralExpr(nil, 2)))
}

func TestConstDeclaration(t *testing.T) {
	testStmts(t, `<?php const PI = 3.141, ZERO = 0;`, []ast.IStatement{
		ast.NewConstDeclarationStmt(nil, "PI", ast.NewFloatingLiteralExpr(nil, 3.141)),
		ast.NewConstDeclarationStmt(nil, "ZERO", ast.NewIntegerLiteralExpr(nil, 0)),
	})
}

func TestEqualityExpression(t *testing.T) {
	testExpr(t, `<?php "234" !== true;`, ast.NewEqualityExpr(
		ast.NewStringLiteralExpr(nil, "234", ast.DoubleQuotedString), "!==", ast.NewBooleanLiteralExpr(nil, true),
	))
	testExpr(t, `<?php "234" == true;`, ast.NewEqualityExpr(
		ast.NewStringLiteralExpr(nil, "234", ast.DoubleQuotedString), "==", ast.NewBooleanLiteralExpr(nil, true),
	))
}

func TestOperatorExpression(t *testing.T) {
	// Shift
	testExpr(t, `<?php "234" << 12;`,
		ast.NewBinaryOpExpr(ast.NewStringLiteralExpr(nil, "234", ast.DoubleQuotedString), "<<", ast.NewIntegerLiteralExpr(nil, 12)),
	)

	// Additive
	testExpr(t, `<?php "234" + 12;`,
		ast.NewBinaryOpExpr(ast.NewStringLiteralExpr(nil, "234", ast.DoubleQuotedString), "+", ast.NewIntegerLiteralExpr(nil, 12)),
	)

	// Multiplicative
	testExpr(t, `<?php "234" * 12;`,
		ast.NewBinaryOpExpr(ast.NewStringLiteralExpr(nil, "234", ast.DoubleQuotedString), "*", ast.NewIntegerLiteralExpr(nil, 12)),
	)

	// Exponentiation
	testExpr(t, `<?php "234" ** 12;`,
		ast.NewBinaryOpExpr(ast.NewStringLiteralExpr(nil, "234", ast.DoubleQuotedString), "**", ast.NewIntegerLiteralExpr(nil, 12)),
	)

	// Logical not
	testExpr(t, `<?php !true;`, ast.NewLogicalNotExpr(nil, ast.NewBooleanLiteralExpr(nil, true)))

	// Logical inc or
	testExpr(t, `<?php "234" || 12;`,
		ast.NewBinaryOpExpr(ast.NewStringLiteralExpr(nil, "234", ast.DoubleQuotedString), "||", ast.NewIntegerLiteralExpr(nil, 12)),
	)

	// Logical and
	testExpr(t, `<?php "234" && 12;`,
		ast.NewBinaryOpExpr(ast.NewStringLiteralExpr(nil, "234", ast.DoubleQuotedString), "&&", ast.NewIntegerLiteralExpr(nil, 12)),
	)

	// Bitwise inc or
	testExpr(t, `<?php "234" | 12;`,
		ast.NewBinaryOpExpr(ast.NewStringLiteralExpr(nil, "234", ast.DoubleQuotedString), "|", ast.NewIntegerLiteralExpr(nil, 12)),
	)

	// Bitwise exc or
	testExpr(t, `<?php "234" ^ 12;`,
		ast.NewBinaryOpExpr(ast.NewStringLiteralExpr(nil, "234", ast.DoubleQuotedString), "^", ast.NewIntegerLiteralExpr(nil, 12)),
	)

	// Bitwise and
	testExpr(t, `<?php "234" & 12;`,
		ast.NewBinaryOpExpr(ast.NewStringLiteralExpr(nil, "234", ast.DoubleQuotedString), "&", ast.NewIntegerLiteralExpr(nil, 12)),
	)
}

func TestIntrinsic(t *testing.T) {
	// Die
	testExpr(t, `<?php die(42);`, ast.NewExitIntrinsic(nil, ast.NewIntegerLiteralExpr(nil, 42)))
	testExpr(t, `<?php die();`, ast.NewExitIntrinsic(nil, nil))
	testExpr(t, `<?php die;`, ast.NewExitIntrinsic(nil, nil))
	// Exit
	testExpr(t, `<?php exit(42);`, ast.NewExitIntrinsic(nil, ast.NewIntegerLiteralExpr(nil, 42)))
	testExpr(t, `<?php exit();`, ast.NewExitIntrinsic(nil, nil))
	testExpr(t, `<?php exit;`, ast.NewExitIntrinsic(nil, nil))

	// Empty
	testExpr(t, `<?php empty(false);`, ast.NewEmptyIntrinsic(nil, ast.NewBooleanLiteralExpr(nil, false)))

	// Isset
	testExpr(t, `<?php isset($a);`,
		ast.NewIssetIntrinsic(nil, []ast.IExpression{ast.NewSimpleVariableExpr(ast.NewVariableNameExpr(nil, "$a"))}),
	)

	// Unset
	testExpr(t, `<?php unset($a);`,
		ast.NewUnsetIntrinsic(nil, []ast.IExpression{ast.NewSimpleVariableExpr(ast.NewVariableNameExpr(nil, "$a"))}),
	)
}
