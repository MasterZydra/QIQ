package interpreter

import (
	"GoPHP/cmd/goPHP/common"
	"GoPHP/cmd/goPHP/common/os"
	"GoPHP/cmd/goPHP/config"
	"GoPHP/cmd/goPHP/request"
	"GoPHP/cmd/goPHP/runtime"
	"GoPHP/cmd/goPHP/runtime/values"
	"fmt"
	"math"
	"regexp"
	"strings"
)

func registerPredefinedVariables(environment *Environment, request *request.Request, interpreter runtime.Interpreter) {
	if interpreter == nil {
		registerPredefinedVariableEnv(environment, request, true)
		registerPredefinedVariableGet(environment, request, interpreter, true)
		registerPredefinedVariablePost(environment, request, interpreter, true)
		registerPredefinedVariableCookie(environment, request, true)
		registerPredefinedVariableServer(environment, request, interpreter, true)
		return
	}

	variables_order := interpreter.GetIni().GetStr("variables_order")
	for _, variable := range variables_order {
		switch string(variable) {
		case "E":
			registerPredefinedVariableEnv(environment, request, true)
		case "G":
			registerPredefinedVariableGet(environment, request, interpreter, true)
		case "P":
			registerPredefinedVariablePost(environment, request, interpreter, true)
		case "C":
			registerPredefinedVariableCookie(environment, request, true)
		case "S":
			registerPredefinedVariableServer(environment, request, interpreter, true)
		}
	}

	if !strings.Contains(variables_order, "E") {
		registerPredefinedVariableEnv(environment, request, false)
	}
	if !strings.Contains(variables_order, "G") {
		registerPredefinedVariableGet(environment, request, interpreter, false)
	}
	if !strings.Contains(variables_order, "P") {
		registerPredefinedVariablePost(environment, request, interpreter, false)
	}
	if !strings.Contains(variables_order, "C") {
		registerPredefinedVariableCookie(environment, request, false)
	}
	if !strings.Contains(variables_order, "S") {
		registerPredefinedVariableServer(environment, request, interpreter, false)
	}

	requestVar := values.NewArray()
	mergeArrays(requestVar, environment.predefinedVariables["$_GET"].(*values.Array))
	mergeArrays(requestVar, environment.predefinedVariables["$_POST"].(*values.Array))
	mergeArrays(requestVar, environment.predefinedVariables["$_COOKIE"].(*values.Array))
	environment.predefinedVariables["$_REQUEST"] = requestVar
}

// TODO Replace with std lib func array_merge
func mergeArrays(a, b *values.Array) {
	for _, key := range b.Keys {
		value, _ := b.GetElement(key)
		a.SetElement(key, values.DeepCopy(value))
	}
}

func registerPredefinedVariableEnv(environment *Environment, request *request.Request, init bool) {
	if init {
		environment.predefinedVariables["$_ENV"] = stringMapToArray(request.Env)
	} else {
		environment.predefinedVariables["$_ENV"] = values.NewArray()
	}
}

func registerPredefinedVariableCookie(environment *Environment, request *request.Request, init bool) {
	if init {
		environment.predefinedVariables["$_COOKIE"] = parseCookies(request.Cookie)
	} else {
		environment.predefinedVariables["$_COOKIE"] = values.NewArray()
	}
}

func registerPredefinedVariableGet(environment *Environment, request *request.Request, interpreter runtime.Interpreter, init bool) {
	if init {
		array, err := parseQuery(request.QueryString, interpreter)
		if err != nil {
			println(err.Error())
		}
		environment.predefinedVariables["$_GET"] = array
	} else {
		environment.predefinedVariables["$_GET"] = values.NewArray()
	}
}

func registerPredefinedVariablePost(environment *Environment, request *request.Request, interpreter runtime.Interpreter, init bool) {
	if init {
		post, file, err := parsePost(request.Post, interpreter)
		if err != nil {
			println(err.Error())
		}
		environment.predefinedVariables["$_POST"] = post
		environment.predefinedVariables["$_FILES"] = file
	} else {
		environment.predefinedVariables["$_POST"] = values.NewArray()
		environment.predefinedVariables["$_FILES"] = values.NewArray()
	}
}

func registerPredefinedVariableServer(environment *Environment, request *request.Request, interpreter runtime.Interpreter, init bool) {
	environment.predefinedVariables["$_SERVER"] = values.NewArray()
	if init {
		server := environment.predefinedVariables["$_SERVER"].(*values.Array)
		if len(request.Args) > 0 {
			server.SetElement(values.NewStr("argc"), values.NewInt(int64(len(request.Args))))
			server.SetElement(values.NewStr("argv"), paramToArray(request.Args, interpreter))
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

func paramToArray(params [][]string, interpreter runtime.Interpreter) *values.Array {
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
				if common.IsIntegerLiteral(key, false) {
					intValue, _ := common.IntegerLiteralToInt64(key, false)
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
		env := NewEnvironment(nil, request.NewRequest(), interpreter)
		env.declareVariable("$"+paramName, arrayValue)

		// Execute PHP to store new array values in env
		NewInterpreter(interpreter.GetIni(), request.NewRequest(), "").process(fmt.Sprintf(`<?php $%s = "%s";`, key, value), env, true)

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
	environment.predefinedConstants["DIRECTORY_SEPARATOR"] = values.NewStr(os.DIR_SEP)
	environment.predefinedConstants["FALSE"] = values.NewBool(false)
	environment.predefinedConstants["TRUE"] = values.NewBool(true)
	environment.predefinedConstants["NULL"] = values.NewNull()
	environment.predefinedConstants["PHP_INT_MAX"] = values.NewInt(math.MaxInt64)
	environment.predefinedConstants["PHP_INT_MIN"] = values.NewInt(math.MinInt64)
	environment.predefinedConstants["PHP_INT_SIZE"] = values.NewInt(64 / 8)
	environment.predefinedConstants["PHP_OS"] = values.NewStr(os.Os())
	environment.predefinedConstants["PHP_OS_FAMILY"] = values.NewStr(os.OS_FAMILY)
	environment.predefinedConstants["PHP_EOL"] = values.NewStr(os.EOL)
	environment.predefinedConstants["PHP_VERSION"] = values.NewStr(config.Version)
	environment.predefinedConstants["PHP_MAJOR_VERSION"] = values.NewInt(config.MajorVersion)
	environment.predefinedConstants["PHP_MINOR_VERSION"] = values.NewInt(config.MinorVersion)
	environment.predefinedConstants["PHP_RELEASE_VERSION"] = values.NewInt(config.ReleaseVersion)
	environment.predefinedConstants["PHP_EXTRA_VERSION"] = values.NewStr(config.ExtraVersion)
	environment.predefinedConstants["PHP_VERSION_ID"] = values.NewInt(config.VersionId)
}
