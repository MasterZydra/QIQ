package parser

import (
	"QIQ/cmd/qiq/ast"
	"QIQ/cmd/qiq/common"
	"QIQ/cmd/qiq/ini"
	"QIQ/cmd/qiq/lexer"
	"QIQ/cmd/qiq/phpError"
	"QIQ/cmd/qiq/position"
	"QIQ/cmd/qiq/stats"
	"slices"
	"strings"
)

type Parser struct {
	ini     *ini.Ini
	program *ast.Program
	lexer   *lexer.Lexer
	tokens  []*lexer.Token
	currPos int
	id      int64
}

func NewParser(ini *ini.Ini) *Parser { return &Parser{ini: ini} }

func (parser *Parser) init() {
	parser.program = ast.NewProgram()
	parser.lexer = lexer.NewLexer(parser.ini)
	parser.currPos = 0
}

func (parser *Parser) nextId() int64 {
	parser.id++
	return parser.id
}

func (parser *Parser) ProduceAST(sourceCode string, filename string) (*ast.Program, phpError.Error) {
	parser.init()

	var lexerErr error
	parser.tokens, lexerErr = parser.lexer.Tokenize(sourceCode, filename)
	if lexerErr != nil {
		return parser.program, phpError.NewParseError("%s", lexerErr.Error())
	}

	stat := stats.Start()
	defer stats.StopAndPrint(stat, "Parser")

	PrintParserCallstack("Parser callstack", nil)
	PrintParserCallstack("----------------", nil)

	for !parser.isEof() {
		if parser.isTokenType(lexer.StartTagToken, true) || parser.isTokenType(lexer.EndTagToken, true) ||
			parser.isToken(lexer.OpOrPuncToken, ";", true) {
			continue
		}
		stmt, err := parser.parseStmt()
		if err != nil {
			return parser.program, err
		}
		if stmt.GetKind() != ast.EmptyNode {
			parser.program.Append(stmt)
		}
	}

	return parser.program, nil
}

func (parser *Parser) parseMixedStmt(compoundEndKeywords []string) (ast.IStatement, phpError.Error) {
	return parser.parseMixedStmtRec(compoundEndKeywords, ast.NewCompoundStmt(parser.nextId(), []ast.IStatement{}))
}

func (parser *Parser) parseMixedStmtRec(compoundEndKeywords []string, textExprCompoundStmt *ast.CompoundStatement) (ast.IStatement, phpError.Error) {
	// Resolve text expressions
	if parser.isTextExpression(true) {
		PrintParserCallstack("text-expression (mixed-stmt)", parser)
		textExpr := ast.NewExpressionStmt(parser.nextId(), ast.NewTextExpr(parser.nextId(), parser.eat().Value))
		parser.isTokenType(lexer.StartTagToken, true)

		textExprCompoundStmt.Statements = append(textExprCompoundStmt.Statements, textExpr)

		if parser.isTokenType(lexer.KeywordToken, false) && slices.Contains(compoundEndKeywords, strings.ToLower(parser.at().Value)) {
			return textExprCompoundStmt, nil
		}

		println(parser.at().Value, parser.at().TokenType)

		stmt, err := parser.parseMixedStmtRec(compoundEndKeywords, textExprCompoundStmt)
		for err == nil || parser.isTextExpression(false) {
			println(parser.at().Value, parser.at().TokenType)
			if parser.isTokenType(lexer.KeywordToken, false) && slices.Contains(compoundEndKeywords, strings.ToLower(parser.at().Value)) {
				break
			}
			textExprCompoundStmt.Statements = append(textExprCompoundStmt.Statements, stmt)
			stmt, err = parser.parseMixedStmtRec(compoundEndKeywords, textExprCompoundStmt)
		}

		return textExprCompoundStmt, nil
	}

	return parser.parseStmt()
}

func (parser *Parser) parseStmt() (ast.IStatement, phpError.Error) {
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

	// Resolve text expressions
	if parser.isTextExpression(true) {
		PrintParserCallstack("text-expression", parser)
		stmt := ast.NewExpressionStmt(parser.nextId(), ast.NewTextExpr(parser.nextId(), parser.eat().Value))
		parser.isTokenType(lexer.StartTagToken, true)
		return stmt, nil
	}

	if parser.isTokenType(lexer.TextToken, false) {
		return ast.NewExpressionStmt(parser.nextId(), ast.NewTextExpr(parser.nextId(), parser.eat().Value)), nil
	}

	// -------------------------------------- compound-statement -------------------------------------- MARK: compound-statement

	// Spec: https://phplang.org/spec/11-statements.html#grammar-compound-statement

	// compound-statement:
	//    {   statement-list(opt)   }

	// statement-list:
	//    statement
	//    statement-list   statement

	// Supported statement: compound statement: `{ doThis(); doThat(); }`
	if parser.isToken(lexer.OpOrPuncToken, "{", true) {
		PrintParserCallstack("compound-statement", parser)
		statements := []ast.IStatement{}
		for !parser.isEof() && !parser.isToken(lexer.OpOrPuncToken, "}", false) {
			stmt, err := parser.parseStmt()
			if err != nil {
				return ast.NewEmptyStmt(), err
			}
			statements = append(statements, stmt)
		}

		if !parser.isToken(lexer.OpOrPuncToken, "}", true) {
			return ast.NewEmptyStmt(), NewExpectedError("}", parser.at())
		}
		return ast.NewCompoundStmt(parser.nextId(), statements), nil
	}

	// TODO named-label-statement

	// selection-statement
	if parser.isToken(lexer.KeywordToken, "if", false) || parser.isToken(lexer.KeywordToken, "switch", false) {
		return parser.parseSelectionStmt()
	}

	// iteration-statement
	if parser.isToken(lexer.KeywordToken, "while", false) || parser.isToken(lexer.KeywordToken, "do", false) ||
		parser.isToken(lexer.KeywordToken, "for", false) || parser.isToken(lexer.KeywordToken, "foreach", false) {
		return parser.parseIterationStmt()
	}

	// jump-statement
	if parser.isToken(lexer.KeywordToken, "goto", false) || parser.isToken(lexer.KeywordToken, "continue", false) ||
		parser.isToken(lexer.KeywordToken, "break", false) || parser.isToken(lexer.KeywordToken, "return", false) ||
		parser.isToken(lexer.KeywordToken, "throw", false) {
		return parser.parseJumpStmt()
	}

	// try-statement
	if parser.isToken(lexer.KeywordToken, "try", false) {
		return parser.parseTryStmt()
	}

	// -------------------------------------- declare-statement -------------------------------------- MARK: declare-statement

	// Spec: https://phplang.org/spec/11-statements.html?#the-declare-statement

	// declare-statement:
	//    declare   (   declare-directive   )   statement
	//    declare   (   declare-directive   )   :   statement-list   enddeclare   ;
	//    declare   (   declare-directive   )   ;

	// declare-directive:
	//    ticks   =   literal
	//    encoding   =   literal
	//    strict_types   =   literal

	// Supported statement: declare statement: `declare(strict_types = 1)`
	if parser.isToken(lexer.KeywordToken, "declare", false) {
		PrintParserCallstack("declare-statement", parser)
		pos := parser.eat().Position

		// (
		if !parser.isToken(lexer.OpOrPuncToken, "(", true) {
			return ast.NewEmptyStmt(), NewExpectedError("(", parser.at())
		}

		// declare-directive
		if !parser.isTokenType(lexer.NameToken, false) {
			return ast.NewEmptyStmt(), phpError.NewParseError("Unsupported declare '%s' in %s", parser.at().Value, parser.at().GetPosString())
		}
		if !slices.Contains([]string{"ticks", "encoding", "strict_types"}, parser.at().Value) {
			return ast.NewEmptyStmt(), phpError.NewParseError("Unsupported declare '%s' in %s", parser.at().Value, parser.at().GetPosString())
		}
		directive := parser.eat().Value

		// =
		if !parser.isToken(lexer.OpOrPuncToken, "=", true) {
			return ast.NewEmptyStmt(), NewExpectedError("=", parser.at())
		}

		// literal
		if !lexer.IsLiteral(parser.at()) {
			return ast.NewEmptyStmt(), phpError.NewParseError("Expected a literal. Got %s in %s", parser.at().TokenType, parser.at().GetPosString())
		}
		literal, err := parser.parseLiteral()
		if err != nil {
			return ast.NewEmptyStmt(), err
		}

		// Validate allowed values
		if directive == "strict_types" {
			if literal.GetKind() != ast.IntegerLiteralExpr {
				return ast.NewEmptyStmt(), phpError.NewParseError("Only 0 and 1 allowed as values for the declare directive 'strict_types' in %s", parser.at().GetPosString())
			}
			if literal.(*ast.IntegerLiteralExpression).Value < 0 || literal.(*ast.IntegerLiteralExpression).Value > 1 {
				return ast.NewEmptyStmt(), phpError.NewParseError("Only 0 and 1 allowed as values for the declare directive 'strict_types' in %s", parser.at().GetPosString())
			}
		}

		// )
		if !parser.isToken(lexer.OpOrPuncToken, ")", true) {
			return ast.NewEmptyStmt(), NewExpectedError(")", parser.at())
		}

		// ;
		if !parser.isToken(lexer.OpOrPuncToken, ";", true) {
			return ast.NewEmptyStmt(), NewExpectedError(";", parser.at())
		}

		// TODO declare-statement - statement
		// TODO declare-statement - statement-list

		return ast.NewDeclareStmt(parser.nextId(), pos, directive, literal), nil
	}

	// -------------------------------------- echo-statement -------------------------------------- MARK: echo-statement

	// Spec https://phplang.org/spec/11-statements.html#the-echo-statement

	// echo-statement:
	//    echo   expression-list   ;

	// expression-list:
	//    expression
	//    expression-list   ,   expression

	// Supported statement: echo statement: `echo "abc", 123, true;`
	if parser.isToken(lexer.KeywordToken, "echo", false) {
		PrintParserCallstack("echo-statement", parser)
		pos := parser.eat().Position
		expressions := make([]ast.IExpression, 0)
		for {
			expr, err := parser.parseExpr()
			if err != nil {
				return ast.NewEmptyStmt(), err
			}

			expressions = append(expressions, expr)

			if parser.isToken(lexer.OpOrPuncToken, ",", true) {
				continue
			}
			if parser.isToken(lexer.OpOrPuncToken, ";", true) {
				break
			}
			return ast.NewEmptyStmt(), phpError.NewParseError("Invalid echo statement detected. Got: %s", parser.at())
		}

		if len(expressions) == 0 {
			return ast.NewEmptyStmt(), phpError.NewParseError("Invalid echo statement detected. Got: %s", parser.at())
		}

		return ast.NewEchoStmt(parser.nextId(), pos, expressions), nil
	}

	// -------------------------------------- unset-statement -------------------------------------- MARK: unset-statement

	// Spec: https://phplang.org/spec/11-statements.html#grammar-unset-statement

	// unset-statement:
	//    unset   (   variable-list   ,opt   )   ;

	// variable-list:
	//    variable
	//    variable-list   ,   variable

	// Supported intrinsic: unset intrinsic: `unset($v);`
	if parser.isToken(lexer.KeywordToken, "unset", false) {
		PrintParserCallstack("unset-statement", parser)
		pos := parser.eat().Position
		if !parser.isToken(lexer.OpOrPuncToken, "(", true) {
			return ast.NewEmptyStmt(), NewExpectedError("(", parser.at())
		}
		args := []ast.IExpression{}
		for {
			if len(args) > 0 && parser.isToken(lexer.OpOrPuncToken, ")", true) {
				break
			}

			arg, err := parser.parseExpr()
			if err != nil {
				return ast.NewEmptyStmt(), err
			}
			if !ast.IsVariableExpr(arg) {
				return ast.NewEmptyStmt(), phpError.NewParseError("Fatal error: Cannot use unset() on the result of an expression")
			}
			args = append(args, arg)

			if parser.isToken(lexer.OpOrPuncToken, ",", true) ||
				parser.isToken(lexer.OpOrPuncToken, ")", false) {
				continue
			}
			return ast.NewEmptyStmt(), phpError.NewParseError("Expected \",\" or \")\". Got: %s", parser.at())
		}
		if !parser.isToken(lexer.OpOrPuncToken, ";", true) {
			return ast.NewEmptyStmt(), NewExpectedError(";", parser.at())
		}
		return ast.NewExpressionStmt(parser.nextId(), ast.NewUnsetIntrinsic(parser.nextId(), pos, args)), nil
	}

	// const-declaration
	if parser.isToken(lexer.KeywordToken, "const", false) {
		return parser.parseConstDeclaration()
	}

	// function-definition
	if parser.isToken(lexer.KeywordToken, "function", false) {
		return parser.parseFunctionDefinition()
	}

	// class-declaration
	if ((parser.isTokenType(lexer.KeywordToken, false) && common.IsClassModifierKeyword(parser.at().Value)) &&
		parser.next(0).TokenType == lexer.KeywordToken && parser.next(0).Value == "class") ||
		parser.isToken(lexer.KeywordToken, "class", false) {
		return parser.parseClassDeclaration()
	}

	// interface-declaration
	if parser.isToken(lexer.KeywordToken, "interface", false) {
		return parser.parseInterfaceDeclaration()
	}

	// TODO trait-declaration

	// -------------------------------------- namespace-definition -------------------------------------- MARK: namespace-definition

	// Spec: https://phplang.org/spec/18-namespaces.html#grammar-namespace-definition

	// namespace-definition:
	//    namespace   namespace-name   ;
	//    namespace   namespace-name(opt)   compound-statement

	// Spec: https://phplang.org/spec/09-lexical-structure.html#grammar-namespace-name

	// namespace-name::
	//    name
	//    namespace-name   \   name

	// Supported statement: namespace definition: `namespace My\Name\Space;`
	if parser.isToken(lexer.KeywordToken, "namespace", true) {
		PrintParserCallstack("namespace-definition", parser)
		namespace := []string{}

		for {
			if !parser.isTokenType(lexer.NameToken, false) && !parser.isTokenType(lexer.KeywordToken, false) {
				break
			}

			namespace = append(namespace, parser.eat().Value)

			if parser.isToken(lexer.OpOrPuncToken, "\\", true) {
				continue
			}
			break
		}

		if !parser.isToken(lexer.OpOrPuncToken, ";", true) {
			return ast.NewEmptyStmt(), NewExpectedError(";", parser.at())
		}

		// TODO namespace-declaration - support compound-statment

		// TODO Reuse namespaces within one request
		parser.at().Position.File.Namespace = position.NewNamespace(namespace)
		return ast.NewEmptyStmt(), nil
	}

	// TODO namespace-use-declaration

	// -------------------------------------- global-declaration -------------------------------------- MARK: global-declaration

	// Spec: https://phplang.org/spec/07-variables.html#grammar-global-declaration

	// global-declaration:
	//    global   variable-name-list   ;

	// variable-name-list:
	//    simple-variable
	//    variable-name-list   ,   simple-variable

	// Supported statement: global declaration: `global $var;`
	if parser.isToken(lexer.KeywordToken, "global", false) {
		PrintParserCallstack("global-declaration-statement", parser)
		pos := parser.eat().Position
		variables := []ast.IExpression{}

		for {
			variable, err := parser.parsePrimaryExpr()
			if err != nil {
				return ast.NewEmptyStmt(), err
			}
			if variable.GetKind() != ast.SimpleVariableExpr {
				return ast.NewEmptyStmt(), phpError.NewParseError("Global declaration expected a simple variable but got %s in %s", variable.GetKind(), variable.GetPosString())
			}

			variables = append(variables, variable)

			if parser.isToken(lexer.OpOrPuncToken, ";", false) {
				break
			}
			if parser.isToken(lexer.OpOrPuncToken, ",", true) {
				continue
			}

			return ast.NewEmptyStmt(), phpError.NewParseError("Global declaration - expected \";\" or \"$\" but got token \"%s\" in %s", parser.at(), pos.ToPosString())
		}

		if !parser.isToken(lexer.OpOrPuncToken, ";", true) {
			return ast.NewEmptyStmt(), phpError.NewParseError("Global declaration - unexpected token %s in %s", parser.at(), pos.ToPosString())
		}

		return ast.NewGlobalDeclarationStmt(parser.nextId(), pos, variables), nil
	}

	// TODO function-static-declaration

	// -------------------------------------- expression-statement -------------------------------------- MARK: expression-statement

	// Spec: https://phplang.org/spec/11-statements.html#grammar-expression-statement

	// expression-statement:
	//    expression(opt)   ;

	// If present, expression is evaluated for its side effects, if any, and any resulting value is discarded.
	// If expression is omitted, the statement is a null statement, which has no effect on execution.
	parser.isToken(lexer.OpOrPuncToken, ";", true)

	PrintParserCallstack("expression-statement", parser)
	if expr, err := parser.parseExpr(); err != nil {
		return ast.NewEmptyExpr(), err
	} else {
		if parser.isToken(lexer.OpOrPuncToken, ";", true) {
			return ast.NewExpressionStmt(parser.nextId(), expr), nil
		}
		return ast.NewEmptyExpr(),
			phpError.NewParseError(`Statement must end with a semicolon. Got: "%s" in %s`, parser.at().Value, parser.at().GetPosString())
	}
}

