package ast

import "fmt"

func ToString(stmt IStatement) string {
	if stmt == nil {
		return "nil"
	}
	result, _ := stmt.Process(DumpVisitor{}, nil)
	return result.(string)
}

type DumpVisitor struct {
}

func dumpStatements(statements []IStatement) string {
	stmts := "{"
	for _, statement := range statements {
		stmts += ToString(statement) + ", "
	}
	stmts += "}"
	return stmts
}

func dumpExpressions(expressions []IExpression) string {
	exprs := "{"
	for _, expression := range expressions {
		exprs += ToString(expression) + ", "
	}
	exprs += "}"
	return exprs
}

// ProcessArrayLiteralExpr implements Visitor.
func (visitor DumpVisitor) ProcessArrayLiteralExpr(stmt *ArrayLiteralExpression, _ any) (any, error) {
	elements := "{"
	for _, key := range stmt.Keys {
		elements += ToString(key) + " => " + ToString(stmt.Elements[key]) + ", "
	}
	elements += "}"
	return fmt.Sprintf("{%s - elements: %s }", stmt.GetKind(), elements), nil
}

// ProcessArrayNextKeyExpr implements Visitor.
func (visitor DumpVisitor) ProcessArrayNextKeyExpr(stmt *ArrayNextKeyExpression, _ any) (any, error) {
	return fmt.Sprintf("{%s}", stmt.GetKind()), nil
}

// ProcessBinaryOpExpr implements Visitor.
func (visitor DumpVisitor) ProcessBinaryOpExpr(stmt *BinaryOpExpression, _ any) (any, error) {
	return fmt.Sprintf(
		"{%s - lhs: %s, operator: \"%s\" rhs: %s }",
		stmt.GetKind(), ToString(stmt.Lhs), stmt.Operator, ToString(stmt.Rhs),
	), nil
}

// ProcessBreakStmt implements Visitor.
func (visitor DumpVisitor) ProcessBreakStmt(stmt *BreakStatement, _ any) (any, error) {
	return fmt.Sprintf("{%s - %s}", stmt.GetKind(), ToString(stmt.Expr)), nil
}

// ProcessCastExpr implements Visitor.
func (visitor DumpVisitor) ProcessCastExpr(stmt *CastExpression, _ any) (any, error) {
	return fmt.Sprintf(
		"{%s - operator: \"%s\" expr: %s }",
		stmt.GetKind(), stmt.Operator, ToString(stmt.Expr),
	), nil
}

// ProcessCoalesceExpr implements Visitor.
func (visitor DumpVisitor) ProcessCoalesceExpr(stmt *CoalesceExpression, _ any) (any, error) {
	return fmt.Sprintf(
		"{%s - condition: %s, elseExpr: %s }",
		stmt.GetKind(), ToString(stmt.Cond), ToString(stmt.ElseExpr),
	), nil
}

// ProcessCompoundAssignmentExpr implements Visitor.
func (visitor DumpVisitor) ProcessCompoundAssignmentExpr(stmt *CompoundAssignmentExpression, _ any) (any, error) {
	return fmt.Sprintf(
		"{%s - variable: %s, operator: \"%s\", value: %s }",
		stmt.GetKind(), ToString(stmt.Variable), stmt.Operator, ToString(stmt.Value),
	), nil
}

// ProcessCompoundStmt implements Visitor.
func (visitor DumpVisitor) ProcessCompoundStmt(stmt *CompoundStatement, _ any) (any, error) {
	return fmt.Sprintf("{%s - %s}", stmt.GetKind(), dumpStatements(stmt.Statements)), nil
}

// ProcessConditionalExpr implements Visitor.
func (visitor DumpVisitor) ProcessConditionalExpr(stmt *ConditionalExpression, _ any) (any, error) {
	return fmt.Sprintf(
		"{%s - condition: %s, ifExpr: %s, elseExpr: %s }",
		stmt.GetKind(), ToString(stmt.Cond), ToString(stmt.IfExpr), ToString(stmt.ElseExpr),
	), nil
}

// ProcessConstDeclarationStmt implements Visitor.
func (visitor DumpVisitor) ProcessConstDeclarationStmt(stmt *ConstDeclarationStatement, _ any) (any, error) {
	return fmt.Sprintf("{%s - name: \"%s\" value: %s}", stmt.GetKind(), stmt.Name, ToString(stmt.Value)), nil
}

// ProcessConstantAccessExpr implements Visitor.
func (visitor DumpVisitor) ProcessConstantAccessExpr(stmt *ConstantAccessExpression, _ any) (any, error) {
	return fmt.Sprintf("{%s - constantName: %s}", stmt.GetKind(), stmt.ConstantName), nil
}

