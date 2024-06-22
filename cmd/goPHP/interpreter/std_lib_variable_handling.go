package interpreter

import (
	"GoPHP/cmd/goPHP/common"
	"GoPHP/cmd/goPHP/phpError"
	"fmt"
	"math"
	"slices"
	"strings"
)

func registerNativeVariableHandlingFunctions(environment *Environment) {
	environment.nativeFunctions["boolval"] = nativeFn_boolval
	environment.nativeFunctions["doubleval"] = nativeFn_floatval
	environment.nativeFunctions["floatval"] = nativeFn_floatval
	environment.nativeFunctions["get_debug_type"] = nativeFn_get_debug_type
	environment.nativeFunctions["gettype"] = nativeFn_gettype
	environment.nativeFunctions["intval"] = nativeFn_intval
	environment.nativeFunctions["is_array"] = nativeFn_is_array
	environment.nativeFunctions["is_bool"] = nativeFn_is_bool
	environment.nativeFunctions["is_double"] = nativeFn_is_float
	environment.nativeFunctions["is_float"] = nativeFn_is_float
	environment.nativeFunctions["is_int"] = nativeFn_is_int
	environment.nativeFunctions["is_integer"] = nativeFn_is_int
	environment.nativeFunctions["is_long"] = nativeFn_is_int
	environment.nativeFunctions["is_null"] = nativeFn_is_null
	environment.nativeFunctions["is_scalar"] = nativeFn_is_scalar
	environment.nativeFunctions["is_string"] = nativeFn_is_string
	environment.nativeFunctions["print_r"] = nativeFn_print_r
	environment.nativeFunctions["strval"] = nativeFn_strval
	environment.nativeFunctions["var_dump"] = nativeFn_var_dump
	environment.nativeFunctions["var_export"] = nativeFn_var_export
}

// ------------------- MARK: boolval -------------------

func nativeFn_boolval(args []IRuntimeValue, _ *Interpreter) (IRuntimeValue, phpError.Error) {
	args, err := NewFuncParamValidator("boolval").addParam("$value", []string{"mixed"}, nil).validate(args)
	if err != nil {
		return NewVoidRuntimeValue(), err
	}

	boolean, err := lib_boolval(args[0])
	return NewBooleanRuntimeValue(boolean), err
}

func lib_boolval(runtimeValue IRuntimeValue) (bool, phpError.Error) {
	// Spec: https://phplang.org/spec/08-conversions.html#converting-to-boolean-type
	// Spec: https://www.php.net/manual/en/function.boolval.php

	switch runtimeValue.GetType() {
	case ArrayValue:
		// Spec: https://phplang.org/spec/08-conversions.html#converting-to-boolean-type
		// If the source is an array with zero elements, the result value is FALSE; otherwise, the result value is TRUE.
		return len(runtimeValToArrayRuntimeVal(runtimeValue).GetElements()) != 0, nil
	case BooleanValue:
		return runtimeValToBoolRuntimeVal(runtimeValue).GetValue(), nil
	case IntegerValue:
		// Spec: https://phplang.org/spec/08-conversions.html#converting-to-boolean-type
		// If the source type is int or float,
		// then if the source value tests equal to 0, the result value is FALSE; otherwise, the result value is TRUE.
		return runtimeValToIntRuntimeVal(runtimeValue).GetValue() != 0, nil
	case FloatingValue:
		// Spec: https://phplang.org/spec/08-conversions.html#converting-to-boolean-type
		// If the source type is int or float,
		// then if the source value tests equal to 0, the result value is FALSE; otherwise, the result value is TRUE.
		return math.Abs(runtimeValToFloatRuntimeVal(runtimeValue).GetValue()) != 0, nil
	case NullValue:
		// Spec: https://phplang.org/spec/08-conversions.html#converting-to-boolean-type
		// If the source value is NULL, the result value is FALSE.
		return false, nil
	case StringValue:
		// Spec: https://phplang.org/spec/08-conversions.html#converting-to-boolean-type
		// If the source is an empty string or the string “0”, the result value is FALSE; otherwise, the result value is TRUE.
		str := runtimeValToStrRuntimeVal(runtimeValue).GetValue()
		return str != "" && str != "0", nil
	default:
		return false, phpError.NewError("boolval: Unsupported runtime value %s", runtimeValue.GetType())
	}

	// TODO boolval - object
	// Spec: https://phplang.org/spec/08-conversions.html#converting-to-boolean-type
	// If the source is an object, the result value is TRUE.

	// TODO boolval - resource
	// Spec: https://phplang.org/spec/08-conversions.html#converting-to-boolean-type
	// If the source is a resource, the result value is TRUE.
}

