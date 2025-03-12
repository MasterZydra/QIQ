package errorHandling

import (
	"GoPHP/cmd/goPHP/ini"
	"GoPHP/cmd/goPHP/phpError"
	"GoPHP/cmd/goPHP/runtime"
	"GoPHP/cmd/goPHP/runtime/funcParamValidator"
	"GoPHP/cmd/goPHP/runtime/values"
	"strconv"
)

func Register(environment runtime.Environment) {
	// Category: Error Handling Functions
	environment.AddNativeFunction("error_reporting", nativeFn_error_reporting)
}

// ------------------- MARK: error_reporting -------------------

func nativeFn_error_reporting(args []values.RuntimeValue, context runtime.Context) (values.RuntimeValue, phpError.Error) {
	// Spec: https://www.php.net/manual/en/function.error-reporting.php

	args, err := funcParamValidator.NewValidator("error_reporting").AddParam("$error_level", []string{"int"}, values.NewNull()).Validate(args)
	if err != nil {
		return values.NewVoid(), err
	}

	if args[0].GetType() == values.NullValue {
		return values.NewInt(context.Interpreter.GetIni().GetInt("error_reporting")), nil
	}

	newValue := args[0].(*values.Int).Value
	if newValue == -1 {
		newValue = phpError.E_ALL
	}

	previous := context.Interpreter.GetIni().GetInt("error_reporting")
	context.Interpreter.GetIni().Set("error_reporting", strconv.FormatInt(newValue, 10), ini.INI_USER)

	return values.NewInt(previous), nil
}
