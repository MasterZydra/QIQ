package ast

type NodeType string

const (
	EmptyNode   NodeType = "Empty"
	ProgramNode NodeType = "Program"
	TextNode    NodeType = "Text"
	// Expressions
	IntegerLiteralExpr   NodeType = "IntegerLiteralExpression"
	FloatingLiteralExpr  NodeType = "FloatingLiteralExpression"
	StringLiteralExpr    NodeType = "StringLiteralExpression"
	VariableNameExpr     NodeType = "VariableNameExpression"
	SimpleVariableExpr   NodeType = "SimpleVariableExpression"
	SimpleAssignmentExpr NodeType = "SimpleAssignmentExpression"
	// Statements
	EchoStmt       NodeType = "EchoStatement"
	ExpressionStmt NodeType = "ExpressionStatement"
)
