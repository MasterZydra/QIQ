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

type IFunctionDefinitionStatement interface {
	IStatement
	GetFunctionName() string
	GetParams() []FunctionParameter
	GetBody() ICompoundStatement
	GetReturnType() []string
}

type FunctionDefinitionStatement struct {
	stmt         IStatement
	functionName string
	params       []FunctionParameter
	body         ICompoundStatement
	returnType   []string
}

func NewFunctionDefinitionStmt(id int64, pos *position.Position, functionName string, params []FunctionParameter, body ICompoundStatement, returnType []string) *FunctionDefinitionStatement {
	return &FunctionDefinitionStatement{stmt: NewStmt(id, FunctionDefinitionStmt, pos),
		functionName: functionName, params: params, body: body, returnType: returnType,
	}
}

func (stmt *FunctionDefinitionStatement) GetId() int64 {
	return stmt.stmt.GetId()
}

func (stmt *FunctionDefinitionStatement) GetKind() NodeType {
	return stmt.stmt.GetKind()
}

func (stmt *FunctionDefinitionStatement) GetPosition() *position.Position {
	return stmt.stmt.GetPosition()
}

func (stmt *FunctionDefinitionStatement) GetFunctionName() string {
	return stmt.functionName
}

func (stmt *FunctionDefinitionStatement) GetParams() []FunctionParameter {
	return stmt.params
}

func (stmt *FunctionDefinitionStatement) GetBody() ICompoundStatement {
	return stmt.body
}

func (stmt *FunctionDefinitionStatement) GetReturnType() []string {
	return stmt.returnType
}

func (stmt *FunctionDefinitionStatement) String() string {
	return fmt.Sprintf("{%s - name: %s, params: %s, body: %s, returnType: %s}",
		stmt.GetKind(), stmt.functionName, stmt.params, stmt.body, stmt.returnType,
	)
}

func StmtToFunctionDefinitionStmt(stmt IStatement) IFunctionDefinitionStatement {
	var i interface{} = stmt
	return i.(IFunctionDefinitionStatement)
}

// ------------------- MARK: IfStatement -------------------

type IIfStatement interface {
	IStatement
	GetCondition() IExpression
	GetIfBlock() IStatement
	GetElseIf() []IIfStatement
	GetElseBlock() IStatement
}

type IfStatement struct {
	stmt      IStatement
	condition IExpression
	ifBlock   IStatement
	elseIf    []IIfStatement
	elseBlock IStatement
}

func NewDoStmt(id int64, pos *position.Position, condition IExpression, block IStatement) *IfStatement {
	return &IfStatement{stmt: NewStmt(id, DoStmt, pos), condition: condition, ifBlock: block}
}

func NewWhileStmt(id int64, pos *position.Position, condition IExpression, block IStatement) *IfStatement {
	return &IfStatement{stmt: NewStmt(id, WhileStmt, pos), condition: condition, ifBlock: block}
}

func NewIfStmt(id int64, pos *position.Position, condition IExpression, ifBlock IStatement, elseIf []IIfStatement, elseBlock IStatement) *IfStatement {
	return &IfStatement{stmt: NewStmt(id, IfStmt, pos), condition: condition, ifBlock: ifBlock, elseIf: elseIf, elseBlock: elseBlock}
}

func (stmt *IfStatement) GetId() int64 {
	return stmt.stmt.GetId()
}

func (stmt *IfStatement) GetKind() NodeType {
	return stmt.stmt.GetKind()
}

func (stmt *IfStatement) GetPosition() *position.Position {
	return stmt.stmt.GetPosition()
}

func (stmt *IfStatement) GetCondition() IExpression {
	return stmt.condition
}

func (stmt *IfStatement) GetIfBlock() IStatement {
	return stmt.ifBlock
}

func (stmt *IfStatement) GetElseIf() []IIfStatement {
	return stmt.elseIf
}

func (stmt *IfStatement) GetElseBlock() IStatement {
	return stmt.elseBlock
}

func (stmt *IfStatement) String() string {
	return fmt.Sprintf("{%s - condition: %s, ifBlock: %s, elseIf: %s, else: %s}",
		stmt.GetKind(), stmt.condition, stmt.ifBlock, stmt.elseIf, stmt.elseBlock)
}

