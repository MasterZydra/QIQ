package ast

import (
	"GoPHP/cmd/goPHP/position"
	"fmt"
)

// ------------------- MARK: Expression -------------------

type IExpression interface {
	IStatement
}

type Expression struct {
	id   int64
	kind NodeType
	pos  *position.Position
}

func NewExpr(id int64, kind NodeType, pos *position.Position) *Expression {
	return &Expression{id: id, kind: kind, pos: pos}
}

func (stmt *Expression) GetId() int64 {
	return stmt.id
}

func (expr *Expression) GetKind() NodeType {
	return expr.kind
}

func (expr *Expression) GetPosition() *position.Position {
	return expr.pos
}

func (expr *Expression) String() string {
	return fmt.Sprintf("{%s}", expr.GetKind())
}

func NewEmptyExpr() *Expression {
	return NewExpr(0, EmptyNode, nil)
}

// ------------------- MARK: TextExpression -------------------

type TextExpression struct {
	*Expression
	Value string
}

func NewTextExpr(id int64, value string) *TextExpression {
	return &TextExpression{Expression: NewExpr(id, TextNode, nil), Value: value}
}

func (expr *TextExpression) String() string {
	return fmt.Sprintf("{%s - value: \"%s\" }", expr.GetKind(), expr.Value)
}

// ------------------- MARK: VariableNameExpression -------------------

type VariableNameExpression struct {
	*Expression
	VariableName string
}

func NewVariableNameExpr(id int64, pos *position.Position, variableName string) *VariableNameExpression {
	return &VariableNameExpression{Expression: NewExpr(id, VariableNameExpr, pos), VariableName: variableName}
}

func (expr *VariableNameExpression) String() string {
	return fmt.Sprintf("{%s - variableName: \"%s\" }", expr.GetKind(), expr.VariableName)
}

// ------------------- MARK: SimpleVariableExpression -------------------

type SimpleVariableExpression struct {
	*Expression
	VariableName IExpression
}

func NewSimpleVariableExpr(id int64, variableName IExpression) *SimpleVariableExpression {
	return &SimpleVariableExpression{Expression: NewExpr(id, SimpleVariableExpr, variableName.GetPosition()), VariableName: variableName}
}

func (expr *SimpleVariableExpression) String() string {
	return fmt.Sprintf("{%s - variableName: \"%s\" }", expr.GetKind(), expr.VariableName)
}

// ------------------- MARK: SubscriptExpression -------------------

type SubscriptExpression struct {
	*Expression
	Variable IExpression
	Index    IExpression
}

func NewSubscriptExpr(id int64, variable IExpression, index IExpression) *SubscriptExpression {
	return &SubscriptExpression{Expression: NewExpr(id, SubscriptExpr, variable.GetPosition()), Variable: variable, Index: index}
}

func (expr *SubscriptExpression) String() string {
	return fmt.Sprintf("{%s - variable: %s, index: \"%s\" }", expr.GetKind(), expr.Variable, expr.Index)
}

// ------------------- MARK: FunctionCallExpression -------------------

type FunctionCallExpression struct {
	*Expression
	FunctionName string
	Arguments    []IExpression
}

func NewFunctionCallExpr(id int64, pos *position.Position, functionName string, arguments []IExpression) *FunctionCallExpression {
	return &FunctionCallExpression{Expression: NewExpr(id, FunctionCallExpr, pos), FunctionName: functionName, Arguments: arguments}
}

func (expr *FunctionCallExpression) String() string {
	return fmt.Sprintf("{%s - functionName: \"%s\" arguments: %s}", expr.GetKind(), expr.FunctionName, expr.Arguments)
}

// ------------------- MARK: EmptyIntrinsic -------------------

func NewExitIntrinsic(id int64, pos *position.Position, expression IExpression) *FunctionCallExpression {
	return &FunctionCallExpression{Expression: NewExpr(id, ExitIntrinsicExpr, pos),
		FunctionName: "exit", Arguments: []IExpression{expression},
	}
}

