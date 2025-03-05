package interpreter

import (
	"GoPHP/cmd/goPHP/ini"
	"GoPHP/cmd/goPHP/phpError"
	"GoPHP/cmd/goPHP/runtime"
	"GoPHP/cmd/goPHP/runtime/values"
	"slices"
	"strconv"
)

func registerNativeFunctions(environment *Environment) {
	registerNativeDateTimeFunctions(environment)
	registerNativeMathFunctions(environment)
	registerNativeMiscFunctions(environment)
	registerNativeOptionsInfoFunctions(environment)
	registerNativeOutputControlFunctions(environment)
	registerNativeStringsFunctions(environment)
	registerNativeVariableHandlingFunctions(environment)

	environment.nativeFunctions["array_key_exits"] = nativeFn_array_key_exists
	environment.nativeFunctions["error_reporting"] = nativeFn_error_reporting
	environment.nativeFunctions["key_exits"] = nativeFn_array_key_exists
}

type nativeFunction func([]values.RuntimeValue, runtime.Context) (values.RuntimeValue, phpError.Error)

// ------------------- MARK: array_key_exits -------------------

func nativeFn_array_key_exists(args []values.RuntimeValue, _ runtime.Context) (values.RuntimeValue, phpError.Error) {
	args, err := NewFuncParamValidator("array_key_exists").
		addParam("$key", []string{"string", "int", "float", "bool", "resource", "null"}, nil).
		addParam("$array", []string{"array"}, nil).
		validate(args)
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

// ------------------- MARK: arrayval -------------------

// This is not an official function. But converting different types to array is needed in several places
func lib_arrayval(runtimeValue values.RuntimeValue) (*values.Array, phpError.Error) {
	// Spec: https://phplang.org/spec/08-conversions.html#converting-to-array-type

	// The result type is array.

	if runtimeValue.GetType() == values.NullValue {
		// Spec: https://phplang.org/spec/08-conversions.html#converting-to-array-type
		// If the source value is NULL, the result value is an array of zero elements.
		return values.NewArray(), nil
	}

	// TODO lib_arrayval - resource
	if lib_is_scalar(runtimeValue) {
		// Spec: https://phplang.org/spec/08-conversions.html#converting-to-array-type
		// If the source type is scalar or resource and it is non-NULL,
		// the result value is an array of one element under the key 0 whose value is that of the source.
		array := values.NewArray()
		array.SetElement(nil, runtimeValue)
		return array, nil
	}

	// TODO lib_arrayval - object
	// Spec: https://phplang.org/spec/08-conversions.html#converting-to-array-type
	// If the source is an object, the result is an array of zero or more elements, where the elements are key/value pairs corresponding to the object’s instance properties. The order of insertion of the elements into the array is the lexical order of the instance properties in the class-member-declarations list.

	// TODO lib_arrayval - instance properties
	// Spec: https://phplang.org/spec/08-conversions.html#converting-to-array-type
	// For public instance properties, the keys of the array elements would be the same as the property name.
	// The key for a private instance property has the form “\0class\0name”, where the class is the class name, and the name is the property name.
	// The key for a protected instance property has the form “\0*\0name”, where name is that of the property.
	// The value for each key is that from the corresponding property, or NULL if the property was not initialized.

	return values.NewArray(), phpError.NewError("lib_arrayval: Unsupported type %s", runtimeValue.GetType())
}

// ------------------- MARK: error_reporting -------------------

func nativeFn_error_reporting(args []values.RuntimeValue, context runtime.Context) (values.RuntimeValue, phpError.Error) {
	// Spec: https://www.php.net/manual/en/function.error-reporting.php

	args, err := NewFuncParamValidator("error_reporting").addParam("$error_level", []string{"int"}, values.NewNull()).validate(args)
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
