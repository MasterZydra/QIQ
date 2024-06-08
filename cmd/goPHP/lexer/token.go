package lexer

import (
	"GoPHP/cmd/goPHP/position"
	"fmt"
)

type Token struct {
	TokenType TokenType
	Value     string
	Position  *position.Position
}

func NewToken(tokenType TokenType, value string, position *position.Position) *Token {
	return &Token{TokenType: tokenType, Value: value, Position: position}
}

func (token *Token) String() string {
	return fmt.Sprintf("&{Token - type: %s, value: \"%s\", position: %s}", token.TokenType, token.Value, token.Position)
}

type TokenType string

const (
	EndOfFileToken TokenType = "EOF"
	// Spec: https://phplang.org/spec/04-basic-concepts.html#grammar-start-tag
	TextToken     TokenType = "Text"
	StartTagToken TokenType = "StartTag"
	EndTagToken   TokenType = "EndTag"
	// Spec: https://phplang.org/spec/09-lexical-structure.html#general-1
	VariableNameToken    TokenType = "VariableName"
	NameToken            TokenType = "Name"
	KeywordToken         TokenType = "Keyword"
	IntegerLiteralToken  TokenType = "IntegerLiteral"
	FloatingLiteralToken TokenType = "FloatingLiteral"
	StringLiteralToken   TokenType = "StringLiteral"
	OpOrPuncToken        TokenType = "OperatorOrPunctuator"
)
