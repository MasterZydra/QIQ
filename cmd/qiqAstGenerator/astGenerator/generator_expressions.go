package astGenerator

import "QIQ/cmd/qiq/ast"

// ProcessAnonymousFunctionCreationExpr implements ast.Visitor.
func (generator *AstGenerator) ProcessAnonymousFunctionCreationExpr(stmt *ast.AnonymousFunctionCreationExpression, _ any) (any, error) {
	panic("ProcessAnonymousFunctionCreationExpr unimplemented")
}

// ProcessArrayLiteralExpr implements ast.Visitor.
func (generator *AstGenerator) ProcessArrayLiteralExpr(stmt *ast.ArrayLiteralExpression, _ any) (any, error) {
	generator.print(`ast.NewArrayLiteralExpr(0, nil)`)
	return nil, nil
}

// ProcessArrayNextKeyExpr implements ast.Visitor.
func (generator *AstGenerator) ProcessArrayNextKeyExpr(stmt *ast.ArrayNextKeyExpression, _ any) (any, error) {
	panic("ProcessArrayNextKeyExpr unimplemented")
}

// ProcessBinaryOpExpr implements ast.Visitor.
func (generator *AstGenerator) ProcessBinaryOpExpr(stmt *ast.BinaryOpExpression, _ any) (any, error) {
	panic("ProcessBinaryOpExpr unimplemented")
}

// ProcessCastExpr implements ast.Visitor.
func (generator *AstGenerator) ProcessCastExpr(stmt *ast.CastExpression, _ any) (any, error) {
	panic("ProcessCastExpr unimplemented")
}

// ProcessCoalesceExpr implements ast.Visitor.
func (generator *AstGenerator) ProcessCoalesceExpr(stmt *ast.CoalesceExpression, _ any) (any, error) {
	panic("ProcessCoalesceExpr unimplemented")
}

// ProcessCompoundAssignmentExpr implements ast.Visitor.
func (generator *AstGenerator) ProcessCompoundAssignmentExpr(stmt *ast.CompoundAssignmentExpression, _ any) (any, error) {
	panic("ProcessCompoundAssignmentExpr unimplemented")
}

// ProcessConditionalExpr implements ast.Visitor.
func (generator *AstGenerator) ProcessConditionalExpr(stmt *ast.ConditionalExpression, _ any) (any, error) {
	panic("ProcessConditionalExpr unimplemented")
}

// ProcessConstantAccessExpr implements ast.Visitor.
func (generator *AstGenerator) ProcessConstantAccessExpr(stmt *ast.ConstantAccessExpression, _ any) (any, error) {
	generator.print(`ast.NewConstantAccessExpr(0, nil, "%s")`, stmt.ConstantName)
	return nil, nil
}

// ProcessEmptyIntrinsicExpr implements ast.Visitor.
func (generator *AstGenerator) ProcessEmptyIntrinsicExpr(stmt *ast.EmptyIntrinsicExpression, _ any) (any, error) {
	panic("ProcessEmptyIntrinsicExpr unimplemented")
}

// ProcessEqualityExpr implements ast.Visitor.
func (generator *AstGenerator) ProcessEqualityExpr(stmt *ast.EqualityExpression, _ any) (any, error) {
	panic("ProcessEqualityExpr unimplemented")
}

// ProcessErrorControlExpr implements ast.Visitor.
func (generator *AstGenerator) ProcessErrorControlExpr(stmt *ast.ErrorControlExpression, _ any) (any, error) {
	panic("ProcessErrorControlExpr unimplemented")
}

// ProcessEvalIntrinsicExpr implements ast.Visitor.
func (generator *AstGenerator) ProcessEvalIntrinsicExpr(stmt *ast.EvalIntrinsicExpression, _ any) (any, error) {
	panic("ProcessEvalIntrinsicExpr unimplemented")
}

