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
		actual := program.GetStatements()[index].(*ast.ExpressionStatement).Expr
		if ast.ToString(expect) != ast.ToString(actual) {
			t.Errorf("\nExpected: \"%s\"\nGot       \"%s\"", ast.ToString(expect), ast.ToString(actual))
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
		if ast.ToString(expect) != ast.ToString(actual) {
			t.Errorf("\nExpected: \"%s\"\nGot       \"%s\"", ast.ToString(expect), ast.ToString(actual))
			return
		}
	}
}

func TestText(t *testing.T) {
	testExpr(t, "<html>", ast.NewTextExpr(0, "<html>"))
}

func TestVariable(t *testing.T) {
	// Lookup
	testExpr(t, "<?php $myVar;", ast.NewSimpleVariableExpr(0, ast.NewVariableNameExpr(0, nil, "$myVar")))
	testExpr(t, "<?php $$myVar;", ast.NewSimpleVariableExpr(0, ast.NewSimpleVariableExpr(0, ast.NewVariableNameExpr(0, nil, "$myVar"))))

	// Simple assignment
	testExpr(t, `<?php $variable = "abc";`, ast.NewSimpleAssignmentExpr(0,
		ast.NewSimpleVariableExpr(0, ast.NewVariableNameExpr(0, nil, "$variable")),
		ast.NewStringLiteralExpr(0, nil, "abc", ast.DoubleQuotedString),
	))

	// Compound assignment
	testExprs(t, `<?php $a = 42; $a += 2;`, []ast.IExpression{
		ast.NewSimpleAssignmentExpr(0, ast.NewSimpleVariableExpr(0, ast.NewVariableNameExpr(0, nil, "$a")), ast.NewIntegerLiteralExpr(0, nil, 42)),
		ast.NewCompoundAssignmentExpr(0, ast.NewSimpleVariableExpr(0, ast.NewVariableNameExpr(0, nil, "$a")), "+", ast.NewIntegerLiteralExpr(0, nil, 2)),
	})
}

func TestArray(t *testing.T) {
	// Subscript
	testExpr(t, "<?php $myVar[];", ast.NewSubscriptExpr(0, ast.NewSimpleVariableExpr(0, ast.NewVariableNameExpr(0, nil, "$myVar")), nil))
	testExpr(t, "<?php $myVar[0];",
		ast.NewSubscriptExpr(0, ast.NewSimpleVariableExpr(0, ast.NewVariableNameExpr(0, nil, "$myVar")), ast.NewIntegerLiteralExpr(0, nil, 0)),
	)
}

func TestFunctionCall(t *testing.T) {
	// Without argument
	testExpr(t, "<?php func();", ast.NewFunctionCallExpr(0, nil, ast.NewStringLiteralExpr(0, nil, "func", ast.SingleQuotedString), []ast.IExpression{}))
	// With argument
	testExpr(t, "<?php func(42);", ast.NewFunctionCallExpr(0, nil, ast.NewStringLiteralExpr(0, nil, "func", ast.SingleQuotedString), []ast.IExpression{ast.NewIntegerLiteralExpr(0, nil, 42)}))
}