func StmtToIfStmt(stmt IStatement) IIfStatement {
	var i interface{} = stmt
	return i.(IIfStatement)
}

// ------------------- MARK: CompoundStatement -------------------

type ICompoundStatement interface {
	IStatement
	GetStatements() []IStatement
}

type CompoundStatement struct {
	stmt       IStatement
	statements []IStatement
}

func NewCompoundStmt(id int64, statements []IStatement) *CompoundStatement {
	return &CompoundStatement{stmt: NewStmt(id, CompoundStmt, nil), statements: statements}
}

func (stmt *CompoundStatement) GetId() int64 {
	return stmt.stmt.GetId()
}

func (stmt *CompoundStatement) GetKind() NodeType {
	return stmt.stmt.GetKind()
}

func (stmt *CompoundStatement) GetPosition() *position.Position {
	return stmt.stmt.GetPosition()
}

func (stmt *CompoundStatement) GetStatements() []IStatement {
	return stmt.statements
}

func (stmt *CompoundStatement) String() string {
	return fmt.Sprintf("{%s - %s}", stmt.GetKind(), stmt.statements)
}

func StmtToCompoundStmt(stmt IStatement) ICompoundStatement {
	var i interface{} = stmt
	return i.(ICompoundStatement)
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

func NewEchoStmt(id int64, pos *position.Position, expressions []IExpression) *EchoStatement {
	return &EchoStatement{stmt: NewStmt(id, EchoStmt, pos), expressions: expressions}
}

func (stmt *EchoStatement) GetId() int64 {
	return stmt.stmt.GetId()
}

func (stmt *EchoStatement) GetKind() NodeType {
	return stmt.stmt.GetKind()
}

func (stmt *EchoStatement) GetPosition() *position.Position {
	return stmt.stmt.GetPosition()
}

func (stmt *EchoStatement) GetExpressions() []IExpression {
	return stmt.expressions
}

func (stmt *EchoStatement) String() string {
	return fmt.Sprintf("{%s - %s}", stmt.GetKind(), stmt.expressions)
}

func StmtToEchoStmt(stmt IStatement) IEchoStatement {
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

func NewConstDeclarationStmt(id int64, pos *position.Position, name string, value IExpression) *ConstDeclarationStatement {
	return &ConstDeclarationStatement{stmt: NewStmt(id, ConstDeclarationStmt, pos), name: name, value: value}
}

func (stmt *ConstDeclarationStatement) GetId() int64 {
	return stmt.stmt.GetId()
}

func (stmt *ConstDeclarationStatement) GetKind() NodeType {
	return stmt.stmt.GetKind()
}

func (stmt *ConstDeclarationStatement) GetPosition() *position.Position {
	return stmt.stmt.GetPosition()
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

func StmtToConstDeclStmt(stmt IStatement) IConstDeclarationStatement {
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

func NewReturnStmt(id int64, pos *position.Position, expr IExpression) *ExpressionStatement {
	return &ExpressionStatement{stmt: NewStmt(id, ReturnStmt, pos), expr: expr}
}

func NewExpressionStmt(id int64, expr IExpression) *ExpressionStatement {
	return &ExpressionStatement{stmt: NewStmt(id, ExpressionStmt, expr.GetPosition()), expr: expr}
}

func (stmt *ExpressionStatement) GetId() int64 {
	return stmt.stmt.GetId()
}

func (stmt *ExpressionStatement) GetKind() NodeType {
	return stmt.stmt.GetKind()
}

func (stmt *ExpressionStatement) GetPosition() *position.Position {
	return stmt.stmt.GetPosition()
}

func (stmt *ExpressionStatement) GetExpression() IExpression {
	return stmt.expr
}

func (stmt *ExpressionStatement) String() string {
	return fmt.Sprintf("{%s - %s}", stmt.GetKind(), stmt.expr)
}

func StmtToExprStmt(stmt IStatement) IExpressionStatement {
	var i interface{} = stmt
	return i.(IExpressionStatement)
}
