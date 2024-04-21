package lexer

import (
	"GoPHP/cmd/goPHP/common"
	"strings"
)

func (lexer *Lexer) isEof() bool {
	return lexer.currPos > len(lexer.input)-1
}

func (lexer *Lexer) at() string {
	if lexer.isEof() {
		return ""
	}

	return lexer.input[lexer.currPos : lexer.currPos+1]
}

func (lexer *Lexer) next(offset int) string {
	pos := lexer.currPos + offset + 1
	if pos > len(lexer.input) {
		pos = len(lexer.input)
	}

	return lexer.input[pos : pos+1]
}

func (lexer *Lexer) nextN(length int) string {
	end := lexer.currPos + length
	if end > len(lexer.input) {
		end = len(lexer.input)
	}

	return lexer.input[lexer.currPos:end]
}

func (lexer *Lexer) eat() string {
	if lexer.isEof() {
		return ""
	}

	result := lexer.at()
	lexer.currPos++
	return result
}

func (lexer *Lexer) eatN(length int) string {
	result := ""
	for i := 0; i < length; i++ {
		result += lexer.eat()
	}
	return result
}

func (lexer *Lexer) pushToken(tokenType TokenType, value string) {
	lexer.tokens = append(lexer.tokens, NewToken(tokenType, value))
}

func (lexer *Lexer) pushKeywordToken(keyword string) {
	if common.IsKeyword(keyword) {
		lexer.pushToken(KeywordToken, strings.ToLower(keyword))
		return
	}
	if common.IsCorePredefinedConstants(keyword) || common.IsContextDependentConstants(keyword) {
		lexer.pushToken(KeywordToken, strings.ToUpper(keyword))
		return
	}
	lexer.pushToken(KeywordToken, keyword)
}

func (lexer *Lexer) lastToken() *Token {
	return lexer.tokens[len(lexer.tokens)-1]
}
