package parser

import (
	"GoPHP/cmd/goPHP/ast"
	"GoPHP/cmd/goPHP/common"
	"GoPHP/cmd/goPHP/lexer"
	"GoPHP/cmd/goPHP/phpError"
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

	for {
		// class-const-declaration
		if (parser.isTokenType(lexer.KeywordToken, false) && common.IsVisibilitModifierKeyword(parser.at().Value) &&
			parser.next(0).TokenType == lexer.KeywordToken && parser.next(0).Value == "const") ||
			parser.isToken(lexer.KeywordToken, "const", false) {
			parser.parseClassConstDeclaration(class)
			continue
		}

		// TODO property-declaration

		// TODO method-declaration

		// TODO constructor-declaration

		// TODO destructor-declaration

		// TODO trait-use-clause

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