// ------------------- MARK: floatval -------------------

func nativeFn_floatval(args []IRuntimeValue, _ *Interpreter) (IRuntimeValue, phpError.Error) {
	args, err := NewFuncParamValidator("floatval").addParam("$value", []string{"mixed"}, nil).validate(args)
	if err != nil {
		return NewVoidRuntimeValue(), err
	}

	floating, err := lib_floatval(args[0])
	return NewFloatingRuntimeValue(floating), err
}

func lib_floatval(runtimeValue IRuntimeValue) (float64, phpError.Error) {
	// Spec: https://phplang.org/spec/08-conversions.html#converting-to-floating-point-type
	// Spec: https://www.php.net/manual/en/function.floatval.php

	switch runtimeValue.GetType() {
	case FloatingValue:
		return runtimeValToFloatRuntimeVal(runtimeValue).GetValue(), nil
	case IntegerValue:
		// Spec: https://phplang.org/spec/08-conversions.html#converting-to-floating-point-type
		// If the source type is int,
		// if the precision can be preserved the result value is the closest approximation to the source value;
		// otherwise, the result is undefined.
		return float64(runtimeValToIntRuntimeVal(runtimeValue).GetValue()), nil
	case StringValue:
		// Spec: https://phplang.org/spec/08-conversions.html#converting-to-floating-point-type
		// If the source is a numeric string or leading-numeric string having integer format,
		// the string’s integer value is treated as described above for a conversion from int.
		// If the source is a numeric string or leading-numeric string having floating-point format,
		// the result value is the closest approximation to the string’s floating-point value.
		// The trailing non-numeric characters in leading-numeric strings are ignored.
		// For any other string, the result value is 0.
		intStr := runtimeValToStrRuntimeVal(runtimeValue).GetValue()
		if common.IsFloatingLiteralWithSign(intStr) {
			return common.FloatingLiteralToFloat64WithSign(intStr), nil
		}
		if common.IsIntegerLiteralWithSign(intStr) {
			intValue, err := common.IntegerLiteralToInt64WithSign(intStr)
			if err != nil {
				return 0, phpError.NewError(err.Error())
			}
			return lib_floatval(NewIntegerRuntimeValue(intValue))
		}
		return 0, nil
	default:
		// Spec: https://phplang.org/spec/08-conversions.html#converting-to-floating-point-type
		// For sources of all other types, the conversion result is obtained by first converting
		// the source value to int and then to float.
		intValue, err := lib_intval(runtimeValue)
		if err != nil {
			return 0, err
		}
		return lib_floatval(NewIntegerRuntimeValue(intValue))
	}

	// TODO lib_floatval - object
	// Spec: https://phplang.org/spec/08-conversions.html#converting-to-floating-point-type
	// If the source is an object, if the class defines a conversion function, the result is determined by that function (this is currently available only to internal classes). If not, the conversion is invalid, the result is assumed to be 1.0 and a non-fatal error is produced.
}

// ------------------- MARK: get_debug_type -------------------

func nativeFn_get_debug_type(args []IRuntimeValue, _ *Interpreter) (IRuntimeValue, phpError.Error) {
	args, err := NewFuncParamValidator("get_debug_type").addParam("$value", []string{"mixed"}, nil).validate(args)
	if err != nil {
		return NewVoidRuntimeValue(), err
	}

	typeStr, err := lib_get_debug_type(args[0])
	return NewStringRuntimeValue(typeStr), err
}

func lib_get_debug_type(runtimeValue IRuntimeValue) (string, phpError.Error) {
	// Spec: https://www.php.net/manual/en/function.get-debug-type

	// TODO lib_get_debug_type - object
	// TODO lib_get_debug_type - resource
	// TODO lib_get_debug_type - resource (closed)
	switch runtimeValue.GetType() {
	case ArrayValue:
		return "array", nil
	case BooleanValue:
		return "bool", nil
	case FloatingValue:
		return "float", nil
	case IntegerValue:
		return "int", nil
	case NullValue:
		return "null", nil
	case StringValue:
		return "string", nil
	default:
		return "unknown type", nil
	}
}

