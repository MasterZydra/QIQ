package interpreter

import (
	"GoPHP/cmd/goPHP/common"
	"GoPHP/cmd/goPHP/config"
	"GoPHP/cmd/goPHP/ini"
	"GoPHP/cmd/goPHP/phpError"
	"fmt"
	"math"
	"regexp"
	"strings"
)

func registerPredefinedVariables(environment *Environment, request *Request, ini *ini.Ini) {
	if ini == nil {
		registerPredefinedVariableEnv(environment, request)
		registerPredefinedVariableGet(environment, request)
		registerPredefinedVariablePost(environment, request)
		// TODO Cookie
		registerPredefinedVariableServer(environment, request)
		return
	}

	for _, variable := range ini.GetStr("variables_order") {
		switch string(variable) {
		case "E":
			registerPredefinedVariableEnv(environment, request)
		case "G":
			registerPredefinedVariableGet(environment, request)
		case "P":
			registerPredefinedVariablePost(environment, request)
		// case "C":
		// TODO Cookie
		case "S":
			registerPredefinedVariableServer(environment, request)
		}
	}
}

func registerPredefinedVariableEnv(environment *Environment, request *Request) {
	environment.predefinedVariables["$_ENV"] = stringMapToArray(request.Env)
}

func registerPredefinedVariableGet(environment *Environment, request *Request) {
	environment.predefinedVariables["$_GET"] = paramToArray(request.GetParams)
}

func registerPredefinedVariablePost(environment *Environment, request *Request) {
	environment.predefinedVariables["$_POST"] = paramToArray(request.PostParams)
}

func registerPredefinedVariableServer(environment *Environment, request *Request) {
	environment.predefinedVariables["$_SERVER"] = NewArrayRuntimeValue()
	if len(request.Args) > 0 {
		server := environment.predefinedVariables["$_SERVER"].(*ArrayRuntimeValue)
		server.SetElement(NewStringRuntimeValue("argc"), NewIntegerRuntimeValue(int64(len(request.Args))))
		server.SetElement(NewStringRuntimeValue("argv"), paramToArray(request.Args))
	}
}

func stringMapToArray(stringMap map[string]string) *ArrayRuntimeValue {
	result := NewArrayRuntimeValue()
	for key, value := range stringMap {
		result.SetElement(NewStringRuntimeValue(key), NewStringRuntimeValue(value))
	}
	return result
}

func paramToArray(params [][]string) *ArrayRuntimeValue {
	result := NewArrayRuntimeValue()

	for _, param := range params {
		key := param[0]
		value := param[1]

		// No array
		if !strings.Contains(key, "]") {
			var keyValue IRuntimeValue
			if common.IsIntegerLiteral(key) {
				intValue, _ := common.IntegerLiteralToInt64(key)
				keyValue = NewIntegerRuntimeValue(intValue)
			} else {
				keyValue = NewStringRuntimeValue(key)
			}
			result.SetElement(keyValue, NewStringRuntimeValue(value))
			continue
		}

		// Array

		openingBracket := strings.Index(key, "[")
		// Get name of param without brackets
		paramName := key[:openingBracket]

		// Check if array is already in params
		arrayValue, found := result.GetElement(NewStringRuntimeValue(paramName))
		if !found || arrayValue.GetType() != ArrayValue {
			arrayValue = NewArrayRuntimeValue()
		}

		// Wrap keys that are strings in double quotes
		decimalKeys, _ := regexp.Compile(`\[[0-9]+\]`)
		nondecimalKeys, _ := regexp.Compile(`\[.+\]`)
		matches := nondecimalKeys.FindAllString(decimalKeys.ReplaceAllString(key, ""), -1)
		for _, match := range matches {
			replacement := `["` + match[1:len(match)-1] + `"]`
			key = strings.Replace(key, match, replacement, 1)
		}

		// Prepare environment
		env := NewEnvironment(nil, NewRequest(), nil)
		env.declareVariable("$"+paramName, arrayValue)

		// Execute PHP to store new array values in env
		interpreter := NewInterpreter(ini.NewDefaultIni(), NewRequest(), "")
		interpreter.process(fmt.Sprintf(`<?php $%s = "%s";`, key, value), env)

		// Extract array from environment
		arrayValue = env.variables["$"+paramName]

		result.SetElement(NewStringRuntimeValue(paramName), arrayValue)
		continue
	}
	return result
}

