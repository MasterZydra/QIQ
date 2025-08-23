package values

import (
	"QIQ/cmd/qiq/ast"
)

type Object struct {
	*abstractValue
	Class         *ast.ClassDeclarationStatement
	PropertyNames []string
	Properties    map[string]RuntimeValue
	// TODO methods
	// TODO parent
	// Status
	IsUsed       bool
	IsDestructed bool
}

func NewObject(class *ast.ClassDeclarationStatement) *Object {
	return &Object{abstractValue: newAbstractValue(ObjectValue),
		Class:         class,
		PropertyNames: append([]string(nil), class.PropertieNames...),
		Properties:    map[string]RuntimeValue{},
	}
}

func (object *Object) SetProperty(name string, value RuntimeValue) {
	object.Properties[name] = value
}

func (object *Object) GetProperty(name string) (RuntimeValue, bool) {
	value, found := object.Properties[name]
	if !found {
		return NewNull(), false
	}
	return value, true
}

func (object *Object) GetMethod(name string) (*ast.MethodDefinitionStatement, bool) {
	method, found := object.Class.Methods[name]
	return method, found
}
