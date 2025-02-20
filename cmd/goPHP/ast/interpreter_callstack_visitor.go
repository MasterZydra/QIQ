package ast

import (
	"GoPHP/cmd/goPHP/config"
	"fmt"
)

func PrintInterpreterCallstack(stmt IStatement) {
	if !config.ShowInterpreterCallStack {
		return
	}
	if stmt == nil {
		println("nil")
		return
	}
	result, _ := stmt.Process(InterpreterCallStackVisitor{}, nil)
	println(result.(string))
}

type InterpreterCallStackVisitor struct {
}

func dumpInternalCallstackStatements(statements []IStatement) string {
	stmts := "{"
	for _, statement := range statements {
		stmts += ToString(statement) + ", "
	}
	stmts += "}"
	return stmts
}

func dumpInternalCallstackExpressions(expressions []IExpression) string {
	exprs := "{"
	for _, expression := range expressions {
		exprs += ToString(expression) + ", "
	}
	exprs += "}"
	return exprs
}

// ProcessArrayLiteralExpr implements Visitor.
func (visitor InterpreterCallStackVisitor) ProcessArrayLiteralExpr(stmt *ArrayLiteralExpression, _ any) (any, error) {
	elements := "{"
	for _, key := range stmt.Keys {
		elements += ToString(key) + " => " + ToString(stmt.Elements[key]) + ", "
	}
	elements += "}"
	return fmt.Sprintf("{%s - elements: %s , pos: %s}", stmt.GetKind(), elements, stmt.GetPosition().ToPosString()), nil
}

// ProcessBinaryOpExpr implements Visitor.
func (visitor InterpreterCallStackVisitor) ProcessBinaryOpExpr(stmt *BinaryOpExpression, _ any) (any, error) {
	return fmt.Sprintf(
		"{%s - lhs: %s, operator: \"%s\" rhs: %s, pos: %s }",
		stmt.GetKind(), ToString(stmt.Lhs), stmt.Operator, ToString(stmt.Rhs), stmt.GetPosition().ToPosString(),
	), nil
}

// ProcessBreakStmt implements Visitor.
func (visitor InterpreterCallStackVisitor) ProcessBreakStmt(stmt *BreakStatement, _ any) (any, error) {
	return fmt.Sprintf("{%s - %s, pos: %s}", stmt.GetKind(), ToString(stmt.Expr), stmt.GetPosition().ToPosString()), nil
}

// ProcessCastExpr implements Visitor.
func (visitor InterpreterCallStackVisitor) ProcessCastExpr(stmt *CastExpression, _ any) (any, error) {
	return fmt.Sprintf(
		"{%s - operator: \"%s\" expr: %s, pos: %s }",
		stmt.GetKind(), stmt.Operator, ToString(stmt.Expr), stmt.GetPosition().ToPosString(),
	), nil
}

// ProcessCoalesceExpr implements Visitor.
func (visitor InterpreterCallStackVisitor) ProcessCoalesceExpr(stmt *CoalesceExpression, _ any) (any, error) {
	return fmt.Sprintf(
		"{%s - condition: %s, elseExpr: %s, pos: %s }",
		stmt.GetKind(), ToString(stmt.Cond), ToString(stmt.ElseExpr), stmt.GetPosition().ToPosString(),
	), nil
}

// ProcessCompoundAssignmentExpr implements Visitor.
func (visitor InterpreterCallStackVisitor) ProcessCompoundAssignmentExpr(stmt *CompoundAssignmentExpression, _ any) (any, error) {
	return fmt.Sprintf(
		"{%s - variable: %s, operator: \"%s\", value: %s, pos: %s }",
		stmt.GetKind(), ToString(stmt.Variable), stmt.Operator, ToString(stmt.Value), stmt.GetPosition().ToPosString(),
	), nil
}

// ProcessCompoundStmt implements Visitor.
func (visitor InterpreterCallStackVisitor) ProcessCompoundStmt(stmt *CompoundStatement, _ any) (any, error) {
	return fmt.Sprintf("{%s - %s}", stmt.GetKind(), dumpInternalCallstackStatements(stmt.Statements)), nil
}

