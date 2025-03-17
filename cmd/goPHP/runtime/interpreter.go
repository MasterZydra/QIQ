package runtime

import (
	"GoPHP/cmd/goPHP/ini"
	"GoPHP/cmd/goPHP/phpError"
	"GoPHP/cmd/goPHP/runtime/outputBuffer"
)

type Interpreter interface {
	GetIni() *ini.Ini
	GetOutputBufferStack() *outputBuffer.Stack
	Print(str string)
	Println(str string)
	PrintError(err phpError.Error)
	WriteResult(str string)
}