func (parser *Parser) parseConstDeclaration() (ast.IStatement, phpError.Error) {
	// -------------------------------------- const-declaration -------------------------------------- MARK: const-declaration

	// Spec: https://phplang.org/spec/14-classes.html#grammar-const-declaration

	// const-declaration:
	//    const   const-elements   ;

	// const-elements:
	//    const-element
	//    const-elements   ,   const-element

	// const-element:
	//    name   =   constant-expression

	// Supported statement: const statement: `const TRUTH = 42;`
	PrintParserCallstack("const-statement", parser)
	pos := parser.eat().Position
	if parser.at().TokenType != lexer.NameToken && parser.at().TokenType != lexer.KeywordToken {
		if err := parser.expectTokenType(lexer.NameToken, false); err != nil {
			return ast.NewEmptyStmt(), err
		}
	}
	for {
		name := parser.eat().Value
		if err := parser.expect(lexer.OpOrPuncToken, "=", true); err != nil {
			return ast.NewEmptyStmt(), err
		}
		// TODO parse constant-expression
		value, err := parser.parseExpr()
		if err != nil {
			return ast.NewEmptyStmt(), err
		}

		stmt := ast.NewConstDeclarationStmt(parser.nextId(), pos, name, value)
		if parser.isToken(lexer.OpOrPuncToken, ",", true) {
			parser.program.Append(stmt)
			continue
		}
		if parser.isToken(lexer.OpOrPuncToken, ";", true) {
			return stmt, nil
		}
		return ast.NewEmptyStmt(), phpError.NewParseError("Const declaration - unexpected token %s", parser.at())
	}
}

func (parser *Parser) parseSelectionStmt() (ast.IStatement, phpError.Error) {
	// -------------------------------------- selection-statement -------------------------------------- MARK: selection-statement

	// Spec: https://phplang.org/spec/11-statements.html#grammar-selection-statement

	// selection-statement:
	//    if-statement
	//    switch-statement

	// Supported statement: if statement: `if (true) { ... } elseif (false) { ... } else { ... }`
	if parser.isToken(lexer.KeywordToken, "if", false) {
		// Spec: https://phplang.org/spec/11-statements.html#grammar-if-statement

		// if-statement:
		//    if   (   expression   )   statement   elseif-clauses-1(opt)   else-clause-1(opt)
		//    if   (   expression   )   :   statement-list   elseif-clauses-2(opt)   else-clause-2(opt)   endif   ;

		// elseif-clauses-1:
		//    elseif-clause-1
		//    elseif-clauses-1   elseif-clause-1

		// elseif-clause-1:
		//    elseif   (   expression   )   statement

		// else-clause-1:
		//    else   statement

		// elseif-clauses-2:
		//    elseif-clause-2
		//    elseif-clauses-2   elseif-clause-2

		// elseif-clause-2:
		//    elseif   (   expression   )   :   statement-list

		// else-clause-2:
		//    else   :   statement-list

		PrintParserCallstack("if-statement", parser)

		// if
		ifPos := parser.eat().Position
		if !parser.isToken(lexer.OpOrPuncToken, "(", true) {
			return ast.NewEmptyStmt(), NewExpectedError("(", parser.at())
		}

		condition, err := parser.parseExpr()
		if err != nil {
			return ast.NewEmptyStmt(), err
		}

		if !parser.isToken(lexer.OpOrPuncToken, ")", true) {
			return ast.NewEmptyStmt(), NewExpectedError(")", parser.at())
		}

		isAltSytax := parser.isToken(lexer.OpOrPuncToken, ":", true)

		var ifBlock ast.IStatement
		if !isAltSytax {
			ifBlock, err = parser.parseStmt()
			if err != nil {
				return ast.NewEmptyStmt(), err
			}
		} else {
			statements := []ast.IStatement{}
			for !parser.isToken(lexer.KeywordToken, "elseif", false) && !parser.isToken(lexer.KeywordToken, "else", false) &&
				!parser.isToken(lexer.KeywordToken, "endif", false) {
				statement, err := parser.parseMixedStmt([]string{"elseif", "else", "endif"})
				if err != nil {
					return ast.NewEmptyStmt(), err
				}
				statements = append(statements, statement)
			}
			ifBlock = ast.NewCompoundStmt(parser.nextId(), statements)
		}

		// elseif
		elseIf := []*ast.IfStatement{}
		for parser.isToken(lexer.KeywordToken, "elseif", false) {
			elseIfPos := parser.eat().Position
			if !parser.isToken(lexer.OpOrPuncToken, "(", true) {
				return ast.NewEmptyStmt(), NewExpectedError("(", parser.at())
			}

			elseIfCondition, err := parser.parseExpr()
			if err != nil {
				return ast.NewEmptyStmt(), err
			}

			if !parser.isToken(lexer.OpOrPuncToken, ")", true) {
				return ast.NewEmptyStmt(), NewExpectedError(")", parser.at())
			}

			if isAltSytax && !parser.isToken(lexer.OpOrPuncToken, ":", true) {
				return ast.NewEmptyStmt(), NewExpectedError(":", parser.at())
			}

			var elseIfBlock ast.IStatement
			if !isAltSytax {
				elseIfBlock, err = parser.parseStmt()
				if err != nil {
					return ast.NewEmptyStmt(), err
				}
			} else {
				statements := []ast.IStatement{}
				for !parser.isToken(lexer.KeywordToken, "elseif", false) && !parser.isToken(lexer.KeywordToken, "else", false) &&
					!parser.isToken(lexer.KeywordToken, "endif", false) {
					statement, err := parser.parseMixedStmt([]string{"elseif", "else", "endif"})
					if err != nil {
						return ast.NewEmptyStmt(), err
					}
					statements = append(statements, statement)
				}
				elseIfBlock = ast.NewCompoundStmt(parser.nextId(), statements)
			}

			elseIf = append(elseIf, ast.NewIfStmt(parser.nextId(), elseIfPos, elseIfCondition, elseIfBlock, nil, nil))
		}

		// else
		var elseBlock ast.IStatement = nil
		if parser.isToken(lexer.KeywordToken, "else", true) {
			if isAltSytax && !parser.isToken(lexer.OpOrPuncToken, ":", true) {
				return ast.NewEmptyStmt(), NewExpectedError(":", parser.at())
			}

			if !isAltSytax {
				elseBlock, err = parser.parseStmt()
				if err != nil {
					return ast.NewEmptyStmt(), err
				}
			} else {
				statements := []ast.IStatement{}
				for !parser.isToken(lexer.KeywordToken, "endif", false) {
					statement, err := parser.parseMixedStmt([]string{"endif"})
					if err != nil {
						return ast.NewEmptyStmt(), err
					}
					statements = append(statements, statement)
				}
				elseBlock = ast.NewCompoundStmt(parser.nextId(), statements)
			}
		}

		if isAltSytax && !parser.isToken(lexer.KeywordToken, "endif", true) {
			return ast.NewEmptyStmt(), NewExpectedError("endif", parser.at())
		}
		if isAltSytax && !parser.isToken(lexer.OpOrPuncToken, ";", true) {
			return ast.NewEmptyStmt(), NewExpectedError(";", parser.at())
		}

		return ast.NewIfStmt(parser.nextId(), ifPos, condition, ifBlock, elseIf, elseBlock), nil
	}

	// TODO switch-statement
	// if parser.isToken(lexer.KeywordToken, "switch", false) {
	// }

	return ast.NewEmptyStmt(), phpError.NewParseError("Unsupported selection statement %s", parser.at())
}

