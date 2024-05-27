package interpreter

import (
	"GoPHP/cmd/goPHP/ast"
	"math"
	"regexp"
	"slices"
	"strings"
)

func (interpreter *Interpreter) print(str string) {
	interpreter.result += str
}

func (interpreter *Interpreter) println(str string) {
	interpreter.print(str + "\n")
}

func (interpreter *Interpreter) processCondition(expr ast.IExpression, env *Environment) (IRuntimeValue, bool, Error) {
	runtimeValue, err := interpreter.processStmt(expr, env)
	if err != nil {
		return runtimeValue, false, err
	}

	boolean, err := lib_boolval(runtimeValue)
	return runtimeValue, boolean, err
}

func (interpreter *Interpreter) lookupVariable(expr ast.IExpression, env *Environment, suppressWarning bool) (IRuntimeValue, Error) {
	variableName, err := interpreter.varExprToVarName(expr, env)
	if err != nil {
		return NewVoidRuntimeValue(), err
	}

	runtimeValue, err := env.lookupVariable(variableName)
	if !suppressWarning && err != nil {
		interpreter.printError(err)
	}
	return runtimeValue, nil
}

// Convert a variable expression into the interpreted variable name
func (interpreter *Interpreter) varExprToVarName(expr ast.IExpression, env *Environment) (string, Error) {
	switch expr.GetKind() {
	case ast.SimpleVariableExpr:
		variableNameExpr := ast.ExprToSimpleVarExpr(expr).GetVariableName()

		if variableNameExpr.GetKind() == ast.VariableNameExpr {
			return ast.ExprToVarNameExpr(variableNameExpr).GetVariableName(), nil
		}

		if variableNameExpr.GetKind() == ast.SimpleVariableExpr {
			variableName, err := interpreter.varExprToVarName(variableNameExpr, env)
			if err != nil {
				return "", err
			}
			runtimeValue, err := env.lookupVariable(variableName)
			if err != nil {
				interpreter.printError(err)
			}
			valueStr, err := lib_strval(runtimeValue)
			if err != nil {
				return "", err
			}
			return "$" + valueStr, nil
		}

		return "", NewError("varExprToVarName - SimpleVariableExpr: Unsupported expression: %s", expr)
	case ast.SubscriptExpr:
		return interpreter.varExprToVarName(ast.ExprToSubscriptExpr(expr).GetVariable(), env)
	default:
		return "", NewError("varExprToVarName: Unsupported expression: %s", expr)
	}
}

func (interpreter *Interpreter) printError(err Error) {
	if err.GetErrorType() == WarningPhpError && interpreter.config.ErrorReporting&E_WARNING != 0 {
		interpreter.println("Warning: " + err.GetMessage())
	}
	// TODO implement
	// Depending on interpreter.config.ErrorReporting
}

// ------------------- MARK: RuntimeValue -------------------

func (interpreter *Interpreter) exprToRuntimeValue(expr ast.IExpression, env *Environment) (IRuntimeValue, Error) {
	switch expr.GetKind() {
	case ast.ArrayLiteralExpr:
		arrayRuntimeValue := NewArrayRuntimeValue()
		for _, key := range ast.ExprToArrayLitExpr(expr).GetKeys() {
			keyValue, err := interpreter.processStmt(key, env)
			if err != nil {
				return NewVoidRuntimeValue(), err
			}
			elementValue, err := interpreter.processStmt(ast.ExprToArrayLitExpr(expr).GetElements()[key], env)
			if err != nil {
				return NewVoidRuntimeValue(), err
			}
			arrayRuntimeValue.SetElement(keyValue, elementValue)
		}
		return arrayRuntimeValue, nil
	case ast.IntegerLiteralExpr:
		return NewIntegerRuntimeValue(ast.ExprToIntLitExpr(expr).GetValue()), nil
	case ast.FloatingLiteralExpr:
		return NewFloatingRuntimeValue(ast.ExprToFloatLitExpr(expr).GetValue()), nil
	case ast.StringLiteralExpr:
		str := ast.ExprToStrLitExpr(expr).GetValue()
		// variable substitution
		// TODO move to area where it is called before printing it
		if ast.ExprToStrLitExpr(expr).GetStringType() == ast.DoubleQuotedString {
			r, _ := regexp.Compile(`{\$[A-Za-z_][A-Za-z0-9_]*['A-Za-z0-9\[\]]*}`)
			matches := r.FindAllString(str, -1)
			for _, match := range matches {
				exprStr := "<?= " + match[1:len(match)-1] + ";"
				result, err := NewInterpreter(interpreter.config, interpreter.request).process(exprStr, env)
				if err != nil {
					return NewVoidRuntimeValue(), err
				}
				str = strings.Replace(str, match, result, 1)
			}
		}
		return NewStringRuntimeValue(str), nil
	default:
		return NewVoidRuntimeValue(), NewError("exprToRuntimeValue: Unsupported expression: %s", expr)
	}
}

