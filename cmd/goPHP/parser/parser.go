package parser

import (
	"GoPHP/cmd/goPHP/ast"
	"GoPHP/cmd/goPHP/common"
	"GoPHP/cmd/goPHP/lexer"
	"fmt"
	"strings"
)

type Parser struct {
	program *ast.Program
	lexer   *lexer.Lexer
	tokens  []*lexer.Token
	currPos int
}

func NewParser() *Parser {
	return &Parser{}
}

func (parser *Parser) init() {
	parser.program = ast.NewProgram()
	parser.lexer = lexer.NewLexer()
	parser.currPos = 0
}

func (parser *Parser) ProduceAST(sourceCode string) (*ast.Program, error) {
	parser.init()

	var err error
	parser.tokens, err = parser.lexer.Tokenize(sourceCode)
	if err != nil {
		return parser.program, err
	}

	for !parser.isEof() {
		if parser.isTokenType(lexer.StartTagToken, true) || parser.isTokenType(lexer.EndTagToken, true) {
			continue
		}
		stmt, err := parser.parseStatement()
		if err != nil {
			return parser.program, err
		}
		parser.program.Append(stmt)
	}

	return parser.program, err
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

	if parser.isTokenType(lexer.TextToken, false) {
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

	if parser.isToken(lexer.KeywordToken, "echo", true) {
		expressions := make([]ast.IExpression, 0)
		for {
			expr, err := parser.parseExpression()
			if err != nil {
				return ast.NewEmptyStatement(), err
			}

			expressions = append(expressions, expr)

			if parser.isToken(lexer.OperatorOrPunctuatorToken, ",", true) {
				continue
			}
			if parser.isToken(lexer.OperatorOrPunctuatorToken, ";", true) {
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

	// ------------------- MARK: const-declaration -------------------

	// Spec: https://phplang.org/spec/14-classes.html#grammar-const-declaration

	// const-declaration:
	//    const   const-elements   ;

	// const-elements:
	//    const-element
	//    const-elements   ,   const-element

	// const-element:
	//    name   =   constant-expression

	if parser.isToken(lexer.KeywordToken, "const", true) {
		if err := parser.expectTokenType(lexer.NameToken, false); err != nil {
			return ast.NewEmptyStatement(), err
		}
		for {
			name := parser.eat().Value
			if err := parser.expect(lexer.OperatorOrPunctuatorToken, "=", true); err != nil {
				return ast.NewEmptyStatement(), err
			}
			// TODO parse constant-expression
			value, err := parser.parseExpression()
			if err != nil {
				return ast.NewEmptyStatement(), err
			}

			stmt := ast.NewConstDeclarationStatement(name, value)
			if parser.isToken(lexer.OperatorOrPunctuatorToken, ",", true) {
				parser.program.Append(stmt)
				continue
			}
			if parser.isToken(lexer.OperatorOrPunctuatorToken, ";", true) {
				return stmt, nil
			}
			return ast.NewEmptyStatement(), fmt.Errorf("Parser error: Const declaration - unexpected token %s", parser.at())
		}
	}

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
	parser.isToken(lexer.OperatorOrPunctuatorToken, ";", true)

	if expr, err := parser.parseExpression(); err != nil {
		return ast.NewEmptyExpression(), err
	} else {
		if parser.isToken(lexer.OperatorOrPunctuatorToken, ";", true) {
			return ast.NewExpressionStatement(expr), nil
		}
		return ast.NewEmptyExpression(), fmt.Errorf("Parser error: Statement must end with a semicolon")
	}
}

func (parser *Parser) parseExpression() (ast.IExpression, error) {
	// Spec: https://phplang.org/spec/10-expressions.html#grammar-expression

	// expression:
	//    logical-inc-OR-expression-2
	//    include-expression
	//    include-once-expression
	//    require-expression
	//    require-once-expression
	// Spec-Fix: So that by following assignment-expression the primary-expression is reachable
	//    assignment-expression

	// TODO logical-inc-OR-expression-2
	// TODO include-expression
	// TODO include-once-expression
	// TODO require-expression
	// TODO require-once-expression

	// assignment-expression
	return parser.parseAssignmentExpression()
}

func (parser *Parser) parseAssignmentExpression() (ast.IExpression, error) {
	// Spec: https://phplang.org/spec/10-expressions.html#grammar-assignment-expression

	// assignment-expression:
	//    conditional-expression
	//    simple-assignment-expression
	//    compound-assignment-expression

	// conditional-expression

	// Spec: https://phplang.org/spec/10-expressions.html#grammar-conditional-expression

	// conditional-expression:
	//    coalesce-expression
	//    conditional-expression   ?   expression(opt)   :   coalesce-expression

	// coalesce-expression
	expr, err := parser.parseCoalesceExpression()
	if err != nil {
		return ast.NewEmptyExpression(), err
	}

	// conditional-expression   ?   expression(opt)   :   coalesce-expression
	for parser.isToken(lexer.OperatorOrPunctuatorToken, "?", true) {
		var ifExpr ast.IExpression = nil
		if !parser.isToken(lexer.OperatorOrPunctuatorToken, ":", false) {
			ifExpr, err = parser.parseExpression()
			if err != nil {
				return ast.NewEmptyExpression(), err
			}
		}
		if err := parser.expect(lexer.OperatorOrPunctuatorToken, ":", true); err != nil {
			return ast.NewEmptyExpression(), err
		}
		elseExpr, err := parser.parseExpression()
		if err != nil {
			return ast.NewEmptyExpression(), err
		}
		expr = ast.NewConditionalExpression(expr, ifExpr, elseExpr)
	}

	// Spec: https://phplang.org/spec/10-expressions.html#grammar-simple-assignment-expression

	// simple-assignment-expression:
	//    variable   =   assignment-expression
	//    list-intrinsic   =   assignment-expression

	// TODO simple-assignment-expression - list-intrinsic
	if ast.IsVariableExpression(expr) && parser.isToken(lexer.OperatorOrPunctuatorToken, "=", true) {
		value, err := parser.parseAssignmentExpression()
		if err != nil {
			return ast.NewEmptyExpression(), err
		}
		return ast.NewSimpleAssignmentExpression(expr, value), nil
	}

	// Spec: https://phplang.org/spec/10-expressions.html#grammar-compound-assignment-expression

	// compound-assignment-expression:
	//    variable   compound-assignment-operator   assignment-expression

	// compound-assignment-operator: one of
	//    **=   *=   /=   %=   +=   -=   .=   <<=   >>=   &=   ^=   |=

	if ast.IsVariableExpression(expr) &&
		parser.isTokenType(lexer.OperatorOrPunctuatorToken, false) && common.IsCompoundAssignmentOperator(parser.at().Value) {
		operatorStr := strings.ReplaceAll(parser.eat().Value, "=", "")
		value, err := parser.parseAssignmentExpression()
		if err != nil {
			return ast.NewEmptyExpression(), nil
		}
		return ast.NewCompoundAssignmentExpression(expr, operatorStr, value), nil
	}

	return expr, nil
}

func (parser *Parser) parseCoalesceExpression() (ast.IExpression, error) {
	// Spec: https://phplang.org/spec/10-expressions.html#grammar-coalesce-expression

	// coalesce-expression:
	//    logical-inc-OR-expression-1
	//    logical-inc-OR-expression-1   ??   coalesce-expression

	// logical-inc-OR-expression-1
	expr, err := parser.parseLogicalIncOrExpression1()
	if err != nil {
		return ast.NewEmptyExpression(), err
	}

	// logical-inc-OR-expression-1   ??   coalesce-expression
	if parser.isToken(lexer.OperatorOrPunctuatorToken, "??", true) {
		elseExpr, err := parser.parseCoalesceExpression()
		if err != nil {
			return ast.NewEmptyExpression(), err
		}

		return ast.NewCoalesceExpression(expr, elseExpr), nil
	}

	return expr, nil
}

func (parser *Parser) parseLogicalIncOrExpression1() (ast.IExpression, error) {
	// Spec: https://phplang.org/spec/10-expressions.html#grammar-logical-inc-OR-expression-1

	// logical-inc-OR-expression-1:
	//    logical-AND-expression-1
	//    logical-inc-OR-expression-1   ||   logical-AND-expression-1

	// TODO logical-inc-OR-expression-1
	return parser.parseLogicalAndExpression1()
}

func (parser *Parser) parseLogicalAndExpression1() (ast.IExpression, error) {
	// Spec: https://phplang.org/spec/10-expressions.html#grammar-logical-AND-expression-1

	// logical-AND-expression-1:
	//    bitwise-inc-OR-expression
	//    logical-AND-expression-1   &&   bitwise-inc-OR-expression

	// TODO logical-AND-expression-1
	return parser.parseBitwiseIncOrExpression()
}

func (parser *Parser) parseBitwiseIncOrExpression() (ast.IExpression, error) {
	// Spec: https://phplang.org/spec/10-expressions.html#grammar-bitwise-inc-OR-expression

	// bitwise-inc-OR-expression:
	//    bitwise-exc-OR-expression
	//    bitwise-inc-OR-expression   |   bitwise-exc-OR-expression

	// TODO bitwise-inc-OR-expression
	return parser.parseBitwiseExcOrExpression()
}

func (parser *Parser) parseBitwiseExcOrExpression() (ast.IExpression, error) {
	// Spec: https://phplang.org/spec/10-expressions.html#grammar-bitwise-exc-OR-expression

	// bitwise-exc-OR-expression:
	//    bitwise-AND-expression
	//    bitwise-exc-OR-expression   ^   bitwise-AND-expression

	// TODO bitwise-exc-OR-expression
	return parser.parseBitwiseAndExpression()
}

func (parser *Parser) parseBitwiseAndExpression() (ast.IExpression, error) {
	// Spec: https://phplang.org/spec/10-expressions.html#grammar-bitwise-AND-expression

	// bitwise-AND-expression:
	//    equality-expression
	//    bitwise-AND-expression   &   equality-expression

	// TODO bitwise-AND-expression
	return parser.parseEqualityExpression()
}

func (parser *Parser) parseEqualityExpression() (ast.IExpression, error) {
	// Spec: https://phplang.org/spec/10-expressions.html#grammar-equality-expression

	// equality-expression:
	//    relational-expression
	//    equality-expression   ==   relational-expression
	//    equality-expression   !=   relational-expression
	//    equality-expression   <>   relational-expression
	//    equality-expression   ===   relational-expression
	//    equality-expression   !==   relational-expression

	lhs, err := parser.parserRelationalExpression()
	if err != nil {
		return ast.NewEmptyExpression(), err
	}

	for parser.isTokenType(lexer.OperatorOrPunctuatorToken, false) && common.IsEqualityOperator(parser.at().Value) {
		operator := parser.eat().Value
		rhs, err := parser.parserRelationalExpression()
		if err != nil {
			return ast.NewEmptyExpression(), err
		}
		lhs = ast.NewEqualityExpression(lhs, operator, rhs)
	}
	return lhs, nil
}

func (parser *Parser) parserRelationalExpression() (ast.IExpression, error) {
	// Spec: https://phplang.org/spec/10-expressions.html#grammar-relational-expression

	// relational-expression:
	//    shift-expression
	//    relational-expression   <   shift-expression
	//    relational-expression   >   shift-expression
	//    relational-expression   <=   shift-expression
	//    relational-expression   >=   shift-expression
	//    relational-expression   <=>   shift-expression

	// TODO relational-expression
	return parser.parseShiftExpression()
}

func (parser *Parser) parseShiftExpression() (ast.IExpression, error) {
	// Spec: https://phplang.org/spec/10-expressions.html#grammar-shift-expression

	// shift-expression:
	//    additive-expression
	//    shift-expression   <<   additive-expression
	//    shift-expression   >>   additive-expression

	// TODO shift-expression
	return parser.parseAdditiveExpression()
}

func (parser *Parser) parseAdditiveExpression() (ast.IExpression, error) {
	// Spec: https://phplang.org/spec/10-expressions.html#grammar-additive-expression

	// additive-expression:
	//    multiplicative-expression
	//    additive-expression   +   multiplicative-expression
	//    additive-expression   -   multiplicative-expression
	//    additive-expression   .   multiplicative-expression

	lhs, err := parser.parseMultiplicativeExpression()
	if err != nil {
		return ast.NewEmptyExpression(), err
	}

	for parser.isTokenType(lexer.OperatorOrPunctuatorToken, false) && common.IsAdditiveOperator(parser.at().Value) {
		operator := parser.eat().Value
		rhs, err := parser.parseMultiplicativeExpression()
		if err != nil {
			return ast.NewEmptyExpression(), err
		}
		lhs = ast.NewAdditiveExpression(lhs, operator, rhs)
	}

	return lhs, nil
}

func (parser *Parser) parseMultiplicativeExpression() (ast.IExpression, error) {
	// Spec: https://phplang.org/spec/10-expressions.html#grammar-multiplicative-expression

	// multiplicative-expression:
	//    logical-NOT-expression
	//    multiplicative-expression   *   logical-NOT-expression
	//    multiplicative-expression   /   logical-NOT-expression
	//    multiplicative-expression   %   logical-NOT-expression

	lhs, err := parser.parseLogicalNotExpression()
	if err != nil {
		return ast.NewEmptyExpression(), err
	}

	for parser.isTokenType(lexer.OperatorOrPunctuatorToken, false) && common.IsMultiplicativeOperator(parser.at().Value) {
		operator := parser.eat().Value
		rhs, err := parser.parseLogicalNotExpression()
		if err != nil {
			return ast.NewEmptyExpression(), err
		}
		lhs = ast.NewMultiplicativeExpression(lhs, operator, rhs)
	}
	return lhs, nil
}

func (parser *Parser) parseLogicalNotExpression() (ast.IExpression, error) {
	// Spec: https://phplang.org/spec/10-expressions.html#grammar-logical-NOT-expression

	// logical-NOT-expression:
	//    instanceof-expression
	//    !   instanceof-expression

	isNotExpression := parser.isToken(lexer.OperatorOrPunctuatorToken, "!", true)

	expr, err := parser.parseInstanceofExpression()
	if err != nil {
		return ast.NewEmptyExpression(), err
	}

	if isNotExpression {
		return ast.NewLogicalNotExpression(expr), nil
	}
	return expr, nil
}

func (parser *Parser) parseInstanceofExpression() (ast.IExpression, error) {
	// Spec: https://phplang.org/spec/10-expressions.html#grammar-instanceof-expression

	// instanceof-expression:
	//    unary-expression
	//    instanceof-subject   instanceof   class-type-designator

	// instanceof-subject:
	//    instanceof-expression

	// TODO instanceof-expression
	return parser.parseUnaryExpression()
}

func (parser *Parser) parseUnaryExpression() (ast.IExpression, error) {
	// Spec: https://phplang.org/spec/10-expressions.html#grammar-unary-expression

	// unary-expression:
	//    exponentiation-expression
	//    unary-op-expression
	//    error-control-expression
	//    cast-expression

	// These operators associate right-to-left.

	// TODO exponentiation-expression

	// unary-op-expression

	// Spec: https://phplang.org/spec/10-expressions.html#grammar-unary-op-expression

	// unary-op-expression:
	//    unary-operator   unary-expression

	// unary-operator: one of
	//    +   -   ~

	// TODO unary-op-expression - constraints
	if parser.isTokenType(lexer.OperatorOrPunctuatorToken, false) && common.IsUnaryOperator(parser.at().Value) {
		// TODO unary-op-expression
	}

	// TODO error-control-expression
	// TODO cast-expression

	return parser.parsePrimaryExpression()
}

func (parser *Parser) parsePrimaryExpression() (ast.IExpression, error) {
	// Spec: https://phplang.org/spec/10-expressions.html#grammar-primary-expression

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

	// Spec: https://phplang.org/spec/10-expressions.html#grammar-variable

	// variable:
	//    callable-variable
	//    scoped-property-access-expression
	//    member-access-expression

	// ------------------- MARK: callable-variable -------------------

	// Spec: https://phplang.org/spec/10-expressions.html#grammar-callable-variable

	// callable-variable:
	//    simple-variable
	//    subscript-expression
	//    member-call-expression
	//    scoped-call-expression
	//    function-call-expression

	// ------------------- MARK: simple-variable -------------------

	// Spec: https://phplang.org/spec/10-expressions.html#simple-variable

	// simple-variable:
	//    variable-name

	if parser.isTokenType(lexer.VariableNameToken, false) {
		return ast.NewSimpleVariableExpression(ast.NewVariableNameExpression(parser.eat().Value)), nil
	}

	// Spec: https://phplang.org/spec/10-expressions.html#simple-variable

	// simple-variable:
	//    variable-name
	//    $   {   expression   }

	if parser.isToken(lexer.OperatorOrPunctuatorToken, "$", false) &&
		parser.next(0).TokenType == lexer.OperatorOrPunctuatorToken && parser.next(0).Value == "{" {
		parser.eatN(2)
		// Get expression
		expr, err := parser.parsePrimaryExpression()
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

	if parser.isToken(lexer.OperatorOrPunctuatorToken, "$", true) {
		if expr, err := parser.parsePrimaryExpression(); err != nil {
			return ast.NewEmptyExpression(), err
		} else {
			return ast.NewSimpleVariableExpression(expr), nil
		}
	}

	// TODO subscript-expression
	// TODO member-call-expression
	// TODO scoped-call-expression

	// ------------------- MARK: function-call-expression -------------------

	// Spec: https://phplang.org/spec/10-expressions.html#grammar-function-call-expression

	// function-call-expression:
	//    qualified-name   (   argument-expression-list(opt)   )
	//    qualified-name   (   argument-expression-list   ,   )
	//    callable-expression   (   argument-expression-list(opt)   )
	//    callable-expression   (   argument-expression-list   ,   )

	// argument-expression-list:
	//    argument-expression
	//    argument-expression-list   ,   argument-expression

	// argument-expression:
	//    variadic-unpacking
	//    expression

	// variadic-unpacking:
	//    ...   expression

	if parser.isTokenType(lexer.NameToken, false) &&
		parser.next(0).TokenType == lexer.OperatorOrPunctuatorToken && parser.next(0).Value == "(" {
		functionName := parser.eat().Value
		args := []ast.IExpression{}
		parser.eat() // Eat opening parentheses
		for {
			if parser.isToken(lexer.OperatorOrPunctuatorToken, ")", true) {
				break
			}

			arg, err := parser.parseExpression()
			if err != nil {
				return ast.NewEmptyExpression(), err
			}
			args = append(args, arg)

			if parser.isToken(lexer.OperatorOrPunctuatorToken, ",", true) ||
				parser.isToken(lexer.OperatorOrPunctuatorToken, ")", false) {
				continue
			}
			return ast.NewEmptyExpression(), fmt.Errorf("Expected \",\" or \")\". Got: %s", parser.at())
		}
		return ast.NewFunctionCallExpression(functionName, args), nil
	}
	// TODO function-call-expression
	// TODO function-call-expression - qualified-name

	// TODO scoped-property-access-expression
	// TODO member-access-expression

	// TODO class-constant-access-expression

	// ------------------- MARK: constant-access-expression -------------------

	// Spec: https://phplang.org/spec/10-expressions.html#grammar-constant-access-expression

	// constant-access-expression:
	//    qualified-name

	// A constant-access-expression evaluates to the value of the constant with name qualified-name.

	// Spec: https://phplang.org/spec/09-lexical-structure.html#grammar-qualified-name

	// qualified-name::
	//    namespace-name-as-a-prefix(opt)   name

	if parser.isTokenType(lexer.NameToken, false) {
		// TODO constant-access-expression - namespace-name-as-a-prefix
		// TODO constant-access-expression - check if name is a defined constant here or in interpreter
		return ast.NewConstantAccessExpression(parser.eat().Value), nil
	}

	// ------------------- MARK: literal -------------------

	if parser.isToken(lexer.KeywordToken, "FALSE", false) || parser.isToken(lexer.KeywordToken, "TRUE", false) ||
		parser.isTokenType(lexer.IntegerLiteralToken, false) || parser.isTokenType(lexer.FloatingLiteralToken, false) ||
		parser.isTokenType(lexer.StringLiteralToken, false) || parser.isToken(lexer.KeywordToken, "NULL", false) {
		return parser.parseLiteral()
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

func (parser *Parser) parseLiteral() (ast.IExpression, error) {
	// Spec: https://phplang.org/spec/10-expressions.html#grammar-literal

	// literal:
	//    integer-literal
	//    floating-literal
	//    string-literal

	// A literal evaluates to its value, as specified in the lexical specification for literals.

	// boolean-literal
	if parser.isToken(lexer.KeywordToken, "FALSE", true) {
		return ast.NewBooleanLiteralExpression(false), nil
	}
	if parser.isToken(lexer.KeywordToken, "TRUE", true) {
		return ast.NewBooleanLiteralExpression(true), nil
	}

	// integer-literal
	if parser.isTokenType(lexer.IntegerLiteralToken, false) {
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
	if parser.isTokenType(lexer.FloatingLiteralToken, false) {
		if common.IsFloatingLiteral(parser.at().Value) {
			return ast.NewFloatingLiteralExpression(common.FloatingLiteralToFloat64(parser.eat().Value)), nil
		}

		return ast.NewEmptyExpression(), fmt.Errorf("Parser error: Unsupported floating literal \"%s\"", parser.at().Value)
	}

	// string-literal
	if parser.isTokenType(lexer.StringLiteralToken, false) {
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

	// null-literal
	if parser.isToken(lexer.KeywordToken, "NULL", true) {
		return ast.NewNullLiteralExpression(), nil
	}

	return ast.NewEmptyExpression(), fmt.Errorf("parseLiteral: Unsupported literal: %s", parser.at())
}
