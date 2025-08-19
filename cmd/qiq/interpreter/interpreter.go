package interpreter

import (
	"QIQ/cmd/qiq/ast"
	"QIQ/cmd/qiq/common"
	"QIQ/cmd/qiq/ini"
	"QIQ/cmd/qiq/parser"
	"QIQ/cmd/qiq/phpError"
	"QIQ/cmd/qiq/request"
	"QIQ/cmd/qiq/runtime/outputBuffer"
	"QIQ/cmd/qiq/runtime/values"
	"QIQ/cmd/qiq/stats"
)

var _ ast.Visitor = &Interpreter{}

type Interpreter struct {
	filename              string
	includedFiles         []string
	classDeclarations     map[string]*ast.ClassDeclarationStatement
	interfaceDeclarations map[string]*ast.InterfaceDeclarationStatement
	ini                   *ini.Ini
	request               *request.Request
	response              *request.Response
	parser                *parser.Parser
	env                   *Environment
	cache                 map[int64]values.RuntimeValue
	outputBufferStack     *outputBuffer.Stack
	result                string
	resultRuntimeValue    values.RuntimeValue
	// Status
	suppressWarning bool
	exitCalled      bool
}

func NewInterpreter(ini *ini.Ini, r *request.Request, filename string) (*Interpreter, phpError.Error) {
	interpreter := &Interpreter{
		filename:              filename,
		includedFiles:         []string{},
		classDeclarations:     map[string]*ast.ClassDeclarationStatement{},
		interfaceDeclarations: map[string]*ast.InterfaceDeclarationStatement{},
		ini:                   ini,
		request:               r,
		response:              request.NewResponse(),
		parser:                parser.NewParser(ini),
		cache:                 map[int64]values.RuntimeValue{},
		outputBufferStack:     outputBuffer.NewStack(),
	}

	var err phpError.Error
	interpreter.env, err = NewEnvironment(nil, r, interpreter)
	if err != nil {
		return interpreter, err
	}

	interpreter.AddClass("stdclass", ast.NewClassDeclarationStmt(0, nil, "stdClass", false, false))

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

func (interpreter *Interpreter) GetResponse() *request.Response {
	return interpreter.response
}

func (interpreter *Interpreter) GetOutputBufferStack() *outputBuffer.Stack {
	return interpreter.outputBufferStack
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

	if !interpreter.exitCalled {
		interpreter.ProcessExitIntrinsicExpr(ast.NewExitIntrinsic(0, nil, ast.NewIntegerLiteralExpr(0, nil, 0)), interpreter.env)
	}

	return runtimeValue, nil
}

func (interperter *Interpreter) ProcessStatement(stmt ast.IStatement, env any) (values.RuntimeValue, phpError.Error) {
	return interperter.processStmt(stmt, env)
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

	// Destruct unused objects
	if stmt.GetKind() != ast.ObjectCreationExpr {
		for index, object := range env.(*Environment).objects {
			if object.IsUsed {
				continue
			}
			if err := interpreter.destructObject(object, env.(*Environment)); err != nil {
				interpreter.PrintError(err)
			}
			env.(*Environment).objects = common.RemoveIndex(env.(*Environment).objects, index)
		}
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