func TestLiteral(t *testing.T) {
	// Array literal
	expected := ast.NewArrayLiteralExpr(0, nil)
	expected.AddElement(ast.NewIntegerLiteralExpr(0, nil, 0), ast.NewIntegerLiteralExpr(0, nil, 1))
	expected.AddElement(ast.NewIntegerLiteralExpr(0, nil, 1), ast.NewStringLiteralExpr(0, nil, "a", ast.DoubleQuotedString))
	expected.AddElement(ast.NewIntegerLiteralExpr(0, nil, 2), ast.NewBooleanLiteralExpr(0, nil, false))
	testExpr(t, `<?php [1, "a", false];`, expected)
	testExpr(t, `<?php array(1, "a", false);`, expected)

	// Boolean literal
	testExpr(t, "<?php true;", ast.NewBooleanLiteralExpr(0, nil, true))
	testExpr(t, "<?php false;", ast.NewBooleanLiteralExpr(0, nil, false))

	// Null literal
	testExpr(t, "<?php null;", ast.NewNullLiteralExpr(0, nil))

	// Integer literal
	// decimal-literal
	testExpr(t, "<?php 42;", ast.NewIntegerLiteralExpr(0, nil, 42))
	// octal-literal
	testExpr(t, "<?php 042;", ast.NewIntegerLiteralExpr(0, nil, 34))
	// hexadecimal-literal
	testExpr(t, "<?php 0x42;", ast.NewIntegerLiteralExpr(0, nil, 66))
	// binary-literal
	testExpr(t, "<?php 0b110110101;", ast.NewIntegerLiteralExpr(0, nil, 437))

	// Floating literal
	testExpr(t, "<?php .5;", ast.NewFloatingLiteralExpr(0, nil, 0.5))
	testExpr(t, "<?php 1.2;", ast.NewFloatingLiteralExpr(0, nil, 1.2))
	testExpr(t, "<?php .5e-4;", ast.NewFloatingLiteralExpr(0, nil, 0.5e-4))
	testExpr(t, "<?php 2.5e+4;", ast.NewFloatingLiteralExpr(0, nil, 2.5e+4))
	testExpr(t, "<?php 2e4;", ast.NewFloatingLiteralExpr(0, nil, 2e4))

	// String literal
	testExpr(t, `<?php b'A "single quoted" \'string\'';`, ast.NewStringLiteralExpr(0, nil, `A "single quoted" 'string'`, ast.SingleQuotedString))
	testExpr(t, `<?php 'A "single quoted" \'string\'';`, ast.NewStringLiteralExpr(0, nil, `A "single quoted" 'string'`, ast.SingleQuotedString))
	testExpr(t, `<?php b"A \"double quoted\" 'string'";`, ast.NewStringLiteralExpr(0, nil, `A "double quoted" 'string'`, ast.DoubleQuotedString))
	testExpr(t, `<?php "A \"double quoted\" 'string'";`, ast.NewStringLiteralExpr(0, nil, `A "double quoted" 'string'`, ast.DoubleQuotedString))
}

func TestEchoStatement(t *testing.T) {
	testStmt(t, `<?php echo 12, "abc", $var;`, ast.NewEchoStmt(0, nil, []ast.IExpression{
		ast.NewIntegerLiteralExpr(0, nil, 12), ast.NewStringLiteralExpr(0, nil, "abc", ast.DoubleQuotedString),
		ast.NewSimpleVariableExpr(0, ast.NewVariableNameExpr(0, nil, "$var")),
	}))

	// Print
	testExpr(t, `<?php print "abc";`, ast.NewPrintExpr(0, nil, ast.NewStringLiteralExpr(0, nil, "abc", ast.DoubleQuotedString)))
}

func TestConditional(t *testing.T) {
	// Conditional
	testExpr(t, `<?php 1 ? "a" : "b";`, ast.NewConditionalExpr(0,
		ast.NewIntegerLiteralExpr(0, nil, 1), ast.NewStringLiteralExpr(0, nil, "a", ast.DoubleQuotedString),
		ast.NewStringLiteralExpr(0, nil, "b", ast.DoubleQuotedString),
	))

	testExpr(t, `<?php 1 ?: "b";`,
		ast.NewConditionalExpr(0, ast.NewIntegerLiteralExpr(0, nil, 1), nil, ast.NewStringLiteralExpr(0, nil, "b", ast.DoubleQuotedString)),
	)

	testExpr(t, `<?php 1 ? "a" : 2 ? "b": "c";`, ast.NewConditionalExpr(0,
		ast.NewIntegerLiteralExpr(0, nil, 1), ast.NewStringLiteralExpr(0, nil, "a", ast.DoubleQuotedString),
		ast.NewConditionalExpr(0,
			ast.NewIntegerLiteralExpr(0, nil, 2), ast.NewStringLiteralExpr(0, nil, "b", ast.DoubleQuotedString),
			ast.NewStringLiteralExpr(0, nil, "c", ast.DoubleQuotedString),
		),
	))

	// Coalesce
	testExpr(t, `<?php "a" ?? "b";`, ast.NewCoalesceExpr(0,
		ast.NewStringLiteralExpr(0, nil, "a", ast.DoubleQuotedString), ast.NewStringLiteralExpr(0, nil, "b", ast.DoubleQuotedString),
	))

	testExpr(t, `<?php "a" ?? "b" ?? "c";`, ast.NewCoalesceExpr(0,
		ast.NewStringLiteralExpr(0, nil, "a", ast.DoubleQuotedString),
		ast.NewCoalesceExpr(0,
			ast.NewStringLiteralExpr(0, nil, "b", ast.DoubleQuotedString), ast.NewStringLiteralExpr(0, nil, "c", ast.DoubleQuotedString),
		),
	))
}

