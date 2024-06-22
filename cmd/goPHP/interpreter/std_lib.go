package interpreter

import (
	"GoPHP/cmd/goPHP/phpError"
	"fmt"
	"slices"
)

func registerNativeFunctions(environment *Environment) {
	registerNativeVariableHandlingFunctions(environment)

	environment.nativeFunctions["array_key_exits"] = nativeFn_array_key_exists
	environment.nativeFunctions["error_reporting"] = nativeFn_error_reporting
	environment.nativeFunctions["getenv"] = nativeFn_getenv
	environment.nativeFunctions["ini_get"] = nativeFn_ini_get
	environment.nativeFunctions["key_exits"] = nativeFn_array_key_exists
	environment.nativeFunctions["strlen"] = nativeFn_strlen
}

type nativeFunction func([]IRuntimeValue, *Interpreter) (IRuntimeValue, phpError.Error)

// ------------------- MARK: array_key_exits -------------------

func nativeFn_array_key_exists(args []IRuntimeValue, _ *Interpreter) (IRuntimeValue, phpError.Error) {
	args, err := NewFuncParamValidator("array_key_exists").
		addParam("$key", []string{"string", "int", "float", "bool", "resource", "null"}, nil).
		addParam("$array", []string{"array"}, nil).
		validate(args)
	if err != nil {
		return NewVoidRuntimeValue(), err
	}

	boolean, err := lib_array_key_exists(args[0], runtimeValToArrayRuntimeVal(args[1]))
	return NewBooleanRuntimeValue(boolean), err
}

func lib_array_key_exists(key IRuntimeValue, array IArrayRuntimeValue) (bool, phpError.Error) {
	// Spec: https://www.php.net/manual/en/function.array-key-exists.php

	// TODO lib_array_key_exists - allowedKeyTypes - resource
	allowedKeyTypes := []ValueType{StringValue, IntegerValue, FloatingValue, BooleanValue, NullValue}

	if !slices.Contains(allowedKeyTypes, key.GetType()) {
		return false, phpError.NewError("Values of type %s are not allowed as array key", key.GetType())
	}

	_, ok := array.GetElement(key)
	return ok, nil
}

// ------------------- MARK: arrayval -------------------

// This is not an official function. But converting different types to array is needed in several places
func lib_arrayval(runtimeValue IRuntimeValue) (IArrayRuntimeValue, phpError.Error) {
	// Spec: https://phplang.org/spec/08-conversions.html#converting-to-array-type

	// The result type is array.

	if runtimeValue.GetType() == NullValue {
		// Spec: https://phplang.org/spec/08-conversions.html#converting-to-array-type
		// If the source value is NULL, the result value is an array of zero elements.
		return NewArrayRuntimeValue(), nil
	}

	// TODO lib_arrayval - resource
	if lib_is_scalar(runtimeValue) {
		// Spec: https://phplang.org/spec/08-conversions.html#converting-to-array-type
		// If the source type is scalar or resource and it is non-NULL,
		// the result value is an array of one element under the key 0 whose value is that of the source.
		array := NewArrayRuntimeValue()
		array.SetElement(NewIntegerRuntimeValue(0), runtimeValue)
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

	return NewArrayRuntimeValue(), phpError.NewError("lib_arrayval: Unsupported type %s", runtimeValue.GetType())
}

// ------------------- MARK: error_reporting -------------------

func nativeFn_error_reporting(args []IRuntimeValue, interpreter *Interpreter) (IRuntimeValue, phpError.Error) {
	// Spec: https://www.php.net/manual/en/function.error-reporting.php

	args, err := NewFuncParamValidator("error_reporting").addParam("$error_level", []string{"int"}, NewNullRuntimeValue()).validate(args)
	if err != nil {
		return NewVoidRuntimeValue(), err
	}

	if args[0].GetType() == NullValue {
		return NewIntegerRuntimeValue(interpreter.ini.ErrorReporting), nil
	}

	newValue := runtimeValToIntRuntimeVal(args[0]).GetValue()
	if newValue == -1 {
		newValue = phpError.E_ALL
	}

	previous := interpreter.ini.ErrorReporting
	interpreter.ini.ErrorReporting = newValue

	return NewIntegerRuntimeValue(previous), nil
}

// ------------------- MARK: getenv -------------------

func nativeFn_getenv(args []IRuntimeValue, interpreter *Interpreter) (IRuntimeValue, phpError.Error) {
	// Spec: https://www.php.net/manual/en/function.getenv.php

	//  getenv(?string $name = null, bool $local_only = false): string|array|false

	// Returns the value of the environment variable name, or false if the environment variable name does not exist.
	// If name is null, all environment variables are returned as an associative array.

	// TODO getenv - add support for $local_only
	args, err := NewFuncParamValidator("getenv").addParam("$name", []string{"string"}, NewNullRuntimeValue()).validate(args)
	if err != nil {
		return NewVoidRuntimeValue(), err
	}

	if args[0].GetType() == NullValue {
		return interpreter.env.lookupVariable("$_ENV")
	}

	envVars, err := interpreter.env.lookupVariable("$_ENV")
	if err != nil {
		return envVars, err
	}
	envArray := runtimeValToArrayRuntimeVal(envVars)
	value, found := envArray.GetElement(args[0])
	if !found {
		return NewBooleanRuntimeValue(false), nil
	}
	return value, nil
}

// ------------------- MARK: ini_get -------------------

func nativeFn_ini_get(args []IRuntimeValue, interpreter *Interpreter) (IRuntimeValue, phpError.Error) {
	// Spec: https://www.php.net/manual/en/function.ini-get

	args, err := NewFuncParamValidator("ini_get").addParam("$option", []string{"string"}, nil).validate(args)
	if err != nil {
		return NewVoidRuntimeValue(), err
	}

	switch runtimeValToStrRuntimeVal(args[0]).GetValue() {
	case "error_reporting":
		return NewStringRuntimeValue(fmt.Sprintf("%d", interpreter.ini.ErrorReporting)), nil
	case "short_open_tag":
		if interpreter.ini.ShortOpenTag {
			return NewStringRuntimeValue("1"), nil
		}
		return NewStringRuntimeValue("0"), nil
	default:
		return NewBooleanRuntimeValue(false), nil
	}
}

// ------------------- MARK: strlen -------------------

func nativeFn_strlen(args []IRuntimeValue, _ *Interpreter) (IRuntimeValue, phpError.Error) {
	// Spec: https://www.php.net/manual/en/function.strlen

	args, err := NewFuncParamValidator("strlen").addParam("$string", []string{"string"}, nil).validate(args)
	if err != nil {
		return NewVoidRuntimeValue(), err
	}

	return NewIntegerRuntimeValue(int64(len(runtimeValToStrRuntimeVal(args[0]).GetValue()))), nil
}