// ProcessConditionalExpr implements Visitor.
func (visitor InterpreterCallStackVisitor) ProcessConditionalExpr(stmt *ConditionalExpression, _ any) (any, error) {
	return fmt.Sprintf(
		"{%s - condition: %s, ifExpr: %s, elseExpr: %s, pos: %s }",
		stmt.GetKind(), ToString(stmt.Cond), ToString(stmt.IfExpr), ToString(stmt.ElseExpr), stmt.GetPosition().ToPosString(),
	), nil
}

// ProcessConstDeclarationStmt implements Visitor.
func (visitor InterpreterCallStackVisitor) ProcessConstDeclarationStmt(stmt *ConstDeclarationStatement, _ any) (any, error) {
	return fmt.Sprintf("{%s - name: \"%s\" value: %s, pos: %s}", stmt.GetKind(), stmt.Name, ToString(stmt.Value), stmt.GetPosition().ToPosString()), nil
}

// ProcessConstantAccessExpr implements Visitor.
func (visitor InterpreterCallStackVisitor) ProcessConstantAccessExpr(stmt *ConstantAccessExpression, _ any) (any, error) {
	return fmt.Sprintf("{%s - constantName: %s, pos: %s}", stmt.GetKind(), stmt.ConstantName, stmt.GetPosition().ToPosString()), nil
}

// ProcessContinueStmt implements Visitor.
func (visitor InterpreterCallStackVisitor) ProcessContinueStmt(stmt *ContinueStatement, _ any) (any, error) {
	return fmt.Sprintf("{%s - %s, pos: %s}", stmt.GetKind(), ToString(stmt.Expr), stmt.GetPosition().ToPosString()), nil
}

// ProcessDoStmt implements Visitor.
func (visitor InterpreterCallStackVisitor) ProcessDoStmt(stmt *DoStatement, _ any) (any, error) {
	return fmt.Sprintf(
		"{%s - condition: %s, block: %s, pos: %s }",
		stmt.GetKind(), ToString(stmt.Condition), ToString(stmt.Block), stmt.GetPosition().ToPosString(),
	), nil
}

// ProcessEchoStmt implements Visitor.
func (visitor InterpreterCallStackVisitor) ProcessEchoStmt(stmt *EchoStatement, _ any) (any, error) {
	return fmt.Sprintf("{%s - %s, pos: %s}", stmt.GetKind(), dumpInternalCallstackExpressions(stmt.Expressions), stmt.GetPosition().ToPosString()), nil
}

// ProcessEmptyIntrinsicExpr implements Visitor.
func (visitor InterpreterCallStackVisitor) ProcessEmptyIntrinsicExpr(stmt *EmptyIntrinsicExpression, _ any) (any, error) {
	return fmt.Sprintf(
		"{%s - functionName: \"%s\", arguments: %s, pos: %s}",
		stmt.GetKind(), stmt.FunctionName, dumpInternalCallstackExpressions(stmt.Arguments), stmt.GetPosition().ToPosString(),
	), nil
}

// ProcessEqualityExpr implements Visitor.
func (visitor InterpreterCallStackVisitor) ProcessEqualityExpr(stmt *EqualityExpression, _ any) (any, error) {
	return fmt.Sprintf(
		"{%s - lhs: %s, operator: \"%s\", rhs: %s, pos: %s }",
		stmt.GetKind(), ToString(stmt.Lhs), stmt.Operator, ToString(stmt.Rhs), stmt.GetPosition().ToPosString(),
	), nil
}

// ProcessEvalIntrinsicExpr implements Visitor.
func (visitor InterpreterCallStackVisitor) ProcessEvalIntrinsicExpr(stmt *EvalIntrinsicExpression, _ any) (any, error) {
	return fmt.Sprintf(
		"{%s - functionName: \"%s\", arguments: %s, pos: %s}",
		stmt.GetKind(), stmt.FunctionName, dumpInternalCallstackExpressions(stmt.Arguments), stmt.GetPosition().ToPosString(),
	), nil
}

// ProcessExitIntrinsicExpr implements Visitor.
func (visitor InterpreterCallStackVisitor) ProcessExitIntrinsicExpr(stmt *ExitIntrinsicExpression, _ any) (any, error) {
	return fmt.Sprintf(
		"{%s - functionName: \"%s\", arguments: %s, pos: %s}",
		stmt.GetKind(), stmt.FunctionName, dumpInternalCallstackExpressions(stmt.Arguments), stmt.GetPosition().ToPosString(),
	), nil
}