// ------------------- MARK: gettype -------------------

func nativeFn_gettype(args []IRuntimeValue, _ *Interpreter) (IRuntimeValue, phpError.Error) {
	args, err := NewFuncParamValidator("gettype").addParam("$value", []string{"mixed"}, nil).validate(args)
	if err != nil {
		return NewVoidRuntimeValue(), err
	}

	typeStr, err := lib_gettype(args[0])
	return NewStringRuntimeValue(typeStr), err
}

func lib_gettype(runtimeValue IRuntimeValue) (string, phpError.Error) {
	// Spec: https://www.php.net/manual/en/function.gettype.php

	// TODO lib_gettype - object
	// TODO lib_gettype - resource
	// TODO lib_gettype - resource (closed)
	switch runtimeValue.GetType() {
	case ArrayValue:
		return "array", nil
	case BooleanValue:
		return "boolean", nil
	case FloatingValue:
		return "double", nil
	case IntegerValue:
		return "integer", nil
	case NullValue:
		return "NULL", nil
	case StringValue:
		return "string", nil
	default:
		return "unknown type", nil
	}
}

// ------------------- MARK: intval -------------------

func nativeFn_intval(args []IRuntimeValue, _ *Interpreter) (IRuntimeValue, phpError.Error) {
	args, err := NewFuncParamValidator("intval").addParam("$value", []string{"mixed"}, nil).validate(args)
	if err != nil {
		return NewVoidRuntimeValue(), err
	}

	integer, err := lib_intval(args[0])
	return NewIntegerRuntimeValue(integer), err
}

func lib_intval(runtimeValue IRuntimeValue) (int64, phpError.Error) {
	// Spec: https://phplang.org/spec/08-conversions.html#converting-to-integer-type

	switch runtimeValue.GetType() {
	case ArrayValue:
		// Spec: https://phplang.org/spec/08-conversions.html#converting-to-integer-type
		// If the source is an array with zero elements, the result value is 0; otherwise, the result value is 1.
		if len(runtimeValToArrayRuntimeVal(runtimeValue).GetElements()) == 0 {
			return 0, nil
		}
		return 1, nil
	case BooleanValue:
		// Spec: https://phplang.org/spec/08-conversions.html#converting-to-integer-type
		// If the source type is bool, then if the source value is FALSE, the result value is 0; otherwise, the result value is 1.
		if runtimeValToBoolRuntimeVal(runtimeValue).GetValue() {
			return 1, nil
		}
		return 0, nil
	case FloatingValue:
		// Spec: https://phplang.org/spec/08-conversions.html#converting-to-integer-type
		// If the source type is float, for the values INF, -INF, and NAN, the result value is zero.
		// For all other values, if the precision can be preserved (that is, the float is within the range of an integer),
		// the fractional part is rounded towards zero.
		return int64(runtimeValToFloatRuntimeVal(runtimeValue).GetValue()), nil
	case IntegerValue:
		return runtimeValToIntRuntimeVal(runtimeValue).GetValue(), nil
	case NullValue:
		// Spec: https://phplang.org/spec/08-conversions.html#converting-to-integer-type
		// If the source value is NULL, the result value is 0.
		return 0, nil
	case StringValue:
		// Spec: https://phplang.org/spec/08-conversions.html#converting-to-integer-type
		// If the source is a numeric string or leading-numeric string having integer format,
		// if the precision can be preserved the result value is that string’s integer value;
		// otherwise, the result is undefined.
		// If the source is a numeric string or leading-numeric string having floating-point format,
		// the string’s floating-point value is treated as described above for a conversion from float.
		// The trailing non-numeric characters in leading-numeric strings are ignored.
		// For any other string, the result value is 0.
		intStr := runtimeValToStrRuntimeVal(runtimeValue).GetValue()
		if common.IsFloatingLiteralWithSign(intStr) {
			return lib_intval(NewFloatingRuntimeValue(common.FloatingLiteralToFloat64WithSign(intStr)))
		}
		if common.IsIntegerLiteralWithSign(intStr) {
			intValue, err := common.IntegerLiteralToInt64WithSign(intStr)
			if err != nil {
				return 0, phpError.NewError(err.Error())
			}
			return intValue, nil
		}
		return 0, nil
	default:
		return 0, phpError.NewError("lib_intval: Unsupported runtime value %s", runtimeValue.GetType())
	}

	// TODO lib_intval - object
	// Spec: https://phplang.org/spec/08-conversions.html#converting-to-integer-type
	// If the source is an object, if the class defines a conversion function, the result is determined by that function (this is currently available only to internal classes). If not, the conversion is invalid, the result is assumed to be 1 and a non-fatal error is produced.

	// TODO lib_intval - resource
	// Spec: https://phplang.org/spec/08-conversions.html#converting-to-integer-type
	// If the source is a resource, the result is the resource’s unique ID.
}

