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

	boolean, err := runtimeValueToBool(runtimeValue)
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

func runtimeValueToString(runtimeValue IRuntimeValue) (string, error) {
	switch runtimeValue.GetType() {
	case IntegerValue:
		return fmt.Sprintf("%d", runtimeValToIntRuntimeVal(runtimeValue).GetValue()), nil
	case FloatingValue:
		return fmt.Sprintf("%g", runtimeValToFloatRuntimeVal(runtimeValue).GetValue()), nil
	case StringValue:
		return runtimeValToStrRuntimeVal(runtimeValue).GetValue(), nil
	default:
		return "", fmt.Errorf("exprToString: Unsupported runtime value: %s", runtimeValue.GetType())
	}
}

func runtimeValueToValueType(valueType ValueType, runtimeValue IRuntimeValue) (IRuntimeValue, error) {
	switch valueType {
	case BooleanValue:
		value, err := runtimeValueToBool(runtimeValue)
		if err != nil {
			return NewVoidRuntimeValue(), err
		}
		return NewBooleanRuntimeValue(value), nil
	case FloatingValue:
		value, err := runtimeValueToFloat(runtimeValue)
		if err != nil {
			return NewVoidRuntimeValue(), err
		}
		return NewFloatingRuntimeValue(value), nil
	case IntegerValue:
		value, err := runtimeValueToInt(runtimeValue)
		if err != nil {
			return NewVoidRuntimeValue(), err
		}
		return NewIntegerRuntimeValue(value), nil
	default:
		return NewVoidRuntimeValue(), fmt.Errorf("runtimeValueToValueType: Unsupported runtime value: %s", valueType)
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

func runtimeValueToFloat(runtimeValue IRuntimeValue) (float64, error) {
	// Spec: https://phplang.org/spec/08-conversions.html#converting-to-floating-point-type

	switch runtimeValue.GetType() {
	case FloatingValue:
		return runtimeValToFloatRuntimeVal(runtimeValue).GetValue(), nil
	case IntegerValue:
		// Spec: https://phplang.org/spec/08-conversions.html#converting-to-floating-point-type
		// If the source type is int,
		// if the precision can be preserved the result value is the closest approximation to the source value;
		// otherwise, the result is undefined.
		return float64(runtimeValToIntRuntimeVal(runtimeValue).GetValue()), nil
	default:
		// Spec: https://phplang.org/spec/08-conversions.html#converting-to-floating-point-type
		// For sources of all other types, the conversion result is obtained by first converting
		// the source value to int and then to float.
		intValue, err := runtimeValueToInt(runtimeValue)
		if err != nil {
			return 0, err
		}
		return runtimeValueToFloat(NewIntegerRuntimeValue(intValue))
	}

	// TODO runtimeValueToFloat - string
	// Spec: https://phplang.org/spec/08-conversions.html#converting-to-floating-point-type
	// If the source is a numeric string or leading-numeric string having integer format, the string’s integer value is treated as described above for a conversion from int. If the source is a numeric string or leading-numeric string having floating-point format, the result value is the closest approximation to the string’s floating-point value. The trailing non-numeric characters in leading-numeric strings are ignored. For any other string, the result value is 0.

	// TODO runtimeValueToFloat - object
	// Spec: https://phplang.org/spec/08-conversions.html#converting-to-floating-point-type
	// If the source is an object, if the class defines a conversion function, the result is determined by that function (this is currently available only to internal classes). If not, the conversion is invalid, the result is assumed to be 1.0 and a non-fatal error is produced.

}

func runtimeValueToInt(runtimeValue IRuntimeValue) (int64, error) {
	// Spec: https://phplang.org/spec/08-conversions.html#converting-to-integer-type

	switch runtimeValue.GetType() {
	case BooleanValue:
		// Spec: https://phplang.org/spec/08-conversions.html#converting-to-integer-type
		// If the source type is bool, then if the source value is FALSE, the result value is 0; otherwise, the result value is 1.
		if runtimeValToBoolRuntimeVal(runtimeValue).GetValue() {
			return 1, nil
		}
		return 0, nil
	case IntegerValue:
		return runtimeValToIntRuntimeVal(runtimeValue).GetValue(), nil
	case NullValue:
		// Spec: https://phplang.org/spec/08-conversions.html#converting-to-integer-type
		// If the source value is NULL, the result value is 0.
		return 0, nil
	default:
		return 0, fmt.Errorf("runtimeValueToInt: Unsupported runtime value %s", runtimeValue.GetType())
	}
	// TODO runtimeValueToInt - float
	// Spec: https://phplang.org/spec/08-conversions.html#converting-to-integer-type
	// If the source type is float, for the values INF, -INF, and NAN, the result value is zero. For all other values, if the precision can be preserved (that is, the float is within the range of an integer), the fractional part is rounded towards zero. If the precision cannot be preserved, the following conversion algorithm is used, where X is defined as two to the power of the number of bits in an integer (for example, 2 to the power of 32, i.e. 4294967296):
	// 1. We take the floating point remainder (wherein the remainder has the same sign as the dividend) of dividing the float by X, rounded towards zero.
	// 2. If the remainder is less than zero, it is rounded towards infinity and X is added.
	// 3. This result is converted to an unsigned integer.
	// 4. This result is converted to a signed integer by treating the unsigned integer as a two’s complement representation of the signed integer.
	// Implementations may implement this conversion differently (for example, on some architectures there may be hardware support for this specific conversion mode) so long as the result is the same.

	// TODO runtimeValueToInt - string
	// Spec: https://phplang.org/spec/08-conversions.html#converting-to-integer-type
	// If the source is a numeric string or leading-numeric string having integer format, if the precision can be preserved the result value is that string’s integer value; otherwise, the result is undefined. If the source is a numeric string or leading-numeric string having floating-point format, the string’s floating-point value is treated as described above for a conversion from float. The trailing non-numeric characters in leading-numeric strings are ignored. For any other string, the result value is 0.

	// TODO runtimeValueToInt - array
	// Spec: https://phplang.org/spec/08-conversions.html#converting-to-integer-type
	// If the source is an array with zero elements, the result value is 0; otherwise, the result value is 1.

	// TODO runtimeValueToInt - object
	// Spec: https://phplang.org/spec/08-conversions.html#converting-to-integer-type
	// If the source is an object, if the class defines a conversion function, the result is determined by that function (this is currently available only to internal classes). If not, the conversion is invalid, the result is assumed to be 1 and a non-fatal error is produced.

	// TODO runtimeValueToInt - resource
	// Spec: https://phplang.org/spec/08-conversions.html#converting-to-integer-type
	// If the source is a resource, the result is the resource’s unique ID.
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