// ProcessExitIntrinsicExpr implements ast.Visitor.
func (generator *AstGenerator) ProcessExitIntrinsicExpr(stmt *ast.ExitIntrinsicExpression, _ any) (any, error) {
	panic("ProcessExitIntrinsicExpr unimplemented")
}

// ProcessExpr implements ast.Visitor.
func (generator *AstGenerator) ProcessExpr(stmt *ast.Expression, _ any) (any, error) {
	panic("ProcessExpr unimplemented")
}

// ProcessFloatingLiteralExpr implements ast.Visitor.
func (generator *AstGenerator) ProcessFloatingLiteralExpr(stmt *ast.FloatingLiteralExpression, _ any) (any, error) {
	panic("ProcessFloatingLiteralExpr unimplemented")
}

// ProcessFunctionCallExpr implements ast.Visitor.
func (generator *AstGenerator) ProcessFunctionCallExpr(stmt *ast.FunctionCallExpression, _ any) (any, error) {
	generator.print("ast.NewFunctionCallExpr(0, nil, ")
	generator.processStmt(stmt.FunctionName)
	generator.print(", []ast.IExpression{")
	for i, expr := range stmt.Arguments {
		if i > 0 {
			generator.print(", ")
		}
		generator.processStmt(expr)
	}
	generator.print("})")
	return nil, nil
}

// ProcessIncludeExpr implements ast.Visitor.
func (generator *AstGenerator) ProcessIncludeExpr(stmt *ast.IncludeExpression, _ any) (any, error) {
	panic("ProcessIncludeExpr unimplemented")
}

// ProcessIncludeOnceExpr implements ast.Visitor.
func (generator *AstGenerator) ProcessIncludeOnceExpr(stmt *ast.IncludeOnceExpression, _ any) (any, error) {
	panic("ProcessIncludeOnceExpr unimplemented")
}

// ProcessIntegerLiteralExpr implements ast.Visitor.
func (generator *AstGenerator) ProcessIntegerLiteralExpr(stmt *ast.IntegerLiteralExpression, _ any) (any, error) {
	panic("ProcessIntegerLiteralExpr unimplemented")
}

// ProcessIssetIntrinsicExpr implements ast.Visitor.
func (generator *AstGenerator) ProcessIssetIntrinsicExpr(stmt *ast.IssetIntrinsicExpression, _ any) (any, error) {
	panic("ProcessIssetIntrinsicExpr unimplemented")
}

// ProcessLogicalExpr implements ast.Visitor.
func (generator *AstGenerator) ProcessLogicalExpr(stmt *ast.LogicalExpression, _ any) (any, error) {
	panic("ProcessLogicalExpr unimplemented")
}

// ProcessLogicalNotExpr implements ast.Visitor.
func (generator *AstGenerator) ProcessLogicalNotExpr(stmt *ast.LogicalNotExpression, _ any) (any, error) {
	panic("ProcessLogicalNotExpr unimplemented")
}

// ProcessMemberAccessExpr implements ast.Visitor.
func (generator *AstGenerator) ProcessMemberAccessExpr(stmt *ast.MemberAccessExpression, _ any) (any, error) {
	generator.print("ast.NewMemberAccessExpr(0, nil, ")
	generator.processStmt(stmt.Object)
	generator.print(", ")
	generator.processStmt(stmt.Member)
	generator.print(")")
	return nil, nil
}

// ProcessObjectCreationExpr implements ast.Visitor.
func (generator *AstGenerator) ProcessObjectCreationExpr(stmt *ast.ObjectCreationExpression, _ any) (any, error) {
	panic("ProcessObjectCreationExpr unimplemented")
}

// ProcessParenthesizedExpr implements ast.Visitor.
func (generator *AstGenerator) ProcessParenthesizedExpr(stmt *ast.ParenthesizedExpression, _ any) (any, error) {
	panic("ProcessParenthesizedExpr unimplemented")
}