func (parser *Parser) parseIterationStmt() (ast.IStatement, phpError.Error) {
	// -------------------------------------- iteration-statement -------------------------------------- MARK: iteration-statement

	// Spec: https://phplang.org/spec/11-statements.html#grammar-iteration-statement

	// iteration-statement:
	//    while-statement
	//    do-statement
	//    for-statement
	//    foreach-statement

	// Supported statement: while statement: `while (true) { ... }`
	if parser.isToken(lexer.KeywordToken, "while", false) {
		// Spec: https://phplang.org/spec/11-statements.html#grammar-iteration-statement

		// while-statement:
		//    while   (   expression   )   statement
		//    while   (   expression   )   :   statement-list   endwhile   ;

		PrintParserCallstack("while-statement", parser)

		// condition
		whilePos := parser.eat().Position
		if !parser.isToken(lexer.OpOrPuncToken, "(", true) {
			return ast.NewEmptyStmt(), NewExpectedError("(", parser.at())
		}

		condition, err := parser.parseExpr()
		if err != nil {
			return ast.NewEmptyStmt(), err
		}

		if !parser.isToken(lexer.OpOrPuncToken, ")", true) {
			return ast.NewEmptyStmt(), NewExpectedError(")", parser.at())
		}

		isAltSytax := parser.isToken(lexer.OpOrPuncToken, ":", true)

		var block ast.IStatement
		if !isAltSytax {
			block, err = parser.parseStmt()
			if err != nil {
				return ast.NewEmptyStmt(), err
			}
		} else {
			statements := []ast.IStatement{}
			for !parser.isToken(lexer.KeywordToken, "endwhile", false) {
				statement, err := parser.parseMixedStmt([]string{"endwhile"})
				if err != nil {
					return ast.NewEmptyStmt(), err
				}
				statements = append(statements, statement)
			}
			block = ast.NewCompoundStmt(parser.nextId(), statements)
		}

		if isAltSytax && !parser.isToken(lexer.KeywordToken, "endwhile", true) {
			return ast.NewEmptyStmt(), NewExpectedError("endwhile", parser.at())
		}
		if isAltSytax && !parser.isToken(lexer.OpOrPuncToken, ";", true) {
			return ast.NewEmptyStmt(), NewExpectedError(";", parser.at())
		}

		return ast.NewWhileStmt(parser.nextId(), whilePos, condition, block), nil
	}

	// Supported statement: do statement: `do { ... } while (true);`
	if parser.isToken(lexer.KeywordToken, "do", false) {
		// Spec: https://phplang.org/spec/11-statements.html#grammar-do-statement

		// do-statement:
		//    do   statement   while   (   expression   )   ;

		PrintParserCallstack("do-statement", parser)

		doPos := parser.eat().Position

		// statement
		block, err := parser.parseStmt()
		if err != nil {
			return ast.NewEmptyStmt(), err
		}

		// condition
		if !parser.isToken(lexer.KeywordToken, "while", true) {
			return ast.NewEmptyStmt(), NewExpectedError("while", parser.at())
		}

		if !parser.isToken(lexer.OpOrPuncToken, "(", true) {
			return ast.NewEmptyStmt(), NewExpectedError("(", parser.at())
		}

		condition, err := parser.parseExpr()
		if err != nil {
			return ast.NewEmptyStmt(), err
		}

		if !parser.isToken(lexer.OpOrPuncToken, ")", true) {
			return ast.NewEmptyStmt(), NewExpectedError(")", parser.at())
		}

		if !parser.isToken(lexer.OpOrPuncToken, ";", true) {
			return ast.NewEmptyStmt(), NewExpectedError(";", parser.at())
		}

		return ast.NewDoStmt(parser.nextId(), doPos, condition, block), nil
	}

	// Supported statement: for statement: `for (...; ...; ...) { ... }`
	if parser.isToken(lexer.KeywordToken, "for", false) {
		// Spec: https://phplang.org/spec/11-statements.html#grammar-for-statement

		// for-statement:
		//    for   (   for-initializer(opt)   ;   for-control(opt)   ;   for-end-of-loop(opt)   )   statement
		//    for   (   for-initializer(opt)   ;   for-control(opt)   ;   for-end-of-loop(opt)   )   :   statement-list   endfor   ;

		// for-initializer:
		//    for-expression-group

		// for-control:
		//    for-expression-group

		// for-end-of-loop:
		//    for-expression-group

		// for-expression-group:
		//    expression
		//    for-expression-group   ,   expression

		PrintParserCallstack("for-statement", parser)

		pos := parser.eat().Position

		if !parser.isToken(lexer.OpOrPuncToken, "(", true) {
			return ast.NewEmptyStmt(), NewExpectedError("(", parser.at())
		}

		var err phpError.Error = nil
		parseExprGroup := func() (*ast.CompoundStatement, phpError.Error) {
			stmts := []ast.IStatement{}
			for {
				expr, err := parser.parseExpr()
				if err != nil {
					return nil, err
				}
				stmts = append(stmts, expr)

				if !parser.isToken(lexer.OpOrPuncToken, ",", true) {
					break
				}
			}
			return ast.NewCompoundStmt(parser.nextId(), stmts), nil
		}

		// for-initializer
		var forInitializer *ast.CompoundStatement = nil
		if !parser.isToken(lexer.OpOrPuncToken, ";", false) {
			forInitializer, err = parseExprGroup()
			if err != nil {
				return ast.NewEmptyExpr(), err
			}
		}

		if !parser.isToken(lexer.OpOrPuncToken, ";", true) {
			return ast.NewEmptyStmt(), NewExpectedError(";", parser.at())
		}

		// for-control
		var forControl *ast.CompoundStatement = nil
		if !parser.isToken(lexer.OpOrPuncToken, ";", false) {
			forControl, err = parseExprGroup()
			if err != nil {
				return ast.NewEmptyExpr(), err
			}
		}

		if !parser.isToken(lexer.OpOrPuncToken, ";", true) {
			return ast.NewEmptyStmt(), NewExpectedError(";", parser.at())
		}

		// for-end-of-loop
		var forEndOfLoop *ast.CompoundStatement = nil
		if !parser.isToken(lexer.OpOrPuncToken, ")", false) {
			forEndOfLoop, err = parseExprGroup()
			if err != nil {
				return ast.NewEmptyExpr(), err
			}
		}

		if !parser.isToken(lexer.OpOrPuncToken, ")", true) {
			return ast.NewEmptyStmt(), NewExpectedError(")", parser.at())
		}

		isAltSytax := parser.isToken(lexer.OpOrPuncToken, ":", true)

		var block ast.IStatement
		if !isAltSytax {
			block, err = parser.parseStmt()
			if err != nil {
				return ast.NewEmptyStmt(), err
			}
		} else {
			statements := []ast.IStatement{}
			for !parser.isToken(lexer.KeywordToken, "endfor", false) {
				statement, err := parser.parseStmt()
				if err != nil {
					return ast.NewEmptyStmt(), err
				}
				statements = append(statements, statement)
			}
			block = ast.NewCompoundStmt(parser.nextId(), statements)
		}

		if isAltSytax && !parser.isToken(lexer.KeywordToken, "endfor", true) {
			return ast.NewEmptyStmt(), NewExpectedError("endfor", parser.at())
		}
		if isAltSytax && !parser.isToken(lexer.OpOrPuncToken, ";", true) {
			return ast.NewEmptyStmt(), NewExpectedError(";", parser.at())
		}

		return ast.NewForStmt(parser.nextId(), pos, forInitializer, forControl, forEndOfLoop, block), nil
	}

	// Supported statement: foreach statement: `foreach ($entries as $key => $entry) { ... }`
	if parser.isToken(lexer.KeywordToken, "foreach", false) {
		// Spec: https://phplang.org/spec/11-statements.html#the-foreach-statement

		// foreach-statement:
		//    foreach   (   foreach-collection-name   as   foreach-key(opt)   foreach-value   )   statement
		//    foreach   (   foreach-collection-name   as   foreach-key(opt)   foreach-value   )   :   statement-list   endforeach   ;

		// foreach-collection-name:
		//    expression

		// foreach-key:
		//    expression   =>

		// foreach-value:
		//    &(opt)   expression
		//    list-intrinsic

		PrintParserCallstack("foreach-statement", parser)

		pos := parser.eat().Position

		if !parser.isToken(lexer.OpOrPuncToken, "(", true) {
			return ast.NewEmptyStmt(), NewExpectedError("(", parser.at())
		}

		collection, err := parser.parseExpr()
		if err != nil {
			return ast.NewEmptyStmt(), err
		}

		if !parser.isToken(lexer.KeywordToken, "as", true) {
			return ast.NewEmptyStmt(), NewExpectedError("as", parser.at())
		}

		byRef := parser.isToken(lexer.OpOrPuncToken, "&", false)
		var byRefPos *position.Position = nil
		if byRef {
			byRefPos = parser.eat().Position
		}

		valuePos := parser.at().GetPosString()
		value, err := parser.parseExpr()
		if err != nil {
			return ast.NewEmptyStmt(), err
		}
		// Check if value is a variable name
		if value.GetKind() != ast.SimpleVariableExpr {
			return ast.NewEmptyStmt(), phpError.NewParseError("Syntax error, unexpected token \"%s\", expecting variable name in %s", parser.at().Value, valuePos)
		}

		var key ast.IExpression = nil
		if parser.isToken(lexer.OpOrPuncToken, "=>", true) {
			if byRef {
				return ast.NewEmptyStmt(), phpError.NewParseError("Syntax error, key cannot be by reference in %s", byRefPos.ToPosString())
			}
			byRef = parser.isToken(lexer.OpOrPuncToken, "&", true)

			key = value
			valuePos = parser.at().GetPosString()
			value, err = parser.parseExpr()
			if err != nil {
				return ast.NewEmptyStmt(), err
			}
			// Check if value is a variable name
			if value.GetKind() != ast.SimpleVariableExpr {
				return ast.NewEmptyStmt(), phpError.NewParseError("Syntax error, unexpected token \"%s\", expecting variable name in %s", parser.at().Value, valuePos)
			}
		}

		if !parser.isToken(lexer.OpOrPuncToken, ")", true) {
			return ast.NewEmptyStmt(), NewExpectedError(")", parser.at())
		}

		isAltSytax := parser.isToken(lexer.OpOrPuncToken, ":", true)

		var block ast.IStatement
		if !isAltSytax {
			block, err = parser.parseStmt()
			if err != nil {
				return ast.NewEmptyStmt(), err
			}
		} else {
			statements := []ast.IStatement{}
			for !parser.isToken(lexer.KeywordToken, "endforeach", false) {
				statement, err := parser.parseStmt()
				if err != nil {
					return ast.NewEmptyStmt(), err
				}
				statements = append(statements, statement)
			}
			block = ast.NewCompoundStmt(parser.nextId(), statements)
		}

		if isAltSytax && !parser.isToken(lexer.KeywordToken, "endforeach", true) {
			return ast.NewEmptyStmt(), NewExpectedError("endforeach", parser.at())
		}
		if isAltSytax && !parser.isToken(lexer.OpOrPuncToken, ";", true) {
			return ast.NewEmptyStmt(), NewExpectedError(";", parser.at())
		}

		return ast.NewForeachStmt(parser.nextId(), pos, collection, key, value, byRef, block), nil
	}

	return ast.NewEmptyStmt(), phpError.NewParseError("Unsupported iteration statement '%s' in %s", parser.at().Value, parser.at().GetPosString())
}

