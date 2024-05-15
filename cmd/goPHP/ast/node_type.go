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
	NullLiteralExpr        NodeType = "NullLiteralExpression"
	VariableNameExpr       NodeType = "VariableNameExpression"
	SimpleVariableExpr     NodeType = "SimpleVariableExpression"
	FunctionCallExpr       NodeType = "FunctionCallExpression"
	SimpleAssignmentExpr   NodeType = "SimpleAssignmentExpression"
	ConstantAccessExpr     NodeType = "ConstantAccessExpression"
	CompoundAssignmentExpr NodeType = "CompoundAssignmentExpression"
	ConditionalExpr        NodeType = "ConditionalExpression"
	CoalesceExpr           NodeType = "CoalesceExpression"
	EqualityExpr           NodeType = "EqualityExpression"
	AdditiveExpr           NodeType = "AdditiveExpression"
	MultiplicativeExpr     NodeType = "MultiplicativeExpression"
	UnaryOpExpr            NodeType = "UnaryOpExpression"
	LogicalNotExpr         NodeType = "LogicalNotExpression"
	// Statements
	EchoStmt             NodeType = "EchoStatement"
	ConstDeclarationStmt NodeType = "ConstDeclarationStatement"
	ExpressionStmt       NodeType = "ExpressionStatement"
)
