package optionsInfo

import (
	"QIQ/cmd/qiq/ini"
	"QIQ/cmd/qiq/phpError"
	"QIQ/cmd/qiq/runtime"
	"QIQ/cmd/qiq/runtime/funcParamValidator"
	"QIQ/cmd/qiq/runtime/stdlib/variableHandling"
	"QIQ/cmd/qiq/runtime/values"
)

func Register(environment runtime.Environment) {
	// Category: Options/Info Functions
	environment.AddNativeFunction("getenv", nativeFn_getenv)
	environment.AddNativeFunction("ini_get", nativeFn_ini_get)
	environment.AddNativeFunction("ini_set", nativeFn_ini_set)

	// Const Category: Options/Info Constants
	// Spec: https://www.php.net/manual/en/info.constants.php
	environment.AddPredefinedConstant("INI_USER", values.NewInt(int64(ini.INI_USER)))
	environment.AddPredefinedConstant("INI_PERDIR", values.NewInt(int64(ini.INI_PERDIR)))
	environment.AddPredefinedConstant("INI_SYSTEM", values.NewInt(int64(ini.INI_SYSTEM)))
	environment.AddPredefinedConstant("INI_ALL", values.NewInt(int64(ini.INI_ALL)))
}

// -------------------------------------- getenv -------------------------------------- MARK: getenv

func nativeFn_getenv(args []values.RuntimeValue, context runtime.Context) (values.RuntimeValue, phpError.Error) {
	// Spec: https://www.php.net/manual/en/function.getenv.php

	//  getenv(?string $name = null, bool $local_only = false): string|array|false

	// Returns the value of the environment variable name, or false if the environment variable name does not exist.
	// If name is null, all environment variables are returned as an associative array.

	// TODO getenv - add support for $local_only
	args, err := funcParamValidator.NewValidator("getenv").AddParam("$name", []string{"string"}, values.NewNull()).Validate(args)
	if err != nil {
		return values.NewVoid(), err
	}

	if args[0].GetType() == values.NullValue {
		return context.Env.LookupVariable("$_ENV")
	}

	envVars, err := context.Env.LookupVariable("$_ENV")
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

// -------------------------------------- ini_get -------------------------------------- MARK: ini_get

func nativeFn_ini_get(args []values.RuntimeValue, context runtime.Context) (values.RuntimeValue, phpError.Error) {
	// Spec: https://www.php.net/manual/en/function.ini-get

	args, err := funcParamValidator.NewValidator("ini_get").AddParam("$option", []string{"string"}, nil).Validate(args)
	if err != nil {
		return values.NewVoid(), err
	}

	value, iniErr := context.Interpreter.GetIni().Get(args[0].(*values.Str).Value)
	if iniErr != nil {
		return values.NewBool(false), nil
	}
	return values.NewStr(value), nil
}

// -------------------------------------- ini_set -------------------------------------- MARK: ini_set

func nativeFn_ini_set(args []values.RuntimeValue, context runtime.Context) (values.RuntimeValue, phpError.Error) {
	// Spec: https://www.php.net/manual/en/function.ini-set

	args, err := funcParamValidator.NewValidator("ini_set").
		AddParam("$option", []string{"string"}, nil).
		AddParam("$value", []string{"string", "int", "float", "bool", "null"}, nil).
		Validate(args)
	if err != nil {
		return values.NewVoid(), err
	}

	value, err := variableHandling.StrVal(args[1])
	if err != nil {
		return values.NewVoid(), err
	}

	option := args[0].(*values.Str).Value

	oldValue, err := context.Interpreter.GetIni().Get(option)
	if err != nil {
		return values.NewBool(false), nil
	}
	err = context.Interpreter.GetIni().Set(option, value, ini.INI_USER)
	if err != nil {
		return values.NewBool(false), nil
	}

	return values.NewStr(oldValue), nil
}
