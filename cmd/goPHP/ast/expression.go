package ast

import (
	"GoPHP/cmd/goPHP/position"
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

func NewEmptyExpr() *Expression {
	return NewExpr(0, EmptyNode, nil)
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
	if expr.pos == nil {
		return &position.Position{}
	}
	return expr.pos
}

func (stmt *Expression) Process(visitor Visitor, context any) (any, error) {
	return visitor.ProcessExpr(stmt, context)
}

// ------------------- MARK: ParenthesizedExpression -------------------

type ParenthesizedExpression struct {
	*Expression
	Expr IExpression
}

func NewParenthesizedExpr(id int64, pos *position.Position, expr IExpression) *ParenthesizedExpression {
	return &ParenthesizedExpression{Expression: NewExpr(id, ParenthesizedExpr, pos), Expr: expr}
}

func (stmt *ParenthesizedExpression) Process(visitor Visitor, context any) (any, error) {
	return visitor.ProcessParenthesizedExpr(stmt, context)
}

// ------------------- MARK: PrintExpression -------------------

type PrintExpression struct {
	*ParenthesizedExpression
}

func NewPrintExpr(id int64, pos *position.Position, expr IExpression) *PrintExpression {
	return &PrintExpression{&ParenthesizedExpression{Expression: NewExpr(id, PrintExpr, pos), Expr: expr}}
}

func (stmt *PrintExpression) Process(visitor Visitor, context any) (any, error) {
	return visitor.ProcessPrintExpr(stmt, context)
}

// ------------------- MARK: TextExpression -------------------

type TextExpression struct {
	*Expression
	Value string
}

func NewTextExpr(id int64, value string) *TextExpression {
	return &TextExpression{Expression: NewExpr(id, TextNode, nil), Value: value}
}

func (stmt *TextExpression) Process(visitor Visitor, context any) (any, error) {
	return visitor.ProcessTextExpr(stmt, context)
}

// ------------------- MARK: VariableNameExpression -------------------

type VariableNameExpression struct {
	*Expression
	VariableName string
}

func NewVariableNameExpr(id int64, pos *position.Position, variableName string) *VariableNameExpression {
	return &VariableNameExpression{Expression: NewExpr(id, VariableNameExpr, pos), VariableName: variableName}
}

func (stmt *VariableNameExpression) Process(visitor Visitor, context any) (any, error) {
	return visitor.ProcessVariableNameExpr(stmt, context)
}

// ------------------- MARK: SimpleVariableExpression -------------------

type SimpleVariableExpression struct {
	*Expression
	VariableName IExpression
}

func NewSimpleVariableExpr(id int64, variableName IExpression) *SimpleVariableExpression {
	return &SimpleVariableExpression{Expression: NewExpr(id, SimpleVariableExpr, variableName.GetPosition()), VariableName: variableName}
}

func (stmt *SimpleVariableExpression) Process(visitor Visitor, context any) (any, error) {
	return visitor.ProcessSimpleVariableExpr(stmt, context)
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

func (stmt *SubscriptExpression) Process(visitor Visitor, context any) (any, error) {
	return visitor.ProcessSubscriptExpr(stmt, context)
}

// ------------------- MARK: FunctionCallExpression -------------------

type FunctionCallExpression struct {
	*Expression
	FunctionName IExpression
	Arguments    []IExpression
}

func NewFunctionCallExpr(id int64, pos *position.Position, functionName IExpression, arguments []IExpression) *FunctionCallExpression {
	return &FunctionCallExpression{Expression: NewExpr(id, FunctionCallExpr, pos), FunctionName: functionName, Arguments: arguments}
}

func (stmt *FunctionCallExpression) Process(visitor Visitor, context any) (any, error) {
	return visitor.ProcessFunctionCallExpr(stmt, context)
}

// ------------------- MARK: ExitIntrinsic -------------------

type ExitIntrinsicExpression struct {
	*FunctionCallExpression
}

func NewExitIntrinsic(id int64, pos *position.Position, expression IExpression) *ExitIntrinsicExpression {
	return &ExitIntrinsicExpression{FunctionCallExpression: &FunctionCallExpression{
		Expression:   NewExpr(id, ExitIntrinsicExpr, pos),
		FunctionName: NewStringLiteralExpr(id, pos, "exit", SingleQuotedString),
		Arguments:    []IExpression{expression},
	}}
}

func (stmt *ExitIntrinsicExpression) Process(visitor Visitor, context any) (any, error) {
	return visitor.ProcessExitIntrinsicExpr(stmt, context)
}

// ------------------- MARK: EmptyIntrinsic -------------------

type EmptyIntrinsicExpression struct {
	*FunctionCallExpression
}

func NewEmptyIntrinsic(id int64, pos *position.Position, expression IExpression) *EmptyIntrinsicExpression {
	return &EmptyIntrinsicExpression{&FunctionCallExpression{
		Expression:   NewExpr(id, EmptyIntrinsicExpr, pos),
		FunctionName: NewStringLiteralExpr(id, pos, "empty", SingleQuotedString),
		Arguments:    []IExpression{expression},
	}}
}

func (stmt *EmptyIntrinsicExpression) Process(visitor Visitor, context any) (any, error) {
	return visitor.ProcessEmptyIntrinsicExpr(stmt, context)
}

// ------------------- MARK: EvalIntrinsic -------------------

type EvalIntrinsicExpression struct {
	*FunctionCallExpression
}

func NewEvalIntrinsic(id int64, pos *position.Position, expression IExpression) *EvalIntrinsicExpression {
	return &EvalIntrinsicExpression{&FunctionCallExpression{
		Expression:   NewExpr(id, EvalIntrinsicExpr, pos),
		FunctionName: NewStringLiteralExpr(id, pos, "eval", SingleQuotedString),
		Arguments:    []IExpression{expression},
	}}
}

func (stmt *EvalIntrinsicExpression) Process(visitor Visitor, context any) (any, error) {
	return visitor.ProcessEvalIntrinsicExpr(stmt, context)
}

// ------------------- MARK: IssetIntrinsic -------------------

type IssetIntrinsicExpression struct {
	*FunctionCallExpression
}

func NewIssetIntrinsic(id int64, pos *position.Position, arguments []IExpression) *IssetIntrinsicExpression {
	return &IssetIntrinsicExpression{&FunctionCallExpression{
		Expression:   NewExpr(id, IssetIntrinsicExpr, pos),
		FunctionName: NewStringLiteralExpr(id, pos, "isset", SingleQuotedString),
		Arguments:    arguments,
	}}
}

func (stmt *IssetIntrinsicExpression) Process(visitor Visitor, context any) (any, error) {
	return visitor.ProcessIssetIntrinsicExpr(stmt, context)
}

// ------------------- MARK: UnsetIntrinsic -------------------

type UnsetIntrinsicExpression struct {
	*FunctionCallExpression
}

func NewUnsetIntrinsic(id int64, pos *position.Position, arguments []IExpression) *UnsetIntrinsicExpression {
	return &UnsetIntrinsicExpression{&FunctionCallExpression{
		Expression:   NewExpr(id, UnsetIntrinsicExpr, pos),
		FunctionName: NewStringLiteralExpr(id, pos, "unset", SingleQuotedString),
		Arguments:    arguments,
	}}
}

func (stmt *UnsetIntrinsicExpression) Process(visitor Visitor, context any) (any, error) {
	return visitor.ProcessUnsetIntrinsicExpr(stmt, context)
}

// ------------------- MARK: ConstantAccessExpression -------------------

type ConstantAccessExpression struct {
	*Expression
	ConstantName string
}

func NewConstantAccessExpr(id int64, pos *position.Position, constantName string) *ConstantAccessExpression {
	return &ConstantAccessExpression{Expression: NewExpr(id, ConstantAccessExpr, pos), ConstantName: constantName}
}

func (stmt *ConstantAccessExpression) Process(visitor Visitor, context any) (any, error) {
	return visitor.ProcessConstantAccessExpr(stmt, context)
}

// ------------------- MARK: ArrayLiteralExpression -------------------

type ArrayLiteralExpression struct {
	*Expression
	Keys           []IExpression
	Elements       map[IExpression]IExpression
	arrayNextKeyId int64
}

func NewArrayLiteralExpr(id int64, pos *position.Position) *ArrayLiteralExpression {
	return &ArrayLiteralExpression{
		Expression: NewExpr(id, ArrayLiteralExpr, pos),
		Keys:       []IExpression{},
		Elements:   map[IExpression]IExpression{},
	}
}

func (expr *ArrayLiteralExpression) AddElement(key IExpression, value IExpression) {
	if key == nil {
		key = NewArrayNextKey(expr.arrayNextKeyId)
		expr.arrayNextKeyId++
	}
	expr.Keys = append(expr.Keys, key)
	expr.Elements[key] = value
}

func (expr *ArrayLiteralExpression) Process(visitor Visitor, context any) (any, error) {
	return visitor.ProcessArrayLiteralExpr(expr, context)
}

type ArrayNextKeyExpression struct {
	id int64
}

func NewArrayNextKey(id int64) *ArrayNextKeyExpression {
	return &ArrayNextKeyExpression{id: id}
}

// GetId implements IExpression.
func (expr *ArrayNextKeyExpression) GetId() int64 {
	return expr.id
}

// GetKind implements IExpression.
func (expr *ArrayNextKeyExpression) GetKind() NodeType {
	return ArrayNextKeyExpr
}

// GetPosition implements IExpression.
func (expr *ArrayNextKeyExpression) GetPosition() *position.Position {
	return &position.Position{}
}

// Process implements IExpression.
func (expr *ArrayNextKeyExpression) Process(visitor Visitor, context any) (any, error) {
	return visitor.ProcessArrayNextKeyExpr(expr, context)
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

func (stmt *IntegerLiteralExpression) Process(visitor Visitor, context any) (any, error) {
	return visitor.ProcessIntegerLiteralExpr(stmt, context)
}

// ------------------- MARK: FloatingLiteralExpression -------------------

type FloatingLiteralExpression struct {
	*Expression
	Value float64
}

func NewFloatingLiteralExpr(id int64, pos *position.Position, value float64) *FloatingLiteralExpression {
	return &FloatingLiteralExpression{Expression: NewExpr(id, FloatingLiteralExpr, pos), Value: value}
}

func (stmt *FloatingLiteralExpression) Process(visitor Visitor, context any) (any, error) {
	return visitor.ProcessFloatingLiteralExpr(stmt, context)
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

func (stmt *StringLiteralExpression) Process(visitor Visitor, context any) (any, error) {
	return visitor.ProcessStringLiteralExpr(stmt, context)
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

func (stmt *SimpleAssignmentExpression) Process(visitor Visitor, context any) (any, error) {
	return visitor.ProcessSimpleAssignmentExpr(stmt, context)
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

func (stmt *CompoundAssignmentExpression) Process(visitor Visitor, context any) (any, error) {
	return visitor.ProcessCompoundAssignmentExpr(stmt, context)
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

func (stmt *ConditionalExpression) Process(visitor Visitor, context any) (any, error) {
	return visitor.ProcessConditionalExpr(stmt, context)
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

func (stmt *CoalesceExpression) Process(visitor Visitor, context any) (any, error) {
	return visitor.ProcessCoalesceExpr(stmt, context)
}

// ------------------- MARK: BinaryOpExpression -------------------

type BinaryOpExpression struct {
	*Expression
	Lhs      IExpression
	Operator string
	Rhs      IExpression
}

func NewBinaryOpExpr(id int64, lhs IExpression, operator string, rhs IExpression) *BinaryOpExpression {
	return &BinaryOpExpression{Expression: NewExpr(id, BinaryOpExpr, lhs.GetPosition()), Lhs: lhs, Operator: operator, Rhs: rhs}
}

func (stmt *BinaryOpExpression) Process(visitor Visitor, context any) (any, error) {
	return visitor.ProcessBinaryOpExpr(stmt, context)
}

// ------------------- MARK: LogicalExpression -------------------

type LogicalExpression struct {
	*BinaryOpExpression
}

func NewLogicalExpr(id int64, lhs IExpression, operator string, rhs IExpression) *LogicalExpression {
	return &LogicalExpression{&BinaryOpExpression{Expression: NewExpr(id, BinaryOpExpr, lhs.GetPosition()), Lhs: lhs, Operator: operator, Rhs: rhs}}
}

func (stmt *LogicalExpression) Process(visitor Visitor, context any) (any, error) {
	return visitor.ProcessLogicalExpr(stmt, context)
}

// ------------------- MARK: RelationalExpression -------------------

type RelationalExpression struct {
	*BinaryOpExpression
}

func NewRelationalExpr(id int64, lhs IExpression, operator string, rhs IExpression) *RelationalExpression {
	return &RelationalExpression{&BinaryOpExpression{Expression: NewExpr(id, RelationalExpr, lhs.GetPosition()), Lhs: lhs, Operator: operator, Rhs: rhs}}
}

func (stmt *RelationalExpression) Process(visitor Visitor, context any) (any, error) {
	return visitor.ProcessRelationalExpr(stmt, context)
}

// ------------------- MARK: EqualityExpression -------------------

type EqualityExpression struct {
	*BinaryOpExpression
}

func NewEqualityExpr(id int64, lhs IExpression, operator string, rhs IExpression) *EqualityExpression {
	return &EqualityExpression{&BinaryOpExpression{Expression: NewExpr(id, EqualityExpr, lhs.GetPosition()), Lhs: lhs, Operator: operator, Rhs: rhs}}
}

func (stmt *EqualityExpression) Process(visitor Visitor, context any) (any, error) {
	return visitor.ProcessEqualityExpr(stmt, context)
}

// ------------------- MARK: UnaryOpExpression -------------------

type UnaryOpExpression struct {
	*Expression
	Operator string
	Expr     IExpression
}

func NewUnaryOpExpr(id int64, pos *position.Position, operator string, expression IExpression) *UnaryOpExpression {
	return &UnaryOpExpression{Expression: NewExpr(id, UnaryOpExpr, pos), Operator: operator, Expr: expression}
}

func (stmt *UnaryOpExpression) Process(visitor Visitor, context any) (any, error) {
	return visitor.ProcessUnaryExpr(stmt, context)
}

// ------------------- MARK: PrefixIncExpression -------------------

type PrefixIncExpression struct {
	*UnaryOpExpression
}

func NewPrefixIncExpr(id int64, pos *position.Position, expression IExpression, operator string) *PrefixIncExpression {
	return &PrefixIncExpression{&UnaryOpExpression{Expression: NewExpr(id, PrefixIncExpr, pos), Operator: operator, Expr: expression}}
}

func (stmt *PrefixIncExpression) Process(visitor Visitor, context any) (any, error) {
	return visitor.ProcessPrefixIncExpr(stmt, context)
}

// ------------------- MARK: PostfixIncExpression -------------------

type PostfixIncExpression struct {
	*UnaryOpExpression
}

func NewPostfixIncExpr(id int64, pos *position.Position, expression IExpression, operator string) *PostfixIncExpression {
	return &PostfixIncExpression{&UnaryOpExpression{Expression: NewExpr(id, PostfixIncExpr, pos), Operator: operator, Expr: expression}}
}

func (stmt *PostfixIncExpression) Process(visitor Visitor, context any) (any, error) {
	return visitor.ProcessPostfixIncExpr(stmt, context)
}

// ------------------- MARK: LogicalNotExpression -------------------

type LogicalNotExpression struct {
	*UnaryOpExpression
}

func NewLogicalNotExpr(id int64, pos *position.Position, expression IExpression) *LogicalNotExpression {
	return &LogicalNotExpression{&UnaryOpExpression{Expression: NewExpr(id, LogicalNotExpr, pos), Operator: "!", Expr: expression}}
}

func (stmt *LogicalNotExpression) Process(visitor Visitor, context any) (any, error) {
	return visitor.ProcessLogicalNotExpr(stmt, context)
}

// ------------------- MARK: CastExpression -------------------

type CastExpression struct {
	*UnaryOpExpression
}

func NewCastExpr(id int64, pos *position.Position, castType string, expression IExpression) *CastExpression {
	return &CastExpression{&UnaryOpExpression{Expression: NewExpr(id, CastExpr, pos), Operator: castType, Expr: expression}}
}

func (stmt *CastExpression) Process(visitor Visitor, context any) (any, error) {
	return visitor.ProcessCastExpr(stmt, context)
}

// ------------------- MARK: IncludeExpression -------------------

type IncludeExpression struct {
	*Expression
	Expr IExpression
}

func NewIncludeExpr(id int64, pos *position.Position, expression IExpression) *IncludeExpression {
	return &IncludeExpression{Expression: NewExpr(id, IncludeExpr, pos), Expr: expression}
}

func (stmt *IncludeExpression) Process(visitor Visitor, context any) (any, error) {
	return visitor.ProcessIncludeExpr(stmt, context)
}

// ------------------- MARK: IncludeOnceExpression -------------------

type IncludeOnceExpression struct {
	*IncludeExpression
}

func NewIncludeOnceExpr(id int64, pos *position.Position, expression IExpression) *IncludeOnceExpression {
	return &IncludeOnceExpression{&IncludeExpression{Expression: NewExpr(id, IncludeOnceExpr, pos), Expr: expression}}
}

func (stmt *IncludeOnceExpression) Process(visitor Visitor, context any) (any, error) {
	return visitor.ProcessIncludeOnceExpr(stmt, context)
}

// ------------------- MARK: RequireExpression -------------------

type RequireExpression struct {
	*IncludeExpression
}

func NewRequireExpr(id int64, pos *position.Position, expression IExpression) *RequireExpression {
	return &RequireExpression{&IncludeExpression{Expression: NewExpr(id, RequireExpr, pos), Expr: expression}}
}

func (stmt *RequireExpression) Process(visitor Visitor, context any) (any, error) {
	return visitor.ProcessRequireExpr(stmt, context)
}

// ------------------- MARK: RequireOnceExpression -------------------

type RequireOnceExpression struct {
	*IncludeExpression
}

func NewRequireOnceExpr(id int64, pos *position.Position, expression IExpression) *RequireOnceExpression {
	return &RequireOnceExpression{&IncludeExpression{Expression: NewExpr(id, RequireOnceExpr, pos), Expr: expression}}
}

func (stmt *RequireOnceExpression) Process(visitor Visitor, context any) (any, error) {
	return visitor.ProcessRequireOnceExpr(stmt, context)
}
