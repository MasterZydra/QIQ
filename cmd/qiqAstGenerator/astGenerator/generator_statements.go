package astGenerator

import (
	"QIQ/cmd/qiq/ast"
)

// ProcessBreakStmt implements ast.Visitor.
func (generator *AstGenerator) ProcessBreakStmt(stmt *ast.BreakStatement, _ any) (any, error) {
	panic("ProcessBreakStmt is unimplemented")
}

// ProcessClassDeclarationStmt implements ast.Visitor.
func (generator *AstGenerator) ProcessClassDeclarationStmt(stmt *ast.ClassDeclarationStatement, _ any) (any, error) {
	variableName := toGoVarName(stmt.Name)
	// Create new class declaration stmt
	generator.println(`%s := ast.NewClassDeclarationStmt(0, nil, "%s", %s, %s)`, variableName, stmt.Name, toBoolStr(stmt.IsAbstract), toBoolStr(stmt.IsFinal))

	for _, interfaceName := range stmt.Interfaces {
		generator.println(`%s.Interfaces = append(%s.Interfaces, "%s")`, variableName, variableName, interfaceName)
	}

	for _, propertyName := range stmt.PropertieNames {
		property := stmt.Properties[propertyName]
		generator.print(`%s.AddProperty(ast.NewPropertyDeclarationStmt(0, nil, "%s", "%s", %s, %s, `, variableName, property.Name, property.Visibility, toBoolStr(property.IsStatic), toStringSlice(property.Type))
		generator.processStmt(property.InitialValue)
		generator.println("))")
	}

	for _, methodName := range stmt.MethodNames {
		method, _ := stmt.GetMethod(methodName)
		generator.print(`%s.AddMethod(`, variableName)
		generator.print(`ast.NewMethodDefinitionStmt(0, nil, "%s", %s, %s, `,
			method.Name, toStringSlice(method.Modifiers), funcParamArrayToStr(method.Params),
		)
		generator.processStmt(method.Body)
		generator.println(`, %s))`, toStringSlice(method.ReturnType))
	}

	generator.println("")
	return nil, nil
}

// ProcessCompoundStmt implements ast.Visitor.
func (generator *AstGenerator) ProcessCompoundStmt(stmt *ast.CompoundStatement, _ any) (any, error) {
	generator.print("ast.NewCompoundStmt(0, []ast.IStatement{")
	for i, statement := range stmt.Statements {
		if i > 0 {
			generator.print(", ")
		}
		generator.processStmt(statement)
	}
	generator.print("})")
	return nil, nil
}

// ProcessConstDeclarationStmt implements ast.Visitor.
func (generator *AstGenerator) ProcessConstDeclarationStmt(stmt *ast.ConstDeclarationStatement, _ any) (any, error) {
	panic("ProcessConstDeclarationStmt is unimplemented")
}

// ProcessContinueStmt implements ast.Visitor.
func (generator *AstGenerator) ProcessContinueStmt(stmt *ast.ContinueStatement, _ any) (any, error) {
	panic("ProcessContinueStmt is unimplemented")
}

// ProcessDeclareStmt implements ast.Visitor.
func (generator *AstGenerator) ProcessDeclareStmt(stmt *ast.DeclareStatement, _ any) (any, error) {
	panic("ProcessDeclareStmt is unimplemented")
}

// ProcessDoStmt implements ast.Visitor.
func (generator *AstGenerator) ProcessDoStmt(stmt *ast.DoStatement, _ any) (any, error) {
	panic("ProcessDoStmt is unimplemented")
}

// ProcessEchoStmt implements ast.Visitor.
func (generator *AstGenerator) ProcessEchoStmt(stmt *ast.EchoStatement, _ any) (any, error) {
	panic("ProcessEchoStmt is unimplemented")
}

