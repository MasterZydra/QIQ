package interpreter

import (
	"GoPHP/cmd/goPHP/phpError"
	"GoPHP/cmd/goPHP/runtime/values"
)

func registerNativeMiscFunctions(environment *Environment) {
	environment.nativeFunctions["constant"] = nativeFn_constant
	environment.nativeFunctions["defined"] = nativeFn_defined
}

// ------------------- MARK: constant -------------------

func nativeFn_constant(args []values.RuntimeValue, interpreter *Interpreter) (values.RuntimeValue, phpError.Error) {
	// Spec: https://www.php.net/manual/en/function.constant.php

	args, err := NewFuncParamValidator("constant").addParam("$name", []string{"string"}, nil).validate(args)
	if err != nil {
		return values.NewVoid(), err
	}

	constantValue, err := interpreter.env.lookupConstant(args[0].(*values.Str).Value)
	if err != nil {
		return values.NewVoid(), err
	}

	return constantValue, nil
}

// ------------------- MARK: defined -------------------

func nativeFn_defined(args []values.RuntimeValue, interpreter *Interpreter) (values.RuntimeValue, phpError.Error) {
	// Spec: https://www.php.net/manual/en/function.defined.php

	args, err := NewFuncParamValidator("defined").addParam("$name", []string{"string"}, nil).validate(args)
	if err != nil {
		return values.NewVoid(), err
	}

	_, err = interpreter.env.lookupConstant(args[0].(*values.Str).Value)
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