// ------------------- MARK: EmptyIntrinsic -------------------

func NewEmptyIntrinsic(id int64, pos *position.Position, expression IExpression) *FunctionCallExpression {
	return &FunctionCallExpression{Expression: NewExpr(id, EmptyIntrinsicExpr, pos),
		FunctionName: "empty", Arguments: []IExpression{expression},
	}
}

// ------------------- MARK: IssetIntrinsic -------------------

func NewIssetIntrinsic(id int64, pos *position.Position, arguments []IExpression) *FunctionCallExpression {
	return &FunctionCallExpression{Expression: NewExpr(id, IssetIntrinsicExpr, pos),
		FunctionName: "isset", Arguments: arguments,
	}
}

// ------------------- MARK: UnsetIntrinsic -------------------

func NewUnsetIntrinsic(id int64, pos *position.Position, arguments []IExpression) *FunctionCallExpression {
	return &FunctionCallExpression{Expression: NewExpr(id, UnsetIntrinsicExpr, pos),
		FunctionName: "unset", Arguments: arguments,
	}
}

// ------------------- MARK: ConstantAccessExpression -------------------

type ConstantAccessExpression struct {
	*Expression
	ConstantName string
}

func NewConstantAccessExpr(id int64, pos *position.Position, constantName string) *ConstantAccessExpression {
	return &ConstantAccessExpression{Expression: NewExpr(id, ConstantAccessExpr, pos), ConstantName: constantName}
}

func (expr *ConstantAccessExpression) String() string {
	return fmt.Sprintf("{%s - constantName: %s}", expr.GetKind(), expr.ConstantName)
}

// ------------------- MARK: ArrayLiteralExpression -------------------

type ArrayLiteralExpression struct {
	*Expression
	Keys     []IExpression
	Elements map[IExpression]IExpression
}

func NewArrayLiteralExpr(id int64, pos *position.Position) *ArrayLiteralExpression {
	return &ArrayLiteralExpression{
		Expression: NewExpr(id, ArrayLiteralExpr, pos),
		Keys:       []IExpression{},
		Elements:   map[IExpression]IExpression{},
	}
}

func (expr *ArrayLiteralExpression) AddElement(key IExpression, value IExpression) {
	expr.Keys = append(expr.Keys, key)
	expr.Elements[key] = value
}

func (expr *ArrayLiteralExpression) String() string {
	return fmt.Sprintf("{%s - elements: %s }", expr.GetKind(), expr.Elements)
}

// ------------------- MARK: BooleanLiteralExpression -------------------

func NewBooleanLiteralExpr(id int64, pos *position.Position, value bool) *ConstantAccessExpression {
	if value {
		return NewConstantAccessExpr(id, pos, "TRUE")
	}
	return NewConstantAccessExpr(id, pos, "FALSE")
}

// ------------------- MARK: IntegerLiteralExpression -------------------

type IntegerLiteralExpression struct {
	*Expression
	Value int64
}

func NewIntegerLiteralExpr(id int64, pos *position.Position, value int64) *IntegerLiteralExpression {
	return &IntegerLiteralExpression{Expression: NewExpr(id, IntegerLiteralExpr, pos), Value: value}
}

func (expr *IntegerLiteralExpression) String() string {
	return fmt.Sprintf("{%s - value: %d }", expr.GetKind(), expr.Value)
}

// ------------------- MARK: FloatingLiteralExpression -------------------

type FloatingLiteralExpression struct {
	*Expression
	Value float64
}

func NewFloatingLiteralExpr(id int64, pos *position.Position, value float64) *FloatingLiteralExpression {
	return &FloatingLiteralExpression{Expression: NewExpr(id, FloatingLiteralExpr, pos), Value: value}
}

func (expr *FloatingLiteralExpression) String() string {
	return fmt.Sprintf("{%s - value: %f }", expr.GetKind(), expr.Value)
}

