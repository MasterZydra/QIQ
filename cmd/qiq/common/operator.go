package common

import (
	"slices"
)

// Spec: https://phplang.org/spec/10-expressions.html#grammar-relational-expression
var relationalExpressionOps = []string{"<", ">", "<=", ">=", "<=>"}

// Spec: https://phplang.org/spec/10-expressions.html#grammar-relational-expression
func IsRelationalExpressionOp(op string) bool { return slices.Contains(relationalExpressionOps, op) }

// Spec: https://phplang.org/spec/10-expressions.html#grammar-compound-assignment-operator
var compoundAssignmentOps = []string{
	"**=", "*=", "/=", "%=", "+=", "-=", ".=", "<<=", ">>=", "&=", "^=", "|=",
}

// Spec: https://phplang.org/spec/10-expressions.html#grammar-compound-assignment-operator
func IsCompoundAssignmentOp(op string) bool { return slices.Contains(compoundAssignmentOps, op) }

// Spec: https://phplang.org/spec/10-expressions.html#grammar-equality-expression
var equalityOps = []string{"==", "!=", "<>", "===", "!=="}

// Spec: https://phplang.org/spec/10-expressions.html#grammar-equality-expression
func IsEqualityOp(op string) bool { return slices.Contains(equalityOps, op) }

// Spec: https://phplang.org/spec/10-expressions.html#grammar-additive-expression
var additiveOps = []string{"+", "-", "."}

// Spec: https://phplang.org/spec/10-expressions.html#grammar-additive-expression
func IsAdditiveOp(op string) bool { return slices.Contains(additiveOps, op) }

// Spec: https://phplang.org/spec/10-expressions.html#grammar-multiplicative-expression
var multiplicativeOps = []string{"*", "/", "%"}

// Spec: https://phplang.org/spec/10-expressions.html#grammar-multiplicative-expression
func IsMultiplicativeOp(op string) bool { return slices.Contains(multiplicativeOps, op) }

// Spec: https://phplang.org/spec/10-expressions.html#grammar-unary-operator
var unaryOps = []string{"+", "-", "~"}

// Spec: https://phplang.org/spec/10-expressions.html#grammar-unary-operator
func IsUnaryOp(op string) bool { return slices.Contains(unaryOps, op) }
