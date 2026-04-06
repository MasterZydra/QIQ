package classes

import (
	"QIQ/cmd/qiq/ast"
	"QIQ/cmd/qiq/runtime"
)

func RegisterDefaultClasses(interpreter runtime.Interpreter) {
	// -------------------------------------- Exception -------------------------------------- MARK: Exception

	// Spec: https://www.php.net/manual/en/class.exception.php
	Exception := ast.NewClassDeclarationStmt(0, nil, "Exception", false, false)
	Exception.Interfaces = append(Exception.Interfaces, "Throwable")
	Exception.AddProperty(ast.NewPropertyDeclarationStmt(0, nil, "$message", "protected", false, []string{"string"}, ast.NewStringLiteralExpr(0, nil, "", ast.DoubleQuotedString)))
	Exception.AddProperty(ast.NewPropertyDeclarationStmt(0, nil, "$string", "private", false, []string{"string"}, ast.NewStringLiteralExpr(0, nil, "", ast.DoubleQuotedString)))
	Exception.AddProperty(ast.NewPropertyDeclarationStmt(0, nil, "$code", "protected", false, []string{"int"}, nil))
	Exception.AddProperty(ast.NewPropertyDeclarationStmt(0, nil, "$file", "protected", false, []string{"string"}, ast.NewStringLiteralExpr(0, nil, "", ast.DoubleQuotedString)))
	Exception.AddProperty(ast.NewPropertyDeclarationStmt(0, nil, "$line", "protected", false, []string{"int"}, nil))
	Exception.AddProperty(ast.NewPropertyDeclarationStmt(0, nil, "$trace", "private", false, []string{"array"}, ast.NewArrayLiteralExpr(0, nil)))
	Exception.AddProperty(ast.NewPropertyDeclarationStmt(0, nil, "$previous", "private", false, []string{"null", "Throwable"}, ast.NewConstantAccessExpr(0, nil, "NULL")))
	Exception.AddMethod(ast.NewMethodDefinitionStmt(0, nil, "__construct", []string{"public"}, []ast.FunctionParameter{{Name: "$message", Type: []string{"string"}, DefaultValue: ast.NewStringLiteralExpr(0, nil, "", ast.DoubleQuotedString)}, {Name: "$code", Type: []string{"int"}, DefaultValue: ast.NewIntegerLiteralExpr(0, nil, 0)}, {Name: "$previous", Type: []string{"null", "Throwable"}, DefaultValue: ast.NewConstantAccessExpr(0, nil, "NULL")}}, ast.NewCompoundStmt(0, []ast.IStatement{ast.NewExpressionStmt(0, ast.NewSimpleAssignmentExpr(0, ast.NewMemberAccessExpr(0, nil, ast.NewSimpleVariableExpr(0, ast.NewVariableNameExpr(0, nil, "$this")), ast.NewConstantAccessExpr(0, nil, "message")), ast.NewSimpleVariableExpr(0, ast.NewVariableNameExpr(0, nil, "$message")))), ast.NewExpressionStmt(0, ast.NewSimpleAssignmentExpr(0, ast.NewMemberAccessExpr(0, nil, ast.NewSimpleVariableExpr(0, ast.NewVariableNameExpr(0, nil, "$this")), ast.NewConstantAccessExpr(0, nil, "code")), ast.NewSimpleVariableExpr(0, ast.NewVariableNameExpr(0, nil, "$code")))), ast.NewExpressionStmt(0, ast.NewSimpleAssignmentExpr(0, ast.NewMemberAccessExpr(0, nil, ast.NewSimpleVariableExpr(0, ast.NewVariableNameExpr(0, nil, "$this")), ast.NewConstantAccessExpr(0, nil, "previous")), ast.NewSimpleVariableExpr(0, ast.NewVariableNameExpr(0, nil, "$previous"))))}), []string{}))
	Exception.AddMethod(ast.NewMethodDefinitionStmt(0, nil, "getMessage", []string{"public", "final"}, []ast.FunctionParameter{}, ast.NewCompoundStmt(0, []ast.IStatement{ast.NewReturnStmt(0, nil, ast.NewMemberAccessExpr(0, nil, ast.NewSimpleVariableExpr(0, ast.NewVariableNameExpr(0, nil, "$this")), ast.NewConstantAccessExpr(0, nil, "message")))}), []string{"string"}))
	Exception.AddMethod(ast.NewMethodDefinitionStmt(0, nil, "getPrevious", []string{"public", "final"}, []ast.FunctionParameter{}, ast.NewCompoundStmt(0, []ast.IStatement{ast.NewReturnStmt(0, nil, ast.NewMemberAccessExpr(0, nil, ast.NewSimpleVariableExpr(0, ast.NewVariableNameExpr(0, nil, "$this")), ast.NewConstantAccessExpr(0, nil, "previous")))}), []string{"null", "Throwable"}))
	Exception.AddMethod(ast.NewMethodDefinitionStmt(0, nil, "getCode", []string{"public", "final"}, []ast.FunctionParameter{}, ast.NewCompoundStmt(0, []ast.IStatement{ast.NewReturnStmt(0, nil, ast.NewMemberAccessExpr(0, nil, ast.NewSimpleVariableExpr(0, ast.NewVariableNameExpr(0, nil, "$this")), ast.NewConstantAccessExpr(0, nil, "code")))}), []string{"int"}))
	Exception.AddMethod(ast.NewMethodDefinitionStmt(0, nil, "getFile", []string{"public", "final"}, []ast.FunctionParameter{}, ast.NewCompoundStmt(0, []ast.IStatement{ast.NewReturnStmt(0, nil, ast.NewMemberAccessExpr(0, nil, ast.NewSimpleVariableExpr(0, ast.NewVariableNameExpr(0, nil, "$this")), ast.NewConstantAccessExpr(0, nil, "file")))}), []string{"string"}))
	Exception.AddMethod(ast.NewMethodDefinitionStmt(0, nil, "getLine", []string{"public", "final"}, []ast.FunctionParameter{}, ast.NewCompoundStmt(0, []ast.IStatement{ast.NewReturnStmt(0, nil, ast.NewMemberAccessExpr(0, nil, ast.NewSimpleVariableExpr(0, ast.NewVariableNameExpr(0, nil, "$this")), ast.NewConstantAccessExpr(0, nil, "line")))}), []string{"int"}))
	Exception.AddMethod(ast.NewMethodDefinitionStmt(0, nil, "getTrace", []string{"public", "final"}, []ast.FunctionParameter{}, ast.NewCompoundStmt(0, []ast.IStatement{ast.NewReturnStmt(0, nil, ast.NewMemberAccessExpr(0, nil, ast.NewSimpleVariableExpr(0, ast.NewVariableNameExpr(0, nil, "$this")), ast.NewConstantAccessExpr(0, nil, "trace")))}), []string{"array"}))
	Exception.AddMethod(ast.NewMethodDefinitionStmt(0, nil, "getTraceAsString", []string{"public", "final"}, []ast.FunctionParameter{}, ast.NewCompoundStmt(0, []ast.IStatement{ast.NewReturnStmt(0, nil, ast.NewFunctionCallExpr(0, nil, ast.NewStringLiteralExpr(0, nil, "implode", ast.DoubleQuotedString), []ast.IExpression{ast.NewConstantAccessExpr(0, nil, "PHP_EOL"), ast.NewMemberAccessExpr(0, nil, ast.NewSimpleVariableExpr(0, ast.NewVariableNameExpr(0, nil, "$this")), ast.NewConstantAccessExpr(0, nil, "trace"))}))}), []string{"string"}))
	Exception.AddMethod(ast.NewMethodDefinitionStmt(0, nil, "__toString", []string{"public"}, []ast.FunctionParameter{}, ast.NewCompoundStmt(0, []ast.IStatement{ast.NewReturnStmt(0, nil, ast.NewMemberAccessExpr(0, nil, ast.NewSimpleVariableExpr(0, ast.NewVariableNameExpr(0, nil, "$this")), ast.NewConstantAccessExpr(0, nil, "string")))}), []string{"string"}))
	Exception.AddMethod(ast.NewMethodDefinitionStmt(0, nil, "__clone", []string{"private"}, []ast.FunctionParameter{}, ast.NewCompoundStmt(0, []ast.IStatement{}), []string{"void"}))

	interpreter.AddClass(Exception.Name, Exception)

	// -------------------------------------- ErrorException -------------------------------------- MARK: ErrorException

	// Spec: https://www.php.net/manual/en/class.errorexception.php
	ErrorException := ast.NewClassDeclarationStmt(0, nil, "ErrorException", false, false)
	ErrorException.BaseClass = "Exception"
	ErrorException.AddProperty(ast.NewPropertyDeclarationStmt(0, nil, "$severity", "protected", false, []string{"int"}, ast.NewConstantAccessExpr(0, nil, "E_ERROR")))
	ErrorException.AddMethod(ast.NewMethodDefinitionStmt(0, nil, "__construct", []string{"public"}, []ast.FunctionParameter{{Name: "$message", Type: []string{"string"}, DefaultValue: ast.NewStringLiteralExpr(0, nil, "", ast.DoubleQuotedString)}, {Name: "$code", Type: []string{"int"}, DefaultValue: ast.NewIntegerLiteralExpr(0, nil, 0)}, {Name: "$severity", Type: []string{"int"}, DefaultValue: ast.NewConstantAccessExpr(0, nil, "E_ERROR")}, {Name: "$filename", Type: []string{"null", "string"}, DefaultValue: ast.NewConstantAccessExpr(0, nil, "NULL")}, {Name: "$line", Type: []string{"null", "int"}, DefaultValue: ast.NewConstantAccessExpr(0, nil, "NULL")}, {Name: "$previous", Type: []string{"null", "Throwable"}, DefaultValue: ast.NewConstantAccessExpr(0, nil, "NULL")}}, ast.NewCompoundStmt(0, []ast.IStatement{ast.NewExpressionStmt(0, ast.NewMemberAccessExpr(0, nil, ast.NewConstantAccessExpr(0, nil, "parent"), ast.NewFunctionCallExpr(0, nil, ast.NewStringLiteralExpr(0, nil, "__construct", ast.DoubleQuotedString), []ast.IExpression{ast.NewSimpleVariableExpr(0, ast.NewVariableNameExpr(0, nil, "$message")), ast.NewSimpleVariableExpr(0, ast.NewVariableNameExpr(0, nil, "$code")), ast.NewSimpleVariableExpr(0, ast.NewVariableNameExpr(0, nil, "$previous"))}))), ast.NewExpressionStmt(0, ast.NewSimpleAssignmentExpr(0, ast.NewMemberAccessExpr(0, nil, ast.NewSimpleVariableExpr(0, ast.NewVariableNameExpr(0, nil, "$this")), ast.NewConstantAccessExpr(0, nil, "severity")), ast.NewSimpleVariableExpr(0, ast.NewVariableNameExpr(0, nil, "$severity")))), ast.NewExpressionStmt(0, ast.NewSimpleAssignmentExpr(0, ast.NewMemberAccessExpr(0, nil, ast.NewSimpleVariableExpr(0, ast.NewVariableNameExpr(0, nil, "$this")), ast.NewConstantAccessExpr(0, nil, "file")), ast.NewCoalesceExpr(0, ast.NewSimpleVariableExpr(0, ast.NewVariableNameExpr(0, nil, "$filename")), ast.NewStringLiteralExpr(0, nil, "", ast.DoubleQuotedString)))), ast.NewExpressionStmt(0, ast.NewSimpleAssignmentExpr(0, ast.NewMemberAccessExpr(0, nil, ast.NewSimpleVariableExpr(0, ast.NewVariableNameExpr(0, nil, "$this")), ast.NewConstantAccessExpr(0, nil, "line")), ast.NewCoalesceExpr(0, ast.NewSimpleVariableExpr(0, ast.NewVariableNameExpr(0, nil, "$line")), ast.NewIntegerLiteralExpr(0, nil, 0))))}), []string{}))
	ErrorException.AddMethod(ast.NewMethodDefinitionStmt(0, nil, "getSeverity", []string{"public", "final"}, []ast.FunctionParameter{}, ast.NewCompoundStmt(0, []ast.IStatement{ast.NewReturnStmt(0, nil, ast.NewMemberAccessExpr(0, nil, ast.NewSimpleVariableExpr(0, ast.NewVariableNameExpr(0, nil, "$this")), ast.NewConstantAccessExpr(0, nil, "severity")))}), []string{"int"}))

	interpreter.AddClass(ErrorException.Name, ErrorException)

	// -------------------------------------- Error -------------------------------------- MARK: Error

	// Spec: https://www.php.net/manual/en/class.error.php
	Error := ast.NewClassDeclarationStmt(0, nil, "Error", false, false)
	Error.Interfaces = append(Error.Interfaces, "Throwable")
	Error.AddProperty(ast.NewPropertyDeclarationStmt(0, nil, "$message", "protected", false, []string{"string"}, ast.NewStringLiteralExpr(0, nil, "", ast.DoubleQuotedString)))
	Error.AddProperty(ast.NewPropertyDeclarationStmt(0, nil, "$string", "private", false, []string{"string"}, ast.NewStringLiteralExpr(0, nil, "", ast.DoubleQuotedString)))
	Error.AddProperty(ast.NewPropertyDeclarationStmt(0, nil, "$code", "protected", false, []string{"int"}, nil))
	Error.AddProperty(ast.NewPropertyDeclarationStmt(0, nil, "$file", "protected", false, []string{"string"}, ast.NewStringLiteralExpr(0, nil, "", ast.DoubleQuotedString)))
	Error.AddProperty(ast.NewPropertyDeclarationStmt(0, nil, "$line", "protected", false, []string{"int"}, nil))
	Error.AddProperty(ast.NewPropertyDeclarationStmt(0, nil, "$trace", "private", false, []string{"array"}, ast.NewArrayLiteralExpr(0, nil)))
	Error.AddProperty(ast.NewPropertyDeclarationStmt(0, nil, "$previous", "private", false, []string{"null", "Throwable"}, ast.NewConstantAccessExpr(0, nil, "NULL")))
	Error.AddMethod(ast.NewMethodDefinitionStmt(0, nil, "__construct", []string{"public"}, []ast.FunctionParameter{{Name: "$message", Type: []string{"string"}, DefaultValue: ast.NewStringLiteralExpr(0, nil, "", ast.DoubleQuotedString)}, {Name: "$code", Type: []string{"int"}, DefaultValue: ast.NewIntegerLiteralExpr(0, nil, 0)}, {Name: "$previous", Type: []string{"null", "Throwable"}, DefaultValue: ast.NewConstantAccessExpr(0, nil, "NULL")}}, ast.NewCompoundStmt(0, []ast.IStatement{ast.NewExpressionStmt(0, ast.NewSimpleAssignmentExpr(0, ast.NewMemberAccessExpr(0, nil, ast.NewSimpleVariableExpr(0, ast.NewVariableNameExpr(0, nil, "$this")), ast.NewConstantAccessExpr(0, nil, "message")), ast.NewSimpleVariableExpr(0, ast.NewVariableNameExpr(0, nil, "$message")))), ast.NewExpressionStmt(0, ast.NewSimpleAssignmentExpr(0, ast.NewMemberAccessExpr(0, nil, ast.NewSimpleVariableExpr(0, ast.NewVariableNameExpr(0, nil, "$this")), ast.NewConstantAccessExpr(0, nil, "code")), ast.NewSimpleVariableExpr(0, ast.NewVariableNameExpr(0, nil, "$code")))), ast.NewExpressionStmt(0, ast.NewSimpleAssignmentExpr(0, ast.NewMemberAccessExpr(0, nil, ast.NewSimpleVariableExpr(0, ast.NewVariableNameExpr(0, nil, "$this")), ast.NewConstantAccessExpr(0, nil, "previous")), ast.NewSimpleVariableExpr(0, ast.NewVariableNameExpr(0, nil, "$previous"))))}), []string{}))
	Error.AddMethod(ast.NewMethodDefinitionStmt(0, nil, "getMessage", []string{"public", "final"}, []ast.FunctionParameter{}, ast.NewCompoundStmt(0, []ast.IStatement{ast.NewReturnStmt(0, nil, ast.NewMemberAccessExpr(0, nil, ast.NewSimpleVariableExpr(0, ast.NewVariableNameExpr(0, nil, "$this")), ast.NewConstantAccessExpr(0, nil, "message")))}), []string{"string"}))
	Error.AddMethod(ast.NewMethodDefinitionStmt(0, nil, "getPrevious", []string{"public", "final"}, []ast.FunctionParameter{}, ast.NewCompoundStmt(0, []ast.IStatement{ast.NewReturnStmt(0, nil, ast.NewMemberAccessExpr(0, nil, ast.NewSimpleVariableExpr(0, ast.NewVariableNameExpr(0, nil, "$this")), ast.NewConstantAccessExpr(0, nil, "previous")))}), []string{"null", "Throwable"}))
	Error.AddMethod(ast.NewMethodDefinitionStmt(0, nil, "getCode", []string{"public", "final"}, []ast.FunctionParameter{}, ast.NewCompoundStmt(0, []ast.IStatement{ast.NewReturnStmt(0, nil, ast.NewMemberAccessExpr(0, nil, ast.NewSimpleVariableExpr(0, ast.NewVariableNameExpr(0, nil, "$this")), ast.NewConstantAccessExpr(0, nil, "code")))}), []string{"int"}))
	Error.AddMethod(ast.NewMethodDefinitionStmt(0, nil, "getFile", []string{"public", "final"}, []ast.FunctionParameter{}, ast.NewCompoundStmt(0, []ast.IStatement{ast.NewReturnStmt(0, nil, ast.NewMemberAccessExpr(0, nil, ast.NewSimpleVariableExpr(0, ast.NewVariableNameExpr(0, nil, "$this")), ast.NewConstantAccessExpr(0, nil, "file")))}), []string{"string"}))
	Error.AddMethod(ast.NewMethodDefinitionStmt(0, nil, "getLine", []string{"public", "final"}, []ast.FunctionParameter{}, ast.NewCompoundStmt(0, []ast.IStatement{ast.NewReturnStmt(0, nil, ast.NewMemberAccessExpr(0, nil, ast.NewSimpleVariableExpr(0, ast.NewVariableNameExpr(0, nil, "$this")), ast.NewConstantAccessExpr(0, nil, "line")))}), []string{"int"}))
	Error.AddMethod(ast.NewMethodDefinitionStmt(0, nil, "getTrace", []string{"public", "final"}, []ast.FunctionParameter{}, ast.NewCompoundStmt(0, []ast.IStatement{ast.NewReturnStmt(0, nil, ast.NewMemberAccessExpr(0, nil, ast.NewSimpleVariableExpr(0, ast.NewVariableNameExpr(0, nil, "$this")), ast.NewConstantAccessExpr(0, nil, "trace")))}), []string{"array"}))
	Error.AddMethod(ast.NewMethodDefinitionStmt(0, nil, "getTraceAsString", []string{"public", "final"}, []ast.FunctionParameter{}, ast.NewCompoundStmt(0, []ast.IStatement{ast.NewReturnStmt(0, nil, ast.NewFunctionCallExpr(0, nil, ast.NewStringLiteralExpr(0, nil, "implode", ast.DoubleQuotedString), []ast.IExpression{ast.NewConstantAccessExpr(0, nil, "PHP_EOL"), ast.NewMemberAccessExpr(0, nil, ast.NewSimpleVariableExpr(0, ast.NewVariableNameExpr(0, nil, "$this")), ast.NewConstantAccessExpr(0, nil, "trace"))}))}), []string{"string"}))
	Error.AddMethod(ast.NewMethodDefinitionStmt(0, nil, "__toString", []string{"public"}, []ast.FunctionParameter{}, ast.NewCompoundStmt(0, []ast.IStatement{ast.NewReturnStmt(0, nil, ast.NewMemberAccessExpr(0, nil, ast.NewSimpleVariableExpr(0, ast.NewVariableNameExpr(0, nil, "$this")), ast.NewConstantAccessExpr(0, nil, "string")))}), []string{"string"}))
	Error.AddMethod(ast.NewMethodDefinitionStmt(0, nil, "__clone", []string{"private"}, []ast.FunctionParameter{}, ast.NewCompoundStmt(0, []ast.IStatement{}), []string{"void"}))

	interpreter.AddClass(Error.Name, Error)

	// -------------------------------------- CompileError -------------------------------------- MARK: CompileError

	// Spec: https://www.php.net/manual/en/class.compileerror.php
	CompileError := ast.NewClassDeclarationStmt(0, nil, "CompileError", false, false)
	CompileError.BaseClass = "Error"

	interpreter.AddClass(CompileError.Name, CompileError)

	// -------------------------------------- ParseError -------------------------------------- MARK: ParseError

	// Spec: https://www.php.net/manual/en/class.parseerror.php
	ParseError := ast.NewClassDeclarationStmt(0, nil, "ParseError", false, false)
	ParseError.BaseClass = "CompileError"

	interpreter.AddClass(ParseError.Name, ParseError)

	// -------------------------------------- TypeError -------------------------------------- MARK: TypeError

	// Spec: https://www.php.net/manual/en/class.typeerror.php
	TypeError := ast.NewClassDeclarationStmt(0, nil, "TypeError", false, false)
	TypeError.BaseClass = "Error"

	interpreter.AddClass(TypeError.Name, TypeError)

	// -------------------------------------- ArgumentCountError -------------------------------------- MARK: ArgumentCountError

	// Spec: https://www.php.net/manual/en/class.argumentcounterror.php
	ArgumentCountError := ast.NewClassDeclarationStmt(0, nil, "ArgumentCountError", false, false)
	ArgumentCountError.BaseClass = "TypeError"

	interpreter.AddClass(ArgumentCountError.Name, ArgumentCountError)

	// -------------------------------------- ValueError -------------------------------------- MARK: ValueError

	// Spec: https://www.php.net/manual/en/class.valueerror.php
	ValueError := ast.NewClassDeclarationStmt(0, nil, "ValueError", false, false)
	ValueError.BaseClass = "Error"

	interpreter.AddClass(ValueError.Name, ValueError)

	// -------------------------------------- ArithmeticError -------------------------------------- MARK: ArithmeticError

	// Spec: https://www.php.net/manual/en/class.arithmeticerror.php

	ArithmeticError := ast.NewClassDeclarationStmt(0, nil, "ArithmeticError", false, false)
	ArithmeticError.BaseClass = "Error"

	interpreter.AddClass(ArithmeticError.Name, ArithmeticError)

	// -------------------------------------- DivisionByZeroError -------------------------------------- MARK: DivisionByZeroError

	// Spec: https://www.php.net/manual/en/class.divisionbyzeroerror.php
	DivisionByZeroError := ast.NewClassDeclarationStmt(0, nil, "DivisionByZeroError", false, false)
	DivisionByZeroError.BaseClass = "ArithmeticError"

	interpreter.AddClass(DivisionByZeroError.Name, DivisionByZeroError)

	// -------------------------------------- UnhandledMatchError -------------------------------------- MARK: UnhandledMatchError

	// Spec: https://www.php.net/manual/en/class.unhandledmatcherror.php
	UnhandledMatchError := ast.NewClassDeclarationStmt(0, nil, "UnhandledMatchError", false, false)
	UnhandledMatchError.BaseClass = "Error"

	interpreter.AddClass(UnhandledMatchError.Name, UnhandledMatchError)

	// -------------------------------------- RequestParseBodyException -------------------------------------- MARK: RequestParseBodyException

	// Spec: https://www.php.net/manual/en/class.requestparsebodyexception.php
	RequestParseBodyException := ast.NewClassDeclarationStmt(0, nil, "RequestParseBodyException", false, false)
	RequestParseBodyException.BaseClass = "Exception"

	interpreter.AddClass(RequestParseBodyException.Name, RequestParseBodyException)

	// -------------------------------------- ClosedGeneratorException -------------------------------------- MARK: ClosedGeneratorException

	// Spec: https://www.php.net/manual/en/class.closedgeneratorexception.php
	ClosedGeneratorException := ast.NewClassDeclarationStmt(0, nil, "ClosedGeneratorException", false, false)
	ClosedGeneratorException.BaseClass = "Exception"

	interpreter.AddClass(ClosedGeneratorException.Name, ClosedGeneratorException)

	// -------------------------------------- FiberError -------------------------------------- MARK: FiberError

	// Spec: https://www.php.net/manual/en/class.fibererror.php
	FiberError := ast.NewClassDeclarationStmt(0, nil, "FiberError", false, false)
	FiberError.BaseClass = "Error"

	interpreter.AddClass(FiberError.Name, FiberError)

	// -------------------------------------- stdClass -------------------------------------- MARK: stdClass

	// Spec: https://www.php.net/manual/en/class.stdclass.php
	stdClass := ast.NewClassDeclarationStmt(0, nil, "stdClass", false, false)

	interpreter.AddClass(stdClass.Name, stdClass)

	// -------------------------------------- JsonException -------------------------------------- MARK: JsonException

	// Spec: https://www.php.net/manual/en/class.jsonexception.php
	JsonException := ast.NewClassDeclarationStmt(0, nil, "JsonException", false, false)
	JsonException.BaseClass = "Exception"

	interpreter.AddClass(JsonException.Name, JsonException)

	// -------------------------------------- ReflectionException -------------------------------------- MARK: ReflectionException

	// Spec: https://www.php.net/manual/en/class.reflectionexception.php
	ReflectionException := ast.NewClassDeclarationStmt(0, nil, "ReflectionException", false, false)
	ReflectionException.BaseClass = "Exception"

	interpreter.AddClass(ReflectionException.Name, ReflectionException)

	// -------------------------------------- LogicException -------------------------------------- MARK: LogicException

	// Spec: https://www.php.net/manual/en/class.logicexception.php
	LogicException := ast.NewClassDeclarationStmt(0, nil, "LogicException", false, false)
	LogicException.BaseClass = "Exception"

	interpreter.AddClass(LogicException.Name, LogicException)

	// -------------------------------------- BadFunctionCallException -------------------------------------- MARK: BadFunctionCallException

	// Spec: https://www.php.net/manual/en/class.badfunctioncallexception.php
	BadFunctionCallException := ast.NewClassDeclarationStmt(0, nil, "BadFunctionCallException", false, false)
	BadFunctionCallException.BaseClass = "LogicException"

	interpreter.AddClass(BadFunctionCallException.Name, BadFunctionCallException)

	// -------------------------------------- BadMethodCallException -------------------------------------- MARK: BadMethodCallException

	// Spec: https://www.php.net/manual/en/class.badmethodcallexception.php
	BadMethodCallException := ast.NewClassDeclarationStmt(0, nil, "BadMethodCallException", false, false)
	BadMethodCallException.BaseClass = "BadFunctionCallException"

	interpreter.AddClass(BadMethodCallException.Name, BadMethodCallException)

	// -------------------------------------- DomainException -------------------------------------- MARK: DomainException

	// Spec: https://www.php.net/manual/en/class.domainexception.php
	DomainException := ast.NewClassDeclarationStmt(0, nil, "DomainException", false, false)
	DomainException.BaseClass = "LogicException"

	interpreter.AddClass(DomainException.Name, DomainException)

	// -------------------------------------- InvalidArgumentException -------------------------------------- MARK: InvalidArgumentException

	// Spec: https://www.php.net/manual/en/class.invalidargumentexception.php
	InvalidArgumentException := ast.NewClassDeclarationStmt(0, nil, "InvalidArgumentException", false, false)
	InvalidArgumentException.BaseClass = "LogicException"

	interpreter.AddClass(InvalidArgumentException.Name, InvalidArgumentException)

	// -------------------------------------- LengthException -------------------------------------- MARK: LengthException

	// Spec: https://www.php.net/manual/en/class.lengthexception.php
	LengthException := ast.NewClassDeclarationStmt(0, nil, "LengthException", false, false)
	LengthException.BaseClass = "LogicException"

	interpreter.AddClass(LengthException.Name, LengthException)

	// -------------------------------------- OutOfRangeException -------------------------------------- MARK: OutOfRangeException

	// Spec: https://www.php.net/manual/en/class.outofrangeexception.php
	OutOfRangeException := ast.NewClassDeclarationStmt(0, nil, "OutOfRangeException", false, false)
	OutOfRangeException.BaseClass = "LogicException"

	interpreter.AddClass(OutOfRangeException.Name, OutOfRangeException)

	// -------------------------------------- RuntimeException -------------------------------------- MARK: RuntimeException

	// Spec: https://www.php.net/manual/en/class.runtimeexception.php
	RuntimeException := ast.NewClassDeclarationStmt(0, nil, "RuntimeException", false, false)
	RuntimeException.BaseClass = "Exception"

	interpreter.AddClass(RuntimeException.Name, RuntimeException)

	// -------------------------------------- OutOfBoundsException -------------------------------------- MARK: OutOfBoundsException

	// Spec: https://www.php.net/manual/en/class.outofboundsexception.php
	OutOfBoundsException := ast.NewClassDeclarationStmt(0, nil, "OutOfBoundsException", false, false)
	OutOfBoundsException.BaseClass = "RuntimeException"

	interpreter.AddClass(OutOfBoundsException.Name, OutOfBoundsException)

	// -------------------------------------- OverflowException -------------------------------------- MARK: OverflowException

	// Spec: https://www.php.net/manual/en/class.overflowexception.php
	OverflowException := ast.NewClassDeclarationStmt(0, nil, "OverflowException", false, false)
	OverflowException.BaseClass = "RuntimeException"

	interpreter.AddClass(OverflowException.Name, OverflowException)

	// -------------------------------------- RangeException -------------------------------------- MARK: RangeException

	// Spec: https://www.php.net/manual/en/class.rangeexception.php
	RangeException := ast.NewClassDeclarationStmt(0, nil, "RangeException", false, false)
	RangeException.BaseClass = "RuntimeException"

	interpreter.AddClass(RangeException.Name, RangeException)

	// -------------------------------------- UnderflowException -------------------------------------- MARK: UnderflowException

	// Spec: https://www.php.net/manual/en/class.underflowexception.php
	UnderflowException := ast.NewClassDeclarationStmt(0, nil, "UnderflowException", false, false)
	UnderflowException.BaseClass = "RuntimeException"

	interpreter.AddClass(UnderflowException.Name, UnderflowException)

	// -------------------------------------- UnexpectedValueException -------------------------------------- MARK: UnexpectedValueException

	// Spec: https://www.php.net/manual/en/class.unexpectedvalueexception.php
	UnexpectedValueException := ast.NewClassDeclarationStmt(0, nil, "UnexpectedValueException", false, false)
	UnexpectedValueException.BaseClass = "RuntimeException"

	interpreter.AddClass(UnexpectedValueException.Name, UnexpectedValueException)

	// -------------------------------------- AssertionError -------------------------------------- MARK: AssertionError

	// Spec: https://www.php.net/manual/en/class.assertionerror.php
	AssertionError := ast.NewClassDeclarationStmt(0, nil, "AssertionError", false, false)
	AssertionError.BaseClass = "Error"

	interpreter.AddClass(AssertionError.Name, AssertionError)
}