// ProcessExpr implements Visitor.
func (visitor InterpreterCallStackVisitor) ProcessExpr(stmt *Expression, _ any) (any, error) {
	return fmt.Sprintf("{%s - pos: %s}", stmt.GetKind(), stmt.GetPosition().ToPosString()), nil
}

// ProcessExpressionStmt implements Visitor.
func (visitor InterpreterCallStackVisitor) ProcessExpressionStmt(stmt *ExpressionStatement, _ any) (any, error) {
	return fmt.Sprintf("{%s - %s, pos: %s}", stmt.GetKind(), ToString(stmt.Expr), stmt.GetPosition().ToPosString()), nil
}

// ProcessFloatingLiteralExpr implements Visitor.
func (visitor InterpreterCallStackVisitor) ProcessFloatingLiteralExpr(stmt *FloatingLiteralExpression, _ any) (any, error) {
	return fmt.Sprintf("{%s - value: %f, pos: %s }", stmt.GetKind(), stmt.Value, stmt.GetPosition().ToPosString()), nil
}

// ProcessForStmt implements Visitor.
func (visitor InterpreterCallStackVisitor) ProcessForStmt(stmt *ForStatement, context any) (any, error) {
	return fmt.Sprintf(
		"{%s - initializer: %s, control: %s, endOfLoop: %s, block: %s, pos: %s}",
		stmt.GetKind(), ToString(stmt.Initializer), ToString(stmt.Control), ToString(stmt.EndOfLoop), ToString(stmt.Block), stmt.GetPosition().ToPosString(),
	), nil
}

// ProcessFunctionCallExpr implements Visitor.
func (visitor InterpreterCallStackVisitor) ProcessFunctionCallExpr(stmt *FunctionCallExpression, _ any) (any, error) {
	return fmt.Sprintf(
		"{%s - functionName: \"%s\", arguments: %s, pos: %s}",
		stmt.GetKind(), stmt.FunctionName, dumpInternalCallstackExpressions(stmt.Arguments), stmt.GetPosition().ToPosString(),
	), nil
}

// ProcessFunctionDefinitionStmt implements Visitor.
func (visitor InterpreterCallStackVisitor) ProcessFunctionDefinitionStmt(stmt *FunctionDefinitionStatement, _ any) (any, error) {
	return fmt.Sprintf("{%s - name: %s, params: %s, body: %s, returnType: %s, pos: %s}",
		stmt.GetKind(), stmt.FunctionName, stmt.Params, ToString(stmt.Body), stmt.ReturnType, stmt.GetPosition().ToPosString(),
	), nil
}

// ProcessIfStmt implements Visitor.
func (visitor InterpreterCallStackVisitor) ProcessIfStmt(stmt *IfStatement, _ any) (any, error) {
	elseIf := "{"
	for _, elseIfStmt := range stmt.ElseIf {
		elseIf += ToString(elseIfStmt) + ", "
	}
	elseIf += "}"
	return fmt.Sprintf(
		"{%s - condition: %s, ifBlock: %s, elseIf: %s, else: %s, pos: %s}",
		stmt.GetKind(), ToString(stmt.Condition), ToString(stmt.IfBlock), elseIf, ToString(stmt.ElseBlock), stmt.GetPosition().ToPosString(),
	), nil
}

// ProcessIncludeExpr implements Visitor.
func (visitor InterpreterCallStackVisitor) ProcessIncludeExpr(stmt *IncludeExpression, _ any) (any, error) {
	return fmt.Sprintf("{%s - expr: %s, pos: %s }", stmt.GetKind(), ToString(stmt.Expr), stmt.GetPosition().ToPosString()), nil
}

// ProcessIncludeOnceExpr implements Visitor.
func (visitor InterpreterCallStackVisitor) ProcessIncludeOnceExpr(stmt *IncludeOnceExpression, _ any) (any, error) {
	return fmt.Sprintf("{%s - expr: %s, pos: %s }", stmt.GetKind(), ToString(stmt.Expr), stmt.GetPosition().ToPosString()), nil
}

// ProcessIntegerLiteralExpr implements Visitor.
func (visitor InterpreterCallStackVisitor) ProcessIntegerLiteralExpr(stmt *IntegerLiteralExpression, _ any) (any, error) {
	return fmt.Sprintf("{%s - value: %d, pos: %s }", stmt.GetKind(), stmt.Value, stmt.GetPosition().ToPosString()), nil
}

