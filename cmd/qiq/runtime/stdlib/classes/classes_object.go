package classes

import (
	"QIQ/cmd/qiq/ast"
	"QIQ/cmd/qiq/phpError"
	"QIQ/cmd/qiq/runtime"
	"QIQ/cmd/qiq/runtime/funcParamValidator"
	"QIQ/cmd/qiq/runtime/values"
	"strings"
)

func Register(environment runtime.Environment) {
	// Category: Classes/Object Functions
	environment.AddNativeFunction("class_exists", nativeFn_class_exists)
	environment.AddNativeFunction("get_class", nativeFn_get_class)
	environment.AddNativeFunction("get_parent_class", nativeFn_get_parent_class)
	environment.AddNativeFunction("is_subclass_of", nativeFn_is_subclass_of)
}

// -------------------------------------- class_exists -------------------------------------- MARK: class_exists

func nativeFn_class_exists(args []values.RuntimeValue, context runtime.Context) (values.RuntimeValue, phpError.Error) {
	// https://www.php.net/manual/en/function.class-exists.php

	args, err := funcParamValidator.NewValidator("class_exists").
		AddParam("$class", []string{"string"}, nil).
		Validate(args)
	if err != nil {
		return values.NewVoid(), err
	}

	_, found := context.Interpreter.GetClass(args[0].(*values.Str).Value)

	return values.NewBool(found), nil
}

// -------------------------------------- get_class -------------------------------------- MARK: get_class

func nativeFn_get_class(args []values.RuntimeValue, _ runtime.Context) (values.RuntimeValue, phpError.Error) {
	// Spec: https://www.php.net/manual/en/function.get-class.php

	args, err := funcParamValidator.NewValidator("get_class").
		AddParam("$object", []string{"object"}, nil).
		Validate(args)
	if err != nil {
		return values.NewVoid(), err
	}

	object := args[0].(*values.Object)

	// Spec: https://www.php.net/manual/en/function.get-class.php
	// If the object is an instance of a class which exists in a namespace,
	// the qualified namespaced name of that class is returned.
	namespace := object.Class.GetPosition().File.GetNamespaceStr()
	return values.NewStr(namespace + object.Class.Name), nil
}

// -------------------------------------- get_parent_class -------------------------------------- MARK: get_class

func nativeFn_get_parent_class(args []values.RuntimeValue, context runtime.Context) (values.RuntimeValue, phpError.Error) {
	// Spec: https://www.php.net/manual/en/function.get-parent-class.php

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
	// Spec: https://www.php.net/manual/en/function.is-subclass-of.php

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

	return values.NewBool(strings.EqualFold(classDecl.BaseClass, class)), nil
}

// TODO class_alias
// TODO enum_exists
// TODO get_called_class
// TODO get_class_methods
// TODO get_class_vars
// TODO get_declared_classes
// TODO get_declared_interfaces
// TODO get_declared_traits
// TODO get_mangled_object_vars
// TODO get_object_vars
// TODO interface_exists
// TODO is_a
// TODO method_exists
// TODO property_exists
// TODO trait_exists
