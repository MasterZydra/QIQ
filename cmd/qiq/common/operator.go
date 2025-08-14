package common

import (
	"slices"
)

// Spec: https://phplang.org/spec/10-expressions.html#grammar-relational-expression
var relationalExpressionOps = []string{"<", ">", "<=", ">=", "<=>"}

func IsRelationalExpressionOp(op string) bool {
	// Spec: https://phplang.org/spec/10-expressions.html#grammar-relational-expression
	return slices.Contains(relationalExpressionOps, op)
}

// Spec: https://phplang.org/spec/10-expressions.html#grammar-compound-assignment-operator
var compoundAssignmentOps = []string{
	"**=", "*=", "/=", "%=", "+=", "-=", ".=", "<<=", ">>=", "&=", "^=", "|=",
}

func IsCompoundAssignmentOp(op string) bool {
	// Spec: https://phplang.org/spec/10-expressions.html#grammar-compound-assignment-operator
	return slices.Contains(compoundAssignmentOps, op)
}

// Spec: https://phplang.org/spec/10-expressions.html#grammar-equality-expression
var equalityOps = []string{"==", "!=", "<>", "===", "!=="}

func IsEqualityOp(op string) bool {
	// Spec: https://phplang.org/spec/10-expressions.html#grammar-equality-expression
	return slices.Contains(equalityOps, op)
}

// Spec: https://phplang.org/spec/10-expressions.html#grammar-additive-expression
var additiveOps = []string{"+", "-", "."}

func IsAdditiveOp(op string) bool {
	// Spec: https://phplang.org/spec/10-expressions.html#grammar-additive-expression
	return slices.Contains(additiveOps, op)
}

// Spec: https://phplang.org/spec/10-expressions.html#grammar-multiplicative-expression
var multiplicativeOps = []string{"*", "/", "%"}

func IsMultiplicativeOp(op string) bool {
	// Spec: https://phplang.org/spec/10-expressions.html#grammar-multiplicative-expression
	return slices.Contains(multiplicativeOps, op)
}

// Spec: https://phplang.org/spec/10-expressions.html#grammar-unary-operator
var unaryOps = []string{"+", "-", "~"}

func IsUnaryOp(op string) bool {
	// Spec: https://phplang.org/spec/10-expressions.html#grammar-unary-operator
	return slices.Contains(unaryOps, op)
}
