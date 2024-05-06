package parser

import (
	"GoPHP/cmd/goPHP/ast"
	"GoPHP/cmd/goPHP/common"
	"GoPHP/cmd/goPHP/lexer"
	"fmt"
)

type Parser struct {
	lexer   *lexer.Lexer
	tokens  []*lexer.Token
	currPos int
}

func NewParser() *Parser {
	return &Parser{}
}

func (parser *Parser) init() {
	parser.lexer = lexer.NewLexer()
	parser.currPos = 0
}

func (parser *Parser) ProduceAST(sourceCode string) (*ast.Program, error) {
	parser.init()

	program := ast.NewProgram()

	var err error
	parser.tokens, err = parser.lexer.Tokenize(sourceCode)
	if err != nil {
		return program, err
	}

	for !parser.isEof() {
		if parser.at().TokenType == lexer.StartTagToken || parser.at().TokenType == lexer.EndTagToken {
			parser.eat()
			continue
		}
		stmt, err := parser.parseStatement()
		if err != nil {
			return program, err
		}
		program.Append(stmt)
	}

	return program, err
}

func (parser *Parser) parseStatement() (ast.IStatement, error) {
	// Spec: https://phplang.org/spec/11-statements.html#general

	// statement:
	//    compound-statement
	//    named-label-statement
	//    expression-statement
	//    selection-statement
	//    iteration-statement
	//    jump-statement
	//    try-statement
	//    declare-statement
	//    echo-statement
	//    unset-statement
	//    const-declaration
	//    function-definition
	//    class-declaration
	//    interface-declaration
	//    trait-declaration
	//    namespace-definition
	//    namespace-use-declaration
	//    global-declaration
	//    function-static-declaration

	if parser.at().TokenType == lexer.TextToken {
		return ast.NewExpressionStatement(ast.NewTextExpression(parser.eat().Value)), nil
	}

	// TODO compound-statement
	// TODO named-label-statement
	// TODO selection-statement
	// TODO iteration-statement
	// TODO jump-statement
	// TODO try-statement
	// TODO declare-statement

	// ------------------- MARK: echo-statement -------------------

	// Spec https://phplang.org/spec/11-statements.html#the-echo-statement

	// echo-statement:
	//    echo   expression-list   ;

	// expression-list:
	//    expression
	//    expression-list   ,   expression

	if parser.at().TokenType == lexer.KeywordToken && parser.at().Value == "echo" {
		parser.eat()

		expressions := make([]ast.IExpression, 0)
		for {
			expr, err := parser.parseExpression()
			if err != nil {
				return ast.NewEmptyStatement(), err
			}

			expressions = append(expressions, expr)

			if parser.at().TokenType == lexer.OperatorOrPunctuatorToken && parser.at().Value == "," {
				parser.eat()
				continue
			}
			if parser.at().TokenType == lexer.OperatorOrPunctuatorToken && parser.at().Value == ";" {
				parser.eat()
				break
			}
			return ast.NewEmptyStatement(), fmt.Errorf("Parser error: Invalid echo statement detected")
		}

		if len(expressions) == 0 {
			return ast.NewEmptyStatement(), fmt.Errorf("Parser error: Invalid echo statement detected")
		}

		return ast.NewEchoStatement(expressions), nil
	}

	// TODO unset-statement
	// TODO const-declaration
	// TODO function-definition
	// TODO class-declaration
	// TODO interface-declaration
	// TODO trait-declaration
	// TODO namespace-definition
	// TODO namespace-use-declaration
	// TODO global-declaration
	// TODO function-static-declaration

	// ------------------- MARK: expression-statement -------------------

	// Spec: https://phplang.org/spec/11-statements.html#grammar-expression-statement

	// expression-statement:
	//    expression(opt)   ;

	// If present, expression is evaluated for its side effects, if any, and any resulting value is discarded.
	// If expression is omitted, the statement is a null statement, which has no effect on execution.
	if parser.at().TokenType == lexer.OperatorOrPunctuatorToken && parser.at().Value == ";" {
		parser.eat()
	}

	if expr, err := parser.parseExpression(); err != nil {
		return ast.NewEmptyExpression(), err
	} else {
		if parser.at().TokenType == lexer.OperatorOrPunctuatorToken && parser.at().Value == ";" {
			parser.eat()
			return ast.NewExpressionStatement(expr), nil
		}
		return ast.NewEmptyExpression(), fmt.Errorf("Parser error: Statement must end with a semicolon")
	}
}