// ------------------- MARK: is_array -------------------

func nativeFn_is_array(args []IRuntimeValue, _ *Interpreter) (IRuntimeValue, phpError.Error) {
	args, err := NewFuncParamValidator("is_array").addParam("$value", []string{"mixed"}, nil).validate(args)
	if err != nil {
		return NewVoidRuntimeValue(), err
	}

	return NewBooleanRuntimeValue(lib_is_array(args[0])), nil
}

func lib_is_array(runtimeValue IRuntimeValue) bool {
	// Spec: https://www.php.net/manual/en/function.is-array
	return runtimeValue.GetType() == ArrayValue
}

// ------------------- MARK: is_bool -------------------

func nativeFn_is_bool(args []IRuntimeValue, _ *Interpreter) (IRuntimeValue, phpError.Error) {
	args, err := NewFuncParamValidator("is_bool").addParam("$value", []string{"mixed"}, nil).validate(args)
	if err != nil {
		return NewVoidRuntimeValue(), err
	}

	return NewBooleanRuntimeValue(lib_is_bool(args[0])), nil
}

func lib_is_bool(runtimeValue IRuntimeValue) bool {
	// Spec: https://www.php.net/manual/en/function.is-bool
	return runtimeValue.GetType() == BooleanValue
}

// ------------------- MARK: is_float -------------------

func nativeFn_is_float(args []IRuntimeValue, _ *Interpreter) (IRuntimeValue, phpError.Error) {
	args, err := NewFuncParamValidator("is_float").addParam("$value", []string{"mixed"}, nil).validate(args)
	if err != nil {
		return NewVoidRuntimeValue(), err
	}

	return NewBooleanRuntimeValue(lib_is_float(args[0])), nil
}

func lib_is_float(runtimeValue IRuntimeValue) bool {
	// Spec: https://www.php.net/manual/en/function.is-float
	return runtimeValue.GetType() == FloatingValue
}

// ------------------- MARK: is_int -------------------

func nativeFn_is_int(args []IRuntimeValue, _ *Interpreter) (IRuntimeValue, phpError.Error) {
	args, err := NewFuncParamValidator("is_int").addParam("$value", []string{"mixed"}, nil).validate(args)
	if err != nil {
		return NewVoidRuntimeValue(), err
	}

	return NewBooleanRuntimeValue(lib_is_int(args[0])), nil
}

func lib_is_int(runtimeValue IRuntimeValue) bool {
	// Spec: https://www.php.net/manual/en/function.is-int
	return runtimeValue.GetType() == IntegerValue
}

// ------------------- MARK: is_null -------------------

func nativeFn_is_null(args []IRuntimeValue, _ *Interpreter) (IRuntimeValue, phpError.Error) {
	args, err := NewFuncParamValidator("is_null").addParam("$value", []string{"mixed"}, nil).validate(args)
	if err != nil {
		return NewVoidRuntimeValue(), err
	}

	return NewBooleanRuntimeValue(lib_is_null(args[0])), nil
}

func lib_is_null(runtimeValue IRuntimeValue) bool {
	// Spec: https://www.php.net/manual/en/function.is-null.php
	return runtimeValue.GetType() == NullValue
}

// ------------------- MARK: is_scalar -------------------

func nativeFn_is_scalar(args []IRuntimeValue, _ *Interpreter) (IRuntimeValue, phpError.Error) {
	args, err := NewFuncParamValidator("is_scalar").addParam("$value", []string{"mixed"}, nil).validate(args)
	if err != nil {
		return NewVoidRuntimeValue(), err
	}

	return NewBooleanRuntimeValue(lib_is_scalar(args[0])), nil
}

