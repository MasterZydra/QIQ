package interpreter

import (
	"GoPHP/cmd/goPHP/ast"
	"GoPHP/cmd/goPHP/parser"
	"fmt"
)

type Interpreter struct {
	parser *parser.Parser
	env    *Environment
	result string
}

func NewInterpreter() *Interpreter {
	return &Interpreter{parser: parser.NewParser(), env: NewEnvironment(nil)}
}

func (interpreter *Interpreter) Process(sourceCode string) (string, error) {
	interpreter.result = ""
	program, err := interpreter.parser.ProduceAST(sourceCode)
	if err != nil {
		return interpreter.result, err
	}

	for _, stmt := range program.GetStatements() {
		if _, err := interpreter.process(stmt); err != nil {
			return interpreter.result, err
		}
	}

	return interpreter.result, nil
}

func (interpreter *Interpreter) process(stmt ast.IStatement) (IRuntimeValue, error) {
	switch stmt.GetKind() {
	// Statements
	case ast.ConstDeclarationStmt:
		return interpreter.processConstDeclarationStatement(ast.StmtToConstDeclStatement(stmt))
	case ast.ExpressionStmt:
		return interpreter.process(ast.StmtToExprStatement(stmt).GetExpression())
	case ast.EchoStmt:
		return interpreter.processEchoStatement(ast.StmtToEchoStatement(stmt))

	// Expressions
	case ast.BooleanLiteralExpr, ast.IntegerLiteralExpr, ast.FloatingLiteralExpr, ast.StringLiteralExpr:
		return exprToRuntimeValue(stmt)
	case ast.TextNode:
		interpreter.print(ast.ExprToTextExpr(stmt).GetValue())
		return NewVoidRuntimeValue(), nil
	case ast.SimpleVariableExpr:
		return interpreter.processSimpleVariableExpression(ast.ExprToSimpleVarExpr(stmt))
	case ast.SimpleAssignmentExpr:
		return interpreter.processSimpleAssignmentExpression(ast.ExprToSimpleAssignExpr(stmt))
	case ast.ConstantAccessExpr:
		return interpreter.processConstantAccessExpression(ast.ExprToConstAccessExpr(stmt))
	case ast.CompoundAssignmentExpr:
		return interpreter.processCompoundAssignmentExpression(ast.ExprToCompoundAssignExpr(stmt))
	case ast.ConditionalExpr:
		return interpreter.processConditionalExpression(ast.ExprToCondExpr(stmt))

	default:
		return NewVoidRuntimeValue(), fmt.Errorf("Interpreter error: Unsupported statement or expression: %s", stmt)
	}
}

func (interpreter *Interpreter) processConstDeclarationStatement(stmt ast.IConstDeclarationStatement) (IRuntimeValue, error) {
	value, err := interpreter.process(stmt.GetValue())
	if err != nil {
		return NewVoidRuntimeValue(), err
	}
	return interpreter.env.declareConstant(stmt.GetName(), value)
}

func (interpreter *Interpreter) processEchoStatement(stmt ast.IEchoStatement) (IRuntimeValue, error) {
	for _, expr := range stmt.GetExpressions() {
		if runtimeValue, err := interpreter.process(expr); err != nil {
			return NewVoidRuntimeValue(), err
		} else {
			var str string
			str, err = runtimeValueToString(runtimeValue)
			if err != nil {
				return NewVoidRuntimeValue(), err
			}
			interpreter.print(str)
		}
	}
	return NewVoidRuntimeValue(), nil
}

func (interpreter *Interpreter) processSimpleVariableExpression(expr ast.ISimpleVariableExpression) (IRuntimeValue, error) {
	variableName, err := interpreter.varExprToVarName(expr)
	if err != nil {
		return NewVoidRuntimeValue(), err
	}

	return interpreter.env.lookupVariable(variableName)
}

func (interpreter *Interpreter) processSimpleAssignmentExpression(expr ast.ISimpleAssignmentExpression) (IRuntimeValue, error) {
	if !ast.IsVariableExpression(expr.GetVariable()) {
		return NewVoidRuntimeValue(),
			fmt.Errorf("Interpreter error: processSimpleAssignmentExpression: Invalid variable: %s", expr.GetVariable())
	}

	value, err := interpreter.process(expr.GetValue())
	if err != nil {
		return NewVoidRuntimeValue(), err
	}

	variableName, err := interpreter.varExprToVarName(expr.GetVariable())
	if err != nil {
		return NewVoidRuntimeValue(), err
	}

	return interpreter.env.declareVariable(variableName, value)
}

func (interpreter *Interpreter) processConstantAccessExpression(expr ast.IConstantAccessExpression) (IRuntimeValue, error) {
	return interpreter.env.lookupConstant(expr.GetConstantName())
}

func (interpreter *Interpreter) processCompoundAssignmentExpression(expr ast.ICompoundAssignmentExpression) (IRuntimeValue, error) {
	if !ast.IsVariableExpression(expr.GetVariable()) {
		return NewVoidRuntimeValue(),
			fmt.Errorf("Interpreter error: processCompoundAssignmentExpression: Invalid variable: %s", expr.GetVariable())
	}

	operand1, err := interpreter.process(expr.GetVariable())
	if err != nil {
		return NewVoidRuntimeValue(), err
	}

	operand2, err := interpreter.process(expr.GetValue())
	if err != nil {
		return NewVoidRuntimeValue(), err
	}

	newValue, err := calculate(operand1, expr.GetOperator(), operand2)
	if err != nil {
		return NewVoidRuntimeValue(), err
	}

	variableName, err := interpreter.varExprToVarName(expr.GetVariable())
	if err != nil {
		return NewVoidRuntimeValue(), err
	}

	return interpreter.env.declareVariable(variableName, newValue)
}

func (interpreter *Interpreter) processConditionalExpression(expr ast.IConditionalExpression) (IRuntimeValue, error) {
	isConditionTrue, err := interpreter.processCondition(expr.GetCondition())
	if err != nil {
		return NewVoidRuntimeValue(), err
	}

	if !isConditionTrue {
		return interpreter.process(expr.GetElseExpr())
	}
	if expr.GetIfExpr() == nil {
		return interpreter.process(expr.GetCondition())
	} else {
		return interpreter.process(expr.GetIfExpr())
	}
}
