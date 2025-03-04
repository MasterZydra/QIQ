package interpreter

import (
	"GoPHP/cmd/goPHP/common"
	"GoPHP/cmd/goPHP/config"
	"GoPHP/cmd/goPHP/ini"
	"GoPHP/cmd/goPHP/phpError"
	"GoPHP/cmd/goPHP/runtime/values"
	"fmt"
	"math"
	"regexp"
	"strings"
)

func registerPredefinedVariables(environment *Environment, request *Request, ini *ini.Ini) {
	if ini == nil {
		registerPredefinedVariableEnv(environment, request, true)
		registerPredefinedVariableGet(environment, request, ini, true)
		registerPredefinedVariablePost(environment, request, ini, true)
		// TODO Cookie
		registerPredefinedVariableServer(environment, request, ini, true)
		return
	}

	variables_order := ini.GetStr("variables_order")
	for _, variable := range variables_order {
		switch string(variable) {
		case "E":
			registerPredefinedVariableEnv(environment, request, true)
		case "G":
			registerPredefinedVariableGet(environment, request, ini, true)
		case "P":
			registerPredefinedVariablePost(environment, request, ini, true)
		// TODO Cookie
		// case "C":
		// 	registerPredefinedVariableCookie(environment, request, true)
		case "S":
			registerPredefinedVariableServer(environment, request, ini, true)
		}
	}

	if !strings.Contains(variables_order, "E") {
		registerPredefinedVariableEnv(environment, request, false)
	}
	if !strings.Contains(variables_order, "G") {
		registerPredefinedVariableGet(environment, request, ini, false)
	}
	if !strings.Contains(variables_order, "P") {
		registerPredefinedVariablePost(environment, request, ini, false)
	}
	// TODO Cookie
	// if !strings.Contains(variables_order, "C") {
	// 	registerPredefinedVariableCookie(environment, request, false)
	// }
	if !strings.Contains(variables_order, "S") {
		registerPredefinedVariableServer(environment, request, ini, false)
	}

	requestVar := values.NewArray()
	mergeArrays(requestVar, environment.predefinedVariables["$_GET"].(*values.Array))
	mergeArrays(requestVar, environment.predefinedVariables["$_POST"].(*values.Array))
	// TODO Cookie
	// mergeArrays(requestVar, environment.predefinedVariables["$_COOKIE"].(*values.Array))
	environment.predefinedVariables["$_REQUEST"] = requestVar
}

// TODO Replace with std lib func array_merge
func mergeArrays(a, b *values.Array) {
	for _, key := range b.Keys {
		value, _ := b.GetElement(key)
		a.SetElement(key, deepCopy(value))
	}
}

func registerPredefinedVariableEnv(environment *Environment, request *Request, init bool) {
	if init {
		environment.predefinedVariables["$_ENV"] = stringMapToArray(request.Env)
	} else {
		environment.predefinedVariables["$_ENV"] = values.NewArray()
	}
}

func registerPredefinedVariableGet(environment *Environment, request *Request, ini *ini.Ini, init bool) {
	if init {
		array, err := parseQuery(request.QueryString, ini)
		if err != nil {
			println(err.Error())
		}
		environment.predefinedVariables["$_GET"] = array
	} else {
		environment.predefinedVariables["$_GET"] = values.NewArray()
	}
}

func registerPredefinedVariablePost(environment *Environment, request *Request, ini *ini.Ini, init bool) {
	if init {
		array, err := parseQuery(request.Post, ini)
		if err != nil {
			println(err.Error())
		}
		environment.predefinedVariables["$_POST"] = array
	} else {
		environment.predefinedVariables["$_POST"] = values.NewArray()
	}
}