func lib_is_scalar(runtimeValue IRuntimeValue) bool {
	// Spec: https://www.php.net/manual/en/function.is-scalar.php
	return slices.Contains([]ValueType{BooleanValue, IntegerValue, FloatingValue, StringValue}, runtimeValue.GetType())
}

// ------------------- MARK: is_string -------------------

func nativeFn_is_string(args []IRuntimeValue, _ *Interpreter) (IRuntimeValue, phpError.Error) {
	args, err := NewFuncParamValidator("is_string").addParam("$value", []string{"mixed"}, nil).validate(args)
	if err != nil {
		return NewVoidRuntimeValue(), err
	}

	return NewBooleanRuntimeValue(lib_is_string(args[0])), nil
}

func lib_is_string(runtimeValue IRuntimeValue) bool {
	// Spec: https://www.php.net/manual/en/function.is-string
	return runtimeValue.GetType() == StringValue
}

// ------------------- MARK: print_r -------------------

func nativeFn_print_r(args []IRuntimeValue, interpreter *Interpreter) (IRuntimeValue, phpError.Error) {
	args, err := NewFuncParamValidator("print_r").
		addParam("$value", []string{"mixed"}, nil).
		addParam("$return", []string{"bool"}, NewBooleanRuntimeValue(false)).
		validate(args)
	if err != nil {
		return NewVoidRuntimeValue(), err
	}

	return lib_print_r(args[0], runtimeValToBoolRuntimeVal(args[1]).GetValue(), interpreter)
}

func lib_print_r(value IRuntimeValue, returnValue bool, interpreter *Interpreter) (IRuntimeValue, phpError.Error) {
	result, err := lib_print_r_var(value, 4)
	if err != nil {
		return NewVoidRuntimeValue(), err
	}

	if returnValue {
		return NewStringRuntimeValue(result), nil
	} else {
		interpreter.print(result)
		return NewBooleanRuntimeValue(true), nil
	}
}

func lib_print_r_var(value IRuntimeValue, depth int) (string, phpError.Error) {
	result := ""
	var err phpError.Error
	switch value.GetType() {
	case ArrayValue:
		keys := runtimeValToArrayRuntimeVal(value).GetKeys()
		elements := runtimeValToArrayRuntimeVal(value).GetElements()
		result = fmt.Sprintf("Array\n%s(\n", strings.Repeat(" ", depth-4))
		for _, key := range keys {
			keyStr, err := lib_print_r_var(key, depth+4)
			if err != nil {
				return "", err
			}
			valueStr, err := lib_print_r_var(elements[key], depth+8)
			if err != nil {
				return "", err
			}

			result += fmt.Sprintf("%s[%s] => %s\n", strings.Repeat(" ", depth), keyStr, valueStr)
		}
		result += fmt.Sprintf("%s)", strings.Repeat(" ", depth-4))
	case BooleanValue:
		if runtimeValToBoolRuntimeVal(value).GetValue() {
			result = "1"
		} else {
			result = ""
		}
	case FloatingValue, IntegerValue:
		result, err = lib_strval(value)
		if err != nil {
			return "", err
		}
	case NullValue:
		result = ""
	case StringValue:
		result = runtimeValToStrRuntimeVal(value).GetValue()
	default:
		return "", phpError.NewError("lib_print_r_var: Unsupported runtime value %s", value.GetType())
	}
	return result, nil
}

// ------------------- MARK: strval -------------------

func nativeFn_strval(args []IRuntimeValue, _ *Interpreter) (IRuntimeValue, phpError.Error) {
	args, err := NewFuncParamValidator("strval").addParam("$value", []string{"mixed"}, nil).validate(args)
	if err != nil {
		return NewVoidRuntimeValue(), err
	}

	str, err := lib_strval(args[0])
	return NewStringRuntimeValue(str), err
}

