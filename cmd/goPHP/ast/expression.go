package ast

import "fmt"

// ------------------- MARK: Expression -------------------

type IExpression interface {
	IStatement
}

type Expression struct {
	kind NodeType
}

func NewExpression(kind NodeType) *Expression {
	return &Expression{kind: kind}
}

func (expr *Expression) GetKind() NodeType {
	return expr.kind
}

func (expr *Expression) String() string {
	return fmt.Sprintf("{%s}", expr.GetKind())
}

func NewEmptyExpression() *Expression {
	return NewExpression(EmptyNode)
}

// ------------------- MARK: TextExpression -------------------

type ITextExpression interface {
	IExpression
	GetValue() string
}

type TextExpression struct {
	expr  IExpression
	value string
}

func NewTextExpression(value string) *TextExpression {
	return &TextExpression{expr: NewExpression(TextNode), value: value}
}

func (expr *TextExpression) GetKind() NodeType {
	return expr.expr.GetKind()
}

func (expr *TextExpression) GetValue() string {
	return expr.value
}

func (expr *TextExpression) String() string {
	return fmt.Sprintf("{%s - value: \"%s\" }", expr.GetKind(), expr.value)
}

func ExprToTextExpr(expr IExpression) ITextExpression {
	var i interface{} = expr
	return i.(ITextExpression)
}

// ------------------- MARK: VariableNameExpression -------------------

type IVariableNameExpression interface {
	IExpression
	GetVariableName() string
}

type VariableNameExpression struct {
	expr         IExpression
	variableName string
}

func NewVariableNameExpression(variableName string) *VariableNameExpression {
	return &VariableNameExpression{expr: NewExpression(VariableNameExpr), variableName: variableName}
}

func (expr *VariableNameExpression) GetKind() NodeType {
	return expr.expr.GetKind()
}

func (expr *VariableNameExpression) GetVariableName() string {
	return expr.variableName
}

func (expr *VariableNameExpression) String() string {
	return fmt.Sprintf("{%s - variableName: \"%s\" }", expr.GetKind(), expr.variableName)
}

func ExprToVarNameExpr(expr IExpression) IVariableNameExpression {
	var i interface{} = expr
	return i.(IVariableNameExpression)
}

// ------------------- MARK: SimpleVariableExpression -------------------

type ISimpleVariableExpression interface {
	IExpression
	GetVariableName() IExpression
}

type SimpleVariableExpression struct {
	expr         IExpression
	variableName IExpression
}

func NewSimpleVariableExpression(variableName IExpression) *SimpleVariableExpression {
	return &SimpleVariableExpression{expr: NewExpression(SimpleVariableExpr), variableName: variableName}
}

func (expr *SimpleVariableExpression) GetKind() NodeType {
	return expr.expr.GetKind()
}

func (expr *SimpleVariableExpression) GetVariableName() IExpression {
	return expr.variableName
}

func (expr *SimpleVariableExpression) String() string {
	return fmt.Sprintf("{%s - variableName: \"%s\" }", expr.GetKind(), expr.variableName)
}

func ExprToSimpleVarExpr(expr IExpression) ISimpleVariableExpression {
	var i interface{} = expr
	return i.(ISimpleVariableExpression)
}

// ------------------- MARK: BooleanLiteralExpression -------------------

type IBooleanLiteralExpression interface {
	IExpression
	GetValue() bool
}

type BooleanLiteralExpression struct {
	expr  IExpression
	value bool
}

func NewBooleanLiteralExpression(value bool) *BooleanLiteralExpression {
	return &BooleanLiteralExpression{expr: NewExpression(BooleanLiteralExpr), value: value}
}

func (expr *BooleanLiteralExpression) GetKind() NodeType {
	return expr.expr.GetKind()
}

func (expr *BooleanLiteralExpression) GetValue() bool {
	return expr.value
}

func (expr *BooleanLiteralExpression) String() string {
	return fmt.Sprintf("{%s - value: %t }", expr.GetKind(), expr.value)
}

func ExprToBoolLitExpr(expr IExpression) IBooleanLiteralExpression {
	var i interface{} = expr
	return i.(IBooleanLiteralExpression)
}

// ------------------- MARK: IntegerLiteralExpression -------------------

type IIntegerLiteralExpression interface {
	IExpression
	GetValue() int64
}

type IntegerLiteralExpression struct {
	expr  IExpression
	value int64
}

func NewIntegerLiteralExpression(value int64) *IntegerLiteralExpression {
	return &IntegerLiteralExpression{expr: NewExpression(IntegerLiteralExpr), value: value}
}

func (expr *IntegerLiteralExpression) GetKind() NodeType {
	return expr.expr.GetKind()
}

func (expr *IntegerLiteralExpression) GetValue() int64 {
	return expr.value
}

func (expr *IntegerLiteralExpression) String() string {
	return fmt.Sprintf("{%s - value: %d }", expr.GetKind(), expr.value)
}

func ExprToIntLitExpr(expr IExpression) IIntegerLiteralExpression {
	var i interface{} = expr
	return i.(IIntegerLiteralExpression)
}

// ------------------- MARK: FloatingLiteralExpression -------------------

type IFloatingLiteralExpression interface {
	IExpression
	GetValue() float64
}

type FloatingLiteralExpression struct {
	expr  IExpression
	value float64
}

func NewFloatingLiteralExpression(value float64) *FloatingLiteralExpression {
	return &FloatingLiteralExpression{expr: NewExpression(FloatingLiteralExpr), value: value}
}

func (expr *FloatingLiteralExpression) GetKind() NodeType {
	return expr.expr.GetKind()
}

func (expr *FloatingLiteralExpression) GetValue() float64 {
	return expr.value
}