// ProcessIssetIntrinsicExpr implements Visitor.
func (visitor InterpreterCallStackVisitor) ProcessIssetIntrinsicExpr(stmt *IssetIntrinsicExpression, _ any) (any, error) {
	return fmt.Sprintf(
		"{%s - functionName: \"%s\" arguments: %s, pos: %s }",
		stmt.GetKind(), stmt.FunctionName, dumpInternalCallstackExpressions(stmt.Arguments), stmt.GetPosition().ToPosString(),
	), nil
}

// ProcessLogicalExpr implements Visitor.
func (visitor InterpreterCallStackVisitor) ProcessLogicalExpr(stmt *LogicalExpression, _ any) (any, error) {
	return fmt.Sprintf(
		"{%s - lhs: %s, operator: \"%s\", rhs: %s, pos: %s }",
		stmt.GetKind(), ToString(stmt.Lhs), stmt.Operator, ToString(stmt.Rhs), stmt.GetPosition().ToPosString(),
	), nil
}

// ProcessLogicalNotExpr implements Visitor.
func (visitor InterpreterCallStackVisitor) ProcessLogicalNotExpr(stmt *LogicalNotExpression, _ any) (any, error) {
	return fmt.Sprintf("{%s - operator: \"%s\", expr: %s, pos: %s }", stmt.GetKind(), stmt.Operator, ToString(stmt.Expr), stmt.GetPosition().ToPosString()), nil
}

// ProcessParenthesizedExpr implements Visitor.
func (visitor InterpreterCallStackVisitor) ProcessParenthesizedExpr(stmt *ParenthesizedExpression, _ any) (any, error) {
	return fmt.Sprintf("{%s - expr: %s, pos: %s }", stmt.GetKind(), ToString(stmt.Expr), stmt.GetPosition().ToPosString()), nil
}

// ProcessPostfixIncExpr implements Visitor.
func (visitor InterpreterCallStackVisitor) ProcessPostfixIncExpr(stmt *PostfixIncExpression, _ any) (any, error) {
	return fmt.Sprintf("{%s - operator: \"%s\", expr: %s, pos: %s }", stmt.GetKind(), stmt.Operator, ToString(stmt.Expr), stmt.GetPosition().ToPosString()), nil
}

// ProcessPrefixIncExpr implements Visitor.
func (visitor InterpreterCallStackVisitor) ProcessPrefixIncExpr(stmt *PrefixIncExpression, _ any) (any, error) {
	return fmt.Sprintf("{%s - operator: \"%s\", expr: %s, pos: %s }", stmt.GetKind(), stmt.Operator, ToString(stmt.Expr), stmt.GetPosition().ToPosString()), nil
}

// ProcessPrintExpr implements Visitor.
func (visitor InterpreterCallStackVisitor) ProcessPrintExpr(stmt *PrintExpression, context any) (any, error) {
	return fmt.Sprintf("{%s - expr: %s, pos: %s }", stmt.GetKind(), ToString(stmt.Expr), stmt.GetPosition().ToPosString()), nil
}

// ProcessRelationalExpr implements Visitor.
func (visitor InterpreterCallStackVisitor) ProcessRelationalExpr(stmt *RelationalExpression, _ any) (any, error) {
	return fmt.Sprintf(
		"{%s - lhs: %s, operator: \"%s\", rhs: %s, pos: %s }",
		stmt.GetKind(), ToString(stmt.Lhs), stmt.Operator, ToString(stmt.Rhs), stmt.GetPosition().ToPosString(),
	), nil
}

// ProcessRequireExpr implements Visitor.
func (visitor InterpreterCallStackVisitor) ProcessRequireExpr(stmt *RequireExpression, _ any) (any, error) {
	return fmt.Sprintf("{%s - expr: %s, pos: %s }", stmt.GetKind(), ToString(stmt.Expr), stmt.GetPosition().ToPosString()), nil
}

// ProcessRequireOnceExpr implements Visitor.
func (visitor InterpreterCallStackVisitor) ProcessRequireOnceExpr(stmt *RequireOnceExpression, _ any) (any, error) {
	return fmt.Sprintf("{%s - expr: %s, pos: %s }", stmt.GetKind(), ToString(stmt.Expr), stmt.GetPosition().ToPosString()), nil
}

