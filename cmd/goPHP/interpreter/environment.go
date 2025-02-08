package interpreter

import (
	"GoPHP/cmd/goPHP/ast"
	"GoPHP/cmd/goPHP/ini"
	"GoPHP/cmd/goPHP/phpError"
	"strings"
)

type Environment struct {
	parent    *Environment
	variables map[string]IRuntimeValue
	constants map[string]IRuntimeValue
	functions map[string]*ast.FunctionDefinitionStatement
	// StdLib
	predefinedVariables map[string]IRuntimeValue
	predefinedConstants map[string]IRuntimeValue
	nativeFunctions     map[string]nativeFunction
	// Context
	CurrentFunction *ast.FunctionDefinitionStatement
}

func NewEnvironment(parentEnv *Environment, request *Request, ini *ini.Ini) *Environment {
	env := &Environment{
		parent:    parentEnv,
		variables: map[string]IRuntimeValue{},
		constants: map[string]IRuntimeValue{},
		functions: map[string]*ast.FunctionDefinitionStatement{},
		// StdLib
		predefinedVariables: map[string]IRuntimeValue{},
		predefinedConstants: map[string]IRuntimeValue{},
		nativeFunctions:     map[string]nativeFunction{},
	}

	if parentEnv == nil {
		registerNativeFunctions(env)
		registerPredefinedVariables(env, request, ini)
		registerPredefinedConstants(env)
	}

	return env
}

// ------------------- MARK: Variables -------------------

func (env *Environment) declareVariable(variableName string, value IRuntimeValue) (IRuntimeValue, phpError.Error) {
	env.variables[variableName] = deepCopy(value)

	return value, nil
}

func (env *Environment) resolvePredefinedVariable(variableName string) (*Environment, phpError.Error) {
	if _, ok := env.predefinedVariables[variableName]; ok {
		return env, nil
	}

	if env.parent == nil {
		return nil, phpError.NewWarning("Undefined variable %s", variableName)
	}

	return env.parent.resolvePredefinedVariable(variableName)
}

func (env *Environment) resolveVariable(variableName string) (*Environment, phpError.Error) {
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

func (env *Environment) lookupVariable(variableName string) (IRuntimeValue, phpError.Error) {
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
	return NewNullRuntimeValue(), phpError.NewWarning("Undefined variable %s", variableName)
}

func (env *Environment) unsetVariable(variableName string) {
	environment, err := env.resolveVariable(variableName)
	if err != nil {
		return
	}
	delete(environment.variables, variableName)
}

// ------------------- MARK: Constants -------------------

func (env *Environment) declareConstant(constantName string, value IRuntimeValue) (IRuntimeValue, phpError.Error) {
	// Get "global" environment
	var environment *Environment = env
	for environment.parent != nil {
		environment = environment.parent
	}

	if _, err := environment.lookupConstant(constantName); err == nil {
		return NewVoidRuntimeValue(), phpError.NewWarning("Constant %s already defined", constantName)
	}

	environment.constants[constantName] = value

	return value, nil
}

func (env *Environment) lookupConstant(constantName string) (IRuntimeValue, phpError.Error) {
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
	return NewVoidRuntimeValue(), phpError.NewError("Undefined constant \"%s\"", constantName)
}

// ------------------- MARK: Native functions -------------------

func (env *Environment) resolveNativeFunction(functionName string) (*Environment, phpError.Error) {
	if _, ok := env.nativeFunctions[functionName]; ok {
		return env, nil
	}

	if env.parent == nil {
		return nil, phpError.NewError("Call to undefined function %s()", functionName)
	}

	return env.parent.resolveNativeFunction(functionName)
}

func (env *Environment) lookupNativeFunction(functionName string) (nativeFunction, phpError.Error) {
	functionName = strings.ToLower(functionName)

	environment, err := env.resolveNativeFunction(functionName)
	if err != nil {
		return nil, err
	}

	value, ok := environment.nativeFunctions[functionName]
	if !ok {
		return nil, phpError.NewError("Call to undefined function %s()", functionName)
	}
	return value, nil
}

// ------------------- MARK: User functions -------------------

func (env *Environment) defineUserFunction(function *ast.FunctionDefinitionStatement) phpError.Error {
	_, err := env.lookupNativeFunction(function.FunctionName)
	if err == nil {
		return phpError.NewError("Cannot redeclare %s()", function.FunctionName)
	}
	_, err = env.lookupUserFunction(function.FunctionName)
	if err == nil {
		return phpError.NewError("Cannot redeclare %s()", function.FunctionName)
	}

	functionName := strings.ToLower(function.FunctionName)

	env.functions[functionName] = function

	return nil
}

func (env *Environment) resolveUserFunction(functionName string) (*Environment, phpError.Error) {
	if _, ok := env.functions[functionName]; ok {
		return env, nil
	}

	if env.parent == nil {
		return nil, phpError.NewError("Call to undefined function %s()", functionName)
	}

	return env.parent.resolveUserFunction(functionName)
}

func (env *Environment) lookupUserFunction(functionName string) (*ast.FunctionDefinitionStatement, phpError.Error) {
	functionName = strings.ToLower(functionName)

	environment, err := env.resolveUserFunction(functionName)
	if err != nil {
		return &ast.FunctionDefinitionStatement{}, err
	}

	value, ok := environment.functions[functionName]
	if !ok {
		return &ast.FunctionDefinitionStatement{}, phpError.NewError("Call to undefined function %s()", functionName)
	}
	return value, nil
}
