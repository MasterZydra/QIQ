package interpreter

import (
	"GoPHP/cmd/goPHP/ast"
	"fmt"
	"math"
	"slices"
)

func (interpreter *Interpreter) print(str string) {
	interpreter.result += str
}

func (interpreter *Interpreter) processCondition(expr ast.IExpression) (IRuntimeValue, bool, error) {
	runtimeValue, err := interpreter.process(expr)
	if err != nil {
		return runtimeValue, false, err
	}

	boolean, err := lib_boolval(runtimeValue)
	return runtimeValue, boolean, err
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
			valueStr, err := lib_strval(runtimeValue)
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

// ------------------- MARK: RuntimeValue -------------------

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
	case ast.NullLiteralExpr:
		return NewNullRuntimeValue(), nil
	default:
		return NewVoidRuntimeValue(), fmt.Errorf("exprToRuntimeValue: Unsupported expression: %s", expr)
	}
}

func runtimeValueToValueType(valueType ValueType, runtimeValue IRuntimeValue) (IRuntimeValue, error) {
	switch valueType {
	case BooleanValue:
		return nativeFn_boolval([]IRuntimeValue{runtimeValue}, nil)
	case FloatingValue:
		return nativeFn_floatval([]IRuntimeValue{runtimeValue}, nil)
	case IntegerValue:
		return nativeFn_intval([]IRuntimeValue{runtimeValue}, nil)
	case StringValue:
		return nativeFn_strval([]IRuntimeValue{runtimeValue}, nil)
	default:
		return NewVoidRuntimeValue(), fmt.Errorf("runtimeValueToValueType: Unsupported runtime value: %s", valueType)
	}
}

// ------------------- MARK: calculation -------------------

func calculate(operand1 IRuntimeValue, operator string, operand2 IRuntimeValue) (IRuntimeValue, error) {
	resultType := VoidValue
	if slices.Contains([]string{"."}, operator) {
		resultType = StringValue
	} else {
		resultType = IntegerValue
		if operand1.GetType() == FloatingValue || operand2.GetType() == FloatingValue {
			resultType = FloatingValue
		}
	}

	var err error
	operand1, err = runtimeValueToValueType(resultType, operand1)
	if err != nil {
		return NewVoidRuntimeValue(), err
	}
	operand2, err = runtimeValueToValueType(resultType, operand2)
	if err != nil {
		return NewVoidRuntimeValue(), err
	}
	// TODO testing how PHP behavious: var_dump(1.0 + 2); var_dump(1 + 2.0); var_dump("1" + 2);
	// var_dump("1" + "2"); => int
	// var_dump("1" . 2); => str
	// type order "string" - "int" - "float"

	// Testen
	//   true + 2
	//   true && 3

	switch resultType {
	case IntegerValue:
		return calculateInteger(runtimeValToIntRuntimeVal(operand1), operator, runtimeValToIntRuntimeVal(operand2))
	case FloatingValue:
		return calculateFloating(runtimeValToFloatRuntimeVal(operand1), operator, runtimeValToFloatRuntimeVal(operand2))
	default:
		return NewVoidRuntimeValue(), fmt.Errorf("calculate: Type \"%s\" not implemented", operator)
	}
}

func calculateFloating(operand1 IFloatingRuntimeValue, operator string, operand2 IFloatingRuntimeValue) (IFloatingRuntimeValue, error) {
	switch operator {
	case "+":
		return NewFloatingRuntimeValue(operand1.GetValue() + operand2.GetValue()), nil
	case "-":
		return NewFloatingRuntimeValue(operand1.GetValue() - operand2.GetValue()), nil
	case "*":
		return NewFloatingRuntimeValue(operand1.GetValue() * operand2.GetValue()), nil
	case "/":
		return NewFloatingRuntimeValue(operand1.GetValue() / operand2.GetValue()), nil
	case "**":
		return NewFloatingRuntimeValue(math.Pow(operand1.GetValue(), float64(operand2.GetValue()))), nil
	default:
		return NewFloatingRuntimeValue(0), fmt.Errorf("calculateInteger: Operator \"%s\" not implemented", operator)
	}
}

func calculateInteger(operand1 IIntegerRuntimeValue, operator string, operand2 IIntegerRuntimeValue) (IIntegerRuntimeValue, error) {
	switch operator {
	case "+":
		return NewIntegerRuntimeValue(operand1.GetValue() + operand2.GetValue()), nil
	case "-":
		return NewIntegerRuntimeValue(operand1.GetValue() - operand2.GetValue()), nil
	case "*":
		return NewIntegerRuntimeValue(operand1.GetValue() * operand2.GetValue()), nil
	case "/":
		return NewIntegerRuntimeValue(operand1.GetValue() / operand2.GetValue()), nil
	case "%":
		return NewIntegerRuntimeValue(operand1.GetValue() % operand2.GetValue()), nil
	case "**":
		return NewIntegerRuntimeValue(int64(math.Pow(float64(operand1.GetValue()), float64(operand2.GetValue())))), nil
	default:
		return NewIntegerRuntimeValue(0), fmt.Errorf("calculateInteger: Operator \"%s\" not implemented", operator)
	}
}