// ProcessContinueStmt implements Visitor.
func (visitor DumpVisitor) ProcessContinueStmt(stmt *ContinueStatement, _ any) (any, error) {
	return fmt.Sprintf("{%s - %s}", stmt.GetKind(), ToString(stmt.Expr)), nil
}

// ProcessDoStmt implements Visitor.
func (visitor DumpVisitor) ProcessDoStmt(stmt *DoStatement, _ any) (any, error) {
	return fmt.Sprintf(
		"{%s - condition: %s, block: %s}",
		stmt.GetKind(), ToString(stmt.Condition), ToString(stmt.Block),
	), nil
}

// ProcessEchoStmt implements Visitor.
func (visitor DumpVisitor) ProcessEchoStmt(stmt *EchoStatement, _ any) (any, error) {
	return fmt.Sprintf("{%s - %s}", stmt.GetKind(), dumpExpressions(stmt.Expressions)), nil
}

// ProcessEmptyIntrinsicExpr implements Visitor.
func (visitor DumpVisitor) ProcessEmptyIntrinsicExpr(stmt *EmptyIntrinsicExpression, _ any) (any, error) {
	return fmt.Sprintf(
		"{%s - functionName: \"%s\" arguments: %s}",
		stmt.GetKind(), ToString(stmt.FunctionName), dumpExpressions(stmt.Arguments),
	), nil
}

// ProcessEqualityExpr implements Visitor.
func (visitor DumpVisitor) ProcessEqualityExpr(stmt *EqualityExpression, _ any) (any, error) {
	return fmt.Sprintf(
		"{%s - lhs: %s, operator: \"%s\" rhs: %s }",
		stmt.GetKind(), ToString(stmt.Lhs), stmt.Operator, ToString(stmt.Rhs),
	), nil
}

// ProcessEvalIntrinsicExpr implements Visitor.
func (visitor DumpVisitor) ProcessEvalIntrinsicExpr(stmt *EvalIntrinsicExpression, _ any) (any, error) {
	return fmt.Sprintf(
		"{%s - functionName: \"%s\" arguments: %s}",
		stmt.GetKind(), ToString(stmt.FunctionName), dumpExpressions(stmt.Arguments),
	), nil
}

// ProcessExitIntrinsicExpr implements Visitor.
func (visitor DumpVisitor) ProcessExitIntrinsicExpr(stmt *ExitIntrinsicExpression, _ any) (any, error) {
	return fmt.Sprintf(
		"{%s - functionName: \"%s\" arguments: %s}",
		stmt.GetKind(), ToString(stmt.FunctionName), dumpExpressions(stmt.Arguments),
	), nil
}

// ProcessExpr implements Visitor.
func (visitor DumpVisitor) ProcessExpr(stmt *Expression, _ any) (any, error) {
	return fmt.Sprintf("{%s}", stmt.GetKind()), nil
}

// ProcessExpressionStmt implements Visitor.
func (visitor DumpVisitor) ProcessExpressionStmt(stmt *ExpressionStatement, _ any) (any, error) {
	return fmt.Sprintf("{%s - %s}", stmt.GetKind(), ToString(stmt.Expr)), nil
}

// ProcessFloatingLiteralExpr implements Visitor.
func (visitor DumpVisitor) ProcessFloatingLiteralExpr(stmt *FloatingLiteralExpression, _ any) (any, error) {
	return fmt.Sprintf("{%s - value: %f }", stmt.GetKind(), stmt.Value), nil
}

// ProcessForStmt implements Visitor.
func (visitor DumpVisitor) ProcessForStmt(stmt *ForStatement, context any) (any, error) {
	return fmt.Sprintf(
		"{%s - initializer: %s, control: %s, endOfLoop: %s, block: %s}",
		stmt.GetKind(), ToString(stmt.Initializer), ToString(stmt.Control), ToString(stmt.EndOfLoop), ToString(stmt.Block),
	), nil
}

// ProcessFunctionCallExpr implements Visitor.
func (visitor DumpVisitor) ProcessFunctionCallExpr(stmt *FunctionCallExpression, _ any) (any, error) {
	return fmt.Sprintf(
		"{%s - functionName: \"%s\", arguments: %s}",
		stmt.GetKind(), ToString(stmt.FunctionName), dumpExpressions(stmt.Arguments),
	), nil
}

