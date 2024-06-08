package interpreter

import (
	"GoPHP/cmd/goPHP/ast"
	"fmt"
	"regexp"
	"strings"
)

type Environment struct {
	parent    *Environment
	variables map[string]IRuntimeValue
	constants map[string]IRuntimeValue
	functions map[string]userFunction
	// StdLib
	predefinedVariables map[string]IRuntimeValue
	predefinedConstants map[string]IRuntimeValue
	nativeFunctions     map[string]nativeFunction
}

func NewEnvironment(parentEnv *Environment, request *Request) *Environment {
	env := &Environment{
		parent:    parentEnv,
		variables: map[string]IRuntimeValue{},
		constants: map[string]IRuntimeValue{},
		functions: map[string]userFunction{},
		// StdLib
		predefinedVariables: map[string]IRuntimeValue{},
		predefinedConstants: map[string]IRuntimeValue{},
		nativeFunctions:     map[string]nativeFunction{},
	}

	if parentEnv == nil {
		registerNativeFunctions(env)
		registerPredefinedVariables(env, request)
		registerPredefinedConstants(env)
	}

	return env
}

func registerPredefinedVariables(environment *Environment, request *Request) {
	environment.predefinedVariables["$_ENV"] = stringMapToArray(request.Env)
	environment.predefinedVariables["$_SERVER"] = environment.predefinedVariables["$_ENV"]
	environment.predefinedVariables["$_GET"] = paramToArray(request.GetParams)
	environment.predefinedVariables["$_POST"] = paramToArray(request.PostParams)
}

func stringMapToArray(stringMap map[string]string) IArrayRuntimeValue {
	result := NewArrayRuntimeValue()
	for key, value := range stringMap {
		result.SetElement(NewStringRuntimeValue(key), NewStringRuntimeValue(value))
	}
	return result
}

func paramToArray(params [][]string) IArrayRuntimeValue {
	result := NewArrayRuntimeValue()

	for _, param := range params {
		key := param[0]
		value := param[1]

		// No array
		if !strings.Contains(key, "]") {
			result.SetElement(NewStringRuntimeValue(key), NewStringRuntimeValue(value))
			continue
		}

		// Array

		openingBracket := strings.Index(key, "[")
		// Get name of param without brackets
		paramName := key[:openingBracket]

		// Check if array is already in params
		arrayValue, found := result.GetElement(NewStringRuntimeValue(paramName))
		if !found || arrayValue.GetType() != ArrayValue {
			arrayValue = NewArrayRuntimeValue()
		}

		// Wrap keys that are strings in double quotes
		decimalKeys, _ := regexp.Compile(`\[[0-9]+\]`)
		nondecimalKeys, _ := regexp.Compile(`\[.+\]`)
		matches := nondecimalKeys.FindAllString(decimalKeys.ReplaceAllString(key, ""), -1)
		for _, match := range matches {
			replacement := `["` + match[1:len(match)-1] + `"]`
			key = strings.Replace(key, match, replacement, 1)
		}

		// Prepare environment
		env := NewEnvironment(nil, NewRequest())
		env.declareVariable("$"+paramName, arrayValue)

		// Execute PHP to store new array values in env
		interpreter := NewInterpreter(NewProdConfig(), NewRequest(), "")
		interpreter.process(fmt.Sprintf(`<?php $%s = "%s";`, key, value), env)

		// Extract array from environment
		arrayValue = env.variables["$"+paramName]

		result.SetElement(NewStringRuntimeValue(paramName), arrayValue)
		continue
	}
	return result
}

func registerPredefinedConstants(environment *Environment) {
	// Spec: https://phplang.org/spec/06-constants.html#core-predefined-constants
	environment.predefinedConstants["FALSE"] = NewBooleanRuntimeValue(false)
	environment.predefinedConstants["TRUE"] = NewBooleanRuntimeValue(true)
	environment.predefinedConstants["NULL"] = NewNullRuntimeValue()
	environment.predefinedConstants["PHP_OS"] = NewStringRuntimeValue(getPhpOs())
	environment.predefinedConstants["PHP_OS_FAMILY"] = NewStringRuntimeValue(getPhpOsFamily())
	if getPhpOs() == "Windows" {
		environment.predefinedConstants["PHP_EOL"] = NewStringRuntimeValue("\r\n")
	} else {
		environment.predefinedConstants["PHP_EOL"] = NewStringRuntimeValue("\n")
	}

	// Spec: https://www.php.net/manual/en/errorfunc.constants.php
	environment.predefinedConstants["E_ERROR"] = NewIntegerRuntimeValue(E_ERROR)
	environment.predefinedConstants["E_WARNING"] = NewIntegerRuntimeValue(E_WARNING)
	environment.predefinedConstants["E_PARSE"] = NewIntegerRuntimeValue(E_PARSE)
	environment.predefinedConstants["E_NOTICE"] = NewIntegerRuntimeValue(E_NOTICE)
	environment.predefinedConstants["E_CORE_ERROR"] = NewIntegerRuntimeValue(E_CORE_ERROR)
	environment.predefinedConstants["E_CORE_WARNING"] = NewIntegerRuntimeValue(E_CORE_WARNING)
	environment.predefinedConstants["E_COMPILE_ERROR"] = NewIntegerRuntimeValue(E_COMPILE_ERROR)
	environment.predefinedConstants["E_COMPILE_WARNING"] = NewIntegerRuntimeValue(E_COMPILE_WARNING)
	environment.predefinedConstants["E_USER_ERROR"] = NewIntegerRuntimeValue(E_USER_ERROR)
	environment.predefinedConstants["E_USER_WARNING"] = NewIntegerRuntimeValue(E_USER_WARNING)
	environment.predefinedConstants["E_USER_NOTICE"] = NewIntegerRuntimeValue(E_USER_NOTICE)
	environment.predefinedConstants["E_STRICT"] = NewIntegerRuntimeValue(E_STRICT)
	environment.predefinedConstants["E_RECOVERABLE_ERROR"] = NewIntegerRuntimeValue(E_RECOVERABLE_ERROR)
	environment.predefinedConstants["E_DEPRECATED"] = NewIntegerRuntimeValue(E_DEPRECATED)
	environment.predefinedConstants["E_USER_DEPRECATED"] = NewIntegerRuntimeValue(E_USER_DEPRECATED)
	environment.predefinedConstants["E_ALL"] = NewIntegerRuntimeValue(E_ALL)
}

