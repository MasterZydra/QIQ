package runtime

import (
	"QIQ/cmd/qiq/phpError"
	"QIQ/cmd/qiq/runtime/values"
)

type NativeFunction func([]values.RuntimeValue, Context) (values.RuntimeValue, phpError.Error)
