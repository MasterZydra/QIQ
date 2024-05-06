package ast

type NodeType string

const (
	EmptyNode   NodeType = "Empty"
	ProgramNode NodeType = "Program"
	TextNode    NodeType = "Text"
	// Expressions
	VariableNameExpr    NodeType = "VariableNameExpression"
	SimpleVariableExpr  NodeType = "SimpleVariableExpression"
	IntegerLiteralExpr  NodeType = "IntegerLiteralExpression"
	FloatingLiteralExpr NodeType = "FloatingLiteralExpression"
	StringLiteralExpr   NodeType = "StringLiteralExpression"
	// Statements
	EchoStmt       NodeType = "EchoStatement"
	ExpressionStmt NodeType = "ExpressionStatement"
)
