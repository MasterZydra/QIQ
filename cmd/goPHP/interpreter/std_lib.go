package interpreter

import (
	"fmt"
	"math"
	"slices"
)

func registerNativeFunctions(environment *Environment) {
	environment.nativeFunctions["array_key_exits"] = nativeFn_array_key_exists
	environment.nativeFunctions["boolval"] = nativeFn_boolval
	environment.nativeFunctions["error_reporting"] = nativeFn_error_reporting
	environment.nativeFunctions["floatval"] = nativeFn_floatval
	environment.nativeFunctions["intval"] = nativeFn_intval
	environment.nativeFunctions["is_null"] = nativeFn_is_null
	environment.nativeFunctions["key_exits"] = nativeFn_array_key_exists
	environment.nativeFunctions["strval"] = nativeFn_strval
}

type nativeFunction func([]IRuntimeValue, *Environment) (IRuntimeValue, error)

// ------------------- MARK: array_key_exits -------------------

func nativeFn_array_key_exists(args []IRuntimeValue, env *Environment) (IRuntimeValue, error) {
	if len(args) != 2 {
		return NewVoidRuntimeValue(), fmt.Errorf(
			"Uncaught ArgumentCountError: array_key_exists() expects exactly 2 arguments, %d given", len(args),
		)
	}

	if args[1].GetType() != ArrayValue {
		return NewVoidRuntimeValue(), fmt.Errorf("Uncaught TypeError: $array must be of type array")
	}

	boolean, err := lib_array_key_exists(args[0], runtimeValToArrayRuntimeVal(args[1]))
	return NewBooleanRuntimeValue(boolean), err
}

func lib_array_key_exists(key IRuntimeValue, array IArrayRuntimeValue) (bool, error) {
	// Spec: https://www.php.net/manual/en/function.array-key-exists.php

	// TODO lib_array_key_exists - allowedKeyTypes - resource
	allowedKeyTypes := []ValueType{StringValue, IntegerValue, FloatingValue, BooleanValue, NullValue}

	if !slices.Contains(allowedKeyTypes, key.GetType()) {
		return false, fmt.Errorf("Values of type %s are not allowed as array key", key.GetType())
	}

	_, ok := array.GetElement(key)
	return ok, nil
}

// ------------------- MARK: boolval -------------------

func nativeFn_boolval(args []IRuntimeValue, env *Environment) (IRuntimeValue, error) {
	if len(args) != 1 {
		return NewVoidRuntimeValue(), fmt.Errorf(
			"Uncaught ArgumentCountError: boolval() expects exactly 1 argument, %d given", len(args),
		)
	}

	boolean, err := lib_boolval(args[0])
	return NewBooleanRuntimeValue(boolean), err
}

func lib_boolval(runtimeValue IRuntimeValue) (bool, error) {
	// Spec: https://phplang.org/spec/08-conversions.html#converting-to-boolean-type

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
		return false, fmt.Errorf("boolval: Unsupported runtime value %s", runtimeValue.GetType())
	}

	// TODO boolval - object
	// Spec: https://phplang.org/spec/08-conversions.html#converting-to-boolean-type
	// If the source is an object, the result value is TRUE.

	// TODO boolval - resource
	// Spec: https://phplang.org/spec/08-conversions.html#converting-to-boolean-type
	// If the source is a resource, the result value is TRUE.
}

// ------------------- MARK: error_reporting -------------------

func nativeFn_error_reporting(args []IRuntimeValue, env *Environment) (IRuntimeValue, error) {
	// Spec: https://www.php.net/manual/en/function.error-reporting.php

	// TODO implement nativeFn_error_reporting after adding configuration for error level
	return NewIntegerRuntimeValue(0), nil
}

// ------------------- MARK: floatval -------------------

func nativeFn_floatval(args []IRuntimeValue, env *Environment) (IRuntimeValue, error) {
	if len(args) != 1 {
		return NewVoidRuntimeValue(), fmt.Errorf(
			"Uncaught ArgumentCountError: floatval() expects exactly 1 argument, %d given", len(args),
		)
	}

	floating, err := lib_floatval(args[0])
	return NewFloatingRuntimeValue(floating), err
}