func (parser *Parser) parseJumpStmt() (ast.IStatement, phpError.Error) {
	// -------------------------------------- jump-statement -------------------------------------- MARK: jump-statement

	// Spec: https://phplang.org/spec/11-statements.html#grammar-jump-statement

	// jump-statement:
	//    goto-statement
	//    continue-statement
	//    break-statement
	//    return-statement
	//    throw-statement

	// TODO goto-statement

	// Supported statement: continue statement: `continue (2);`
	if parser.isToken(lexer.KeywordToken, "continue", false) {
		// Spec: https://phplang.org/spec/11-statements.html#grammar-continue-statement

		// continue-statement:
		//    continue   breakout-level(opt)   ;

		// breakout-level:
		//   integer-literal
		//   (   breakout-level   )

		PrintParserCallstack("continue-statement", parser)

		pos := parser.eat().Position

		var expr ast.IExpression = nil
		var err phpError.Error

		if !parser.isToken(lexer.OpOrPuncToken, ";", false) {
			isParenthesized := parser.isToken(lexer.OpOrPuncToken, "(", true)

			expr, err = parser.parseExpr()
			if err != nil {
				return ast.NewEmptyStmt(), err
			}

			if isParenthesized && !parser.isToken(lexer.OpOrPuncToken, ")", true) {
				return ast.NewEmptyExpr(), phpError.NewError("Expected closing parentheses. Got %s", parser.at())
			}
		}

		if !parser.isToken(lexer.OpOrPuncToken, ";", true) {
			return ast.NewEmptyStmt(), NewExpectedError(";", parser.at())
		}

		if expr == nil {
			expr = ast.NewIntegerLiteralExpr(parser.nextId(), nil, 1)
		}

		return ast.NewContinueStmt(parser.nextId(), pos, expr), nil
	}

	// Supported statement: break statement: `break 1;`
	if parser.isToken(lexer.KeywordToken, "break", false) {
		// Spec: https://phplang.org/spec/11-statements.html#grammar-break-statement

		// break-statement:
		//    break   breakout-level(opt)   ;

		// breakout-level:
		//   integer-literal
		//   (   breakout-level   )

		PrintParserCallstack("break-statement", parser)

		pos := parser.eat().Position

		var expr ast.IExpression = nil
		var err phpError.Error

		if !parser.isToken(lexer.OpOrPuncToken, ";", false) {
			isParenthesized := parser.isToken(lexer.OpOrPuncToken, "(", true)

			expr, err = parser.parseExpr()
			if err != nil {
				return ast.NewEmptyStmt(), err
			}

			if isParenthesized && !parser.isToken(lexer.OpOrPuncToken, ")", true) {
				return ast.NewEmptyExpr(), phpError.NewError("Expected closing parentheses. Got %s", parser.at())
			}
		}

		if !parser.isToken(lexer.OpOrPuncToken, ";", true) {
			return ast.NewEmptyStmt(), NewExpectedError(";", parser.at())
		}

		if expr == nil {
			expr = ast.NewIntegerLiteralExpr(parser.nextId(), nil, 1)
		}

		return ast.NewBreakStmt(parser.nextId(), pos, expr), nil
	}

	// -------------------------------------- return-statement -------------------------------------- MARK: return-statement

	// Spec: https://phplang.org/spec/11-statements.html#grammar-return-statement

	// return-statement:
	//    return   expressionopt   ;

	// Supported statement: return statement: `return 42;`
	if parser.isToken(lexer.KeywordToken, "return", false) {
		PrintParserCallstack("return-statement", parser)
		pos := parser.eat().Position
		var expr ast.IExpression = nil
		if !parser.isToken(lexer.OpOrPuncToken, ";", false) {
			var err phpError.Error
			expr, err = parser.parseExpr()
			if err != nil {
				return ast.NewEmptyStmt(), err
			}
		}
		if !parser.isToken(lexer.OpOrPuncToken, ";", true) {
			return ast.NewEmptyStmt(), phpError.NewParseError("Expected: \";\". Got: \"%s\"", parser.at())
		}

		return ast.NewReturnStmt(parser.nextId(), pos, expr), nil
	}

	// -------------------------------------- throw-statement -------------------------------------- MARK: throw-statement

	// Spec: https://phplang.org/spec/11-statements.html#grammar-throw-statement

	// throw-statement:
	//    throw   expression   ;

	// Supported statement: throw statement: `throw new Exception();`
	if parser.isToken(lexer.KeywordToken, "throw", false) {
		PrintParserCallstack("throw-statement", parser)
		pos := parser.eat().Position
		expr, err := parser.parseExpr()
		if err != nil {
			return ast.NewEmptyStmt(), err
		}
		if !parser.isToken(lexer.OpOrPuncToken, ";", true) {
			return ast.NewEmptyStmt(), phpError.NewParseError("Expected: \";\". Got: \"%s\"", parser.at())
		}
		return ast.NewThrowStmt(parser.nextId(), pos, expr), nil
	}

	return ast.NewEmptyStmt(), phpError.NewParseError("Unsupported jump statement '%s' in %s", parser.at().Value, parser.at().GetPosString())
}

func (parser *Parser) parseTryStmt() (ast.IStatement, phpError.Error) {
	// -------------------------------------- try-statement -------------------------------------- MARK: try-statement
	// Spec: https://phplang.org/spec/11-statements.html#the-try-statement

	// try-statement:
	//    try   compound-statement   catch-clauses
	//    try   compound-statement   finally-clause
	//    try   compound-statement   catch-clauses   finally-clause

	// catch-clauses:
	//    catch-clause
	//    catch-clauses   catch-clause

	// catch-clause:
	//    catch   (   catch-name-list   variable-name(opt)   )   compound-statement
	// Spec-Fix: Since PHP 8 the variable name is optional.

	// catch-name-list:
	//    qualified-name
	//    catch-name-list   |   qualified-name

	// finally-clause:
	//    finally   compound-statement

	// Supported statement: try statement: `try { ... } catch (...) { ... } finally { ... }`
	PrintParserCallstack("try-statement", parser)

	// Eat "try"
	pos := parser.eat().Position

	// Body
	if !parser.isToken(lexer.OpOrPuncToken, "{", false) {
		return ast.NewEmptyStmt(), NewExpectedError("{", parser.at())
	}
	body, err := parser.parseStmt()
	if err != nil {
		return ast.NewEmptyStmt(), err
	}
	if body.GetKind() != ast.CompoundStmt {
		return ast.NewEmptyStmt(), phpError.NewParseError("Expected compound statement. Got %s in %s", body.GetKind(), body.GetPosString())
	}

	tryStmt := ast.NewTryStmt(parser.nextId(), pos, body.(*ast.CompoundStatement))

	// Catch
	for parser.isToken(lexer.KeywordToken, "catch", true) {
		if !parser.isToken(lexer.OpOrPuncToken, "(", true) {
			return ast.NewEmptyStmt(), NewExpectedError("(", parser.at())
		}

		catchNames := []string{}
		for {
			if !common.IsQualifiedName(parser.at().Value) {
				return ast.NewEmptyStmt(), phpError.NewParseError("Expected qualified name in %s", parser.at().GetPosString())
			}
			catchNames = append(catchNames, parser.eat().Value)

			if parser.isToken(lexer.OpOrPuncToken, "|", true) {
				continue
			}
			break
		}

		variableName := ""
		if !parser.isToken(lexer.OpOrPuncToken, ")", false) {
			if !common.IsVariableName(parser.at().Value) {
				return ast.NewEmptyStmt(), phpError.NewParseError("Expected variable name in %s", parser.at().GetPosString())
			}
			variableName = parser.eat().Value
		}
		if !parser.isToken(lexer.OpOrPuncToken, ")", true) {
			return ast.NewEmptyStmt(), NewExpectedError(")", parser.at())
		}

		if !parser.isToken(lexer.OpOrPuncToken, "{", false) {
			return ast.NewEmptyStmt(), NewExpectedError("{", parser.at())
		}
		body, err := parser.parseStmt()
		if err != nil {
			return ast.NewEmptyStmt(), err
		}
		if body.GetKind() != ast.CompoundStmt {
			return ast.NewEmptyStmt(), phpError.NewParseError("Expected compound statement. Got %s in %s", body.GetKind(), body.GetPosString())
		}

		tryStmt.AddCatch(ast.CatchStatement{ErrorType: catchNames, VariableName: variableName, Body: body.(*ast.CompoundStatement)})
	}

	// Finally
	if parser.isToken(lexer.KeywordToken, "finally", true) {
		if !parser.isToken(lexer.OpOrPuncToken, "{", false) {
			return ast.NewEmptyStmt(), NewExpectedError("{", parser.at())
		}
		finally, err := parser.parseStmt()
		if err != nil {
			return ast.NewEmptyStmt(), err
		}
		if finally.GetKind() != ast.CompoundStmt {
			return ast.NewEmptyStmt(), phpError.NewParseError("Expected compound statement. Got %s in %s", finally.GetKind(), finally.GetPosString())
		}
		tryStmt.Finally = finally.(*ast.CompoundStatement)
	}

	if tryStmt.Finally == nil && len(tryStmt.Catches) == 0 {
		return ast.NewEmptyStmt(), phpError.NewError("Cannot use try without catch of finally in %s", tryStmt.GetPosString())
	}

	return tryStmt, nil
}

func (parser *Parser) parseFunctionDefinition() (ast.IStatement, phpError.Error) {
	// -------------------------------------- function-definition -------------------------------------- MARK: function-definition
	// Spec: https://phplang.org/spec/13-functions.html#grammar-function-definition

	// function-definition:
	//    function-definition-header   compound-statement

	// function-definition-header:
	//    function   &(opt)   name   (   parameter-declaration-list(opt)   )   return-type(opt)

	// parameter-declaration-list:
	//    simple-parameter-declaration-list
	//    variadic-declaration-list

	// simple-parameter-declaration-list:
	//    parameter-declaration
	//    parameter-declaration-list   ,   parameter-declaration

	// variadic-declaration-list:
	//    simple-parameter-declaration-list   ,   variadic-parameter
	//    variadic-parameter

	// parameter-declaration:
	//    type-declaration(opt)   &(opt)   variable-name   default-argument-specifier(opt)

	// variadic-parameter:
	//    type-declaration(opt)   &(opt)   ...   variable-name

	// type-declaration:
	//    ?(opt)   base-type-declaration

	// return-type:
	//    :   type-declaration
	//    :   void

	// base-type-declaration:
	//    array
	//    callable
	//    iterable
	//    scalar-type
	//    qualified-name

	// scalar-type:
	//    bool
	//    float
	//    int
	//    string

	// default-argument-specifier:
	//    =   constant-expression

	// Supported statement: function definition: `function func1($param1) { ... }`
	PrintParserCallstack("function-definition", parser)
	if !parser.isToken(lexer.KeywordToken, "function", false) {
		return ast.NewEmptyStmt(), NewExpectedError("function", parser.at())
	}

	pos := parser.eat().Position

	// TODO byRef: function-definition - &(opt)

	if parser.at().TokenType != lexer.NameToken {
		return ast.NewEmptyStmt(), phpError.NewParseError("Function name expected. Got %s", parser.at().TokenType)
	}

	functionName := parser.at().Value
	functionNamePos := parser.eat().GetPosString()
	if !common.IsName(functionName) {
		return ast.NewEmptyExpr(), phpError.NewParseError("\"%s\" is not a valid function name in %s", functionName, functionNamePos)
	}
	if common.IsReservedName(functionName) {
		return ast.NewEmptyStmt(), phpError.NewError("Cannot use \"%s\" as a function name as it is reserved in %s", functionName, functionNamePos)
	}

	if !parser.isToken(lexer.OpOrPuncToken, "(", true) {
		return ast.NewEmptyStmt(), NewExpectedError("(", parser.at())
	}

	parameters, err := parser.parseFunctionParameters()
	if err != nil {
		return ast.NewEmptyStmt(), err
	}

	if !parser.isToken(lexer.OpOrPuncToken, ")", true) {
		return ast.NewEmptyStmt(), NewExpectedError(")", parser.at())
	}

	returnTypes := []string{}
	if parser.isToken(lexer.OpOrPuncToken, ":", true) {
		for parser.at().TokenType == lexer.KeywordToken && common.IsReturnTypeKeyword(parser.at().Value) {
			returnTypes = append(returnTypes, strings.ToLower(parser.eat().Value))
			if parser.isToken(lexer.OpOrPuncToken, "|", true) {
				continue
			}
			break
		}
	}

	body, err := parser.parseStmt()
	if err != nil {
		return ast.NewEmptyStmt(), err
	}
	if body.GetKind() != ast.CompoundStmt {
		return ast.NewEmptyStmt(), phpError.NewParseError("Expected compound statement. Got %s", body.GetKind())
	}

	return ast.NewFunctionDefinitionStmt(parser.nextId(), pos, functionName, parameters, body.(*ast.CompoundStatement), returnTypes), nil
}

func (parser *Parser) parseFunctionParameters() ([]ast.FunctionParameter, phpError.Error) {
	parameters := []ast.FunctionParameter{}
	if !parser.isToken(lexer.OpOrPuncToken, ")", false) {
		for {
			// Allow trailing comma
			if parser.isToken(lexer.OpOrPuncToken, ")", false) {
				break
			}

			// type-declaration
			paramTypes := []string{}
			if parser.isPhpType(parser.at()) {
				var err phpError.Error
				paramTypes, err = parser.getTypes(true)
				if err != nil {
					return parameters, err
				}
			}

			byRef := parser.isToken(lexer.OpOrPuncToken, "&", true)

			if parser.at().TokenType != lexer.VariableNameToken {
				return parameters, phpError.NewParseError("Expected variable. Got \"%s\" (%s) in %s", parser.at().Value, parser.at().TokenType, parser.at().GetPosString())
			}

			paramName := parser.eat().Value

			// TODO parse constant-expression
			var defaultValue ast.IExpression = nil
			if parser.isToken(lexer.OpOrPuncToken, "=", true) {
				var err phpError.Error
				defaultValue, err = parser.parseExpr()
				if err != nil {
					return parameters, err
				}
			}

			parameters = append(parameters, ast.NewFunctionParam(byRef, paramName, paramTypes, defaultValue))

			if parser.isToken(lexer.OpOrPuncToken, ",", true) {
				continue
			}
			if parser.isToken(lexer.OpOrPuncToken, ")", false) {
				break
			}
			return parameters, phpError.NewParseError("Expected \",\" or \")\". Got %s", parser.at())
		}

		// TODO function-definition - variadic-parameter
	}

	return parameters, nil
}

