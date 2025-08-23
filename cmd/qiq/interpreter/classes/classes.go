package classes

import (
	"QIQ/cmd/qiq/ast"
	"QIQ/cmd/qiq/runtime"
)

func RegisterDefaultClasses(interpreter runtime.Interpreter) {
	// -------------------------------------- stdClass -------------------------------------- MARK: stdClass

	// Spec: https://www.php.net/manual/en/class.stdclass.php
	stdClass := ast.NewClassDeclarationStmt(0, nil, "stdClass", false, false)

	interpreter.AddClass(stdClass.Name, stdClass)
}