func (parser *Parser) parseExpression() (ast.IExpression, error) {
	// Spec: https://phplang.org/spec/10-expressions.html#general-1

	// primary-expression:
	//    variable
	//    class-constant-access-expression
	//    constant-access-expression
	//    literal
	//    array-creation-expression
	//    intrinsic
	//    anonymous-function-creation-expression
	//    object-creation-expression
	//    postfix-increment-expression
	//    postfix-decrement-expression
	//    prefix-increment-expression
	//    prefix-decrement-expression
	//    byref-assignment-expression
	//    shell-command-expression
	//    (   expression   )

	// ------------------- MARK: variable -------------------

	// Spec: https://phplang.org/spec/10-expressions.html#simple-variable

	// simple-variable:
	//    variable-name

	if parser.at().TokenType == lexer.VariableNameToken {
		return ast.NewSimpleVariableExpression(ast.NewVariableNameExpression(parser.eat().Value)), nil
	}

	// Spec: https://phplang.org/spec/10-expressions.html#simple-variable

	// simple-variable:
	//    variable-name
	//    $   {   expression   }

	if parser.at().TokenType == lexer.OperatorOrPunctuatorToken && parser.at().Value == "$" &&
		parser.next(0).TokenType == lexer.OperatorOrPunctuatorToken && parser.next(0).Value == "{" {
		parser.eatN(2)
		// Get expression
		expr, err := parser.parseExpression()
		if err != nil {
			return ast.NewEmptyExpression(), err
		}

		if parser.at().Value == "}" {
			parser.eat()
			return ast.NewSimpleVariableExpression(expr), nil
		}

		return ast.NewEmptyExpression(), fmt.Errorf("Parser error: End of simple variable expression not detected")
	}

	// Spec: https://phplang.org/spec/10-expressions.html#simple-variable

	// simple-variable:
	//    variable-name
	//    $   simple-variable

	if parser.at().TokenType == lexer.OperatorOrPunctuatorToken && parser.at().Value == "$" {
		parser.eat()
		if expr, err := parser.parseExpression(); err != nil {
			return ast.NewEmptyExpression(), err
		} else {
			return ast.NewSimpleVariableExpression(expr), nil
		}
	}

	// TODO class-constant-access-expression
	// TODO constant-access-expression

	// ------------------- MARK: literal -------------------

	// Spec: https://phplang.org/spec/10-expressions.html#grammar-literal

	// literal:
	//    integer-literal
	//    floating-literal
	//    string-literal

	// A literal evaluates to its value, as specified in the lexical specification for literals.

	// integer-literal
	if parser.at().TokenType == lexer.IntegerLiteralToken {
		// decimal-literal
		if common.IsDecimalLiteral(parser.at().Value) {
			return ast.NewIntegerLiteralExpression(common.DecimalLiteralToInt64(parser.eat().Value)), nil
		}

		// octal-literal
		if common.IsOctalLiteral(parser.at().Value) {
			return ast.NewIntegerLiteralExpression(common.OctalLiteralToInt64(parser.eat().Value)), nil
		}

		// hexadecimal-literal
		if common.IsHexadecimalLiteral(parser.at().Value) {
			return ast.NewIntegerLiteralExpression(common.HexadecimalLiteralToInt64(parser.eat().Value)), nil
		}

		// binary-literal
		if common.IsBinaryLiteral(parser.at().Value) {
			return ast.NewIntegerLiteralExpression(common.BinaryLiteralToInt64(parser.eat().Value)), nil
		}

		return ast.NewEmptyExpression(), fmt.Errorf("Parser error: Unsupported integer literal \"%s\"", parser.at().Value)
	}

	// floating-literal
	if parser.at().TokenType == lexer.FloatingLiteralToken {
		if common.IsFloatingLiteral(parser.at().Value) {
			return ast.NewFloatingLiteralExpression(common.FloatingLiteralToFloat64(parser.eat().Value)), nil
		}

		return ast.NewEmptyExpression(), fmt.Errorf("Parser error: Unsupported floating literal \"%s\"", parser.at().Value)
	}

	// string-literal
	if parser.at().TokenType == lexer.StringLiteralToken {
		// single-quoted-string-literal
		if common.IsSingleQuotedStringLiteral(parser.at().Value) {
			return ast.NewStringLiteralExpression(
					common.SingleQuotedStringLiteralToString(parser.eat().Value), ast.SingleQuotedString),
				nil
		}

		// TODO double-quoted-string-literal
		if common.IsDoubleQuotedStringLiteral(parser.at().Value) {
			return ast.NewStringLiteralExpression(
					common.DoubleQuotedStringLiteralToString(parser.eat().Value), ast.DoubleQuotedString),
				nil
		}

		// TODO heredoc-string-literal
		// TODO nowdoc-string-literal
	}

	// TODO array-creation-expression
	// TODO intrinsic
	// TODO anonymous-function-creation-expression
	// TODO object-creation-expression
	// TODO postfix-increment-expression
	// TODO postfix-decrement-expression
	// TODO prefix-increment-expression
	// TODO prefix-decrement-expression
	// TODO byref-assignment-expression
	// TODO shell-command-expression
	// TODO (   expression   )

	return ast.NewEmptyExpression(), fmt.Errorf("Parser error: Unsupported expression type: %s", parser.at())
}