func runtimeValueToValueType(valueType ValueType, runtimeValue IRuntimeValue) (IRuntimeValue, Error) {
	switch valueType {
	case BooleanValue:
		boolean, err := lib_boolval(runtimeValue)
		return NewBooleanRuntimeValue(boolean), err
	case FloatingValue:
		floating, err := lib_floatval(runtimeValue)
		return NewFloatingRuntimeValue(floating), err
	case IntegerValue:
		integer, err := lib_intval(runtimeValue)
		return NewIntegerRuntimeValue(integer), err
	case StringValue:
		str, err := lib_strval(runtimeValue)
		return NewStringRuntimeValue(str), err
	default:
		return NewVoidRuntimeValue(), NewError("runtimeValueToValueType: Unsupported runtime value: %s", valueType)
	}
}

// ------------------- MARK: unary-op-calculation -------------------

func calculateUnary(operator string, operand IRuntimeValue) (IRuntimeValue, Error) {
	switch operand.GetType() {
	case BooleanValue:
		return calculateUnaryBoolean(operator, runtimeValToBoolRuntimeVal(operand))
	case IntegerValue:
		return calculateUnaryInteger(operator, runtimeValToIntRuntimeVal(operand))
	case FloatingValue:
		return calculateUnaryFloating(operator, runtimeValToFloatRuntimeVal(operand))
	case NullValue:
		// Spec: https://phplang.org/spec/10-expressions.html#unary-arithmetic-operators
		// For a unary + or unary - operator used with a NULL-valued operand, the value of the result is zero and the type is int.
		return NewIntegerRuntimeValue(0), nil
	default:
		return NewVoidRuntimeValue(), NewError("calculateUnary: Type \"%s\" not implemented", operand.GetType())
	}

	// TODO calculateUnary - string
	// Spec: https://phplang.org/spec/10-expressions.html#unary-arithmetic-operators
	// For a unary + or - operator used with a numeric string or a leading-numeric string, the string is first converted to an int or float, as appropriate, after which it is handled as an arithmetic operand. The trailing non-numeric characters in leading-numeric strings are ignored. With a non-numeric string, the result has type int and value 0. If the string was leading-numeric or non-numeric, a non-fatal error MUST be produced.
	// For a unary ~ operator used with a string, the result is the string with each byte being bitwise complement of the corresponding byte of the source string.

	// TODO calculateUnary - object
	// If the operand has an object type supporting the operation, then the object semantics defines the result. Otherwise, for ~ the fatal error is issued and for + and - the object is converted to int.
}

func calculateUnaryBoolean(operator string, operand IBooleanRuntimeValue) (IIntegerRuntimeValue, Error) {
	switch operator {
	case "+":
		// Spec: https://phplang.org/spec/10-expressions.html#unary-arithmetic-operators
		// For a unary "+" operator used with a TRUE-valued operand, the value of the result is 1 and the type is int.
		// When used with a FALSE-valued operand, the value of the result is zero and the type is int.
		if runtimeValToBoolRuntimeVal(operand).GetValue() {
			return NewIntegerRuntimeValue(1), nil
		}
		return NewIntegerRuntimeValue(0), nil

	case "-":
		// Spec: https://phplang.org/spec/10-expressions.html#unary-arithmetic-operators
		// For a unary "-" operator used with a TRUE-valued operand, the value of the result is -1 and the type is int.
		// When used with a FALSE-valued operand, the value of the result is zero and the type is int.
		if runtimeValToBoolRuntimeVal(operand).GetValue() {
			return NewIntegerRuntimeValue(-1), nil
		}
		return NewIntegerRuntimeValue(0), nil

	default:
		return NewIntegerRuntimeValue(0), NewError("calculateUnaryBoolean: Operator \"%s\" not implemented", operator)
	}
}

