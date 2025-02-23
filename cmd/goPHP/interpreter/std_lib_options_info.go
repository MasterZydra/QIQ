package interpreter

import (
	"GoPHP/cmd/goPHP/ini"
	"GoPHP/cmd/goPHP/phpError"
)

func registerNativeOptionsInfoFunctions(environment *Environment) {
	environment.nativeFunctions["getenv"] = nativeFn_getenv
	environment.nativeFunctions["ini_get"] = nativeFn_ini_get
	environment.nativeFunctions["ini_set"] = nativeFn_ini_set
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
	envArray := envVars.(*ArrayRuntimeValue)
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

	value, iniErr := interpreter.ini.Get(args[0].(*StringRuntimeValue).Value)
	if iniErr != nil {
		return NewBooleanRuntimeValue(false), nil
	}
	return NewStringRuntimeValue(value), nil
}

// ------------------- MARK: ini_set -------------------

func nativeFn_ini_set(args []IRuntimeValue, interpreter *Interpreter) (IRuntimeValue, phpError.Error) {
	// Spec: https://www.php.net/manual/en/function.ini-set

	args, err := NewFuncParamValidator("ini_set").
		addParam("$option", []string{"string"}, nil).
		addParam("$value", []string{"string", "int", "float", "bool", "null"}, nil).
		validate(args)
	if err != nil {
		return NewVoidRuntimeValue(), err
	}

	value, err := lib_strval(args[1])
	if err != nil {
		return NewVoidRuntimeValue(), err
	}

	option := args[0].(*StringRuntimeValue).Value

	oldValue, err := interpreter.ini.Get(option)
	if err != nil {
		return NewBooleanRuntimeValue(false), nil
	}
	err = interpreter.ini.Set(option, value, ini.INI_USER)
	if err != nil {
		return NewBooleanRuntimeValue(false), nil
	}

	return NewStringRuntimeValue(oldValue), nil
}
