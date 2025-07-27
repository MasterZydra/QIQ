package parser

import (
	"GoPHP/cmd/goPHP/common"
	"GoPHP/cmd/goPHP/config"
	"GoPHP/cmd/goPHP/lexer"
	"GoPHP/cmd/goPHP/phpError"
	"GoPHP/cmd/goPHP/position"
	"fmt"
	"strings"
)

func PrintParserCallstack(function string, parser *Parser) {
	if config.ShowParserCallStack {
		if parser != nil {
			fmt.Printf("%s (%s)\n", function, parser.at().Position.ToPosString())
		} else {
			println(function)
		}
	}
}

// -------------------------------------- Common -------------------------------------- MARK: Common

func (parser *Parser) isEof() bool {
	return parser.currPos > len(parser.tokens)-1
}

func (parser *Parser) at() *lexer.Token {
	if parser.isEof() {
		lastPos := parser.tokens[parser.currPos-1].Position
		eofPos := position.NewPosition(lastPos.File, lastPos.Line, lastPos.Column+1)
		return lexer.NewToken(lexer.EndOfFileToken, "EOF", eofPos)
	}

	return parser.tokens[parser.currPos]
}

func (parser *Parser) next(offset int) *lexer.Token {
	pos := parser.currPos + offset + 1
	if parser.isEof() || pos >= len(parser.tokens) {
		return lexer.NewToken(lexer.EndOfFileToken, "", nil)
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

// -------------------------------------- Tokens -------------------------------------- MARK: Tokens

func (parser *Parser) isTextExpression(eat bool) bool {
	isTextExpression := (parser.isToken(lexer.OpOrPuncToken, ";", false) &&
		parser.next(0).TokenType == lexer.EndTagToken &&
		parser.next(1).TokenType == lexer.TextToken) || (parser.isTokenType(lexer.EndTagToken, false) &&
		parser.next(0).TokenType == lexer.TextToken)

	if isTextExpression && eat {
		parser.isToken(lexer.OpOrPuncToken, ";", true)
		parser.eat()
	}

	return isTextExpression
}

func (parser *Parser) isPhpType(token *lexer.Token) bool {
	return token.TokenType == lexer.OpOrPuncToken && token.Value == "?" ||
		token.TokenType == lexer.KeywordToken && common.IsReturnTypeKeyword(token.Value)
}

func (parser *Parser) getTypes(eat bool) ([]string, phpError.Error) {
	types, _, err := parser.getTypesWithOffset(eat, -1)
	return types, err
}

func (parser *Parser) getTypesWithOffset(eat bool, offset int) ([]string, int, phpError.Error) {
	types := []string{}

	token := func() *lexer.Token {
		return parser.next(offset)
	}

	if token().TokenType == lexer.OpOrPuncToken && token().Value == "?" {
		types = append(types, "null")
		offset++
	}

	if !parser.isPhpType(token()) {
		return types, offset, phpError.NewParseError("Expected a type. Got \"%s\"s at %s", token().Value, token().Position.ToPosString())
	}

	for parser.isPhpType(token()) {
		types = append(types, strings.ToLower(token().Value))
		offset++

		if token().TokenType == lexer.OpOrPuncToken && token().Value == "|" {
			offset++
			continue
		}

		break
	}

	if eat {
		parser.eatN(offset + 1)
	}

	if len(types) == 0 {
		types = append(types, "mixed")
	}

	return types, offset, nil
}
