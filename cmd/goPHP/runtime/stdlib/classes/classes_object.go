package classes

import (
	"GoPHP/cmd/goPHP/ast"
	"GoPHP/cmd/goPHP/phpError"
	"GoPHP/cmd/goPHP/runtime"
	"GoPHP/cmd/goPHP/runtime/funcParamValidator"
	"GoPHP/cmd/goPHP/runtime/values"
)

func Register(environment runtime.Environment) {
	// Category: Classes/Object Functions
	environment.AddNativeFunction("get_class", nativeFn_get_class)
	environment.AddNativeFunction("get_parent_class", nativeFn_get_parent_class)
	environment.AddNativeFunction("is_subclass_of", nativeFn_is_subclass_of)
}

// -------------------------------------- get_class -------------------------------------- MARK: get_class

func nativeFn_get_class(args []values.RuntimeValue, _ runtime.Context) (values.RuntimeValue, phpError.Error) {
	args, err := funcParamValidator.NewValidator("get_class").
		AddParam("$object", []string{"object"}, nil).
		Validate(args)
	if err != nil {
		return values.NewVoid(), err
	}

	// TODO return qualified namespaced name

	// Spec: https://www.php.net/manual/en/function.get-class.php
	// If the object is an instance of a class which exists in a namespace,
	// the qualified namespaced name of that class is returned.

	return values.NewStr(args[0].(*values.Object).Class.Name), nil
}

// -------------------------------------- get_parent_class -------------------------------------- MARK: get_class

func nativeFn_get_parent_class(args []values.RuntimeValue, context runtime.Context) (values.RuntimeValue, phpError.Error) {
	args, err := funcParamValidator.NewValidator("get_parent_class").
		AddParam("$object_or_class", []string{"object", "string"}, nil).
		Validate(args)
	if err != nil {
		return values.NewVoid(), err
	}

	objectOrClass := args[0]

	var class *ast.ClassDeclarationStatement = nil

	if objectOrClass.GetType() == values.StrValue {
		className := objectOrClass.(*values.Str).Value
		var found bool
		class, found = context.Interpreter.GetClass(className)
		if !found {
			return values.NewBool(false), nil
		}
	}

	if objectOrClass.GetType() == values.ObjectValue {
		class = objectOrClass.(*values.Object).Class
	}

	if class.BaseClass == "" {
		return values.NewBool(false), nil
	}
	return values.NewStr(class.BaseClass), nil
}

// -------------------------------------- is_subclass_of -------------------------------------- MARK: is_subclass_of

func nativeFn_is_subclass_of(args []values.RuntimeValue, context runtime.Context) (values.RuntimeValue, phpError.Error) {
	args, err := funcParamValidator.NewValidator("is_subclass_of").
		AddParam("$object_or_class", []string{"object", "string"}, nil).
		AddParam("$class", []string{"string"}, nil).
		AddParam("$allow_string", []string{"bool"}, values.NewBool(true)).
		Validate(args)
	if err != nil {
		return values.NewVoid(), err
	}

	objectOrClass := args[0]
	class := args[1].(*values.Str).Value
	allowString := args[2].(*values.Bool).Value

	if !allowString && objectOrClass.GetType() == values.StrValue {
		return values.NewVoid(), phpError.NewError("is_subclass_of: $object_or_class must be an object")
	}

	var classDecl *ast.ClassDeclarationStatement = nil

	if objectOrClass.GetType() == values.StrValue {
		className := objectOrClass.(*values.Str).Value
		var found bool
		classDecl, found = context.Interpreter.GetClass(className)
		if !found {
			return values.NewBool(false), nil
		}
	}

	if objectOrClass.GetType() == values.ObjectValue {
		classDecl = objectOrClass.(*values.Object).Class
	}

	// TODO is_subclass_of - check if it or parent implements it
	// TODO is_subclass_of - check if some parent is this class

	return values.NewBool(classDecl.BaseClass == class), nil
}
