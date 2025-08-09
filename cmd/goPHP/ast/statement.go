package ast

import (
	"GoPHP/cmd/goPHP/position"
	"slices"
)

// -------------------------------------- Statement -------------------------------------- MARK: Statement

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
	if stmt.pos == nil {
		return &position.Position{}
	}
	return stmt.pos
}

func (stmt *Statement) Process(visitor Visitor, context any) (any, error) {
	return visitor.ProcessStmt(stmt, context)
}

// -------------------------------------- MethodDefinitionStatement -------------------------------------- MARK: MethodDefinitionStatement

type MethodDefinitionStatement struct {
	*Statement
	Modifiers  []string
	Name       string
	Params     []FunctionParameter
	Body       *CompoundStatement
	ReturnType []string
}

func NewMethodDefinitionStmt(id int64, pos *position.Position, name string, modifiers []string, params []FunctionParameter, body *CompoundStatement, returnType []string) *MethodDefinitionStatement {
	return &MethodDefinitionStatement{Statement: NewStmt(id, MethodDefinitionStmt, pos),
		Name:       name,
		Modifiers:  modifiers,
		Params:     params,
		Body:       body,
		ReturnType: returnType,
	}
}

func (stmt *MethodDefinitionStatement) Process(visitor Visitor, context any) (any, error) {
	panic("MethodDefinitionStatement.Process should not be called")
}

// -------------------------------------- FunctionDefinitionStatement -------------------------------------- MARK: FunctionDefinitionStatement

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

// -------------------------------------- ForStatement -------------------------------------- MARK: ForStatement

type ForStatement struct {
	*Statement
	Initializer *CompoundStatement
	Control     *CompoundStatement
	EndOfLoop   *CompoundStatement
	Block       IStatement
}

func NewForStmt(id int64, pos *position.Position, initializer *CompoundStatement, control *CompoundStatement, endOfLoop *CompoundStatement, block IStatement) *ForStatement {
	return &ForStatement{Statement: NewStmt(id, ForStmt, pos), Initializer: initializer, Control: control, EndOfLoop: endOfLoop, Block: block}
}

func (stmt *ForStatement) Process(visitor Visitor, context any) (any, error) {
	return visitor.ProcessForStmt(stmt, context)
}

// -------------------------------------- WhileStatement -------------------------------------- MARK: WhileStatement

type WhileStatement struct {
	*Statement
	Condition IExpression
	Block     IStatement
}

func NewWhileStmt(id int64, pos *position.Position, condition IExpression, block IStatement) *WhileStatement {
	return &WhileStatement{Statement: NewStmt(id, WhileStmt, pos), Condition: condition, Block: block}
}

func (stmt *WhileStatement) Process(visitor Visitor, context any) (any, error) {
	return visitor.ProcessWhileStmt(stmt, context)
}

// -------------------------------------- DoStatement -------------------------------------- MARK: DoStatement

type DoStatement struct {
	*WhileStatement
}

func NewDoStmt(id int64, pos *position.Position, condition IExpression, block IStatement) *DoStatement {
	return &DoStatement{&WhileStatement{Statement: NewStmt(id, DoStmt, pos), Condition: condition, Block: block}}
}

func (stmt *DoStatement) Process(visitor Visitor, context any) (any, error) {
	return visitor.ProcessDoStmt(stmt, context)
}

// -------------------------------------- IfStatement -------------------------------------- MARK: IfStatement

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

// -------------------------------------- CompoundStatement -------------------------------------- MARK: CompoundStatement

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

// -------------------------------------- EchoStatement -------------------------------------- MARK: EchoStatement

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

// -------------------------------------- ClassConstDeclarationStatement -------------------------------------- MARK: ClassConstDeclarationStatement

type ClassConstDeclarationStatement struct {
	*Statement
	Name      string
	Value     IExpression
	Visiblity string
}

func NewClassConstDeclarationStmt(id int64, pos *position.Position, name string, value IExpression, visibility string) *ClassConstDeclarationStatement {
	return &ClassConstDeclarationStatement{Statement: NewStmt(id, ClassConstDeclarationStmt, pos), Name: name, Value: value, Visiblity: visibility}
}