func (expr *FloatingLiteralExpression) String() string {
	return fmt.Sprintf("{%s - value: %f }", expr.GetKind(), expr.value)
}

func ExprToFloatLitExpr(expr IExpression) IFloatingLiteralExpression {
	var i interface{} = expr
	return i.(IFloatingLiteralExpression)
}

// ------------------- MARK: StringLiteralExpression -------------------

type StringType string

const (
	SingleQuotedString StringType = "SingleQuotedString"
	DoubleQuotedString StringType = "DoubleQuotedString"
)

type IStringLiteralExpression interface {
	IExpression
	GetStringType() StringType
	GetValue() string
}

type StringLiteralExpression struct {
	expr       IExpression
	stringType StringType
	value      string
}

func NewStringLiteralExpression(value string, stringType StringType) *StringLiteralExpression {
	return &StringLiteralExpression{expr: NewExpression(StringLiteralExpr), value: value, stringType: stringType}
}

func (expr *StringLiteralExpression) GetKind() NodeType {
	return expr.expr.GetKind()
}

func (expr *StringLiteralExpression) GetStringType() StringType {
	return expr.stringType
}

func (expr *StringLiteralExpression) GetValue() string {
	return expr.value
}

func (expr *StringLiteralExpression) String() string {
	return fmt.Sprintf("{%s - type: \"%s\" value: \"%s\" }", expr.GetKind(), expr.stringType, expr.value)
}

func ExprToStrLitExpr(expr IExpression) IStringLiteralExpression {
	var i interface{} = expr
	return i.(IStringLiteralExpression)
}

// ------------------- MARK: SimpleAssignmentExpression -------------------

type ISimpleAssignmentExpression interface {
	IExpression
	GetVariable() IExpression
	GetValue() IExpression
}

type SimpleAssignmentExpression struct {
	expr     IExpression
	variable IExpression
	value    IExpression
}

func NewSimpleAssignmentExpression(variable IExpression, value IExpression) *SimpleAssignmentExpression {
	return &SimpleAssignmentExpression{expr: NewExpression(SimpleAssignmentExpr), variable: variable, value: value}
}

func (expr *SimpleAssignmentExpression) GetKind() NodeType {
	return expr.expr.GetKind()
}

func (expr *SimpleAssignmentExpression) GetVariable() IExpression {
	return expr.variable
}

func (expr *SimpleAssignmentExpression) GetValue() IExpression {
	return expr.value
}

func (expr *SimpleAssignmentExpression) String() string {
	return fmt.Sprintf("{%s - variable: %s, value: %s }", expr.GetKind(), expr.variable, expr.value)
}

func ExprToSimpleAssignExpr(expr IExpression) ISimpleAssignmentExpression {
	var i interface{} = expr
	return i.(ISimpleAssignmentExpression)
}

// ------------------- MARK: CompoundAssignmentExpression -------------------

type ICompoundAssignmentExpression interface {
	IExpression
	GetVariable() IExpression
	GetOperator() string
	GetValue() IExpression
}

type CompoundAssignmentExpression struct {
	expr     IExpression
	variable IExpression
	operator string
	value    IExpression
}

func NewCompoundAssignmentExpression(variable IExpression, operator string, value IExpression) *CompoundAssignmentExpression {
	return &CompoundAssignmentExpression{
		expr: NewExpression(CompoundAssignmentExpr), variable: variable, operator: operator, value: value,
	}
}

func (expr *CompoundAssignmentExpression) GetKind() NodeType {
	return expr.expr.GetKind()
}

func (expr *CompoundAssignmentExpression) GetVariable() IExpression {
	return expr.variable
}

func (expr *CompoundAssignmentExpression) GetOperator() string {
	return expr.operator
}

func (expr *CompoundAssignmentExpression) GetValue() IExpression {
	return expr.value
}

func (expr *CompoundAssignmentExpression) String() string {
	return fmt.Sprintf(
		"{%s - variable: %s, operator: \"%s\", value: %s }",
		expr.GetKind(), expr.variable, expr.operator, expr.value,
	)
}

func ExprToCompoundAssignExpr(expr IExpression) ICompoundAssignmentExpression {
	var i interface{} = expr
	return i.(ICompoundAssignmentExpression)
}

// ------------------- MARK: ConditionalExpression -------------------

type IConditionalExpression interface {
	IExpression
	GetCondition() IExpression
	GetIfExpr() IExpression
	GetElseExpr() IExpression
}

type ConditionalExpression struct {
	expr     IExpression
	cond     IExpression
	ifExpr   IExpression
	elseExpr IExpression
}

func NewConditionalExpression(cond IExpression, ifExpr IExpression, elseExpr IExpression) *ConditionalExpression {
	return &ConditionalExpression{expr: NewExpression(ConditionalExpr), cond: cond, ifExpr: ifExpr, elseExpr: elseExpr}
}

func (expr *ConditionalExpression) GetKind() NodeType {
	return expr.expr.GetKind()
}

func (expr *ConditionalExpression) GetCondition() IExpression {
	return expr.cond
}

func (expr *ConditionalExpression) GetIfExpr() IExpression {
	return expr.ifExpr
}

func (expr *ConditionalExpression) GetElseExpr() IExpression {
	return expr.elseExpr
}

func (expr *ConditionalExpression) String() string {
	return fmt.Sprintf("{%s - condition: %s, ifExpr: %s, elseExpr: %s }", expr.GetKind(), expr.cond, expr.ifExpr, expr.elseExpr)
}

func ExprToCondExpr(expr IExpression) IConditionalExpression {
	var i interface{} = expr
	return i.(IConditionalExpression)
}