func lib_strval(runtimeValue IRuntimeValue) (string, phpError.Error) {
	// Spec: https://phplang.org/spec/08-conversions.html#converting-to-string-type

	switch runtimeValue.GetType() {
	case ArrayValue:
		// Spec: https://phplang.org/spec/08-conversions.html#converting-to-string-type
		// If the source is an array, the conversion is invalid. The result value is the string “Array” and a non-fatal error is produced.
		return "Array", nil
		// TODO lib_strval - array non-fatal error: "Warning: Array to string conversion in /home/user/scripts/code.php on line 2" - only if E_ALL | E_WARNING
	case BooleanValue:
		// Spec: https://phplang.org/spec/08-conversions.html#converting-to-string-type
		// If the source type is bool, then if the source value is FALSE, the result value is the empty string;
		// otherwise, the result value is “1”.
		if !runtimeValToBoolRuntimeVal(runtimeValue).GetValue() {
			return "", nil
		}
		return "1", nil
	case FloatingValue:
		// Spec: https://phplang.org/spec/08-conversions.html#converting-to-string-type
		// If the source type is int or float, then the result value is a string containing the textual representation
		// of the source value (as specified by the library function sprintf).
		return fmt.Sprintf("%g", runtimeValToFloatRuntimeVal(runtimeValue).GetValue()), nil
	case IntegerValue:
		// Spec: https://phplang.org/spec/08-conversions.html#converting-to-string-type
		// If the source type is int or float, then the result value is a string containing the textual representation
		// of the source value (as specified by the library function sprintf).
		return fmt.Sprintf("%d", runtimeValToIntRuntimeVal(runtimeValue).GetValue()), nil
	case NullValue:
		// Spec: https://phplang.org/spec/08-conversions.html#converting-to-string-type
		// If the source value is NULL, the result value is the empty string.
		return "", nil
	case StringValue:
		return runtimeValToStrRuntimeVal(runtimeValue).GetValue(), nil
	default:
		return "", phpError.NewError("lib_strval: Unsupported runtime value %s", runtimeValue.GetType())
	}

	// TODO lib_strval - object
	// Spec: https://phplang.org/spec/08-conversions.html#converting-to-string-type
	// If the source is an object, then if that object’s class has a __toString method, the result value is the string returned by that method; otherwise, the conversion is invalid and a fatal error is produced.

	// TODO lib_strval - resource
	// Spec: https://phplang.org/spec/08-conversions.html#converting-to-string-type
	// If the source is a resource, the result value is an implementation-defined string.
}

// ------------------- MARK: var_dump -------------------

func nativeFn_var_dump(args []IRuntimeValue, interpreter *Interpreter) (IRuntimeValue, phpError.Error) {
	// Spec: https://www.php.net/manual/en/function.var-dump

	args, err := NewFuncParamValidator("var_dump").
		addParam("$value", []string{"mixed"}, nil).addVariableLenParam("$values", []string{"mixed"}).
		validate(args)
	if err != nil {
		return NewVoidRuntimeValue(), err
	}

	if err := lib_var_dump_var(interpreter, args[0], 2); err != nil {
		return NewVoidRuntimeValue(), err
	}

	if len(args) == 2 {
		values := runtimeValToArrayRuntimeVal(args[1])
		for _, key := range values.GetKeys() {
			argValue, _ := values.GetElement(key)
			if err := lib_var_dump_var(interpreter, argValue, 2); err != nil {
				return NewVoidRuntimeValue(), err
			}
		}
	}

	return NewVoidRuntimeValue(), nil
}