func calculateUnaryFloating(operator string, operand IFloatingRuntimeValue) (IRuntimeValue, Error) {
	switch operator {
	case "+":
		// Spec: https://phplang.org/spec/10-expressions.html#unary-arithmetic-operators
		// For a unary "+" operator used with an arithmetic operand, the type and value of the result is the type and value of the operand.
		return operand, nil

	case "-":
		// Spec: https://phplang.org/spec/10-expressions.html#unary-arithmetic-operators
		// For a unary - operator used with an arithmetic operand, the value of the result is the negated value of the operand.
		// However, if an int operand’s original value is the smallest representable for that type,
		// the operand is treated as if it were float and the result will be float.
		return NewFloatingRuntimeValue(-runtimeValToFloatRuntimeVal(operand).GetValue()), nil

	case "~":
		// Spec: https://phplang.org/spec/10-expressions.html#unary-arithmetic-operators
		// For a unary ~ operator used with a float operand, the value of the operand is first converted to int before the bitwise complement is computed.
		intRuntimeValue, err := runtimeValueToValueType(IntegerValue, operand)
		if err != nil {
			return NewFloatingRuntimeValue(0), err
		}
		return calculateUnaryInteger(operator, runtimeValToIntRuntimeVal(intRuntimeValue))

	default:
		return NewFloatingRuntimeValue(0), NewError("calculateUnaryFloating: Operator \"%s\" not implemented", operator)
	}
}

func calculateUnaryInteger(operator string, operand IIntegerRuntimeValue) (IIntegerRuntimeValue, Error) {
	switch operator {
	case "+":
		// Spec: https://phplang.org/spec/10-expressions.html#unary-arithmetic-operators
		// For a unary "+" operator used with an arithmetic operand, the type and value of the result is the type and value of the operand.
		return operand, nil

	case "-":
		// Spec: https://phplang.org/spec/10-expressions.html#unary-arithmetic-operators
		// For a unary - operator used with an arithmetic operand, the value of the result is the negated value of the operand.
		// However, if an int operand’s original value is the smallest representable for that type,
		// the operand is treated as if it were float and the result will be float.
		return NewIntegerRuntimeValue(-runtimeValToIntRuntimeVal(operand).GetValue()), nil

	case "~":
		// Spec: https://phplang.org/spec/10-expressions.html#unary-arithmetic-operators
		// For a unary ~ operator used with an int operand, the type of the result is int.
		// The value of the result is the bitwise complement of the value of the operand
		// (that is, each bit in the result is set if and only if the corresponding bit in the operand is clear).
		return NewIntegerRuntimeValue(^runtimeValToIntRuntimeVal(operand).GetValue()), nil
	default:
		return NewIntegerRuntimeValue(0), NewError("calculateUnaryInteger: Operator \"%s\" not implemented", operator)
	}
}

// ------------------- MARK: binary-op-calculation -------------------

func calculate(operand1 IRuntimeValue, operator string, operand2 IRuntimeValue) (IRuntimeValue, Error) {
	resultType := VoidValue
	if slices.Contains([]string{"."}, operator) {
		resultType = StringValue
	} else if slices.Contains([]string{"&&", "||"}, operator) {
		resultType = BooleanValue
	} else if slices.Contains([]string{"&", "|", "^", "<<", ">>"}, operator) {
		resultType = IntegerValue
	} else {
		resultType = IntegerValue
		if operand1.GetType() == FloatingValue || operand2.GetType() == FloatingValue {
			resultType = FloatingValue
		}
	}

	var err Error
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
	case BooleanValue:
		return calculateBoolean(runtimeValToBoolRuntimeVal(operand1), operator, runtimeValToBoolRuntimeVal(operand2))
	case IntegerValue:
		return calculateInteger(runtimeValToIntRuntimeVal(operand1), operator, runtimeValToIntRuntimeVal(operand2))
	case FloatingValue:
		return calculateFloating(runtimeValToFloatRuntimeVal(operand1), operator, runtimeValToFloatRuntimeVal(operand2))
	case StringValue:
		return calculateString(runtimeValToStrRuntimeVal(operand1), operator, runtimeValToStrRuntimeVal(operand2))
	default:
		return NewVoidRuntimeValue(), NewError("calculate: Type \"%s\" not implemented", operator)
	}
}

func calculateBoolean(operand1 IBooleanRuntimeValue, operator string, operand2 IBooleanRuntimeValue) (IBooleanRuntimeValue, Error) {
	switch operator {
	case "&&":
		return NewBooleanRuntimeValue(operand1.GetValue() && operand2.GetValue()), nil
	case "||":
		return NewBooleanRuntimeValue(operand1.GetValue() || operand2.GetValue()), nil
	default:
		return NewBooleanRuntimeValue(false), NewError("calculateBoolean: Operator \"%s\" not implemented", operator)
	}
}

