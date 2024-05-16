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
	case ast.BooleanLiteralExpr, ast.IntegerLiteralExpr, ast.FloatingLiteralExpr, ast.StringLiteralExpr, ast.NullLiteralExpr:
		return exprToRuntimeValue(stmt)
	case ast.TextNode:
		interpreter.print(ast.ExprToTextExpr(stmt).GetValue())
		return NewVoidRuntimeValue(), nil
	case ast.SimpleVariableExpr:
		return interpreter.processSimpleVariableExpression(ast.ExprToSimpleVarExpr(stmt))
	case ast.SimpleAssignmentExpr:
		return interpreter.processSimpleAssignmentExpression(ast.ExprToSimpleAssignExpr(stmt))
	case ast.FunctionCallExpr:
		return interpreter.processFunctionCallExpression(ast.ExprToFuncCallExpr(stmt))
	case ast.ConstantAccessExpr:
		return interpreter.processConstantAccessExpression(ast.ExprToConstAccessExpr(stmt))
	case ast.CompoundAssignmentExpr:
		return interpreter.processCompoundAssignmentExpression(ast.ExprToCompoundAssignExpr(stmt))
	case ast.ConditionalExpr:
		return interpreter.processConditionalExpression(ast.ExprToCondExpr(stmt))
	case ast.CoalesceExpr:
		return interpreter.processCoalesceExpression(ast.ExprToCoalesceExpr(stmt))
	case ast.EqualityExpr:
		return interpreter.processEqualityExpression(ast.ExprToEqualExpr(stmt))
	case ast.AdditiveExpr:
		return interpreter.processAdditiveExpression(ast.ExprToEqualExpr(stmt))
	case ast.MultiplicativeExpr:
		return interpreter.processMultiplicativeExpression(ast.ExprToEqualExpr(stmt))
	case ast.LogicalNotExpr:
		return interpreter.processLogicalNotExpression(ast.ExprToUnaryOpExpr(stmt))

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
			str, err = lib_strval(runtimeValue)
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

func (interpreter *Interpreter) processFunctionCallExpression(expr ast.IFunctionCallExpression) (IRuntimeValue, error) {
	nativeFunction, err := interpreter.env.lookupNativeFunction(expr.GetFunctionName())
	if err != nil {
		return NewVoidRuntimeValue(), err
	}

	functionArguments := make([]IRuntimeValue, len(expr.GetArguments()))
	for index, arg := range expr.GetArguments() {
		runtimeValue, err := interpreter.process(arg)
		if err != nil {
			return NewVoidRuntimeValue(), err
		}
		functionArguments[index] = runtimeValue
	}
	return nativeFunction(functionArguments, interpreter.env)
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
	// Spec: https://phplang.org/spec/10-expressions.html#grammar-conditional-expression

	// Spec: https://phplang.org/spec/10-expressions.html#grammar-conditional-expression
	// Given the expression "e1 ? e2 : e3", e1 is evaluated first and converted to bool if it has another type.
	runtimeValue, isConditionTrue, err := interpreter.processCondition(expr.GetCondition())
	if err != nil {
		return NewVoidRuntimeValue(), err
	}

	if isConditionTrue {
		if expr.GetIfExpr() != nil {
			// Spec: https://phplang.org/spec/10-expressions.html#grammar-conditional-expression
			// If the result is TRUE, then and only then is e2 evaluated, and the result and its type become the result
			// and type of the whole expression.
			return interpreter.process(expr.GetIfExpr())
		} else {
			// Spec: https://phplang.org/spec/10-expressions.html#grammar-conditional-expression
			// There is a sequence point after the evaluation of e1.
			// If e2 is omitted, the result and type of the whole expression is the value
			// and type of e1 (before the conversion to bool).
			return runtimeValue, nil
		}
	} else {
		// Spec: https://phplang.org/spec/10-expressions.html#grammar-conditional-expression
		// Otherwise, then and only then is e3 evaluated, and the result and its type become the result
		// and type of the whole expression.
		return interpreter.process(expr.GetElseExpr())
	}
}

func (interpreter *Interpreter) processCoalesceExpression(expr ast.ICoalesceExpression) (IRuntimeValue, error) {
	// Spec: https://phplang.org/spec/10-expressions.html#grammar-coalesce-expression

	cond, err := interpreter.process(expr.GetCondition())
	if err != nil {
		return NewVoidRuntimeValue(), err
	}
	// Spec: https://phplang.org/spec/10-expressions.html#grammar-coalesce-expression
	// Given the expression e1 ?? e2, if e1 is set and not NULL (i.e. TRUE for isset), then the result is e1.
	if cond.GetType() != NullValue {
		return cond, nil
	}

	// Spec: https://phplang.org/spec/10-expressions.html#grammar-coalesce-expression
	// Otherwise, then and only then is e2 evaluated, and the result becomes the result of the whole expression.
	// There is a sequence point after the evaluation of e1.
	return interpreter.process(expr.GetElseExpr())

	// TODO processCoalesceExpression - handle uninitialized variables
	// Spec: https://phplang.org/spec/10-expressions.html#grammar-coalesce-expression
	// Note that the semantics of ?? is similar to isset so that uninitialized variables will not produce warnings when used in e1.
	// TODO use isset here
}

func (interpreter *Interpreter) processEqualityExpression(expr ast.IEqualityExpression) (IRuntimeValue, error) {
	lhs, err := interpreter.process(expr.GetLHS())
	if err != nil {
		return NewVoidRuntimeValue(), err
	}

	rhs, err := interpreter.process(expr.GetRHS())
	if err != nil {
		return NewVoidRuntimeValue(), err
	}
	return compare(lhs, expr.GetOperator(), rhs)
}

func (interpreter *Interpreter) processAdditiveExpression(expr ast.IEqualityExpression) (IRuntimeValue, error) {
	lhs, err := interpreter.process(expr.GetLHS())
	if err != nil {
		return NewVoidRuntimeValue(), err
	}

	rhs, err := interpreter.process(expr.GetRHS())
	if err != nil {
		return NewVoidRuntimeValue(), err
	}
	return calculate(lhs, expr.GetOperator(), rhs)
}

func (interpreter *Interpreter) processMultiplicativeExpression(expr ast.IEqualityExpression) (IRuntimeValue, error) {
	return interpreter.processAdditiveExpression(expr)
}

func (interpreter *Interpreter) processLogicalNotExpression(expr ast.IUnaryOpExpression) (IRuntimeValue, error) {
	runtimeValue, err := interpreter.process(expr.GetExpression())
	if err != nil {
		return NewVoidRuntimeValue(), err
	}
	boolValue, err := lib_boolval(runtimeValue)
	if err != nil {
		return NewVoidRuntimeValue(), err
	}
	return NewBooleanRuntimeValue(!boolValue), nil
}
