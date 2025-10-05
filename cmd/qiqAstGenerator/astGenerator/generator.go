package astGenerator

import (
	"QIQ/cmd/qiq/ast"
	"QIQ/cmd/qiq/ini"
	"QIQ/cmd/qiq/parser"
	"QIQ/cmd/qiq/phpError"
)

type AstGenerator struct {
	output string
}

func NewAstGenerator() *AstGenerator { return &AstGenerator{} }

func (generator *AstGenerator) Process(sourceCode, filename string) (string, phpError.Error) {
	program, err := parser.NewParser(ini.NewDevIni()).ProduceAST(sourceCode, filename)
	if err != nil {
		return "", err
	}

	err = generator.processProgram(program)
	if err != nil {
		return "", err
	}

	return generator.output, nil
}

func (generator *AstGenerator) processProgram(program *ast.Program) phpError.Error {
	for _, stmt := range program.GetStatements() {
		if err := generator.processStmt(stmt); err != nil {
			return err
		}
	}
	return nil
}

func (generator *AstGenerator) processStmt(stmt ast.IStatement) phpError.Error {
	if stmt == nil {
		generator.print("nil")
		return nil
	}
	_, err := stmt.Process(generator, nil)
	if err != nil {
		return err.(phpError.Error)
	}
	return nil
}
