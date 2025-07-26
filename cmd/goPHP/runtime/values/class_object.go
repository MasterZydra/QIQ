package values

import (
	"GoPHP/cmd/goPHP/ast"
)

type Object struct {
	*abstractValue
	Class      *ast.ClassDeclarationStatement
	Properties map[string]RuntimeValue
	// TODO methods
	// TODO parent
}

func NewObject(class *ast.ClassDeclarationStatement) *Object {
	return &Object{abstractValue: newAbstractValue(ObjectValue),
		Class: class, Properties: map[string]RuntimeValue{}}
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
	// TODO search parent classes if not contained in this one
	method, found := object.Class.Methods[name]
	return method, found
}