// ProcessReturnStmt implements Visitor.
func (visitor InterpreterCallStackVisitor) ProcessReturnStmt(stmt *ReturnStatement, _ any) (any, error) {
	return fmt.Sprintf("{%s - %s, pos: %s }", stmt.GetKind(), ToString(stmt.Expr), stmt.GetPosition().ToPosString()), nil
}

// ProcessSimpleAssignmentExpr implements Visitor.
func (visitor InterpreterCallStackVisitor) ProcessSimpleAssignmentExpr(stmt *SimpleAssignmentExpression, _ any) (any, error) {
	return fmt.Sprintf(
		"{%s - variable: %s, value: %s, pos: %s }",
		stmt.GetKind(), ToString(stmt.Variable), ToString(stmt.Value), stmt.GetPosition().ToPosString(),
	), nil
}

// ProcessSimpleVariableExpr implements Visitor.
func (visitor InterpreterCallStackVisitor) ProcessSimpleVariableExpr(stmt *SimpleVariableExpression, _ any) (any, error) {
	return fmt.Sprintf("{%s - variableName: \"%s\", pos: %s }", stmt.GetKind(), ToString(stmt.VariableName), stmt.GetPosition().ToPosString()), nil
}

// ProcessStmt implements Visitor.
func (visitor InterpreterCallStackVisitor) ProcessStmt(stmt *Statement, _ any) (any, error) {
	return fmt.Sprintf("{%s - pos: %s}", stmt.GetKind(), stmt.GetPosition().ToPosString()), nil
}

// ProcessStringLiteralExpr implements Visitor.
func (visitor InterpreterCallStackVisitor) ProcessStringLiteralExpr(stmt *StringLiteralExpression, _ any) (any, error) {
	return fmt.Sprintf("{%s - type: \"%s\", value: \"%s\", pos: %s }", stmt.GetKind(), stmt.StringType, stmt.Value, stmt.GetPosition().ToPosString()), nil
}

// ProcessSubscriptExpr implements Visitor.
func (visitor InterpreterCallStackVisitor) ProcessSubscriptExpr(stmt *SubscriptExpression, _ any) (any, error) {
	return fmt.Sprintf(
		"{%s - variable: %s, index: \"%s\", pos: %s }",
		stmt.GetKind(), ToString(stmt.Variable), ToString(stmt.Index), stmt.GetPosition().ToPosString(),
	), nil
}

// ProcessTextExpr implements Visitor.
func (visitor InterpreterCallStackVisitor) ProcessTextExpr(stmt *TextExpression, _ any) (any, error) {
	return fmt.Sprintf("{%s - value: \"%s\" }", stmt.GetKind(), stmt.Value), nil
}

// ProcessUnaryExpr implements Visitor.
func (visitor InterpreterCallStackVisitor) ProcessUnaryExpr(stmt *UnaryOpExpression, _ any) (any, error) {
	return fmt.Sprintf(
		"{%s - operator: \"%s\", expr: %s, pos: %s }",
		stmt.GetKind(), stmt.Operator, ToString(stmt.Expr), stmt.GetPosition().ToPosString(),
	), nil
}

// ProcessUnsetIntrinsicExpr implements Visitor.
func (visitor InterpreterCallStackVisitor) ProcessUnsetIntrinsicExpr(stmt *UnsetIntrinsicExpression, _ any) (any, error) {
	return fmt.Sprintf(
		"{%s - functionName: \"%s\", arguments: %s, pos: %s }",
		stmt.GetKind(), stmt.FunctionName, dumpInternalCallstackExpressions(stmt.Arguments), stmt.GetPosition().ToPosString(),
	), nil
}

// ProcessWhileStmt implements Visitor.
func (visitor InterpreterCallStackVisitor) ProcessWhileStmt(stmt *WhileStatement, _ any) (any, error) {
	return fmt.Sprintf(
		"{%s - condition: %s, block: %s, pos: %s }",
		stmt.GetKind(), ToString(stmt.Condition), ToString(stmt.Block), stmt.GetPosition().ToPosString(),
	), nil
}

// ProcessVariableNameExpr implements Visitor.
func (visitor InterpreterCallStackVisitor) ProcessVariableNameExpr(stmt *VariableNameExpression, _ any) (any, error) {
	return fmt.Sprintf("{%s - variableName: \"%s\", pos: %s }", stmt.GetKind(), stmt.VariableName, stmt.GetPosition().ToPosString()), nil
}
