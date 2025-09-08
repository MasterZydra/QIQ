package interfaces

import (
	"QIQ/cmd/qiq/ast"
	"QIQ/cmd/qiq/runtime"
)

func RegisterDefaultInterfaces(interpreter runtime.Interpreter) {
	// -------------------------------------- Traversable -------------------------------------- MARK: Traversable

	// Spec: https://www.php.net/manual/en/class.traversable.php
	Traversable := ast.NewInterfaceDeclarationStmt(0, nil, "Traversable")

	interpreter.AddInterface(Traversable.Name, Traversable)

	// -------------------------------------- IteratorAggregate -------------------------------------- MARK: IteratorAggregate

	// Spec: https://www.php.net/manual/en/class.iteratoraggregate.php
	IteratorAggregate := ast.NewInterfaceDeclarationStmt(0, nil, "IteratorAggregate")
	IteratorAggregate.Parents = append(IteratorAggregate.Parents, "Traversable")
	IteratorAggregate.AddMethod(ast.NewMethodDefinitionStmt(0, nil, "getIterator", []string{"public"}, []ast.FunctionParameter{}, nil, []string{"Traversable"}))

	interpreter.AddInterface(IteratorAggregate.Name, IteratorAggregate)

	// -------------------------------------- Iterator -------------------------------------- MARK: Iterator

	// Spec: https://www.php.net/manual/en/class.iterator.php
	Iterator := ast.NewInterfaceDeclarationStmt(0, nil, "Iterator")
	Iterator.Parents = append(Iterator.Parents, "Traversable")
	Iterator.AddMethod(ast.NewMethodDefinitionStmt(0, nil, "current", []string{"public"}, []ast.FunctionParameter{}, nil, []string{"mixed"}))
	Iterator.AddMethod(ast.NewMethodDefinitionStmt(0, nil, "key", []string{"public"}, []ast.FunctionParameter{}, nil, []string{"mixed"}))
	Iterator.AddMethod(ast.NewMethodDefinitionStmt(0, nil, "next", []string{"public"}, []ast.FunctionParameter{}, nil, []string{"void"}))
	Iterator.AddMethod(ast.NewMethodDefinitionStmt(0, nil, "rewind", []string{"public"}, []ast.FunctionParameter{}, nil, []string{"void"}))
	Iterator.AddMethod(ast.NewMethodDefinitionStmt(0, nil, "valid", []string{"public"}, []ast.FunctionParameter{}, nil, []string{"bool"}))

	interpreter.AddInterface(Iterator.Name, Iterator)

	// -------------------------------------- Serializable -------------------------------------- MARK: Serializable

	// Spec: https://www.php.net/manual/en/class.serializable.php
	Serializable := ast.NewInterfaceDeclarationStmt(0, nil, "Serializable")
	Serializable.AddMethod(ast.NewMethodDefinitionStmt(0, nil, "serialize", []string{"public"}, []ast.FunctionParameter{}, nil, []string{"null", "string"}))
	Serializable.AddMethod(ast.NewMethodDefinitionStmt(0, nil, "unserialize", []string{"public"}, []ast.FunctionParameter{ast.NewFunctionParam(false, "$data", []string{"string"}, nil)}, nil, []string{"void"}))

	interpreter.AddInterface(Serializable.Name, Serializable)

	// -------------------------------------- ArrayAccess -------------------------------------- MARK: ArrayAccess

	// Spec: https://www.php.net/manual/en/class.arrayaccess.php
	ArrayAccess := ast.NewInterfaceDeclarationStmt(0, nil, "ArrayAccess")
	ArrayAccess.AddMethod(ast.NewMethodDefinitionStmt(0, nil, "offsetExists", []string{"public"}, []ast.FunctionParameter{ast.NewFunctionParam(false, "$offset", []string{"mixed"}, nil)}, nil, []string{"bool"}))
	ArrayAccess.AddMethod(ast.NewMethodDefinitionStmt(0, nil, "offsetGet", []string{"public"}, []ast.FunctionParameter{ast.NewFunctionParam(false, "$offset", []string{"mixed"}, nil)}, nil, []string{"mixed"}))
	ArrayAccess.AddMethod(ast.NewMethodDefinitionStmt(0, nil, "offsetSet", []string{"public"}, []ast.FunctionParameter{ast.NewFunctionParam(false, "$offset", []string{"mixed"}, nil), {Name: "$value", Type: []string{"mixed"}}}, nil, []string{"void"}))
	ArrayAccess.AddMethod(ast.NewMethodDefinitionStmt(0, nil, "offsetUnset", []string{"public"}, []ast.FunctionParameter{ast.NewFunctionParam(false, "$offset", []string{"mixed"}, nil)}, nil, []string{"void"}))

	interpreter.AddInterface(ArrayAccess.Name, ArrayAccess)

	// -------------------------------------- Countable -------------------------------------- MARK: Countable

	// Spec: https://www.php.net/manual/en/class.countable.php
	Countable := ast.NewInterfaceDeclarationStmt(0, nil, "Countable")
	Countable.AddMethod(ast.NewMethodDefinitionStmt(0, nil, "count", []string{"public"}, []ast.FunctionParameter{}, nil, []string{"int"}))

	interpreter.AddInterface(Countable.Name, Countable)

	// -------------------------------------- Stringable -------------------------------------- MARK: Stringable

	// Spec: https://www.php.net/manual/en/class.stringable.php
	Stringable := ast.NewInterfaceDeclarationStmt(0, nil, "Stringable")
	Stringable.AddMethod(ast.NewMethodDefinitionStmt(0, nil, "__toString", []string{"public"}, []ast.FunctionParameter{}, nil, []string{"string"}))

	interpreter.AddInterface(Stringable.Name, Stringable)

	// -------------------------------------- Throwable -------------------------------------- MARK: Throwable

	// Spec: https://www.php.net/manual/en/class.throwable.php
	Throwable := ast.NewInterfaceDeclarationStmt(0, nil, "Throwable")
	Throwable.Parents = append(Throwable.Parents, "Stringable")
	Throwable.AddMethod(ast.NewMethodDefinitionStmt(0, nil, "getMessage", []string{"public"}, []ast.FunctionParameter{}, nil, []string{"string"}))
	Throwable.AddMethod(ast.NewMethodDefinitionStmt(0, nil, "getCode", []string{"public"}, []ast.FunctionParameter{}, nil, []string{"int"}))
	Throwable.AddMethod(ast.NewMethodDefinitionStmt(0, nil, "getFile", []string{"public"}, []ast.FunctionParameter{}, nil, []string{"string"}))
	Throwable.AddMethod(ast.NewMethodDefinitionStmt(0, nil, "getLine", []string{"public"}, []ast.FunctionParameter{}, nil, []string{"int"}))
	Throwable.AddMethod(ast.NewMethodDefinitionStmt(0, nil, "getTrace", []string{"public"}, []ast.FunctionParameter{}, nil, []string{"array"}))
	Throwable.AddMethod(ast.NewMethodDefinitionStmt(0, nil, "getTraceAsString", []string{"public"}, []ast.FunctionParameter{}, nil, []string{"string"}))
	Throwable.AddMethod(ast.NewMethodDefinitionStmt(0, nil, "getPrevious", []string{"public"}, []ast.FunctionParameter{}, nil, []string{"null", "Throwable"}))

	interpreter.AddInterface(Throwable.Name, Throwable)

	// -------------------------------------- UnitEnum -------------------------------------- MARK: UnitEnum

	// Spec: https://www.php.net/manual/en/class.unitenum.php
	UnitEnum := ast.NewInterfaceDeclarationStmt(0, nil, "UnitEnum")
	UnitEnum.AddMethod(ast.NewMethodDefinitionStmt(0, nil, "cases", []string{"public", "static"}, []ast.FunctionParameter{}, nil, []string{"array"}))

	interpreter.AddInterface(UnitEnum.Name, UnitEnum)
}