// ------------------- MARK: StringLiteralExpression -------------------

type StringType string

const (
	SingleQuotedString StringType = "SingleQuotedString"
	DoubleQuotedString StringType = "DoubleQuotedString"
)

type StringLiteralExpression struct {
	*Expression
	StringType StringType
	Value      string
}

func NewStringLiteralExpr(id int64, pos *position.Position, value string, stringType StringType) *StringLiteralExpression {
	return &StringLiteralExpression{Expression: NewExpr(id, StringLiteralExpr, pos), Value: value, StringType: stringType}
}

func (expr *StringLiteralExpression) String() string {
	return fmt.Sprintf("{%s - type: \"%s\" value: \"%s\" }", expr.GetKind(), expr.StringType, expr.Value)
}

// ------------------- MARK: NullLiteralExpression -------------------

func NewNullLiteralExpr(id int64, pos *position.Position) *ConstantAccessExpression {
	return NewConstantAccessExpr(id, pos, "NULL")
}

// ------------------- MARK: SimpleAssignmentExpression -------------------

type SimpleAssignmentExpression struct {
	*Expression
	Variable IExpression
	Value    IExpression
}

func NewSimpleAssignmentExpr(id int64, variable IExpression, value IExpression) *SimpleAssignmentExpression {
	return &SimpleAssignmentExpression{Expression: NewExpr(id, SimpleAssignmentExpr, variable.GetPosition()), Variable: variable, Value: value}
}

func (expr *SimpleAssignmentExpression) String() string {
	return fmt.Sprintf("{%s - variable: %s, value: %s }", expr.GetKind(), expr.Variable, expr.Value)
}

// ------------------- MARK: CompoundAssignmentExpression -------------------

type CompoundAssignmentExpression struct {
	*Expression
	Variable IExpression
	Operator string
	Value    IExpression
}

func NewCompoundAssignmentExpr(id int64, variable IExpression, operator string, value IExpression) *CompoundAssignmentExpression {
	return &CompoundAssignmentExpression{
		Expression: NewExpr(id, CompoundAssignmentExpr, variable.GetPosition()), Variable: variable, Operator: operator, Value: value,
	}
}

func (expr *CompoundAssignmentExpression) String() string {
	return fmt.Sprintf(
		"{%s - variable: %s, operator: \"%s\", value: %s }",
		expr.GetKind(), expr.Variable, expr.Operator, expr.Value,
	)
}

// ------------------- MARK: ConditionalExpression -------------------

type ConditionalExpression struct {
	*Expression
	Cond     IExpression
	IfExpr   IExpression
	ElseExpr IExpression
}

func NewConditionalExpr(id int64, cond IExpression, ifExpr IExpression, elseExpr IExpression) *ConditionalExpression {
	return &ConditionalExpression{Expression: NewExpr(id, ConditionalExpr, cond.GetPosition()), Cond: cond, IfExpr: ifExpr, ElseExpr: elseExpr}
}

func (expr *ConditionalExpression) String() string {
	return fmt.Sprintf("{%s - condition: %s, ifExpr: %s, elseExpr: %s }", expr.GetKind(), expr.Cond, expr.IfExpr, expr.ElseExpr)
}

// ------------------- MARK: CoalesceExpression -------------------

type CoalesceExpression struct {
	*Expression
	Cond     IExpression
	ElseExpr IExpression
}

func NewCoalesceExpr(id int64, cond IExpression, elseExpr IExpression) *CoalesceExpression {
	return &CoalesceExpression{Expression: NewExpr(id, CoalesceExpr, cond.GetPosition()), Cond: cond, ElseExpr: elseExpr}
}

func (expr *CoalesceExpression) String() string {
	return fmt.Sprintf("{%s - condition: %s, elseExpr: %s }", expr.GetKind(), expr.Cond, expr.ElseExpr)
}

// ------------------- MARK: BinaryOpExpression -------------------

type BinaryOpExpression struct {
	*Expression
	Lhs      IExpression
	Operator string
	Rhs      IExpression
}

