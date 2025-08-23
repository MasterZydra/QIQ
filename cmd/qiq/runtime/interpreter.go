package runtime

import (
	"QIQ/cmd/qiq/ast"
	"QIQ/cmd/qiq/ini"
	"QIQ/cmd/qiq/phpError"
	"QIQ/cmd/qiq/request"
	"QIQ/cmd/qiq/runtime/outputBuffer"
	"QIQ/cmd/qiq/runtime/values"
)

type Interpreter interface {
	ProcessStatement(stmt ast.IStatement, env any) (values.RuntimeValue, phpError.Error)
	GetIni() *ini.Ini
	// Request
	GetRequest() *request.Request
	GetResponse() *request.Response
	// Class declarations
	AddClass(class string, classDecl *ast.ClassDeclarationStatement)
	GetClass(class string) (*ast.ClassDeclarationStatement, bool)
	// Interface declarations
	AddInterface(interfaceName string, interfaceDecl *ast.InterfaceDeclarationStatement)
	// Output
	GetOutputBufferStack() *outputBuffer.Stack
	Print(str string)
	Println(str string)
	PrintError(err phpError.Error)
	WriteResult(str string)
}