func lib_var_dump_var(interpreter *Interpreter, value IRuntimeValue, depth int) phpError.Error {
	switch value.GetType() {
	case ArrayValue:
		keys := runtimeValToArrayRuntimeVal(value).GetKeys()
		elements := runtimeValToArrayRuntimeVal(value).GetElements()
		interpreter.println(fmt.Sprintf("array(%d) {", len(keys)))
		for _, key := range keys {
			switch key.GetType() {
			case IntegerValue:
				keyValue := runtimeValToIntRuntimeVal(key).GetValue()
				interpreter.println(fmt.Sprintf("%s[%d]=>", strings.Repeat(" ", depth), keyValue))
			case StringValue:
				keyValue := runtimeValToStrRuntimeVal(key).GetValue()
				interpreter.println(fmt.Sprintf(`%s["%s"]=>`, strings.Repeat(" ", depth), keyValue))
			default:
				return phpError.NewError("lib_var_dump_var: Unsupported array key type %s", key.GetType())
			}
			interpreter.print(strings.Repeat(" ", depth))
			if err := lib_var_dump_var(interpreter, elements[key], depth+2); err != nil {
				return err
			}
		}
		interpreter.println(strings.Repeat(" ", depth-2) + "}")
	case BooleanValue:
		if runtimeValToBoolRuntimeVal(value).GetValue() {
			interpreter.println("bool(true)")
		} else {
			interpreter.println("bool(false)")
		}
	case FloatingValue:
		strVal, err := lib_strval(value)
		if err != nil {
			return err
		}
		interpreter.println("float(" + strVal + ")")
	case IntegerValue:
		strVal, err := lib_strval(value)
		if err != nil {
			return err
		}
		interpreter.println("int(" + strVal + ")")
	case NullValue:
		interpreter.println("NULL")
	case StringValue:
		strVal := runtimeValToStrRuntimeVal(value).GetValue()
		interpreter.println(fmt.Sprintf("string(%d) \"%s\"", len(strVal), strVal))
	default:
		return phpError.NewError("lib_var_dump_var: Unsupported runtime value %s", value.GetType())
	}
	return nil
}

// ------------------- MARK: var_export -------------------

func nativeFn_var_export(args []IRuntimeValue, interpreter *Interpreter) (IRuntimeValue, phpError.Error) {
	// Spec: https://www.php.net/manual/en/function.var-export

	args, err := NewFuncParamValidator("var_dump").
		addParam("$value", []string{"mixed"}, nil).
		addParam("$return", []string{"bool"}, NewBooleanRuntimeValue(false)).
		validate(args)
	if err != nil {
		return NewVoidRuntimeValue(), err
	}

	return lib_var_export(args[0], runtimeValToBoolRuntimeVal(args[1]).GetValue(), interpreter)
}

func lib_var_export(value IRuntimeValue, returnValue bool, interpreter *Interpreter) (IRuntimeValue, phpError.Error) {
	result, err := lib_var_export_var(value, 2)
	if err != nil {
		return NewVoidRuntimeValue(), err
	}

	if returnValue {
		return NewStringRuntimeValue(result), nil
	} else {
		interpreter.print(result)
		return NewNullRuntimeValue(), nil
	}
}

func lib_var_export_var(value IRuntimeValue, depth int) (string, phpError.Error) {
	result := ""
	var err phpError.Error
	switch value.GetType() {
	case ArrayValue:
		keys := runtimeValToArrayRuntimeVal(value).GetKeys()
		elements := runtimeValToArrayRuntimeVal(value).GetElements()
		result = fmt.Sprintf("%sarray (\n", strings.Repeat(" ", depth-2))
		for _, key := range keys {
			keyStr, err := lib_var_export_var(key, depth+2)
			if err != nil {
				return "", err
			}
			valueStr, err := lib_var_export_var(elements[key], depth+2)
			if err != nil {
				return "", err
			}

			if elements[key].GetType() == ArrayValue {
				result += fmt.Sprintf("%s%s => \n%s%s,\n",
					strings.Repeat(" ", depth), keyStr,
					strings.Repeat(" ", depth-2), valueStr,
				)
			} else {
				result += fmt.Sprintf("%s%s => %s,\n", strings.Repeat(" ", depth), keyStr, valueStr)
			}
		}
		result += fmt.Sprintf("%s)", strings.Repeat(" ", depth-2))
	case BooleanValue:
		if runtimeValToBoolRuntimeVal(value).GetValue() {
			result = "true"
		} else {
			result = "false"
		}
	case FloatingValue, IntegerValue:
		result, err = lib_strval(value)
		if err != nil {
			return "", err
		}
	case NullValue:
		result = "NULL"
	case StringValue:
		result = "'" + runtimeValToStrRuntimeVal(value).GetValue() + "'"
	default:
		return "", phpError.NewError("lib_var_export: Unsupported runtime value %s", value.GetType())
	}
	return result, nil
}

// TODO debug_​zval_​dump
// TODO get_​defined_​vars
// TODO get_​resource_​id
// TODO get_​resource_​type
// TODO is_​callable
// TODO is_​countable
// TODO is_​iterable
// TODO is_​numeric
// TODO is_​object
// TODO is_​resource
// TODO serialize
// TODO settype
// TODO unserialize