func lib_floatval(runtimeValue IRuntimeValue) (float64, error) {
	// Spec: https://phplang.org/spec/08-conversions.html#converting-to-floating-point-type

	switch runtimeValue.GetType() {
	case FloatingValue:
		return runtimeValToFloatRuntimeVal(runtimeValue).GetValue(), nil
	case IntegerValue:
		// Spec: https://phplang.org/spec/08-conversions.html#converting-to-floating-point-type
		// If the source type is int,
		// if the precision can be preserved the result value is the closest approximation to the source value;
		// otherwise, the result is undefined.
		return float64(runtimeValToIntRuntimeVal(runtimeValue).GetValue()), nil
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

	// TODO lib_floatval - string
	// Spec: https://phplang.org/spec/08-conversions.html#converting-to-floating-point-type
	// If the source is a numeric string or leading-numeric string having integer format, the string’s integer value is treated as described above for a conversion from int. If the source is a numeric string or leading-numeric string having floating-point format, the result value is the closest approximation to the string’s floating-point value. The trailing non-numeric characters in leading-numeric strings are ignored. For any other string, the result value is 0.

	// TODO lib_floatval - object
	// Spec: https://phplang.org/spec/08-conversions.html#converting-to-floating-point-type
	// If the source is an object, if the class defines a conversion function, the result is determined by that function (this is currently available only to internal classes). If not, the conversion is invalid, the result is assumed to be 1.0 and a non-fatal error is produced.
}

// ------------------- MARK: intval -------------------

func nativeFn_intval(args []IRuntimeValue, env *Environment) (IRuntimeValue, error) {
	if len(args) != 1 {
		return NewVoidRuntimeValue(), fmt.Errorf(
			"Uncaught ArgumentCountError: intval() expects exactly 1 argument, %d given", len(args),
		)
	}

	integer, err := lib_intval(args[0])
	return NewIntegerRuntimeValue(integer), err
}

func lib_intval(runtimeValue IRuntimeValue) (int64, error) {
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
	case IntegerValue:
		return runtimeValToIntRuntimeVal(runtimeValue).GetValue(), nil
	case NullValue:
		// Spec: https://phplang.org/spec/08-conversions.html#converting-to-integer-type
		// If the source value is NULL, the result value is 0.
		return 0, nil
	default:
		return 0, fmt.Errorf("lib_intval: Unsupported runtime value %s", runtimeValue.GetType())
	}
	// TODO lib_intval - float
	// Spec: https://phplang.org/spec/08-conversions.html#converting-to-integer-type
	// If the source type is float, for the values INF, -INF, and NAN, the result value is zero. For all other values, if the precision can be preserved (that is, the float is within the range of an integer), the fractional part is rounded towards zero. If the precision cannot be preserved, the following conversion algorithm is used, where X is defined as two to the power of the number of bits in an integer (for example, 2 to the power of 32, i.e. 4294967296):
	// 1. We take the floating point remainder (wherein the remainder has the same sign as the dividend) of dividing the float by X, rounded towards zero.
	// 2. If the remainder is less than zero, it is rounded towards infinity and X is added.
	// 3. This result is converted to an unsigned integer.
	// 4. This result is converted to a signed integer by treating the unsigned integer as a two’s complement representation of the signed integer.
	// Implementations may implement this conversion differently (for example, on some architectures there may be hardware support for this specific conversion mode) so long as the result is the same.

	// TODO lib_intval - string
	// Spec: https://phplang.org/spec/08-conversions.html#converting-to-integer-type
	// If the source is a numeric string or leading-numeric string having integer format, if the precision can be preserved the result value is that string’s integer value; otherwise, the result is undefined. If the source is a numeric string or leading-numeric string having floating-point format, the string’s floating-point value is treated as described above for a conversion from float. The trailing non-numeric characters in leading-numeric strings are ignored. For any other string, the result value is 0.

	// TODO lib_intval - object
	// Spec: https://phplang.org/spec/08-conversions.html#converting-to-integer-type
	// If the source is an object, if the class defines a conversion function, the result is determined by that function (this is currently available only to internal classes). If not, the conversion is invalid, the result is assumed to be 1 and a non-fatal error is produced.

	// TODO lib_intval - resource
	// Spec: https://phplang.org/spec/08-conversions.html#converting-to-integer-type
	// If the source is a resource, the result is the resource’s unique ID.
}

// ------------------- MARK: is_null -------------------

func nativeFn_is_null(args []IRuntimeValue, env *Environment) (IRuntimeValue, error) {
	if len(args) != 1 {
		return NewVoidRuntimeValue(), fmt.Errorf(
			"Uncaught ArgumentCountError: is_null() expects exactly 1 argument, %d given", len(args),
		)
	}

	return NewBooleanRuntimeValue(lib_is_null(args[0])), nil
}

func lib_is_null(runtimeValue IRuntimeValue) bool {
	// Spec: https://www.php.net/manual/en/function.is-null.php
	return runtimeValue.GetType() == NullValue
}

// ------------------- MARK: strval -------------------

func nativeFn_strval(args []IRuntimeValue, env *Environment) (IRuntimeValue, error) {
	if len(args) != 1 {
		return NewVoidRuntimeValue(), fmt.Errorf(
			"Uncaught ArgumentCountError: strval() expects exactly 1 argument,  %d given", len(args),
		)
	}

	str, err := lib_strval(args[0])
	return NewStringRuntimeValue(str), err
}

func lib_strval(runtimeValue IRuntimeValue) (string, error) {
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
		return "", fmt.Errorf("lib_strval: Unsupported runtime value %s", runtimeValue.GetType())
	}

	// TODO lib_strval - object
	// Spec: https://phplang.org/spec/08-conversions.html#converting-to-string-type
	// If the source is an object, then if that object’s class has a __toString method, the result value is the string returned by that method; otherwise, the conversion is invalid and a fatal error is produced.

	// TODO lib_strval - resource
	// Spec: https://phplang.org/spec/08-conversions.html#converting-to-string-type
	// If the source is a resource, the result value is an implementation-defined string.
}
