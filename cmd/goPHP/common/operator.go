package common

import (
	"slices"
)

// Spec: https://phplang.org/spec/10-expressions.html#grammar-compound-assignment-operator
var compoundAssignmentOperators = []string{
	"**=", "*=", "/=", "%=", "+=", "-=", ".=", "<<=", ">>=", "&=", "^=", "|=",
}

func IsCompoundAssignmentOperator(op string) bool {
	// Spec: https://phplang.org/spec/10-expressions.html#grammar-compound-assignment-operator
	return slices.Contains(compoundAssignmentOperators, op)
}

// Spec: https://phplang.org/spec/10-expressions.html#grammar-equality-expression
var equalityOperators = []string{"==", "!=", "<>", "===", "!=="}

func IsEqualityOperator(op string) bool {
	// Spec: https://phplang.org/spec/10-expressions.html#grammar-equality-expression
	return slices.Contains(equalityOperators, op)
}

// Spec: https://phplang.org/spec/10-expressions.html#grammar-additive-expression
var additiveOperators = []string{"+", "-", "."}

func IsAdditiveOperator(op string) bool {
	// Spec: https://phplang.org/spec/10-expressions.html#grammar-additive-expression
	return slices.Contains(additiveOperators, op)
}

// Spec: https://phplang.org/spec/10-expressions.html#grammar-multiplicative-expression
var multiplicativeOperators = []string{"*", "/", "%"}

func IsMultiplicativeOperator(op string) bool {
	// Spec: https://phplang.org/spec/10-expressions.html#grammar-multiplicative-expression
	return slices.Contains(multiplicativeOperators, op)
}

// Spec: https://phplang.org/spec/10-expressions.html#grammar-unary-operator
var unaryOperators = []string{"+", "-", "~"}

func IsUnaryOperator(op string) bool {
	// Spec: https://phplang.org/spec/10-expressions.html#grammar-unary-operator
	return slices.Contains(unaryOperators, op)
}