func TestCastExpression(t *testing.T) {
	testExpr(t, `<?php (string)42;`, ast.NewCastExpr(0, nil, "string", ast.NewIntegerLiteralExpr(0, nil, 42)))
}

func TestParenthesizedExpression(t *testing.T) {
	testExpr(t, `<?php (1+2);`,
		ast.NewParenthesizedExpr(0, nil, ast.NewBinaryOpExpr(0, ast.NewIntegerLiteralExpr(0, nil, 1), "+", ast.NewIntegerLiteralExpr(0, nil, 2))),
	)
}

func TestConstDeclaration(t *testing.T) {
	testStmts(t, `<?php const PI = 3.141, ZERO = 0;`, []ast.IStatement{
		ast.NewConstDeclarationStmt(0, nil, "PI", ast.NewFloatingLiteralExpr(0, nil, 3.141)),
		ast.NewConstDeclarationStmt(0, nil, "ZERO", ast.NewIntegerLiteralExpr(0, nil, 0)),
	})
}

func TestEqualityExpression(t *testing.T) {
	testExpr(t, `<?php "234" !== true;`, ast.NewEqualityExpr(0,
		ast.NewStringLiteralExpr(0, nil, "234", ast.DoubleQuotedString), "!==", ast.NewBooleanLiteralExpr(0, nil, true),
	))
	testExpr(t, `<?php "234" == true;`, ast.NewEqualityExpr(0,
		ast.NewStringLiteralExpr(0, nil, "234", ast.DoubleQuotedString), "==", ast.NewBooleanLiteralExpr(0, nil, true),
	))
}

