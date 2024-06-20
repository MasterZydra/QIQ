package interpreter

import (
	"GoPHP/cmd/goPHP/phpError"
	"testing"
)

func TestTooManyParams(t *testing.T) {
	validator := NewFuncParamValidator("testFn")
	_, err := validator.validate([]IRuntimeValue{NewIntegerRuntimeValue(42), NewIntegerRuntimeValue(43)})
	expectedErr := phpError.NewError("Uncaught ArgumentCountError: testFn() expects exactly 0 argument, 2 given")
	if err.GetMessage() != expectedErr.GetMessage() {
		t.Errorf("\nExpected: \"%s\"\nGot: \"%s\"", expectedErr, err)
	}

	validator = NewFuncParamValidator("testFn").addParam("paramA", []string{"int"}, nil)
	_, err = validator.validate([]IRuntimeValue{NewIntegerRuntimeValue(42), NewIntegerRuntimeValue(43)})
	expectedErr = phpError.NewError("Uncaught ArgumentCountError: testFn() expects exactly 1 argument, 2 given")
	if err.GetMessage() != expectedErr.GetMessage() {
		t.Errorf("\nExpected: \"%s\"\nGot: \"%s\"", expectedErr, err)
	}

	validator = NewFuncParamValidator("testFn").addParam("paramA", []string{"int"}, NewIntegerRuntimeValue(0))
	_, err = validator.validate([]IRuntimeValue{NewIntegerRuntimeValue(42), NewIntegerRuntimeValue(43)})
	expectedErr = phpError.NewError("Uncaught ArgumentCountError: testFn() expects most 1 argument, 2 given")
	if err.GetMessage() != expectedErr.GetMessage() {
		t.Errorf("\nExpected: \"%s\"\nGot: \"%s\"", expectedErr, err)
	}
}

func TestTooFewParams(t *testing.T) {
	validator := NewFuncParamValidator("testFn").addParam("paramA", []string{"int"}, nil)
	_, err := validator.validate([]IRuntimeValue{})
	expectedErr := phpError.NewError("Uncaught ArgumentCountError: Too few arguments to function testFn(), 0 passed and at least 1 expected")
	if err.GetMessage() != expectedErr.GetMessage() {
		t.Errorf("\nExpected: \"%s\"\nGot: \"%s\"", expectedErr, err)
	}
}

func TestWrongParamType(t *testing.T) {
	validator := NewFuncParamValidator("testFn").addParam("paramA", []string{"int"}, nil)
	_, err := validator.validate([]IRuntimeValue{NewStringRuntimeValue("abc")})
	expectedErr := phpError.NewError("Uncaught TypeError: testFn(): Argument #1 (paramA) must be of type int, string given")
	if err.GetMessage() != expectedErr.GetMessage() {
		t.Errorf("\nExpected: \"%s\"\nGot: \"%s\"", expectedErr, err)
	}

	validator = NewFuncParamValidator("testFn").addParam("paramA", []string{"int", "float"}, nil)
	_, err = validator.validate([]IRuntimeValue{NewStringRuntimeValue("abc")})
	expectedErr = phpError.NewError("Uncaught TypeError: testFn(): Argument #1 (paramA) must be of type int|float, string given")
	if err.GetMessage() != expectedErr.GetMessage() {
		t.Errorf("\nExpected: \"%s\"\nGot: \"%s\"", expectedErr, err)
	}

	validator = NewFuncParamValidator("testFn").addVariableLenParam("paramA", []string{"int", "float"})
	_, err = validator.validate([]IRuntimeValue{NewIntegerRuntimeValue(42), NewStringRuntimeValue("abc")})
	expectedErr = phpError.NewError("Uncaught TypeError: testFn(): Argument #1 (paramA) must be of type int|float, string given")
	if err == nil || err.GetMessage() != expectedErr.GetMessage() {
		t.Errorf("\nExpected: \"%s\"\nGot: \"%s\"", expectedErr, err)
	}
}

func TestCorrectParamType(t *testing.T) {
	validator := NewFuncParamValidator("testFn").addParam("paramA", []string{"int"}, nil)
	_, err := validator.validate([]IRuntimeValue{NewIntegerRuntimeValue(42)})
	if err != nil {
		t.Errorf("\nExpected: nil\nGot: \"%s\"", err)
	}

	validator = NewFuncParamValidator("testFn").addParam("paramA", []string{"int"}, NewIntegerRuntimeValue(42))
	_, err = validator.validate([]IRuntimeValue{NewIntegerRuntimeValue(42)})
	if err != nil {
		t.Errorf("\nExpected: nil\nGot: \"%s\"", err)
	}

	validator = NewFuncParamValidator("testFn").addParam("paramA", []string{"int"}, NewIntegerRuntimeValue(42))
	_, err = validator.validate([]IRuntimeValue{})
	if err != nil {
		t.Errorf("\nExpected: nil\nGot: \"%s\"", err)
	}

	validator = NewFuncParamValidator("testFn").
		addParam("paramA", []string{"int"}, nil).
		addParam("paramB", []string{"int"}, NewIntegerRuntimeValue(42))
	_, err = validator.validate([]IRuntimeValue{NewIntegerRuntimeValue(0)})
	if err != nil {
		t.Errorf("\nExpected: nil\nGot: \"%s\"", err)
	}

	validator = NewFuncParamValidator("testFn").
		addParam("paramA", []string{"int"}, nil).
		addVariableLenParam("paramB", []string{"string"})
	_, err = validator.validate([]IRuntimeValue{NewIntegerRuntimeValue(42), NewStringRuntimeValue("abc"), NewStringRuntimeValue("abc")})
	if err != nil {
		t.Errorf("\nExpected: nil\nGot: \"%s\"", err)
	}

	validator = NewFuncParamValidator("testFn").
		addVariableLenParam("paramA", []string{"mixed"})
	_, err = validator.validate([]IRuntimeValue{NewIntegerRuntimeValue(42), NewBooleanRuntimeValue(true), NewStringRuntimeValue("abc")})
	if err != nil {
		t.Errorf("\nExpected: nil\nGot: \"%s\"", err)
	}
}