func registerPredefinedVariableServer(environment *Environment, request *Request, ini *ini.Ini, init bool) {
	environment.predefinedVariables["$_SERVER"] = values.NewArray()
	if init {
		server := environment.predefinedVariables["$_SERVER"].(*values.Array)
		if len(request.Args) > 0 {
			server.SetElement(values.NewStr("argc"), values.NewInt(int64(len(request.Args))))
			server.SetElement(values.NewStr("argv"), paramToArray(request.Args, ini))
		}
		server.SetElement(values.NewStr("DOCUMENT_ROOT"), values.NewStr(request.DocumentRoot))
		server.SetElement(values.NewStr("QUERY_STRING"), values.NewStr(request.QueryString))
		server.SetElement(values.NewStr("REMOTE_ADDR"), values.NewStr(request.RemoteAddr))
		server.SetElement(values.NewStr("REMOTE_PORT"), values.NewStr(request.RemotePort))
		server.SetElement(values.NewStr("REQUEST_METHOD"), values.NewStr(request.Method))
		server.SetElement(values.NewStr("REQUEST_TIME_FLOAT"), values.NewFloat(float64(request.RequestTime.UnixMicro())/math.Pow(10, 6)))
		server.SetElement(values.NewStr("REQUEST_TIME"), values.NewInt(request.RequestTime.Unix()))
		server.SetElement(values.NewStr("REQUEST_URI"), values.NewStr(request.RequestURI))
		server.SetElement(values.NewStr("SCRIPT_FILENAME"), values.NewStr(request.ScriptFilename))
		server.SetElement(values.NewStr("SCRIPT_NAME"), values.NewStr(strings.Replace(request.ScriptFilename, request.DocumentRoot, "", 1)))
		server.SetElement(values.NewStr("SERVER_ADDR"), values.NewStr(request.ServerAddr))
		server.SetElement(values.NewStr("SERVER_PORT"), values.NewStr(request.ServerPort))
		server.SetElement(values.NewStr("SERVER_PROTOCOL"), values.NewStr(request.Protocol))
		server.SetElement(values.NewStr("SERVER_SOFTWARE"), values.NewStr(config.SoftwareVersion))
	}
}

func stringMapToArray(stringMap map[string]string) *values.Array {
	result := values.NewArray()
	for key, value := range stringMap {
		result.SetElement(values.NewStr(key), values.NewStr(value))
	}
	return result
}

func paramToArray(params [][]string, ini *ini.Ini) *values.Array {
	result := values.NewArray()

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
			var keyValue values.RuntimeValue
			if hasKey {
				if common.IsIntegerLiteral(key) {
					intValue, _ := common.IntegerLiteralToInt64(key)
					keyValue = values.NewInt(intValue)
				} else {
					keyValue = values.NewStr(key)
				}
			}
			result.SetElement(keyValue, values.NewStr(value))
			continue
		}

		// Array

		openingBracket := strings.Index(key, "[")
		// Get name of param without brackets
		paramName := key[:openingBracket]

		// Check if array is already in params
		arrayValue, found := result.GetElement(values.NewStr(paramName))
		if !found || arrayValue.GetType() != values.ArrayValue {
			arrayValue = values.NewArray()
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
		env := NewEnvironment(nil, NewRequest(), ini)
		env.declareVariable("$"+paramName, arrayValue)

		// Execute PHP to store new array values in env
		interpreter := NewInterpreter(ini, NewRequest(), "")
		interpreter.process(fmt.Sprintf(`<?php $%s = "%s";`, key, value), env, true)

		// Extract array from environment
		arrayValue = env.variables["$"+paramName]

		result.SetElement(values.NewStr(paramName), arrayValue)
		continue
	}
	return result
}

