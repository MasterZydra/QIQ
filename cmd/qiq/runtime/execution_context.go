package runtime

import (
	"QIQ/cmd/qiq/ast"
	"QIQ/cmd/qiq/runtime/values"
	"strings"
)

type ExecutionContext struct {
	// Classes
	classNames        []string
	classDeclarations map[string]*ast.ClassDeclarationStatement
	// Interfaces
	interfaceNames        []string
	interfaceDeclarations map[string]*ast.InterfaceDeclarationStatement
	// Objects
	objects map[string][]*values.Object
}

func NewExecutionContext() *ExecutionContext {
	return &ExecutionContext{
		// Classes
		classNames:        []string{},
		classDeclarations: map[string]*ast.ClassDeclarationStatement{},
		// Interfaces
		interfaceNames:        []string{},
		interfaceDeclarations: map[string]*ast.InterfaceDeclarationStatement{},
		// Objects
		objects: map[string][]*values.Object{},
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

// -------------------------------------- Interfaces -------------------------------------- MARK: Interfaces

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

// -------------------------------------- Objects -------------------------------------- MARK: Objects

func (executionContext *ExecutionContext) AddObject(className string, object *values.Object) {
	_, found := executionContext.objects[className]
	if !found {
		executionContext.objects[className] = []*values.Object{object}
		return
	}
	executionContext.objects[className] = append(executionContext.objects[className], object)
}

func (executionContext *ExecutionContext) CountObjects(className string) int {
	_, found := executionContext.objects[className]
	if !found {
		return 0
	}
	return len(executionContext.objects[className])
}