func (stmt *ClassConstDeclarationStatement) Process(visitor Visitor, context any) (any, error) {
	panic("ClassConstDeclarationStatement.Process should not be called")
}

// -------------------------------------- TraitUseStatement -------------------------------------- MARK: TraitUseStatement

type TraitUseStatement struct {
	*Statement
	Name string
}

func NewTraitUseStmt(id int64, pos *position.Position, name string) *TraitUseStatement {
	return &TraitUseStatement{Statement: NewStmt(id, ClassConstDeclarationStmt, pos), Name: name}
}

func (stmt *TraitUseStatement) Process(visitor Visitor, context any) (any, error) {
	panic("TraitUseStatement.Process should not be called")
}

// -------------------------------------- ConstDeclarationStatement -------------------------------------- MARK: ConstDeclarationStatement

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

// -------------------------------------- ExpressionStatement -------------------------------------- MARK: ExpressionStatement

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

// -------------------------------------- BreakStatement -------------------------------------- MARK: BreakStatement

type BreakStatement struct {
	*ExpressionStatement
}

func NewBreakStmt(id int64, pos *position.Position, expr IExpression) *BreakStatement {
	return &BreakStatement{&ExpressionStatement{Statement: NewStmt(id, BreakStmt, pos), Expr: expr}}
}

func (stmt *BreakStatement) Process(visitor Visitor, context any) (any, error) {
	return visitor.ProcessBreakStmt(stmt, context)
}

// -------------------------------------- ContinueStatement -------------------------------------- MARK: ContinueStatement

type ContinueStatement struct {
	*ExpressionStatement
}

func NewContinueStmt(id int64, pos *position.Position, expr IExpression) *ContinueStatement {
	return &ContinueStatement{&ExpressionStatement{Statement: NewStmt(id, ContinueStmt, pos), Expr: expr}}
}

func (stmt *ContinueStatement) Process(visitor Visitor, context any) (any, error) {
	return visitor.ProcessContinueStmt(stmt, context)
}

// -------------------------------------- ReturnStatement -------------------------------------- MARK: ReturnStatement

type ReturnStatement struct {
	*ExpressionStatement
}

func NewReturnStmt(id int64, pos *position.Position, expr IExpression) *ReturnStatement {
	return &ReturnStatement{&ExpressionStatement{Statement: NewStmt(id, ReturnStmt, pos), Expr: expr}}
}

func (stmt *ReturnStatement) Process(visitor Visitor, context any) (any, error) {
	return visitor.ProcessReturnStmt(stmt, context)
}

// -------------------------------------- GlobalDeclarationStatement -------------------------------------- MARK: GlobalDeclarationStatement

type GlobalDeclarationStatement struct {
	*Statement
	Variables []IExpression
}

func NewGlobalDeclarationStmt(id int64, pos *position.Position, variables []IExpression) *GlobalDeclarationStatement {
	return &GlobalDeclarationStatement{Statement: NewStmt(id, GlobalDeclarationStmt, pos), Variables: variables}
}

func (stmt *GlobalDeclarationStatement) Process(visitor Visitor, context any) (any, error) {
	return visitor.ProcessGlobalDeclarationStmt(stmt, context)
}

// -------------------------------------- PropertyDeclarationStatement -------------------------------------- MARK: PropertyDeclarationStatement

type PropertyDeclarationStatement struct {
	*Statement
	Visibility   string
	IsStatic     bool
	Name         string
	Type         []string
	InitialValue IExpression
}

func NewPropertyDeclarationStmt(id int64, pos *position.Position, name, visibility string, isStatic bool, pType []string, initialValue IExpression) *PropertyDeclarationStatement {
	return &PropertyDeclarationStatement{
		Statement:    NewStmt(id, PropertyDeclarationStmt, pos),
		Name:         name,
		Visibility:   visibility,
		IsStatic:     isStatic,
		Type:         pType,
		InitialValue: initialValue,
	}
}

