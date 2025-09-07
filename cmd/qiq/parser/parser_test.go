package parser

import (
	"QIQ/cmd/qiq/ast"
	"QIQ/cmd/qiq/common/os"
	"QIQ/cmd/qiq/ini"
	"QIQ/cmd/qiq/phpError"
	"testing"
)

var TEST_FILE_NAME string = getTestFileName()
var TEST_FILE_PATH string = getTestFilePath()

func getTestFileName() string {
	if os.IS_WIN {
		return `C:\Users\admin\test.php`
	} else {
		return "/home/admin/test.php"
	}
}

func getTestFilePath() string {
	if os.IS_WIN {
		return `C:\Users\admin`
	} else {
		return "/home/admin"
	}
}

func testForError(t *testing.T, php string, expected phpError.Error) {
	_, err := NewParser(ini.NewDevIni()).ProduceAST(php, TEST_FILE_NAME)
	if err == nil {
		t.Errorf("\nCode: \"%s\"\nExpected error, but got \"nil\"", php)
		return
	}
	if err.GetErrorType() != expected.GetErrorType() || err.GetMessage() != expected.GetMessage() {
		t.Errorf("\nCode: \"%s\"\nExpected: %s\nGot:      %s", php, expected, err)
	}
}

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
			t.Errorf("\nExpected: \"%s\"\nGot:      \"%s\"", ast.ToString(expect), ast.ToString(actual))
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
		t.Errorf("Code: \"%s\"\nUnexpected error: \"%s\"\n", php, err)
		return
	}
	for index, expect := range expected {
		actual := program.GetStatements()[index]
		if ast.ToString(expect) != ast.ToString(actual) {
			t.Errorf(
				"Code: \"%s\"\nExpected: \"%s\"\nGot:      \"%s\"\n",
				php, ast.ToString(expect), ast.ToString(actual),
			)
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
	expected.AddElement(nil, ast.NewIntegerLiteralExpr(0, nil, 1))
	expected.AddElement(nil, ast.NewStringLiteralExpr(0, nil, "a", ast.DoubleQuotedString))
	expected.AddElement(nil, ast.NewBooleanLiteralExpr(0, nil, false))
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
	testExpr(t, "<?php b<<<   ID\nSome text\nover\nmutiple lines\nID;", ast.NewStringLiteralExpr(0, nil, "Some text\nover\nmutiple lines", ast.HeredocString))
	testExpr(t, "<?php <<<EOF\nEOF;", ast.NewStringLiteralExpr(0, nil, "", ast.HeredocString))
	testExpr(t, "<?php <<<   ID\nSome text\nover\nmutiple lines\nID;", ast.NewStringLiteralExpr(0, nil, "Some text\nover\nmutiple lines", ast.HeredocString))
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

func TestErrorControlExpression(t *testing.T) {
	testExpr(t, `<?php @func();`,
		ast.NewErrorControlExpr(0, nil, ast.NewFunctionCallExpr(0, nil, ast.NewStringLiteralExpr(0, nil, "func", ast.SingleQuotedString), []ast.IExpression{})),
	)
}

func TestGlobalDeclaration(t *testing.T) {
	testStmt(t, `<?php global $foo, $bar;`,
		ast.NewGlobalDeclarationStmt(0, nil, []ast.IExpression{
			ast.NewSimpleVariableExpr(0, ast.NewVariableNameExpr(0, nil, "$foo")),
			ast.NewSimpleVariableExpr(0, ast.NewVariableNameExpr(0, nil, "$bar")),
		}),
	)
}

// -------------------------------------- Loops -------------------------------------- MARK: Loops

func TestLoops(t *testing.T) {
	// While
	testStmt(t, `<?php while (true) {}`,
		ast.NewWhileStmt(0, nil, ast.NewBooleanLiteralExpr(0, nil, true), ast.NewCompoundStmt(0, []ast.IStatement{})),
	)
	testStmt(t, `<?php while (true) { func(); }`,
		ast.NewWhileStmt(0, nil, ast.NewBooleanLiteralExpr(0, nil, true), ast.NewCompoundStmt(0, []ast.IStatement{
			ast.NewExpressionStmt(0, ast.NewFunctionCallExpr(0, nil, ast.NewStringLiteralExpr(0, nil, "func", ast.SingleQuotedString), []ast.IExpression{})),
		})),
	)
	testStmt(t, `<?php while (true): endwhile;`,
		ast.NewWhileStmt(0, nil, ast.NewBooleanLiteralExpr(0, nil, true), ast.NewCompoundStmt(0, []ast.IStatement{})),
	)
	testStmt(t, `<?php while (true): func(); endwhile;`,
		ast.NewWhileStmt(0, nil, ast.NewBooleanLiteralExpr(0, nil, true), ast.NewCompoundStmt(0, []ast.IStatement{
			ast.NewExpressionStmt(0, ast.NewFunctionCallExpr(0, nil, ast.NewStringLiteralExpr(0, nil, "func", ast.SingleQuotedString), []ast.IExpression{})),
		})),
	)

	// Do while
	testStmt(t, `<?php do {} while (true);`,
		ast.NewDoStmt(0, nil, ast.NewBooleanLiteralExpr(0, nil, true), ast.NewCompoundStmt(0, []ast.IStatement{})),
	)
	testStmt(t, `<?php do { func(); } while (true);`,
		ast.NewDoStmt(0, nil, ast.NewBooleanLiteralExpr(0, nil, true), ast.NewCompoundStmt(0, []ast.IStatement{
			ast.NewExpressionStmt(0, ast.NewFunctionCallExpr(0, nil, ast.NewStringLiteralExpr(0, nil, "func", ast.SingleQuotedString), []ast.IExpression{})),
		})),
	)

	// For
	testStmt(t, `<?php for ($i = 1; $i <= 10; ++$i) { func(); }`,
		ast.NewForStmt(0, nil,
			ast.NewCompoundStmt(0, []ast.IStatement{ast.NewSimpleAssignmentExpr(0, ast.NewSimpleVariableExpr(0, ast.NewVariableNameExpr(0, nil, "$i")), ast.NewIntegerLiteralExpr(0, nil, 1))}),
			ast.NewCompoundStmt(0, []ast.IStatement{ast.NewRelationalExpr(0, ast.NewSimpleVariableExpr(0, ast.NewVariableNameExpr(0, nil, "$i")), "<=", ast.NewIntegerLiteralExpr(0, nil, 10))}),
			ast.NewCompoundStmt(0, []ast.IStatement{ast.NewPrefixIncExpr(0, nil, ast.NewSimpleVariableExpr(0, ast.NewVariableNameExpr(0, nil, "$i")), "++")}),
			ast.NewCompoundStmt(0, []ast.IStatement{ast.NewExpressionStmt(0, ast.NewFunctionCallExpr(0, nil, ast.NewStringLiteralExpr(0, nil, "func", ast.SingleQuotedString), []ast.IExpression{}))}),
		),
	)
	// Omit 1st and 3rd expressions
	testStmt(t, `<?php for (; $i <= 10;): func(); endfor;`,
		ast.NewForStmt(0, nil,
			nil,
			ast.NewCompoundStmt(0, []ast.IStatement{ast.NewRelationalExpr(0, ast.NewSimpleVariableExpr(0, ast.NewVariableNameExpr(0, nil, "$i")), "<=", ast.NewIntegerLiteralExpr(0, nil, 10))}),
			nil,
			ast.NewCompoundStmt(0, []ast.IStatement{ast.NewExpressionStmt(0, ast.NewFunctionCallExpr(0, nil, ast.NewStringLiteralExpr(0, nil, "func", ast.SingleQuotedString), []ast.IExpression{}))}),
		),
	)
	// Omit all 3 expressions
	testStmt(t, `<?php for (;;) { func(); }`,
		ast.NewForStmt(0, nil,
			nil,
			nil,
			nil,
			ast.NewCompoundStmt(0, []ast.IStatement{ast.NewExpressionStmt(0, ast.NewFunctionCallExpr(0, nil, ast.NewStringLiteralExpr(0, nil, "func", ast.SingleQuotedString), []ast.IExpression{}))}),
		),
	)
	// Use groups of expressions
	testStmt(t, `<?php for ($a = 100, $i = 1; ++$i, $i <= 10; ++$i, $a -= 10) { func(); }`,
		ast.NewForStmt(0, nil,
			ast.NewCompoundStmt(0, []ast.IStatement{
				ast.NewSimpleAssignmentExpr(0, ast.NewSimpleVariableExpr(0, ast.NewVariableNameExpr(0, nil, "$a")), ast.NewIntegerLiteralExpr(0, nil, 100)),
				ast.NewSimpleAssignmentExpr(0, ast.NewSimpleVariableExpr(0, ast.NewVariableNameExpr(0, nil, "$i")), ast.NewIntegerLiteralExpr(0, nil, 1)),
			}),
			ast.NewCompoundStmt(0, []ast.IStatement{
				ast.NewPrefixIncExpr(0, nil, ast.NewSimpleVariableExpr(0, ast.NewVariableNameExpr(0, nil, "$i")), "++"),
				ast.NewRelationalExpr(0, ast.NewSimpleVariableExpr(0, ast.NewVariableNameExpr(0, nil, "$i")), "<=", ast.NewIntegerLiteralExpr(0, nil, 10)),
			}),
			ast.NewCompoundStmt(0, []ast.IStatement{
				ast.NewPrefixIncExpr(0, nil, ast.NewSimpleVariableExpr(0, ast.NewVariableNameExpr(0, nil, "$i")), "++"),
				ast.NewCompoundAssignmentExpr(0, ast.NewSimpleVariableExpr(0, ast.NewVariableNameExpr(0, nil, "$a")), "-", ast.NewIntegerLiteralExpr(0, nil, 10)),
			}),
			ast.NewCompoundStmt(0, []ast.IStatement{ast.NewExpressionStmt(0, ast.NewFunctionCallExpr(0, nil, ast.NewStringLiteralExpr(0, nil, "func", ast.SingleQuotedString), []ast.IExpression{}))}),
		),
	)

	// Foreach
	testStmt(t, `<?php foreach ([] as $value) {}`,
		ast.NewForeachStmt(0, nil,
			ast.NewArrayLiteralExpr(0, nil),
			nil,
			ast.NewSimpleVariableExpr(0, ast.NewVariableNameExpr(0, nil, "$value")),
			false,
			ast.NewCompoundStmt(0, []ast.IStatement{}),
		),
	)
	testStmt(t, `<?php foreach ([] as $key => $value) { func(); }`,
		ast.NewForeachStmt(0, nil,
			ast.NewArrayLiteralExpr(0, nil),
			ast.NewSimpleVariableExpr(0, ast.NewVariableNameExpr(0, nil, "$key")),
			ast.NewSimpleVariableExpr(0, ast.NewVariableNameExpr(0, nil, "$value")),
			false,
			ast.NewCompoundStmt(0, []ast.IStatement{ast.NewExpressionStmt(0, ast.NewFunctionCallExpr(0, nil, ast.NewStringLiteralExpr(0, nil, "func", ast.SingleQuotedString), []ast.IExpression{}))}),
		),
	)
	// byRef
	testStmt(t, `<?php foreach ([] as &$value) {}`,
		ast.NewForeachStmt(0, nil,
			ast.NewArrayLiteralExpr(0, nil),
			nil,
			ast.NewSimpleVariableExpr(0, ast.NewVariableNameExpr(0, nil, "$value")),
			true,
			ast.NewCompoundStmt(0, []ast.IStatement{}),
		),
	)
	testStmt(t, `<?php foreach ([] as $key => &$value) { func(); }`,
		ast.NewForeachStmt(0, nil,
			ast.NewArrayLiteralExpr(0, nil),
			ast.NewSimpleVariableExpr(0, ast.NewVariableNameExpr(0, nil, "$key")),
			ast.NewSimpleVariableExpr(0, ast.NewVariableNameExpr(0, nil, "$value")),
			true,
			ast.NewCompoundStmt(0, []ast.IStatement{ast.NewExpressionStmt(0, ast.NewFunctionCallExpr(0, nil, ast.NewStringLiteralExpr(0, nil, "func", ast.SingleQuotedString), []ast.IExpression{}))}),
		),
	)
	testForError(t, `<?php foreach ([] as &$key => $value) { }`, phpError.NewParseError("Syntax error, key cannot be by reference at /home/admin/test.php:1:22"))
}

// -------------------------------------- Class -------------------------------------- MARK: Class

func TestClassDeclaration(t *testing.T) {
	//  Object creation
	stmt := ast.NewExpressionStmt(0, ast.NewObjectCreationExpr(0, nil, "stdClass", []ast.IExpression{}))
	testStmt(t, `<?php new stdClass;`, stmt)
	//  Object creation with namespace
	stmt = ast.NewExpressionStmt(0, ast.NewObjectCreationExpr(0, nil, `\stdClass`, []ast.IExpression{}))
	testStmt(t, `<?php new \stdClass;`, stmt)
	//  Object creation with namespaces
	stmt = ast.NewExpressionStmt(0, ast.NewObjectCreationExpr(0, nil, `my\name\space\stdClass`, []ast.IExpression{}))
	testStmt(t, `<?php new my\name\space\stdClass;`, stmt)

	// Reserved name
	testForError(t, `<?php class Parent {}`, phpError.NewError(`Cannot use "Parent" as a class name as it is reserved in %s:1:13`, TEST_FILE_NAME))

	// Empty class
	class := ast.NewClassDeclarationStmt(0, nil, "c", false, false)
	testStmt(t, `<?php class c { }`, class)

	// Abstract class
	class = ast.NewClassDeclarationStmt(0, nil, "c", true, false)
	testStmt(t, `<?php abstract class c { }`, class)

	// Final class
	class = ast.NewClassDeclarationStmt(0, nil, "c", false, true)
	testStmt(t, `<?php final class c { }`, class)

	// Abstract and extended class
	class = ast.NewClassDeclarationStmt(0, nil, "c", true, false)
	class.BaseClass = "b"
	testStmt(t, `<?php abstract class c extends b { }`, class)

	// Final class with interfaces
	class = ast.NewClassDeclarationStmt(0, nil, "c", false, true)
	class.Interfaces = append(class.Interfaces, "i")
	testStmt(t, `<?php final class c implements i { }`, class)

	// Final class with multiple interfaces
	class = ast.NewClassDeclarationStmt(0, nil, "c", false, true)
	class.Interfaces = append(class.Interfaces, "i")
	class.Interfaces = append(class.Interfaces, "j")
	testStmt(t, `<?php final class c implements i, j { }`, class)

	// Class with constants
	class = ast.NewClassDeclarationStmt(0, nil, "c", false, true)
	class.AddConst(ast.NewClassConstDeclarationStmt(0, nil, "a", ast.NewStringLiteralExpr(0, nil, "a", ast.DoubleQuotedString), "public"))
	class.AddConst(ast.NewClassConstDeclarationStmt(0, nil, "b", ast.NewStringLiteralExpr(0, nil, "b", ast.DoubleQuotedString), "private"))
	class.AddConst(ast.NewClassConstDeclarationStmt(0, nil, "c", ast.NewIntegerLiteralExpr(0, nil, 3), "private"))
	testStmt(t, `<?php final class c { const a="a"; private const b="b", c=3; }`, class)

	// Class with trait
	class = ast.NewClassDeclarationStmt(0, nil, "c", false, false)
	class.AddTrait(ast.NewTraitUseStmt(0, nil, "MyTrait"))
	testStmt(t, `<?php class c { use MyTrait; }`, class)

	// Class with multiple traits
	class = ast.NewClassDeclarationStmt(0, nil, "c", false, false)
	class.AddTrait(ast.NewTraitUseStmt(0, nil, "MyTrait"))
	class.AddTrait(ast.NewTraitUseStmt(0, nil, "MySecondTrait"))
	testStmt(t, `<?php class c { use MyTrait, MySecondTrait; }`, class)

	// Class with constructor
	class = ast.NewClassDeclarationStmt(0, nil, "c", false, false)
	class.AddMethod(ast.NewMethodDefinitionStmt(0, nil, "__construct", []string{}, []ast.FunctionParameter{}, ast.NewCompoundStmt(0, []ast.IStatement{}), []string{}))
	testStmt(t, `<?php class c { function __construct() {} }`, class)

	// Class with private constructor
	class = ast.NewClassDeclarationStmt(0, nil, "c", false, false)
	class.AddMethod(ast.NewMethodDefinitionStmt(0, nil, "__construct", []string{"private"}, []ast.FunctionParameter{}, ast.NewCompoundStmt(0, []ast.IStatement{}), []string{}))
	testStmt(t, `<?php class c { private function __construct() {} }`, class)

	// Class with final private constructor
	class = ast.NewClassDeclarationStmt(0, nil, "c", false, false)
	class.AddMethod(ast.NewMethodDefinitionStmt(0, nil, "__construct", []string{"private", "final"}, []ast.FunctionParameter{}, ast.NewCompoundStmt(0, []ast.IStatement{}), []string{}))
	testStmt(t, `<?php class c { final private function __construct() {} }`, class)

	// Class with constructor with parameters
	class = ast.NewClassDeclarationStmt(0, nil, "c", false, false)
	class.AddMethod(ast.NewMethodDefinitionStmt(0, nil, "__construct", []string{}, []ast.FunctionParameter{{Name: "$name", Type: []string{"string"}}}, ast.NewCompoundStmt(0, []ast.IStatement{}), []string{}))
	testStmt(t, `<?php class c { function __construct(string $name) {} }`, class)

	// Class with constructor with body
	class = ast.NewClassDeclarationStmt(0, nil, "c", false, false)
	class.AddMethod(ast.NewMethodDefinitionStmt(0, nil, "__construct", []string{}, []ast.FunctionParameter{}, ast.NewCompoundStmt(0, []ast.IStatement{ast.NewExpressionStmt(0, ast.NewExitIntrinsic(0, nil, nil))}), []string{}))
	testStmt(t, `<?php class c { function __construct() { exit(); } }`, class)

	// Class with destructor
	class = ast.NewClassDeclarationStmt(0, nil, "c", false, false)
	class.AddMethod(ast.NewMethodDefinitionStmt(0, nil, "__destruct", []string{}, []ast.FunctionParameter{}, ast.NewCompoundStmt(0, []ast.IStatement{}), []string{}))
	testStmt(t, `<?php class c { function __destruct() {} }`, class)

	// Class with private destructor
	class = ast.NewClassDeclarationStmt(0, nil, "c", false, false)
	class.AddMethod(ast.NewMethodDefinitionStmt(0, nil, "__destruct", []string{"private"}, []ast.FunctionParameter{}, ast.NewCompoundStmt(0, []ast.IStatement{}), []string{}))
	testStmt(t, `<?php class c { private function __destruct() {} }`, class)

	// Class with final private destructor
	class = ast.NewClassDeclarationStmt(0, nil, "c", false, false)
	class.AddMethod(ast.NewMethodDefinitionStmt(0, nil, "__destruct", []string{"private", "final"}, []ast.FunctionParameter{}, ast.NewCompoundStmt(0, []ast.IStatement{}), []string{}))
	testStmt(t, `<?php class c { final private function __destruct() {} }`, class)

	// Class with destructor with body
	class = ast.NewClassDeclarationStmt(0, nil, "c", false, false)
	class.AddMethod(ast.NewMethodDefinitionStmt(0, nil, "__destruct", []string{}, []ast.FunctionParameter{}, ast.NewCompoundStmt(0, []ast.IStatement{ast.NewExpressionStmt(0, ast.NewExitIntrinsic(0, nil, nil))}), []string{}))
	testStmt(t, `<?php class c { function __destruct() { exit(); } }`, class)

	// Class with private static method
	class = ast.NewClassDeclarationStmt(0, nil, "c", false, false)
	class.AddMethod(ast.NewMethodDefinitionStmt(0, nil, "myFunction", []string{"private", "static"}, []ast.FunctionParameter{}, ast.NewCompoundStmt(0, []ast.IStatement{}), []string{}))
	testStmt(t, `<?php class c { private static function myFunction() {} }`, class)

	// Class with method with parameters
	class = ast.NewClassDeclarationStmt(0, nil, "c", false, false)
	class.AddMethod(ast.NewMethodDefinitionStmt(0, nil, "myFunction", []string{}, []ast.FunctionParameter{ast.NewFunctionParam(false, "$name", []string{"string"})}, ast.NewCompoundStmt(0, []ast.IStatement{}), []string{"void"}))
	testStmt(t, `<?php class c { function myFunction(string $name): void {} }`, class)

	// Class with method with return type
	class = ast.NewClassDeclarationStmt(0, nil, "c", false, false)
	class.AddMethod(ast.NewMethodDefinitionStmt(0, nil, "myFunction", []string{}, []ast.FunctionParameter{}, ast.NewCompoundStmt(0, []ast.IStatement{ast.NewExpressionStmt(0, ast.NewExitIntrinsic(0, nil, nil))}), []string{"null", "int"}))
	testStmt(t, `<?php class c { function myFunction(): ?int { exit(); } }`, class)

	// Class with method with body
	class = ast.NewClassDeclarationStmt(0, nil, "c", false, false)
	class.AddMethod(ast.NewMethodDefinitionStmt(0, nil, "myFunction", []string{}, []ast.FunctionParameter{}, ast.NewCompoundStmt(0, []ast.IStatement{ast.NewExpressionStmt(0, ast.NewExitIntrinsic(0, nil, nil))}), []string{"int", "float"}))
	testStmt(t, `<?php class c { function myFunction(): int|float { exit(); } }`, class)

	// Class with property
	class = ast.NewClassDeclarationStmt(0, nil, "c", false, false)
	class.AddProperty(ast.NewPropertyDeclarationStmt(0, nil, "$member", "public", false, []string{}, nil))
	testStmt(t, `<?php class c { public $member; }`, class)

	// Class with protected property and initial value
	class = ast.NewClassDeclarationStmt(0, nil, "c", false, false)
	class.AddProperty(ast.NewPropertyDeclarationStmt(0, nil, "$member", "protected", false, []string{}, ast.NewIntegerLiteralExpr(0, nil, 42)))
	testStmt(t, `<?php class c { protected $member = 42; }`, class)

	// Class with multiple properties
	class = ast.NewClassDeclarationStmt(0, nil, "c", false, false)
	class.AddProperty(ast.NewPropertyDeclarationStmt(0, nil, "$a", "private", false, []string{}, nil))
	class.AddProperty(ast.NewPropertyDeclarationStmt(0, nil, "$b", "protected", false, []string{}, nil))
	class.AddProperty(ast.NewPropertyDeclarationStmt(0, nil, "$c", "public", false, []string{"null", "int"}, ast.NewIntegerLiteralExpr(0, nil, 42)))
	testStmt(t, `<?php class c { private $a; protected $b; public ?int $c = 42; }`, class)
}

// -------------------------------------- Interface -------------------------------------- MARK: Interface

func TestInterfaceDeclaration(t *testing.T) {
	// Reserved name
	testForError(t, `<?php interface Parent {}`, phpError.NewError(`Cannot use "Parent" as an interface name as it is reserved in %s:1:17`, TEST_FILE_NAME))

	// Empty interface
	interfaceDecl := ast.NewInterfaceDeclarationStmt(0, nil, "i")
	testStmt(t, `<?php interface i { }`, interfaceDecl)

	// Interface with parent
	interfaceDecl = ast.NewInterfaceDeclarationStmt(0, nil, "i")
	interfaceDecl.Parents = append(interfaceDecl.Parents, "j")
	testStmt(t, `<?php interface i extends j { }`, interfaceDecl)

	// Interface with multiple parents
	interfaceDecl = ast.NewInterfaceDeclarationStmt(0, nil, "i")
	interfaceDecl.Parents = append(interfaceDecl.Parents, "j")
	interfaceDecl.Parents = append(interfaceDecl.Parents, "k")
	testStmt(t, `<?php interface i extends j, k { }`, interfaceDecl)

	// Interface with methods
	interfaceDecl = ast.NewInterfaceDeclarationStmt(0, nil, "i")
	interfaceDecl.AddMethod(ast.NewMethodDefinitionStmt(0, nil, "f", []string{}, []ast.FunctionParameter{{Name: "$p", Type: []string{"string"}}}, nil, []string{"void"}))
	testStmt(t, `<?php interface i { function f (string $p): void; }`, interfaceDecl)

	// Interface with constants
	interfaceDecl = ast.NewInterfaceDeclarationStmt(0, nil, "i")
	interfaceDecl.AddConst(ast.NewClassConstDeclarationStmt(0, nil, "a", ast.NewStringLiteralExpr(0, nil, "a", ast.DoubleQuotedString), "public"))
	interfaceDecl.AddConst(ast.NewClassConstDeclarationStmt(0, nil, "b", ast.NewStringLiteralExpr(0, nil, "b", ast.DoubleQuotedString), "private"))
	interfaceDecl.AddConst(ast.NewClassConstDeclarationStmt(0, nil, "c", ast.NewIntegerLiteralExpr(0, nil, 3), "private"))
	testStmt(t, `<?php interface i { const a="a"; private const b="b", c=3; }`, interfaceDecl)
}

// -------------------------------------- Anonymous functions -------------------------------------- MARK: Anonymous functions

func TestAnonymousFunctions(t *testing.T) {
	// Empty anonymous function
	stmt := ast.NewExpressionStmt(0, ast.NewSimpleAssignmentExpr(0,
		ast.NewSimpleVariableExpr(0, ast.NewVariableNameExpr(0, nil, "$f")),
		ast.NewAnonymousFunctionCreationExpr(0, nil, []ast.FunctionParameter{}, ast.NewCompoundStmt(0, []ast.IStatement{}), []string{}),
	))
	testStmt(t, `<?php $f = function() {};`, stmt)

	// Anonymous function
	stmt = ast.NewExpressionStmt(0, ast.NewSimpleAssignmentExpr(0,
		ast.NewSimpleVariableExpr(0, ast.NewVariableNameExpr(0, nil, "$f")),
		ast.NewAnonymousFunctionCreationExpr(0, nil, []ast.FunctionParameter{}, ast.NewCompoundStmt(0, []ast.IStatement{
			ast.NewExpressionStmt(0, ast.NewFunctionCallExpr(0, nil, ast.NewStringLiteralExpr(0, nil, "do_smth", ast.SingleQuotedString), []ast.IExpression{})),
		}), []string{}),
	))
	testStmt(t, `<?php $f = function() { do_smth(); };`, stmt)

	// Anonymous function with byRef param
	stmt = ast.NewExpressionStmt(0, ast.NewSimpleAssignmentExpr(0,
		ast.NewSimpleVariableExpr(0, ast.NewVariableNameExpr(0, nil, "$f")),
		ast.NewAnonymousFunctionCreationExpr(0, nil, []ast.FunctionParameter{ast.NewFunctionParam(true, "$a", []string{})}, ast.NewCompoundStmt(0, []ast.IStatement{}), []string{}),
	))
	testStmt(t, `<?php $f = function(&$a) {};`, stmt)
}
