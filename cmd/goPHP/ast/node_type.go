package ast

type NodeType string

const (
	EmptyNode   NodeType = "Empty"
	ProgramNode NodeType = "Program"
	TextNode    NodeType = "Text"
	// Expressions
	ArrayLiteralExpr       NodeType = "ArrayLiteralExpression"
	IntegerLiteralExpr     NodeType = "IntegerLiteralExpression"
	FloatingLiteralExpr    NodeType = "FloatingLiteralExpression"
	StringLiteralExpr      NodeType = "StringLiteralExpression"
	VariableNameExpr       NodeType = "VariableNameExpression"
	SimpleVariableExpr     NodeType = "SimpleVariableExpression"
	SubscriptExpr          NodeType = "SubscriptExpression"
	FunctionCallExpr       NodeType = "FunctionCallExpression"
	EmptyIntrinsicExpr     NodeType = "EmptyIntrinsicExpression"
	IssetIntrinsicExpr     NodeType = "IssetIntrinsicExpression"
	UnsetIntrinsicExpr     NodeType = "UnsetIntrinsicExpression"
	SimpleAssignmentExpr   NodeType = "SimpleAssignmentExpression"
	ConstantAccessExpr     NodeType = "ConstantAccessExpression"
	CompoundAssignmentExpr NodeType = "CompoundAssignmentExpression"
	ConditionalExpr        NodeType = "ConditionalExpression"
	CoalesceExpr           NodeType = "CoalesceExpression"
	BinaryOpExpr           NodeType = "BinaryOpExpression"
	EqualityExpr           NodeType = "EqualityExpression"
	RelationalExpr         NodeType = "RelationalExpression"
	ShiftExpr              NodeType = "ShiftExpression"
	UnaryOpExpr            NodeType = "UnaryOpExpression"
	PostfixIncExpr         NodeType = "PostfixIncExpression"
	PrefixIncExpr          NodeType = "PrefixIncExpression"
	LogicalNotExpr         NodeType = "LogicalNotExpression"
	// Statements
	IfStmt               NodeType = "IfStatement"
	CompoundStmt         NodeType = "CompoundStatement"
	EchoStmt             NodeType = "EchoStatement"
	ConstDeclarationStmt NodeType = "ConstDeclarationStatement"
	ExpressionStmt       NodeType = "ExpressionStatement"
)
