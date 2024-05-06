package interpreter

import (
	"GoPHP/cmd/goPHP/ast"
	"fmt"
)

func (interpreter *Interpreter) print(str string) {
	interpreter.result += str
}

func exprToRuntimeValue(expr ast.IExpression) (IRuntimeValue, error) {
	switch expr.GetKind() {
	case ast.IntegerLiteralExpr:
		return NewIntegerRuntimeValue(ast.ExprToIntLitExpr(expr).GetValue()), nil
	case ast.FloatingLiteralExpr:
		return NewFloatingRuntimeValue(ast.ExprToFloatLitExpr(expr).GetValue()), nil
	case ast.StringLiteralExpr:
		return NewStringRuntimeValue(ast.ExprToStrLitExpr(expr).GetValue()), nil
	default:
		return NewVoidRuntimeValue(), fmt.Errorf("exprToRuntimeValue: Unsupported expression: %s", expr)
	}
}

func runtimeValueToString(runtimeValue IRuntimeValue) (string, error) {
	switch runtimeValue.GetType() {
	case IntegerValue:
		return fmt.Sprintf("%d", runtimeValToIntRuntimeVal(runtimeValue).GetValue()), nil
	case FloatingValue:
		return fmt.Sprintf("%f", runtimeValToFloatRuntimeVal(runtimeValue)), nil
	case StringValue:
		return runtimeValToStrRuntimeVal(runtimeValue).GetValue(), nil
	default:
		return "", fmt.Errorf("exprToString: Unsupported runtime value: %s", runtimeValue.GetType())
	}
}

// Convert a variable expression into the interpreted variable name
func (interpreter *Interpreter) varExprToVarName(expr ast.IExpression) (string, error) {
	switch expr.GetKind() {
	case ast.SimpleVariableExpr:
		variableNameExpr := ast.ExprToSimpleVarExpr(expr).GetVariableName()
		name, err := interpreter.varExprToVarName(variableNameExpr)
		if err != nil {
			return "", fmt.Errorf("varExprToVarName - SimpleVariableExpr: %s", err)
		}
		return name, nil
	case ast.VariableNameExpr:
		return ast.ExprToVarNameExpr(expr).GetVariableName(), nil
	default:
		return "", fmt.Errorf("varExprToVarName: Unsupported expression: %s", expr)
	}
}
