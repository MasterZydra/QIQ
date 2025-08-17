package runtime

import (
	"QIQ/cmd/qiq/ast"
	"QIQ/cmd/qiq/ini"
	"QIQ/cmd/qiq/phpError"
	"QIQ/cmd/qiq/request"
	"QIQ/cmd/qiq/runtime/outputBuffer"
)

type Interpreter interface {
	GetRequest() *request.Request
	GetResponse() *request.Response
	GetIni() *ini.Ini
	GetOutputBufferStack() *outputBuffer.Stack
	AddClass(class string, classDecl *ast.ClassDeclarationStatement)
	GetClass(class string) (*ast.ClassDeclarationStatement, bool)
	Print(str string)
	Println(str string)
	PrintError(err phpError.Error)
	WriteResult(str string)
}
