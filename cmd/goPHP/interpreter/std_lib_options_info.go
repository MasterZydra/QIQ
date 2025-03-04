package interpreter

import (
	"GoPHP/cmd/goPHP/ini"
	"GoPHP/cmd/goPHP/phpError"
	"GoPHP/cmd/goPHP/runtime/values"
)

func registerNativeOptionsInfoFunctions(environment *Environment) {
	environment.nativeFunctions["getenv"] = nativeFn_getenv
	environment.nativeFunctions["ini_get"] = nativeFn_ini_get
	environment.nativeFunctions["ini_set"] = nativeFn_ini_set
}

// ------------------- MARK: getenv -------------------

func nativeFn_getenv(args []values.RuntimeValue, interpreter *Interpreter) (values.RuntimeValue, phpError.Error) {
	// Spec: https://www.php.net/manual/en/function.getenv.php

	//  getenv(?string $name = null, bool $local_only = false): string|array|false

	// Returns the value of the environment variable name, or false if the environment variable name does not exist.
	// If name is null, all environment variables are returned as an associative array.

	// TODO getenv - add support for $local_only
	args, err := NewFuncParamValidator("getenv").addParam("$name", []string{"string"}, values.NewNull()).validate(args)
	if err != nil {
		return values.NewVoid(), err
	}

	if args[0].GetType() == values.NullValue {
		return interpreter.env.lookupVariable("$_ENV")
	}

	envVars, err := interpreter.env.lookupVariable("$_ENV")
	if err != nil {
		return envVars, err
	}
	envArray := envVars.(*values.Array)
	value, found := envArray.GetElement(args[0])
	if !found {
		return values.NewBool(false), nil
	}
	return value, nil
}

// ------------------- MARK: ini_get -------------------

func nativeFn_ini_get(args []values.RuntimeValue, interpreter *Interpreter) (values.RuntimeValue, phpError.Error) {
	// Spec: https://www.php.net/manual/en/function.ini-get

	args, err := NewFuncParamValidator("ini_get").addParam("$option", []string{"string"}, nil).validate(args)
	if err != nil {
		return values.NewVoid(), err
	}

	value, iniErr := interpreter.ini.Get(args[0].(*values.Str).Value)
	if iniErr != nil {
		return values.NewBool(false), nil
	}
	return values.NewStr(value), nil
}

// ------------------- MARK: ini_set -------------------

func nativeFn_ini_set(args []values.RuntimeValue, interpreter *Interpreter) (values.RuntimeValue, phpError.Error) {
	// Spec: https://www.php.net/manual/en/function.ini-set

	args, err := NewFuncParamValidator("ini_set").
		addParam("$option", []string{"string"}, nil).
		addParam("$value", []string{"string", "int", "float", "bool", "null"}, nil).
		validate(args)
	if err != nil {
		return values.NewVoid(), err
	}

	value, err := lib_strval(args[1])
	if err != nil {
		return values.NewVoid(), err
	}

	option := args[0].(*values.Str).Value

	oldValue, err := interpreter.ini.Get(option)
	if err != nil {
		return values.NewBool(false), nil
	}
	err = interpreter.ini.Set(option, value, ini.INI_USER)
	if err != nil {
		return values.NewBool(false), nil
	}

	return values.NewStr(oldValue), nil
}
