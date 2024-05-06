package ast

type NodeType string

const (
	EmptyNode   NodeType = "Empty"
	ProgramNode NodeType = "Program"
	TextNode    NodeType = "Text"
	// Expressions
	BooleanLiteralExpr     NodeType = "BooleanLiteralExpression"
	IntegerLiteralExpr     NodeType = "IntegerLiteralExpression"
	FloatingLiteralExpr    NodeType = "FloatingLiteralExpression"
	StringLiteralExpr      NodeType = "StringLiteralExpression"
	VariableNameExpr       NodeType = "VariableNameExpression"
	SimpleVariableExpr     NodeType = "SimpleVariableExpression"
	SimpleAssignmentExpr   NodeType = "SimpleAssignmentExpression"
	CompoundAssignmentExpr NodeType = "CompoundAssignmentExpression"
	ConditionalExpr        NodeType = "ConditionalExpression"
	// Statements
	EchoStmt       NodeType = "EchoStatement"
	ExpressionStmt NodeType = "ExpressionStatement"
)
