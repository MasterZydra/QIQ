package math

import (
	"GoPHP/cmd/goPHP/phpError"
	"GoPHP/cmd/goPHP/runtime"
	"GoPHP/cmd/goPHP/runtime/funcParamValidator"
	"GoPHP/cmd/goPHP/runtime/values"
	goMath "math"
)

func Register(environment runtime.Environment) {
	// Category: Math Functions
	environment.AddNativeFunction("abs", nativeFn_abs)
	environment.AddNativeFunction("acos", nativeFn_acos)
	environment.AddNativeFunction("acosh", nativeFn_acosh)
	environment.AddNativeFunction("asin", nativeFn_asin)
	environment.AddNativeFunction("asinh", nativeFn_asinh)
	environment.AddNativeFunction("pi", nativeFn_pi)
}

// ------------------- MARK: abs -------------------

func nativeFn_abs(args []values.RuntimeValue, _ runtime.Context) (values.RuntimeValue, phpError.Error) {
	args, err := funcParamValidator.NewValidator("abs").AddParam("$num", []string{"int", "float"}, nil).Validate(args)
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
	numValue = goMath.Abs(numValue)
	if numType == values.FloatValue {
		return values.NewFloat(numValue), nil
	}
	if numType == values.IntValue {
		return values.NewInt(int64(numValue)), nil
	}

	return values.NewVoid(), phpError.NewError("Unsupported value type %s", numType)
}

// ------------------- MARK: acos -------------------

func nativeFn_acos(args []values.RuntimeValue, _ runtime.Context) (values.RuntimeValue, phpError.Error) {
	args, err := funcParamValidator.NewValidator("acos").AddParam("$num", []string{"float"}, nil).Validate(args)
	if err != nil {
		return values.NewVoid(), err
	}

	// Spec: https://www.php.net/manual/en/function.acos.php
	return values.NewFloat(goMath.Acos(args[0].(*values.Float).Value)), nil
}

// ------------------- MARK: acosh -------------------

func nativeFn_acosh(args []values.RuntimeValue, _ runtime.Context) (values.RuntimeValue, phpError.Error) {
	args, err := funcParamValidator.NewValidator("acosh").AddParam("$num", []string{"float"}, nil).Validate(args)
	if err != nil {
		return values.NewVoid(), err
	}

	// Spec: https://www.php.net/manual/en/function.acosh.php
	return values.NewFloat(goMath.Acosh(args[0].(*values.Float).Value)), nil
}

// ------------------- MARK: asin -------------------

func nativeFn_asin(args []values.RuntimeValue, _ runtime.Context) (values.RuntimeValue, phpError.Error) {
	args, err := funcParamValidator.NewValidator("asin").AddParam("$num", []string{"float"}, nil).Validate(args)
	if err != nil {
		return values.NewVoid(), err
	}

	// Spec: https://www.php.net/manual/en/function.sin.php
	return values.NewFloat(goMath.Asin(args[0].(*values.Float).Value)), nil
}

// ------------------- MARK: asinh -------------------

func nativeFn_asinh(args []values.RuntimeValue, _ runtime.Context) (values.RuntimeValue, phpError.Error) {
	args, err := funcParamValidator.NewValidator("asinh").AddParam("$num", []string{"float"}, nil).Validate(args)
	if err != nil {
		return values.NewVoid(), err
	}

	// Spec: https://www.php.net/manual/en/function.sinh.php
	return values.NewFloat(goMath.Asinh(args[0].(*values.Float).Value)), nil
}

// ------------------- MARK: pi -------------------

func nativeFn_pi(args []values.RuntimeValue, context runtime.Context) (values.RuntimeValue, phpError.Error) {
	_, err := funcParamValidator.NewValidator("pi").Validate(args)
	if err != nil {
		return values.NewVoid(), err
	}

	// Spec: https://www.php.net/manual/en/function.pi.php
	return context.Env.LookupConstant("M_PI")
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
