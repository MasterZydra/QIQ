package common

import (
	"slices"
)

// Spec: https://phplang.org/spec/10-expressions.html#grammar-compound-assignment-operator
var compoundAssignmentOperator = []string{
	"**=", "*=", "/=", "%=", "+=", "-=", ".=", "<<=", ">>=", "&=", "^=", "|=",
}

func IsCompoundAssignmentOperator(op string) bool {
	// Spec: https://phplang.org/spec/10-expressions.html#grammar-compound-assignment-operator

	return slices.Contains(compoundAssignmentOperator, op)
}
