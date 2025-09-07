package interpreter

import (
	"QIQ/cmd/qiq/ast"
	"QIQ/cmd/qiq/phpError"
	"QIQ/cmd/qiq/request"
	"QIQ/cmd/qiq/runtime"
	"QIQ/cmd/qiq/runtime/stdlib"
	"QIQ/cmd/qiq/runtime/values"
	"slices"
	"strings"
)

type Environment struct {
	parent          *Environment
	globalVariables []string
	variables       map[string]*values.Slot
	constants       map[string]*values.Slot
	functions       map[string]*ast.FunctionDefinitionStatement
	objects         []*values.Object
	// StdLib
	predefinedVariables map[string]*values.Slot
	predefinedConstants map[string]*values.Slot
	nativeFunctions     map[string]runtime.NativeFunction
	// Context
	CurrentFunction *ast.FunctionDefinitionStatement
	CurrentObject   *values.Object
	CurrentMethod   *ast.MethodDefinitionStatement
}

func NewEnvironment(parentEnv *Environment, request *request.Request, interpreter runtime.Interpreter) (*Environment, phpError.Error) {
	env := &Environment{
		parent:    parentEnv,
		variables: map[string]*values.Slot{},
		constants: map[string]*values.Slot{},
		functions: map[string]*ast.FunctionDefinitionStatement{},
		objects:   []*values.Object{},
		// StdLib
		predefinedVariables: map[string]*values.Slot{},
		predefinedConstants: map[string]*values.Slot{},
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

	if _, found := env.variables[variableName]; found {
		env.variables[variableName].Value = value
	} else {
		env.variables[variableName] = values.NewSlot(values.DeepCopy(value))
	}
	return value, nil
}

func (env *Environment) declareVariableByRef(variableName string, slot *values.Slot) (values.RuntimeValue, phpError.Error) {
	if slices.Contains(env.globalVariables, variableName) {
		return env.parent.declareVariableByRef(variableName, slot)
	}

	env.variables[variableName] = slot
	return slot.Value, nil
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
	if slot, ok := environment.predefinedVariables[variableName]; ok {
		return slot.Value, nil
	}
	if slices.Contains(env.globalVariables, variableName) {
		slot, ok := environment.variables[variableName]
		if ok {
			return slot.Value, nil
		}
		return values.NewNull(), nil
	}
	if slot, ok := environment.variables[variableName]; ok {
		return slot.Value, nil
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
		if variable.Value.GetType() != values.ObjectValue {
			continue
		}
		objects = append(objects, variable.Value.(*values.Object))
	}
	return objects
}

// -------------------------------------- Objects -------------------------------------- MARK: Objects

func (env *Environment) AddObject(object *values.Object) {
	env.objects = append(env.objects, object)
}

// -------------------------------------- Constants -------------------------------------- MARK: Constants

func (env *Environment) AddConstant(name string, value values.RuntimeValue) {
	env.constants[name] = values.NewSlot(value)
}

func (env *Environment) AddPredefinedConstant(name string, value values.RuntimeValue) {
	env.predefinedConstants[name] = values.NewSlot(value)
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

	environment.constants[constantName] = values.NewSlot(value)

	return value, nil
}

func (env *Environment) LookupConstant(constantName string) (values.RuntimeValue, phpError.Error) {
	// Get "global" environment
	var environment *Environment = env
	for environment.parent != nil {
		environment = environment.parent
	}

	if slot, ok := environment.predefinedConstants[constantName]; ok {
		return slot.Value, nil
	}
	if slot, ok := environment.constants[constantName]; ok {
		return slot.Value, nil
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
