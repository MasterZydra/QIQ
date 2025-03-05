package runtime

type Context struct {
	Interpreter Interpreter
	Env         Environment
}

func NewContext(interpreter Interpreter, env Environment) Context {
	return Context{Interpreter: interpreter, Env: env}
}
