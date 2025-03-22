package interpreter

import (
	"GoPHP/cmd/goPHP/ast"
	"GoPHP/cmd/goPHP/ini"
	"GoPHP/cmd/goPHP/parser"
	"GoPHP/cmd/goPHP/phpError"
	"GoPHP/cmd/goPHP/request"
	"GoPHP/cmd/goPHP/runtime/outputBuffer"
	"GoPHP/cmd/goPHP/runtime/values"
	"GoPHP/cmd/goPHP/stats"
)

type Interpreter struct {
	filename           string
	includedFiles      []string
	ini                *ini.Ini
	request            *request.Request
	parser             *parser.Parser
	env                *Environment
	cache              map[int64]values.RuntimeValue
	outputBufferStack  *outputBuffer.Stack
	result             string
	resultRuntimeValue values.RuntimeValue
	exitCode           int64
	// Status
	suppressWarning bool
}

func NewInterpreter(ini *ini.Ini, request *request.Request, filename string) (*Interpreter, phpError.Error) {
	interpreter := &Interpreter{
		filename: filename, includedFiles: []string{}, ini: ini, request: request, parser: parser.NewParser(ini),
		cache:             map[int64]values.RuntimeValue{},
		outputBufferStack: outputBuffer.NewStack(),
		exitCode:          0,
	}
	var err phpError.Error
	interpreter.env, err = NewEnvironment(nil, request, interpreter)
	if err != nil {
		return interpreter, err
	}

	if ini.GetBool("register_argc_argv") {
		server := interpreter.env.predefinedVariables["$_SERVER"].(*values.Array)
		if !server.Contains(values.NewStr("argc")) && !server.Contains(values.NewStr("argv")) {
			server.SetElement(values.NewStr("argv"), interpreter.env.predefinedVariables["$_GET"])
			server.SetElement(values.NewStr("argc"), values.NewInt(int64(len(interpreter.env.predefinedVariables["$_GET"].(*values.Array).Keys))))
		}
	}

	return interpreter, nil
}

func (interpreter *Interpreter) GetIni() *ini.Ini {
	return interpreter.ini
}

func (interpreter *Interpreter) GetRequest() *request.Request {
	return interpreter.request
}

func (interpreter *Interpreter) GetOutputBufferStack() *outputBuffer.Stack {
	return interpreter.outputBufferStack
}

func (interpreter *Interpreter) GetExitCode() int {
	return int(interpreter.exitCode)
}

func (interpreter *Interpreter) Process(sourceCode string) (string, phpError.Error) {
	return interpreter.process(sourceCode, interpreter.env, false)
}

func (interpreter *Interpreter) process(sourceCode string, env *Environment, resetResult bool) (string, phpError.Error) {
	if resetResult {
		interpreter.result = ""
	}

	program, err := interpreter.parser.ProduceAST(sourceCode, interpreter.filename)
	if err != nil {
		return interpreter.result, err
	}

	stat := stats.Start()
	defer stats.StopAndPrint(stat, "Interpreter")

	parser.PrintParserCallstack("", nil)
	parser.PrintParserCallstack("Interpreter callstack", nil)
	parser.PrintParserCallstack("---------------------", nil)

	interpreter.resultRuntimeValue = nil
	interpreter.resultRuntimeValue, err = interpreter.processProgram(program, env)

	return interpreter.result, err
}

func (interpreter *Interpreter) processProgram(program *ast.Program, env *Environment) (values.RuntimeValue, phpError.Error) {
	err := interpreter.scanForFunctionDefinition(program.GetStatements(), env)
	if err != nil {
		return values.NewVoid(), err
	}

	defer interpreter.flushOutputBuffers()

	var runtimeValue values.RuntimeValue = values.NewVoid()
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

func (interpreter *Interpreter) processStmt(stmt ast.IStatement, env any) (value values.RuntimeValue, phpErr phpError.Error) {
	defer func() {
		if r := recover(); r != nil {
			value = r.(ValueOrError).Value
			phpErr = r.(ValueOrError).Error.(phpError.Error)
		}
	}()

	ast.PrintInterpreterCallstack(stmt)
	runtimeValue, err := stmt.Process(interpreter, env)
	if err != nil {
		phpErr = err.(phpError.Error)
	}
	return runtimeValue.(values.RuntimeValue), phpErr
}

type ValueOrError struct {
	Value values.RuntimeValue
	Error error
}

func must(value values.RuntimeValue, err error) values.RuntimeValue {
	if err != nil {
		panic(ValueOrError{Value: value, Error: err})
	}
	return value
}

// Return Void if error is not nil
func mustOrVoid[V any](value V, err error) V {
	if err != nil {
		panic(ValueOrError{Value: values.NewVoid(), Error: err})
	}
	return value
}
