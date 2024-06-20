package interpreter

import (
	"GoPHP/cmd/goPHP/phpError"
	"slices"
	"strings"
)

type funcParam struct {
	name          string
	paramType     []string
	isVariableLen bool
	defaultValue  IRuntimeValue
}

type funcParamValidator struct {
	funcName string
	params   []funcParam
}

func NewFuncParamValidator(funcName string) *funcParamValidator {
	return &funcParamValidator{funcName: funcName, params: []funcParam{}}
}

// Add parameter value
func (validator *funcParamValidator) addParam(name string, paramType []string, defaultValue IRuntimeValue) *funcParamValidator {
	// Add type of default value to allowed types
	// e.g. function test(int $a = null) => Allowed types: int|null
	if defaultValue != nil {
		defaultValueType, found := paramTypeRuntimeValue[defaultValue.GetType()]
		if found && !slices.Contains(paramType, defaultValueType) && !slices.Contains(paramType, "mixed") {
			paramType = append(paramType, defaultValueType)
		}
	}

	validator.params = append(validator.params, funcParam{name: name, paramType: paramType, defaultValue: defaultValue})
	return validator
}

// Add parameter with variable length (e.g. "mixed ...$args")
func (validator *funcParamValidator) addVariableLenParam(name string, paramType []string) *funcParamValidator {
	validator.params = append(validator.params, funcParam{
		name: name, paramType: paramType, isVariableLen: true, defaultValue: NewArrayRuntimeValue(),
	})
	return validator
}

// Validate the given arguments
func (validator *funcParamValidator) validate(args []IRuntimeValue) ([]IRuntimeValue, phpError.Error) {

	typeMatches := func(param funcParam, arg IRuntimeValue) bool {
		typeStr, found := paramTypeRuntimeValue[arg.GetType()]
		if !found {
			return false
		}
		return slices.Contains(param.paramType, "mixed") || slices.Contains(param.paramType, typeStr)
	}

	lastArgIndex := 0
	allArgsValidated := false
	validatedArgs := []IRuntimeValue{}
	for paramIndex, param := range validator.params {
		lastArgIndex = paramIndex
		if paramIndex >= len(args) {
			allArgsValidated = true
		}

		if allArgsValidated {
			if param.defaultValue == nil {
				return args, phpError.NewError(
					"Uncaught ArgumentCountError: Too few arguments to function %s(), %d passed and at least %d expected",
					validator.funcName, len(args), validator.getLeastExpectedParams(),
				)
			}
			validatedArgs = append(validatedArgs, param.defaultValue)
			continue
		}

		arg := args[paramIndex]
		if !param.isVariableLen {
			if typeMatches(param, arg) {
				validatedArgs = append(validatedArgs, arg)
				continue
			}

			typeStr, found := paramTypeRuntimeValue[arg.GetType()]
			if !found {
				return args, phpError.NewError("validate: No mapping for type %s", arg.GetType())
			}

			// Type mismatch
			return args, phpError.NewError(
				"Uncaught TypeError: %s(): Argument #%d (%s) must be of type %s, %s given",
				validator.funcName, paramIndex+1, param.name,
				strings.Join(param.paramType, "|"), typeStr,
			)
		}

		// Variable length parameter
		argIndex := paramIndex
		varLenArg := NewArrayRuntimeValue()
		for argIndex < len(args) {
			arg := args[argIndex]
			if typeMatches(param, arg) {
				varLenArg.SetElement(NewIntegerRuntimeValue(int64(argIndex-paramIndex)), arg)
				argIndex++
				continue
			}

			typeStr, found := paramTypeRuntimeValue[arg.GetType()]
			if !found {
				return args, phpError.NewError("validate: No mapping for type %s", arg.GetType())
			}

			// Type mismatch
			return args, phpError.NewError(
				"Uncaught TypeError: %s(): Argument #%d (%s) must be of type %s, %s given",
				validator.funcName, paramIndex+1, param.name,
				strings.Join(param.paramType, "|"), typeStr,
			)
		}
		validatedArgs = append(validatedArgs, varLenArg)
		return validatedArgs, nil
	}

	// Too many arguments given
	if len(args) > lastArgIndex+1 {
		if len(validator.params) > 0 && validator.params[len(validator.params)-1].defaultValue != nil {
			// Optional arguments at the end
			return args, phpError.NewError(
				"Uncaught ArgumentCountError: %s() expects most %d argument, %d given",
				validator.funcName, len(validator.params), len(args),
			)
		} else {
			// No optional arguments
			return args, phpError.NewError(
				"Uncaught ArgumentCountError: %s() expects exactly %d argument, %d given",
				validator.funcName, len(validator.params), len(args),
			)
		}
	}
	return validatedArgs, nil
}

func (validator *funcParamValidator) getLeastExpectedParams() int {
	leastParams := 0
	for _, param := range validator.params {
		if param.defaultValue == nil {
			leastParams++
			continue
		}
		return leastParams
	}
	return leastParams
}
