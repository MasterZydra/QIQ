package interpreter

import (
	"strings"
)

type Environment struct {
	parent    *Environment
	variables map[string]IRuntimeValue
	constants map[string]IRuntimeValue
	// StdLib
	predefinedVariables map[string]IRuntimeValue
	predefinedConstants map[string]IRuntimeValue
	nativeFunctions     map[string]nativeFunction
}

func NewEnvironment(parentEnv *Environment, request *Request) *Environment {
	env := &Environment{
		parent:    parentEnv,
		variables: make(map[string]IRuntimeValue),
		constants: make(map[string]IRuntimeValue),
		// StdLib
		predefinedVariables: make(map[string]IRuntimeValue),
		predefinedConstants: make(map[string]IRuntimeValue),
		nativeFunctions:     make(map[string]nativeFunction),
	}

	if parentEnv == nil {
		registerNativeFunctions(env)
		registerPredefinedVariables(env, request)
		registerPredefinedConstants(env)
	}

	return env
}

func registerPredefinedVariables(environment *Environment, request *Request) {
	environment.predefinedVariables["$_GET"] = NewArrayRuntimeValueFromMap(request.GetParams)
	environment.predefinedVariables["$_POST"] = NewArrayRuntimeValueFromMap(request.PostParams)
}

func registerPredefinedConstants(environment *Environment) {
	// Spec: https://phplang.org/spec/06-constants.html#core-predefined-constants
	environment.predefinedConstants["FALSE"] = NewBooleanRuntimeValue(false)
	environment.predefinedConstants["TRUE"] = NewBooleanRuntimeValue(true)
	environment.predefinedConstants["NULL"] = NewNullRuntimeValue()

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
	env.variables[variableName] = value

	return value, nil
}

func (env *Environment) resolveVariable(variableName string) (*Environment, Error) {
	if _, ok := env.predefinedVariables[variableName]; ok {
		return env, nil
	}
	if _, ok := env.variables[variableName]; ok {
		return env, nil
	}

	if env.parent == nil {
		return nil, NewError("Interpreter error: Cannot resolve variable '%s' as it does not exist", variableName)
	}

	return env.parent.resolveVariable(variableName)
}

func (env *Environment) lookupVariable(variableName string) (IRuntimeValue, Error) {
	environment, err := env.resolveVariable(variableName)
	if err != nil {
		return NewNullRuntimeValue(), NewWarning("Undefined variable %s", variableName)
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
		environment = env.parent
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
		environment = env.parent
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
