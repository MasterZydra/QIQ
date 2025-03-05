package stdlib

import (
	"GoPHP/cmd/goPHP/ini"
	"GoPHP/cmd/goPHP/phpError"
	"GoPHP/cmd/goPHP/runtime"
	"GoPHP/cmd/goPHP/runtime/funcParamValidator"
	"GoPHP/cmd/goPHP/runtime/stdlib/dateTime"
	"GoPHP/cmd/goPHP/runtime/stdlib/math"
	"GoPHP/cmd/goPHP/runtime/stdlib/misc"
	"GoPHP/cmd/goPHP/runtime/stdlib/optionsInfo"
	"GoPHP/cmd/goPHP/runtime/stdlib/outputControl"
	"GoPHP/cmd/goPHP/runtime/stdlib/strings"
	"GoPHP/cmd/goPHP/runtime/stdlib/variableHandling"
	"GoPHP/cmd/goPHP/runtime/values"
	"slices"
	"strconv"
)

func RegisterNativeFunctions(environment runtime.Environment) {
	dateTime.Register(environment)
	math.Register(environment)
	misc.Register(environment)
	optionsInfo.Register(environment)
	outputControl.Register(environment)
	strings.Register(environment)
	variableHandling.Register(environment)

	environment.AddNativeFunction("array_key_exits", nativeFn_array_key_exists)
	environment.AddNativeFunction("error_reporting", nativeFn_error_reporting)
	environment.AddNativeFunction("key_exits", nativeFn_array_key_exists)
}

// ------------------- MARK: array_key_exits -------------------

func nativeFn_array_key_exists(args []values.RuntimeValue, _ runtime.Context) (values.RuntimeValue, phpError.Error) {
	args, err := funcParamValidator.NewValidator("array_key_exists").
		AddParam("$key", []string{"string", "int", "float", "bool", "resource", "null"}, nil).
		AddParam("$array", []string{"array"}, nil).
		Validate(args)
	if err != nil {
		return values.NewVoid(), err
	}

	boolean, err := lib_array_key_exists(args[0], args[1].(*values.Array))
	return values.NewBool(boolean), err
}

func lib_array_key_exists(key values.RuntimeValue, array *values.Array) (bool, phpError.Error) {
	// Spec: https://www.php.net/manual/en/function.array-key-exists.php

	// TODO lib_array_key_exists - allowedKeyTypes - resource
	allowedKeyTypes := []values.ValueType{values.StrValue, values.IntValue, values.FloatValue, values.BoolValue, values.NullValue}

	if !slices.Contains(allowedKeyTypes, key.GetType()) {
		return false, phpError.NewError("Values of type %s are not allowed as array key", key.GetType())
	}

	_, ok := array.GetElement(key)
	return ok, nil
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
