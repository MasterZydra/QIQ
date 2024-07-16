package ast

import (
	"GoPHP/cmd/goPHP/position"
	"fmt"
)

// ------------------- MARK: Statement -------------------

type IStatement interface {
	GetId() int64
	GetKind() NodeType
	GetPosition() *position.Position
	String() string
}

type Statement struct {
	id   int64
	kind NodeType
	pos  *position.Position
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

func (stmt *Statement) String() string {
	return fmt.Sprintf("{%s}", stmt.GetKind())
}

func NewEmptyStmt() *Statement {
	return &Statement{kind: EmptyNode}
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

func (stmt *FunctionDefinitionStatement) String() string {
	return fmt.Sprintf("{%s - name: %s, params: %s, body: %s, returnType: %s}",
		stmt.GetKind(), stmt.FunctionName, stmt.Params, stmt.Body, stmt.ReturnType,
	)
}

// ------------------- MARK: IfStatement -------------------

type IfStatement struct {
	*Statement
	Condition IExpression
	IfBlock   IStatement
	ElseIf    []*IfStatement
	ElseBlock IStatement
}

func NewDoStmt(id int64, pos *position.Position, condition IExpression, block IStatement) *IfStatement {
	return &IfStatement{Statement: NewStmt(id, DoStmt, pos), Condition: condition, IfBlock: block}
}

func NewWhileStmt(id int64, pos *position.Position, condition IExpression, block IStatement) *IfStatement {
	return &IfStatement{Statement: NewStmt(id, WhileStmt, pos), Condition: condition, IfBlock: block}
}

func NewIfStmt(id int64, pos *position.Position, condition IExpression, ifBlock IStatement, elseIf []*IfStatement, elseBlock IStatement) *IfStatement {
	return &IfStatement{Statement: NewStmt(id, IfStmt, pos), Condition: condition, IfBlock: ifBlock, ElseIf: elseIf, ElseBlock: elseBlock}
}

func (stmt *IfStatement) String() string {
	return fmt.Sprintf("{%s - condition: %s, ifBlock: %s, elseIf: %s, else: %s}",
		stmt.GetKind(), stmt.Condition, stmt.IfBlock, stmt.ElseIf, stmt.ElseBlock)
}

// ------------------- MARK: CompoundStatement -------------------

type CompoundStatement struct {
	*Statement
	Statements []IStatement
}

func NewCompoundStmt(id int64, statements []IStatement) *CompoundStatement {
	return &CompoundStatement{Statement: NewStmt(id, CompoundStmt, nil), Statements: statements}
}

func (stmt *CompoundStatement) String() string {
	return fmt.Sprintf("{%s - %s}", stmt.GetKind(), stmt.Statements)
}

// ------------------- MARK: EchoStatement -------------------

type EchoStatement struct {
	*Statement
	Expressions []IExpression
}

func NewEchoStmt(id int64, pos *position.Position, expressions []IExpression) *EchoStatement {
	return &EchoStatement{Statement: NewStmt(id, EchoStmt, pos), Expressions: expressions}
}

func (stmt *EchoStatement) String() string {
	return fmt.Sprintf("{%s - %s}", stmt.GetKind(), stmt.Expressions)
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

func (stmt *ConstDeclarationStatement) String() string {
	return fmt.Sprintf("{%s - name: \"%s\" value: %s}", stmt.GetKind(), stmt.Name, stmt.Value)
}

// ------------------- MARK: ExpressionStatement -------------------

type ExpressionStatement struct {
	*Statement
	Expr IExpression
}

func NewBreakStmt(id int64, pos *position.Position, expr IExpression) *ExpressionStatement {
	return &ExpressionStatement{Statement: NewStmt(id, BreakStmt, pos), Expr: expr}
}

func NewContinueStmt(id int64, pos *position.Position, expr IExpression) *ExpressionStatement {
	return &ExpressionStatement{Statement: NewStmt(id, ContinueStmt, pos), Expr: expr}
}

func NewReturnStmt(id int64, pos *position.Position, expr IExpression) *ExpressionStatement {
	return &ExpressionStatement{Statement: NewStmt(id, ReturnStmt, pos), Expr: expr}
}

func NewExpressionStmt(id int64, expr IExpression) *ExpressionStatement {
	return &ExpressionStatement{Statement: NewStmt(id, ExpressionStmt, expr.GetPosition()), Expr: expr}
}

func (stmt *ExpressionStatement) String() string {
	return fmt.Sprintf("{%s - %s}", stmt.GetKind(), stmt.Expr)
}
