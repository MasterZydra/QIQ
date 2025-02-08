package interpreter

import (
	"GoPHP/cmd/goPHP/ast"
	"GoPHP/cmd/goPHP/ini"
	"GoPHP/cmd/goPHP/parser"
	"GoPHP/cmd/goPHP/phpError"
)

type Interpreter struct {
	filename      string
	includedFiles []string
	ini           *ini.Ini
	request       *Request
	parser        *parser.Parser
	env           *Environment
	cache         map[int64]IRuntimeValue
	result        string
	exitCode      int64
}

func NewInterpreter(ini *ini.Ini, request *Request, filename string) *Interpreter {
	interpreter := &Interpreter{
		filename: filename, includedFiles: []string{}, ini: ini, request: request, parser: parser.NewParser(ini),
		env: NewEnvironment(nil, request, ini), cache: map[int64]IRuntimeValue{},
		exitCode: 0,
	}

	if ini.GetBool("register_argc_argv") {
		server := interpreter.env.predefinedVariables["$_SERVER"].(*ArrayRuntimeValue)
		_, argvAlreadyDefined := server.findKey(NewStringRuntimeValue("argv"))
		_, argcAlreadyDefined := server.findKey(NewStringRuntimeValue("argc"))
		if !argcAlreadyDefined && !argvAlreadyDefined {
			server.SetElement(NewStringRuntimeValue("argv"), interpreter.env.predefinedVariables["$_GET"])
			server.SetElement(NewStringRuntimeValue("argc"), NewIntegerRuntimeValue(int64(len(interpreter.env.predefinedVariables["$_GET"].(*ArrayRuntimeValue).Keys))))
		}
	}

	return interpreter
}

func (interpreter *Interpreter) GetExitCode() int {
	return int(interpreter.exitCode)
}

func (interpreter *Interpreter) Process(sourceCode string) (string, phpError.Error) {
	return interpreter.process(sourceCode, interpreter.env)
}

func (interpreter *Interpreter) process(sourceCode string, env *Environment) (string, phpError.Error) {
	interpreter.result = ""
	program, parserErr := interpreter.parser.ProduceAST(sourceCode, interpreter.filename)
	if parserErr != nil {
		return interpreter.result, parserErr
	}

	_, err := interpreter.processProgram(program, env)

	return interpreter.result, err
}

func (interpreter *Interpreter) processProgram(program *ast.Program, env *Environment) (IRuntimeValue, phpError.Error) {
	err := interpreter.scanForFunctionDefinition(program.GetStatements(), env)
	if err != nil {
		return NewVoidRuntimeValue(), err
	}

	var runtimeValue IRuntimeValue = NewVoidRuntimeValue()
	for _, stmt := range program.GetStatements() {
		if runtimeValue, err = interpreter.processStmt(stmt, env); err != nil {
			// Handle exit event - Stop code execution
			if err.GetErrorType() == phpError.EventError && err.GetMessage() == phpError.ExitEvent {
				break
			}
			return runtimeValue, err
		}
	}
	return runtimeValue, nil
}

func (interpreter *Interpreter) processStmt(stmt ast.IStatement, env any) (value IRuntimeValue, phpErr phpError.Error) {
	defer func() {
		if r := recover(); r != nil {
			value = r.(ValueOrError).Value
			phpErr = r.(ValueOrError).Error.(phpError.Error)
		}
	}()
	runtimeValue, err := stmt.Process(interpreter, env)
	if err != nil {
		phpErr = err.(phpError.Error)
	}
	return runtimeValue.(IRuntimeValue), phpErr
}

type ValueOrError struct {
	Value IRuntimeValue
	Error error
}

func must(value IRuntimeValue, err error) IRuntimeValue {
	if err != nil {
		panic(ValueOrError{Value: value, Error: err})
	}
	return value
}

// Return VoidRuntimeValue if error is not nil
func mustOrVoid[V any](value V, err error) V {
	if err != nil {
		panic(ValueOrError{Value: NewVoidRuntimeValue(), Error: err})
	}
	return value
}
