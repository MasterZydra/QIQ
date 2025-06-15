package values

import "GoPHP/cmd/goPHP/ast"

type Object struct {
	class *ast.ClassDeclarationStatement
	*abstractValue
	// TODO properties
	// TODO methods
	// TODO parent
}

func NewObject(class *ast.ClassDeclarationStatement) *Object {
	return &Object{abstractValue: newAbstractValue(ObjectValue), class: class}
}
