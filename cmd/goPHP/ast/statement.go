package ast

import (
	"GoPHP/cmd/goPHP/position"
)

// ------------------- MARK: Statement -------------------

type IStatement interface {
	GetId() int64
	GetKind() NodeType
	GetPosition() *position.Position
	Process(visitor Visitor, context any) (any, error)
}

type Statement struct {
	id   int64
	kind NodeType
	pos  *position.Position
}

func NewEmptyStmt() *Statement {
	return &Statement{kind: EmptyNode}
}

func NewStmt(id int64, kind NodeType, pos *position.Position) *Statement {
	return &Statement{id: id, kind: kind, pos: pos}
}

func (stmt *Statement) GetId() int64 {
	return stmt.id
}

func (stmt *Statement) GetKind() NodeType {
	return stmt.kind
}

func (stmt *Statement) GetPosition() *position.Position {
	return stmt.pos
}

func (stmt *Statement) Process(visitor Visitor, context any) (any, error) {
	return visitor.ProcessStmt(stmt, context)
}

// ------------------- MARK: FunctionDefinitionStatement -------------------

type FunctionParameter struct {
	Type []string
	Name string
}

type FunctionDefinitionStatement struct {
	*Statement
	FunctionName string
	Params       []FunctionParameter
	Body         *CompoundStatement
	ReturnType   []string
}

func NewFunctionDefinitionStmt(id int64, pos *position.Position, functionName string, params []FunctionParameter, body *CompoundStatement, returnType []string) *FunctionDefinitionStatement {
	return &FunctionDefinitionStatement{Statement: NewStmt(id, FunctionDefinitionStmt, pos),
		FunctionName: functionName, Params: params, Body: body, ReturnType: returnType,
	}
}

func (stmt *FunctionDefinitionStatement) Process(visitor Visitor, context any) (any, error) {
	return visitor.ProcessFunctionDefinitionStmt(stmt, context)
}

// ------------------- MARK: IfStatement -------------------

type IfStatement struct {
	*Statement
	Condition IExpression
	IfBlock   IStatement
	ElseIf    []*IfStatement
	ElseBlock IStatement
}

func NewIfStmt(id int64, pos *position.Position, condition IExpression, ifBlock IStatement, elseIf []*IfStatement, elseBlock IStatement) *IfStatement {
	return &IfStatement{Statement: NewStmt(id, IfStmt, pos), Condition: condition, IfBlock: ifBlock, ElseIf: elseIf, ElseBlock: elseBlock}
}

func (stmt *IfStatement) Process(visitor Visitor, context any) (any, error) {
	return visitor.ProcessIfStmt(stmt, context)
}

// ------------------- MARK: WhileStatement -------------------

type WhileStatement struct {
	*IfStatement
}

func NewWhileStmt(id int64, pos *position.Position, condition IExpression, block IStatement) *WhileStatement {
	return &WhileStatement{&IfStatement{Statement: NewStmt(id, WhileStmt, pos), Condition: condition, IfBlock: block}}
}

func (stmt *WhileStatement) Process(visitor Visitor, context any) (any, error) {
	return visitor.ProcessWhileStmt(stmt, context)
}

// ------------------- MARK: DoStatement -------------------

type DoStatement struct {
	*IfStatement
}

func NewDoStmt(id int64, pos *position.Position, condition IExpression, block IStatement) *DoStatement {
	return &DoStatement{&IfStatement{Statement: NewStmt(id, DoStmt, pos), Condition: condition, IfBlock: block}}
}

func (stmt *DoStatement) Process(visitor Visitor, context any) (any, error) {
	return visitor.ProcessDoStmt(stmt, context)
}

// ------------------- MARK: CompoundStatement -------------------

type CompoundStatement struct {
	*Statement
	Statements []IStatement
}

func NewCompoundStmt(id int64, statements []IStatement) *CompoundStatement {
	return &CompoundStatement{Statement: NewStmt(id, CompoundStmt, nil), Statements: statements}
}

func (stmt *CompoundStatement) Process(visitor Visitor, context any) (any, error) {
	return visitor.ProcessCompoundStmt(stmt, context)
}

// ------------------- MARK: EchoStatement -------------------

type EchoStatement struct {
	*Statement
	Expressions []IExpression
}

func NewEchoStmt(id int64, pos *position.Position, expressions []IExpression) *EchoStatement {
	return &EchoStatement{Statement: NewStmt(id, EchoStmt, pos), Expressions: expressions}
}

func (stmt *EchoStatement) Process(visitor Visitor, context any) (any, error) {
	return visitor.ProcessEchoStmt(stmt, context)
}

// ------------------- MARK: ConstDeclarationStatement -------------------

type ConstDeclarationStatement struct {
	*Statement
	Name  string
	Value IExpression
}

func NewConstDeclarationStmt(id int64, pos *position.Position, name string, value IExpression) *ConstDeclarationStatement {
	return &ConstDeclarationStatement{Statement: NewStmt(id, ConstDeclarationStmt, pos), Name: name, Value: value}
}

func (stmt *ConstDeclarationStatement) Process(visitor Visitor, context any) (any, error) {
	return visitor.ProcessConstDeclarationStmt(stmt, context)
}

// ------------------- MARK: ExpressionStatement -------------------

type ExpressionStatement struct {
	*Statement
	Expr IExpression
}

func NewExpressionStmt(id int64, expr IExpression) *ExpressionStatement {
	return &ExpressionStatement{Statement: NewStmt(id, ExpressionStmt, expr.GetPosition()), Expr: expr}
}

func (stmt *ExpressionStatement) Process(visitor Visitor, context any) (any, error) {
	return visitor.ProcessExpressionStmt(stmt, context)
}

// ------------------- MARK: BreakStatement -------------------

type BreakStatement struct {
	*ExpressionStatement
}

func NewBreakStmt(id int64, pos *position.Position, expr IExpression) *BreakStatement {
	return &BreakStatement{&ExpressionStatement{Statement: NewStmt(id, BreakStmt, pos), Expr: expr}}
}

func (stmt *BreakStatement) Process(visitor Visitor, context any) (any, error) {
	return visitor.ProcessBreakStmt(stmt, context)
}

// ------------------- MARK: ContinueStatement -------------------

type ContinueStatement struct {
	*ExpressionStatement
}

func NewContinueStmt(id int64, pos *position.Position, expr IExpression) *ContinueStatement {
	return &ContinueStatement{&ExpressionStatement{Statement: NewStmt(id, ContinueStmt, pos), Expr: expr}}
}

func (stmt *ContinueStatement) Process(visitor Visitor, context any) (any, error) {
	return visitor.ProcessContinueStmt(stmt, context)
}

// ------------------- MARK: ReturnStatement -------------------

type ReturnStatement struct {
	*ExpressionStatement
}

func NewReturnStmt(id int64, pos *position.Position, expr IExpression) *ReturnStatement {
	return &ReturnStatement{&ExpressionStatement{Statement: NewStmt(id, ReturnStmt, pos), Expr: expr}}
}

func (stmt *ReturnStatement) Process(visitor Visitor, context any) (any, error) {
	return visitor.ProcessReturnStmt(stmt, context)
}
