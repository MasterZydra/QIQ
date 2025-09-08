package funcParamValidator

import (
	"QIQ/cmd/qiq/phpError"
	"QIQ/cmd/qiq/runtime/values"
	"slices"
	"strings"
)

type funcParam struct {
	name          string
	paramType     []string
	isVariableLen bool
	defaultValue  values.RuntimeValue
}

type Validator struct {
	funcName string
	params   []funcParam
}

func NewValidator(funcName string) *Validator {
	return &Validator{funcName: funcName, params: []funcParam{}}
}

// Add parameter value
func (validator *Validator) AddParam(name string, paramType []string, defaultValue values.RuntimeValue) *Validator {
	// Add type of default value to allowed types
	// e.g. function test(int $a = null) => Allowed types: int|null
	if defaultValue != nil {
		defaultValueType := values.ToPhpType(defaultValue)
		if defaultValueType == "" && !slices.Contains(paramType, defaultValueType) && !slices.Contains(paramType, "mixed") {
			paramType = append(paramType, defaultValueType)
		}
	}

	validator.params = append(validator.params, funcParam{name: name, paramType: paramType, defaultValue: defaultValue})
	return validator
}

// Add parameter with variable length (e.g. "mixed ...$args")
func (validator *Validator) AddVariableLenParam(name string, paramType []string) *Validator {
	validator.params = append(validator.params, funcParam{
		name: name, paramType: paramType, isVariableLen: true, defaultValue: values.NewArray(),
	})
	return validator
}

// Validate the given arguments
func (validator *Validator) Validate(args []values.RuntimeValue) ([]values.RuntimeValue, phpError.Error) {
	typeMatches := func(param funcParam, arg values.RuntimeValue) bool {
		typeStr := values.ToPhpType(arg)
		if typeStr == "" {
			return false
		}
		if typeStr == "NULL" {
			typeStr = "null"
		}
		return len(param.paramType) == 0 || slices.Contains(param.paramType, "mixed") || slices.Contains(param.paramType, typeStr)
	}

	lastArgIndex := 0
	allArgsValidated := false
	validatedArgs := []values.RuntimeValue{}
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

			typeStr := values.ToPhpType(arg)
			if typeStr == "" {
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
		varLenArg := values.NewArray()
		for argIndex < len(args) {
			arg := args[argIndex]
			if typeMatches(param, arg) {
				varLenArg.SetElement(values.NewInt(int64(argIndex-paramIndex)), arg)
				argIndex++
				continue
			}

			typeStr := values.ToPhpType(arg)
			if typeStr == "" {
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
	if (!allArgsValidated && len(args) > len(validator.params)) || (allArgsValidated && len(args) > lastArgIndex+1) {
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

func (validator *Validator) getLeastExpectedParams() int {
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
