package parser

import (
	"GoPHP/cmd/goPHP/ast"
	"GoPHP/cmd/goPHP/common"
	"GoPHP/cmd/goPHP/lexer"
	"GoPHP/cmd/goPHP/phpError"
	"GoPHP/cmd/goPHP/position"
)

func (parser *Parser) parseObjectCreationExpression() (ast.IExpression, phpError.Error) {
	// -------------------------------------- object-creation-expression -------------------------------------- MARK: object-creation-expression

	// Spec: https://phplang.org/spec/10-expressions.html#grammar-object-creation-expression

	// object-creation-expression:
	//    new   class-type-designator   (   argument-expression-list(opt)   )
	//    new   class-type-designator   (   argument-expression-list   ,(opt)   )
	//    new   class-type-designator
	//    new   class   (   argument-expression-list(opt)   )   class-base-clause(opt)   class-interface-clause(opt)   {   class-member-declarations(opt)   }
	//    new   class   class-base-clause(opt)   class-interface-clause(opt)   {   class-member-declarations(opt)   }

	// class-type-designator:
	//    qualified-name
	//    new-variable

	// new-variable:
	//    simple-variable
	//    new-variable   [   expression(opt)   ]
	//    new-variable   {   expression   }
	//    new-variable   ->   member-name
	//    qualified-name   ::   simple-variable
	//    relative-scope   ::   simple-variable
	//    new-variable   ::   simple-variable

	// TODO object-creation-expression - variants

	// Supported expression: object creation expression: `new myClass;`
	if !parser.isToken(lexer.KeywordToken, "new", false) {
		return ast.NewEmptyExpr(), phpError.NewParseError("Expected keyword \"new\". Got %s", parser.at())
	}

	pos := parser.eat().Position

	if !parser.isTokenType(lexer.NameToken, false) {
		return ast.NewEmptyExpr(), phpError.NewParseError("parseObjectCreationExpression: Only qualified name as designator allowed")
	}
	if !common.IsQualifiedName(parser.at().Value) {
		return ast.NewEmptyExpr(), phpError.NewParseError("parseObjectCreationExpression: Only qualified name as designator allowed")
	}
	designator := parser.eat().Value

	hasParenthese := parser.isToken(lexer.OpOrPuncToken, "(", true)

	args := []ast.IExpression{}
	if hasParenthese {
		for {
			if parser.isToken(lexer.OpOrPuncToken, ")", false) {
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
	}

	if hasParenthese && !parser.isToken(lexer.OpOrPuncToken, ")", true) {
		return ast.NewEmptyExpr(), NewExpectedError(")", parser.at())
	}

	return ast.NewObjectCreationExpr(parser.nextId(), pos, designator, args), nil
}

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

	// Supported statement: class declaration: `class MyClass extends ParentC implements I, J {}`
	PrintParserCallstack("class-declaration", parser)

	// class-modifier
	isAbstract := parser.isToken(lexer.KeywordToken, "abstract", true)
	isFinal := parser.isToken(lexer.KeywordToken, "final", true)

	pos := parser.eat().Position

	// class name
	className := parser.at().Value
	classNamePos := parser.eat().Position
	if !common.IsName(className) {
		return ast.NewEmptyExpr(), phpError.NewParseError("\"%s\" is not a valid class name at %s", className, classNamePos.ToPosString())
	}

	class := ast.NewClassDeclarationStmt(parser.nextId(), pos, className, isAbstract, isFinal)

	// class-base-clause
	if parser.isToken(lexer.KeywordToken, "extends", true) {
		class.BaseClass = parser.at().Value
		baseClassPos := parser.eat().Position
		if !common.IsQualifiedName(class.BaseClass) {
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
		return ast.NewEmptyStmt(), NewExpectedError("{", parser.at())
	}

	if err := parser.parseClassMemberDeclaration(class); err != nil {
		return ast.NewEmptyStmt(), err
	}

	if !parser.isToken(lexer.OpOrPuncToken, "}", true) {
		return ast.NewEmptyStmt(), NewExpectedError("}", parser.at())
	}

	return class, nil
}

func (parser *Parser) parseClassMemberDeclaration(class *ast.ClassDeclarationStatement) phpError.Error {
	// -------------------------------------- class-member-declarations -------------------------------------- MARK: class-member-declarations

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
		// End of class-member-declarations
		if parser.isToken(lexer.OpOrPuncToken, "}", false) {
			return nil
		}

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

		// constructor-declaration
		isConstructorDeclaration, err := parser.parseClassConstrutorDeclaration(class)
		if isConstructorDeclaration && err != nil {
			return err
		}
		if isConstructorDeclaration {
			continue
		}

		// destructor-declaration
		isDestructorDeclaration, err := parser.parseClassDestrutorDeclaration(class)
		if isDestructorDeclaration && err != nil {
			return err
		}
		if isDestructorDeclaration {
			continue
		}

		// method-declaration
		isMethodDeclaration, err := parser.parseClassMethodDeclaration(class)
		if isMethodDeclaration && err != nil {
			return err
		}
		if isMethodDeclaration {
			continue
		}

		// property-declaration
		isPropertyDeclaration, err := parser.parseClassPropertyDeclaration(class)
		if isPropertyDeclaration && err != nil {
			return err
		}
		if isPropertyDeclaration {
			continue
		}

		return phpError.NewParseError("parseClassMemberDeclaration: Unexpected token: %s", parser.at())
	}
}

func (parser *Parser) parseClassConstDeclaration(class *ast.ClassDeclarationStatement) phpError.Error {
	// -------------------------------------- class-const-declaration -------------------------------------- MARK: class-const-declaration

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
	// -------------------------------------- trait-use-clause -------------------------------------- MARK: trait-use-clause

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
			return NewExpectedError(";", parser.at())
		}
		// TODO trait-select-and-alias-clauses(opt)
		return nil
	}
}

