package runtime

import (
	"QIQ/cmd/qiq/ast"
	"strings"
)

type ExecutionContext struct {
	classNames            []string
	classDeclarations     map[string]*ast.ClassDeclarationStatement
	interfaceNames        []string
	interfaceDeclarations map[string]*ast.InterfaceDeclarationStatement
}

func NewExecutionContext() *ExecutionContext {
	return &ExecutionContext{
		classNames:            []string{},
		classDeclarations:     map[string]*ast.ClassDeclarationStatement{},
		interfaceNames:        []string{},
		interfaceDeclarations: map[string]*ast.InterfaceDeclarationStatement{},
	}
}

// -------------------------------------- Classes -------------------------------------- MARK: Classes

func (executionContext *ExecutionContext) AddClass(class string, classDecl *ast.ClassDeclarationStatement) {
	executionContext.classNames = append(executionContext.classNames, class)
	// TODO check if class already exists and return error that re-declaration is not possible
	executionContext.classDeclarations[strings.ToLower(class)] = classDecl
}

func (executionContext *ExecutionContext) GetClass(class string) (*ast.ClassDeclarationStatement, bool) {
	classDeclaration, found := executionContext.classDeclarations[strings.ToLower(class)]
	if !found {
		return nil, false
	}
	return classDeclaration, true
}

func (executionContext *ExecutionContext) GetClasses() []string {
	return executionContext.classNames
}

// -------------------------------------- Interface -------------------------------------- MARK: Interface

func (executionContext *ExecutionContext) AddInterface(interfaceName string, interfaceDecl *ast.InterfaceDeclarationStatement) {
	executionContext.interfaceNames = append(executionContext.interfaceNames, interfaceName)
	// TODO check if class already exists and return error that re-declaration is not possible
	executionContext.interfaceDeclarations[strings.ToLower(interfaceName)] = interfaceDecl
}

func (executionContext *ExecutionContext) GetInterface(interfaceName string) (*ast.InterfaceDeclarationStatement, bool) {
	interfaceDecl, found := executionContext.interfaceDeclarations[strings.ToLower(interfaceName)]
	if !found {
		return nil, false
	}
	return interfaceDecl, true
}

func (executionContext *ExecutionContext) GetInterfaces() []string {
	return executionContext.interfaceNames
}