// ProcessExpressionStmt implements ast.Visitor.
func (generator *AstGenerator) ProcessExpressionStmt(stmt *ast.ExpressionStatement, _ any) (any, error) {
	generator.print("ast.NewExpressionStmt(0, ")
	generator.processStmt(stmt.Expr)
	generator.print(")")
	return nil, nil
}

// ProcessForStmt implements ast.Visitor.
func (generator *AstGenerator) ProcessForStmt(stmt *ast.ForStatement, _ any) (any, error) {
	panic("ProcessForStmt is unimplemented")
}

// ProcessForeachStmt implements ast.Visitor.
func (generator *AstGenerator) ProcessForeachStmt(stmt *ast.ForeachStatement, _ any) (any, error) {
	panic("ProcessForeachStmt is unimplemented")
}

// ProcessFunctionDefinitionStmt implements ast.Visitor.
func (generator *AstGenerator) ProcessFunctionDefinitionStmt(stmt *ast.FunctionDefinitionStatement, _ any) (any, error) {
	panic("ProcessFunctionDefinitionStmt is unimplemented")
}

// ProcessGlobalDeclarationStmt implements ast.Visitor.
func (generator *AstGenerator) ProcessGlobalDeclarationStmt(stmt *ast.GlobalDeclarationStatement, _ any) (any, error) {
	panic("ProcessGlobalDeclarationStmt is unimplemented")
}

// ProcessIfStmt implements ast.Visitor.
func (generator *AstGenerator) ProcessIfStmt(stmt *ast.IfStatement, _ any) (any, error) {
	panic("ProcessIfStmt is unimplemented")
}

// ProcessInterfaceDeclarationStmt implements ast.Visitor.
func (generator *AstGenerator) ProcessInterfaceDeclarationStmt(stmt *ast.InterfaceDeclarationStatement, _ any) (any, error) {
	variableName := toGoVarName(stmt.Name)
	// Create new interface declaration stmt
	generator.println(`%s := ast.NewInterfaceDeclarationStmt(0, nil, "%s")`, variableName, stmt.Name)

	// Add all parents
	for _, parent := range stmt.Parents {
		generator.println(`%s.Parents = append(%s.Parents, "%s")`, variableName, variableName, parent)
	}

	// Add all methods
	for _, methodName := range stmt.MethodNames {
		methodDecl, _ := stmt.GetMethod(methodName)
		generator.println(
			`%s.AddMethod(ast.NewMethodDefinitionStmt(0, nil, "%s", %s, %s, nil, %s))`,
			variableName, methodDecl.Name, toStringSlice(methodDecl.Modifiers), funcParamArrayToStr(methodDecl.Params), toStringSlice(methodDecl.ReturnType),
		)
	}

	generator.println("")
	return nil, nil
}

// ProcessReturnStmt implements ast.Visitor.
func (generator *AstGenerator) ProcessReturnStmt(stmt *ast.ReturnStatement, _ any) (any, error) {
	generator.print("ast.NewReturnStmt(0, nil, ")
	generator.processStmt(stmt.Expr)
	generator.print(")")
	return nil, nil
}

// ProcessStmt implements ast.Visitor.
func (generator *AstGenerator) ProcessStmt(stmt *ast.Statement, _ any) (any, error) {
	panic("ProcessStmt is unimplemented")
}

// ProcessThrowStmt implements ast.Visitor.
func (generator *AstGenerator) ProcessThrowStmt(stmt *ast.ThrowStatement, _ any) (any, error) {
	panic("ProcessThrowStmt is unimplemented")
}

// ProcessTryStmt implements ast.Visitor.
func (generator *AstGenerator) ProcessTryStmt(stmt *ast.TryStatement, _ any) (any, error) {
	panic("ProcessTryStmt is unimplemented")
}

// ProcessWhileStmt implements ast.Visitor.
func (generator *AstGenerator) ProcessWhileStmt(stmt *ast.WhileStatement, _ any) (any, error) {
	panic("ProcessWhileStmt is unimplemented")
}