func (parser *Parser) parseClassConstrutorDeclaration(class *ast.ClassDeclarationStatement) (bool, phpError.Error) {
	// -------------------------------------- constructor-declaration -------------------------------------- MARK: constructor-declaration

	// Spec: https://phplang.org/spec/14-classes.html#grammar-constructor-declaration

	// constructor-declaration:
	//    method-modifiers   function   &(opt)   __construct   (   parameter-declaration-list(opt)   )   compound-statement

	// Check if the following tokens result in a valid constructor definition
	isConstructor, offset, visibilityModifierKeyword, classModifierKeyword, staticModifierKeyword, staticModifierKeywordPos := parser.isMethod("__construct")

	// Return if it is not a constructor declaration
	if !isConstructor {
		return isConstructor, nil
	}

	PrintParserCallstack("constructor-declaration", parser)

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
		return isConstructor, NewExpectedError("(", parser.at())
	}
	parameters, err := parser.parseFunctionParameters()
	if err != nil {
		return isConstructor, err
	}

	if !parser.isToken(lexer.OpOrPuncToken, ")", true) {
		return isConstructor, NewExpectedError(")", parser.at())
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
		"__construct", modifiers, parameters, body.(*ast.CompoundStatement), []string{"void"},
	))

	return isConstructor, nil
}