func calculateFloating(operand1 IFloatingRuntimeValue, operator string, operand2 IFloatingRuntimeValue) (IFloatingRuntimeValue, Error) {
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
		return NewFloatingRuntimeValue(math.Pow(operand1.GetValue(), operand2.GetValue())), nil
	default:
		return NewFloatingRuntimeValue(0), NewError("calculateInteger: Operator \"%s\" not implemented", operator)
	}
}

func calculateInteger(operand1 IIntegerRuntimeValue, operator string, operand2 IIntegerRuntimeValue) (IIntegerRuntimeValue, Error) {
	switch operator {
	case "<<":
		return NewIntegerRuntimeValue(operand1.GetValue() << operand2.GetValue()), nil
	case ">>":
		return NewIntegerRuntimeValue(operand1.GetValue() >> operand2.GetValue()), nil
	case "^":
		return NewIntegerRuntimeValue(operand1.GetValue() ^ operand2.GetValue()), nil
	case "|":
		return NewIntegerRuntimeValue(operand1.GetValue() | operand2.GetValue()), nil
	case "&":
		return NewIntegerRuntimeValue(operand1.GetValue() & operand2.GetValue()), nil
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
		return NewIntegerRuntimeValue(0), NewError("calculateInteger: Operator \"%s\" not implemented", operator)
	}
}

func calculateString(operand1 IStringRuntimeValue, operator string, operand2 IStringRuntimeValue) (IStringRuntimeValue, Error) {
	switch operator {
	case ".":
		return NewStringRuntimeValue(operand1.GetValue() + operand2.GetValue()), nil
	default:
		return NewStringRuntimeValue(""), NewError("calculateString: Operator \"%s\" not implemented", operator)
	}
}

// ------------------- MARK: comparison -------------------

func compare(lhs IRuntimeValue, operator string, rhs IRuntimeValue) (IBooleanRuntimeValue, Error) {
	// Spec: https://phplang.org/spec/10-expressions.html#grammar-equality-expression

	// TODO compare - "==", "!=", "<>"
	// Spec: https://phplang.org/spec/10-expressions.html#grammar-equality-expression
	// Operator == represents value equality, operators != and <> are equivalent and represent value inequality.
	// For operators ==, !=, and <>, the operands of different types are converted and compared according to the same rules as in relational operators. Two objects of different types are always not equal.

	// Spec: https://phplang.org/spec/10-expressions.html#grammar-equality-expression
	// Operator === represents same type and value equality, or identity, comparison,
	// and operator !== represents the opposite of ===.
	// The values are considered identical if they have the same type and compare as equal, with the additional conditions below:
	//    When comparing two objects, identity operators check to see if the two operands are the exact same object,
	//    not two different objects of the same type and value.
	//    Arrays must have the same elements in the same order to be considered identical.
	//    Strings are identical if they contain the same characters, unlike value comparison operators no conversions are performed for numeric strings.
	if operator == "===" || operator == "!==" {
		result := lhs.GetType() == rhs.GetType()
		if result {
			switch lhs.GetType() {
			case BooleanValue:
				result = runtimeValToBoolRuntimeVal(lhs).GetValue() == runtimeValToBoolRuntimeVal(rhs).GetValue()
			case FloatingValue:
				result = runtimeValToFloatRuntimeVal(lhs).GetValue() == runtimeValToFloatRuntimeVal(rhs).GetValue()
			case IntegerValue:
				result = runtimeValToIntRuntimeVal(lhs).GetValue() == runtimeValToIntRuntimeVal(rhs).GetValue()
			case NullValue:
				result = true
			case StringValue:
				result = runtimeValToStrRuntimeVal(lhs).GetValue() == runtimeValToStrRuntimeVal(rhs).GetValue()
			default:
				return NewBooleanRuntimeValue(false), NewError("compare: Runtime type %s for operator \"===\" not implemented", lhs.GetType())
			}
		}

		if operator == "!==" {
			return NewBooleanRuntimeValue(!result), nil
		} else {
			return NewBooleanRuntimeValue(result), nil
		}
	}

	return NewBooleanRuntimeValue(false), NewError("compare: Operator \"%s\" not implemented", operator)
}