func (parser *Parser) parseAnonymousFunctionCreationExpression() (ast.IExpression, phpError.Error) {
	// -------------------------------------- anonymous-function-creation-expression -------------------------------------- MARK: anonymous-function-creation-expression

	// Spec: https://phplang.org/spec/10-expressions.html#anonymous-function-creation

	// anonymous-function-creation-expression:
	//    static(opt)   function   &(opt)   (   parameter-declaration-list(opt)   )   anonymous-function-use-clause(opt)   return-type(opt)   compound-statement

	// anonymous-function-use-clause:
	//    use   (   use-variable-name-list   )

	// use-variable-name-list:
	//    &(opt)   variable-name
	//    use-variable-name-list   ,   &(opt)   variable-name

	// Spec: https://phplang.org/spec/10-expressions.html#anonymous-function-creation
	// This operator returns an object of type Closure, or a derived type thereof, that encapsulates the anonymous function defined within.

	// TODO anonymous-function-creation-expression - static
	// Supported statement: anonymous function creation: `function ($param1) { ... }`
	PrintParserCallstack("anonymous-function-creation", parser)
	if !parser.isToken(lexer.KeywordToken, "function", false) {
		return ast.NewEmptyStmt(), NewExpectedError("function", parser.at())
	}

	pos := parser.eat().Position

	// TODO byRef: anonymous-function-creation-expression - &(opt)

	if !parser.isToken(lexer.OpOrPuncToken, "(", true) {
		return ast.NewEmptyStmt(), NewExpectedError("(", parser.at())
	}

	parameters, err := parser.parseFunctionParameters()
	if err != nil {
		return ast.NewEmptyStmt(), err
	}

	if !parser.isToken(lexer.OpOrPuncToken, ")", true) {
		return ast.NewEmptyStmt(), NewExpectedError(")", parser.at())
	}

	// TODO anonymous-function-creation-expression - anonymous-function-use-clause(opt)

	returnTypes := []string{}
	if parser.isToken(lexer.OpOrPuncToken, ":", true) {
		for parser.at().TokenType == lexer.KeywordToken && common.IsReturnTypeKeyword(parser.at().Value) {
			returnTypes = append(returnTypes, strings.ToLower(parser.eat().Value))
			if parser.isToken(lexer.OpOrPuncToken, "|", true) {
				continue
			}
			break
		}
	}

	body, err := parser.parseStmt()
	if err != nil {
		return ast.NewEmptyStmt(), err
	}
	if body.GetKind() != ast.CompoundStmt {
		return ast.NewEmptyStmt(), phpError.NewParseError("Expected compound statement. Got %s", body.GetKind())
	}

	return ast.NewAnonymousFunctionCreationExpr(parser.nextId(), pos, parameters, body.(*ast.CompoundStatement), returnTypes), nil
}

func (parser *Parser) parseExpr() (ast.IExpression, phpError.Error) {
	// Spec: https://phplang.org/spec/10-expressions.html#grammar-expression

	// expression:
	//    logical-inc-OR-expression-2
	//    include-expression
	//    include-once-expression
	//    require-expression
	//    require-once-expression

	// -------------------------------------- include-expression -------------------------------------- MARK: include-expression

	// Spec: https://phplang.org/spec/10-expressions.html#grammar-include-expression

	// include-expression:
	//    include   expression

	// Supported expression: include expression: `include 'lib.php';`
	if parser.isToken(lexer.KeywordToken, "include", false) {
		PrintParserCallstack("include-expression", parser)
		pos := parser.eat().Position

		expr, err := parser.parseExpr()
		if err != nil {
			return ast.NewEmptyExpr(), err
		}

		return ast.NewIncludeExpr(parser.nextId(), pos, expr), nil
	}

	// -------------------------------------- include-once-expression -------------------------------------- MARK: include-once-expression

	// Spec: https://phplang.org/spec/10-expressions.html#grammar-include-once-expression

	// include-once-expression:
	//    include_once   expression

	// Supported expression: include_once expression: `include_once 'lib.php';`
	if parser.isToken(lexer.KeywordToken, "include_once", false) {
		PrintParserCallstack("include-once-expression", parser)
		pos := parser.eat().Position

		expr, err := parser.parseExpr()
		if err != nil {
			return ast.NewEmptyExpr(), err
		}

		return ast.NewIncludeOnceExpr(parser.nextId(), pos, expr), nil
	}

	// -------------------------------------- require-expression -------------------------------------- MARK: require-expression

	// Spec: https://phplang.org/spec/10-expressions.html#grammar-require-expression

	// require-expression:
	//    require   expression

	// Supported expression: require expression: `require 'lib.php';`
	if parser.isToken(lexer.KeywordToken, "require", false) {
		PrintParserCallstack("require-expression", parser)
		pos := parser.eat().Position

		expr, err := parser.parseExpr()
		if err != nil {
			return ast.NewEmptyExpr(), err
		}

		return ast.NewRequireExpr(parser.nextId(), pos, expr), nil
	}

	// -------------------------------------- require-once-expression -------------------------------------- MARK: require-once-expression

	// Spec: https://phplang.org/spec/10-expressions.html#grammar-require-once-expression

	// require-once-expression:
	//    require_once   expression

	// Supported expression: require_once expression: `require_once 'lib.php';`
	if parser.isToken(lexer.KeywordToken, "require_once", false) {
		PrintParserCallstack("require-once-expression", parser)
		pos := parser.eat().Position

		expr, err := parser.parseExpr()
		if err != nil {
			return ast.NewEmptyExpr(), err
		}

		return ast.NewRequireOnceExpr(parser.nextId(), pos, expr), nil
	}

	return parser.parseLogicalIncOrExpr2()
}

func (parser *Parser) parseLogicalIncOrExpr2() (ast.IExpression, phpError.Error) {
	// Spec: https://phplang.org/spec/10-expressions.html#grammar-logical-inc-OR-expression-2

	// logical-inc-OR-expression-2:
	//    logical-exc-OR-expression
	//    logical-inc-OR-expression-2   or   logical-exc-OR-expression

	lhs, err := parser.parseLogicalExcOrExpr()
	if err != nil {
		return ast.NewEmptyExpr(), err
	}

	// Supported expression: logical inc or expression 2: `$var or 8;`
	for parser.isToken(lexer.KeywordToken, "or", true) {
		PrintParserCallstack("logical-inc-OR-expression-2", parser)
		rhs, err := parser.parseLogicalIncOrExpr2()
		if err != nil {
			return ast.NewEmptyExpr(), err
		}
		lhs = ast.NewLogicalExpr(parser.nextId(), lhs, "||", rhs)
	}
	return lhs, nil
}

func (parser *Parser) parseLogicalExcOrExpr() (ast.IExpression, phpError.Error) {
	// Spec: https://phplang.org/spec/10-expressions.html#grammar-logical-exc-OR-expression

	// logical-exc-OR-expression:
	//    logical-AND-expression-2
	//    logical-exc-OR-expression   xor   logical-AND-expression-2

	lhs, err := parser.parseLogicalAndExpr2()
	if err != nil {
		return ast.NewEmptyExpr(), err
	}

	// Supported expression: logical exc or expression: `$var xor 8;`
	for parser.isToken(lexer.KeywordToken, "xor", true) {
		PrintParserCallstack("logical-exc-OR-expression", parser)
		rhs, err := parser.parseLogicalExcOrExpr()
		if err != nil {
			return ast.NewEmptyExpr(), err
		}
		lhs = ast.NewLogicalExpr(parser.nextId(), lhs, "xor", rhs)
	}
	return lhs, nil
}

func (parser *Parser) parseLogicalAndExpr2() (ast.IExpression, phpError.Error) {
	// Spec: https://phplang.org/spec/10-expressions.html?#grammar-logical-AND-expression-2

	// logical-AND-expression-2:
	//    print-expression
	//    logical-AND-expression-2   and   yield-expression

	// print-expression

	// Spec: https://phplang.org/spec/10-expressions.html#grammar-print-expression

	// print-expression:
	//    yield-expression
	//    print   print-expression
	// Spec-Fix: So that by following assignment-expression the primary-expression is reachable
	//    assignment-expression

	var lhs ast.IExpression
	var err phpError.Error

	// Supported statement: print statement: `print "abc";`
	if parser.isToken(lexer.KeywordToken, "print", false) {
		PrintParserCallstack("print-expression", parser)
		pos := parser.eat().Position

		lhs, err = parser.parseLogicalAndExpr2()
		if err != nil {
			return ast.NewEmptyExpr(), err
		}
		lhs = ast.NewPrintExpr(parser.nextId(), pos, lhs)
	} else {
		lhs, err = parser.parseAssignmentExpr()
		if err != nil {
			return ast.NewEmptyExpr(), err
		}
	}

	// Supported expression: logical and expression 2: `$var and 8;`
	for parser.isToken(lexer.KeywordToken, "and", true) {
		PrintParserCallstack("logical-AND-expression-2", parser)
		rhs, err := parser.parseLogicalAndExpr2()
		if err != nil {
			return ast.NewEmptyExpr(), err
		}
		lhs = ast.NewLogicalExpr(parser.nextId(), lhs, "&&", rhs)
	}
	return lhs, nil

	// TODO yield-expression
}

func (parser *Parser) parseAssignmentExpr() (ast.IExpression, phpError.Error) {
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
	expr, err := parser.parseCoalesceExpr()
	if err != nil {
		return ast.NewEmptyExpr(), err
	}

	// Supported expression: conditional expression: `$var ? $a : "b";`
	// conditional-expression   ?   expression(opt)   :   coalesce-expression
	for parser.isToken(lexer.OpOrPuncToken, "?", true) {
		PrintParserCallstack("conditional-expression", parser)
		var ifExpr ast.IExpression = nil
		if !parser.isToken(lexer.OpOrPuncToken, ":", false) {
			ifExpr, err = parser.parseExpr()
			if err != nil {
				return ast.NewEmptyExpr(), err
			}
		}
		if err := parser.expect(lexer.OpOrPuncToken, ":", true); err != nil {
			return ast.NewEmptyExpr(), err
		}
		elseExpr, err := parser.parseExpr()
		if err != nil {
			return ast.NewEmptyExpr(), err
		}
		expr = ast.NewConditionalExpr(parser.nextId(), expr, ifExpr, elseExpr)
	}

	// Spec: https://phplang.org/spec/10-expressions.html#grammar-simple-assignment-expression

	// simple-assignment-expression:
	//    variable   =   assignment-expression
	//    list-intrinsic   =   assignment-expression

	// Supported expression: simple assignment expression: `$v = "abc";`
	// TODO simple-assignment-expression - list-intrinsic
	if ast.IsVariableExpr(expr) && parser.isToken(lexer.OpOrPuncToken, "=", true) {
		PrintParserCallstack("simple-assignment-expression", parser)
		value, err := parser.parseAssignmentExpr()
		if err != nil {
			return ast.NewEmptyExpr(), err
		}
		if expr.GetKind() == ast.SimpleVariableExpr && expr.(*ast.SimpleVariableExpression).VariableName.GetKind() == ast.VariableNameExpr &&
			expr.(*ast.SimpleVariableExpression).VariableName.(*ast.VariableNameExpression).VariableName == "$this" {
			return ast.NewEmptyExpr(), phpError.NewError("Cannot re-assign $this in %s", expr.GetPosString())
		}
		return ast.NewSimpleAssignmentExpr(parser.nextId(), expr, value), nil
	}

	// Spec: https://phplang.org/spec/10-expressions.html#grammar-compound-assignment-expression

	// compound-assignment-expression:
	//    variable   compound-assignment-operator   assignment-expression

	// compound-assignment-operator: one of
	//    **=   *=   /=   %=   +=   -=   .=   <<=   >>=   &=   ^=   |=

	// Supported expression: compound assignment expression: `$v += 2; $w &= 8;`
	if ast.IsVariableExpr(expr) &&
		parser.isTokenType(lexer.OpOrPuncToken, false) && common.IsCompoundAssignmentOp(parser.at().Value) {
		PrintParserCallstack("compound-assignment-operator", parser)
		operatorStr := strings.ReplaceAll(parser.eat().Value, "=", "")
		value, err := parser.parseAssignmentExpr()
		if err != nil {
			return ast.NewEmptyExpr(), nil
		}
		return ast.NewCompoundAssignmentExpr(parser.nextId(), expr, operatorStr, value), nil
	}

	return expr, nil
}

func (parser *Parser) parseCoalesceExpr() (ast.IExpression, phpError.Error) {
	// Spec: https://phplang.org/spec/10-expressions.html#grammar-coalesce-expression

	// coalesce-expression:
	//    logical-inc-OR-expression-1
	//    logical-inc-OR-expression-1   ??   coalesce-expression

	// logical-inc-OR-expression-1
	expr, err := parser.parseLogicalIncOrExpr1()
	if err != nil {
		return ast.NewEmptyExpr(), err
	}

	// Supported expression: coalesce expression: `$var ?? "b";`
	// logical-inc-OR-expression-1   ??   coalesce-expression
	if parser.isToken(lexer.OpOrPuncToken, "??", true) {
		PrintParserCallstack("coalesce-expression", parser)
		elseExpr, err := parser.parseCoalesceExpr()
		if err != nil {
			return ast.NewEmptyExpr(), err
		}

		return ast.NewCoalesceExpr(parser.nextId(), expr, elseExpr), nil
	}

	return expr, nil
}

