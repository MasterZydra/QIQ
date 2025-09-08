package functionHandling

import (
	"QIQ/cmd/qiq/phpError"
	"QIQ/cmd/qiq/runtime"
	"QIQ/cmd/qiq/runtime/funcParamValidator"
	"QIQ/cmd/qiq/runtime/values"
)

func Register(environment runtime.Environment) {
	// Category: Function Handling Functions
	environment.AddNativeFunction("function_exists", nativeFn_function_exists)
}

// -------------------------------------- function_exists -------------------------------------- MARK: function_exists

func nativeFn_function_exists(args []values.RuntimeValue, context runtime.Context) (values.RuntimeValue, phpError.Error) {
	// Spec: https://www.php.net/manual/function.function-exists.php

	args, err := funcParamValidator.NewValidator("function_exists").AddParam("$function", []string{"string"}, nil).Validate(args)
	if err != nil {
		return values.NewVoid(), err
	}

	return values.NewBool(context.Env.FunctionExists(args[0].(*values.Str).Value)), nil
}
