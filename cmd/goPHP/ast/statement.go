package ast

import "fmt"

// ------------------- MARK: Statement -------------------

type IStatement interface {
	GetKind() NodeType
	String() string
}

type Statement struct {
	kind NodeType
}

func NewStatement(kind NodeType) *Statement {
	return &Statement{kind: kind}
}

func (stmt *Statement) GetKind() NodeType {
	return stmt.kind
}

func (stmt *Statement) String() string {
	return fmt.Sprintf("{%s}", stmt.GetKind())
}

func NewEmptyStatement() *Statement {
	return &Statement{kind: EmptyNode}
}

// ------------------- MARK: EchoStatement -------------------

type IEchoStatement interface {
	IStatement
	GetExpressions() []IExpression
}

type EchoStatement struct {
	stmt        IStatement
	expressions []IExpression
}

func NewEchoStatement(expressions []IExpression) *EchoStatement {
	return &EchoStatement{stmt: NewStatement(EchoStmt), expressions: expressions}
}

func (stmt *EchoStatement) GetKind() NodeType {
	return stmt.stmt.GetKind()
}

func (stmt *EchoStatement) GetExpressions() []IExpression {
	return stmt.expressions
}

func (stmt *EchoStatement) String() string {
	return fmt.Sprintf("{%s - %s}", stmt.GetKind(), stmt.expressions)
}

func StmtToEchoStatement(stmt IStatement) IEchoStatement {
	var i interface{} = stmt
	return i.(IEchoStatement)
}

// ------------------- MARK: ConstDeclarationStatement -------------------

type IConstDeclarationStatement interface {
	IStatement
	GetName() string
	GetValue() IExpression
}

type ConstDeclarationStatement struct {
	stmt  IStatement
	name  string
	value IExpression
}

func NewConstDeclarationStatement(name string, value IExpression) *ConstDeclarationStatement {
	return &ConstDeclarationStatement{stmt: NewStatement(ConstDeclarationStmt), name: name, value: value}
}

func (stmt *ConstDeclarationStatement) GetKind() NodeType {
	return stmt.stmt.GetKind()
}

func (stmt *ConstDeclarationStatement) GetName() string {
	return stmt.name
}

func (stmt *ConstDeclarationStatement) GetValue() IExpression {
	return stmt.value
}

func (stmt *ConstDeclarationStatement) String() string {
	return fmt.Sprintf("{%s - name: \"%s\" value: %s}", stmt.GetKind(), stmt.name, stmt.value)
}

func StmtToConstDeclStatement(stmt IStatement) IConstDeclarationStatement {
	var i interface{} = stmt
	return i.(IConstDeclarationStatement)
}

// ------------------- MARK: ExpressionStatement -------------------

type IExpressionStatement interface {
	IStatement
	GetExpression() IExpression
}

type ExpressionStatement struct {
	stmt IStatement
	expr IExpression
}

func NewExpressionStatement(expr IExpression) *ExpressionStatement {
	return &ExpressionStatement{stmt: NewStatement(ExpressionStmt), expr: expr}
}

func (stmt *ExpressionStatement) GetKind() NodeType {
	return stmt.stmt.GetKind()
}

func (stmt *ExpressionStatement) GetExpression() IExpression {
	return stmt.expr
}

func (stmt *ExpressionStatement) String() string {
	return fmt.Sprintf("{%s - %s}", stmt.GetKind(), stmt.expr)
}

func StmtToExprStatement(stmt IStatement) IExpressionStatement {
	var i interface{} = stmt
	return i.(IExpressionStatement)
}
