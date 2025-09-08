package ast

type AddGetMethod interface {
	AddMethod(method *MethodDefinitionStatement)
	GetMethod(name string) (*MethodDefinitionStatement, bool)
	// Required for error messages
	GetQualifiedName() string
}

type AddGetConst interface {
	AddConst(constStmt *ClassConstDeclarationStatement)
	GetConst(name string) (*ClassConstDeclarationStatement, bool)
	// Required for error messages
	GetQualifiedName() string
}
