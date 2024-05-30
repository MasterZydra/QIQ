package interpreter

import (
	"fmt"
	"math"
	"slices"
	"strings"
)

func registerNativeFunctions(environment *Environment) {
	environment.nativeFunctions["array_key_exits"] = nativeFn_array_key_exists
	environment.nativeFunctions["boolval"] = nativeFn_boolval
	environment.nativeFunctions["error_reporting"] = nativeFn_error_reporting
	environment.nativeFunctions["floatval"] = nativeFn_floatval
	environment.nativeFunctions["intval"] = nativeFn_intval
	environment.nativeFunctions["is_null"] = nativeFn_is_null
	environment.nativeFunctions["is_scalar"] = nativeFn_is_scalar
	environment.nativeFunctions["key_exits"] = nativeFn_array_key_exists
	environment.nativeFunctions["strval"] = nativeFn_strval
	environment.nativeFunctions["var_dump"] = nativeFn_var_dump
}

type nativeFunction func([]IRuntimeValue, *Interpreter) (IRuntimeValue, Error)

// ------------------- MARK: array_key_exits -------------------

func nativeFn_array_key_exists(args []IRuntimeValue, _ *Interpreter) (IRuntimeValue, Error) {
	if len(args) != 2 {
		return NewVoidRuntimeValue(),
			NewError("Uncaught ArgumentCountError: array_key_exists() expects exactly 2 arguments, %d given", len(args))
	}

	if args[1].GetType() != ArrayValue {
		return NewVoidRuntimeValue(), NewError("Uncaught TypeError: $array must be of type array")
	}

	boolean, err := lib_array_key_exists(args[0], runtimeValToArrayRuntimeVal(args[1]))
	return NewBooleanRuntimeValue(boolean), err
}

func lib_array_key_exists(key IRuntimeValue, array IArrayRuntimeValue) (bool, Error) {
	// Spec: https://www.php.net/manual/en/function.array-key-exists.php

	// TODO lib_array_key_exists - allowedKeyTypes - resource
	allowedKeyTypes := []ValueType{StringValue, IntegerValue, FloatingValue, BooleanValue, NullValue}

	if !slices.Contains(allowedKeyTypes, key.GetType()) {
		return false, NewError("Values of type %s are not allowed as array key", key.GetType())
	}

	_, ok := array.GetElement(key)
	return ok, nil
}

// ------------------- MARK: arrayval -------------------

// This is not an official function. But converting different types to array is needed in several places
func lib_arrayval(runtimeValue IRuntimeValue) (IArrayRuntimeValue, Error) {
	// Spec: https://phplang.org/spec/08-conversions.html#converting-to-array-type

	// The result type is array.

	if runtimeValue.GetType() == NullValue {
		// Spec: https://phplang.org/spec/08-conversions.html#converting-to-array-type
		// If the source value is NULL, the result value is an array of zero elements.
		return NewArrayRuntimeValue(), nil
	}

	// TODO lib_arrayval - resource
	if lib_is_scalar(runtimeValue) {
		// Spec: https://phplang.org/spec/08-conversions.html#converting-to-array-type
		// If the source type is scalar or resource and it is non-NULL,
		// the result value is an array of one element under the key 0 whose value is that of the source.
		array := NewArrayRuntimeValue()
		array.SetElement(NewIntegerRuntimeValue(0), runtimeValue)
		return array, nil
	}

	// TODO lib_arrayval - object
	// Spec: https://phplang.org/spec/08-conversions.html#converting-to-array-type
	// If the source is an object, the result is an array of zero or more elements, where the elements are key/value pairs corresponding to the object’s instance properties. The order of insertion of the elements into the array is the lexical order of the instance properties in the class-member-declarations list.

	// TODO lib_arrayval - instance properties
	// Spec: https://phplang.org/spec/08-conversions.html#converting-to-array-type
	// For public instance properties, the keys of the array elements would be the same as the property name.
	// The key for a private instance property has the form “\0class\0name”, where the class is the class name, and the name is the property name.
	// The key for a protected instance property has the form “\0*\0name”, where name is that of the property.
	// The value for each key is that from the corresponding property, or NULL if the property was not initialized.

	return NewArrayRuntimeValue(), NewError("lib_arrayval: Unsupported type %s", runtimeValue.GetType())
}

// ------------------- MARK: boolval -------------------

func nativeFn_boolval(args []IRuntimeValue, _ *Interpreter) (IRuntimeValue, Error) {
	if len(args) != 1 {
		return NewVoidRuntimeValue(),
			NewError("Uncaught ArgumentCountError: boolval() expects exactly 1 argument, %d given", len(args))
	}

	boolean, err := lib_boolval(args[0])
	return NewBooleanRuntimeValue(boolean), err
}

func lib_boolval(runtimeValue IRuntimeValue) (bool, Error) {
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
		return false, NewError("boolval: Unsupported runtime value %s", runtimeValue.GetType())
	}

	// TODO boolval - object
	// Spec: https://phplang.org/spec/08-conversions.html#converting-to-boolean-type
	// If the source is an object, the result value is TRUE.

	// TODO boolval - resource
	// Spec: https://phplang.org/spec/08-conversions.html#converting-to-boolean-type
	// If the source is a resource, the result value is TRUE.
}

// ------------------- MARK: error_reporting -------------------

