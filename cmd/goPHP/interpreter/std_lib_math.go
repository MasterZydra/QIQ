package interpreter

import (
	"GoPHP/cmd/goPHP/phpError"
	"GoPHP/cmd/goPHP/runtime/values"
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

func nativeFn_abs(args []values.RuntimeValue, _ *Interpreter) (values.RuntimeValue, phpError.Error) {
	args, err := NewFuncParamValidator("abs").addParam("$num", []string{"int", "float"}, nil).validate(args)
	if err != nil {
		return values.NewVoid(), err
	}

	// Spec: https://www.php.net/manual/en/function.abs.php

	numType := args[0].GetType()

	var numValue float64 = 0
	if numType == values.FloatValue {
		numValue = args[0].(*values.Float).Value
	}
	if numType == values.IntValue {
		numValue = float64(args[0].(*values.Int).Value)
	}
	numValue = math.Abs(numValue)
	if numType == values.FloatValue {
		return values.NewFloat(numValue), nil
	}
	if numType == values.IntValue {
		return values.NewInt(int64(numValue)), nil
	}

	return values.NewVoid(), phpError.NewError("Unsupported value type %s", numType)
}

// ------------------- MARK: acos -------------------

func nativeFn_acos(args []values.RuntimeValue, _ *Interpreter) (values.RuntimeValue, phpError.Error) {
	args, err := NewFuncParamValidator("acos").addParam("$num", []string{"float"}, nil).validate(args)
	if err != nil {
		return values.NewVoid(), err
	}

	// Spec: https://www.php.net/manual/en/function.acos.php
	return values.NewFloat(math.Acos(args[0].(*values.Float).Value)), nil
}

// ------------------- MARK: acosh -------------------

func nativeFn_acosh(args []values.RuntimeValue, _ *Interpreter) (values.RuntimeValue, phpError.Error) {
	args, err := NewFuncParamValidator("acosh").addParam("$num", []string{"float"}, nil).validate(args)
	if err != nil {
		return values.NewVoid(), err
	}

	// Spec: https://www.php.net/manual/en/function.acosh.php
	return values.NewFloat(math.Acosh(args[0].(*values.Float).Value)), nil
}

// ------------------- MARK: asin -------------------

func nativeFn_asin(args []values.RuntimeValue, _ *Interpreter) (values.RuntimeValue, phpError.Error) {
	args, err := NewFuncParamValidator("asin").addParam("$num", []string{"float"}, nil).validate(args)
	if err != nil {
		return values.NewVoid(), err
	}

	// Spec: https://www.php.net/manual/en/function.sin.php
	return values.NewFloat(math.Asin(args[0].(*values.Float).Value)), nil
}

// ------------------- MARK: asinh -------------------

func nativeFn_asinh(args []values.RuntimeValue, _ *Interpreter) (values.RuntimeValue, phpError.Error) {
	args, err := NewFuncParamValidator("asinh").addParam("$num", []string{"float"}, nil).validate(args)
	if err != nil {
		return values.NewVoid(), err
	}

	// Spec: https://www.php.net/manual/en/function.sinh.php
	return values.NewFloat(math.Asinh(args[0].(*values.Float).Value)), nil
}

// ------------------- MARK: pi -------------------

func nativeFn_pi(args []values.RuntimeValue, interpreter *Interpreter) (values.RuntimeValue, phpError.Error) {
	_, err := NewFuncParamValidator("pi").validate(args)
	if err != nil {
		return values.NewVoid(), err
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