func registerPredefinedConstants(environment *Environment) {
	// Spec: https://phplang.org/spec/06-constants.html#core-predefined-constants
	// Spec: https://www.php.net/manual/en/reserved.constants.php
	environment.predefinedConstants["FALSE"] = NewBooleanRuntimeValue(false)
	environment.predefinedConstants["TRUE"] = NewBooleanRuntimeValue(true)
	environment.predefinedConstants["NULL"] = NewNullRuntimeValue()
	environment.predefinedConstants["PHP_INT_MAX"] = NewIntegerRuntimeValue(math.MaxInt64)
	environment.predefinedConstants["PHP_INT_MIN"] = NewIntegerRuntimeValue(math.MinInt64)
	environment.predefinedConstants["PHP_INT_SIZE"] = NewIntegerRuntimeValue(64 / 8)
	environment.predefinedConstants["PHP_OS"] = NewStringRuntimeValue(getPhpOs())
	environment.predefinedConstants["PHP_OS_FAMILY"] = NewStringRuntimeValue(getPhpOsFamily())
	if getPhpOs() == "Windows" {
		environment.predefinedConstants["PHP_EOL"] = NewStringRuntimeValue("\r\n")
	} else {
		environment.predefinedConstants["PHP_EOL"] = NewStringRuntimeValue("\n")
	}
	environment.predefinedConstants["PHP_VERSION"] = NewStringRuntimeValue(config.Version)
	environment.predefinedConstants["PHP_MAJOR_VERSION"] = NewIntegerRuntimeValue(config.MajorVersion)
	environment.predefinedConstants["PHP_MINOR_VERSION"] = NewIntegerRuntimeValue(config.MinorVersion)
	environment.predefinedConstants["PHP_RELEASE_VERSION"] = NewIntegerRuntimeValue(config.ReleaseVersion)
	environment.predefinedConstants["PHP_EXTRA_VERSION"] = NewStringRuntimeValue(config.ExtraVersion)
	environment.predefinedConstants["PHP_VERSION_ID"] = NewIntegerRuntimeValue(config.VersionId)

	// Spec: https://www.php.net/manual/en/math.constants.php
	environment.predefinedConstants["M_1_PI"] = NewFloatingRuntimeValue(1 / math.Pi)
	environment.predefinedConstants["M_2_PI"] = NewFloatingRuntimeValue(2 / math.Pi)
	environment.predefinedConstants["M_2_SQRTPI"] = NewFloatingRuntimeValue(2 / math.SqrtPi)
	environment.predefinedConstants["M_E"] = NewFloatingRuntimeValue(math.E)
	environment.predefinedConstants["M_EULER"] = NewFloatingRuntimeValue(math.E)
	environment.predefinedConstants["M_LN10"] = NewFloatingRuntimeValue(math.Ln10)
	environment.predefinedConstants["M_LN2"] = NewFloatingRuntimeValue(math.Ln2)
	environment.predefinedConstants["M_LNPI"] = NewFloatingRuntimeValue(math.Log(math.Pi))
	environment.predefinedConstants["M_LOG10E"] = NewFloatingRuntimeValue(math.Log10E)
	environment.predefinedConstants["M_LOG2E"] = NewFloatingRuntimeValue(math.Log2E)
	environment.predefinedConstants["M_PI"] = NewFloatingRuntimeValue(math.Pi)
	environment.predefinedConstants["M_PI_2"] = NewFloatingRuntimeValue(math.Pi / 2)
	environment.predefinedConstants["M_PI_4"] = NewFloatingRuntimeValue(math.Pi / 4)
	environment.predefinedConstants["M_SQRT1_2"] = NewFloatingRuntimeValue(1 / math.Sqrt2)
	environment.predefinedConstants["M_SQRT2"] = NewFloatingRuntimeValue(math.Sqrt2)
	environment.predefinedConstants["M_SQRT3"] = NewFloatingRuntimeValue(math.Sqrt(3))
	environment.predefinedConstants["M_SQRTPI"] = NewFloatingRuntimeValue(math.SqrtPi)
	environment.predefinedConstants["PHP_ROUND_HALF_UP"] = NewIntegerRuntimeValue(1)
	environment.predefinedConstants["PHP_ROUND_HALF_DOWN"] = NewIntegerRuntimeValue(2)
	environment.predefinedConstants["PHP_ROUND_HALF_EVEN"] = NewIntegerRuntimeValue(3)
	environment.predefinedConstants["PHP_ROUND_HALF_ODD"] = NewIntegerRuntimeValue(4)

	// Spec: https://www.php.net/manual/en/errorfunc.constants.php
	environment.predefinedConstants["E_ERROR"] = NewIntegerRuntimeValue(phpError.E_ERROR)
	environment.predefinedConstants["E_WARNING"] = NewIntegerRuntimeValue(phpError.E_WARNING)
	environment.predefinedConstants["E_PARSE"] = NewIntegerRuntimeValue(phpError.E_PARSE)
	environment.predefinedConstants["E_NOTICE"] = NewIntegerRuntimeValue(phpError.E_NOTICE)
	environment.predefinedConstants["E_CORE_ERROR"] = NewIntegerRuntimeValue(phpError.E_CORE_ERROR)
	environment.predefinedConstants["E_CORE_WARNING"] = NewIntegerRuntimeValue(phpError.E_CORE_WARNING)
	environment.predefinedConstants["E_COMPILE_ERROR"] = NewIntegerRuntimeValue(phpError.E_COMPILE_ERROR)
	environment.predefinedConstants["E_COMPILE_WARNING"] = NewIntegerRuntimeValue(phpError.E_COMPILE_WARNING)
	environment.predefinedConstants["E_USER_ERROR"] = NewIntegerRuntimeValue(phpError.E_USER_ERROR)
	environment.predefinedConstants["E_USER_WARNING"] = NewIntegerRuntimeValue(phpError.E_USER_WARNING)
	environment.predefinedConstants["E_USER_NOTICE"] = NewIntegerRuntimeValue(phpError.E_USER_NOTICE)
	environment.predefinedConstants["E_STRICT"] = NewIntegerRuntimeValue(phpError.E_STRICT)
	environment.predefinedConstants["E_RECOVERABLE_ERROR"] = NewIntegerRuntimeValue(phpError.E_RECOVERABLE_ERROR)
	environment.predefinedConstants["E_DEPRECATED"] = NewIntegerRuntimeValue(phpError.E_DEPRECATED)
	environment.predefinedConstants["E_USER_DEPRECATED"] = NewIntegerRuntimeValue(phpError.E_USER_DEPRECATED)
	environment.predefinedConstants["E_ALL"] = NewIntegerRuntimeValue(phpError.E_ALL)
}
