package parser

import (
	"GoPHP/cmd/goPHP/ast"
	"GoPHP/cmd/goPHP/common"
	"GoPHP/cmd/goPHP/lexer"
	"GoPHP/cmd/goPHP/phpError"
	"GoPHP/cmd/goPHP/position"
)

func (parser *Parser) parseClassDeclaration() (ast.IStatement, phpError.Error) {
	// -------------------------------------- class-declaration -------------------------------------- MARK: class-declaration

	// Spec: https://phplang.org/spec/14-classes.html#class-declarations

	// class-declaration:
	//    class-modifier(opt)   class   name   class-base-clause(opt)   class-interface-clause(opt)   {   class-member-declarations(opt)   }

	// class-modifier:
	//    abstract
	//    final

	// class-base-clause:
	//    extends   qualified-name

	// class-interface-clause:
	//    implements   qualified-name
	//    class-interface-clause   ,   qualified-name

	PrintParserCallstack("class-declaration", parser)

	// class-modifier
	isAbstract := parser.isToken(lexer.KeywordToken, "abstract", true)
	isFinal := parser.isToken(lexer.KeywordToken, "final", true)

	class := ast.NewClassDeclarationStmt(parser.nextId(), parser.eat().Position, isAbstract, isFinal)

	// class name
	class.Name = parser.at().Value
	classNamePos := parser.eat().Position
	if !common.IsName(class.Name) {
		return ast.NewEmptyExpr(), phpError.NewParseError("\"%s\" is not a valid class name at %s", class.Name, classNamePos.ToPosString())
	}

	// class-base-clause
	if parser.isToken(lexer.KeywordToken, "extends", true) {
		class.BaseClass = parser.at().Value
		baseClassPos := parser.eat().Position
		if !common.IsQualifiedName(class.Name) {
			return ast.NewEmptyExpr(), phpError.NewParseError("\"%s\" is not a valid class name at %s", class.Name, baseClassPos.ToPosString())
		}
	}

	// class-interface-clause
	if parser.isToken(lexer.KeywordToken, "implements", true) {
		for {
			interfaceName := parser.at().Value
			interfaceNamePos := parser.eat().Position
			if !common.IsQualifiedName(interfaceName) {
				return ast.NewEmptyExpr(), phpError.NewParseError("\"%s\" is not a valid interface name at %s", class.Name, interfaceNamePos.ToPosString())
			}

			class.Interfaces = append(class.Interfaces, interfaceName)

			if parser.isToken(lexer.OpOrPuncToken, ",", true) {
				continue
			}
			break
		}
	}

	if !parser.isToken(lexer.OpOrPuncToken, "{", true) {
		return ast.NewEmptyStmt(), phpError.NewParseError("Expected \"{\". Got %s", parser.at())
	}

	if err := parser.parseClassMemberDeclaration(class); err != nil {
		return ast.NewEmptyStmt(), err
	}

	if !parser.isToken(lexer.OpOrPuncToken, "}", true) {
		return ast.NewEmptyStmt(), phpError.NewParseError("Expected \"{\". Got %s", parser.at())
	}

	return class, nil
}

func (parser *Parser) parseClassMemberDeclaration(class *ast.ClassDeclarationStatement) phpError.Error {
	// Spec: https://phplang.org/spec/14-classes.html#grammar-class-member-declarations

	// class-member-declarations:
	//    class-member-declaration
	//    class-member-declarations   class-member-declaration

	// class-member-declaration:
	//    class-const-declaration
	//    property-declaration
	//    method-declaration
	//    constructor-declaration
	//    destructor-declaration
	//    trait-use-clause

	PrintParserCallstack("class-member-declarations", parser)

	for {
		// trait-use-clause
		if parser.isToken(lexer.KeywordToken, "use", false) {
			if err := parser.parserTraitUseClause(class); err != nil {
				return err
			}
			continue
		}

		// class-const-declaration
		if (parser.isTokenType(lexer.KeywordToken, false) && common.IsVisibilitModifierKeyword(parser.at().Value) &&
			parser.next(0).TokenType == lexer.KeywordToken && parser.next(0).Value == "const") ||
			parser.isToken(lexer.KeywordToken, "const", false) {
			if err := parser.parseClassConstDeclaration(class); err != nil {
				return err
			}
			continue
		}

		// TODO property-declaration

		// TODO method-declaration

		// constructor-declaration
		isConstructorDeclaration, err := parser.parseClassConstrutorDeclaration(class)
		if isConstructorDeclaration && err != nil {
			return err
		}
		if isConstructorDeclaration {
			continue
		}

		// TODO destructor-declaration

		// End of class-member-declarations
		if parser.isToken(lexer.OpOrPuncToken, "}", false) {
			return nil
		}

		return phpError.NewParseError("parseClassMemberDeclaration: Unexpected token: %s", parser.at())
	}
}