// ------------------- MARK: Variables -------------------

func (env *Environment) declareVariable(variableName string, value IRuntimeValue) (IRuntimeValue, Error) {
	env.variables[variableName] = deepCopy(value)

	return value, nil
}

func (env *Environment) resolvePredefinedVariable(variableName string) (*Environment, Error) {
	if _, ok := env.predefinedVariables[variableName]; ok {
		return env, nil
	}

	if env.parent == nil {
		return nil, NewWarning("Undefined variable %s", variableName)
	}

	return env.parent.resolvePredefinedVariable(variableName)
}

func (env *Environment) resolveVariable(variableName string) (*Environment, Error) {
	environment, err := env.resolvePredefinedVariable(variableName)
	if err != nil {
		if _, ok := env.variables[variableName]; ok {
			return env, nil
		} else {
			return nil, err
		}
	}
	return environment, nil
}

func (env *Environment) lookupVariable(variableName string) (IRuntimeValue, Error) {
	environment, err := env.resolveVariable(variableName)
	if err != nil {
		return NewNullRuntimeValue(), err
	}
	if value, ok := environment.predefinedVariables[variableName]; ok {
		return value, nil
	}
	if value, ok := environment.variables[variableName]; ok {
		return value, nil
	}
	return NewNullRuntimeValue(), NewWarning("Undefined variable %s", variableName)
}

func (env *Environment) unsetVariable(variableName string) {
	environment, err := env.resolveVariable(variableName)
	if err != nil {
		return
	}
	delete(environment.variables, variableName)
}

// ------------------- MARK: Constants -------------------

func (env *Environment) declareConstant(constantName string, value IRuntimeValue) (IRuntimeValue, Error) {
	// Get "global" environment
	var environment *Environment = env
	for environment.parent != nil {
		environment = environment.parent
	}

	if _, err := environment.lookupConstant(constantName); err == nil {
		return NewVoidRuntimeValue(), NewWarning("Constant %s already defined", constantName)
	}

	environment.constants[constantName] = value

	return value, nil
}

func (env *Environment) lookupConstant(constantName string) (IRuntimeValue, Error) {
	// Get "global" environment
	var environment *Environment = env
	for environment.parent != nil {
		environment = environment.parent
	}

	if value, ok := environment.predefinedConstants[constantName]; ok {
		return value, nil
	}
	if value, ok := environment.constants[constantName]; ok {
		return value, nil
	}
	return NewVoidRuntimeValue(), NewError("Undefined constant \"%s\"", constantName)
}

// ------------------- MARK: Native functions -------------------

func (env *Environment) resolveNativeFunction(functionName string) (*Environment, Error) {
	if _, ok := env.nativeFunctions[functionName]; ok {
		return env, nil
	}

	if env.parent == nil {
		return nil, NewError("Call to undefined function %s()", functionName)
	}

	return env.parent.resolveNativeFunction(functionName)
}

func (env *Environment) lookupNativeFunction(functionName string) (nativeFunction, Error) {
	functionName = strings.ToLower(functionName)

	environment, err := env.resolveNativeFunction(functionName)
	if err != nil {
		return nil, err
	}

	value, ok := environment.nativeFunctions[functionName]
	if !ok {
		return nil, NewError("Call to undefined function %s()", functionName)
	}
	return value, nil
}

// ------------------- MARK: User functions -------------------

type userFunction struct {
	FunctionName string
	Parameters   []ast.FunctionParameter
	Body         ast.ICompoundStatement
	ReturnType   []string
}

func (env *Environment) defineUserFunction(function userFunction) Error {
	_, err := env.lookupNativeFunction(function.FunctionName)
	if err == nil {
		return NewError("Cannot redeclare %s()", function.FunctionName)
	}
	_, err = env.lookupUserFunction(function.FunctionName)
	if err == nil {
		return NewError("Cannot redeclare %s()", function.FunctionName)
	}

	functionName := strings.ToLower(function.FunctionName)

	env.functions[functionName] = function

	return nil
}

func (env *Environment) resolveUserFunction(functionName string) (*Environment, Error) {
	if _, ok := env.functions[functionName]; ok {
		return env, nil
	}

	if env.parent == nil {
		return nil, NewError("Call to undefined function %s()", functionName)
	}

	return env.parent.resolveUserFunction(functionName)
}

func (env *Environment) lookupUserFunction(functionName string) (userFunction, Error) {
	functionName = strings.ToLower(functionName)

	environment, err := env.resolveUserFunction(functionName)
	if err != nil {
		return userFunction{}, err
	}

	value, ok := environment.functions[functionName]
	if !ok {
		return userFunction{}, NewError("Call to undefined function %s()", functionName)
	}
	return value, nil
}
