package astGenerator

import (
	"QIQ/cmd/qiq/ast"
)

// ProcessBreakStmt implements ast.Visitor.
func (generator *AstGenerator) ProcessBreakStmt(stmt *ast.BreakStatement, context any) (any, error) {
	panic("ProcessBreakStmt is unimplemented")
}

// ProcessClassDeclarationStmt implements ast.Visitor.
func (generator *AstGenerator) ProcessClassDeclarationStmt(stmt *ast.ClassDeclarationStatement, context any) (any, error) {
	variableName := toGoVarName(stmt.Name)
	// Create new class declaration stmt
	generator.println(`%s := ast.NewClassDeclarationStmt(0, nil, "%s", %s, %s)`, variableName, stmt.Name, toBoolStr(stmt.IsAbstract), toBoolStr(stmt.IsFinal))

	generator.println("")
	return nil, nil
}

// ProcessCompoundStmt implements ast.Visitor.
func (generator *AstGenerator) ProcessCompoundStmt(stmt *ast.CompoundStatement, context any) (any, error) {
	panic("ProcessCompoundStmt is unimplemented")
}

// ProcessConstDeclarationStmt implements ast.Visitor.
func (generator *AstGenerator) ProcessConstDeclarationStmt(stmt *ast.ConstDeclarationStatement, context any) (any, error) {
	panic("ProcessConstDeclarationStmt is unimplemented")
}

// ProcessContinueStmt implements ast.Visitor.
func (generator *AstGenerator) ProcessContinueStmt(stmt *ast.ContinueStatement, context any) (any, error) {
	panic("ProcessContinueStmt is unimplemented")
}

// ProcessDeclareStmt implements ast.Visitor.
func (generator *AstGenerator) ProcessDeclareStmt(stmt *ast.DeclareStatement, context any) (any, error) {
	panic("ProcessDeclareStmt is unimplemented")
}

// ProcessDoStmt implements ast.Visitor.
func (generator *AstGenerator) ProcessDoStmt(stmt *ast.DoStatement, context any) (any, error) {
	panic("ProcessDoStmt is unimplemented")
}

// ProcessEchoStmt implements ast.Visitor.
func (generator *AstGenerator) ProcessEchoStmt(stmt *ast.EchoStatement, context any) (any, error) {
	panic("ProcessEchoStmt is unimplemented")
}

// ProcessExpressionStmt implements ast.Visitor.
func (generator *AstGenerator) ProcessExpressionStmt(stmt *ast.ExpressionStatement, context any) (any, error) {
	panic("ProcessExpressionStmt is unimplemented")
}

// ProcessForStmt implements ast.Visitor.
func (generator *AstGenerator) ProcessForStmt(stmt *ast.ForStatement, context any) (any, error) {
	panic("ProcessForStmt is unimplemented")
}

// ProcessForeachStmt implements ast.Visitor.
func (generator *AstGenerator) ProcessForeachStmt(stmt *ast.ForeachStatement, context any) (any, error) {
	panic("ProcessForeachStmt is unimplemented")
}

// ProcessFunctionDefinitionStmt implements ast.Visitor.
func (generator *AstGenerator) ProcessFunctionDefinitionStmt(stmt *ast.FunctionDefinitionStatement, context any) (any, error) {
	panic("ProcessFunctionDefinitionStmt is unimplemented")
}

// ProcessGlobalDeclarationStmt implements ast.Visitor.
func (generator *AstGenerator) ProcessGlobalDeclarationStmt(stmt *ast.GlobalDeclarationStatement, context any) (any, error) {
	panic("ProcessGlobalDeclarationStmt is unimplemented")
}

// ProcessIfStmt implements ast.Visitor.
func (generator *AstGenerator) ProcessIfStmt(stmt *ast.IfStatement, context any) (any, error) {
	panic("ProcessIfStmt is unimplemented")
}

// ProcessInterfaceDeclarationStmt implements ast.Visitor.
func (generator *AstGenerator) ProcessInterfaceDeclarationStmt(stmt *ast.InterfaceDeclarationStatement, context any) (any, error) {
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
func (generator *AstGenerator) ProcessReturnStmt(stmt *ast.ReturnStatement, context any) (any, error) {
	panic("ProcessReturnStmt is unimplemented")
}

// ProcessStmt implements ast.Visitor.
func (generator *AstGenerator) ProcessStmt(stmt *ast.Statement, context any) (any, error) {
	panic("ProcessStmt is unimplemented")
}

// ProcessThrowStmt implements ast.Visitor.
func (generator *AstGenerator) ProcessThrowStmt(stmt *ast.ThrowStatement, context any) (any, error) {
	panic("ProcessThrowStmt is unimplemented")
}

// ProcessWhileStmt implements ast.Visitor.
func (generator *AstGenerator) ProcessWhileStmt(stmt *ast.WhileStatement, context any) (any, error) {
	panic("ProcessWhileStmt is unimplemented")
}
