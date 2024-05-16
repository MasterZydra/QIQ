package interpreter

import (
	"fmt"
	"strings"
)

type Environment struct {
	parent    *Environment
	variables map[string]IRuntimeValue
	constants map[string]IRuntimeValue
	// StdLib
	nativeFunctions map[string]nativeFunction
}

func NewEnvironment(parentEnv *Environment) *Environment {
	env := &Environment{
		parent:    parentEnv,
		variables: make(map[string]IRuntimeValue),
		constants: make(map[string]IRuntimeValue),
		// StdLib
		nativeFunctions: make(map[string]nativeFunction),
	}

	if parentEnv == nil {
		registerNativeFunctions(env)
	}

	return env
}

// ------------------- MARK: Variables -------------------

func (env *Environment) declareVariable(variableName string, value IRuntimeValue) (IRuntimeValue, error) {
	env.variables[variableName] = value

	return value, nil
}

func (env *Environment) resolveVariable(variableName string) (*Environment, error) {
	if _, ok := env.variables[variableName]; ok {
		return env, nil
	}

	if env.parent == nil {
		return nil, fmt.Errorf("Interpreter error: Cannot resolve variable '%s' as it does not exist", variableName)
	}

	return env.parent.resolveVariable(variableName)
}

func (env *Environment) lookupVariable(variableName string) (IRuntimeValue, error) {
	environment, err := env.resolveVariable(variableName)
	if err != nil {
		return NewNullRuntimeValue(), fmt.Errorf("Warning: Undefined variable %s", variableName)
	}
	value, ok := environment.variables[variableName]
	if !ok {
		return NewNullRuntimeValue(), fmt.Errorf("Warning: Undefined variable %s", variableName)
	}
	return value, nil
}

// ------------------- MARK: Constants -------------------

func (env *Environment) declareConstant(constantName string, value IRuntimeValue) (IRuntimeValue, error) {
	if _, err := env.lookupConstant(constantName); err == nil {
		return NewVoidRuntimeValue(), fmt.Errorf("Cannot redefine an exisiting constant: \"%s\"", constantName)
	}

	env.constants[constantName] = value

	return value, nil
}

func (env *Environment) resolveConstant(constantName string) (*Environment, error) {
	if _, ok := env.constants[constantName]; ok {
		return env, nil
	}

	if env.parent == nil {
		return nil, fmt.Errorf("Interpreter error: Cannot resolve constant \"%s\" as it does not exist", constantName)
	}

	return env.parent.resolveConstant(constantName)
}

func (env *Environment) lookupConstant(constantName string) (IRuntimeValue, error) {
	environment, err := env.resolveConstant(constantName)
	if err != nil {
		return NewVoidRuntimeValue(), err
	}
	value, ok := environment.constants[constantName]
	if !ok {
		return NewVoidRuntimeValue(), fmt.Errorf("Interpreter error: Undefined constant \"%s\"", constantName)
	}
	return value, nil
}

// ------------------- MARK: Native functions -------------------

func (env *Environment) resolveNativeFunction(functionName string) (*Environment, error) {
	if _, ok := env.nativeFunctions[functionName]; ok {
		return env, nil
	}

	if env.parent == nil {
		return nil, fmt.Errorf("Interpreter error: Cannot resolve native function \"%s\" as it does not exist", functionName)
	}

	return env.parent.resolveNativeFunction(functionName)
}

func (env *Environment) lookupNativeFunction(functionName string) (nativeFunction, error) {
	functionName = strings.ToLower(functionName)

	environment, err := env.resolveNativeFunction(functionName)
	if err != nil {
		return nil, err
	}

	value, ok := environment.nativeFunctions[functionName]
	if !ok {
		return nil, fmt.Errorf("Cannot call undefined function %s()", functionName)
	}
	return value, nil
}
