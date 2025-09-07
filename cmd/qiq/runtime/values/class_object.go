package values

import (
	"QIQ/cmd/qiq/ast"
)

type Object struct {
	*abstractValue
	Class         *ast.ClassDeclarationStatement
	PropertyNames []string
	Properties    map[string]*Slot
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
		Properties:    map[string]*Slot{},
	}
}

func (object *Object) SetProperty(name string, value RuntimeValue) {
	if _, found := object.Properties[name]; found {
		object.Properties[name].Value = value
	} else {
		object.Properties[name] = NewSlot(value)
	}
}

func (object *Object) GetPropertySlot(name string) (*Slot, bool) {
	slot, found := object.Properties[name]
	if slot == nil || !found {
		return nil, false
	}
	return slot, true
}

func (object *Object) GetProperty(name string) (RuntimeValue, bool) {
	slot, found := object.GetPropertySlot(name)
	if slot == nil || !found {
		return NewNull(), false
	}
	return slot.Value, true
}

func (object *Object) GetMethod(name string) (*ast.MethodDefinitionStatement, bool) {
	method, found := object.Class.Methods[name]
	return method, found
}