// ProcessPostfixIncExpr implements ast.Visitor.
func (generator *AstGenerator) ProcessPostfixIncExpr(stmt *ast.PostfixIncExpression, _ any) (any, error) {
	panic("ProcessPostfixIncExpr unimplemented")
}

// ProcessPrefixIncExpr implements ast.Visitor.
func (generator *AstGenerator) ProcessPrefixIncExpr(stmt *ast.PrefixIncExpression, _ any) (any, error) {
	panic("ProcessPrefixIncExpr unimplemented")
}

// ProcessPrintExpr implements ast.Visitor.
func (generator *AstGenerator) ProcessPrintExpr(stmt *ast.PrintExpression, _ any) (any, error) {
	panic("ProcessPrintExpr unimplemented")
}

// ProcessRelationalExpr implements ast.Visitor.
func (generator *AstGenerator) ProcessRelationalExpr(stmt *ast.RelationalExpression, _ any) (any, error) {
	panic("ProcessRelationalExpr unimplemented")
}

// ProcessRequireExpr implements ast.Visitor.
func (generator *AstGenerator) ProcessRequireExpr(stmt *ast.RequireExpression, _ any) (any, error) {
	panic("ProcessRequireExpr unimplemented")
}

// ProcessRequireOnceExpr implements ast.Visitor.
func (generator *AstGenerator) ProcessRequireOnceExpr(stmt *ast.RequireOnceExpression, _ any) (any, error) {
	panic("ProcessRequireOnceExpr unimplemented")
}

// ProcessSimpleAssignmentExpr implements ast.Visitor.
func (generator *AstGenerator) ProcessSimpleAssignmentExpr(stmt *ast.SimpleAssignmentExpression, _ any) (any, error) {
	generator.print("ast.NewSimpleAssignmentExpr(0, ")
	generator.processStmt(stmt.Variable)
	generator.print(", ")
	generator.processStmt(stmt.Value)
	generator.print(")")
	return nil, nil
}

// ProcessSimpleVariableExpr implements ast.Visitor.
func (generator *AstGenerator) ProcessSimpleVariableExpr(stmt *ast.SimpleVariableExpression, _ any) (any, error) {
	generator.print("ast.NewSimpleVariableExpr(0, ")
	generator.processStmt(stmt.VariableName)
	generator.print(")")
	return nil, nil
}

// ProcessStringLiteralExpr implements ast.Visitor.
func (generator *AstGenerator) ProcessStringLiteralExpr(stmt *ast.StringLiteralExpression, _ any) (any, error) {
	generator.print(`ast.NewStringLiteralExpr(0, nil, "%s", ast.DoubleQuotedString)`, stmt.Value)
	return nil, nil
}

// ProcessSubscriptExpr implements ast.Visitor.
func (generator *AstGenerator) ProcessSubscriptExpr(stmt *ast.SubscriptExpression, _ any) (any, error) {
	panic("ProcessSubscriptExpr unimplemented")
}

// ProcessTextExpr implements ast.Visitor.
func (generator *AstGenerator) ProcessTextExpr(stmt *ast.TextExpression, _ any) (any, error) {
	panic("ProcessTextExpr unimplemented")
}

// ProcessUnaryExpr implements ast.Visitor.
func (generator *AstGenerator) ProcessUnaryExpr(stmt *ast.UnaryOpExpression, _ any) (any, error) {
	panic("ProcessUnaryExpr unimplemented")
}

// ProcessUnsetIntrinsicExpr implements ast.Visitor.
func (generator *AstGenerator) ProcessUnsetIntrinsicExpr(stmt *ast.UnsetIntrinsicExpression, _ any) (any, error) {
	panic("ProcessUnsetIntrinsicExpr unimplemented")
}

// ProcessVariableNameExpr implements ast.Visitor.
func (generator *AstGenerator) ProcessVariableNameExpr(stmt *ast.VariableNameExpression, _ any) (any, error) {
	generator.print(`ast.NewVariableNameExpr(0, nil, "%s")`, stmt.VariableName)
	return nil, nil
}
