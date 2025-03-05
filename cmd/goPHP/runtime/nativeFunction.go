package runtime

import (
	"GoPHP/cmd/goPHP/phpError"
	"GoPHP/cmd/goPHP/runtime/values"
)

type NativeFunction func([]values.RuntimeValue, Context) (values.RuntimeValue, phpError.Error)