func (stmt *PropertyDeclarationStatement) Process(visitor Visitor, context any) (any, error) {
	panic("PropertyDeclarationStatement.Process should not be called")
}

// -------------------------------------- ClassDeclarationStatement -------------------------------------- MARK: ClassDeclarationStatement

type ClassDeclarationStatement struct {
	*Statement
	IsAbstract     bool
	IsFinal        bool
	Name           string
	BaseClass      string
	Interfaces     []string
	Constants      map[string]*ClassConstDeclarationStatement
	Methods        map[string]*MethodDefinitionStatement
	PropertieNames []string
	Properties     map[string]*PropertyDeclarationStatement
	Traits         []*TraitUseStatement
}

func NewClassDeclarationStmt(id int64, pos *position.Position, name string, isAbstract, isFinal bool) *ClassDeclarationStatement {
	return &ClassDeclarationStatement{
		Statement:      NewStmt(id, ClassDeclarationStmt, pos),
		Name:           name,
		IsAbstract:     isAbstract,
		IsFinal:        isFinal,
		Interfaces:     []string{},
		Constants:      map[string]*ClassConstDeclarationStatement{},
		Methods:        map[string]*MethodDefinitionStatement{},
		PropertieNames: []string{},
		Properties:     map[string]*PropertyDeclarationStatement{},
		Traits:         []*TraitUseStatement{},
	}
}

func (stmt *ClassDeclarationStatement) Process(visitor Visitor, context any) (any, error) {
	return visitor.ProcessClassDeclarationStmt(stmt, context)
}

func (stmt *ClassDeclarationStatement) AddConst(constStmt *ClassConstDeclarationStatement) {
	stmt.Constants[constStmt.Name] = constStmt
}

func (stmt *ClassDeclarationStatement) AddMethod(method *MethodDefinitionStatement) {
	stmt.Methods[method.Name] = method
}

func (stmt *ClassDeclarationStatement) AddProperty(property *PropertyDeclarationStatement) {
	if !slices.Contains(stmt.PropertieNames, property.Name) {
		stmt.PropertieNames = append(stmt.PropertieNames, property.Name)
	}
	stmt.Properties[property.Name] = property
}

func (stmt *ClassDeclarationStatement) AddTrait(trait *TraitUseStatement) {
	stmt.Traits = append(stmt.Traits, trait)
}

// -------------------------------------- ThrowStatement -------------------------------------- MARK: ThrowStatement

type ThrowStatement struct {
	*Statement
	Expr IExpression
}

func NewThrowStmt(id int64, pos *position.Position, expr IExpression) *ThrowStatement {
	return &ThrowStatement{Statement: NewStmt(id, ThrowStmt, pos), Expr: expr}
}

func (stmt *ThrowStatement) Process(visitor Visitor, context any) (any, error) {
	return visitor.ProcessThrowStmt(stmt, context)
}

// -------------------------------------- DeclareStatement -------------------------------------- MARK: DeclareStatement

type DeclareStatement struct {
	*Statement
	Directive string
	Literal   IExpression
}

func NewDeclareStmt(id int64, pos *position.Position, directive string, literal IExpression) *DeclareStatement {
	return &DeclareStatement{Statement: NewStmt(id, DeclareStmt, pos), Directive: directive, Literal: literal}
}

func (stmt *DeclareStatement) Process(visitor Visitor, context any) (any, error) {
	return visitor.ProcessDeclareStmt(stmt, context)
}

// -------------------------------------- ForeachStatement -------------------------------------- MARK: ForeachStatement

type ForeachStatement struct {
	*Statement
	Collection IExpression
	Key        IExpression
	Value      IExpression
	Block      IStatement
}

func NewForeachStmt(id int64, pos *position.Position, collection, key, value IExpression, block IStatement) *ForeachStatement {
	return &ForeachStatement{Statement: NewStmt(id, ForeachStmt, pos), Collection: collection, Key: key, Value: value, Block: block}
}

func (stmt *ForeachStatement) Process(visitor Visitor, context any) (any, error) {
	return visitor.ProcessForeachStmt(stmt, context)
}