func (parser *Parser) parseClassConstDeclaration(class *ast.ClassDeclarationStatement) phpError.Error {
	// Spec: https://phplang.org/spec/14-classes.html#grammar-class-const-declaration

	// class-const-declaration:
	//    visibility-modifier(opt)   const   const-elements   ;

	// const-elements:
	//    const-element
	//    const-elements   ,   const-element

	// const-element:
	//    name   =   constant-expression

	// Spec: https://phplang.org/spec/14-classes.html#constants
	// If visibility-modifier for a class constant is omitted, public is assumed. The visibility-modifier applies to all constants defined in the const-elements list.
	visibility := "public"

	if parser.isTokenType(lexer.KeywordToken, false) && common.IsVisibilitModifierKeyword(parser.at().Value) {
		visibility = parser.eat().Value
	}

	PrintParserCallstack("class-const-statement", parser)
	pos := parser.eat().Position
	if err := parser.expectTokenType(lexer.NameToken, false); err != nil {
		return err
	}
	for {
		name := parser.eat().Value
		if err := parser.expect(lexer.OpOrPuncToken, "=", true); err != nil {
			return err
		}
		// TODO parse constant-expression
		value, err := parser.parseExpr()
		if err != nil {
			return err
		}

		class.AddConst(ast.NewClassConstDeclarationStmt(parser.nextId(), pos, name, value, visibility))
		if parser.isToken(lexer.OpOrPuncToken, ",", true) {
			continue
		}
		if parser.isToken(lexer.OpOrPuncToken, ";", true) {
			return nil
		}
		return phpError.NewParseError("Class const declaration - unexpected token %s", parser.at())
	}
}

func (parser *Parser) parserTraitUseClause(class *ast.ClassDeclarationStatement) phpError.Error {
	// Spec: https://phplang.org/spec/16-traits.html#grammar-trait-use-clause

	// trait-use-clause:
	//    use   trait-name-list   trait-use-specification

	// trait-name-list:
	//    qualified-name
	//    trait-name-list   ,   qualified-name

	// trait-use-specification:
	//    ;
	//    {   trait-select-and-alias-clauses(opt)   }

	// trait-select-and-alias-clauses:
	//    trait-select-and-alias-clause
	//    trait-select-and-alias-clauses   trait-select-and-alias-clause

	// trait-select-and-alias-clause:
	//    trait-select-insteadof-clause   ;
	//    trait-alias-as-clause   ;

	// trait-select-insteadof-clause:
	//    qualified-name   ::   name   insteadof   trait-name-list

	// trait-alias-as-clause:
	//    name   as   visibility-modifier(opt)   name
	//    name   as   visibility-modifier   name(opt)

	PrintParserCallstack("trait-use-clause", parser)

	// Eat "use"
	parser.eat()

	for {
		// trait-name-list
		traitName := parser.at().Value
		traitNamePos := parser.eat().Position
		if !common.IsQualifiedName(traitName) {
			return phpError.NewParseError("\"%s\" is not a valid trait name at %s", traitName, traitNamePos.ToPosString())
		}
		class.AddTrait(ast.NewTraitUseStmt(parser.nextId(), traitNamePos, traitName))

		if parser.isToken(lexer.OpOrPuncToken, ",", true) {
			continue
		}

		// trait-use-specification
		if !parser.isToken(lexer.OpOrPuncToken, ";", true) {
			return phpError.NewParseError("parserTraitUseClause: Expected \";\". Got: %s", parser.at())
		}
		// TODO trait-select-and-alias-clauses(opt)
		return nil
	}
}

