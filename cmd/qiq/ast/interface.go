package ast

type AddMethod interface {
	AddMethod(method *MethodDefinitionStatement)
}

type AddConst interface {
	AddConst(constStmt *ClassConstDeclarationStatement)
}
