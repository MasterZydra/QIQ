package ast

type Program struct {
	statements []IStatement
}

func NewProgram() *Program {
	return &Program{}
}

func (program *Program) Append(stmt IStatement) {
	program.statements = append(program.statements, stmt)
}

func (program *Program) GetStatements() []IStatement {
	return program.statements
}
