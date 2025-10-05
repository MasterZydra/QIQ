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

	// -------------------------------------- stdClass -------------------------------------- MARK: stdClass

	// Spec: https://www.php.net/manual/en/class.stdclass.php
	stdClass := ast.NewClassDeclarationStmt(0, nil, "stdClass", false, false)

	interpreter.AddClass(stdClass.Name, stdClass)
}
