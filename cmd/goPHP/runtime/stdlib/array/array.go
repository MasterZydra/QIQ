package array

import (
	"GoPHP/cmd/goPHP/phpError"
	"GoPHP/cmd/goPHP/runtime"
	"GoPHP/cmd/goPHP/runtime/funcParamValidator"
	"GoPHP/cmd/goPHP/runtime/values"
	"slices"
)

func Register(environment runtime.Environment) {
	// Category: Array Functions
	environment.AddNativeFunction("array_key_exists", nativeFn_array_key_exists)
	environment.AddNativeFunction("key_exits", nativeFn_array_key_exists)
}

// ------------------- MARK: array_key_exists -------------------

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

	return array.Contains(key), nil
}