func nativeFn_error_reporting(args []IRuntimeValue, interpreter *Interpreter) (IRuntimeValue, Error) {
	// Spec: https://www.php.net/manual/en/function.error-reporting.php

	if len(args) == 0 || len(args) == 1 && args[0].GetType() == NullValue {
		return NewIntegerRuntimeValue(interpreter.config.ErrorReporting), nil
	}

	if len(args) != 1 {
		return NewVoidRuntimeValue(),
			NewError("Uncaught ArgumentCountError: error_reporting() expects most 1 argument, %d given", len(args))
	}

	if args[0].GetType() != IntegerValue {
		return NewVoidRuntimeValue(),
			NewError("Uncaught ArgumentError: error_reporting() expects most integer for $error_level argument, %s given", args[0].GetType())
	}

	newValue := runtimeValToIntRuntimeVal(args[0]).GetValue()
	if newValue == -1 {
		newValue = E_ALL
	}

	previous := interpreter.config.ErrorReporting
	interpreter.config.ErrorReporting = newValue

	return NewIntegerRuntimeValue(previous), nil
}

// ------------------- MARK: floatval -------------------

func nativeFn_floatval(args []IRuntimeValue, _ *Interpreter) (IRuntimeValue, Error) {
	if len(args) != 1 {
		return NewVoidRuntimeValue(),
			NewError("Uncaught ArgumentCountError: floatval() expects exactly 1 argument, %d given", len(args))
	}

	floating, err := lib_floatval(args[0])
	return NewFloatingRuntimeValue(floating), err
}

func lib_floatval(runtimeValue IRuntimeValue) (float64, Error) {
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

func nativeFn_intval(args []IRuntimeValue, _ *Interpreter) (IRuntimeValue, Error) {
	if len(args) != 1 {
		return NewVoidRuntimeValue(),
			NewError("Uncaught ArgumentCountError: intval() expects exactly 1 argument, %d given", len(args))
	}

	integer, err := lib_intval(args[0])
	return NewIntegerRuntimeValue(integer), err
}

func lib_intval(runtimeValue IRuntimeValue) (int64, Error) {
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
	default:
		return 0, NewError("lib_intval: Unsupported runtime value %s", runtimeValue.GetType())
	}

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

func nativeFn_is_null(args []IRuntimeValue, _ *Interpreter) (IRuntimeValue, Error) {
	if len(args) != 1 {
		return NewVoidRuntimeValue(),
			NewError("Uncaught ArgumentCountError: is_null() expects exactly 1 argument, %d given", len(args))
	}

	return NewBooleanRuntimeValue(lib_is_null(args[0])), nil
}

func lib_is_null(runtimeValue IRuntimeValue) bool {
	// Spec: https://www.php.net/manual/en/function.is-null.php
	return runtimeValue.GetType() == NullValue
}

// ------------------- MARK: is_scalar -------------------

func nativeFn_is_scalar(args []IRuntimeValue, _ *Interpreter) (IRuntimeValue, Error) {
	if len(args) != 1 {
		return NewVoidRuntimeValue(),
			NewError("Uncaught ArgumentCountError: is_scalar() expects exactly 1 argument, %d given", len(args))
	}

	return NewBooleanRuntimeValue(lib_is_scalar(args[0])), nil
}

func lib_is_scalar(runtimeValue IRuntimeValue) bool {
	// Spec: https://www.php.net/manual/en/function.is-scalar.php
	return slices.Contains([]ValueType{BooleanValue, IntegerValue, FloatingValue, StringValue}, runtimeValue.GetType())
}

// ------------------- MARK: strval -------------------

func nativeFn_strval(args []IRuntimeValue, _ *Interpreter) (IRuntimeValue, Error) {
	if len(args) != 1 {
		return NewVoidRuntimeValue(),
			NewError("Uncaught ArgumentCountError: strval() expects exactly 1 argument,  %d given", len(args))
	}

	str, err := lib_strval(args[0])
	return NewStringRuntimeValue(str), err
}

func lib_strval(runtimeValue IRuntimeValue) (string, Error) {
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
		return "", NewError("lib_strval: Unsupported runtime value %s", runtimeValue.GetType())
	}

	// TODO lib_strval - object
	// Spec: https://phplang.org/spec/08-conversions.html#converting-to-string-type
	// If the source is an object, then if that object’s class has a __toString method, the result value is the string returned by that method; otherwise, the conversion is invalid and a fatal error is produced.

	// TODO lib_strval - resource
	// Spec: https://phplang.org/spec/08-conversions.html#converting-to-string-type
	// If the source is a resource, the result value is an implementation-defined string.
}

// ------------------- MARK: var_dump -------------------

func nativeFn_var_dump(args []IRuntimeValue, interpreter *Interpreter) (IRuntimeValue, Error) {
	// Spec: https://www.php.net/manual/en/function.var-dump.php

	if len(args) == 0 {
		return NewVoidRuntimeValue(),
			NewError("Uncaught ArgumentCountError: var_dump() expects at least 1 argument,  %d given", len(args))
	}

	for _, arg := range args {
		if err := lib_var_dump_var(interpreter, arg, 2); err != nil {
			return NewVoidRuntimeValue(), err
		}
	}

	return NewVoidRuntimeValue(), nil
}

func lib_var_dump_var(interpreter *Interpreter, value IRuntimeValue, depth int) Error {
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
				return NewError("lib_var_dump_var: Unsupported array key type %s", key.GetType())
			}
			interpreter.print(strings.Repeat(" ", depth))
			if err := lib_var_dump_var(interpreter, elements[key], depth+2); err != nil {
				return err
			}
		}
		interpreter.println(strings.Repeat(" ", depth-2) + "}")
	case BooleanValue:
		boolean := runtimeValToBoolRuntimeVal(value).GetValue()
		if boolean {
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
		return NewError("lib_var_dump_var: Unsupported runtime value %s", value.GetType())
	}
	return nil
}