func (parser *Parser) parseLogicalIncOrExpr1() (ast.IExpression, phpError.Error) {
	// Spec: https://phplang.org/spec/10-expressions.html#grammar-logical-inc-OR-expression-1

	// logical-inc-OR-expression-1:
	//    logical-AND-expression-1
	//    logical-inc-OR-expression-1   ||   logical-AND-expression-1

	lhs, err := parser.parseLogicalAndExpr1()
	if err != nil {
		return ast.NewEmptyExpr(), err
	}

	// Supported expression: logical inc or expression: `$var || 8;`
	for parser.isToken(lexer.OpOrPuncToken, "||", true) {
		PrintParserCallstack("logical-inc-OR-expression-1", parser)
		rhs, err := parser.parseLogicalIncOrExpr1()
		if err != nil {
			return ast.NewEmptyExpr(), err
		}
		lhs = ast.NewLogicalExpr(parser.nextId(), lhs, "||", rhs)
	}
	return lhs, nil
}

func (parser *Parser) parseLogicalAndExpr1() (ast.IExpression, phpError.Error) {
	// Spec: https://phplang.org/spec/10-expressions.html#grammar-logical-AND-expression-1

	// logical-AND-expression-1:
	//    bitwise-inc-OR-expression
	//    logical-AND-expression-1   &&   bitwise-inc-OR-expression

	lhs, err := parser.parseBitwiseIncOrExpr()
	if err != nil {
		return ast.NewEmptyExpr(), err
	}

	// Supported expression: logical and expression: `$var && 8;`
	for parser.isToken(lexer.OpOrPuncToken, "&&", true) {
		PrintParserCallstack("logical-AND-expression-1", parser)
		rhs, err := parser.parseLogicalAndExpr1()
		if err != nil {
			return ast.NewEmptyExpr(), err
		}
		lhs = ast.NewLogicalExpr(parser.nextId(), lhs, "&&", rhs)
	}
	return lhs, nil
}

func (parser *Parser) parseBitwiseIncOrExpr() (ast.IExpression, phpError.Error) {
	// Spec: https://phplang.org/spec/10-expressions.html#grammar-bitwise-inc-OR-expression

	// bitwise-inc-OR-expression:
	//    bitwise-exc-OR-expression
	//    bitwise-inc-OR-expression   |   bitwise-exc-OR-expression

	lhs, err := parser.parseBitwiseExcOrExpr()
	if err != nil {
		return ast.NewEmptyExpr(), err
	}

	// Supported expression: bitwise inc or expression: `$var | 8;`
	for parser.isToken(lexer.OpOrPuncToken, "|", true) {
		PrintParserCallstack("bitwise-inc-OR-expression", parser)
		rhs, err := parser.parseBitwiseIncOrExpr()
		if err != nil {
			return ast.NewEmptyExpr(), err
		}
		lhs = ast.NewBinaryOpExpr(parser.nextId(), lhs, "|", rhs)
	}
	return lhs, nil
}

func (parser *Parser) parseBitwiseExcOrExpr() (ast.IExpression, phpError.Error) {
	// Spec: https://phplang.org/spec/10-expressions.html#grammar-bitwise-exc-OR-expression

	// bitwise-exc-OR-expression:
	//    bitwise-AND-expression
	//    bitwise-exc-OR-expression   ^   bitwise-AND-expression

	lhs, err := parser.parseBitwiseAndExpr()
	if err != nil {
		return ast.NewEmptyExpr(), err
	}

	// Supported expression: bitwise exc or expression: `$var ^ 8;`
	for parser.isToken(lexer.OpOrPuncToken, "^", true) {
		PrintParserCallstack("bitwise-exc-OR-expression", parser)
		rhs, err := parser.parseBitwiseExcOrExpr()
		if err != nil {
			return ast.NewEmptyExpr(), err
		}
		lhs = ast.NewBinaryOpExpr(parser.nextId(), lhs, "^", rhs)
	}
	return lhs, nil
}

func (parser *Parser) parseBitwiseAndExpr() (ast.IExpression, phpError.Error) {
	// Spec: https://phplang.org/spec/10-expressions.html#grammar-bitwise-AND-expression

	// bitwise-AND-expression:
	//    equality-expression
	//    bitwise-AND-expression   &   equality-expression

	lhs, err := parser.parseEqualityExpr()
	if err != nil {
		return ast.NewEmptyExpr(), err
	}

	// Supported expression: bitwise and expression: `$var & 8;`
	for parser.isToken(lexer.OpOrPuncToken, "&", true) {
		PrintParserCallstack("bitwise-AND-expression", parser)
		rhs, err := parser.parseBitwiseAndExpr()
		if err != nil {
			return ast.NewEmptyExpr(), err
		}
		lhs = ast.NewBinaryOpExpr(parser.nextId(), lhs, "&", rhs)
	}
	return lhs, nil
}

func (parser *Parser) parseEqualityExpr() (ast.IExpression, phpError.Error) {
	// Spec: https://phplang.org/spec/10-expressions.html#grammar-equality-expression

	// equality-expression:
	//    relational-expression
	//    equality-expression   ==   relational-expression
	//    equality-expression   !=   relational-expression
	//    equality-expression   <>   relational-expression
	//    equality-expression   ===   relational-expression
	//    equality-expression   !==   relational-expression

	lhs, err := parser.parserRelationalExpr()
	if err != nil {
		return ast.NewEmptyExpr(), err
	}

	// Supported expression: equality expression: `$var === 42;`
	for parser.isTokenType(lexer.OpOrPuncToken, false) && common.IsEqualityOp(parser.at().Value) {
		PrintParserCallstack("equality-expression", parser)
		operator := parser.eat().Value
		rhs, err := parser.parserRelationalExpr()
		if err != nil {
			return ast.NewEmptyExpr(), err
		}
		lhs = ast.NewEqualityExpr(parser.nextId(), lhs, operator, rhs)
	}
	return lhs, nil
}

func (parser *Parser) parserRelationalExpr() (ast.IExpression, phpError.Error) {
	// Spec: https://phplang.org/spec/10-expressions.html#grammar-relational-expression

	// relational-expression:
	//    shift-expression
	//    relational-expression   <   shift-expression
	//    relational-expression   >   shift-expression
	//    relational-expression   <=   shift-expression
	//    relational-expression   >=   shift-expression
	//    relational-expression   <=>   shift-expression

	lhs, err := parser.parseShiftExpr()
	if err != nil {
		return ast.NewEmptyExpr(), err
	}

	// Supported expression: relational expression: `$var >= 42;`
	for parser.isTokenType(lexer.OpOrPuncToken, false) && common.IsRelationalExpressionOp(parser.at().Value) {
		PrintParserCallstack("relational-expression", parser)
		operator := parser.eat().Value
		rhs, err := parser.parseShiftExpr()
		if err != nil {
			return ast.NewEmptyExpr(), err
		}
		lhs = ast.NewRelationalExpr(parser.nextId(), lhs, operator, rhs)
	}
	return lhs, nil
}

func (parser *Parser) parseShiftExpr() (ast.IExpression, phpError.Error) {
	// Spec: https://phplang.org/spec/10-expressions.html#grammar-shift-expression

	// shift-expression:
	//    additive-expression
	//    shift-expression   <<   additive-expression
	//    shift-expression   >>   additive-expression

	lhs, err := parser.parseAdditiveExpr()
	if err != nil {
		return ast.NewEmptyExpr(), err
	}

	// Supported expression: shift expression: `$var << 8;`
	for parser.isTokenType(lexer.OpOrPuncToken, false) && slices.Contains([]string{"<<", ">>"}, parser.at().Value) {
		PrintParserCallstack("shift-expression", parser)
		operator := parser.eat().Value
		rhs, err := parser.parseShiftExpr()
		if err != nil {
			return ast.NewEmptyExpr(), err
		}
		lhs = ast.NewBinaryOpExpr(parser.nextId(), lhs, operator, rhs)
	}
	return lhs, nil
}

func (parser *Parser) parseAdditiveExpr() (ast.IExpression, phpError.Error) {
	// Spec: https://phplang.org/spec/10-expressions.html#grammar-additive-expression

	// additive-expression:
	//    multiplicative-expression
	//    additive-expression   +   multiplicative-expression
	//    additive-expression   -   multiplicative-expression
	//    additive-expression   .   multiplicative-expression

	lhs, err := parser.parseMultiplicativeExpr()
	if err != nil {
		return ast.NewEmptyExpr(), err
	}

	// Supported expression: additive expression: `$var + 42; $var - 42; "a" . "b";`
	for parser.isTokenType(lexer.OpOrPuncToken, false) && common.IsAdditiveOp(parser.at().Value) {
		PrintParserCallstack("additive-expression", parser)
		operator := parser.eat().Value
		rhs, err := parser.parseMultiplicativeExpr()
		if err != nil {
			return ast.NewEmptyExpr(), err
		}
		lhs = ast.NewBinaryOpExpr(parser.nextId(), lhs, operator, rhs)
	}

	return lhs, nil
}

func (parser *Parser) parseMultiplicativeExpr() (ast.IExpression, phpError.Error) {
	// Spec: https://phplang.org/spec/10-expressions.html#grammar-multiplicative-expression

	// multiplicative-expression:
	//    logical-NOT-expression
	//    multiplicative-expression   *   logical-NOT-expression
	//    multiplicative-expression   /   logical-NOT-expression
	//    multiplicative-expression   %   logical-NOT-expression

	lhs, err := parser.parseLogicalNotExpr()
	if err != nil {
		return ast.NewEmptyExpr(), err
	}

	// Supported expression: multiplicative expression: `$var * 42; $var / 42; $var % 42;`
	for parser.isTokenType(lexer.OpOrPuncToken, false) && common.IsMultiplicativeOp(parser.at().Value) {
		PrintParserCallstack("multiplicative-expression", parser)
		operator := parser.eat().Value
		rhs, err := parser.parseLogicalNotExpr()
		if err != nil {
			return ast.NewEmptyExpr(), err
		}
		lhs = ast.NewBinaryOpExpr(parser.nextId(), lhs, operator, rhs)
	}
	return lhs, nil
}

func (parser *Parser) parseLogicalNotExpr() (ast.IExpression, phpError.Error) {
	// Spec: https://phplang.org/spec/10-expressions.html#grammar-logical-NOT-expression

	// logical-NOT-expression:
	//    instanceof-expression
	//    !   instanceof-expression

	// Supported expression: logical not expression: `!$var;`
	isNotExpression := parser.isToken(lexer.OpOrPuncToken, "!", false)
	var pos *position.Position = nil
	if isNotExpression {
		PrintParserCallstack("logical-NOT-expression", parser)
		pos = parser.eat().Position
	}

	expr, err := parser.parseInstanceofExpr()
	if err != nil {
		return ast.NewEmptyExpr(), err
	}

	if isNotExpression {
		return ast.NewLogicalNotExpr(parser.nextId(), pos, expr), nil
	}
	return expr, nil
}

func (parser *Parser) parseInstanceofExpr() (ast.IExpression, phpError.Error) {
	// Spec: https://phplang.org/spec/10-expressions.html#grammar-instanceof-expression

	// instanceof-expression:
	//    unary-expression
	//    instanceof-subject   instanceof   class-type-designator

	// instanceof-subject:
	//    instanceof-expression

	// TODO instanceof-expression

	return parser.parseUnaryExpr()
}

func (parser *Parser) parseUnaryExpr() (ast.IExpression, phpError.Error) {
	// Spec: https://phplang.org/spec/10-expressions.html#grammar-unary-expression

	// unary-expression:
	//    exponentiation-expression
	//    unary-op-expression
	//    error-control-expression
	//    cast-expression

	// These operators associate right-to-left.

	// unary-op-expression

	// Spec: https://phplang.org/spec/10-expressions.html#grammar-unary-op-expression

	// unary-op-expression:
	//    unary-operator   unary-expression

	// unary-operator: one of
	//    +   -   ~

	// Supported expression: unary expression: `-1; +1; ~1;`
	if parser.isTokenType(lexer.OpOrPuncToken, false) && common.IsUnaryOp(parser.at().Value) {
		PrintParserCallstack("unary-op-expression", parser)
		pos := parser.at().Position
		operator := parser.eat().Value
		expr, err := parser.parseUnaryExpr()
		if err != nil {
			return ast.NewEmptyExpr(), err
		}

		return ast.NewUnaryOpExpr(parser.nextId(), pos, operator, expr), nil
	}

	// -------------------------------------- error-control-expression -------------------------------------- MARK: error-control-expression

	// Spec: https://phplang.org/spec/10-expressions.html#error-control-operator

	// error-control-expression:
	//    @   unary-expression

	// Supported expression: error control expression: `@func();`
	if parser.isToken(lexer.OpOrPuncToken, "@", false) {
		PrintParserCallstack("error-control-expression", parser)
		pos := parser.eat().Position
		expr, err := parser.parseUnaryExpr()
		if err != nil {
			return ast.NewEmptyExpr(), err
		}
		return ast.NewErrorControlExpr(parser.nextId(), pos, expr), nil
	}

	// cast-expression

	// Spec: https://phplang.org/spec/10-expressions.html#grammar-cast-expression

	// cast-expression:
	//    (   cast-type   )   unary-expression

	// Supported expression: cast expression: `(int)$a;(string)$a;`
	if parser.isToken(lexer.OpOrPuncToken, "(", false) &&
		parser.next(0).TokenType == lexer.KeywordToken && common.IsCastTypeKeyword(parser.next(0).Value) {
		PrintParserCallstack("cast-expression", parser)
		pos := parser.eat().Position
		castType := parser.eat().Value
		if !parser.isToken(lexer.OpOrPuncToken, ")", true) {
			return ast.NewEmptyExpr(), NewExpectedError(")", parser.at())
		}
		expr, err := parser.parseUnaryExpr()
		if err != nil {
			return ast.NewEmptyExpr(), err
		}
		return ast.NewCastExpr(parser.nextId(), pos, castType, expr), nil
	}

	return parser.parseExponentiationExpr()
}