// ProcessFunctionDefinitionStmt implements Visitor.
func (visitor DumpVisitor) ProcessFunctionDefinitionStmt(stmt *FunctionDefinitionStatement, _ any) (any, error) {
	return fmt.Sprintf("{%s - name: %s, params: %s, body: %s, returnType: %s}",
		stmt.GetKind(), stmt.FunctionName, stmt.Params, ToString(stmt.Body), stmt.ReturnType,
	), nil
}

// ProcessIfStmt implements Visitor.
func (visitor DumpVisitor) ProcessIfStmt(stmt *IfStatement, _ any) (any, error) {
	elseIf := "{"
	for _, elseIfStmt := range stmt.ElseIf {
		elseIf += ToString(elseIfStmt) + ", "
	}
	elseIf += "}"
	return fmt.Sprintf(
		"{%s - condition: %s, ifBlock: %s, elseIf: %s, else: %s}",
		stmt.GetKind(), ToString(stmt.Condition), ToString(stmt.IfBlock), elseIf, ToString(stmt.ElseBlock),
	), nil
}

// ProcessIncludeExpr implements Visitor.
func (visitor DumpVisitor) ProcessIncludeExpr(stmt *IncludeExpression, _ any) (any, error) {
	return fmt.Sprintf("{%s - expr: %s }", stmt.GetKind(), ToString(stmt.Expr)), nil
}

// ProcessIncludeOnceExpr implements Visitor.
func (visitor DumpVisitor) ProcessIncludeOnceExpr(stmt *IncludeOnceExpression, _ any) (any, error) {
	return fmt.Sprintf("{%s - expr: %s }", stmt.GetKind(), ToString(stmt.Expr)), nil
}

// ProcessIntegerLiteralExpr implements Visitor.
func (visitor DumpVisitor) ProcessIntegerLiteralExpr(stmt *IntegerLiteralExpression, _ any) (any, error) {
	return fmt.Sprintf("{%s - value: %d }", stmt.GetKind(), stmt.Value), nil
}

// ProcessIssetIntrinsicExpr implements Visitor.
func (visitor DumpVisitor) ProcessIssetIntrinsicExpr(stmt *IssetIntrinsicExpression, _ any) (any, error) {
	return fmt.Sprintf(
		"{%s - functionName: \"%s\" arguments: %s}",
		stmt.GetKind(), ToString(stmt.FunctionName), dumpExpressions(stmt.Arguments),
	), nil
}

// ProcessLogicalExpr implements Visitor.
func (visitor DumpVisitor) ProcessLogicalExpr(stmt *LogicalExpression, _ any) (any, error) {
	return fmt.Sprintf(
		"{%s - lhs: %s, operator: \"%s\" rhs: %s }",
		stmt.GetKind(), ToString(stmt.Lhs), stmt.Operator, ToString(stmt.Rhs),
	), nil
}

// ProcessLogicalNotExpr implements Visitor.
func (visitor DumpVisitor) ProcessLogicalNotExpr(stmt *LogicalNotExpression, _ any) (any, error) {
	return fmt.Sprintf("{%s - operator: \"%s\" expr: %s }", stmt.GetKind(), stmt.Operator, ToString(stmt.Expr)), nil
}

// ProcessParenthesizedExpr implements Visitor.
func (visitor DumpVisitor) ProcessParenthesizedExpr(stmt *ParenthesizedExpression, _ any) (any, error) {
	return fmt.Sprintf("{%s - expr: %s}", stmt.GetKind(), ToString(stmt.Expr)), nil
}

// ProcessPostfixIncExpr implements Visitor.
func (visitor DumpVisitor) ProcessPostfixIncExpr(stmt *PostfixIncExpression, _ any) (any, error) {
	return fmt.Sprintf("{%s - operator: \"%s\" expr: %s }", stmt.GetKind(), stmt.Operator, ToString(stmt.Expr)), nil
}

// ProcessPrefixIncExpr implements Visitor.
func (visitor DumpVisitor) ProcessPrefixIncExpr(stmt *PrefixIncExpression, _ any) (any, error) {
	return fmt.Sprintf("{%s - operator: \"%s\" expr: %s }", stmt.GetKind(), stmt.Operator, ToString(stmt.Expr)), nil
}

// ProcessPrintExpr implements Visitor.
func (visitor DumpVisitor) ProcessPrintExpr(stmt *PrintExpression, context any) (any, error) {
	return fmt.Sprintf("{%s - expr: %s }", stmt.GetKind(), ToString(stmt.Expr)), nil
}