func (parser *Parser) parseClassDestrutorDeclaration(class *ast.ClassDeclarationStatement) (bool, phpError.Error) {
	// -------------------------------------- destructor-declaration -------------------------------------- MARK: destructor-declaration

	// Spec: https://phplang.org/spec/14-classes.html#grammar-destructor-declaration

	// destructor-declaration:
	//    method-modifiers   function   &(opt)   __destruct   (   )   compound-statement

	// Check if the following tokens result in a valid destructor definition
	isDestructor, offset, visibilityModifierKeyword, classModifierKeyword, staticModifierKeyword, staticModifierKeywordPos := parser.isMethod("__destruct")

	// Return if it is not a destructor declaration
	if !isDestructor {
		return isDestructor, nil
	}

	PrintParserCallstack("destructor-declaration", parser)

	// Static modifier is not allowed for destructor
	if staticModifierKeyword != "" {
		return isDestructor, phpError.NewError(
			"Method %s::__destruct cannot be static in %s",
			class.Name, staticModifierKeywordPos.ToPosString(),
		)
	}

	// Eat all tokens to get the name token "__destruct"
	parser.eatN(offset + 1)

	// Store position of "__destruct"
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

	// (   )
	if !parser.isToken(lexer.OpOrPuncToken, "(", true) {
		return isDestructor, NewExpectedError("(", parser.at())
	}
	if !parser.isToken(lexer.OpOrPuncToken, ")", true) {
		return isDestructor, NewExpectedError(")", parser.at())
	}

	// compound-statement
	body, err := parser.parseStmt()
	if err != nil {
		return isDestructor, err
	}
	if body.GetKind() != ast.CompoundStmt {
		return isDestructor, phpError.NewParseError("Expected compound statement. Got %s", body.GetKind())
	}

	class.AddMethod(ast.NewMethodDefinitionStmt(
		parser.nextId(), pos,
		"__destruct", modifiers, []ast.FunctionParameter{}, body.(*ast.CompoundStatement), []string{"void"},
	))

	return isDestructor, nil
}

func (parser *Parser) parseClassMethodDeclaration(class *ast.ClassDeclarationStatement) (bool, phpError.Error) {
	// -------------------------------------- method-declaration -------------------------------------- MARK: method-declaration

	// Spec: https://phplang.org/spec/14-classes.html#grammar-method-declaration

	// method-declaration:
	//    method-modifiers(opt)   function-definition
	//    method-modifiers   function-definition-header   ;

	// method-modifiers:
	//    method-modifier
	//    method-modifiers   method-modifier

	// method-modifier:
	//    visibility-modifier
	//    static-modifier
	//    class-modifier

	// 	function-definition:
	//    function-definition-header   compound-statement

	// function-definition-header:
	//    function   &(opt)   name   (   parameter-declaration-list(opt)   )   return-type(opt)

	// Check if the following tokens result in a valid constructor definition
	isMethod, offset, visibilityModifierKeyword, classModifierKeyword, staticModifierKeyword, _ := parser.isMethod("")

	// Return if it is not a method declaration
	if !isMethod {
		return isMethod, nil
	}

	PrintParserCallstack("method-declaration", parser)

	// Eat all tokens to get the name token
	parser.eatN(offset + 1)

	// Store position of name token
	name := parser.at().Value
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
	if staticModifierKeyword != "" {
		modifiers = append(modifiers, staticModifierKeyword)
	}

	// (   parameter-declaration-list(opt)   )
	if !parser.isToken(lexer.OpOrPuncToken, "(", true) {
		return isMethod, NewExpectedError("(", parser.at())
	}
	parameters, err := parser.parseFunctionParameters()
	if err != nil {
		return isMethod, err
	}

	if !parser.isToken(lexer.OpOrPuncToken, ")", true) {
		return isMethod, NewExpectedError(")", parser.at())
	}

	// return-type
	var returnTypes []string = []string{"mixed"}
	if parser.isToken(lexer.OpOrPuncToken, ":", true) {
		returnTypes, err = parser.getTypes(true)
		if err != nil {
			return isMethod, err
		}
	}

	// compound-statement
	body, err := parser.parseStmt()
	if err != nil {
		return isMethod, err
	}
	if body.GetKind() != ast.CompoundStmt {
		return isMethod, phpError.NewParseError("Expected compound statement. Got %s", body.GetKind())
	}

	class.AddMethod(ast.NewMethodDefinitionStmt(
		parser.nextId(), pos,
		name, modifiers, parameters, body.(*ast.CompoundStatement), returnTypes,
	))

	return isMethod, nil
}