func (parser *Parser) parseExponentiationExpr() (ast.IExpression, phpError.Error) {
	// Spec: https://phplang.org/spec/10-expressions.html#grammar-exponentiation-expression

	// exponentiation-expression:
	//    clone-expression
	//    clone-expression   **   exponentiation-expression

	lhs, err := parser.parseCloneExpr()
	if err != nil {
		return ast.NewEmptyExpr(), err
	}

	// Supported expression: exponentiation expression: `$var ** 42;`
	if parser.isToken(lexer.OpOrPuncToken, "**", true) {
		PrintParserCallstack("exponentiation-expression", parser)
		rhs, err := parser.parseExponentiationExpr()
		if err != nil {
			return ast.NewEmptyExpr(), err
		}
		return ast.NewBinaryOpExpr(parser.nextId(), lhs, "**", rhs), nil
	}
	return lhs, nil
}

func (parser *Parser) parseCloneExpr() (ast.IExpression, phpError.Error) {
	// Spec: https://phplang.org/spec/10-expressions.html#grammar-clone-expression

	// clone-expression:
	//	primary-expression
	//	clone   primary-expression

	// TODO clone-expression
	return parser.parsePrimaryExpr()
}

func (parser *Parser) parsePrimaryExpr() (ast.IExpression, phpError.Error) {
	PrintParserCallstack("primary-expression", parser)

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

	// -------------------------------------- variable -------------------------------------- MARK: variable

	// Spec: https://phplang.org/spec/10-expressions.html#grammar-variable

	// variable:
	//    callable-variable
	//    scoped-property-access-expression
	//    member-access-expression

	var variable ast.IExpression

	// -------------------------------------- callable-variable -------------------------------------- MARK: callable-variable

	// Spec: https://phplang.org/spec/10-expressions.html#grammar-callable-variable

	// callable-variable:
	//    simple-variable
	//    subscript-expression
	//    member-call-expression
	//    scoped-call-expression
	//    function-call-expression

	// -------------------------------------- simple-variable -------------------------------------- MARK: simple-variable

	// Spec: https://phplang.org/spec/10-expressions.html#simple-variable

	// simple-variable:
	//    variable-name

	if parser.isTokenType(lexer.VariableNameToken, false) {
		PrintParserCallstack("simple-variable", parser)
		variable = ast.NewSimpleVariableExpr(parser.nextId(), ast.NewVariableNameExpr(parser.nextId(), parser.at().Position, parser.eat().Value))
	}

	// Spec: https://phplang.org/spec/10-expressions.html#simple-variable

	// simple-variable:
	//    $   {   expression   }

	// Supported expression: variable access: `echo $v;`
	if parser.isToken(lexer.OpOrPuncToken, "$", false) &&
		parser.next(0).TokenType == lexer.OpOrPuncToken && parser.next(0).Value == "{" {
		PrintParserCallstack("simple-variable", parser)
		parser.eatN(2)
		// Get expression
		expr, err := parser.parseExpr()
		if err != nil {
			return ast.NewEmptyExpr(), err
		}

		if parser.at().Value == "}" {
			parser.eat()
			variable = ast.NewSimpleVariableExpr(parser.nextId(), expr)
		} else {
			return ast.NewEmptyExpr(), phpError.NewParseError("End of simple variable expression not detected in %s", parser.at().GetPosString())
		}
	}

	// Spec: https://phplang.org/spec/10-expressions.html#simple-variable

	// simple-variable:
	//    $   simple-variable

	if parser.isToken(lexer.OpOrPuncToken, "$", true) {
		PrintParserCallstack("simple-variable", parser)
		if expr, err := parser.parsePrimaryExpr(); err != nil {
			return ast.NewEmptyExpr(), err
		} else {
			variable = ast.NewSimpleVariableExpr(parser.nextId(), expr)
		}
	}

	if ast.IsVariableExpr(variable) &&
		!parser.isToken(lexer.OpOrPuncToken, "(", false) && !parser.isToken(lexer.OpOrPuncToken, "[", false) &&
		!parser.isToken(lexer.OpOrPuncToken, "{", false) && !parser.isToken(lexer.OpOrPuncToken, "++", false) &&
		!parser.isToken(lexer.OpOrPuncToken, "--", false) && !parser.isToken(lexer.OpOrPuncToken, "->", false) {
		return variable, nil
	}

	// TODO The following expressions can occur multiple times: $a[0]()["abc"]()()...

	// -------------------------------------- subscript-expression -------------------------------------- MARK: subscript-expression

	// Spec: https://phplang.org/spec/10-expressions.html#grammar-subscript-expression

	// subscript-expression:
	//    dereferencable-expression   [   expression(opt)   ]
	//    dereferencable-expression   {   expression   }   <b>[Deprecated form]</b>

	// dereferencable-expression:
	//    variable
	//    (   expression   )
	//    array-creation-expression
	//    string-literal

	// Supported expression: subscript expression: `$a[1];`
	if ast.IsVariableExpr(variable) && parser.isToken(lexer.OpOrPuncToken, "[", false) {
		PrintParserCallstack("subscript-expression", parser)
		for ast.IsVariableExpr(variable) && parser.isToken(lexer.OpOrPuncToken, "[", true) {
			var err phpError.Error
			var index ast.IExpression
			if !parser.isToken(lexer.OpOrPuncToken, "]", false) {
				index, err = parser.parseExpr()
				if err != nil {
					return ast.NewEmptyExpr(), err
				}
			}
			if !parser.isToken(lexer.OpOrPuncToken, "]", true) {
				return ast.NewEmptyExpr(), NewExpectedError("]", parser.at())
			}
			variable = ast.NewSubscriptExpr(parser.nextId(), variable, index)
		}
		return variable, nil
	}

	// TODO member-call-expression
	// TODO scoped-call-expression

	// -------------------------------------- function-call-expression -------------------------------------- MARK: function-call-expression

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

	// Supported expression: function call expression: `func(42);`
	if (parser.isTokenType(lexer.NameToken, false) && parser.next(0).TokenType == lexer.OpOrPuncToken && parser.next(0).Value == "(") ||
		(ast.IsVariableExpr(variable) && parser.isToken(lexer.OpOrPuncToken, "(", false)) {
		PrintParserCallstack("function-call-expression", parser)

		var pos *position.Position
		var functionName ast.IExpression

		if parser.isTokenType(lexer.NameToken, false) && parser.next(0).TokenType == lexer.OpOrPuncToken && parser.next(0).Value == "(" {
			pos = parser.at().Position
			functionName = ast.NewStringLiteralExpr(parser.nextId(), pos, parser.eat().Value, ast.SingleQuotedString)
		} else {
			pos = variable.GetPosition()
			functionName = variable
		}

		args := []ast.IExpression{}
		// Eat opening parentheses
		parser.eat()
		for {
			if parser.isToken(lexer.OpOrPuncToken, ")", true) {
				break
			}

			arg, err := parser.parseExpr()
			if err != nil {
				return ast.NewEmptyExpr(), err
			}
			args = append(args, arg)

			if parser.isToken(lexer.OpOrPuncToken, ",", true) || parser.isToken(lexer.OpOrPuncToken, ")", false) {
				continue
			}
			return ast.NewEmptyExpr(), phpError.NewParseError("Expected \",\" or \")\". Got: %s", parser.at())
		}
		return ast.NewFunctionCallExpr(parser.nextId(), pos, functionName, args), nil
	}

	// TODO scoped-property-access-expression

	// -------------------------------------- member-access-expression -------------------------------------- MARK: member-access-expression

	// Spec: https://phplang.org/spec/10-expressions.html#member-access-operator

	// member-access-expression:
	//    dereferencable-expression   ->   member-name

	// member-name:
	//    name
	//    simple-variable
	//    {   expression   }

	// TODO member-access-expression - check if it is a "dereferencable-expression"

	// Supported expression: member access expression: `$obj->member`
	if ast.IsVariableExpr(variable) && parser.isToken(lexer.OpOrPuncToken, "->", false) {
		PrintParserCallstack("member-access-expression", parser)

		pos := parser.eat().Position

		member, err := parser.parseExpr()
		if err != nil {
			return ast.NewEmptyExpr(), err
		}

		// TODO member-access-expression - check member name type

		return ast.NewMemberAccessExpr(parser.nextId(), pos, variable, member), nil
	}

	// TODO class-constant-access-expression

	// literal
	if parser.isTokenType(lexer.IntegerLiteralToken, false) || parser.isTokenType(lexer.FloatingLiteralToken, false) ||
		parser.isTokenType(lexer.StringLiteralToken, false) {
		return parser.parseLiteral()
	}

	// array-creation-expression
	if (parser.isToken(lexer.KeywordToken, "array", false) &&
		parser.next(0).TokenType == lexer.OpOrPuncToken && parser.next(0).Value == "(") ||
		parser.isToken(lexer.OpOrPuncToken, "[", false) {
		return parser.parseArrayCreationExpr()
	}

	// intrinsic
	if parser.isToken(lexer.KeywordToken, "empty", false) || parser.isToken(lexer.KeywordToken, "eval", false) ||
		parser.isToken(lexer.KeywordToken, "exit", false) || parser.isToken(lexer.KeywordToken, "die", false) ||
		parser.isToken(lexer.KeywordToken, "isset", false) {
		return parser.parseIntrinsic()
	}

	// object-creation-expression
	if parser.isToken(lexer.KeywordToken, "new", false) {
		return parser.parseObjectCreationExpression()
	}

	// anonymous-function-creation-expression
	if parser.isToken(lexer.KeywordToken, "function", false) {
		return parser.parseAnonymousFunctionCreationExpression()
	}

	// -------------------------------------- postfix-increment-expression -------------------------------------- MARK: postfix-increment-expression

	// Spec: https://phplang.org/spec/10-expressions.html#grammar-postfix-increment-expression

	// postfix-increment-expression:
	//    variable   ++

	// Spec: https://phplang.org/spec/10-expressions.html#grammar-postfix-decrement-expression

	// postfix-decrement-expression:
	//    variable   --

	// Supported expression: postfix (in/de)crease expression: `$var++; $var--;`
	if ast.IsVariableExpr(variable) && (parser.isToken(lexer.OpOrPuncToken, "++", false) ||
		parser.isToken(lexer.OpOrPuncToken, "--", false)) {
		PrintParserCallstack("postfix-decrement-expression", parser)
		return ast.NewPostfixIncExpr(parser.nextId(), parser.at().Position, variable, parser.eat().Value), nil
	}

	// -------------------------------------- prefix-increment-expression -------------------------------------- MARK: prefix-increment-expression

	// Spec: https://phplang.org/spec/10-expressions.html#grammar-prefix-increment-expression

	// prefix-increment-expression:
	//    ++   variable

	// Spec: https://phplang.org/spec/10-expressions.html#grammar-prefix-decrement-expression
	// prefix-decrement-expression:
	//    --   variable

	// Supported expression: prefix (in/de)crease expression: `++$var; --$var;`
	if parser.isToken(lexer.OpOrPuncToken, "++", false) ||
		parser.isToken(lexer.OpOrPuncToken, "--", false) {
		PrintParserCallstack("prefix-decrement-expression", parser)
		pos := parser.at().Position
		operator := parser.eat().Value
		variable, err := parser.parsePrimaryExpr()
		if err != nil {
			return ast.NewEmptyExpr(), err
		}
		if !ast.IsVariableExpr(variable) {
			return ast.NewEmptyExpr(), phpError.NewParseError("Syntax error, unexpected %s", variable)
		}
		return ast.NewPrefixIncExpr(parser.nextId(), pos, variable, operator), nil
	}

	// TODO byref-assignment-expression
	// TODO shell-command-expression

	// -------------------------------------- (   expression   ) -------------------------------------- MARK: (   expression   )

	// Supported expression: parenthesized expression: `(1 + 2) * 3;`
	if parser.isToken(lexer.OpOrPuncToken, "(", false) {
		PrintParserCallstack("parenthesized-expression", parser)
		pos := parser.eat().Position
		expr, err := parser.parseExpr()
		if err != nil {
			return ast.NewEmptyExpr(), err
		}
		if parser.isToken(lexer.OpOrPuncToken, ")", true) {
			return ast.NewParenthesizedExpr(parser.nextId(), pos, expr), nil
		} else {
			return ast.NewEmptyExpr(), NewExpectedError(")", parser.at())
		}
	}

	// -------------------------------------- constant-access-expression -------------------------------------- MARK: constant-access-expression

	// Spec: https://phplang.org/spec/10-expressions.html#grammar-constant-access-expression

	// constant-access-expression:
	//    qualified-name

	// A constant-access-expression evaluates to the value of the constant with name qualified-name.

	// Spec: https://phplang.org/spec/09-lexical-structure.html#grammar-qualified-name

	// qualified-name::
	//    namespace-name-as-a-prefix(opt)   name

	if parser.isTokenType(lexer.NameToken, false) || parser.isTokenType(lexer.KeywordToken, false) {
		// TODO constant-access-expression - namespace-name-as-a-prefix
		// TODO constant-access-expression - check if name is a defined constant here or in interpreter
		PrintParserCallstack("constant-access-expression", parser)
		constantName := parser.at().Value
		// TODO Find a way to reduce "is..constant" to just one time
		if common.IsCorePredefinedConstant(constantName) || common.IsContextDependentConstant(constantName) {
			constantName = strings.ToUpper(constantName)
		}
		return ast.NewConstantAccessExpr(parser.nextId(), parser.eat().Position, constantName), nil
	}

	return ast.NewEmptyExpr(), phpError.NewParseError("Unsupported expression type '%s', value: '%s' in %s", parser.at().TokenType, parser.at().Value, parser.at().GetPosString())
}

