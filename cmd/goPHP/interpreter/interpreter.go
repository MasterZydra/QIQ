package interpreter

import (
	"GoPHP/cmd/goPHP/ast"
	"GoPHP/cmd/goPHP/parser"
	"strings"
)

type Interpreter struct {
	config   *Config
	request  *Request
	parser   *parser.Parser
	env      *Environment
	cache    map[int64]IRuntimeValue
	result   string
	exitCode int64
}

func NewInterpreter(config *Config, request *Request) *Interpreter {
	return &Interpreter{config: config, request: request, parser: parser.NewParser(),
		env: NewEnvironment(nil, request), cache: map[int64]IRuntimeValue{},
		exitCode: 0,
	}
}

func (interpreter *Interpreter) GetExitCode() int {
	return int(interpreter.exitCode)
}

func (interpreter *Interpreter) Process(sourceCode string) (string, Error) {
	return interpreter.process(sourceCode, interpreter.env)
}

func (interpreter *Interpreter) process(sourceCode string, env *Environment) (string, Error) {
	interpreter.result = ""
	program, err := interpreter.parser.ProduceAST(sourceCode)
	if err != nil {
		return interpreter.result, NewParseError(err)
	}

	runtimeErr := interpreter.scanForFunctionDefinition(program.GetStatements(), env)
	if runtimeErr != nil {
		return interpreter.result, runtimeErr
	}

	for _, stmt := range program.GetStatements() {
		if _, err := interpreter.processStmt(stmt, env); err != nil {
			// Handle exit event - Stop code execution
			if err.GetErrorType() == EventError && err.GetMessage() == ExitEvent {
				break
			}
			return interpreter.result, err
		}
	}

	return strings.ReplaceAll(interpreter.result, "\n\n", "\n"), nil
}

func (interpreter *Interpreter) processStmt(stmt ast.IStatement, env *Environment) (IRuntimeValue, Error) {
	switch stmt.GetKind() {
	// Statements
	case ast.ConstDeclarationStmt:
		return interpreter.processConstDeclarationStmt(ast.StmtToConstDeclStatement(stmt), env)
	case ast.CompoundStmt:
		return interpreter.processCompoundStmt(ast.StmtToCompoundStatement(stmt), env)
	case ast.EchoStmt:
		return interpreter.processEchoStmt(ast.StmtToEchoStatement(stmt), env)
	case ast.ExpressionStmt:
		return interpreter.processStmt(ast.StmtToExprStatement(stmt).GetExpression(), env)
	case ast.FunctionDefinitionStmt:
		return interpreter.processFunctionDefinitionStmt(ast.StmtToFunctionDefinitionStatement(stmt), env)
	case ast.IfStmt:
		return interpreter.processIfStmt(ast.StmtToIfStatement(stmt), env)

	// Expressions
	case ast.ArrayLiteralExpr, ast.IntegerLiteralExpr, ast.FloatingLiteralExpr, ast.StringLiteralExpr:
		return interpreter.exprToRuntimeValue(stmt, env)
	case ast.TextNode:
		interpreter.print(ast.ExprToTextExpr(stmt).GetValue())
		return NewVoidRuntimeValue(), nil
	case ast.SimpleVariableExpr:
		return interpreter.processSimpleVariableExpression(ast.ExprToSimpleVarExpr(stmt), env)
	case ast.SimpleAssignmentExpr:
		return interpreter.processSimpleAssignmentExpression(ast.ExprToSimpleAssignExpr(stmt), env)
	case ast.SubscriptExpr:
		return interpreter.processSubscriptExpression(ast.ExprToSubscriptExpr(stmt), env)
	case ast.FunctionCallExpr:
		return interpreter.processFunctionCallExpression(ast.ExprToFuncCallExpr(stmt), env)
	case ast.EmptyIntrinsicExpr:
		return interpreter.processEmptyIntrinsicExpression(ast.ExprToFuncCallExpr(stmt), env)
	case ast.ExitIntrinsicExpr:
		return interpreter.processExitIntrinsicExpression(ast.ExprToFuncCallExpr(stmt), env)
	case ast.IssetIntrinsicExpr:
		return interpreter.processIssetExpression(ast.ExprToFuncCallExpr(stmt), env)
	case ast.UnsetIntrinsicExpr:
		return interpreter.processUnsetExpression(ast.ExprToFuncCallExpr(stmt), env)
	case ast.ConstantAccessExpr:
		return interpreter.processConstantAccessExpression(ast.ExprToConstAccessExpr(stmt), env)
	case ast.CompoundAssignmentExpr:
		return interpreter.processCompoundAssignmentExpression(ast.ExprToCompoundAssignExpr(stmt), env)
	case ast.ConditionalExpr:
		return interpreter.processConditionalExpression(ast.ExprToCondExpr(stmt), env)
	case ast.CoalesceExpr:
		return interpreter.processCoalesceExpression(ast.ExprToCoalesceExpr(stmt), env)
	case ast.RelationalExpr:
		return interpreter.processRelationalExpression(ast.ExprToBinOpExpr(stmt), env)
	case ast.EqualityExpr:
		return interpreter.processEqualityExpression(ast.ExprToBinOpExpr(stmt), env)
	case ast.BinaryOpExpr:
		return interpreter.processAdditiveExpression(ast.ExprToBinOpExpr(stmt), env)
	case ast.UnaryOpExpr:
		return interpreter.processUnaryExpression(ast.ExprToUnaryOpExpr(stmt), env)
	case ast.CastExpr:
		return interpreter.processCastExpression(ast.ExprToUnaryOpExpr(stmt), env)
	case ast.LogicalNotExpr:
		return interpreter.processLogicalNotExpression(ast.ExprToUnaryOpExpr(stmt), env)
	case ast.PostfixIncExpr:
		return interpreter.processPostfixIncExpression(ast.ExprToUnaryOpExpr(stmt), env)
	case ast.PrefixIncExpr:
		return interpreter.processPrefixIncExpression(ast.ExprToUnaryOpExpr(stmt), env)

	default:
		return NewVoidRuntimeValue(), NewError("Unsupported statement or expression: %s", stmt)
	}
}

