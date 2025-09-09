package directory

import (
	"QIQ/cmd/qiq/phpError"
	"QIQ/cmd/qiq/runtime"
	"QIQ/cmd/qiq/runtime/funcParamValidator"
	"QIQ/cmd/qiq/runtime/values"
)

func Register(environment runtime.Environment) {
	// Category: Directory Functions
	environment.AddNativeFunction("getcwd", nativeFn_getcwd)
}

// -------------------------------------- getcwd -------------------------------------- MARK: getcwd

func nativeFn_getcwd(args []values.RuntimeValue, context runtime.Context) (values.RuntimeValue, phpError.Error) {
	// Spec: https://www.php.net/manual/en/function.getcwd.php

	_, err := funcParamValidator.NewValidator("getcwd").Validate(args)
	if err != nil {
		return values.NewVoid(), err
	}

	if context.Interpreter.GetFilename() == "" {
		return values.NewBool(false), nil
	}

	return values.NewStr(context.Interpreter.GetWorkingDir()), nil
}

// TODO chdir
// TODO chroot
// TODO closedir
// TODO dir
// TODO opendir
// TODO readdir
// TODO rewinddir
// TODO scandir