func TestOperatorExpression(t *testing.T) {
	// Shift
	testExpr(t, `<?php "234" << 12;`,
		ast.NewBinaryOpExpr(0, ast.NewStringLiteralExpr(0, nil, "234", ast.DoubleQuotedString), "<<", ast.NewIntegerLiteralExpr(0, nil, 12)),
	)

	// Additive
	testExpr(t, `<?php "234" + 12;`,
		ast.NewBinaryOpExpr(0, ast.NewStringLiteralExpr(0, nil, "234", ast.DoubleQuotedString), "+", ast.NewIntegerLiteralExpr(0, nil, 12)),
	)

	// Multiplicative
	testExpr(t, `<?php "234" * 12;`,
		ast.NewBinaryOpExpr(0, ast.NewStringLiteralExpr(0, nil, "234", ast.DoubleQuotedString), "*", ast.NewIntegerLiteralExpr(0, nil, 12)),
	)

	// Exponentiation
	testExpr(t, `<?php "234" ** 12;`,
		ast.NewBinaryOpExpr(0, ast.NewStringLiteralExpr(0, nil, "234", ast.DoubleQuotedString), "**", ast.NewIntegerLiteralExpr(0, nil, 12)),
	)

	// Logical not
	testExpr(t, `<?php !true;`, ast.NewLogicalNotExpr(0, nil, ast.NewBooleanLiteralExpr(0, nil, true)))

	// Logical inc or
	testExpr(t, `<?php "234" || 12;`,
		ast.NewLogicalExpr(0, ast.NewStringLiteralExpr(0, nil, "234", ast.DoubleQuotedString), "||", ast.NewIntegerLiteralExpr(0, nil, 12)),
	)
	// Logical inc or 2
	testExpr(t, `<?php "234" or 12;`,
		ast.NewLogicalExpr(0, ast.NewStringLiteralExpr(0, nil, "234", ast.DoubleQuotedString), "||", ast.NewIntegerLiteralExpr(0, nil, 12)),
	)

	// Logical exc or
	testExpr(t, `<?php "234" xor 12;`,
		ast.NewLogicalExpr(0, ast.NewStringLiteralExpr(0, nil, "234", ast.DoubleQuotedString), "xor", ast.NewIntegerLiteralExpr(0, nil, 12)),
	)

	// Logical and
	testExpr(t, `<?php "234" && 12;`,
		ast.NewLogicalExpr(0, ast.NewStringLiteralExpr(0, nil, "234", ast.DoubleQuotedString), "&&", ast.NewIntegerLiteralExpr(0, nil, 12)),
	)
	// Logical and 2
	testExpr(t, `<?php "234" and 12;`,
		ast.NewLogicalExpr(0, ast.NewStringLiteralExpr(0, nil, "234", ast.DoubleQuotedString), "&&", ast.NewIntegerLiteralExpr(0, nil, 12)),
	)

	// Bitwise inc or
	testExpr(t, `<?php "234" | 12;`,
		ast.NewLogicalExpr(0, ast.NewStringLiteralExpr(0, nil, "234", ast.DoubleQuotedString), "|", ast.NewIntegerLiteralExpr(0, nil, 12)),
	)

	// Bitwise exc or
	testExpr(t, `<?php "234" ^ 12;`,
		ast.NewLogicalExpr(0, ast.NewStringLiteralExpr(0, nil, "234", ast.DoubleQuotedString), "^", ast.NewIntegerLiteralExpr(0, nil, 12)),
	)

	// Bitwise and
	testExpr(t, `<?php "234" & 12;`,
		ast.NewLogicalExpr(0, ast.NewStringLiteralExpr(0, nil, "234", ast.DoubleQuotedString), "&", ast.NewIntegerLiteralExpr(0, nil, 12)),
	)
}

func TestIntrinsic(t *testing.T) {
	// Die
	testExpr(t, `<?php die(42);`, ast.NewExitIntrinsic(0, nil, ast.NewIntegerLiteralExpr(0, nil, 42)))
	testExpr(t, `<?php die();`, ast.NewExitIntrinsic(0, nil, nil))
	testExpr(t, `<?php die;`, ast.NewExitIntrinsic(0, nil, nil))
	// Exit
	testExpr(t, `<?php exit(42);`, ast.NewExitIntrinsic(0, nil, ast.NewIntegerLiteralExpr(0, nil, 42)))
	testExpr(t, `<?php exit();`, ast.NewExitIntrinsic(0, nil, nil))
	testExpr(t, `<?php exit;`, ast.NewExitIntrinsic(0, nil, nil))

	// Empty
	testExpr(t, `<?php empty(false);`, ast.NewEmptyIntrinsic(0, nil, ast.NewBooleanLiteralExpr(0, nil, false)))

	// Isset
	testExpr(t, `<?php isset($a);`,
		ast.NewIssetIntrinsic(0, nil, []ast.IExpression{ast.NewSimpleVariableExpr(0, ast.NewVariableNameExpr(0, nil, "$a"))}),
	)

	// Unset
	testExpr(t, `<?php unset($a);`,
		ast.NewUnsetIntrinsic(0, nil, []ast.IExpression{ast.NewSimpleVariableExpr(0, ast.NewVariableNameExpr(0, nil, "$a"))}),
	)
}
