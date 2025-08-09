package interpreter

import (
	"GoPHP/cmd/goPHP/ast"
	"GoPHP/cmd/goPHP/phpError"
	"GoPHP/cmd/goPHP/request"
	"GoPHP/cmd/goPHP/runtime"
	"GoPHP/cmd/goPHP/runtime/stdlib"
	"GoPHP/cmd/goPHP/runtime/values"
	"slices"
	"strings"
)

type Environment struct {
	parent          *Environment
	globalVariables []string
	variables       map[string]values.RuntimeValue
	constants       map[string]values.RuntimeValue
	functions       map[string]*ast.FunctionDefinitionStatement
	// StdLib
	predefinedVariables map[string]values.RuntimeValue
	predefinedConstants map[string]values.RuntimeValue
	nativeFunctions     map[string]runtime.NativeFunction
	// Context
	CurrentFunction *ast.FunctionDefinitionStatement
	CurrentObject   *values.Object
	CurrentMethod   *ast.MethodDefinitionStatement
}

func NewEnvironment(parentEnv *Environment, request *request.Request, interpreter runtime.Interpreter) (*Environment, phpError.Error) {
	env := &Environment{
		parent:    parentEnv,
		variables: map[string]values.RuntimeValue{},
		constants: map[string]values.RuntimeValue{},
		functions: map[string]*ast.FunctionDefinitionStatement{},
		// StdLib
		predefinedVariables: map[string]values.RuntimeValue{},
		predefinedConstants: map[string]values.RuntimeValue{},
		nativeFunctions:     map[string]runtime.NativeFunction{},
	}

	if parentEnv == nil {
		stdlib.Register(env)
		if err := registerPredefinedVariables(env, request, interpreter); err != nil {
			return env, err
		}
		registerPredefinedConstants(env)
	}

	return env, nil
}

// -------------------------------------- Variables -------------------------------------- MARK: Variables

func (env *Environment) declareVariable(variableName string, value values.RuntimeValue) (values.RuntimeValue, phpError.Error) {
	if slices.Contains(env.globalVariables, variableName) {
		return env.parent.declareVariable(variableName, value)
	}

	env.variables[variableName] = values.DeepCopy(value)
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
	if err == nil {
		return environment, nil
	}

	if slices.Contains(env.globalVariables, variableName) {
		rootEnv := env
		for rootEnv.parent != nil {
			rootEnv = rootEnv.parent
		}
		return rootEnv, nil
	}

	if _, ok := env.variables[variableName]; ok {
		return env, nil
	}

	return nil, err
}

func (env *Environment) LookupVariable(variableName string) (values.RuntimeValue, phpError.Error) {
	environment, err := env.resolveVariable(variableName)
	if err != nil {
		return values.NewNull(), err
	}
	if value, ok := environment.predefinedVariables[variableName]; ok {
		return value, nil
	}
	if slices.Contains(env.globalVariables, variableName) {
		value, ok := environment.variables[variableName]
		if ok {
			return value, nil
		}
		return values.NewNull(), nil
	}
	if value, ok := environment.variables[variableName]; ok {
		return value, nil
	}
	return values.NewNull(), phpError.NewWarning("Undefined variable %s", variableName)
}

func (env *Environment) unsetVariable(variableName string) {
	environment, err := env.resolveVariable(variableName)
	if err != nil {
		return
	}
	delete(environment.variables, variableName)
}

func (env *Environment) addGlobalVariable(variableName string) {
	if env.parent == nil {
		return
	}
	if slices.Contains(env.globalVariables, variableName) {
		return
	}
	env.globalVariables = append(env.globalVariables, variableName)
}

func (env *Environment) getAllObjects() []*values.Object {
	objects := []*values.Object{}
	for _, variable := range env.variables {
		if variable.GetType() != values.ObjectValue {
			continue
		}
		objects = append(objects, variable.(*values.Object))
	}
	return objects
}

// -------------------------------------- Constants -------------------------------------- MARK: Constants

func (env *Environment) AddConstant(name string, value values.RuntimeValue) {
	env.constants[name] = value
}

func (env *Environment) AddPredefinedConstant(name string, value values.RuntimeValue) {
	env.predefinedConstants[name] = value
}

func (env *Environment) declareConstant(constantName string, value values.RuntimeValue) (values.RuntimeValue, phpError.Error) {
	// Get "global" environment
	var environment *Environment = env
	for environment.parent != nil {
		environment = environment.parent
	}

	if _, err := environment.LookupConstant(constantName); err == nil {
		return values.NewVoid(), phpError.NewWarning("Constant %s already defined", constantName)
	}

	environment.constants[constantName] = value

	return value, nil
}

func (env *Environment) LookupConstant(constantName string) (values.RuntimeValue, phpError.Error) {
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
	return values.NewVoid(), phpError.NewError("Undefined constant \"%s\"", constantName)
}

// -------------------------------------- Functions -------------------------------------- MARK: Functions

func (env *Environment) FunctionExists(functionName string) bool {
	if _, err := env.resolveNativeFunction(functionName); err == nil {
		return true
	}

	_, err := env.resolveUserFunction(strings.ToLower(functionName))
	return err == nil
}

// -------------------------------------- Native functions -------------------------------------- MARK: Native functions

func (env *Environment) AddNativeFunction(functionName string, function runtime.NativeFunction) {
	env.nativeFunctions[functionName] = function
}

func (env *Environment) resolveNativeFunction(functionName string) (*Environment, phpError.Error) {
	if _, ok := env.nativeFunctions[functionName]; ok {
		return env, nil
	}

	if env.parent == nil {
		return nil, phpError.NewError("Call to undefined function %s()", functionName)
	}

	return env.parent.resolveNativeFunction(functionName)
}

func (env *Environment) lookupNativeFunction(functionName string) (runtime.NativeFunction, phpError.Error) {
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

// -------------------------------------- User functions -------------------------------------- MARK: User functions

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
