package runtime

import "QIQ/cmd/qiq/ast"

type Context struct {
	Interpreter Interpreter
	Env         Environment
	Stmt        ast.IStatement
}

func NewContext(interpreter Interpreter, env Environment, stmt ast.IStatement) Context {
	return Context{Interpreter: interpreter, Env: env, Stmt: stmt}
}
