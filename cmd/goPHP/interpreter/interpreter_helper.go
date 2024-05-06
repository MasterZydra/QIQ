package interpreter

import (
	"GoPHP/cmd/goPHP/ast"
	"fmt"
	"math"
)

func (interpreter *Interpreter) print(str string) {
	interpreter.result += str
}

func exprToRuntimeValue(expr ast.IExpression) (IRuntimeValue, error) {
	switch expr.GetKind() {
	case ast.BooleanLiteralExpr:
		return NewBooleanRuntimeValue(ast.ExprToBoolLitExpr(expr).GetValue()), nil
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

func runtimeValueToBool(runtimeValue IRuntimeValue) (bool, error) {
	// Spec: https://phplang.org/spec/08-conversions.html#converting-to-boolean-type

	switch runtimeValue.GetType() {
	case BooleanValue:
		return runtimeValToBoolRuntimeVal(runtimeValue).GetValue(), nil
	case IntegerValue:
		// Spec: https://phplang.org/spec/08-conversions.html#converting-to-boolean-type
		// If the source type is int or float,
		// then if the source value tests equal to 0, the result value is FALSE; otherwise, the result value is TRUE.
		return runtimeValToIntRuntimeVal(runtimeValue).GetValue() != 0, nil
	case FloatingValue:
		// Spec: https://phplang.org/spec/08-conversions.html#converting-to-boolean-type
		// If the source type is int or float,
		// then if the source value tests equal to 0, the result value is FALSE; otherwise, the result value is TRUE.
		return math.Abs(runtimeValToFloatRuntimeVal(runtimeValue).GetValue()) != 0, nil
	case NullValue:
		// Spec: https://phplang.org/spec/08-conversions.html#converting-to-boolean-type
		// If the source value is NULL, the result value is FALSE.
		return false, nil
	case StringValue:
		// Spec: https://phplang.org/spec/08-conversions.html#converting-to-boolean-type
		// If the source is an empty string or the string “0”, the result value is FALSE; otherwise, the result value is TRUE.
		str := runtimeValToStrRuntimeVal(runtimeValue).GetValue()
		return str != "" && str != "0", nil
	default:
		return false, fmt.Errorf("runtimeValueToBool: Unsupported runtime value %s", runtimeValue.GetType())
	}
	// TODO runtimeValueToBool - array
	// Spec: https://phplang.org/spec/08-conversions.html#converting-to-boolean-type
	// If the source is an array with zero elements, the result value is FALSE; otherwise, the result value is TRUE.

	// TODO runtimeValueToBool - object
	// Spec: https://phplang.org/spec/08-conversions.html#converting-to-boolean-type
	// If the source is an object, the result value is TRUE.

	// TODO runtimeValueToBool - resource
	// Spec: https://phplang.org/spec/08-conversions.html#converting-to-boolean-type
	// If the source is a resource, the result value is TRUE.
}

// Convert a variable expression into the interpreted variable name
func (interpreter *Interpreter) varExprToVarName(expr ast.IExpression) (string, error) {
	switch expr.GetKind() {
	case ast.SimpleVariableExpr:
		variableNameExpr := ast.ExprToSimpleVarExpr(expr).GetVariableName()

		if variableNameExpr.GetKind() == ast.VariableNameExpr {
			return ast.ExprToVarNameExpr(variableNameExpr).GetVariableName(), nil
		}

		if variableNameExpr.GetKind() == ast.SimpleVariableExpr {
			variableName, err := interpreter.varExprToVarName(variableNameExpr)
			if err != nil {
				return "", err
			}
			runtimeValue, err := interpreter.env.lookupVariable(variableName)
			if err != nil {
				return "", err
			}
			valueStr, err := runtimeValueToString(runtimeValue)
			if err != nil {
				return "", err
			}
			return "$" + valueStr, nil
		}

		return "", fmt.Errorf("varExprToVarName - SimpleVariableExpr: Unsupported expression: %s", expr)
	default:
		return "", fmt.Errorf("varExprToVarName: Unsupported expression: %s", expr)
	}
}

func calculate(operand1 IRuntimeValue, operator string, operand2 IRuntimeValue) (IRuntimeValue, error) {
	// TODO implement all variable types and operators

	if operand1.GetType() != operand2.GetType() {
		// TODO implement type conversion
		return NewNullRuntimeValue(), fmt.Errorf("calculate: Operand types do not match: %s vs %s", operand1.GetType(), operand2.GetType())
	}

	switch operator {
	case "+":
		result := runtimeValToIntRuntimeVal(operand1).GetValue() + runtimeValToIntRuntimeVal(operand2).GetValue()
		return NewIntegerRuntimeValue(result), nil
	default:
		return NewNullRuntimeValue(), fmt.Errorf("calculate: Operator \"%s\" not implemented", operator)
	}
}

func (interpreter *Interpreter) processCondition(expr ast.IExpression) (bool, error) {
	runtimeValue, err := interpreter.process(expr)
	if err != nil {
		return false, err
	}

	return runtimeValueToBool(runtimeValue)
}
