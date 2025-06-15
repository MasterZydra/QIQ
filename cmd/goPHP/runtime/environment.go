package runtime

import (
	"GoPHP/cmd/goPHP/phpError"
	"GoPHP/cmd/goPHP/runtime/values"
)

type Environment interface {
	LookupConstant(constantName string) (values.RuntimeValue, phpError.Error)
	LookupVariable(variableName string) (values.RuntimeValue, phpError.Error)
	AddNativeFunction(functionName string, function NativeFunction)
	AddPredefinedConstants(name string, value values.RuntimeValue)
	AddConstants(name string, value values.RuntimeValue)
	FunctionExists(functionName string) bool
}