func (parser *Parser) parseLiteral() (ast.IExpression, phpError.Error) {
	// -------------------------------------- literal -------------------------------------- MARK: literal

	// Spec: https://phplang.org/spec/10-expressions.html#grammar-literal

	// literal:
	//    integer-literal
	//    floating-literal
	//    string-literal

	// A literal evaluates to its value, as specified in the lexical specification for literals.

	// integer-literal
	if parser.isTokenType(lexer.IntegerLiteralToken, false) {
		PrintParserCallstack("integer-literal", parser)
		intValue, err := common.IntegerLiteralToInt64(parser.at().Value, false)
		if err != nil {
			return ast.NewEmptyExpr(), phpError.NewParseError("Unsupported integer literal \"%s\"", parser.at().Value)
		}

		return ast.NewIntegerLiteralExpr(parser.nextId(), parser.eat().Position, intValue), nil
	}

	// floating-literal
	if parser.isTokenType(lexer.FloatingLiteralToken, false) {
		PrintParserCallstack("floating-literal", parser)
		if common.IsFloatingLiteral(parser.at().Value, false) {
			return ast.NewFloatingLiteralExpr(parser.nextId(), parser.at().Position, common.FloatingLiteralToFloat64(parser.eat().Value, false)), nil
		}

		return ast.NewEmptyExpr(), phpError.NewParseError("Unsupported floating literal \"%s\"", parser.at().Value)
	}

	// string-literal
	if parser.isTokenType(lexer.StringLiteralToken, false) {
		PrintParserCallstack("string-literal", parser)
		// single-quoted-string-literal
		if common.IsSingleQuotedStringLiteral(parser.at().Value) {
			return ast.NewStringLiteralExpr(
					parser.nextId(), parser.at().Position, common.SingleQuotedStringLiteralToString(parser.eat().Value), ast.SingleQuotedString),
				nil
		}

		// double-quoted-string-literal
		if common.IsDoubleQuotedStringLiteral(parser.at().Value) {
			return ast.NewStringLiteralExpr(
					parser.nextId(), parser.at().Position, common.DoubleQuotedStringLiteralToString(parser.eat().Value), ast.DoubleQuotedString),
				nil
		}

		// heredoc-string-literal
		if common.IsHeredocStringLiteral(parser.at().Value) {
			return ast.NewStringLiteralExpr(
					parser.nextId(), parser.at().Position, common.HeredocStringLiteralToString(parser.eat().Value), ast.HeredocString),
				nil
		}

		// TODO nowdoc-string-literal
	}

	return ast.NewEmptyExpr(), phpError.NewParseError("parseLiteral: Unsupported literal: '%s' in %s", parser.at().Value, parser.at().GetPosString())
}

func (parser *Parser) parseArrayCreationExpr() (ast.IExpression, phpError.Error) {
	// -------------------------------------- array-creation-expression -------------------------------------- MARK: array-creation-expression

	// Spec: https://phplang.org/spec/10-expressions.html#grammar-array-creation-expression

	// array-creation-expression:
	//    array   (   array-initializer(opt)   )
	//    [   array-initializer(opt)   ]

	// array-initializer:
	//    array-initializer-list   ,(opt)

	// array-initializer-list:
	//    array-element-initializer
	//    array-element-initializer   ,   array-initializer-list

	// array-element-initializer:
	//    &(opt)   element-value
	//    element-key   =>   &(opt)   element-value

	// element-key:
	//    expression

	// element-value:
	//    expression

	PrintParserCallstack("array-creation-expression", parser)

	if !((parser.isToken(lexer.KeywordToken, "array", false) &&
		parser.next(0).TokenType == lexer.OpOrPuncToken && parser.next(0).Value == "(") ||
		parser.isToken(lexer.OpOrPuncToken, "[", false)) {
		return ast.NewEmptyExpr(), phpError.NewParseError("Unsupported array creation: %s", parser.at())
	}

	isShortSyntax := true
	var pos *position.Position
	if parser.isToken(lexer.KeywordToken, "array", false) &&
		parser.next(0).TokenType == lexer.OpOrPuncToken && parser.next(0).Value == "(" {
		pos = parser.eat().Position
		parser.eat()
		isShortSyntax = false
	} else {
		pos = parser.eat().Position
	}

	arrayExpr := ast.NewArrayLiteralExpr(parser.nextId(), pos)
	for {
		if (!isShortSyntax && parser.isToken(lexer.OpOrPuncToken, ")", true)) ||
			(isShortSyntax && parser.isToken(lexer.OpOrPuncToken, "]", true)) {
			break
		}

		// TODO byRef: &(opt)   element-value
		// TODO byRef: element-key   =>   &(opt)   element-value

		keyOrValue, err := parser.parseExpr()
		var value ast.IExpression
		if err != nil {
			return ast.NewEmptyExpr(), err
		}

		if parser.isToken(lexer.OpOrPuncToken, "=>", true) {
			value, err = parser.parseExpr()
			if err != nil {
				return ast.NewEmptyExpr(), err
			}
		}

		if value == nil {
			arrayExpr.AddElement(nil, keyOrValue)
		} else {
			arrayExpr.AddElement(keyOrValue, value)
		}

		if parser.isToken(lexer.OpOrPuncToken, ",", true) ||
			(!isShortSyntax && parser.isToken(lexer.OpOrPuncToken, ")", false)) ||
			(isShortSyntax && parser.isToken(lexer.OpOrPuncToken, "]", false)) {
			continue
		}
		if isShortSyntax {
			return ast.NewEmptyExpr(), phpError.NewParseError("Expected \",\" or \"]\". Got: %s", parser.at())
		} else {
			return ast.NewEmptyExpr(), phpError.NewParseError("Expected \",\" or \")\". Got: %s", parser.at())
		}
	}
	return arrayExpr, nil
}

func (parser *Parser) parseIntrinsic() (ast.IExpression, phpError.Error) {
	// -------------------------------------- intrinsic -------------------------------------- MARK: intrinsic

	// Spec: https://phplang.org/spec/10-expressions.html#grammar-intrinsic

	// intrinsic:
	//    empty-intrinsic
	//    eval-intrinsic
	//    exit-intrinsic
	//    isset-intrinsic

	// -------------------------------------- empty-intrinsic -------------------------------------- MARK: empty-intrinsic

	// Spec: https://phplang.org/spec/10-expressions.html#grammar-empty-intrinsic

	// empty-intrinsic:
	//    empty   (   expression   )

	// Supported intrinsic: empty intrinsic: `empty($v);`
	if parser.isToken(lexer.KeywordToken, "empty", false) {
		PrintParserCallstack("empty-intrinsic", parser)
		pos := parser.eat().Position
		if !parser.isToken(lexer.OpOrPuncToken, "(", true) {
			return ast.NewEmptyExpr(), NewExpectedError("(", parser.at())
		}
		expr, err := parser.parseExpr()
		if err != nil {
			return ast.NewEmptyExpr(), err
		}
		if !parser.isToken(lexer.OpOrPuncToken, ")", true) {
			return ast.NewEmptyExpr(), NewExpectedError(")", parser.at())
		}
		return ast.NewEmptyIntrinsic(parser.nextId(), pos, expr), nil
	}

	// -------------------------------------- eval-intrinsic -------------------------------------- MARK: eval-intrinsic

	// Spec: https://phplang.org/spec/10-expressions.html#grammar-eval-intrinsic

	// eval-intrinsic:
	//    eval   (   expression   )

	// Supported intrinsic: eval intrinsic: `eval($v);`
	if parser.isToken(lexer.KeywordToken, "eval", false) {
		PrintParserCallstack("eval-intrinsic", parser)
		pos := parser.eat().Position
		if !parser.isToken(lexer.OpOrPuncToken, "(", true) {
			return ast.NewEmptyExpr(), NewExpectedError("(", parser.at())
		}
		expr, err := parser.parseExpr()
		if err != nil {
			return ast.NewEmptyExpr(), err
		}
		if !parser.isToken(lexer.OpOrPuncToken, ")", true) {
			return ast.NewEmptyExpr(), NewExpectedError(")", parser.at())
		}
		return ast.NewEvalIntrinsic(parser.nextId(), pos, expr), nil
	}

	// -------------------------------------- exit-intrinsic -------------------------------------- MARK: exit-intrinsic

	// Spec: https://phplang.org/spec/10-expressions.html#grammar-exit-intrinsic

	// exit-intrinsic:
	//    exit
	//    exit   (   expression(opt)   )
	//    die
	//    die   (   expression(opt)   )

	// Supported intrinsic: exit intrinsic: `exit(0);`
	// Supported intrinsic: die intrinsic: `die(0);`
	if parser.isToken(lexer.KeywordToken, "exit", false) || parser.isToken(lexer.KeywordToken, "die", false) {
		PrintParserCallstack("exit-intrinsic", parser)
		pos := parser.eat().Position
		var expr ast.IExpression = nil
		if parser.isToken(lexer.OpOrPuncToken, "(", true) {
			if !parser.isToken(lexer.OpOrPuncToken, ")", false) {
				var err phpError.Error
				expr, err = parser.parseExpr()
				if err != nil {
					return ast.NewEmptyExpr(), err
				}
			}
			if !parser.isToken(lexer.OpOrPuncToken, ")", true) {
				return ast.NewEmptyExpr(), NewExpectedError(")", parser.at())
			}
		}
		return ast.NewExitIntrinsic(parser.nextId(), pos, expr), nil
	}

	// -------------------------------------- isset-intrinsic -------------------------------------- MARK: isset-intrinsic

	// Spec: https://phplang.org/spec/10-expressions.html#grammar-isset-intrinsic

	// isset-intrinsic:
	//    isset   (   variable-list   ,(opt)   )

	// variable-list:
	//    variable
	//    variable-list   ,   variable

	// Supported intrinsic: isset intrinsic: `isset($v);`
	if parser.isToken(lexer.KeywordToken, "isset", false) {
		PrintParserCallstack("isset-intrinsic", parser)
		pos := parser.eat().Position
		if !parser.isToken(lexer.OpOrPuncToken, "(", true) {
			return ast.NewEmptyExpr(), NewExpectedError("(", parser.at())
		}
		args := []ast.IExpression{}
		for {
			if len(args) > 0 && parser.isToken(lexer.OpOrPuncToken, ")", true) {
				break
			}

			arg, err := parser.parseExpr()
			if err != nil {
				return ast.NewEmptyExpr(), err
			}
			if !ast.IsVariableExpr(arg) {
				return ast.NewEmptyExpr(), phpError.NewParseError("Fatal error: Cannot use isset() on the result of an expression")
			}
			args = append(args, arg)

			if parser.isToken(lexer.OpOrPuncToken, ",", true) ||
				parser.isToken(lexer.OpOrPuncToken, ")", false) {
				continue
			}
			return ast.NewEmptyExpr(), phpError.NewParseError("Expected \",\" or \")\". Got: %s", parser.at())
		}
		return ast.NewIssetIntrinsic(parser.nextId(), pos, args), nil
	}

	return ast.NewEmptyExpr(), phpError.NewParseError("Unsupported intrinsic: %s", parser.at())
}
