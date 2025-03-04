package interpreter

import (
	"GoPHP/cmd/goPHP/ast"
	"GoPHP/cmd/goPHP/common"
	"GoPHP/cmd/goPHP/config"
	"GoPHP/cmd/goPHP/phpError"
	"GoPHP/cmd/goPHP/runtime/values"
	"math"
	"os"
	"regexp"
	"runtime"
	"slices"
	"strconv"
	"strings"
)

func printDev(str string) {
	if config.IsDevMode {
		println(str)
	}
}

func (interpreter *Interpreter) print(str string) {
	if len(interpreter.outputBuffers) > 0 {
		interpreter.outputBuffers[len(interpreter.outputBuffers)-1].Content += str
	} else {
		interpreter.result += str
	}
}

func (interpreter *Interpreter) flushOutputBuffers() {
	if len(interpreter.outputBuffers) == 0 {
		return
	}

	for len(interpreter.outputBuffers) > 0 {
		nativeFn_ob_end_flush([]values.RuntimeValue{}, interpreter)
	}
}

var PHP_EOL string = getPhpEol()

func getPhpEol() string {
	if getPhpOs() == "Windows" {
		return "\r\n"
	} else {
		return "\n"
	}
}

var DIR_SEP = getPhpDirectorySeparator()

func getPhpDirectorySeparator() string {
	if getPhpOs() == "Windows" {
		return `\`
	} else {
		return "/"
	}
}

func (interpreter *Interpreter) println(str string) {
	interpreter.print(str + PHP_EOL)
}

func (interpreter *Interpreter) processCondition(expr ast.IExpression, env *Environment) (values.RuntimeValue, bool, phpError.Error) {
	runtimeValue, err := interpreter.processStmt(expr, env)
	if err != nil {
		return runtimeValue, false, err
	}

	boolean, err := lib_boolval(runtimeValue)
	return runtimeValue, boolean, err
}

func (interpreter *Interpreter) lookupVariable(expr ast.IExpression, env *Environment) (values.RuntimeValue, phpError.Error) {
	variableName, err := interpreter.varExprToVarName(expr, env)
	if err != nil {
		return values.NewVoid(), err
	}

	runtimeValue, err := env.lookupVariable(variableName)
	if !interpreter.suppressWarning && err != nil {
		interpreter.printError(err)
	}
	return runtimeValue, nil
}

// Convert a variable expression into the interpreted variable name
func (interpreter *Interpreter) varExprToVarName(expr ast.IExpression, env *Environment) (string, phpError.Error) {
	switch expr.GetKind() {
	case ast.SimpleVariableExpr:
		variableNameExpr := expr.(*ast.SimpleVariableExpression).VariableName

		if variableNameExpr.GetKind() == ast.VariableNameExpr {
			return variableNameExpr.(*ast.VariableNameExpression).VariableName, nil
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

		variableName, err := interpreter.processStmt(variableNameExpr, env)
		if err != nil {
			return "", err
		}
		valueStr, err := lib_strval(variableName)
		if err != nil {
			return "", err
		}
		return "$" + valueStr, nil
	case ast.SubscriptExpr:
		return interpreter.varExprToVarName(expr.(*ast.SubscriptExpression).Variable, env)
	default:
		return "", phpError.NewError("varExprToVarName: Unsupported expression: %s", ast.ToString(expr))
	}
}

func (interpreter *Interpreter) ErrorToString(err phpError.Error) string {
	if (err.GetErrorType() == phpError.WarningPhpError && interpreter.ini.GetInt("error_reporting")&phpError.E_WARNING == 0) ||
		(err.GetErrorType() == phpError.ErrorPhpError && interpreter.ini.GetInt("error_reporting")&phpError.E_ERROR == 0) ||
		(err.GetErrorType() == phpError.ParsePhpError && interpreter.ini.GetInt("error_reporting")&phpError.E_PARSE == 0) {
		return ""
	}
	return err.GetMessage()
}

func (interpreter *Interpreter) printError(err phpError.Error) {
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
func (interpreter *Interpreter) scanForFunctionDefinition(statements []ast.IStatement, env *Environment) phpError.Error {
	for _, stmt := range statements {
		if stmt.GetKind() == ast.CompoundStmt {
			interpreter.scanForFunctionDefinition(stmt.(*ast.CompoundStatement).Statements, env)
			continue
		}

		if stmt.GetKind() != ast.FunctionDefinitionStmt {
			continue
		}

		_, err := interpreter.processStmt(stmt, env)
		if err != nil {
			return err
		}
	}
	return nil
}

var literalExprTypeRuntimeValue = map[ast.NodeType]string{
	ast.ArrayLiteralExpr:   "array",
	ast.IntegerLiteralExpr: "int",
	ast.StringLiteralExpr:  "string",
}

func literalExprTypeToRuntimeValue(expr ast.IExpression) (string, phpError.Error) {
	typeStr, found := literalExprTypeRuntimeValue[expr.GetKind()]
	if !found {
		return "", phpError.NewError("literalExprTypeToRuntimeValue: No mapping for type %s", expr.GetKind())
	}
	return typeStr, nil
}

func checkParameterTypes(runtimeValue values.RuntimeValue, expectedTypes []string) phpError.Error {
	typeStr := values.ToPhpType(runtimeValue)
	if typeStr == "" {
		return phpError.NewError("checkParameterTypes: No mapping for type %s", runtimeValue.GetType())
	}

	for _, expectedType := range expectedTypes {
		if expectedType == "mixed" {
			return nil
		}

		if typeStr == expectedType {
			return nil
		}
	}
	return phpError.NewError("Types do not match")
}

func (interpreter *Interpreter) includeFile(filepathExpr ast.IExpression, env *Environment, include bool, once bool) (values.RuntimeValue, phpError.Error) {
	runtimeValue, err := interpreter.processStmt(filepathExpr, env)
	if err != nil {
		return runtimeValue, err
	}
	if runtimeValue.GetType() == values.NullValue {
		return runtimeValue, phpError.NewError("Uncaught ValueError: Path cannot be empty in %s", filepathExpr.GetPosition().ToPosString())
	}

	filename, err := lib_strval(runtimeValue)
	if err != nil {
		return runtimeValue, err
	}

	// Spec: https://phplang.org/spec/10-expressions.html#the-require-operator
	// Once an include file has been included, a subsequent use of require_once on that include file
	// results in a return value of TRUE but nothing else happens.
	if once && slices.Contains(interpreter.includedFiles, filename) && runtime.GOOS != "windows" {
		return values.NewBool(true), nil
	}
	if once && slices.Contains(interpreter.includedFiles, strings.ToLower(filename)) && runtime.GOOS == "windows" {
		return values.NewBool(true), nil
	}

	absFilename := filename
	if !common.IsAbsPath(filename) {
		absFilename = common.GetAbsPathForWorkingDir(common.ExtractPath(filepathExpr.GetPosition().Filename), filename)
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
	getError := func() (values.RuntimeValue, phpError.Error) {
		if include {
			return values.NewVoid(), phpError.NewWarning(
				"%s(): Failed opening '%s' for inclusion (include_path='%s') in %s",
				functionName, filename, common.ExtractPath(filepathExpr.GetPosition().Filename), filepathExpr.GetPosition().ToPosString(),
			)
		} else {
			return values.NewVoid(), phpError.NewError(
				"Uncaught Error: Failed opening required '%s' (include_path='%s') in %s",
				filename, common.ExtractPath(filepathExpr.GetPosition().Filename), filepathExpr.GetPosition().ToPosString(),
			)
		}
	}

	if !common.PathExists(absFilename) {
		interpreter.printError(phpError.NewWarning(
			"%s(%s): Failed to open stream: No such file or directory in %s",
			functionName, filename, filepathExpr.GetPosition().ToPosString(),
		))
		return getError()
	}

	content, fileErr := os.ReadFile(absFilename)
	if fileErr != nil {
		return getError()
	}
	program, parserErr := interpreter.parser.ProduceAST(string(content), filename)

	if runtime.GOOS != "windows" {
		interpreter.includedFiles = append(interpreter.includedFiles, absFilename)
	} else {
		interpreter.includedFiles = append(interpreter.includedFiles, strings.ToLower(absFilename))
	}
	if parserErr != nil {
		return runtimeValue, parserErr
	}
	return interpreter.processProgram(program, env)
}

// ------------------- MARK: Caching -------------------

func (interpreter *Interpreter) isCached(stmt ast.IStatement) bool {
	_, found := interpreter.cache[stmt.GetId()]
	return found
}

func (interpreter *Interpreter) writeCache(stmt ast.IStatement, value values.RuntimeValue) values.RuntimeValue {
	interpreter.cache[stmt.GetId()] = value
	return value
}

// ------------------- MARK: RuntimeValue -------------------

func (interpreter *Interpreter) exprToRuntimeValue(expr ast.IExpression, env *Environment) (values.RuntimeValue, phpError.Error) {
	switch expr.GetKind() {
	case ast.ArrayLiteralExpr:
		Array := values.NewArray()
		for _, key := range expr.(*ast.ArrayLiteralExpression).Keys {
			var keyValue values.RuntimeValue
			var err phpError.Error
			if key.GetKind() != ast.ArrayNextKeyExpr {
				keyValue, err = interpreter.processStmt(key, env)
				if err != nil {
					return values.NewVoid(), err
				}
			}
			elementValue, err := interpreter.processStmt(expr.(*ast.ArrayLiteralExpression).Elements[key], env)
			if err != nil {
				return values.NewVoid(), err
			}
			if err = Array.SetElement(keyValue, elementValue); err != nil {
				return values.NewVoid(), err
			}
		}
		return Array, nil
	case ast.IntegerLiteralExpr:
		return values.NewInt(expr.(*ast.IntegerLiteralExpression).Value), nil
	case ast.FloatingLiteralExpr:
		return values.NewFloat(expr.(*ast.FloatingLiteralExpression).Value), nil
	case ast.StringLiteralExpr:
		str := expr.(*ast.StringLiteralExpression).Value
		if expr.(*ast.StringLiteralExpression).StringType == ast.DoubleQuotedString {
			// variable substitution
			r, _ := regexp.Compile(`({\$[A-Za-z_][A-Za-z0-9_]*['A-Za-z0-9\[\]]*[^}]*})|(\$[A-Za-z_][A-Za-z0-9_]*['A-Za-z0-9\[\]]*)`)
			matches := r.FindAllString(str, -1)
			for _, match := range matches {
				varExpr := match
				if match[0] == '{' {
					// Remove curly braces
					varExpr = match[1 : len(match)-1]
				}
				exprStr := "<?= " + varExpr + ";"
				result, err := NewInterpreter(interpreter.ini, interpreter.request, "").process(exprStr, env, true)
				if err != nil {
					return values.NewVoid(), err
				}
				str = strings.Replace(str, match, result, 1)
			}

			// unicode escape sequence
			r, _ = regexp.Compile(`\\u\{[0-9a-fA-F]+\}`)
			matches = r.FindAllString(str, -1)
			for _, match := range matches {
				unicodeChar, err := strconv.ParseInt(match[3:len(match)-1], 16, 32)
				if err != nil {
					return values.NewVoid(), phpError.NewError(err.Error())
				}
				str = strings.Replace(str, match, string(rune(unicodeChar)), 1)
			}
		}
		return values.NewStr(str), nil
	default:
		return values.NewVoid(), phpError.NewError("exprToRuntimeValue: Unsupported expression: %s", expr)
	}
}

func runtimeValueToValueType(valueType values.ValueType, runtimeValue values.RuntimeValue) (values.RuntimeValue, phpError.Error) {
	switch valueType {
	case values.BoolValue:
		boolean, err := lib_boolval(runtimeValue)
		return values.NewBool(boolean), err
	case values.FloatValue:
		floating, err := lib_floatval(runtimeValue)
		return values.NewFloat(floating), err
	case values.IntValue:
		integer, err := lib_intval(runtimeValue)
		return values.NewInt(integer), err
	case values.StrValue:
		str, err := lib_strval(runtimeValue)
		return values.NewStr(str), err
	default:
		return values.NewVoid(), phpError.NewError("runtimeValueToValueType: Unsupported runtime value: %s", valueType)
	}
}

func deepCopy(value values.RuntimeValue) values.RuntimeValue {
	if value.GetType() != values.ArrayValue {
		return value
	}

	copy := values.NewArray()
	array := value.(*values.Array)
	for _, key := range array.Keys {
		value, _ := array.GetElement(key)
		copy.SetElement(key, deepCopy(value))
	}
	return copy
}

// ------------------- MARK: inc-dec-calculation -------------------

func calculateIncDec(operator string, operand values.RuntimeValue) (values.RuntimeValue, phpError.Error) {
	switch operand.GetType() {
	case values.BoolValue:
		// Spec: https://phplang.org/spec/10-expressions.html#prefix-increment-and-decrement-operators
		// For a prefix ++ or -- operator used with a Boolean-valued operand, there is no side effect, and the result is the operand’s value.
		return operand, nil
	case values.FloatValue:
		return calculateIncDecFloating(operator, operand.(*values.Float))
	case values.IntValue:
		return calculateIncDecInteger(operator, operand.(*values.Int))
	case values.NullValue:
		return calculateIncDecNull(operator)
	case values.StrValue:
		return calculateIncDecString(operator, operand.(*values.Str))
	default:
		return values.NewVoid(), phpError.NewError("calculateIncDec: Type \"%s\" not implemented", operand.GetType())
	}

	// TODO calculateIncDec - object
	// Spec: https://phplang.org/spec/10-expressions.html#prefix-increment-and-decrement-operators
	// If the operand has an object type supporting the operation, then the object semantics defines the result. Otherwise, the operation has no effect and the result is the operand.
}

func calculateIncDecInteger(operator string, operand *values.Int) (values.RuntimeValue, phpError.Error) {
	switch operator {
	case "++":
		// Spec: https://phplang.org/spec/10-expressions.html#prefix-increment-and-decrement-operators
		//For a prefix "++" operator used with an arithmetic operand, the side effect of the operator is to increment the value of the operand by 1.
		// The result is the value of the operand after it has been incremented.
		// If an int operand’s value is the largest representable for that type, the operand is incremented as if it were float.

		// Spec: https://phplang.org/spec/10-expressions.html#prefix-increment-and-decrement-operators
		// For a prefix ++ or -- operator used with an operand having the value INF, -INF, or NAN, there is no side effect, and the result is the operand’s value.
		return calculateInteger(operand, "+", values.NewInt(1))

	case "--":
		// Spec: https://phplang.org/spec/10-expressions.html#prefix-increment-and-decrement-operators
		// For a prefix "--" operator used with an arithmetic operand, the side effect of the operator is to decrement the value of the operand by 1.
		// The result is the value of the operand after it has been decremented.
		// If an int operand’s value is the smallest representable for that type, the operand is decremented as if it were float.

		// Spec: https://phplang.org/spec/10-expressions.html#prefix-increment-and-decrement-operators
		// For a prefix ++ or -- operator used with an operand having the value INF, -INF, or NAN, there is no side effect, and the result is the operand’s value.
		return calculateInteger(operand, "-", values.NewInt(1))

	default:
		return values.NewInt(0), phpError.NewError("calculateIncDecInteger: Operator \"%s\" not implemented", operator)
	}
}

func calculateIncDecFloating(operator string, operand *values.Float) (values.RuntimeValue, phpError.Error) {
	switch operator {
	case "++":
		// Spec: https://phplang.org/spec/10-expressions.html#prefix-increment-and-decrement-operators
		//For a prefix "++" operator used with an arithmetic operand, the side effect of the operator is to increment the value of the operand by 1.
		// The result is the value of the operand after it has been incremented.
		// If an int operand’s value is the largest representable for that type, the operand is incremented as if it were float.

		// Spec: https://phplang.org/spec/10-expressions.html#prefix-increment-and-decrement-operators
		// For a prefix ++ or -- operator used with an operand having the value INF, -INF, or NAN, there is no side effect, and the result is the operand’s value.
		return calculateFloating(operand, "+", values.NewFloat(1))

	case "--":
		// Spec: https://phplang.org/spec/10-expressions.html#prefix-increment-and-decrement-operators
		// For a prefix "--" operator used with an arithmetic operand, the side effect of the operator is to decrement the value of the operand by 1.
		// The result is the value of the operand after it has been decremented.
		// If an int operand’s value is the smallest representable for that type, the operand is decremented as if it were float.

		// Spec: https://phplang.org/spec/10-expressions.html#prefix-increment-and-decrement-operators
		// For a prefix ++ or -- operator used with an operand having the value INF, -INF, or NAN, there is no side effect, and the result is the operand’s value.
		return calculateFloating(operand, "-", values.NewFloat(1))

	default:
		return values.NewInt(0), phpError.NewError("calculateIncDecFloating: Operator \"%s\" not implemented", operator)
	}
}

func calculateIncDecNull(operator string) (values.RuntimeValue, phpError.Error) {
	switch operator {
	case "++":
		// Spec: https://phplang.org/spec/10-expressions.html#prefix-increment-and-decrement-operators
		// For a prefix ++ operator used with a NULL-valued operand, the side effect is that the operand’s type is changed to int,
		// the operand’s value is set to zero, and that value is incremented by 1.
		// The result is the value of the operand after it has been incremented.
		return values.NewInt(1), nil

	case "--":
		// Spec: https://phplang.org/spec/10-expressions.html#prefix-increment-and-decrement-operators
		// For a prefix – operator used with a NULL-valued operand, there is no side effect, and the result is the operand’s value.
		return values.NewNull(), nil

	default:
		return values.NewInt(0), phpError.NewError("calculateIncDecNull: Operator \"%s\" not implemented", operator)
	}
}

func calculateIncDecString(operator string, operand *values.Str) (values.RuntimeValue, phpError.Error) {
	switch operator {
	case "++":
		// Spec: https://phplang.org/spec/10-expressions.html#prefix-increment-and-decrement-operators
		// For a prefix "++" operator used with an operand whose value is an empty string,
		// the side effect is that the operand’s value is changed to the string “1”. The type of the operand is unchanged.
		// The result is the new value of the operand.
		if operand.Value == "" {
			return values.NewStr("1"), nil
		}
		return values.NewVoid(), phpError.NewError("TODO calculateIncDecString")

	case "--":
		// Spec: https://phplang.org/spec/10-expressions.html#prefix-increment-and-decrement-operators
		// For a prefix "--" operator used with an operand whose value is an empty string,
		// the side effect is that the operand’s type is changed to int, the operand’s value is set to zero,
		// and that value is decremented by 1. The result is the value of the operand after it has been incremented.
		if operand.Value == "" {
			return values.NewInt(-1), nil
		}
		return values.NewVoid(), phpError.NewError("TODO calculateIncDecString")

	default:
		return values.NewInt(0), phpError.NewError("calculateIncDecNull: Operator \"%s\" not implemented", operator)
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

func calculateUnary(operator string, operand values.RuntimeValue) (values.RuntimeValue, phpError.Error) {
	switch operand.GetType() {
	case values.BoolValue:
		return calculateUnaryBoolean(operator, operand.(*values.Bool))
	case values.IntValue:
		return calculateUnaryInteger(operator, operand.(*values.Int))
	case values.FloatValue:
		return calculateUnaryFloating(operator, operand.(*values.Float))
	case values.NullValue:
		// Spec: https://phplang.org/spec/10-expressions.html#unary-arithmetic-operators
		// For a unary + or unary - operator used with a NULL-valued operand, the value of the result is zero and the type is int.
		return values.NewInt(0), nil
	default:
		return values.NewVoid(), phpError.NewError("calculateUnary: Type \"%s\" not implemented", operand.GetType())
	}

	// TODO calculateUnary - string
	// Spec: https://phplang.org/spec/10-expressions.html#unary-arithmetic-operators
	// For a unary + or - operator used with a numeric string or a leading-numeric string, the string is first converted to an int or float, as appropriate, after which it is handled as an arithmetic operand. The trailing non-numeric characters in leading-numeric strings are ignored. With a non-numeric string, the result has type int and value 0. If the string was leading-numeric or non-numeric, a non-fatal error MUST be produced.
	// For a unary ~ operator used with a string, the result is the string with each byte being bitwise complement of the corresponding byte of the source string.

	// TODO calculateUnary - object
	// If the operand has an object type supporting the operation, then the object semantics defines the result. Otherwise, for ~ the fatal error is issued and for + and - the object is converted to int.
}

func calculateUnaryBoolean(operator string, operand *values.Bool) (*values.Int, phpError.Error) {
	switch operator {
	case "+":
		// Spec: https://phplang.org/spec/10-expressions.html#unary-arithmetic-operators
		// For a unary "+" operator used with a TRUE-valued operand, the value of the result is 1 and the type is int.
		// When used with a FALSE-valued operand, the value of the result is zero and the type is int.
		if operand.Value {
			return values.NewInt(1), nil
		}
		return values.NewInt(0), nil

	case "-":
		// Spec: https://phplang.org/spec/10-expressions.html#unary-arithmetic-operators
		// For a unary "-" operator used with a TRUE-valued operand, the value of the result is -1 and the type is int.
		// When used with a FALSE-valued operand, the value of the result is zero and the type is int.
		if operand.Value {
			return values.NewInt(-1), nil
		}
		return values.NewInt(0), nil

	default:
		return values.NewInt(0), phpError.NewError("calculateUnaryBoolean: Operator \"%s\" not implemented", operator)
	}
}

func calculateUnaryFloating(operator string, operand *values.Float) (values.RuntimeValue, phpError.Error) {
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
		return values.NewFloat(-operand.Value), nil

	case "~":
		// Spec: https://phplang.org/spec/10-expressions.html#unary-arithmetic-operators
		// For a unary ~ operator used with a float operand, the value of the operand is first converted to int before the bitwise complement is computed.
		intRuntimeValue, err := runtimeValueToValueType(values.IntValue, operand)
		if err != nil {
			return values.NewFloat(0), err
		}
		return calculateUnaryInteger(operator, intRuntimeValue.(*values.Int))

	default:
		return values.NewFloat(0), phpError.NewError("calculateUnaryFloating: Operator \"%s\" not implemented", operator)
	}
}

func calculateUnaryInteger(operator string, operand *values.Int) (*values.Int, phpError.Error) {
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
		return values.NewInt(-operand.Value), nil

	case "~":
		// Spec: https://phplang.org/spec/10-expressions.html#unary-arithmetic-operators
		// For a unary ~ operator used with an int operand, the type of the result is int.
		// The value of the result is the bitwise complement of the value of the operand
		// (that is, each bit in the result is set if and only if the corresponding bit in the operand is clear).
		return values.NewInt(^operand.Value), nil
	default:
		return values.NewInt(0), phpError.NewError("calculateUnaryInteger: Operator \"%s\" not implemented", operator)
	}
}

// ------------------- MARK: binary-op-calculation -------------------

func calculate(operand1 values.RuntimeValue, operator string, operand2 values.RuntimeValue) (values.RuntimeValue, phpError.Error) {
	resultType := values.VoidValue
	if slices.Contains([]string{"."}, operator) {
		resultType = values.StrValue
	} else if slices.Contains([]string{"&", "|", "^", "<<", ">>"}, operator) {
		resultType = values.IntValue
	} else {
		resultType = values.IntValue
		if operand1.GetType() == values.FloatValue || operand2.GetType() == values.FloatValue {
			resultType = values.FloatValue
		}
	}

	var err phpError.Error
	operand1, err = runtimeValueToValueType(resultType, operand1)
	if err != nil {
		return values.NewVoid(), err
	}
	operand2, err = runtimeValueToValueType(resultType, operand2)
	if err != nil {
		return values.NewVoid(), err
	}
	// TODO testing how PHP behavious: var_dump(1.0 + 2); var_dump(1 + 2.0); var_dump("1" + 2);
	// var_dump("1" + "2"); => int
	// var_dump("1" . 2); => str
	// type order "string" - "int" - "float"

	// Testen
	//   true + 2
	//   true && 3

	switch resultType {
	case values.IntValue:
		return calculateInteger(operand1.(*values.Int), operator, operand2.(*values.Int))
	case values.FloatValue:
		return calculateFloating(operand1.(*values.Float), operator, operand2.(*values.Float))
	case values.StrValue:
		return calculateString(operand1.(*values.Str), operator, operand2.(*values.Str))
	default:
		return values.NewVoid(), phpError.NewError("calculate: Type \"%s\" not implemented", resultType)
	}
}

func calculateFloating(operand1 *values.Float, operator string, operand2 *values.Float) (*values.Float, phpError.Error) {
	switch operator {
	case "+":
		return values.NewFloat(operand1.Value + operand2.Value), nil
	case "-":
		return values.NewFloat(operand1.Value - operand2.Value), nil
	case "*":
		return values.NewFloat(operand1.Value * operand2.Value), nil
	case "/":
		return values.NewFloat(operand1.Value / operand2.Value), nil
	case "**":
		return values.NewFloat(math.Pow(operand1.Value, operand2.Value)), nil
	default:
		return values.NewFloat(0), phpError.NewError("calculateFloating: Operator \"%s\" not implemented", operator)
	}
}

func calculateInteger(operand1 *values.Int, operator string, operand2 *values.Int) (*values.Int, phpError.Error) {
	switch operator {
	case "<<":
		return values.NewInt(operand1.Value << operand2.Value), nil
	case ">>":
		return values.NewInt(operand1.Value >> operand2.Value), nil
	case "^":
		return values.NewInt(operand1.Value ^ operand2.Value), nil
	case "|":
		return values.NewInt(operand1.Value | operand2.Value), nil
	case "&":
		return values.NewInt(operand1.Value & operand2.Value), nil
	case "+":
		return values.NewInt(operand1.Value + operand2.Value), nil
	case "-":
		return values.NewInt(operand1.Value - operand2.Value), nil
	case "*":
		return values.NewInt(operand1.Value * operand2.Value), nil
	case "/":
		return values.NewInt(operand1.Value / operand2.Value), nil
	case "%":
		return values.NewInt(operand1.Value % operand2.Value), nil
	case "**":
		return values.NewInt(int64(math.Pow(float64(operand1.Value), float64(operand2.Value)))), nil
	default:
		return values.NewInt(0), phpError.NewError("calculateInteger: Operator \"%s\" not implemented", operator)
	}
}

func calculateString(operand1 *values.Str, operator string, operand2 *values.Str) (*values.Str, phpError.Error) {
	switch operator {
	case ".":
		return values.NewStr(operand1.Value + operand2.Value), nil
	default:
		return values.NewStr(""), phpError.NewError("calculateString: Operator \"%s\" not implemented", operator)
	}
}

// ------------------- MARK: compareRelation -------------------

func compareRelation(lhs values.RuntimeValue, operator string, rhs values.RuntimeValue) (values.RuntimeValue, phpError.Error) {
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
	case values.ArrayValue:
		return compareRelationArray(lhs.(*values.Array), operator, rhs)
	case values.BoolValue:
		return compareRelationBoolean(lhs.(*values.Bool), operator, rhs)
	case values.FloatValue:
		return compareRelationFloating(lhs.(*values.Float), operator, rhs)
	case values.IntValue:
		return compareRelationInteger(lhs.(*values.Int), operator, rhs)
	case values.StrValue:
		return compareRelationString(lhs.(*values.Str), operator, rhs)
	case values.NullValue:
		return compareRelationNull(operator, rhs)
	default:
		return values.NewVoid(), phpError.NewError("compareRelation: Type \"%s\" not implemented", lhs.GetType())
	}

}

func compareRelationArray(lhs *values.Array, operator string, rhs values.RuntimeValue) (values.RuntimeValue, phpError.Error) {
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

	if rhs.GetType() == values.NullValue {
		var err phpError.Error
		rhs, err = lib_arrayval(rhs)
		if err != nil {
			return values.NewVoid(), err
		}
	}

	switch rhs.GetType() {
	case values.ArrayValue:
		rhsArray := rhs.(*values.Array)
		var result int64 = 0
		if len(lhs.Keys) != len(rhsArray.Keys) {
			if len(lhs.Keys) < len(rhsArray.Keys) {
				result = -1
			} else {
				result = 1
			}
		} else {
			for _, key := range lhs.Keys {
				lhsValue, _ := lhs.GetElement(key)
				rhsValue, found := rhsArray.GetElement(key)
				if found {
					equal, err := compare(lhsValue, "===", rhsValue)
					if err != nil {
						return values.NewVoid(), err
					}
					if equal.Value {
						continue
					}
					lessThan, err := compareRelation(lhsValue, operator, rhsValue)
					if err != nil {
						return values.NewVoid(), err
					}
					if lessThan.GetType() == values.BoolValue {
						if lessThan.(*values.Bool).Value {
							result = -1
						} else {
							result = 1
						}
					}
					if lessThan.GetType() == values.IntValue {
						result = lessThan.(*values.Int).Value
					}
				}
			}
		}

		switch operator {
		case "<":
			return values.NewBool(result == -1), nil
		case "<=":
			return values.NewBool(result < 1), nil
		case "<=>":
			return values.NewInt(result), nil
		default:
			return values.NewVoid(), phpError.NewError("compareRelationArray: Operator \"%s\" not implemented", operator)
		}

	case values.BoolValue:
		lhsBoolean, err := lib_boolval(lhs)
		if err != nil {
			return values.NewVoid(), err
		}
		return compareRelationBoolean(values.NewBool(lhsBoolean), operator, rhs)

	case values.FloatValue, values.IntValue, values.StrValue:
		switch operator {
		case "<", "<=":
			return values.NewBool(false), nil
		case "<=>":
			return values.NewInt(1), nil
		default:
			return values.NewVoid(), phpError.NewError("compareRelationArray: Operator \"%s\" not implemented", operator)
		}

	default:
		return values.NewVoid(), phpError.NewError("compareRelationArray: Type \"%s\" not implemented", rhs.GetType())
	}
}

func compareRelationBoolean(lhs *values.Bool, operator string, rhs values.RuntimeValue) (values.RuntimeValue, phpError.Error) {
	// Spec: https://phplang.org/spec/10-expressions.html#grammar-relational-expression

	//       NULL  bool  int  float  string  array  object  resource
	// bool   <-    1     <-   <-     <-      <-     <-      <-

	//   1. If either operand has type bool, the other operand is converted to that type.
	//      The result is the logical comparison of the two operands after conversion, where FALSE is defined to be less than TRUE.

	rhsBoolean, err := lib_boolval(rhs)
	if err != nil {
		return values.NewVoid(), err
	}
	// TODO compareRelationBoolean - object - implement in lib_boolval
	// TODO compareRelationBoolean - resource - implement in lib_boolval

	lhsInt, err := lib_intval(lhs)
	if err != nil {
		return values.NewVoid(), err
	}
	rhsInt, err := lib_intval(values.NewBool(rhsBoolean))
	if err != nil {
		return values.NewVoid(), err
	}

	switch operator {
	case "<":
		return values.NewBool(lhsInt < rhsInt), nil

	case "<=":
		return values.NewBool(lhsInt <= rhsInt), nil

	case "<=>":
		if lhsInt > rhsInt {
			return values.NewInt(1), nil
		}
		if lhsInt == rhsInt {
			return values.NewInt(0), nil
		}
		return values.NewInt(-1), nil

	default:
		return values.NewVoid(), phpError.NewError("compareRelationBoolean: Operator \"%s\" not implemented", operator)
	}
}

func compareRelationFloating(lhs *values.Float, operator string, rhs values.RuntimeValue) (values.RuntimeValue, phpError.Error) {
	// Spec: https://phplang.org/spec/10-expressions.html#grammar-relational-expression

	//        NULL  bool  int  float  string  array  object  resource
	// float   <-    ->    2    2      <-      <      3       <-

	// TODO compareRelationFloating - object
	// TODO compareRelationFloating - resource

	if rhs.GetType() == values.StrValue {
		rhsStr := rhs.(*values.Str).Value
		if strings.Trim(rhsStr, " \t") == "" {
			switch operator {
			case "<", "<=":
				return values.NewBool(false), nil
			case "<=>":
				return values.NewInt(1), nil
			default:
				return values.NewVoid(), phpError.NewError("compareRelationFloating: Operator \"%s\" not implemented for type string", operator)
			}
		}
		if !common.IsIntegerLiteralWithSign(rhsStr) && !common.IsFloatingLiteralWithSign(rhsStr) {
			switch operator {
			case "<", "<=":
				return values.NewBool(true), nil
			case "<=>":
				return values.NewInt(-1), nil
			default:
				return values.NewVoid(), phpError.NewError("compareRelationFloating: Operator \"%s\" not implemented for type string", operator)
			}
		}
	}

	if rhs.GetType() == values.NullValue || rhs.GetType() == values.IntValue || rhs.GetType() == values.StrValue {
		var err phpError.Error
		rhs, err = runtimeValueToValueType(values.FloatValue, rhs)
		if err != nil {
			return values.NewVoid(), err
		}
	}

	switch rhs.GetType() {
	case values.ArrayValue:
		switch operator {
		case "<", "<=":
			return values.NewBool(true), nil
		case "<=>":
			return values.NewInt(-1), nil
		default:
			return values.NewVoid(), phpError.NewError("compareRelationFloating: Operator \"%s\" not implemented for type array", operator)
		}

	case values.BoolValue:
		lhsBoolean, err := lib_boolval(lhs)
		if err != nil {
			return values.NewVoid(), err
		}
		return compareRelationBoolean(values.NewBool(lhsBoolean), operator, rhs)

	case values.FloatValue:
		rhsFloat := rhs.(*values.Float).Value
		switch operator {
		case "<":
			return values.NewBool(lhs.Value < rhsFloat), nil
		case "<=":
			return values.NewBool(lhs.Value <= rhsFloat), nil
		case "<=>":
			if lhs.Value > rhsFloat {
				return values.NewInt(1), nil
			}
			if lhs.Value == rhsFloat {
				return values.NewInt(0), nil
			}
			return values.NewInt(-1), nil
		default:
			return values.NewVoid(), phpError.NewError("compareRelationFloating: Operator \"%s\" not implemented", operator)
		}

	default:
		return values.NewVoid(), phpError.NewError("compareRelationFloating: Type \"%s\" not implemented", rhs.GetType())
	}
}

func compareRelationInteger(lhs *values.Int, operator string, rhs values.RuntimeValue) (values.RuntimeValue, phpError.Error) {
	// Spec: https://phplang.org/spec/10-expressions.html#grammar-relational-expression

	//      NULL  bool  int  float  string  array  object  resource
	// int   <-    ->    2    2      <-      <      3       <-

	// TODO compareRelationInteger - object
	// TODO compareRelationInteger - resource

	if rhs.GetType() == values.StrValue {
		rhsStr := rhs.(*values.Str).Value
		if strings.Trim(rhsStr, " \t") == "" {
			switch operator {
			case "<", "<=":
				return values.NewBool(false), nil
			case "<=>":
				return values.NewInt(1), nil
			default:
				return values.NewVoid(), phpError.NewError("compareRelationInteger: Operator \"%s\" not implemented for type array", operator)
			}
		}
		if !common.IsIntegerLiteralWithSign(rhsStr) && !common.IsFloatingLiteralWithSign(rhsStr) {
			switch operator {
			case "<", "<=":
				return values.NewBool(true), nil
			case "<=>":
				return values.NewInt(-1), nil
			default:
				return values.NewVoid(), phpError.NewError("compareRelationInteger: Operator \"%s\" not implemented for type array", operator)
			}
		}
	}

	if rhs.GetType() == values.NullValue || rhs.GetType() == values.StrValue {
		var err phpError.Error
		rhs, err = runtimeValueToValueType(values.IntValue, rhs)
		if err != nil {
			return values.NewVoid(), err
		}
	}

	switch rhs.GetType() {
	case values.ArrayValue:
		switch operator {
		case "<", "<=":
			return values.NewBool(true), nil
		case "<=>":
			return values.NewInt(-1), nil
		default:
			return values.NewVoid(), phpError.NewError("compareRelationInteger: Operator \"%s\" not implemented for type array", operator)
		}

	case values.BoolValue:
		lhsBoolean, err := lib_boolval(lhs)
		if err != nil {
			return values.NewVoid(), err
		}
		return compareRelationBoolean(values.NewBool(lhsBoolean), operator, rhs)

	case values.FloatValue:
		lhsFloat, err := lib_floatval(lhs)
		if err != nil {
			return values.NewVoid(), err
		}
		return compareRelationFloating(values.NewFloat(lhsFloat), operator, rhs)

	case values.IntValue:
		rhsInt := rhs.(*values.Int).Value
		switch operator {
		case "<":
			return values.NewBool(lhs.Value < rhsInt), nil
		case "<=":
			return values.NewBool(lhs.Value <= rhsInt), nil
		case "<=>":
			if lhs.Value > rhsInt {
				return values.NewInt(1), nil
			}
			if lhs.Value == rhsInt {
				return values.NewInt(0), nil
			}
			return values.NewInt(-1), nil
		default:
			return values.NewVoid(), phpError.NewError("compareRelationInteger: Operator \"%s\" not implemented", operator)
		}

	default:
		return values.NewVoid(), phpError.NewError("compareRelationInteger: Type \"%s\" not implemented", rhs.GetType())
	}
}

func compareRelationNull(operator string, rhs values.RuntimeValue) (values.RuntimeValue, phpError.Error) {
	// Spec: https://phplang.org/spec/10-expressions.html#grammar-relational-expression

	//       NULL  bool  int  float  string  array  object  resource
	// NULL   =     ->    ->   ->     ->      ->     <       <

	// "=" means the result is always “equals”, i.e. strict comparisons are always FALSE and equality comparisons are always TRUE.

	switch rhs.GetType() {
	case values.ArrayValue:
		lhs, err := lib_arrayval(values.NewNull())
		if err != nil {
			return values.NewVoid(), err
		}
		return compareRelationArray(lhs, operator, rhs)

	case values.BoolValue:
		lhs, err := lib_boolval(values.NewNull())
		if err != nil {
			return values.NewVoid(), err
		}
		return compareRelationBoolean(values.NewBool(lhs), operator, rhs)

	case values.FloatValue:
		lhs, err := lib_floatval(values.NewNull())
		if err != nil {
			return values.NewVoid(), err
		}
		return compareRelationFloating(values.NewFloat(lhs), operator, rhs)

	case values.IntValue:
		lhs, err := lib_intval(values.NewNull())
		if err != nil {
			return values.NewVoid(), err
		}
		return compareRelationInteger(values.NewInt(lhs), operator, rhs)

	case values.NullValue:
		switch operator {
		case "<":
			return values.NewBool(false), nil
		case "<=":
			return values.NewBool(true), nil
		case "<=>":
			return values.NewInt(0), nil
		}
		return values.NewVoid(), phpError.NewError("compareRelationNull: Operator \"%s\" not implemented for type NULL", operator)

		// TODO compareRelationNull - object
		// TODO compareRelationNull - resource

	case values.StrValue:
		lhs, err := lib_strval(values.NewNull())
		if err != nil {
			return values.NewVoid(), err
		}
		return compareRelationString(values.NewStr(lhs), operator, rhs)

	default:
		return values.NewVoid(), phpError.NewError("compareRelationNull: Type \"%s\" not implemented", rhs.GetType())
	}
}

func compareRelationString(lhs *values.Str, operator string, rhs values.RuntimeValue) (values.RuntimeValue, phpError.Error) {
	// Spec: https://phplang.org/spec/10-expressions.html#grammar-relational-expression

	//         NULL  bool  int  float  string  array  object  resource
	// string   <-    ->    ->   ->     2, 4    <      3       2

	// TODO compareRelationString - object
	// TODO compareRelationString - resource

	if rhs.GetType() == values.FloatValue || rhs.GetType() == values.IntValue {
		lhsStr := lhs.Value
		if strings.Trim(lhsStr, " \t") == "" {
			switch operator {
			case "<", "<=":
				return values.NewBool(true), nil
			case "<=>":
				return values.NewInt(-1), nil
			default:
				return values.NewVoid(), phpError.NewError("compareRelationInteger: Operator \"%s\" not implemented for type array", operator)
			}
		}
		if !common.IsIntegerLiteralWithSign(lhsStr) && !common.IsFloatingLiteralWithSign(lhsStr) {
			switch operator {
			case "<", "<=":
				return values.NewBool(false), nil
			case "<=>":
				return values.NewInt(1), nil
			default:
				return values.NewVoid(), phpError.NewError("compareRelationInteger: Operator \"%s\" not implemented for type array", operator)
			}
		}
	}

	if rhs.GetType() == values.NullValue {
		var err phpError.Error
		rhs, err = runtimeValueToValueType(values.StrValue, rhs)
		if err != nil {
			return values.NewVoid(), err
		}
	}

	switch rhs.GetType() {
	case values.ArrayValue:
		switch operator {
		case "<", "<=":
			return values.NewBool(true), nil
		case "<=>":
			return values.NewInt(-1), nil
		default:
			return values.NewVoid(), phpError.NewError("compareRelationArray: Operator \"%s\" not implemented", operator)
		}

	case values.BoolValue:
		lhs, err := lib_boolval(lhs)
		if err != nil {
			return values.NewVoid(), err
		}
		return compareRelationBoolean(values.NewBool(lhs), operator, rhs)

	case values.FloatValue:
		lhs, err := lib_floatval(lhs)
		if err != nil {
			return values.NewVoid(), err
		}
		return compareRelationFloating(values.NewFloat(lhs), operator, rhs)

	case values.IntValue:
		lhs, err := lib_intval(lhs)
		if err != nil {
			return values.NewVoid(), err
		}
		return compareRelationInteger(values.NewInt(lhs), operator, rhs)

	case values.StrValue:
		// Spec: https://phplang.org/spec/10-expressions.html#grammar-relational-expression
		//   2. If one of the operands [...] is a [...] numeric string,
		//      which can be represented as int or float without loss of precision,
		//      the operands are converted to the corresponding arithmetic type, with float taking precedence over int,
		//      and resources converting to int. The result is the numerical comparison of the two operands after conversion.
		rhsStr := rhs.(*values.Str).Value
		if common.IsFloatingLiteralWithSign(lhs.Value) && (common.IsIntegerLiteralWithSign(rhsStr) || common.IsFloatingLiteralWithSign(rhsStr)) {
			lhs, err := lib_floatval(lhs)
			if err != nil {
				return values.NewVoid(), err
			}
			return compareRelationFloating(values.NewFloat(lhs), operator, rhs)
		}
		if common.IsIntegerLiteralWithSign(lhs.Value) && (common.IsIntegerLiteralWithSign(rhsStr) || common.IsFloatingLiteralWithSign(rhsStr)) {
			lhs, err := lib_intval(lhs)
			if err != nil {
				return values.NewVoid(), err
			}
			return compareRelationInteger(values.NewInt(lhs), operator, rhs)
		}

		// Spec: https://phplang.org/spec/10-expressions.html#grammar-relational-expression
		//   4. If both operands are non-numeric strings, the result is the lexical comparison of the two operands.
		//      Specifically, the strings are compared byte-by-byte starting with their first byte.
		//      If the two bytes compare equal and there are no more bytes in either string, the strings are equal and the comparison ends;
		//      otherwise, if this is the final byte in one string, the shorter string compares less-than the longer string and the comparison ends.
		//      If the two bytes compare unequal, the string having the lower-valued byte compares less-than the other string, and the comparison ends.
		//      If there are more bytes in the strings, the process is repeated for the next pair of bytes.
		var result int64 = 0
		for index, lhsByte := range []byte(lhs.Value) {
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
		if result == 0 && len(lhs.Value) < len(rhsStr) {
			result = -1
		}
		switch operator {
		case "<":
			return values.NewBool(result == -1), nil
		case "<=":
			return values.NewBool(result < 1), nil
		case "<=>":
			return values.NewInt(result), nil
		default:
			return values.NewVoid(), phpError.NewError("compareRelationString: Operator \"%s\" not implemented", operator)
		}

	default:
		return values.NewVoid(), phpError.NewError("compareRelationString: Type \"%s\" not implemented", rhs.GetType())
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

func compare(lhs values.RuntimeValue, operator string, rhs values.RuntimeValue) (*values.Bool, phpError.Error) {
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
			return values.NewBool(false), err
		}
		result := resultRuntimeValue.(*values.Int).Value == 0

		if operator == "!=" {
			return values.NewBool(!result), nil
		} else {
			return values.NewBool(result), nil
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
			case values.ArrayValue:
				lhsArray := lhs.(*values.Array)
				rhsArray := rhs.(*values.Array)
				if len(lhsArray.Keys) != len(rhsArray.Keys) {
					result = false
				} else {
					for _, key := range lhsArray.Keys {
						lhsValue, found := lhsArray.GetElement(key)
						if !found {
							result = false
							break
						}
						rhsValue, found := rhsArray.GetElement(key)
						if !found {
							result = false
							break
						}
						equal, err := compare(lhsValue, "===", rhsValue)
						if err != nil {
							return values.NewBool(false), err
						}
						if !equal.Value {
							result = false
							break
						}
					}
				}
			case values.BoolValue:
				result = lhs.(*values.Bool).Value == rhs.(*values.Bool).Value
			case values.FloatValue:
				result = lhs.(*values.Float).Value == rhs.(*values.Float).Value
			case values.IntValue:
				result = lhs.(*values.Int).Value == rhs.(*values.Int).Value
			case values.NullValue:
				result = true
			case values.StrValue:
				result = lhs.(*values.Str).Value == rhs.(*values.Str).Value
			default:
				return values.NewBool(false), phpError.NewError("compare: Runtime type %s for operator \"===\" not implemented", lhs.GetType())
			}
		}

		if operator == "!==" {
			return values.NewBool(!result), nil
		} else {
			return values.NewBool(result), nil
		}
	}

	return values.NewBool(false), phpError.NewError("compare: Operator \"%s\" not implemented", operator)
}