func NewRelationalExpr(id int64, lhs IExpression, operator string, rhs IExpression) *BinaryOpExpression {
	return &BinaryOpExpression{Expression: NewExpr(id, RelationalExpr, lhs.GetPosition()), Lhs: lhs, Operator: operator, Rhs: rhs}
}

func NewEqualityExpr(id int64, lhs IExpression, operator string, rhs IExpression) *BinaryOpExpression {
	return &BinaryOpExpression{Expression: NewExpr(id, EqualityExpr, lhs.GetPosition()), Lhs: lhs, Operator: operator, Rhs: rhs}
}

func NewBinaryOpExpr(id int64, lhs IExpression, operator string, rhs IExpression) *BinaryOpExpression {
	return &BinaryOpExpression{Expression: NewExpr(id, BinaryOpExpr, lhs.GetPosition()), Lhs: lhs, Operator: operator, Rhs: rhs}
}

func (expr *BinaryOpExpression) String() string {
	return fmt.Sprintf("{%s - lhs: %s, operator: \"%s\" rhs: %s }", expr.GetKind(), expr.Lhs, expr.Operator, expr.Rhs)
}

// ------------------- MARK: UnaryOpExpression -------------------

type UnaryOpExpression struct {
	*Expression
	Operator string
	Expr     IExpression
}

func NewPrefixIncExpr(id int64, pos *position.Position, expression IExpression, operator string) *UnaryOpExpression {
	return &UnaryOpExpression{Expression: NewExpr(id, PrefixIncExpr, pos), Operator: operator, Expr: expression}
}

func NewPostfixIncExpr(id int64, pos *position.Position, expression IExpression, operator string) *UnaryOpExpression {
	return &UnaryOpExpression{Expression: NewExpr(id, PostfixIncExpr, pos), Operator: operator, Expr: expression}
}

func NewLogicalNotExpr(id int64, pos *position.Position, expression IExpression) *UnaryOpExpression {
	return &UnaryOpExpression{Expression: NewExpr(id, LogicalNotExpr, pos), Operator: "!", Expr: expression}
}

func NewCastExpr(id int64, pos *position.Position, castType string, expression IExpression) *UnaryOpExpression {
	return &UnaryOpExpression{Expression: NewExpr(id, CastExpr, pos), Operator: castType, Expr: expression}
}

func NewUnaryOpExpr(id int64, pos *position.Position, operator string, expression IExpression) *UnaryOpExpression {
	return &UnaryOpExpression{Expression: NewExpr(id, UnaryOpExpr, pos), Operator: operator, Expr: expression}
}

func (expr *UnaryOpExpression) String() string {
	return fmt.Sprintf("{%s - operator: \"%s\" expression: %s }", expr.GetKind(), expr.Operator, expr.Expr)
}

// ------------------- MARK: ExprExpression -------------------

type ExprExpression struct {
	*Expression
	Expr IExpression
}

func NewIncludeExpr(id int64, pos *position.Position, expression IExpression) *ExprExpression {
	return &ExprExpression{Expression: NewExpr(id, IncludeExpr, pos), Expr: expression}
}

func NewIncludeOnceExpr(id int64, pos *position.Position, expression IExpression) *ExprExpression {
	return &ExprExpression{Expression: NewExpr(id, IncludeOnceExpr, pos), Expr: expression}
}

func NewRequireExpr(id int64, pos *position.Position, expression IExpression) *ExprExpression {
	return &ExprExpression{Expression: NewExpr(id, RequireExpr, pos), Expr: expression}
}

func NewRequireOnceExpr(id int64, pos *position.Position, expression IExpression) *ExprExpression {
	return &ExprExpression{Expression: NewExpr(id, RequireOnceExpr, pos), Expr: expression}
}

func (expr *ExprExpression) String() string {
	return fmt.Sprintf("{%s - expression: %s }", expr.GetKind(), expr.Expr)
}
