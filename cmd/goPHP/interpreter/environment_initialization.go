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
		registerPredefinedVariableEnv(environment, request, true)
		registerPredefinedVariableGet(environment, request, true)
		registerPredefinedVariablePost(environment, request, true)
		// TODO Cookie
		registerPredefinedVariableServer(environment, request, true)
		return
	}

	variables_order := ini.GetStr("variables_order")
	for _, variable := range variables_order {
		switch string(variable) {
		case "E":
			registerPredefinedVariableEnv(environment, request, true)
		case "G":
			registerPredefinedVariableGet(environment, request, true)
		case "P":
			registerPredefinedVariablePost(environment, request, true)
		// TODO Cookie
		// case "C":
		// 	registerPredefinedVariableCookie(environment, request, true)
		case "S":
			registerPredefinedVariableServer(environment, request, true)
		}
	}

	if !strings.Contains(variables_order, "E") {
		registerPredefinedVariableEnv(environment, request, false)
	}
	if !strings.Contains(variables_order, "G") {
		registerPredefinedVariableGet(environment, request, false)
	}
	if !strings.Contains(variables_order, "P") {
		registerPredefinedVariablePost(environment, request, false)
	}
	// TODO Cookie
	// if !strings.Contains(variables_order, "C") {
	// 	registerPredefinedVariableCookie(environment, request, false)
	// }
	if !strings.Contains(variables_order, "S") {
		registerPredefinedVariableServer(environment, request, false)
	}
}

func registerPredefinedVariableEnv(environment *Environment, request *Request, init bool) {
	if init {
		environment.predefinedVariables["$_ENV"] = stringMapToArray(request.Env)
	} else {
		environment.predefinedVariables["$_ENV"] = NewArrayRuntimeValue()
	}
}

func registerPredefinedVariableGet(environment *Environment, request *Request, init bool) {
	if init {
		array, err := parseQuery(request.QueryString)
		if err != nil {
			println(err.Error())
		}
		environment.predefinedVariables["$_GET"] = array
	} else {
		environment.predefinedVariables["$_GET"] = NewArrayRuntimeValue()
	}
}

func registerPredefinedVariablePost(environment *Environment, request *Request, init bool) {
	if init {
		environment.predefinedVariables["$_POST"] = paramToArray(request.PostParams)
	} else {
		environment.predefinedVariables["$_POST"] = NewArrayRuntimeValue()
	}
}

func registerPredefinedVariableServer(environment *Environment, request *Request, init bool) {
	environment.predefinedVariables["$_SERVER"] = NewArrayRuntimeValue()
	if init {
		server := environment.predefinedVariables["$_SERVER"].(*ArrayRuntimeValue)
		if len(request.Args) > 0 {
			server.SetElement(NewStringRuntimeValue("argc"), NewIntegerRuntimeValue(int64(len(request.Args))))
			server.SetElement(NewStringRuntimeValue("argv"), paramToArray(request.Args))
		}
		server.SetElement(NewStringRuntimeValue("DOCUMENT_ROOT"), NewStringRuntimeValue(request.DocumentRoot))
		server.SetElement(NewStringRuntimeValue("QUERY_STRING"), NewStringRuntimeValue(request.QueryString))
		server.SetElement(NewStringRuntimeValue("REMOTE_ADDR"), NewStringRuntimeValue(request.RemoteAddr))
		server.SetElement(NewStringRuntimeValue("REMOTE_PORT"), NewStringRuntimeValue(request.RemotePort))
		server.SetElement(NewStringRuntimeValue("REQUEST_METHOD"), NewStringRuntimeValue(request.Method))
		server.SetElement(NewStringRuntimeValue("REQUEST_TIME_FLOAT"), NewFloatingRuntimeValue(float64(request.RequestTime.UnixMicro())/math.Pow(10, 6)))
		server.SetElement(NewStringRuntimeValue("REQUEST_TIME"), NewIntegerRuntimeValue(request.RequestTime.Unix()))
		server.SetElement(NewStringRuntimeValue("REQUEST_URI"), NewStringRuntimeValue(request.RequestURI))
		server.SetElement(NewStringRuntimeValue("SCRIPT_FILENAME"), NewStringRuntimeValue(request.ScriptFilename))
		server.SetElement(NewStringRuntimeValue("SCRIPT_NAME"), NewStringRuntimeValue(strings.Replace(request.ScriptFilename, request.DocumentRoot, "", 1)))
		server.SetElement(NewStringRuntimeValue("SERVER_ADDR"), NewStringRuntimeValue(request.ServerAddr))
		server.SetElement(NewStringRuntimeValue("SERVER_PORT"), NewStringRuntimeValue(request.ServerPort))
		server.SetElement(NewStringRuntimeValue("SERVER_PROTOCOL"), NewStringRuntimeValue(request.Protocol))
		server.SetElement(NewStringRuntimeValue("SERVER_SOFTWARE"), NewStringRuntimeValue(config.SoftwareVersion))
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
		hasKey := false
		var key string
		var value string
		if len(param) == 2 {
			key = param[0]
			value = param[1]
			hasKey = true
		} else if len(param) == 1 {
			value = param[0]
		}

		// No array
		if !strings.Contains(key, "]") {
			var keyValue IRuntimeValue
			if hasKey {
				if common.IsIntegerLiteral(key) {
					intValue, _ := common.IntegerLiteralToInt64(key)
					keyValue = NewIntegerRuntimeValue(intValue)
				} else {
					keyValue = NewStringRuntimeValue(key)
				}
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
