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
	case ast.ExpressionStmt:
		return interpreter.process(ast.StmtToExprStatement(stmt).GetExpression())
	case ast.EchoStmt:
		return interpreter.processEchoStatement(ast.StmtToEchoStatement(stmt))

		// Expressions
	case ast.IntegerLiteralExpr, ast.FloatingLiteralExpr, ast.StringLiteralExpr:
		return exprToRuntimeValue(stmt)
	case ast.TextNode:
		interpreter.print(ast.ExprToTextExpr(stmt).GetValue())
		return NewVoidRuntimeValue(), nil
	case ast.SimpleVariableExpr:
		return interpreter.processSimpleVariableExpression(ast.ExprToSimpleVarExpr(stmt))
	case ast.SimpleAssignmentExpr:
		return interpreter.processSimpleAssignmentExpression(ast.ExprToSimpleAssignExpr(stmt))

	default:
		return NewVoidRuntimeValue(), fmt.Errorf("Interpreter error: Unsupported statement or expression: %s", stmt)
	}
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
