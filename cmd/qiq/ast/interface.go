package ast

type AddGetMethod interface {
	AddMethod(method *MethodDefinitionStatement)
	GetMethod(name string) (*MethodDefinitionStatement, bool)
	// Required for error messages
	GetQualifiedName() string
}

type AddConst interface {
	AddConst(constStmt *ClassConstDeclarationStatement)
}
