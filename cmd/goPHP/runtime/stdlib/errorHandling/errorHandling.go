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

	// Const Category: Error Handling Constants
	// Spec: https://www.php.net/manual/en/errorfunc.constants.php
	environment.AddPredefinedConstants("E_ERROR", values.NewInt(phpError.E_ERROR))
	environment.AddPredefinedConstants("E_WARNING", values.NewInt(phpError.E_WARNING))
	environment.AddPredefinedConstants("E_PARSE", values.NewInt(phpError.E_PARSE))
	environment.AddPredefinedConstants("E_NOTICE", values.NewInt(phpError.E_NOTICE))
	environment.AddPredefinedConstants("E_CORE_ERROR", values.NewInt(phpError.E_CORE_ERROR))
	environment.AddPredefinedConstants("E_CORE_WARNING", values.NewInt(phpError.E_CORE_WARNING))
	environment.AddPredefinedConstants("E_COMPILE_ERROR", values.NewInt(phpError.E_COMPILE_ERROR))
	environment.AddPredefinedConstants("E_COMPILE_WARNING", values.NewInt(phpError.E_COMPILE_WARNING))
	environment.AddPredefinedConstants("E_USER_ERROR", values.NewInt(phpError.E_USER_ERROR))
	environment.AddPredefinedConstants("E_USER_WARNING", values.NewInt(phpError.E_USER_WARNING))
	environment.AddPredefinedConstants("E_USER_NOTICE", values.NewInt(phpError.E_USER_NOTICE))
	environment.AddPredefinedConstants("E_STRICT", values.NewInt(phpError.E_STRICT))
	environment.AddPredefinedConstants("E_RECOVERABLE_ERROR", values.NewInt(phpError.E_RECOVERABLE_ERROR))
	environment.AddPredefinedConstants("E_DEPRECATED", values.NewInt(phpError.E_DEPRECATED))
	environment.AddPredefinedConstants("E_USER_DEPRECATED", values.NewInt(phpError.E_USER_DEPRECATED))
	environment.AddPredefinedConstants("E_ALL", values.NewInt(phpError.E_ALL))
}

// -------------------------------------- error_reporting -------------------------------------- MARK: error_reporting

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
