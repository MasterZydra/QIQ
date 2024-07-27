package interpreter

import (
	"GoPHP/cmd/goPHP/phpError"
	"math"
)

func registerNativeMathFunctions(environment *Environment) {
	environment.nativeFunctions["abs"] = nativeFn_abs
	environment.nativeFunctions["acos"] = nativeFn_acos
	environment.nativeFunctions["acosh"] = nativeFn_acosh
	environment.nativeFunctions["asin"] = nativeFn_asin
	environment.nativeFunctions["asinh"] = nativeFn_asinh
	environment.nativeFunctions["pi"] = nativeFn_pi
}

// ------------------- MARK: abs -------------------

func nativeFn_abs(args []IRuntimeValue, _ *Interpreter) (IRuntimeValue, phpError.Error) {
	args, err := NewFuncParamValidator("abs").addParam("$num", []string{"int", "float"}, nil).validate(args)
	if err != nil {
		return NewVoidRuntimeValue(), err
	}

	// Spec: https://www.php.net/manual/en/function.abs.php

	numType := args[0].GetType()

	var numValue float64 = 0
	if numType == FloatingValue {
		numValue = args[0].(*FloatingRuntimeValue).Value
	}
	if numType == IntegerValue {
		numValue = float64(args[0].(*IntegerRuntimeValue).Value)
	}
	numValue = math.Abs(numValue)
	if numType == FloatingValue {
		return NewFloatingRuntimeValue(numValue), nil
	}
	if numType == IntegerValue {
		return NewIntegerRuntimeValue(int64(numValue)), nil
	}

	return NewVoidRuntimeValue(), phpError.NewError("Unsupported value type %s", numType)
}

// ------------------- MARK: acos -------------------

func nativeFn_acos(args []IRuntimeValue, _ *Interpreter) (IRuntimeValue, phpError.Error) {
	args, err := NewFuncParamValidator("acos").addParam("$num", []string{"float"}, nil).validate(args)
	if err != nil {
		return NewVoidRuntimeValue(), err
	}

	// Spec: https://www.php.net/manual/en/function.acos.php
	return NewFloatingRuntimeValue(math.Acos(args[0].(*FloatingRuntimeValue).Value)), nil
}

// ------------------- MARK: acosh -------------------

func nativeFn_acosh(args []IRuntimeValue, _ *Interpreter) (IRuntimeValue, phpError.Error) {
	args, err := NewFuncParamValidator("acosh").addParam("$num", []string{"float"}, nil).validate(args)
	if err != nil {
		return NewVoidRuntimeValue(), err
	}

	// Spec: https://www.php.net/manual/en/function.acosh.php
	return NewFloatingRuntimeValue(math.Acosh(args[0].(*FloatingRuntimeValue).Value)), nil
}

// ------------------- MARK: asin -------------------

func nativeFn_asin(args []IRuntimeValue, _ *Interpreter) (IRuntimeValue, phpError.Error) {
	args, err := NewFuncParamValidator("asin").addParam("$num", []string{"float"}, nil).validate(args)
	if err != nil {
		return NewVoidRuntimeValue(), err
	}

	// Spec: https://www.php.net/manual/en/function.sin.php
	return NewFloatingRuntimeValue(math.Asin(args[0].(*FloatingRuntimeValue).Value)), nil
}

// ------------------- MARK: asinh -------------------

func nativeFn_asinh(args []IRuntimeValue, _ *Interpreter) (IRuntimeValue, phpError.Error) {
	args, err := NewFuncParamValidator("asinh").addParam("$num", []string{"float"}, nil).validate(args)
	if err != nil {
		return NewVoidRuntimeValue(), err
	}

	// Spec: https://www.php.net/manual/en/function.sinh.php
	return NewFloatingRuntimeValue(math.Asinh(args[0].(*FloatingRuntimeValue).Value)), nil
}

// ------------------- MARK: pi -------------------

func nativeFn_pi(args []IRuntimeValue, interpreter *Interpreter) (IRuntimeValue, phpError.Error) {
	_, err := NewFuncParamValidator("pi").validate(args)
	if err != nil {
		return NewVoidRuntimeValue(), err
	}

	// Spec: https://www.php.net/manual/en/function.pi.php
	return interpreter.env.predefinedConstants["M_PI"], nil
}

// TODO atan
// TODO atan2
// TODO atanh
// TODO base_convert
// TODO bindec
// TODO ceil
// TODO cos
// TODO cosh
// TODO decbin
// TODO dechex
// TODO decoct
// TODO deg2rad
// TODO exp
// TODO expm1
// TODO fdiv
// TODO floor
// TODO fmod
// TODO hexdec
// TODO hypot
// TODO intdiv
// TODO is_finite
// TODO is_infinite
// TODO is_nan
// TODO log
// TODO log10
// TODO log1p
// TODO max
// TODO min
// TODO octdec
// TODO pow
// TODO rad2deg
// TODO round
// TODO sin
// TODO sinh
// TODO sqrt
// TODO tan
// TODO tanh