func (parser *Parser) parseClassConstrutorDeclaration(class *ast.ClassDeclarationStatement) (bool, phpError.Error) {
	// Spec: https://phplang.org/spec/14-classes.html#grammar-constructor-declaration

	// constructor-declaration:
	//    method-modifiers   function   &(opt)   __construct   (   parameter-declaration-list(opt)   )   compound-statement

	PrintParserCallstack("constructor-declaration", parser)

	// Check if the following tokens result in a valid constructor definition
	isConstructor := true
	offset := -1

	visibilityModifierKeyword := ""
	classModifierKeyword := ""
	staticModifierKeyword := ""
	var staticModifierKeywordPos *position.Position = nil
	for {
		token := parser.next(offset)
		// Only allow one visibility modifier keyword
		if visibilityModifierKeyword == "" &&
			token.TokenType == lexer.KeywordToken &&
			common.IsVisibilitModifierKeyword(token.Value) {
			visibilityModifierKeyword = token.Value
			offset++
			continue
		}
		// Allow static modifier even if it will return an error later
		if staticModifierKeyword == "" &&
			token.TokenType == lexer.KeywordToken &&
			token.Value == "static" {
			staticModifierKeyword = token.Value
			staticModifierKeywordPos = token.Position
			offset++
			continue
		}
		// Only allow one class modifier keyword
		if classModifierKeyword == "" &&
			token.TokenType == lexer.KeywordToken &&
			common.IsClassModifierKeyword(token.Value) {
			classModifierKeyword = token.Value
			offset++
			continue
		}

		// TODO &(opt)
		// Check if it is a function with the name "__construct"
		if token.TokenType == lexer.KeywordToken &&
			token.Value == "function" &&
			parser.next(offset+1).TokenType == lexer.NameToken &&
			parser.next(offset+1).Value == "__construct" {
			offset++
			break
		}

		isConstructor = false
		break
	}

	// Return if itis not a constructor declaration
	if !isConstructor {
		return isConstructor, nil
	}

	// Static modifier is not allowed for constructor
	if staticModifierKeyword != "" {
		return isConstructor, phpError.NewError(
			"Method %s::__construct cannot be static in %s",
			class.Name, staticModifierKeywordPos.ToPosString(),
		)
	}

	// Eat all tokens to get the name token "__construct"
	parser.eatN(offset + 1)

	// Store position of "__construct"
	pos := parser.eat().Position

	// Fallback to visibility modifier "public"
	if visibilityModifierKeyword == "" {
		visibilityModifierKeyword = "public"
	}

	// Build modifiers list
	modifiers := []string{visibilityModifierKeyword}
	if classModifierKeyword != "" {
		modifiers = append(modifiers, classModifierKeyword)
	}

	// (   parameter-declaration-list(opt)   )
	if !parser.isToken(lexer.OpOrPuncToken, "(", true) {
		return isConstructor, phpError.NewParseError("Expected \"(\". Got %s", parser.at())
	}
	parameters, err := parser.parseFunctionParameters()
	if err != nil {
		return isConstructor, err
	}

	if !parser.isToken(lexer.OpOrPuncToken, ")", true) {
		return isConstructor, phpError.NewParseError("Expected \")\". Got %s", parser.at())
	}

	// compound-statement
	body, err := parser.parseStmt()
	if err != nil {
		return isConstructor, err
	}
	if body.GetKind() != ast.CompoundStmt {
		return isConstructor, phpError.NewParseError("Expected compound statement. Got %s", body.GetKind())
	}

	class.AddMethod(ast.NewMethodDefinitionStmt(
		parser.nextId(), pos,
		"__construct", modifiers, parameters, body.(*ast.CompoundStatement), []string{"self"},
	))

	return isConstructor, nil
}
