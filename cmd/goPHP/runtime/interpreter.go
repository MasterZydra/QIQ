package runtime

import (
	"GoPHP/cmd/goPHP/ast"
	"GoPHP/cmd/goPHP/ini"
	"GoPHP/cmd/goPHP/phpError"
	"GoPHP/cmd/goPHP/request"
	"GoPHP/cmd/goPHP/runtime/outputBuffer"
)

type Interpreter interface {
	GetRequest() *request.Request
	GetResponse() *request.Response
	GetIni() *ini.Ini
	GetOutputBufferStack() *outputBuffer.Stack
	GetClass(class string) (*ast.ClassDeclarationStatement, bool)
	Print(str string)
	Println(str string)
	PrintError(err phpError.Error)
	WriteResult(str string)
}
