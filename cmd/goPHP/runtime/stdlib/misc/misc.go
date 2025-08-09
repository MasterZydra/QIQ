package misc

import (
	"GoPHP/cmd/goPHP/phpError"
	"GoPHP/cmd/goPHP/runtime"
	"GoPHP/cmd/goPHP/runtime/funcParamValidator"
	"GoPHP/cmd/goPHP/runtime/values"
)

func Register(environment runtime.Environment) {
	// Category: Misc. Functions
	environment.AddNativeFunction("constant", nativeFn_constant)
	environment.AddNativeFunction("define", nativeFn_define)
	environment.AddNativeFunction("defined", nativeFn_defined)
}

// -------------------------------------- constant -------------------------------------- MARK: constant

func nativeFn_constant(args []values.RuntimeValue, context runtime.Context) (values.RuntimeValue, phpError.Error) {
	// Spec: https://www.php.net/manual/en/function.constant.php

	args, err := funcParamValidator.NewValidator("constant").AddParam("$name", []string{"string"}, nil).Validate(args)
	if err != nil {
		return values.NewVoid(), err
	}

	constantValue, err := context.Env.LookupConstant(args[0].(*values.Str).Value)
	if err != nil {
		return values.NewVoid(), err
	}

	return constantValue, nil
}

// -------------------------------------- define -------------------------------------- MARK: define

func nativeFn_define(args []values.RuntimeValue, context runtime.Context) (values.RuntimeValue, phpError.Error) {
	// Spec: https://www.php.net/manual/en/function.define.php

	args, err := funcParamValidator.NewValidator("define").
		AddParam("$constant_name", []string{"string"}, nil).
		AddParam("$value", []string{"mixed"}, nil).
		AddParam("$case_sensitive", []string{"bool"}, values.NewBool(false)).
		Validate(args)
	if err != nil {
		return values.NewVoid(), err
	}

	if args[2].(*values.Bool).Value {
		return values.NewBool(false), phpError.NewError("define: $case_sensitive cannot be changed")
	}

	context.Env.AddConstant(args[0].(*values.Str).Value, args[1])

	return values.NewBool(true), nil
}

// -------------------------------------- defined -------------------------------------- MARK: defined

func nativeFn_defined(args []values.RuntimeValue, context runtime.Context) (values.RuntimeValue, phpError.Error) {
	// Spec: https://www.php.net/manual/en/function.defined.php

	args, err := funcParamValidator.NewValidator("defined").AddParam("$name", []string{"string"}, nil).Validate(args)
	if err != nil {
		return values.NewVoid(), err
	}

	_, err = context.Env.LookupConstant(args[0].(*values.Str).Value)
	return values.NewBool(err == nil), nil
}

// TODO connection_​aborted
// TODO connection_​status
// TODO define
// TODO eval
// TODO get_​browser
// TODO _​_​halt_​compiler
// TODO highlight_​file
// TODO highlight_​string
// TODO hrtime
// TODO ignore_​user_​abort
// TODO pack
// TODO php_​strip_​whitespace
// TODO sapi_​windows_​cp_​conv
// TODO sapi_​windows_​cp_​get
// TODO sapi_​windows_​cp_​is_​utf8
// TODO sapi_​windows_​cp_​set
// TODO sapi_​windows_​generate_​ctrl_​event
// TODO sapi_​windows_​set_​ctrl_​handler
// TODO sapi_​windows_​vt100_​support
// TODO show_​source
// TODO sleep
// TODO sys_​getloadavg
// TODO time_​nanosleep
// TODO time_​sleep_​until
// TODO uniqid
// TODO unpack
// TODO usleep
