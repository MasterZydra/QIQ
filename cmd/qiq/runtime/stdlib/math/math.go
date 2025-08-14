package math

import (
	"QIQ/cmd/qiq/phpError"
	"QIQ/cmd/qiq/runtime"
	"QIQ/cmd/qiq/runtime/funcParamValidator"
	"QIQ/cmd/qiq/runtime/values"
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

	// Const Category: Mathematical Constants
	// Spec: https://www.php.net/manual/en/math.constants.php
	environment.AddPredefinedConstant("M_1_PI", values.NewFloat(1/goMath.Pi))
	environment.AddPredefinedConstant("M_2_PI", values.NewFloat(2/goMath.Pi))
	environment.AddPredefinedConstant("M_2_SQRTPI", values.NewFloat(2/goMath.SqrtPi))
	environment.AddPredefinedConstant("M_E", values.NewFloat(goMath.E))
	environment.AddPredefinedConstant("M_EULER", values.NewFloat(goMath.E))
	environment.AddPredefinedConstant("M_LN10", values.NewFloat(goMath.Ln10))
	environment.AddPredefinedConstant("M_LN2", values.NewFloat(goMath.Ln2))
	environment.AddPredefinedConstant("M_LNPI", values.NewFloat(goMath.Log(goMath.Pi)))
	environment.AddPredefinedConstant("M_LOG10E", values.NewFloat(goMath.Log10E))
	environment.AddPredefinedConstant("M_LOG2E", values.NewFloat(goMath.Log2E))
	environment.AddPredefinedConstant("M_PI", values.NewFloat(goMath.Pi))
	environment.AddPredefinedConstant("M_PI_2", values.NewFloat(goMath.Pi/2))
	environment.AddPredefinedConstant("M_PI_4", values.NewFloat(goMath.Pi/4))
	environment.AddPredefinedConstant("M_SQRT1_2", values.NewFloat(1/goMath.Sqrt2))
	environment.AddPredefinedConstant("M_SQRT2", values.NewFloat(goMath.Sqrt2))
	environment.AddPredefinedConstant("M_SQRT3", values.NewFloat(goMath.Sqrt(3)))
	environment.AddPredefinedConstant("M_SQRTPI", values.NewFloat(goMath.SqrtPi))
	environment.AddPredefinedConstant("PHP_ROUND_HALF_UP", values.NewInt(1))
	environment.AddPredefinedConstant("PHP_ROUND_HALF_DOWN", values.NewInt(2))
	environment.AddPredefinedConstant("PHP_ROUND_HALF_EVEN", values.NewInt(3))
	environment.AddPredefinedConstant("PHP_ROUND_HALF_ODD", values.NewInt(4))
}

// -------------------------------------- abs -------------------------------------- MARK: abs

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

// -------------------------------------- acos -------------------------------------- MARK: acos

func nativeFn_acos(args []values.RuntimeValue, _ runtime.Context) (values.RuntimeValue, phpError.Error) {
	args, err := funcParamValidator.NewValidator("acos").AddParam("$num", []string{"float"}, nil).Validate(args)
	if err != nil {
		return values.NewVoid(), err
	}

	// Spec: https://www.php.net/manual/en/function.acos.php
	return values.NewFloat(goMath.Acos(args[0].(*values.Float).Value)), nil
}

// -------------------------------------- acosh -------------------------------------- MARK: acosh

func nativeFn_acosh(args []values.RuntimeValue, _ runtime.Context) (values.RuntimeValue, phpError.Error) {
	args, err := funcParamValidator.NewValidator("acosh").AddParam("$num", []string{"float"}, nil).Validate(args)
	if err != nil {
		return values.NewVoid(), err
	}

	// Spec: https://www.php.net/manual/en/function.acosh.php
	return values.NewFloat(goMath.Acosh(args[0].(*values.Float).Value)), nil
}

// -------------------------------------- asin -------------------------------------- MARK: asin

func nativeFn_asin(args []values.RuntimeValue, _ runtime.Context) (values.RuntimeValue, phpError.Error) {
	args, err := funcParamValidator.NewValidator("asin").AddParam("$num", []string{"float"}, nil).Validate(args)
	if err != nil {
		return values.NewVoid(), err
	}

	// Spec: https://www.php.net/manual/en/function.sin.php
	return values.NewFloat(goMath.Asin(args[0].(*values.Float).Value)), nil
}

// -------------------------------------- asinh -------------------------------------- MARK: asinh

func nativeFn_asinh(args []values.RuntimeValue, _ runtime.Context) (values.RuntimeValue, phpError.Error) {
	args, err := funcParamValidator.NewValidator("asinh").AddParam("$num", []string{"float"}, nil).Validate(args)
	if err != nil {
		return values.NewVoid(), err
	}

	// Spec: https://www.php.net/manual/en/function.sinh.php
	return values.NewFloat(goMath.Asinh(args[0].(*values.Float).Value)), nil
}

// -------------------------------------- pi -------------------------------------- MARK: pi

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