func registerPredefinedConstants(environment *Environment) {
	// Spec: https://phplang.org/spec/06-constants.html#core-predefined-constants
	// Spec: https://www.php.net/manual/en/reserved.constants.php
	environment.predefinedConstants["DIRECTORY_SEPARATOR"] = values.NewStr(getPhpDirectorySeparator())
	environment.predefinedConstants["FALSE"] = values.NewBool(false)
	environment.predefinedConstants["TRUE"] = values.NewBool(true)
	environment.predefinedConstants["NULL"] = values.NewNull()
	environment.predefinedConstants["PHP_INT_MAX"] = values.NewInt(math.MaxInt64)
	environment.predefinedConstants["PHP_INT_MIN"] = values.NewInt(math.MinInt64)
	environment.predefinedConstants["PHP_INT_SIZE"] = values.NewInt(64 / 8)
	environment.predefinedConstants["PHP_OS"] = values.NewStr(getPhpOs())
	environment.predefinedConstants["PHP_OS_FAMILY"] = values.NewStr(getPhpOsFamily())
	environment.predefinedConstants["PHP_EOL"] = values.NewStr(getPhpEol())
	environment.predefinedConstants["PHP_VERSION"] = values.NewStr(config.Version)
	environment.predefinedConstants["PHP_MAJOR_VERSION"] = values.NewInt(config.MajorVersion)
	environment.predefinedConstants["PHP_MINOR_VERSION"] = values.NewInt(config.MinorVersion)
	environment.predefinedConstants["PHP_RELEASE_VERSION"] = values.NewInt(config.ReleaseVersion)
	environment.predefinedConstants["PHP_EXTRA_VERSION"] = values.NewStr(config.ExtraVersion)
	environment.predefinedConstants["PHP_VERSION_ID"] = values.NewInt(config.VersionId)

	// Spec: https://www.php.net/manual/en/math.constants.php
	environment.predefinedConstants["M_1_PI"] = values.NewFloat(1 / math.Pi)
	environment.predefinedConstants["M_2_PI"] = values.NewFloat(2 / math.Pi)
	environment.predefinedConstants["M_2_SQRTPI"] = values.NewFloat(2 / math.SqrtPi)
	environment.predefinedConstants["M_E"] = values.NewFloat(math.E)
	environment.predefinedConstants["M_EULER"] = values.NewFloat(math.E)
	environment.predefinedConstants["M_LN10"] = values.NewFloat(math.Ln10)
	environment.predefinedConstants["M_LN2"] = values.NewFloat(math.Ln2)
	environment.predefinedConstants["M_LNPI"] = values.NewFloat(math.Log(math.Pi))
	environment.predefinedConstants["M_LOG10E"] = values.NewFloat(math.Log10E)
	environment.predefinedConstants["M_LOG2E"] = values.NewFloat(math.Log2E)
	environment.predefinedConstants["M_PI"] = values.NewFloat(math.Pi)
	environment.predefinedConstants["M_PI_2"] = values.NewFloat(math.Pi / 2)
	environment.predefinedConstants["M_PI_4"] = values.NewFloat(math.Pi / 4)
	environment.predefinedConstants["M_SQRT1_2"] = values.NewFloat(1 / math.Sqrt2)
	environment.predefinedConstants["M_SQRT2"] = values.NewFloat(math.Sqrt2)
	environment.predefinedConstants["M_SQRT3"] = values.NewFloat(math.Sqrt(3))
	environment.predefinedConstants["M_SQRTPI"] = values.NewFloat(math.SqrtPi)
	environment.predefinedConstants["PHP_ROUND_HALF_UP"] = values.NewInt(1)
	environment.predefinedConstants["PHP_ROUND_HALF_DOWN"] = values.NewInt(2)
	environment.predefinedConstants["PHP_ROUND_HALF_EVEN"] = values.NewInt(3)
	environment.predefinedConstants["PHP_ROUND_HALF_ODD"] = values.NewInt(4)

	// Spec: https://www.php.net/manual/en/errorfunc.constants.php
	environment.predefinedConstants["E_ERROR"] = values.NewInt(phpError.E_ERROR)
	environment.predefinedConstants["E_WARNING"] = values.NewInt(phpError.E_WARNING)
	environment.predefinedConstants["E_PARSE"] = values.NewInt(phpError.E_PARSE)
	environment.predefinedConstants["E_NOTICE"] = values.NewInt(phpError.E_NOTICE)
	environment.predefinedConstants["E_CORE_ERROR"] = values.NewInt(phpError.E_CORE_ERROR)
	environment.predefinedConstants["E_CORE_WARNING"] = values.NewInt(phpError.E_CORE_WARNING)
	environment.predefinedConstants["E_COMPILE_ERROR"] = values.NewInt(phpError.E_COMPILE_ERROR)
	environment.predefinedConstants["E_COMPILE_WARNING"] = values.NewInt(phpError.E_COMPILE_WARNING)
	environment.predefinedConstants["E_USER_ERROR"] = values.NewInt(phpError.E_USER_ERROR)
	environment.predefinedConstants["E_USER_WARNING"] = values.NewInt(phpError.E_USER_WARNING)
	environment.predefinedConstants["E_USER_NOTICE"] = values.NewInt(phpError.E_USER_NOTICE)
	environment.predefinedConstants["E_STRICT"] = values.NewInt(phpError.E_STRICT)
	environment.predefinedConstants["E_RECOVERABLE_ERROR"] = values.NewInt(phpError.E_RECOVERABLE_ERROR)
	environment.predefinedConstants["E_DEPRECATED"] = values.NewInt(phpError.E_DEPRECATED)
	environment.predefinedConstants["E_USER_DEPRECATED"] = values.NewInt(phpError.E_USER_DEPRECATED)
	environment.predefinedConstants["E_ALL"] = values.NewInt(phpError.E_ALL)

	// Spec: https://www.php.net/manual/en/info.constants.php
	environment.predefinedConstants["INI_USER"] = values.NewInt(int64(ini.INI_USER))
	environment.predefinedConstants["INI_PERDIR"] = values.NewInt(int64(ini.INI_PERDIR))
	environment.predefinedConstants["INI_SYSTEM"] = values.NewInt(int64(ini.INI_SYSTEM))
	environment.predefinedConstants["INI_ALL"] = values.NewInt(int64(ini.INI_ALL))
}
