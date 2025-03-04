package interpreter

import (
	"GoPHP/cmd/goPHP/phpError"
	"GoPHP/cmd/goPHP/runtime/values"
	"testing"
)

func TestTooManyParams(t *testing.T) {
	validator := NewFuncParamValidator("testFn")
	_, err := validator.validate([]values.RuntimeValue{values.NewInt(42)})
	expectedErr := phpError.NewError("Uncaught ArgumentCountError: testFn() expects exactly 0 argument, 1 given")
	if err.GetMessage() != expectedErr.GetMessage() {
		t.Errorf("\nExpected: \"%s\"\nGot: \"%s\"", expectedErr, err)
	}

	validator = NewFuncParamValidator("testFn")
	_, err = validator.validate([]values.RuntimeValue{values.NewInt(42), values.NewInt(43)})
	expectedErr = phpError.NewError("Uncaught ArgumentCountError: testFn() expects exactly 0 argument, 2 given")
	if err.GetMessage() != expectedErr.GetMessage() {
		t.Errorf("\nExpected: \"%s\"\nGot: \"%s\"", expectedErr, err)
	}

	validator = NewFuncParamValidator("testFn").addParam("paramA", []string{"int"}, nil)
	_, err = validator.validate([]values.RuntimeValue{values.NewInt(42), values.NewInt(43)})
	expectedErr = phpError.NewError("Uncaught ArgumentCountError: testFn() expects exactly 1 argument, 2 given")
	if err.GetMessage() != expectedErr.GetMessage() {
		t.Errorf("\nExpected: \"%s\"\nGot: \"%s\"", expectedErr, err)
	}

	validator = NewFuncParamValidator("testFn").addParam("paramA", []string{"int"}, values.NewInt(0))
	_, err = validator.validate([]values.RuntimeValue{values.NewInt(42), values.NewInt(43)})
	expectedErr = phpError.NewError("Uncaught ArgumentCountError: testFn() expects most 1 argument, 2 given")
	if err.GetMessage() != expectedErr.GetMessage() {
		t.Errorf("\nExpected: \"%s\"\nGot: \"%s\"", expectedErr, err)
	}
}

func TestTooFewParams(t *testing.T) {
	validator := NewFuncParamValidator("testFn").addParam("paramA", []string{"int"}, nil)
	_, err := validator.validate([]values.RuntimeValue{})
	expectedErr := phpError.NewError("Uncaught ArgumentCountError: Too few arguments to function testFn(), 0 passed and at least 1 expected")
	if err.GetMessage() != expectedErr.GetMessage() {
		t.Errorf("\nExpected: \"%s\"\nGot: \"%s\"", expectedErr, err)
	}
}

func TestWrongParamType(t *testing.T) {
	validator := NewFuncParamValidator("testFn").addParam("paramA", []string{"int"}, nil)
	_, err := validator.validate([]values.RuntimeValue{values.NewStr("abc")})
	expectedErr := phpError.NewError("Uncaught TypeError: testFn(): Argument #1 (paramA) must be of type int, string given")
	if err.GetMessage() != expectedErr.GetMessage() {
		t.Errorf("\nExpected: \"%s\"\nGot: \"%s\"", expectedErr, err)
	}

	validator = NewFuncParamValidator("testFn").addParam("paramA", []string{"int", "float"}, nil)
	_, err = validator.validate([]values.RuntimeValue{values.NewStr("abc")})
	expectedErr = phpError.NewError("Uncaught TypeError: testFn(): Argument #1 (paramA) must be of type int|float, string given")
	if err.GetMessage() != expectedErr.GetMessage() {
		t.Errorf("\nExpected: \"%s\"\nGot: \"%s\"", expectedErr, err)
	}

	validator = NewFuncParamValidator("testFn").addVariableLenParam("paramA", []string{"int", "float"})
	_, err = validator.validate([]values.RuntimeValue{values.NewInt(42), values.NewStr("abc")})
	expectedErr = phpError.NewError("Uncaught TypeError: testFn(): Argument #1 (paramA) must be of type int|float, string given")
	if err == nil || err.GetMessage() != expectedErr.GetMessage() {
		t.Errorf("\nExpected: \"%s\"\nGot: \"%s\"", expectedErr, err)
	}
}

func TestCorrectParamType(t *testing.T) {
	validator := NewFuncParamValidator("testFn").addParam("paramA", []string{"int"}, nil)
	_, err := validator.validate([]values.RuntimeValue{values.NewInt(42)})
	if err != nil {
		t.Errorf("\nExpected: nil\nGot: \"%s\"", err)
	}

	validator = NewFuncParamValidator("testFn").addParam("paramA", []string{"int"}, values.NewInt(42))
	_, err = validator.validate([]values.RuntimeValue{values.NewInt(42)})
	if err != nil {
		t.Errorf("\nExpected: nil\nGot: \"%s\"", err)
	}

	validator = NewFuncParamValidator("testFn").addParam("paramA", []string{"int"}, values.NewInt(42))
	_, err = validator.validate([]values.RuntimeValue{})
	if err != nil {
		t.Errorf("\nExpected: nil\nGot: \"%s\"", err)
	}

	validator = NewFuncParamValidator("testFn").
		addParam("paramA", []string{"int"}, nil).
		addParam("paramB", []string{"int"}, values.NewInt(42))
	_, err = validator.validate([]values.RuntimeValue{values.NewInt(0)})
	if err != nil {
		t.Errorf("\nExpected: nil\nGot: \"%s\"", err)
	}

	validator = NewFuncParamValidator("testFn").
		addParam("paramA", []string{"int"}, nil).
		addVariableLenParam("paramB", []string{"string"})
	_, err = validator.validate([]values.RuntimeValue{values.NewInt(42), values.NewStr("abc"), values.NewStr("abc")})
	if err != nil {
		t.Errorf("\nExpected: nil\nGot: \"%s\"", err)
	}

	validator = NewFuncParamValidator("testFn").
		addVariableLenParam("paramA", []string{"mixed"})
	_, err = validator.validate([]values.RuntimeValue{values.NewInt(42), values.NewBool(true), values.NewStr("abc")})
	if err != nil {
		t.Errorf("\nExpected: nil\nGot: \"%s\"", err)
	}
}
