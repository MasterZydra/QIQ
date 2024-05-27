package ast

import (
	"slices"
)

// Spec: https://phplang.org/spec/10-expressions.html#grammar-variable
var variableExpressions = []NodeType{
	SimpleVariableExpr, SubscriptExpr, FunctionCallExpr,
}

func IsVariableExpression(expr IExpression) bool {
	// Spec: https://phplang.org/spec/10-expressions.html#grammar-variable

	// variable:
	//    callable-variable
	//    scoped-property-access-expression
	//    member-access-expression

	// callable-variable:
	//    simple-variable
	//    subscript-expression
	//    member-call-expression
	//    scoped-call-expression
	//    function-call-expression

	if expr == nil {
		return false
	}

	return slices.Contains(variableExpressions, expr.GetKind())
}

var nodeId int64 = 0

func getNextNodeId() int64 {
	result := nodeId
	nodeId += 1
	return result
}
