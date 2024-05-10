package interpreter

import "fmt"

type Environment struct {
	parent    *Environment
	variables map[string]IRuntimeValue
	constants map[string]IRuntimeValue
}

func NewEnvironment(parentEnv *Environment) *Environment {
	return &Environment{
		parent:    parentEnv,
		variables: make(map[string]IRuntimeValue),
		constants: make(map[string]IRuntimeValue),
	}
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
	_, err := env.resolveVariable(variableName)
	if err != nil {
		return NewVoidRuntimeValue(), err
	}
	value, ok := env.variables[variableName]
	if !ok {
		return NewVoidRuntimeValue(), fmt.Errorf("Interpreter error: Undeclared variable \"%s\"", variableName)
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
	_, err := env.resolveConstant(constantName)
	if err != nil {
		return NewVoidRuntimeValue(), err
	}
	value, ok := env.constants[constantName]
	if !ok {
		return NewVoidRuntimeValue(), fmt.Errorf("Interpreter error: Undefined constant \"%s\"", constantName)
	}
	return value, nil
}
