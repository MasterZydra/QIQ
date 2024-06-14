package interpreter

import (
	"GoPHP/cmd/goPHP/ast"
	"GoPHP/cmd/goPHP/common"
	"GoPHP/cmd/goPHP/parser"
	"fmt"
	"math"
	"os"
	"regexp"
	"runtime"
	"slices"
	"strings"
)

func (interpreter *Interpreter) print(str string) {
	interpreter.result += str
}

var PHP_EOL string = ""

func (interpreter *Interpreter) println(str string) {
	if PHP_EOL == "" {
		str, _ := interpreter.env.lookupConstant("PHP_EOL")
		PHP_EOL = runtimeValToStrRuntimeVal(str).GetValue()
	}
	interpreter.print(str + PHP_EOL)
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

func (interpreter *Interpreter) ErrorToString(err Error) string {
	if err.GetErrorType() == WarningPhpError {
		if interpreter.config.ErrorReporting&E_WARNING != 0 {
			return "Warning: " + err.GetMessage()
		}
		return ""
	}
	if err.GetErrorType() == ErrorPhpError {
		if interpreter.config.ErrorReporting&E_ERROR != 0 {
			return "Fatal error: " + err.GetMessage()
		}
		return ""
	}
	return "Unsupported error type \"" + string(err.GetErrorType()) + "\": " + err.GetMessage()
}

func (interpreter *Interpreter) printError(err Error) {
	if errStr := interpreter.ErrorToString(err); errStr == "" {
		return
	} else {
		interpreter.println(errStr)
	}
}

func getPhpOs() string {
	switch runtime.GOOS {
	case "android":
		return "Android"
	case "darwin":
		return "Darwin"
	case "dragonfly":
		return "DragonFly"
	case "freebsd":
		return "FreeBSD"
	case "illumos":
		return "IllumOS"
	case "linux":
		return "Linux"
	case "netbsd":
		return "NetBSD"
	case "openbsd":
		return "OpenBSD"
	case "solaris":
		return "Solaris"
	case "windows":
		return "Windows"
	default:
		return "Unkown"
	}
}

func getPhpOsFamily() string {
	switch runtime.GOOS {
	case "android", "linux":
		return "Linux"
	case "darwin":
		return "Darwin"
	case "dragonfly", "freebsd", "netbsd", "openbsd":
		return "BSD"
	case "solaris":
		return "Solaris"
	case "windows":
		return "Windows"
	default:
		return "Unkown"
	}
}

// Scan and process program for function definitions on root level and in compound statements
func (interpreter *Interpreter) scanForFunctionDefinition(statements []ast.IStatement, env *Environment) Error {
	for _, stmt := range statements {
		if stmt.GetKind() == ast.CompoundStmt {
			interpreter.scanForFunctionDefinition(ast.StmtToCompoundStmt(stmt).GetStatements(), env)
			continue
		}

		if stmt.GetKind() != ast.FunctionDefinitionStmt {
			continue
		}

		_, err := interpreter.processFunctionDefinitionStmt(ast.StmtToFunctionDefinitionStmt(stmt), env)
		if err != nil {
			return err
		}
	}
	return nil
}

var paramTypeRuntimeValue = map[ValueType]string{
	ArrayValue:    "array",
	BooleanValue:  "bool",
	FloatingValue: "float",
	IntegerValue:  "int",
	NullValue:     "NULL",
	StringValue:   "string",
	VoidValue:     "void",
}

func checkParameterTypes(runtimeValue IRuntimeValue, expectedTypes []string) Error {
	typeStr, found := paramTypeRuntimeValue[runtimeValue.GetType()]
	if !found {
		return NewError("checkParameterTypes: No mapping for type %s", runtimeValue.GetType())
	}

	for _, expectedType := range expectedTypes {
		if expectedType == "mixed" {
			return nil
		}

		if typeStr == expectedType {
			return nil
		}
	}
	return NewError("Types do not match")
}

func (interpreter *Interpreter) includeFile(filepathExpr ast.IExpression, env *Environment, include bool, once bool) (IRuntimeValue, Error) {
	runtimeValue, err := interpreter.processStmt(filepathExpr, env)
	if err != nil {
		return runtimeValue, err
	}
	if runtimeValue.GetType() == NullValue {
		return runtimeValue, NewError("Uncaught ValueError: Path cannot be empty in %s", filepathExpr.GetPosition().ToPosString())
	}

	filename, err := lib_strval(runtimeValue)
	if err != nil {
		return runtimeValue, err
	}

	// Spec: https://phplang.org/spec/10-expressions.html#the-require-operator
	// Once an include file has been included, a subsequent use of require_once on that include file
	// results in a return value of TRUE but nothing else happens.
	if once && slices.Contains(interpreter.includedFiles, filename) {
		return NewBooleanRuntimeValue(true), nil
	}

	absFilename := filename
	if !common.IsAbsPath(filename) {
		absFilename = common.GetAbsPathForWorkingDir(common.ExtractPath(filename), filename)
	}

	var functionName string
	if include {
		functionName = "include"
	} else {
		functionName = "require"
	}

	// Spec: https://phplang.org/spec/10-expressions.html#the-require-operator
	// This operator is identical to operator include except that in the case of require,
	// failure to find/open the designated include file produces a fatal error.
	getError := func() (IRuntimeValue, Error) {
		if include {
			return NewVoidRuntimeValue(), NewWarning(
				"%s(): Failed opening '%s' for inclusion (include_path='%s') in %s",
				functionName, filename, common.ExtractPath(filepathExpr.GetPosition().Filename), filepathExpr.GetPosition().ToPosString(),
			)
		} else {
			return NewVoidRuntimeValue(), NewError(
				"Uncaught Error: Failed opening required '%s' (include_path='%s') in %s",
				filename, common.ExtractPath(filepathExpr.GetPosition().Filename), filepathExpr.GetPosition().ToPosString(),
			)
		}
	}

	if !common.PathExists(absFilename) {
		interpreter.printError(NewWarning(
			"%s(%s): Failed to open stream: No such file or directory in %s",
			functionName, filename, filepathExpr.GetPosition().ToPosString(),
		))
		return getError()
	}

	content, fileErr := os.ReadFile(filename)
	if fileErr != nil {
		return getError()
	}
	program, parserErr := parser.NewParser().ProduceAST(string(content), filename)
	fmt.Println(program)
	if parserErr != nil {
		return runtimeValue, NewParseError(parserErr)
	}
	return interpreter.processProgram(program, env)
}

// ------------------- MARK: Caching -------------------

func (interpreter *Interpreter) isCached(stmt ast.IStatement) bool {
	_, found := interpreter.cache[stmt.GetId()]
	return found
}

func (interpreter *Interpreter) writeCache(stmt ast.IStatement, value IRuntimeValue) IRuntimeValue {
	interpreter.cache[stmt.GetId()] = value
	return value
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
				result, err := NewInterpreter(interpreter.config, interpreter.request, "").process(exprStr, env)
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

func deepCopy(value IRuntimeValue) IRuntimeValue {
	if value.GetType() != ArrayValue {
		return value
	}

	copy := NewArrayRuntimeValue()
	array := runtimeValToArrayRuntimeVal(value)
	for _, key := range array.GetKeys() {
		copy.SetElement(key, deepCopy(array.GetElements()[key]))
	}
	return copy
}

// ------------------- MARK: inc-dec-calculation -------------------

func calculateIncDec(operator string, operand IRuntimeValue) (IRuntimeValue, Error) {
	switch operand.GetType() {
	case BooleanValue:
		// Spec: https://phplang.org/spec/10-expressions.html#prefix-increment-and-decrement-operators
		// For a prefix ++ or -- operator used with a Boolean-valued operand, there is no side effect, and the result is the operand’s value.
		return operand, nil
	case FloatingValue:
		return calculateIncDecFloating(operator, runtimeValToFloatRuntimeVal(operand))
	case IntegerValue:
		return calculateIncDecInteger(operator, runtimeValToIntRuntimeVal(operand))
	case NullValue:
		return calculateIncDecNull(operator)
	case StringValue:
		return calculateIncDecString(operator, runtimeValToStrRuntimeVal(operand))
	default:
		return NewVoidRuntimeValue(), NewError("calculateIncDec: Type \"%s\" not implemented", operand.GetType())
	}

	// TODO calculateIncDec - object
	// Spec: https://phplang.org/spec/10-expressions.html#prefix-increment-and-decrement-operators
	// If the operand has an object type supporting the operation, then the object semantics defines the result. Otherwise, the operation has no effect and the result is the operand.
}

func calculateIncDecInteger(operator string, operand IIntegerRuntimeValue) (IRuntimeValue, Error) {
	switch operator {
	case "++":
		// Spec: https://phplang.org/spec/10-expressions.html#prefix-increment-and-decrement-operators
		//For a prefix "++" operator used with an arithmetic operand, the side effect of the operator is to increment the value of the operand by 1.
		// The result is the value of the operand after it has been incremented.
		// If an int operand’s value is the largest representable for that type, the operand is incremented as if it were float.

		// Spec: https://phplang.org/spec/10-expressions.html#prefix-increment-and-decrement-operators
		// For a prefix ++ or -- operator used with an operand having the value INF, -INF, or NAN, there is no side effect, and the result is the operand’s value.
		return calculateInteger(operand, "+", NewIntegerRuntimeValue(1))

	case "--":
		// Spec: https://phplang.org/spec/10-expressions.html#prefix-increment-and-decrement-operators
		// For a prefix "--" operator used with an arithmetic operand, the side effect of the operator is to decrement the value of the operand by 1.
		// The result is the value of the operand after it has been decremented.
		// If an int operand’s value is the smallest representable for that type, the operand is decremented as if it were float.

		// Spec: https://phplang.org/spec/10-expressions.html#prefix-increment-and-decrement-operators
		// For a prefix ++ or -- operator used with an operand having the value INF, -INF, or NAN, there is no side effect, and the result is the operand’s value.
		return calculateInteger(operand, "-", NewIntegerRuntimeValue(1))

	default:
		return NewIntegerRuntimeValue(0), NewError("calculateIncDecInteger: Operator \"%s\" not implemented", operator)
	}
}

func calculateIncDecFloating(operator string, operand IFloatingRuntimeValue) (IRuntimeValue, Error) {
	switch operator {
	case "++":
		// Spec: https://phplang.org/spec/10-expressions.html#prefix-increment-and-decrement-operators
		//For a prefix "++" operator used with an arithmetic operand, the side effect of the operator is to increment the value of the operand by 1.
		// The result is the value of the operand after it has been incremented.
		// If an int operand’s value is the largest representable for that type, the operand is incremented as if it were float.

		// Spec: https://phplang.org/spec/10-expressions.html#prefix-increment-and-decrement-operators
		// For a prefix ++ or -- operator used with an operand having the value INF, -INF, or NAN, there is no side effect, and the result is the operand’s value.
		return calculateFloating(operand, "+", NewFloatingRuntimeValue(1))

	case "--":
		// Spec: https://phplang.org/spec/10-expressions.html#prefix-increment-and-decrement-operators
		// For a prefix "--" operator used with an arithmetic operand, the side effect of the operator is to decrement the value of the operand by 1.
		// The result is the value of the operand after it has been decremented.
		// If an int operand’s value is the smallest representable for that type, the operand is decremented as if it were float.

		// Spec: https://phplang.org/spec/10-expressions.html#prefix-increment-and-decrement-operators
		// For a prefix ++ or -- operator used with an operand having the value INF, -INF, or NAN, there is no side effect, and the result is the operand’s value.
		return calculateFloating(operand, "-", NewFloatingRuntimeValue(1))

	default:
		return NewIntegerRuntimeValue(0), NewError("calculateIncDecFloating: Operator \"%s\" not implemented", operator)
	}
}

func calculateIncDecNull(operator string) (IRuntimeValue, Error) {
	switch operator {
	case "++":
		// Spec: https://phplang.org/spec/10-expressions.html#prefix-increment-and-decrement-operators
		// For a prefix ++ operator used with a NULL-valued operand, the side effect is that the operand’s type is changed to int,
		// the operand’s value is set to zero, and that value is incremented by 1.
		// The result is the value of the operand after it has been incremented.
		return NewIntegerRuntimeValue(1), nil

	case "--":
		// Spec: https://phplang.org/spec/10-expressions.html#prefix-increment-and-decrement-operators
		// For a prefix – operator used with a NULL-valued operand, there is no side effect, and the result is the operand’s value.
		return NewNullRuntimeValue(), nil

	default:
		return NewIntegerRuntimeValue(0), NewError("calculateIncDecNull: Operator \"%s\" not implemented", operator)
	}
}

func calculateIncDecString(operator string, operand IStringRuntimeValue) (IRuntimeValue, Error) {
	switch operator {
	case "++":
		// Spec: https://phplang.org/spec/10-expressions.html#prefix-increment-and-decrement-operators
		// For a prefix "++" operator used with an operand whose value is an empty string,
		// the side effect is that the operand’s value is changed to the string “1”. The type of the operand is unchanged.
		// The result is the new value of the operand.
		if runtimeValToStrRuntimeVal(operand).GetValue() == "" {
			return NewStringRuntimeValue("1"), nil
		}
		return NewVoidRuntimeValue(), NewError("TODO calculateIncDecString")

	case "--":
		// Spec: https://phplang.org/spec/10-expressions.html#prefix-increment-and-decrement-operators
		// For a prefix "--" operator used with an operand whose value is an empty string,
		// the side effect is that the operand’s type is changed to int, the operand’s value is set to zero,
		// and that value is decremented by 1. The result is the value of the operand after it has been incremented.
		if runtimeValToStrRuntimeVal(operand).GetValue() == "" {
			return NewIntegerRuntimeValue(-1), nil
		}
		return NewVoidRuntimeValue(), NewError("TODO calculateIncDecString")

	default:
		return NewIntegerRuntimeValue(0), NewError("calculateIncDecNull: Operator \"%s\" not implemented", operator)
	}

	// TODO calculateIncDecString
	// Spec: https://phplang.org/spec/10-expressions.html#prefix-increment-and-decrement-operators
	/*
		String Operands

		For a prefix -- or ++ operator used with a numeric string, the numeric string is treated as the corresponding int or float value.

		For a prefix -- operator used with a non-numeric string-valued operand, there is no side effect, and the result is the operand’s value.

		For a non-numeric string-valued operand that contains only alphanumeric characters, for a prefix ++ operator, the operand is considered to be a representation of a base-36 number (i.e., with digits 0–9 followed by A–Z or a–z) in which letter case is ignored for value purposes. The right-most digit is incremented by 1. For the digits 0–8, that means going to 1–9. For the letters “A”–“Y” (or “a”–“y”), that means going to “B”–“Z” (or “b”–“z”). For the digit 9, the digit becomes 0, and the carry is added to the next left-most digit, and so on. For the digit “Z” (or “z”), the resulting string has an extra digit “A” (or “a”) appended. For example, when incrementing, “a” -> “b”, “Z” -> “AA”, “AA” -> “AB”, “F29” -> “F30”, “FZ9” -> “GA0”, and “ZZ9” -> “AAA0”. A digit position containing a number wraps modulo-10, while a digit position containing a letter wraps modulo-26.

		For a non-numeric string-valued operand that contains any non-alphanumeric characters, for a prefix ++ operator, all characters up to and including the right-most non-alphanumeric character is passed through to the resulting string, unchanged. Characters to the right of that right-most non-alphanumeric character are treated like a non-numeric string-valued operand that contains only alphanumeric characters, except that the resulting string will not be extended. Instead, a digit position containing a number wraps modulo-10, while a digit position containing a letter wraps modulo-26.
	*/
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
		return NewVoidRuntimeValue(), NewError("calculate: Type \"%s\" not implemented", resultType)
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

// ------------------- MARK: compareRelation -------------------

func compareRelation(lhs IRuntimeValue, operator string, rhs IRuntimeValue) (IRuntimeValue, Error) {
	// Spec: https://phplang.org/spec/10-expressions.html#grammar-relational-expression

	// Note that greater-than semantics is implemented as the reverse of less-than, i.e. "$a > $b" is the same as "$b < $a".
	// This may lead to confusing results if the operands are not well-ordered
	// - such as comparing two objects not having comparison semantics, or comparing arrays.

	// Operator "<=>" represents comparison operator between two expressions,
	// with the result being an integer less than "0" if the expression on the left is less than the expression on the right
	// (i.e. if "$a < $b" would return "TRUE"), as defined below by the semantics of the operator "<",
	// integer "0" if those expressions are equal (as defined by the semantics of the == operator) and
	// integer greater than 0 otherwise.

	// Operator "<" represents less-than, operator ">" represents greater-than, operator "<=" represents less-than-or-equal-to,
	// and operator ">=" represents greater-than-or-equal-to. The type of the result is bool.

	// The following table shows the result for comparison of different types, with the left operand displayed vertically
	// and the right displayed horizontally. The conversions are performed according to type conversion rules.

	// See in compareRelation[Type] ...

	// "<" means that the left operand is always less than the right operand.
	// ">" means that the left operand is always greater than the right operand.
	// "->" means that the left operand is converted to the type of the right operand.
	// "<-" means that the right operand is converted to the type of the left operand.

	// A number means one of the cases below:
	//   2. If one of the operands has arithmetic type, is a resource, or a numeric string,
	//      which can be represented as int or float without loss of precision,
	//      the operands are converted to the corresponding arithmetic type, with float taking precedence over int,
	//      and resources converting to int. The result is the numerical comparison of the two operands after conversion.
	//
	//   3. If only one operand has object type, if the object has comparison handler, that handler defines the result.
	//      Otherwise, if the object can be converted to the other operand’s type, it is converted and the result is used for the comparison.
	//      Otherwise, the object compares greater-than any other operand type.
	//
	//   4. If both operands are non-numeric strings, the result is the lexical comparison of the two operands.
	//      Specifically, the strings are compared byte-by-byte starting with their first byte.
	//      If the two bytes compare equal and there are no more bytes in either string, the strings are equal and the comparison ends;
	//      otherwise, if this is the final byte in one string, the shorter string compares less-than the longer string and the comparison ends.
	//      If the two bytes compare unequal, the string having the lower-valued byte compares less-than the other string, and the comparison ends.
	//      If there are more bytes in the strings, the process is repeated for the next pair of bytes.
	//
	//   6. When comparing two objects, if any of the object types has its own compare semantics, that would define the result,
	//      with the left operand taking precedence. Otherwise, if the objects are of different types, the comparison result is FALSE.
	//      If the objects are of the same type, the properties of the objects are compares using the array comparison described above.

	// Reduce code complexity and duplication by only implementing less-than and less-than-or-equal-to
	switch operator {
	case ">":
		return compareRelation(rhs, "<", lhs)
	case ">=":
		return compareRelation(rhs, "<=", lhs)
	}

	switch lhs.GetType() {
	case ArrayValue:
		return compareRelationArray(runtimeValToArrayRuntimeVal(lhs), operator, rhs)
	case BooleanValue:
		return compareRelationBoolean(runtimeValToBoolRuntimeVal(lhs), operator, rhs)
	case FloatingValue:
		return compareRelationFloating(runtimeValToFloatRuntimeVal(lhs), operator, rhs)
	case IntegerValue:
		return compareRelationInteger(runtimeValToIntRuntimeVal(lhs), operator, rhs)
	case StringValue:
		return compareRelationString(runtimeValToStrRuntimeVal(lhs), operator, rhs)
	case NullValue:
		return compareRelationNull(operator, rhs)
	default:
		return NewVoidRuntimeValue(), NewError("compareRelation: Type \"%s\" not implemented", lhs.GetType())
	}

}

func compareRelationArray(lhs IArrayRuntimeValue, operator string, rhs IRuntimeValue) (IRuntimeValue, Error) {
	// Spec: https://phplang.org/spec/10-expressions.html#grammar-relational-expression

	//        NULL  bool  int  float  string  array  object  resource
	// array   <-    ->    >    >      >       5      3       >

	//   5. If both operands have array type, if the arrays have different numbers of elements,
	//      the one with the fewer is considered less-than the other one, regardless of the keys and values in each, and the comparison ends.
	//      For arrays having the same numbers of elements, the keys from the left operand are considered one by one,
	//      if the next key in the left-hand operand exists in the right-hand operand, the corresponding values are compared.
	//      If they are unequal, the array containing the lesser value is considered less-than the other one, and the comparison ends;
	//      otherwise, the process is repeated with the next element.
	//      If the next key in the left-hand operand does not exist in the right-hand operand, the arrays cannot be compared and FALSE is returned.
	//      If all the values are equal, then the arrays are considered equal.

	// TODO compareRelationArray - object
	// TODO compareRelationArray - resource

	if rhs.GetType() == NullValue {
		var err Error
		rhs, err = lib_arrayval(rhs)
		if err != nil {
			return NewVoidRuntimeValue(), err
		}
	}

	switch rhs.GetType() {
	case ArrayValue:
		rhsArray := runtimeValToArrayRuntimeVal(rhs)
		var result int64 = 0
		if len(lhs.GetKeys()) != len(rhsArray.GetKeys()) {
			if len(lhs.GetKeys()) < len(rhsArray.GetKeys()) {
				result = -1
			} else {
				result = 1
			}
		} else {
			for _, key := range lhs.GetKeys() {
				lhsValue, _ := lhs.GetElement(key)
				rhsValue, found := rhsArray.GetElement(key)
				if found {
					equal, err := compare(lhsValue, "===", rhsValue)
					if err != nil {
						return NewVoidRuntimeValue(), err
					}
					if runtimeValToBoolRuntimeVal(equal).GetValue() {
						continue
					}
					lessThan, err := compareRelation(lhsValue, operator, rhsValue)
					if err != nil {
						return NewVoidRuntimeValue(), err
					}
					if lessThan.GetType() == BooleanValue {
						if runtimeValToBoolRuntimeVal(lessThan).GetValue() {
							result = -1
						} else {
							result = 1
						}
					}
					if lessThan.GetType() == IntegerValue {
						result = runtimeValToIntRuntimeVal(lessThan).GetValue()
					}
				}
			}
		}

		switch operator {
		case "<":
			return NewBooleanRuntimeValue(result == -1), nil
		case "<=":
			return NewBooleanRuntimeValue(result < 1), nil
		case "<=>":
			return NewIntegerRuntimeValue(result), nil
		default:
			return NewVoidRuntimeValue(), NewError("compareRelationArray: Operator \"%s\" not implemented", operator)
		}

	case BooleanValue:
		lhsBoolean, err := lib_boolval(lhs)
		if err != nil {
			return NewVoidRuntimeValue(), err
		}
		return compareRelationBoolean(NewBooleanRuntimeValue(lhsBoolean), operator, rhs)

	case FloatingValue, IntegerValue, StringValue:
		switch operator {
		case "<", "<=":
			return NewBooleanRuntimeValue(false), nil
		case "<=>":
			return NewIntegerRuntimeValue(1), nil
		default:
			return NewVoidRuntimeValue(), NewError("compareRelationArray: Operator \"%s\" not implemented", operator)
		}

	default:
		return NewVoidRuntimeValue(), NewError("compareRelationArray: Type \"%s\" not implemented", rhs.GetType())
	}
}

func compareRelationBoolean(lhs IBooleanRuntimeValue, operator string, rhs IRuntimeValue) (IRuntimeValue, Error) {
	// Spec: https://phplang.org/spec/10-expressions.html#grammar-relational-expression

	//       NULL  bool  int  float  string  array  object  resource
	// bool   <-    1     <-   <-     <-      <-     <-      <-

	//   1. If either operand has type bool, the other operand is converted to that type.
	//      The result is the logical comparison of the two operands after conversion, where FALSE is defined to be less than TRUE.

	rhsBoolean, err := lib_boolval(rhs)
	if err != nil {
		return NewVoidRuntimeValue(), err
	}
	// TODO compareRelationBoolean - object - implement in lib_boolval
	// TODO compareRelationBoolean - resource - implement in lib_boolval

	lhsInt, err := lib_intval(lhs)
	if err != nil {
		return NewVoidRuntimeValue(), err
	}
	rhsInt, err := lib_intval(NewBooleanRuntimeValue(rhsBoolean))
	if err != nil {
		return NewVoidRuntimeValue(), err
	}

	switch operator {
	case "<":
		return NewBooleanRuntimeValue(lhsInt < rhsInt), nil

	case "<=":
		return NewBooleanRuntimeValue(lhsInt <= rhsInt), nil

	case "<=>":
		if lhsInt > rhsInt {
			return NewIntegerRuntimeValue(1), nil
		}
		if lhsInt == rhsInt {
			return NewIntegerRuntimeValue(0), nil
		}
		return NewIntegerRuntimeValue(-1), nil

	default:
		return NewVoidRuntimeValue(), NewError("compareRelationBoolean: Operator \"%s\" not implemented", operator)
	}
}

func compareRelationFloating(lhs IFloatingRuntimeValue, operator string, rhs IRuntimeValue) (IRuntimeValue, Error) {
	// Spec: https://phplang.org/spec/10-expressions.html#grammar-relational-expression

	//        NULL  bool  int  float  string  array  object  resource
	// float   <-    ->    2    2      <-      <      3       <-

	// TODO compareRelationFloating - object
	// TODO compareRelationFloating - resource

	if rhs.GetType() == StringValue {
		rhsStr := runtimeValToStrRuntimeVal(rhs).GetValue()
		if strings.Trim(rhsStr, " \t") == "" {
			switch operator {
			case "<", "<=":
				return NewBooleanRuntimeValue(false), nil
			case "<=>":
				return NewIntegerRuntimeValue(1), nil
			default:
				return NewVoidRuntimeValue(), NewError("compareRelationFloating: Operator \"%s\" not implemented for type string", operator)
			}
		}
		if !common.IsIntegerLiteralWithSign(rhsStr) && !common.IsFloatingLiteralWithSign(rhsStr) {
			switch operator {
			case "<", "<=":
				return NewBooleanRuntimeValue(true), nil
			case "<=>":
				return NewIntegerRuntimeValue(-1), nil
			default:
				return NewVoidRuntimeValue(), NewError("compareRelationFloating: Operator \"%s\" not implemented for type string", operator)
			}
		}
	}

	if rhs.GetType() == NullValue || rhs.GetType() == IntegerValue || rhs.GetType() == StringValue {
		var err Error
		rhs, err = runtimeValueToValueType(FloatingValue, rhs)
		if err != nil {
			return NewVoidRuntimeValue(), err
		}
	}

	switch rhs.GetType() {
	case ArrayValue:
		switch operator {
		case "<", "<=":
			return NewBooleanRuntimeValue(true), nil
		case "<=>":
			return NewIntegerRuntimeValue(-1), nil
		default:
			return NewVoidRuntimeValue(), NewError("compareRelationFloating: Operator \"%s\" not implemented for type array", operator)
		}

	case BooleanValue:
		lhsBoolean, err := lib_boolval(lhs)
		if err != nil {
			return NewVoidRuntimeValue(), err
		}
		return compareRelationBoolean(NewBooleanRuntimeValue(lhsBoolean), operator, rhs)

	case FloatingValue:
		lhsFloat := runtimeValToFloatRuntimeVal(lhs).GetValue()
		rhsFloat := runtimeValToFloatRuntimeVal(rhs).GetValue()
		switch operator {
		case "<":
			return NewBooleanRuntimeValue(lhsFloat < rhsFloat), nil
		case "<=":
			return NewBooleanRuntimeValue(lhsFloat <= rhsFloat), nil
		case "<=>":
			if lhsFloat > rhsFloat {
				return NewIntegerRuntimeValue(1), nil
			}
			if lhsFloat == rhsFloat {
				return NewIntegerRuntimeValue(0), nil
			}
			return NewIntegerRuntimeValue(-1), nil
		default:
			return NewVoidRuntimeValue(), NewError("compareRelationFloating: Operator \"%s\" not implemented", operator)
		}

	default:
		return NewVoidRuntimeValue(), NewError("compareRelationFloating: Type \"%s\" not implemented", rhs.GetType())
	}
}

func compareRelationInteger(lhs IIntegerRuntimeValue, operator string, rhs IRuntimeValue) (IRuntimeValue, Error) {
	// Spec: https://phplang.org/spec/10-expressions.html#grammar-relational-expression

	//      NULL  bool  int  float  string  array  object  resource
	// int   <-    ->    2    2      <-      <      3       <-

	// TODO compareRelationInteger - object
	// TODO compareRelationInteger - resource

	if rhs.GetType() == StringValue {
		rhsStr := runtimeValToStrRuntimeVal(rhs).GetValue()
		if strings.Trim(rhsStr, " \t") == "" {
			switch operator {
			case "<", "<=":
				return NewBooleanRuntimeValue(false), nil
			case "<=>":
				return NewIntegerRuntimeValue(1), nil
			default:
				return NewVoidRuntimeValue(), NewError("compareRelationInteger: Operator \"%s\" not implemented for type array", operator)
			}
		}
		if !common.IsIntegerLiteralWithSign(rhsStr) && !common.IsFloatingLiteralWithSign(rhsStr) {
			switch operator {
			case "<", "<=":
				return NewBooleanRuntimeValue(true), nil
			case "<=>":
				return NewIntegerRuntimeValue(-1), nil
			default:
				return NewVoidRuntimeValue(), NewError("compareRelationInteger: Operator \"%s\" not implemented for type array", operator)
			}
		}
	}

	if rhs.GetType() == NullValue || rhs.GetType() == StringValue {
		var err Error
		rhs, err = runtimeValueToValueType(IntegerValue, rhs)
		if err != nil {
			return NewVoidRuntimeValue(), err
		}
	}

	switch rhs.GetType() {
	case ArrayValue:
		switch operator {
		case "<", "<=":
			return NewBooleanRuntimeValue(true), nil
		case "<=>":
			return NewIntegerRuntimeValue(-1), nil
		default:
			return NewVoidRuntimeValue(), NewError("compareRelationInteger: Operator \"%s\" not implemented for type array", operator)
		}

	case BooleanValue:
		lhsBoolean, err := lib_boolval(lhs)
		if err != nil {
			return NewVoidRuntimeValue(), err
		}
		return compareRelationBoolean(NewBooleanRuntimeValue(lhsBoolean), operator, rhs)

	case FloatingValue:
		lhsFloat, err := lib_floatval(lhs)
		if err != nil {
			return NewVoidRuntimeValue(), err
		}
		return compareRelationFloating(NewFloatingRuntimeValue(lhsFloat), operator, rhs)

	case IntegerValue:
		lhsInt := runtimeValToIntRuntimeVal(lhs).GetValue()
		rhsInt := runtimeValToIntRuntimeVal(rhs).GetValue()
		switch operator {
		case "<":
			return NewBooleanRuntimeValue(lhsInt < rhsInt), nil
		case "<=":
			return NewBooleanRuntimeValue(lhsInt <= rhsInt), nil
		case "<=>":
			if lhsInt > rhsInt {
				return NewIntegerRuntimeValue(1), nil
			}
			if lhsInt == rhsInt {
				return NewIntegerRuntimeValue(0), nil
			}
			return NewIntegerRuntimeValue(-1), nil
		default:
			return NewVoidRuntimeValue(), NewError("compareRelationInteger: Operator \"%s\" not implemented", operator)
		}

	default:
		return NewVoidRuntimeValue(), NewError("compareRelationInteger: Type \"%s\" not implemented", rhs.GetType())
	}
}

func compareRelationNull(operator string, rhs IRuntimeValue) (IRuntimeValue, Error) {
	// Spec: https://phplang.org/spec/10-expressions.html#grammar-relational-expression

	//       NULL  bool  int  float  string  array  object  resource
	// NULL   =     ->    ->   ->     ->      ->     <       <

	// "=" means the result is always “equals”, i.e. strict comparisons are always FALSE and equality comparisons are always TRUE.

	switch rhs.GetType() {
	case ArrayValue:
		lhs, err := lib_arrayval(NewNullRuntimeValue())
		if err != nil {
			return NewVoidRuntimeValue(), err
		}
		return compareRelationArray(lhs, operator, rhs)

	case BooleanValue:
		lhs, err := lib_boolval(NewNullRuntimeValue())
		if err != nil {
			return NewVoidRuntimeValue(), err
		}
		return compareRelationBoolean(NewBooleanRuntimeValue(lhs), operator, rhs)

	case FloatingValue:
		lhs, err := lib_floatval(NewNullRuntimeValue())
		if err != nil {
			return NewVoidRuntimeValue(), err
		}
		return compareRelationFloating(NewFloatingRuntimeValue(lhs), operator, rhs)

	case IntegerValue:
		lhs, err := lib_intval(NewNullRuntimeValue())
		if err != nil {
			return NewVoidRuntimeValue(), err
		}
		return compareRelationInteger(NewIntegerRuntimeValue(lhs), operator, rhs)

	case NullValue:
		switch operator {
		case "<":
			return NewBooleanRuntimeValue(false), nil
		case "<=":
			return NewBooleanRuntimeValue(true), nil
		case "<=>":
			return NewIntegerRuntimeValue(0), nil
		}
		return NewVoidRuntimeValue(), NewError("compareRelationNull: Operator \"%s\" not implemented for type NULL", operator)

		// TODO compareRelationNull - object
		// TODO compareRelationNull - resource

	case StringValue:
		lhs, err := lib_strval(NewNullRuntimeValue())
		if err != nil {
			return NewVoidRuntimeValue(), err
		}
		return compareRelationString(NewStringRuntimeValue(lhs), operator, rhs)

	default:
		return NewVoidRuntimeValue(), NewError("compareRelationNull: Type \"%s\" not implemented", rhs.GetType())
	}
}

func compareRelationString(lhs IStringRuntimeValue, operator string, rhs IRuntimeValue) (IRuntimeValue, Error) {
	// Spec: https://phplang.org/spec/10-expressions.html#grammar-relational-expression

	//         NULL  bool  int  float  string  array  object  resource
	// string   <-    ->    ->   ->     2, 4    <      3       2

	// TODO compareRelationString - object
	// TODO compareRelationString - resource

	if rhs.GetType() == FloatingValue || rhs.GetType() == IntegerValue {
		lhsStr := runtimeValToStrRuntimeVal(lhs).GetValue()
		if strings.Trim(lhsStr, " \t") == "" {
			switch operator {
			case "<", "<=":
				return NewBooleanRuntimeValue(true), nil
			case "<=>":
				return NewIntegerRuntimeValue(-1), nil
			default:
				return NewVoidRuntimeValue(), NewError("compareRelationInteger: Operator \"%s\" not implemented for type array", operator)
			}
		}
		if !common.IsIntegerLiteralWithSign(lhsStr) && !common.IsFloatingLiteralWithSign(lhsStr) {
			switch operator {
			case "<", "<=":
				return NewBooleanRuntimeValue(false), nil
			case "<=>":
				return NewIntegerRuntimeValue(1), nil
			default:
				return NewVoidRuntimeValue(), NewError("compareRelationInteger: Operator \"%s\" not implemented for type array", operator)
			}
		}
	}

	if rhs.GetType() == NullValue {
		var err Error
		rhs, err = runtimeValueToValueType(StringValue, rhs)
		if err != nil {
			return NewVoidRuntimeValue(), err
		}
	}

	switch rhs.GetType() {
	case ArrayValue:
		switch operator {
		case "<", "<=":
			return NewBooleanRuntimeValue(true), nil
		case "<=>":
			return NewIntegerRuntimeValue(-1), nil
		default:
			return NewVoidRuntimeValue(), NewError("compareRelationArray: Operator \"%s\" not implemented", operator)
		}

	case BooleanValue:
		lhs, err := lib_boolval(lhs)
		if err != nil {
			return NewVoidRuntimeValue(), err
		}
		return compareRelationBoolean(NewBooleanRuntimeValue(lhs), operator, rhs)

	case FloatingValue:
		lhs, err := lib_floatval(lhs)
		if err != nil {
			return NewVoidRuntimeValue(), err
		}
		return compareRelationFloating(NewFloatingRuntimeValue(lhs), operator, rhs)

	case IntegerValue:
		lhs, err := lib_intval(lhs)
		if err != nil {
			return NewVoidRuntimeValue(), err
		}
		return compareRelationInteger(NewIntegerRuntimeValue(lhs), operator, rhs)

	case StringValue:
		// Spec: https://phplang.org/spec/10-expressions.html#grammar-relational-expression
		//   2. If one of the operands [...] is a [...] numeric string,
		//      which can be represented as int or float without loss of precision,
		//      the operands are converted to the corresponding arithmetic type, with float taking precedence over int,
		//      and resources converting to int. The result is the numerical comparison of the two operands after conversion.
		rhsStr := runtimeValToStrRuntimeVal(rhs).GetValue()
		if common.IsFloatingLiteralWithSign(lhs.GetValue()) && (common.IsIntegerLiteralWithSign(rhsStr) || common.IsFloatingLiteralWithSign(rhsStr)) {
			lhs, err := lib_floatval(lhs)
			if err != nil {
				return NewVoidRuntimeValue(), err
			}
			return compareRelationFloating(NewFloatingRuntimeValue(lhs), operator, rhs)
		}
		if common.IsIntegerLiteralWithSign(lhs.GetValue()) && (common.IsIntegerLiteralWithSign(rhsStr) || common.IsFloatingLiteralWithSign(rhsStr)) {
			lhs, err := lib_intval(lhs)
			if err != nil {
				return NewVoidRuntimeValue(), err
			}
			return compareRelationInteger(NewIntegerRuntimeValue(lhs), operator, rhs)
		}

		// Spec: https://phplang.org/spec/10-expressions.html#grammar-relational-expression
		//   4. If both operands are non-numeric strings, the result is the lexical comparison of the two operands.
		//      Specifically, the strings are compared byte-by-byte starting with their first byte.
		//      If the two bytes compare equal and there are no more bytes in either string, the strings are equal and the comparison ends;
		//      otherwise, if this is the final byte in one string, the shorter string compares less-than the longer string and the comparison ends.
		//      If the two bytes compare unequal, the string having the lower-valued byte compares less-than the other string, and the comparison ends.
		//      If there are more bytes in the strings, the process is repeated for the next pair of bytes.
		var result int64 = 0
		for index, lhsByte := range []byte(lhs.GetValue()) {
			if index >= len(rhsStr) {
				result = 1
				break
			}
			rhsByte := rhsStr[index]
			if lhsByte > rhsByte {
				result = 1
				break
			}
			if lhsByte < rhsByte {
				result = -1
				break
			}
		}
		if result == 0 && len(lhs.GetValue()) < len(rhsStr) {
			result = -1
		}
		switch operator {
		case "<":
			return NewBooleanRuntimeValue(result == -1), nil
		case "<=":
			return NewBooleanRuntimeValue(result < 1), nil
		case "<=>":
			return NewIntegerRuntimeValue(result), nil
		default:
			return NewVoidRuntimeValue(), NewError("compareRelationString: Operator \"%s\" not implemented", operator)
		}

	default:
		return NewVoidRuntimeValue(), NewError("compareRelationString: Type \"%s\" not implemented", rhs.GetType())
	}
}

// TODO compareRelationObject
// Spec: https://phplang.org/spec/10-expressions.html#grammar-relational-expression
//         NULL  bool  int  float  string  array  object  resource
// object   >     ->    3    3      3       3      6       3

// TODO compareRelationResource
// Spec: https://phplang.org/spec/10-expressions.html#grammar-relational-expression
//           NULL  bool  int  float  string  array  object  resource
// resource   >     ->    ->   ->     2       <      3       2

// ------------------- MARK: comparison -------------------

func compare(lhs IRuntimeValue, operator string, rhs IRuntimeValue) (IBooleanRuntimeValue, Error) {
	// Spec: https://phplang.org/spec/10-expressions.html#grammar-equality-expression

	// Spec: https://phplang.org/spec/10-expressions.html#grammar-equality-expression
	// Operator == represents value equality, operators != and <> are equivalent and represent value inequality.
	// For operators ==, !=, and <>, the operands of different types are converted and compared according to the same rules as in relational operators.
	// Two objects of different types are always not equal.
	if operator == "<>" {
		operator = "!="
	}
	if operator == "==" || operator == "!=" {
		resultRuntimeValue, err := compareRelation(lhs, "<=>", rhs)
		if err != nil {
			return NewBooleanRuntimeValue(false), err
		}
		result := runtimeValToIntRuntimeVal(resultRuntimeValue).GetValue() == 0

		if operator == "!=" {
			return NewBooleanRuntimeValue(!result), nil
		} else {
			return NewBooleanRuntimeValue(result), nil
		}
	}

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
			case ArrayValue:
				lhsArray := runtimeValToArrayRuntimeVal(lhs)
				rhsArray := runtimeValToArrayRuntimeVal(rhs)
				if len(lhsArray.GetKeys()) != len(rhsArray.GetKeys()) {
					result = false
				} else {
					for key, lhsValue := range lhsArray.GetElements() {
						rhsValue, found := rhsArray.GetElement(key)
						if !found {
							result = false
							break
						}
						equal, err := compare(lhsValue, "===", rhsValue)
						if err != nil {
							return NewBooleanRuntimeValue(false), err
						}
						if !runtimeValToBoolRuntimeVal(equal).GetValue() {
							result = false
							break
						}
					}
					result = true
				}
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
