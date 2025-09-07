package runtime

import (
	"QIQ/cmd/qiq/phpError"
	"QIQ/cmd/qiq/runtime/values"
)

type Environment interface {
	// Variables
	LookupVariable(variableName string) (*values.Slot, phpError.Error)
	// Functions
	AddNativeFunction(functionName string, function NativeFunction)
	FunctionExists(functionName string) bool
	// Constants
	LookupConstant(constantName string) (values.RuntimeValue, phpError.Error)
	AddPredefinedConstant(name string, value values.RuntimeValue)
	AddConstant(name string, value values.RuntimeValue)
}
