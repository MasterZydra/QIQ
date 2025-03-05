package interpreter

import (
	"GoPHP/cmd/goPHP/common"
	"GoPHP/cmd/goPHP/phpError"
	"GoPHP/cmd/goPHP/runtime"
	"GoPHP/cmd/goPHP/runtime/values"
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

func nativeFn_boolval(args []values.RuntimeValue, _ runtime.Context) (values.RuntimeValue, phpError.Error) {
	args, err := NewFuncParamValidator("boolval").addParam("$value", []string{"mixed"}, nil).validate(args)
	if err != nil {
		return values.NewVoid(), err
	}

	boolean, err := lib_boolval(args[0])
	return values.NewBool(boolean), err
}

func lib_boolval(runtimeValue values.RuntimeValue) (bool, phpError.Error) {
	// Spec: https://phplang.org/spec/08-conversions.html#converting-to-boolean-type
	// Spec: https://www.php.net/manual/en/function.boolval.php

	switch runtimeValue.GetType() {
	case values.ArrayValue:
		// Spec: https://phplang.org/spec/08-conversions.html#converting-to-boolean-type
		// If the source is an array with zero elements, the result value is FALSE; otherwise, the result value is TRUE.
		return len(runtimeValue.(*values.Array).Elements) != 0, nil
	case values.BoolValue:
		return runtimeValue.(*values.Bool).Value, nil
	case values.IntValue:
		// Spec: https://phplang.org/spec/08-conversions.html#converting-to-boolean-type
		// If the source type is int or float,
		// then if the source value tests equal to 0, the result value is FALSE; otherwise, the result value is TRUE.
		return runtimeValue.(*values.Int).Value != 0, nil
	case values.FloatValue:
		// Spec: https://phplang.org/spec/08-conversions.html#converting-to-boolean-type
		// If the source type is int or float,
		// then if the source value tests equal to 0, the result value is FALSE; otherwise, the result value is TRUE.
		return math.Abs(runtimeValue.(*values.Float).Value) != 0, nil
	case values.NullValue:
		// Spec: https://phplang.org/spec/08-conversions.html#converting-to-boolean-type
		// If the source value is NULL, the result value is FALSE.
		return false, nil
	case values.StrValue:
		// Spec: https://phplang.org/spec/08-conversions.html#converting-to-boolean-type
		// If the source is an empty string or the string “0”, the result value is FALSE; otherwise, the result value is TRUE.
		str := runtimeValue.(*values.Str).Value
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

func nativeFn_floatval(args []values.RuntimeValue, _ runtime.Context) (values.RuntimeValue, phpError.Error) {
	args, err := NewFuncParamValidator("floatval").addParam("$value", []string{"mixed"}, nil).validate(args)
	if err != nil {
		return values.NewVoid(), err
	}

	floating, err := lib_floatval(args[0])
	return values.NewFloat(floating), err
}

func lib_floatval(runtimeValue values.RuntimeValue) (float64, phpError.Error) {
	// Spec: https://phplang.org/spec/08-conversions.html#converting-to-floating-point-type
	// Spec: https://www.php.net/manual/en/function.floatval.php

	switch runtimeValue.GetType() {
	case values.FloatValue:
		return runtimeValue.(*values.Float).Value, nil
	case values.IntValue:
		// Spec: https://phplang.org/spec/08-conversions.html#converting-to-floating-point-type
		// If the source type is int,
		// if the precision can be preserved the result value is the closest approximation to the source value;
		// otherwise, the result is undefined.
		return float64(runtimeValue.(*values.Int).Value), nil
	case values.StrValue:
		// Spec: https://phplang.org/spec/08-conversions.html#converting-to-floating-point-type
		// If the source is a numeric string or leading-numeric string having integer format,
		// the string’s integer value is treated as described above for a conversion from int.
		// If the source is a numeric string or leading-numeric string having floating-point format,
		// the result value is the closest approximation to the string’s floating-point value.
		// The trailing non-numeric characters in leading-numeric strings are ignored.
		// For any other string, the result value is 0.
		intStr := runtimeValue.(*values.Str).Value
		if common.IsFloatingLiteralWithSign(intStr) {
			return common.FloatingLiteralToFloat64WithSign(intStr), nil
		}
		if common.IsIntegerLiteralWithSign(intStr) {
			intValue, err := common.IntegerLiteralToInt64WithSign(intStr)
			if err != nil {
				return 0, phpError.NewError(err.Error())
			}
			return lib_floatval(values.NewInt(intValue))
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
		return lib_floatval(values.NewInt(intValue))
	}

	// TODO lib_floatval - object
	// Spec: https://phplang.org/spec/08-conversions.html#converting-to-floating-point-type
	// If the source is an object, if the class defines a conversion function, the result is determined by that function (this is currently available only to internal classes). If not, the conversion is invalid, the result is assumed to be 1.0 and a non-fatal error is produced.
}

// ------------------- MARK: get_debug_type -------------------

func nativeFn_get_debug_type(args []values.RuntimeValue, _ runtime.Context) (values.RuntimeValue, phpError.Error) {
	args, err := NewFuncParamValidator("get_debug_type").addParam("$value", []string{"mixed"}, nil).validate(args)
	if err != nil {
		return values.NewVoid(), err
	}

	typeStr, err := lib_get_debug_type(args[0])
	return values.NewStr(typeStr), err
}

func lib_get_debug_type(runtimeValue values.RuntimeValue) (string, phpError.Error) {
	// Spec: https://www.php.net/manual/en/function.get-debug-type

	// TODO lib_get_debug_type - object
	// TODO lib_get_debug_type - resource
	// TODO lib_get_debug_type - resource (closed)
	switch runtimeValue.GetType() {
	case values.ArrayValue:
		return "array", nil
	case values.BoolValue:
		return "bool", nil
	case values.FloatValue:
		return "float", nil
	case values.IntValue:
		return "int", nil
	case values.NullValue:
		return "null", nil
	case values.StrValue:
		return "string", nil
	default:
		return "unknown type", nil
	}
}

// ------------------- MARK: gettype -------------------

func nativeFn_gettype(args []values.RuntimeValue, _ runtime.Context) (values.RuntimeValue, phpError.Error) {
	args, err := NewFuncParamValidator("gettype").addParam("$value", []string{"mixed"}, nil).validate(args)
	if err != nil {
		return values.NewVoid(), err
	}

	typeStr, err := lib_gettype(args[0])
	return values.NewStr(typeStr), err
}

func lib_gettype(runtimeValue values.RuntimeValue) (string, phpError.Error) {
	// Spec: https://www.php.net/manual/en/function.gettype.php

	// TODO lib_gettype - object
	// TODO lib_gettype - resource
	// TODO lib_gettype - resource (closed)
	switch runtimeValue.GetType() {
	case values.ArrayValue:
		return "array", nil
	case values.BoolValue:
		return "boolean", nil
	case values.FloatValue:
		return "double", nil
	case values.IntValue:
		return "integer", nil
	case values.NullValue:
		return "NULL", nil
	case values.StrValue:
		return "string", nil
	default:
		return "unknown type", nil
	}
}

// ------------------- MARK: intval -------------------

func nativeFn_intval(args []values.RuntimeValue, _ runtime.Context) (values.RuntimeValue, phpError.Error) {
	args, err := NewFuncParamValidator("intval").addParam("$value", []string{"mixed"}, nil).validate(args)
	if err != nil {
		return values.NewVoid(), err
	}

	integer, err := lib_intval(args[0])
	return values.NewInt(integer), err
}

func lib_intval(runtimeValue values.RuntimeValue) (int64, phpError.Error) {
	// Spec: https://phplang.org/spec/08-conversions.html#converting-to-integer-type

	switch runtimeValue.GetType() {
	case values.ArrayValue:
		// Spec: https://phplang.org/spec/08-conversions.html#converting-to-integer-type
		// If the source is an array with zero elements, the result value is 0; otherwise, the result value is 1.
		if len(runtimeValue.(*values.Array).Elements) == 0 {
			return 0, nil
		}
		return 1, nil
	case values.BoolValue:
		// Spec: https://phplang.org/spec/08-conversions.html#converting-to-integer-type
		// If the source type is bool, then if the source value is FALSE, the result value is 0; otherwise, the result value is 1.
		if runtimeValue.(*values.Bool).Value {
			return 1, nil
		}
		return 0, nil
	case values.FloatValue:
		// Spec: https://phplang.org/spec/08-conversions.html#converting-to-integer-type
		// If the source type is float, for the values INF, -INF, and NAN, the result value is zero.
		// For all other values, if the precision can be preserved (that is, the float is within the range of an integer),
		// the fractional part is rounded towards zero.
		return int64(runtimeValue.(*values.Float).Value), nil
	case values.IntValue:
		return runtimeValue.(*values.Int).Value, nil
	case values.NullValue:
		// Spec: https://phplang.org/spec/08-conversions.html#converting-to-integer-type
		// If the source value is NULL, the result value is 0.
		return 0, nil
	case values.StrValue:
		// Spec: https://phplang.org/spec/08-conversions.html#converting-to-integer-type
		// If the source is a numeric string or leading-numeric string having integer format,
		// if the precision can be preserved the result value is that string’s integer value;
		// otherwise, the result is undefined.
		// If the source is a numeric string or leading-numeric string having floating-point format,
		// the string’s floating-point value is treated as described above for a conversion from float.
		// The trailing non-numeric characters in leading-numeric strings are ignored.
		// For any other string, the result value is 0.
		intStr := runtimeValue.(*values.Str).Value
		if common.IsFloatingLiteralWithSign(intStr) {
			return lib_intval(values.NewFloat(common.FloatingLiteralToFloat64WithSign(intStr)))
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

func nativeFn_is_array(args []values.RuntimeValue, _ runtime.Context) (values.RuntimeValue, phpError.Error) {
	args, err := NewFuncParamValidator("is_array").addParam("$value", []string{"mixed"}, nil).validate(args)
	if err != nil {
		return values.NewVoid(), err
	}

	return values.NewBool(lib_is_array(args[0])), nil
}

func lib_is_array(runtimeValue values.RuntimeValue) bool {
	// Spec: https://www.php.net/manual/en/function.is-array
	return runtimeValue.GetType() == values.ArrayValue
}

// ------------------- MARK: is_bool -------------------

func nativeFn_is_bool(args []values.RuntimeValue, _ runtime.Context) (values.RuntimeValue, phpError.Error) {
	args, err := NewFuncParamValidator("is_bool").addParam("$value", []string{"mixed"}, nil).validate(args)
	if err != nil {
		return values.NewVoid(), err
	}

	return values.NewBool(lib_is_bool(args[0])), nil
}

func lib_is_bool(runtimeValue values.RuntimeValue) bool {
	// Spec: https://www.php.net/manual/en/function.is-bool
	return runtimeValue.GetType() == values.BoolValue
}

// ------------------- MARK: is_float -------------------

func nativeFn_is_float(args []values.RuntimeValue, _ runtime.Context) (values.RuntimeValue, phpError.Error) {
	args, err := NewFuncParamValidator("is_float").addParam("$value", []string{"mixed"}, nil).validate(args)
	if err != nil {
		return values.NewVoid(), err
	}

	return values.NewBool(lib_is_float(args[0])), nil
}

func lib_is_float(runtimeValue values.RuntimeValue) bool {
	// Spec: https://www.php.net/manual/en/function.is-float
	return runtimeValue.GetType() == values.FloatValue
}

// ------------------- MARK: is_int -------------------

func nativeFn_is_int(args []values.RuntimeValue, _ runtime.Context) (values.RuntimeValue, phpError.Error) {
	args, err := NewFuncParamValidator("is_int").addParam("$value", []string{"mixed"}, nil).validate(args)
	if err != nil {
		return values.NewVoid(), err
	}

	return values.NewBool(lib_is_int(args[0])), nil
}

func lib_is_int(runtimeValue values.RuntimeValue) bool {
	// Spec: https://www.php.net/manual/en/function.is-int
	return runtimeValue.GetType() == values.IntValue
}

// ------------------- MARK: is_null -------------------

func nativeFn_is_null(args []values.RuntimeValue, _ runtime.Context) (values.RuntimeValue, phpError.Error) {
	args, err := NewFuncParamValidator("is_null").addParam("$value", []string{"mixed"}, nil).validate(args)
	if err != nil {
		return values.NewVoid(), err
	}

	return values.NewBool(lib_is_null(args[0])), nil
}

func lib_is_null(runtimeValue values.RuntimeValue) bool {
	// Spec: https://www.php.net/manual/en/function.is-null.php
	return runtimeValue.GetType() == values.NullValue
}

// ------------------- MARK: is_scalar -------------------

func nativeFn_is_scalar(args []values.RuntimeValue, _ runtime.Context) (values.RuntimeValue, phpError.Error) {
	args, err := NewFuncParamValidator("is_scalar").addParam("$value", []string{"mixed"}, nil).validate(args)
	if err != nil {
		return values.NewVoid(), err
	}

	return values.NewBool(lib_is_scalar(args[0])), nil
}

func lib_is_scalar(runtimeValue values.RuntimeValue) bool {
	// Spec: https://www.php.net/manual/en/function.is-scalar.php
	return slices.Contains([]values.ValueType{values.BoolValue, values.IntValue, values.FloatValue, values.StrValue}, runtimeValue.GetType())
}

// ------------------- MARK: is_string -------------------

func nativeFn_is_string(args []values.RuntimeValue, _ runtime.Context) (values.RuntimeValue, phpError.Error) {
	args, err := NewFuncParamValidator("is_string").addParam("$value", []string{"mixed"}, nil).validate(args)
	if err != nil {
		return values.NewVoid(), err
	}

	return values.NewBool(lib_is_string(args[0])), nil
}

func lib_is_string(runtimeValue values.RuntimeValue) bool {
	// Spec: https://www.php.net/manual/en/function.is-string
	return runtimeValue.GetType() == values.StrValue
}

// ------------------- MARK: print_r -------------------

func nativeFn_print_r(args []values.RuntimeValue, context runtime.Context) (values.RuntimeValue, phpError.Error) {
	args, err := NewFuncParamValidator("print_r").
		addParam("$value", []string{"mixed"}, nil).
		addParam("$return", []string{"bool"}, values.NewBool(false)).
		validate(args)
	if err != nil {
		return values.NewVoid(), err
	}

	return lib_print_r(args[0], args[1].(*values.Bool).Value, context)
}

func lib_print_r(value values.RuntimeValue, returnValue bool, context runtime.Context) (values.RuntimeValue, phpError.Error) {
	result, err := lib_print_r_var(value, 4)
	if err != nil {
		return values.NewVoid(), err
	}

	if returnValue {
		return values.NewStr(result), nil
	} else {
		context.Interpreter.Print(result)
		return values.NewBool(true), nil
	}
}

func lib_print_r_var(value values.RuntimeValue, depth int) (string, phpError.Error) {
	result := ""
	var err phpError.Error
	switch value.GetType() {
	case values.ArrayValue:
		array := value.(*values.Array)
		result = fmt.Sprintf("Array\n%s(\n", strings.Repeat(" ", depth-4))
		for _, key := range array.Keys {
			keyStr, err := lib_print_r_var(key, depth+4)
			if err != nil {
				return "", err
			}
			elementValue, _ := array.GetElement(key)
			valueStr, err := lib_print_r_var(elementValue, depth+8)
			if err != nil {
				return "", err
			}

			result += fmt.Sprintf("%s[%s] => %s\n", strings.Repeat(" ", depth), keyStr, valueStr)
		}
		result += fmt.Sprintf("%s)\n", strings.Repeat(" ", depth-4))
	case values.BoolValue:
		if value.(*values.Bool).Value {
			result = "1"
		} else {
			result = ""
		}
	case values.FloatValue, values.IntValue:
		result, err = lib_strval(value)
		if err != nil {
			return "", err
		}
	case values.NullValue:
		result = ""
	case values.StrValue:
		result = value.(*values.Str).Value
	default:
		return "", phpError.NewError("lib_print_r_var: Unsupported runtime value %s", value.GetType())
	}
	return result, nil
}

// ------------------- MARK: strval -------------------

func nativeFn_strval(args []values.RuntimeValue, _ runtime.Context) (values.RuntimeValue, phpError.Error) {
	args, err := NewFuncParamValidator("strval").addParam("$value", []string{"mixed"}, nil).validate(args)
	if err != nil {
		return values.NewVoid(), err
	}

	str, err := lib_strval(args[0])
	return values.NewStr(str), err
}

func lib_strval(runtimeValue values.RuntimeValue) (string, phpError.Error) {
	// Spec: https://phplang.org/spec/08-conversions.html#converting-to-string-type

	switch runtimeValue.GetType() {
	case values.ArrayValue:
		// Spec: https://phplang.org/spec/08-conversions.html#converting-to-string-type
		// If the source is an array, the conversion is invalid. The result value is the string “Array” and a non-fatal error is produced.
		return "Array", nil
		// TODO lib_strval - array non-fatal error: "Warning: Array to string conversion in /home/user/scripts/code.php on line 2" - only if E_ALL | E_WARNING
	case values.BoolValue:
		// Spec: https://phplang.org/spec/08-conversions.html#converting-to-string-type
		// If the source type is bool, then if the source value is FALSE, the result value is the empty string;
		// otherwise, the result value is “1”.
		if !runtimeValue.(*values.Bool).Value {
			return "", nil
		}
		return "1", nil
	case values.FloatValue:
		// Spec: https://phplang.org/spec/08-conversions.html#converting-to-string-type
		// If the source type is int or float, then the result value is a string containing the textual representation
		// of the source value (as specified by the library function sprintf).
		return runtimeValue.(*values.Float).ToPhpString(), nil
	case values.IntValue:
		// Spec: https://phplang.org/spec/08-conversions.html#converting-to-string-type
		// If the source type is int or float, then the result value is a string containing the textual representation
		// of the source value (as specified by the library function sprintf).
		return fmt.Sprintf("%d", runtimeValue.(*values.Int).Value), nil
	case values.NullValue:
		// Spec: https://phplang.org/spec/08-conversions.html#converting-to-string-type
		// If the source value is NULL, the result value is the empty string.
		return "", nil
	case values.StrValue:
		return runtimeValue.(*values.Str).Value, nil
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

func nativeFn_var_dump(args []values.RuntimeValue, context runtime.Context) (values.RuntimeValue, phpError.Error) {
	// Spec: https://www.php.net/manual/en/function.var-dump

	args, err := NewFuncParamValidator("var_dump").
		addParam("$value", []string{"mixed"}, nil).addVariableLenParam("$values", []string{"mixed"}).
		validate(args)
	if err != nil {
		return values.NewVoid(), err
	}

	if err := lib_var_dump_var(context, args[0], 2); err != nil {
		return values.NewVoid(), err
	}

	if len(args) == 2 {
		arrayValues := args[1].(*values.Array)
		for _, key := range arrayValues.Keys {
			argValue, _ := arrayValues.GetElement(key)
			if err := lib_var_dump_var(context, argValue, 2); err != nil {
				return values.NewVoid(), err
			}
		}
	}

	return values.NewVoid(), nil
}

func lib_var_dump_var(context runtime.Context, value values.RuntimeValue, depth int) phpError.Error {
	switch value.GetType() {
	case values.ArrayValue:
		array := value.(*values.Array)
		context.Interpreter.Println(fmt.Sprintf("array(%d) {", len(array.Keys)))
		for _, key := range array.Keys {
			switch key.GetType() {
			case values.IntValue:
				keyValue := key.(*values.Int).Value
				context.Interpreter.Println(fmt.Sprintf("%s[%d]=>", strings.Repeat(" ", depth), keyValue))
			case values.StrValue:
				keyValue := key.(*values.Str).Value
				context.Interpreter.Println(fmt.Sprintf(`%s["%s"]=>`, strings.Repeat(" ", depth), keyValue))
			default:
				return phpError.NewError("lib_var_dump_var: Unsupported array key type %s", key.GetType())
			}
			context.Interpreter.Print(strings.Repeat(" ", depth))
			elementValue, _ := array.GetElement(key)
			if err := lib_var_dump_var(context, elementValue, depth+2); err != nil {
				return err
			}
		}
		context.Interpreter.Println(strings.Repeat(" ", depth-2) + "}")
	case values.BoolValue:
		if value.(*values.Bool).Value {
			context.Interpreter.Println("bool(true)")
		} else {
			context.Interpreter.Println("bool(false)")
		}
	case values.FloatValue:
		strVal, err := lib_strval(value)
		if err != nil {
			return err
		}
		context.Interpreter.Println("float(" + strVal + ")")
	case values.IntValue:
		strVal, err := lib_strval(value)
		if err != nil {
			return err
		}
		context.Interpreter.Println("int(" + strVal + ")")
	case values.NullValue:
		context.Interpreter.Println("NULL")
	case values.StrValue:
		strVal := value.(*values.Str).Value
		context.Interpreter.Println(fmt.Sprintf("string(%d) \"%s\"", len(strVal), strVal))
	default:
		return phpError.NewError("lib_var_dump_var: Unsupported runtime value %s", value.GetType())
	}
	return nil
}

// ------------------- MARK: var_export -------------------

func nativeFn_var_export(args []values.RuntimeValue, interpreter runtime.Context) (values.RuntimeValue, phpError.Error) {
	// Spec: https://www.php.net/manual/en/function.var-export

	args, err := NewFuncParamValidator("var_dump").
		addParam("$value", []string{"mixed"}, nil).
		addParam("$return", []string{"bool"}, values.NewBool(false)).
		validate(args)
	if err != nil {
		return values.NewVoid(), err
	}

	return lib_var_export(args[0], args[1].(*values.Bool).Value, interpreter)
}

func lib_var_export(value values.RuntimeValue, returnValue bool, context runtime.Context) (values.RuntimeValue, phpError.Error) {
	result, err := lib_var_export_var(value, 2)
	if err != nil {
		return values.NewVoid(), err
	}

	if returnValue {
		return values.NewStr(result), nil
	} else {
		context.Interpreter.Print(result)
		return values.NewNull(), nil
	}
}

func lib_var_export_var(value values.RuntimeValue, depth int) (string, phpError.Error) {
	result := ""
	var err phpError.Error
	switch value.GetType() {
	case values.ArrayValue:
		array := value.(*values.Array)
		result = fmt.Sprintf("%sarray (\n", strings.Repeat(" ", depth-2))
		for _, key := range array.Keys {
			keyStr, err := lib_var_export_var(key, depth+2)
			if err != nil {
				return "", err
			}
			elementValue, _ := array.GetElement(key)
			valueStr, err := lib_var_export_var(elementValue, depth+2)
			if err != nil {
				return "", err
			}

			if elementValue.GetType() == values.ArrayValue {
				result += fmt.Sprintf("%s%s => \n%s%s,\n",
					strings.Repeat(" ", depth), keyStr,
					strings.Repeat(" ", depth-2), valueStr,
				)
			} else {
				result += fmt.Sprintf("%s%s => %s,\n", strings.Repeat(" ", depth), keyStr, valueStr)
			}
		}
		result += fmt.Sprintf("%s)", strings.Repeat(" ", depth-2))
	case values.BoolValue:
		if value.(*values.Bool).Value {
			result = "true"
		} else {
			result = "false"
		}
	case values.FloatValue, values.IntValue:
		result, err = lib_strval(value)
		if err != nil {
			return "", err
		}
	case values.NullValue:
		result = "NULL"
	case values.StrValue:
		result = "'" + value.(*values.Str).Value + "'"
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