// ProcessRelationalExpr implements Visitor.
func (visitor DumpVisitor) ProcessRelationalExpr(stmt *RelationalExpression, _ any) (any, error) {
	return fmt.Sprintf(
		"{%s - lhs: %s, operator: \"%s\" rhs: %s }",
		stmt.GetKind(), ToString(stmt.Lhs), stmt.Operator, ToString(stmt.Rhs),
	), nil
}

// ProcessRequireExpr implements Visitor.
func (visitor DumpVisitor) ProcessRequireExpr(stmt *RequireExpression, _ any) (any, error) {
	return fmt.Sprintf("{%s - expr: %s }", stmt.GetKind(), ToString(stmt.Expr)), nil
}

// ProcessRequireOnceExpr implements Visitor.
func (visitor DumpVisitor) ProcessRequireOnceExpr(stmt *RequireOnceExpression, _ any) (any, error) {
	return fmt.Sprintf("{%s - expr: %s }", stmt.GetKind(), ToString(stmt.Expr)), nil
}

// ProcessReturnStmt implements Visitor.
func (visitor DumpVisitor) ProcessReturnStmt(stmt *ReturnStatement, _ any) (any, error) {
	return fmt.Sprintf("{%s - %s}", stmt.GetKind(), ToString(stmt.Expr)), nil
}

// ProcessSimpleAssignmentExpr implements Visitor.
func (visitor DumpVisitor) ProcessSimpleAssignmentExpr(stmt *SimpleAssignmentExpression, _ any) (any, error) {
	return fmt.Sprintf(
		"{%s - variable: %s, value: %s }",
		stmt.GetKind(), ToString(stmt.Variable), ToString(stmt.Value),
	), nil
}

// ProcessSimpleVariableExpr implements Visitor.
func (visitor DumpVisitor) ProcessSimpleVariableExpr(stmt *SimpleVariableExpression, _ any) (any, error) {
	return fmt.Sprintf("{%s - variableName: \"%s\" }", stmt.GetKind(), ToString(stmt.VariableName)), nil
}

// ProcessStmt implements Visitor.
func (visitor DumpVisitor) ProcessStmt(stmt *Statement, _ any) (any, error) {
	return fmt.Sprintf("{%s}", stmt.GetKind()), nil
}

// ProcessStringLiteralExpr implements Visitor.
func (visitor DumpVisitor) ProcessStringLiteralExpr(stmt *StringLiteralExpression, _ any) (any, error) {
	return fmt.Sprintf("{%s - type: \"%s\" value: \"%s\" }", stmt.GetKind(), stmt.StringType, stmt.Value), nil
}

// ProcessSubscriptExpr implements Visitor.
func (visitor DumpVisitor) ProcessSubscriptExpr(stmt *SubscriptExpression, _ any) (any, error) {
	return fmt.Sprintf(
		"{%s - variable: %s, index: \"%s\" }",
		stmt.GetKind(), ToString(stmt.Variable), ToString(stmt.Index),
	), nil
}

// ProcessTextExpr implements Visitor.
func (visitor DumpVisitor) ProcessTextExpr(stmt *TextExpression, _ any) (any, error) {
	return fmt.Sprintf("{%s - value: \"%s\" }", stmt.GetKind(), stmt.Value), nil
}

// ProcessUnaryExpr implements Visitor.
func (visitor DumpVisitor) ProcessUnaryExpr(stmt *UnaryOpExpression, _ any) (any, error) {
	return fmt.Sprintf(
		"{%s - operator: \"%s\" expr: %s }",
		stmt.GetKind(), stmt.Operator, ToString(stmt.Expr),
	), nil
}

// ProcessUnsetIntrinsicExpr implements Visitor.
func (visitor DumpVisitor) ProcessUnsetIntrinsicExpr(stmt *UnsetIntrinsicExpression, _ any) (any, error) {
	return fmt.Sprintf(
		"{%s - functionName: \"%s\" arguments: %s}",
		stmt.GetKind(), ToString(stmt.FunctionName), dumpExpressions(stmt.Arguments),
	), nil
}

// ProcessWhileStmt implements Visitor.
func (visitor DumpVisitor) ProcessWhileStmt(stmt *WhileStatement, _ any) (any, error) {
	return fmt.Sprintf(
		"{%s - condition: %s, block: %s}",
		stmt.GetKind(), ToString(stmt.Condition), ToString(stmt.Block),
	), nil
}

// ProcessVariableNameExpr implements Visitor.
func (visitor DumpVisitor) ProcessVariableNameExpr(stmt *VariableNameExpression, _ any) (any, error) {
	return fmt.Sprintf("{%s - variableName: \"%s\" }", stmt.GetKind(), stmt.VariableName), nil
}
