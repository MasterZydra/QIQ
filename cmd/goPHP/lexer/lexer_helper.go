package lexer

import (
	"GoPHP/cmd/goPHP/common"
	"GoPHP/cmd/goPHP/position"
	"strings"
)

func (lexer *Lexer) isEof() bool {
	return lexer.currPos.CurrPos > len(lexer.input)-1
}

func (lexer *Lexer) at() string {
	if lexer.isEof() {
		return ""
	}

	return lexer.input[lexer.currPos.CurrPos : lexer.currPos.CurrPos+1]
}

func (lexer *Lexer) next(offset int) string {
	pos := lexer.currPos.CurrPos + offset + 1
	if pos > len(lexer.input) {
		pos = len(lexer.input)
	}
	if pos+1 > len(lexer.input) {
		return ""
	}
	return lexer.input[pos : pos+1]
}

func (lexer *Lexer) nextN(length int) string {
	end := lexer.currPos.CurrPos + length
	if end > len(lexer.input) {
		end = len(lexer.input)
	}

	return lexer.input[lexer.currPos.CurrPos:end]
}

func (lexer *Lexer) eat() string {
	if lexer.isEof() {
		return ""
	}

	result := lexer.at()

	if lexer.currPos.SearchTokenStart && !lexer.isNewLineChar(result) && !lexer.isWhiteSpaceChar(result) {
		lexer.currPos.CurrTokenLine = lexer.currPos.CurrLine
		lexer.currPos.CurrTokenCol = lexer.currPos.CurrCol
		lexer.currPos.SearchTokenStart = false
	}

	if result == "\n" {
		lexer.currPos.CurrLine += 1
		lexer.currPos.CurrCol = 1
	} else {
		lexer.currPos.CurrCol += 1
	}

	lexer.currPos.CurrPos++
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
	var pos *position.Position = nil
	if tokenType != TextToken {
		pos = position.NewPosition(lexer.filename, lexer.currPos.CurrTokenLine, lexer.currPos.CurrTokenCol)
	}
	lexer.tokens = append(lexer.tokens, NewToken(tokenType, value, pos))
	lexer.currPos.SearchTokenStart = true
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

func (lexer *Lexer) pushSnapShot() {
	lexer.positionSnapShots = append(lexer.positionSnapShots, lexer.currPos)
}

func (lexer *Lexer) popSnapShot(apply bool) {
	if len(lexer.positionSnapShots) == 0 {
		return
	}
	snapShot := lexer.positionSnapShots[len(lexer.positionSnapShots)-1]
	if len(lexer.positionSnapShots) == 1 {
		lexer.positionSnapShots = []PositionSnapshot{}
	} else {
		lexer.positionSnapShots = lexer.positionSnapShots[:len(lexer.positionSnapShots)-1]
	}

	if apply {
		lexer.currPos = snapShot
	}
}