func (interpreter *Interpreter) processConstDeclarationStmt(stmt ast.IConstDeclarationStatement, env *Environment) (IRuntimeValue, Error) {
	value, err := interpreter.processStmt(stmt.GetValue(), env)
	if err != nil {
		return NewVoidRuntimeValue(), err
	}
	return env.declareConstant(stmt.GetName(), value)
}

func (interpreter *Interpreter) processCompoundStmt(stmt ast.ICompoundStatement, env *Environment) (IRuntimeValue, Error) {
	for _, statement := range stmt.GetStatements() {
		_, err := interpreter.processStmt(statement, env)
		if err != nil {
			return NewVoidRuntimeValue(), err
		}
	}
	return NewVoidRuntimeValue(), nil
}

func (interpreter *Interpreter) processEchoStmt(stmt ast.IEchoStatement, env *Environment) (IRuntimeValue, Error) {
	for _, expr := range stmt.GetExpressions() {
		if runtimeValue, err := interpreter.processStmt(expr, env); err != nil {
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

func (interpreter *Interpreter) processFunctionDefinitionStmt(stmt ast.IFunctionDefinitionStatement, env *Environment) (IRuntimeValue, Error) {
	// Check if this function definition was already processed before interpreting the code
	if interpreter.isCached(stmt) {
		return NewVoidRuntimeValue(), nil
	}

	function := userFunction{
		FunctionName: stmt.GetFunctionName(), Parameters: stmt.GetParams(), Body: stmt.GetBody(),
	}
	if err := env.defineUserFunction(function); err != nil {
		return NewVoidRuntimeValue(), err
	}

	return interpreter.writeCache(stmt, NewVoidRuntimeValue()), nil
}

func (interpreter *Interpreter) processIfStmt(stmt ast.IIfStatement, env *Environment) (IRuntimeValue, Error) {
	conditionRuntimeValue, err := interpreter.processStmt(stmt.GetCondition(), env)
	if err != nil {
		return NewVoidRuntimeValue(), err
	}

	condition, err := lib_boolval(conditionRuntimeValue)
	if err != nil {
		return NewVoidRuntimeValue(), err
	}

	if condition {
		_, err = interpreter.processStmt(stmt.GetIfBlock(), env)
		if err != nil {
			return NewVoidRuntimeValue(), err
		}
		return NewVoidRuntimeValue(), nil
	}

	if len(stmt.GetElseIf()) > 0 {
		for _, elseIf := range stmt.GetElseIf() {
			conditionRuntimeValue, err := interpreter.processStmt(elseIf.GetCondition(), env)
			if err != nil {
				return NewVoidRuntimeValue(), err
			}

			condition, err := lib_boolval(conditionRuntimeValue)
			if err != nil {
				return NewVoidRuntimeValue(), err
			}

			if !condition {
				continue
			}

			_, err = interpreter.processStmt(elseIf.GetIfBlock(), env)
			if err != nil {
				return NewVoidRuntimeValue(), err
			}
			return NewVoidRuntimeValue(), nil
		}
	}

	if stmt.GetElseBlock() != nil {
		_, err = interpreter.processStmt(stmt.GetElseBlock(), env)
		if err != nil {
			return NewVoidRuntimeValue(), err
		}
		return NewVoidRuntimeValue(), nil
	}

	return NewVoidRuntimeValue(), nil
}

func (interpreter *Interpreter) processSimpleVariableExpression(expr ast.ISimpleVariableExpression, env *Environment) (IRuntimeValue, Error) {
	return interpreter.lookupVariable(expr, env, false)
}

func (interpreter *Interpreter) processSimpleAssignmentExpression(expr ast.ISimpleAssignmentExpression, env *Environment) (IRuntimeValue, Error) {
	if !ast.IsVariableExpression(expr.GetVariable()) {
		return NewVoidRuntimeValue(),
			NewError("processSimpleAssignmentExpression: Invalid variable: %s", expr.GetVariable())
	}

	value, err := interpreter.processStmt(expr.GetValue(), env)
	if err != nil {
		return NewVoidRuntimeValue(), err
	}

	variableName, err := interpreter.varExprToVarName(expr.GetVariable(), env)
	if err != nil {
		return NewVoidRuntimeValue(), err
	}

	currentValue, _ := env.lookupVariable(variableName)

	if currentValue.GetType() == ArrayValue {
		if expr.GetVariable().GetKind() != ast.SubscriptExpr {
			return NewVoidRuntimeValue(), NewError("processSimpleAssignmentExpression: Unsupported variable type %s", expr.GetVariable().GetKind())
		}

		var key ast.IExpression = ast.ExprToSubscriptExpr(expr.GetVariable()).GetIndex()

		array := runtimeValToArrayRuntimeVal(currentValue)
		if key == nil {
			var lastIndex IRuntimeValue = NewIntegerRuntimeValue(0)
			if len(array.GetKeys()) > 0 {
				lastIndex = array.GetKeys()[len(array.GetKeys())-1]
			}
			if lastIndex.GetType() != IntegerValue {
				return NewVoidRuntimeValue(), NewError("processSimpleAssignmentExpression: Unsupported array key %s", lastIndex.GetType())
			}
			var nextIndex = lastIndex
			if len(array.GetKeys()) > 0 {
				nextIndex = NewIntegerRuntimeValue(runtimeValToIntRuntimeVal(lastIndex).GetValue() + 1)
			}
			array.SetElement(nextIndex, value)
		} else {
			keyValue, err := interpreter.processStmt(key, env)
			if err != nil {
				return NewVoidRuntimeValue(), err
			}
			array.SetElement(keyValue, value)
		}

		return value, nil
	}

	return env.declareVariable(variableName, value)
}

func (interpreter *Interpreter) processSubscriptExpression(expr ast.ISubscriptExpression, env *Environment) (IRuntimeValue, Error) {
	// Spec: https://phplang.org/spec/10-expressions.html#grammar-subscript-expression

	variable, err := interpreter.lookupVariable(expr.GetVariable(), env, false)
	if err != nil {
		return NewVoidRuntimeValue(), err
	}

	// Spec: https://phplang.org/spec/10-expressions.html#grammar-subscript-expression
	// dereferencable-expression designates an array
	if variable.GetType() == ArrayValue {
		// TODO processSubscriptExpression - no key
		// Spec: https://phplang.org/spec/10-expressions.html#grammar-subscript-expression
		// If expression is omitted, a new element is inserted. Its key has type int and is one more than the highest, previously assigned int key for this array. If this is the first element with an int key, key 0 is used. If the largest previously assigned int key is the largest integer value that can be represented, the new element is not added. The result is the added new element, or NULL if the element was not added.

		key, err := interpreter.processStmt(expr.GetIndex(), env)
		if err != nil {
			return NewVoidRuntimeValue(), err
		}

		array := runtimeValToArrayRuntimeVal(variable)

		exists, err := lib_array_key_exists(key, array)
		if err != nil {
			return NewVoidRuntimeValue(), err
		}

		// Spec: https://phplang.org/spec/10-expressions.html#grammar-subscript-expression
		// If expression is present, if the designated element exists,
		// the type and value of the result is the type and value of that element;
		// otherwise, the result is NULL.
		if exists {
			element, _ := array.GetElement(key)
			return element, nil
		} else {
			return NewNullRuntimeValue(), nil
		}

		// TODO processSubscriptExpression
		// If the usage context is as the left-hand side of a simple-assignment-expression, the value of the new element is the value of the right-hand side of that simple-assignment-expression.
		// If the usage context is as the left-hand side of a compound-assignment-expression: the expression e1 op= e2 is evaluated as e1 = NULL op (e2).
		// If the usage context is as the operand of a postfix- or prefix-increment or decrement operator, the value of the new element is considered to be NULL.
	}

	return NewVoidRuntimeValue(), NewError("Unsupported subscript expression: %s", expr)

	/*
	   If dereferencable-expression designates a string, expression must not designate a string.

	   expression can be omitted only if subscript-expression is used in a modifiable-lvalue context and dereferencable-expression does not designate a string. Exception from this is when dereferencable-expression is an empty string - then it is converted to an empty array.

	   If subscript-expression is used in a non-lvalue context, the element being designated must exist.

	   Semantics

	   A subscript-expression designates a (possibly non-existent) element of an array or string. When subscript-expression designates an object of a type that implements ArrayAccess, the minimal semantics are defined below; however, they can be augmented by that object’s methods offsetGet and offsetSet.

	   The element key is designated by expression. If the type of element-key is neither int nor string, keys with float or bool values, or strings whose contents match exactly the pattern of decimal-literal, are converted to integer, and key values of all other types are converted to string.

	   If both dereferencable-expression and expression designate strings, expression is treated as if it specified the int key zero instead and a non-fatal error is produces.

	   A subscript-expression designates a modifiable lvalue if and only if dereferencable-expression designates a modifiable lvalue.

	   dereferencable-expression designates a string

	   The expression is converted to int and the result is the character of the string at the position corresponding to that integer. If the integer is negative, the position is counted backwards from the end of the string. If the position refers to a non-existing offset, the result is an empty string.

	   If the operator is used as the left-hand side of a simple-assignment-expression,

	       If the assigned string is empty, or in case of non-existing negative offset (absolute value larger than string length), a warning is raised and no assignment is performed.
	       If the offset is larger than the current string length, the string is extended to a length equal to the offset value, using space (0x20) padding characters.
	       The value being assigned is converted to string and the character in the specified offset is replaced by the first character of the string.

	   The subscript operator can not be used on a string value in a byRef context or as the operand of the postfix- or prefix-increment or decrement operators or on the left side of compound-assignment-expression, doing so will result in a fatal error.

	   dereferencable-expression designates an object of a type that implements ArrayAccess

	   If expression is present,

	       If subscript-expression is used in a non-lvalue context, the object’s method offsetGet is called with an argument of expression. The return value of the offsetGet is the result.
	       If the usage context is as the left-hand side of a simple-assignment-expression, the object’s method offsetSet is called with a first argument of expression and a second argument that is the value of the right-hand side of that simple-assignment-expression. The value of the right-hand side is the result.
	       If the usage context is as the left-hand side of a compound-assignment-expression, the expression e1[e] op= e2 is evaluated as e1[e] = e1->offsetGet(e) op (e2), which is then processed according to the rules for simple assignment immediately above.
	       If the usage context is as the operand of the postfix- or prefix-increment or decrement operators, the object’s method offsetGet is called with an argument of expression. However, this method has no way of knowing if an increment or decrement operator was used, or whether it was a prefix or postfix operator. In order for the value to be modified by the increment/decrement, offsetGet must return byRef. The result of the subscript operator value returned by offsetGet.

	   If expression is omitted,

	       If the usage context is as the left-hand side of a simple-assignment-expression, the object’s method offsetSet is called with a first argument of NULL and a second argument that is the value of the right-hand side of that simple-assignment-expression. The type and value of the result is the type and value of the right-hand side of that simple-assignment-expression.
	       If the usage context is as the left-hand side of a compound-assignment-expression: The expression e1[] op= e2 is evaluated as e1[] = e1->offsetGet(NULL) op (e2), which is then processed according to the rules for simple assignment immediately above.
	       If the usage context is as the operand of the postfix- or prefix-increment or decrement operators, the object’s method offsetGet is called with an argument of NULL. However, this method has no way of knowing if an increment or decrement operator was used, or whether it was a prefix or postfix operator. In order for the value to be modified by the increment/decrement, offsetGet must return byRef. The result of the subscript operator value returned by offsetGet.

	   Note: The brace ({...}) form of this operator has been deprecated.
	*/

	// Examples
	/*
	   $v = array(10, 20, 30);
	   $v[1] = 1.234;    // change the value (and type) of element [1]
	   $v[-10] = 19;   // insert a new element with int key -10
	   $v["red"] = TRUE; // insert a new element with string key "red"
	   [[2,4,6,8], [5,10], [100,200,300]][0][2]  // designates element with value 6
	   ["black", "white", "yellow"][1][2]  // designates substring "i" in "white"
	   function f() { return [1000, 2000, 3000]; }
	   f()[2];      // designates element with value 3000
	   "red"[1.9];    // designates "e"
	   "red"[-2];    // designates "e"
	   "red"[0][0][0];    // designates "r"
	   // -----------------------------------------
	   class MyVector implements ArrayAccess { /* ... */ /*}
	$vect1 = new MyVector(array(10, 'A' => 2.3, "up"));
	$vect1[10] = 987; // calls Vector::offsetSet(10, 987)
	$vect1[] = "xxx"; // calls Vector::offsetSet(NULL, "xxx")
	$x = $vect1[1];   // calls Vector::offsetGet(1)
	*/
}

func (interpreter *Interpreter) processFunctionCallExpression(expr ast.IFunctionCallExpression, env *Environment) (IRuntimeValue, Error) {
	// Lookup native function
	nativeFunction, err := env.lookupNativeFunction(expr.GetFunctionName())
	if err == nil {
		functionArguments := make([]IRuntimeValue, len(expr.GetArguments()))
		for index, arg := range expr.GetArguments() {
			runtimeValue, err := interpreter.processStmt(arg, env)
			if err != nil {
				return NewVoidRuntimeValue(), err
			}
			functionArguments[index] = deepCopy(runtimeValue)
		}
		// TODO processFunctionCallExpression - Check return type
		return nativeFunction(functionArguments, interpreter)
	}

	// Lookup user function
	userFunction, err := env.lookupUserFunction(expr.GetFunctionName())
	if err != nil {
		return NewVoidRuntimeValue(), err
	}

	functionEnv := NewEnvironment(env, nil)

	if len(userFunction.Parameters) != len(expr.GetArguments()) {
		return NewVoidRuntimeValue(), NewError(
			"Uncaught ArgumentCountError: %s() expects exactly %d arguments, %d given",
			userFunction.FunctionName, len(userFunction.Parameters), len(expr.GetArguments()),
		)
	}
	for index, param := range userFunction.Parameters {
		runtimeValue, err := interpreter.processStmt(expr.GetArguments()[index], env)
		if err != nil {
			return NewVoidRuntimeValue(), err
		}
		// Check if the parameter types match
		err = checkParameterTypes(runtimeValue, param.Type)
		if err != nil && err.GetMessage() == "Types do not match" {
			givenType, err := lib_gettype(runtimeValue)
			if err != nil {
				return NewVoidRuntimeValue(), err
			}
			return NewVoidRuntimeValue(), NewError(
				"Uncaught TypeError: %s(): Argument #%d (%s) must be of type %s, %s given",
				userFunction.FunctionName, index+1, param.Name, strings.Join(param.Type, "|"), givenType,
			)
		}
		// Declare parameter in function environment
		functionEnv.declareVariable(param.Name, runtimeValue)
	}

	runtimeValue, err := interpreter.processStmt(userFunction.Body, functionEnv)
	if err != nil {
		return NewVoidRuntimeValue(), err
	}
	return runtimeValue, nil

}

func (interpreter *Interpreter) processEmptyIntrinsicExpression(expr ast.IFunctionCallExpression, env *Environment) (IRuntimeValue, Error) {
	// Spec: https://phplang.org/spec/10-expressions.html#grammar-empty-intrinsic

	// This intrinsic returns TRUE if the variable or value designated by expression is empty,
	// where empty means that the variable designated by it does not exist, or it exists and its value compares equal to FALSE.
	// Otherwise, the intrinsic returns FALSE.

	// The following values are considered empty:
	// FALSE, 0, 0.0, "" (empty string), "0", NULL, an empty array, and any uninitialized variable.

	// If this intrinsic is used with an expression that designates a dynamic property,
	// then if the class of that property has an __isset, that method is called.
	// If that method returns TRUE, the value of the property is retrieved (which may call __get if defined)
	// and compared to FALSE as described above. Otherwise, the result is FALSE.

	var runtimeValue IRuntimeValue
	var err Error
	if ast.IsVariableExpression(expr.GetArguments()[0]) {
		runtimeValue, err = interpreter.lookupVariable(expr.GetArguments()[0], env, true)
		if err != nil {
			return NewBooleanRuntimeValue(true), nil
		}
	} else {
		runtimeValue, err = interpreter.processStmt(expr.GetArguments()[0], env)
		if err != nil {
			return NewVoidRuntimeValue(), err
		}
	}

	boolean, err := lib_boolval(runtimeValue)
	if err != nil {
		return NewVoidRuntimeValue(), err
	}
	return NewBooleanRuntimeValue(!boolean), nil
}

func (interpreter *Interpreter) processExitIntrinsicExpression(expr ast.IFunctionCallExpression, env *Environment) (IRuntimeValue, Error) {
	// Spec: https://phplang.org/spec/10-expressions.html#grammar-exit-intrinsic

	// "exit" and "die" are equivalent.

	// This intrinsic terminates the current script.
	// If expression designates a string, that string is written to STDOUT.
	// If expression designates an integer, that represents the script’s exit status code. Code 255 is reserved by PHP.
	// Code 0 represents “success”. The exit status code is made available to the execution environment.
	// If expression is omitted or is a string, the exit status code is zero. exit does not have a resulting value.

	// "exit" performs the following operations, in order:
	//   1. Writes the optional string to STDOUT.
	//   2. Calls any functions registered via the library function register_shutdown_function in their order of registration.
	//   3. Invokes destructors for all remaining instances.

	expression := expr.GetArguments()[0]
	if expression != nil {
		exprValue, err := interpreter.processStmt(expression, env)
		if err != nil {
			return NewVoidRuntimeValue(), err
		}
		if exprValue.GetType() == StringValue {
			interpreter.print(runtimeValToStrRuntimeVal(exprValue).GetValue())
		}
		if exprValue.GetType() == IntegerValue {
			exitCode := runtimeValToIntRuntimeVal(exprValue).GetValue()
			if exitCode >= 0 && exitCode < 255 {
				interpreter.exitCode = exitCode
			}
		}
	}

	// TODO processExitIntrinsicExpression - call shutdown functions
	// TODO processExitIntrinsicExpression - call destructors

	return NewVoidRuntimeValue(), NewEvent(ExitEvent)
}

func (interpreter *Interpreter) processIssetExpression(expr ast.IFunctionCallExpression, env *Environment) (IRuntimeValue, Error) {
	// Spec: https://phplang.org/spec/10-expressions.html#grammar-isset-intrinsic

	// This intrinsic returns TRUE if all the variables designated by variabless are set and their values are not NULL.
	// Otherwise, it returns FALSE.

	// If this intrinsic is used with an expression that designate a dynamic property,
	// then if the class of that property has an __isset, that method is called.
	// If that method returns TRUE, the value of the property is retrieved (which may call __get if defined)
	// and if it is not NULL, the result is TRUE. Otherwise, the result is FALSE.

	for _, arg := range expr.GetArguments() {
		runtimeValue, _ := interpreter.lookupVariable(arg, env, true)
		if runtimeValue.GetType() == NullValue {
			return NewBooleanRuntimeValue(false), nil
		}
	}
	return NewBooleanRuntimeValue(true), nil
}

func (interpreter *Interpreter) processUnsetExpression(expr ast.IFunctionCallExpression, env *Environment) (IRuntimeValue, Error) {
	// Spec: https://phplang.org/spec/11-statements.html#grammar-unset-statement

	// This statement unsets the variables designated by each variable in variable-list. No value is returned.
	// An attempt to unset a non-existent variable (such as a non-existent element in an array) is ignored.

	for _, arg := range expr.GetArguments() {
		variableName, err := interpreter.varExprToVarName(arg, env)
		if err != nil {
			return NewVoidRuntimeValue(), err
		}
		env.unsetVariable(variableName)
	}
	return NewVoidRuntimeValue(), nil
}

func (interpreter *Interpreter) processConstantAccessExpression(expr ast.IConstantAccessExpression, env *Environment) (IRuntimeValue, Error) {
	return env.lookupConstant(expr.GetConstantName())
}

func (interpreter *Interpreter) processCompoundAssignmentExpression(expr ast.ICompoundAssignmentExpression, env *Environment) (IRuntimeValue, Error) {
	if !ast.IsVariableExpression(expr.GetVariable()) {
		return NewVoidRuntimeValue(),
			NewError("processCompoundAssignmentExpression: Invalid variable: %s", expr.GetVariable())
	}

	operand1, err := interpreter.processStmt(expr.GetVariable(), env)
	if err != nil {
		return NewVoidRuntimeValue(), err
	}

	operand2, err := interpreter.processStmt(expr.GetValue(), env)
	if err != nil {
		return NewVoidRuntimeValue(), err
	}

	newValue, err := calculate(operand1, expr.GetOperator(), operand2)
	if err != nil {
		return NewVoidRuntimeValue(), err
	}

	variableName, err := interpreter.varExprToVarName(expr.GetVariable(), env)
	if err != nil {
		return NewVoidRuntimeValue(), err
	}

	return env.declareVariable(variableName, newValue)
}

func (interpreter *Interpreter) processConditionalExpression(expr ast.IConditionalExpression, env *Environment) (IRuntimeValue, Error) {
	// Spec: https://phplang.org/spec/10-expressions.html#grammar-conditional-expression

	// Spec: https://phplang.org/spec/10-expressions.html#grammar-conditional-expression
	// Given the expression "e1 ? e2 : e3", e1 is evaluated first and converted to bool if it has another type.
	runtimeValue, isConditionTrue, err := interpreter.processCondition(expr.GetCondition(), env)
	if err != nil {
		return NewVoidRuntimeValue(), err
	}

	if isConditionTrue {
		if expr.GetIfExpr() != nil {
			// Spec: https://phplang.org/spec/10-expressions.html#grammar-conditional-expression
			// If the result is TRUE, then and only then is e2 evaluated, and the result and its type become the result
			// and type of the whole expression.
			return interpreter.processStmt(expr.GetIfExpr(), env)
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
		return interpreter.processStmt(expr.GetElseExpr(), env)
	}
}

func (interpreter *Interpreter) processCoalesceExpression(expr ast.ICoalesceExpression, env *Environment) (IRuntimeValue, Error) {
	// Spec: https://phplang.org/spec/10-expressions.html#grammar-coalesce-expression

	// Store current error reporting
	errorReporting := interpreter.config.ErrorReporting
	// Suppress all errors
	interpreter.config.ErrorReporting = 0

	cond, err := interpreter.processStmt(expr.GetCondition(), env)

	// Restore previous error reporting
	interpreter.config.ErrorReporting = errorReporting

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
	return interpreter.processStmt(expr.GetElseExpr(), env)

	// TODO processCoalesceExpression - handle uninitialized variables
	// Spec: https://phplang.org/spec/10-expressions.html#grammar-coalesce-expression
	// Note that the semantics of ?? is similar to isset so that uninitialized variables will not produce warnings when used in e1.
	// TODO use isset here - Steps: Add caching of expression results - map[exprId]IRuntimeValue
}

func (interpreter *Interpreter) processRelationalExpression(expr ast.IBinaryOpExpression, env *Environment) (IRuntimeValue, Error) {
	lhs, err := interpreter.processStmt(expr.GetLHS(), env)
	if err != nil {
		return NewVoidRuntimeValue(), err
	}

	rhs, err := interpreter.processStmt(expr.GetRHS(), env)
	if err != nil {
		return NewVoidRuntimeValue(), err
	}
	return compareRelation(lhs, expr.GetOperator(), rhs)
}

func (interpreter *Interpreter) processEqualityExpression(expr ast.IBinaryOpExpression, env *Environment) (IRuntimeValue, Error) {
	lhs, err := interpreter.processStmt(expr.GetLHS(), env)
	if err != nil {
		return NewVoidRuntimeValue(), err
	}

	rhs, err := interpreter.processStmt(expr.GetRHS(), env)
	if err != nil {
		return NewVoidRuntimeValue(), err
	}
	return compare(lhs, expr.GetOperator(), rhs)
}

func (interpreter *Interpreter) processAdditiveExpression(expr ast.IBinaryOpExpression, env *Environment) (IRuntimeValue, Error) {
	lhs, err := interpreter.processStmt(expr.GetLHS(), env)
	if err != nil {
		return NewVoidRuntimeValue(), err
	}

	rhs, err := interpreter.processStmt(expr.GetRHS(), env)
	if err != nil {
		return NewVoidRuntimeValue(), err
	}
	return calculate(lhs, expr.GetOperator(), rhs)
}

func (interpreter *Interpreter) processUnaryExpression(expr ast.IUnaryOpExpression, env *Environment) (IRuntimeValue, Error) {
	operand, err := interpreter.processStmt(expr.GetExpression(), env)
	if err != nil {
		return NewVoidRuntimeValue(), err
	}
	return calculateUnary(expr.GetOperator(), operand)
}

func (interpreter *Interpreter) processCastExpression(expr ast.IUnaryOpExpression, env *Environment) (IRuntimeValue, Error) {
	// Spec: https://phplang.org/spec/10-expressions.html#grammar-cast-type

	// A cast-type of "unset" is no longer supported and results in a compile-time error.
	// With the exception of the cast-type unset and binary (see below), the value of the operand cast-expression is converted to the type specified by cast-type, and that is the type and value of the result. This construct is referred to as a cast and is used as the verb, “to cast”. If no conversion is involved, the type and value of the result are the same as those of cast-expression.

	// A cast can result in a loss of information.

	// TODO processCastExpression - object
	// A cast-type of "object" results in a conversion to type "object".

	value, err := interpreter.processStmt(expr.GetExpression(), env)
	if err != nil {
		return NewVoidRuntimeValue(), err
	}

	switch expr.GetOperator() {
	case "array":
		// Spec: https://phplang.org/spec/10-expressions.html#grammar-cast-type
		// A cast-type of "array" results in a conversion to type array.
		return lib_arrayval(value)
	case "binary", "string":
		// Spec: https://phplang.org/spec/10-expressions.html#grammar-cast-type
		// A cast-type of "binary" is reserved for future use in dealing with so-called binary strings. For now, it is fully equivalent to "string" cast.
		// A cast-type of "string" results in a conversion to type "string".
		return runtimeValueToValueType(StringValue, value)
	case "bool", "boolean":
		// Spec: https://phplang.org/spec/10-expressions.html#grammar-cast-type
		// A cast-type of "bool" or "boolean" results in a conversion to type "bool".
		return runtimeValueToValueType(BooleanValue, value)
	case "double", "float", "real":
		// Spec: https://phplang.org/spec/10-expressions.html#grammar-cast-type
		// A cast-type of "float", "double", or "real" results in a conversion to type "float".
		return runtimeValueToValueType(FloatingValue, value)
	case "int", "integer":
		// Spec: https://phplang.org/spec/10-expressions.html#grammar-cast-type
		// A cast-type of "int" or "integer" results in a conversion to type "int".
		return runtimeValueToValueType(IntegerValue, value)
	default:
		return NewVoidRuntimeValue(), NewError("processCastExpression: Unsupported cast type %s", expr.GetOperator())
	}
}

func (interpreter *Interpreter) processLogicalNotExpression(expr ast.IUnaryOpExpression, env *Environment) (IRuntimeValue, Error) {
	runtimeValue, err := interpreter.processStmt(expr.GetExpression(), env)
	if err != nil {
		return NewVoidRuntimeValue(), err
	}
	boolValue, err := lib_boolval(runtimeValue)
	if err != nil {
		return NewVoidRuntimeValue(), err
	}
	return NewBooleanRuntimeValue(!boolValue), nil
}

func (interpreter *Interpreter) processPostfixIncExpression(expr ast.IUnaryOpExpression, env *Environment) (IRuntimeValue, Error) {
	// Spec: https://phplang.org/spec/10-expressions.html#postfix-increment-and-decrement-operators
	// These operators behave like their prefix counterparts except that the value of a postfix ++ or – expression is the value
	// before any increment or decrement takes place.

	previous, err := interpreter.processStmt(expr.GetExpression(), env)
	if err != nil {
		return NewVoidRuntimeValue(), err
	}

	newValue, err := calculateIncDec(expr.GetOperator(), previous)
	if err != nil {
		return NewVoidRuntimeValue(), err
	}

	variableName, err := interpreter.varExprToVarName(expr.GetExpression(), env)
	if err != nil {
		return NewVoidRuntimeValue(), err
	}
	env.declareVariable(variableName, newValue)

	return previous, nil
}

func (interpreter *Interpreter) processPrefixIncExpression(expr ast.IUnaryOpExpression, env *Environment) (IRuntimeValue, Error) {
	previous, err := interpreter.processStmt(expr.GetExpression(), env)
	if err != nil {
		return NewVoidRuntimeValue(), err
	}

	newValue, err := calculateIncDec(expr.GetOperator(), previous)
	if err != nil {
		return NewVoidRuntimeValue(), err
	}

	variableName, err := interpreter.varExprToVarName(expr.GetExpression(), env)
	if err != nil {
		return NewVoidRuntimeValue(), err
	}
	env.declareVariable(variableName, newValue)

	return newValue, nil
}