func (parser *Parser) parseClassPropertyDeclaration(class *ast.ClassDeclarationStatement) (bool, phpError.Error) {
	// -------------------------------------- property-declaration -------------------------------------- MARK: property-declaration

	// Spec: https://phplang.org/spec/14-classes.html#grammar-property-declaration

	// property-declaration:
	//    property-modifier   property-elements   ;

	// property-modifier:
	//    var   // deprecated
	//    visibility-modifier   static-modifier(opt)
	//    static-modifier   visibility-modifier(opt)

	// property-elements:
	//    property-element
	//    property-elements   property-element

	// property-element:
	//    variable-name   property-initializer(opt)   ;

	// property-initializer:
	//    =   constant-expression

	isProperty := false
	offset := -1
	visibilityModifierKeyword := ""
	staticModifierKeyword := ""
	propertyType := []string{"mixed"}

	token := func() *lexer.Token {
		return parser.next(offset)
	}

	step := "modifier"
	for {
		// Only allow one visibility modifier keyword
		if step == "modifier" && visibilityModifierKeyword == "" &&
			token().TokenType == lexer.KeywordToken && common.IsVisibilitModifierKeyword(token().Value) {
			visibilityModifierKeyword = token().Value
			offset++
			continue
		}

		// Allow static modifier even if it will return an error later
		if step == "modifier" && staticModifierKeyword == "" &&
			token().TokenType == lexer.KeywordToken && token().Value == "static" {
			staticModifierKeyword = token().Value
			offset++
			continue
		}

		// Property type
		if step == "modifier" && parser.isPhpType(token()) {
			var err phpError.Error
			propertyType, offset, err = parser.getTypesWithOffset(false, offset)
			if err != nil {
				return isProperty, err
			}
			step = "type"
		}

		// Check if given name is a valid variable name
		if token().TokenType == lexer.VariableNameToken && common.IsVariableName(token().Value) {
			isProperty = true
			break
		}

		break
	}

	// Return if it is not a method declaration
	if !isProperty {
		return isProperty, nil
	}

	PrintParserCallstack("property-declaration", parser)

	// Eat all tokens to get the name token
	parser.eatN(offset + 1)

	// property-element
	pos := parser.at().Position
	name := parser.eat().Value

	// property-initializer
	var initialValue ast.IExpression = nil
	if parser.isToken(lexer.OpOrPuncToken, "=", true) {
		var err phpError.Error
		initialValue, err = parser.parseExpr()
		if err != nil {
			return isProperty, err
		}
		// TODO Check if it is a constant-expression
	}

	// Fallback to visibility modifier "public"
	if visibilityModifierKeyword == "" {
		visibilityModifierKeyword = "public"
	}

	if !parser.isToken(lexer.OpOrPuncToken, ";", true) {
		return isProperty, NewExpectedError(";", parser.at())
	}

	class.AddProperty(ast.NewPropertyDeclarationStmt(parser.nextId(), pos, name, visibilityModifierKeyword, staticModifierKeyword != "", propertyType, initialValue))

	return isProperty, nil
}

func (parser *Parser) isMethod(name string) (
	isFunction bool,
	offset int,
	visibilityModifierKeyword, classModifierKeyword, staticModifierKeyword string,
	staticModifierKeywordPos *position.Position) {
	// -------------------------------------- isMethod -------------------------------------- MARK: isMethod
	offset = -1

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
		// Check if it is a function with the given name
		if token.TokenType == lexer.KeywordToken && token.Value == "function" &&
			((name == "" && parser.next(offset+1).TokenType == lexer.NameToken) ||
				(parser.next(offset+1).TokenType == lexer.NameToken &&
					parser.next(offset+1).Value == name)) {
			isFunction = true
			offset++
			return
		}

		return
	}
}
