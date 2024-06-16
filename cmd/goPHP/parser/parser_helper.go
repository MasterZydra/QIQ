package parser

import (
	"GoPHP/cmd/goPHP/lexer"
	"GoPHP/cmd/goPHP/phpError"
)

func (parser *Parser) isEof() bool {
	return parser.currPos > len(parser.tokens)-1
}

func (parser *Parser) at() *lexer.Token {
	if parser.isEof() {
		return lexer.NewToken(lexer.EndOfFileToken, "", nil)
	}

	return parser.tokens[parser.currPos]
}

func (parser *Parser) next(offset int) *lexer.Token {
	pos := parser.currPos + offset + 1
	if pos > len(parser.tokens) {
		pos = len(parser.tokens)
	}

	return parser.tokens[pos]
}

func (parser *Parser) eat() *lexer.Token {
	if parser.isEof() {
		return lexer.NewToken(lexer.EndOfFileToken, "", nil)
	}

	result := parser.at()
	parser.currPos++
	return result
}

func (parser *Parser) eatN(length int) *lexer.Token {
	var result *lexer.Token
	for i := 0; i < length; i++ {
		result = parser.eat()
	}
	return result
}

func (parser *Parser) isTokenType(tokenType lexer.TokenType, eat bool) bool {
	result := parser.at().TokenType == tokenType
	if result && eat {
		parser.eat()
	}
	return result
}

func (parser *Parser) isToken(tokenType lexer.TokenType, value string, eat bool) bool {
	result := parser.at().TokenType == tokenType && parser.at().Value == value
	if result && eat {
		parser.eat()
	}
	return result
}

func (parser *Parser) expectTokenType(tokenType lexer.TokenType, eat bool) phpError.Error {
	if parser.isTokenType(tokenType, eat) {
		return nil
	}
	return phpError.NewParseError("Unexpected token %s. Expected: %s", parser.at().TokenType, tokenType)
}

func (parser *Parser) expect(tokenType lexer.TokenType, value string, eat bool) phpError.Error {
	if parser.isToken(tokenType, value, eat) {
		return nil
	}
	return phpError.NewParseError("Unexpected token %s. Expected: %s", parser.at(), lexer.NewToken(tokenType, value, nil))
}
