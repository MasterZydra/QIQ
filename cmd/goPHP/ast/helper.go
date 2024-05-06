package ast

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

	return expr.GetKind() == SimpleVariableExpr
}
