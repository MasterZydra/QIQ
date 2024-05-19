package lexer

import "fmt"

type Token struct {
	TokenType TokenType
	Value     string
	Ln        int
	Col       int
}

func NewToken(tokenType TokenType, value string) *Token {
	return &Token{TokenType: tokenType, Value: value}
}

func (token *Token) String() string {
	return fmt.Sprintf("&{Token %s '%s'}", token.TokenType, token.Value)
}

type TokenType string

const (
	EndOfFileToken TokenType = "EOF"
	// Spec: https://phplang.org/spec/04-basic-concepts.html#grammar-start-tag
	TextToken     TokenType = "Text"
	StartTagToken TokenType = "StartTag"
	EndTagToken   TokenType = "EndTag"
	// Spec: https://phplang.org/spec/09-lexical-structure.html#general-1
	VariableNameToken         TokenType = "VariableName"
	NameToken                 TokenType = "Name"
	KeywordToken              TokenType = "Keyword"
	IntegerLiteralToken       TokenType = "IntegerLiteral"
	FloatingLiteralToken      TokenType = "FloatingLiteral"
	StringLiteralToken        TokenType = "StringLiteral"
	OperatorOrPunctuatorToken TokenType = "OperatorOrPunctuator"
)
