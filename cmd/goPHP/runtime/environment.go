package runtime

import (
	"GoPHP/cmd/goPHP/phpError"
	"GoPHP/cmd/goPHP/runtime/values"
)

type Environment interface {
	// Variables
	LookupVariable(variableName string) (values.RuntimeValue, phpError.Error)
	// Functions
	AddNativeFunction(functionName string, function NativeFunction)
	FunctionExists(functionName string) bool
	// Constants
	LookupConstant(constantName string) (values.RuntimeValue, phpError.Error)
	AddPredefinedConstant(name string, value values.RuntimeValue)
	AddConstant(name string, value values.RuntimeValue)
}
