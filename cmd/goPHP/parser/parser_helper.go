package parser

import "GoPHP/cmd/goPHP/lexer"

func (parser *Parser) isEof() bool {
	return parser.currPos > len(parser.tokens)-1
}

func (parser *Parser) at() *lexer.Token {
	if parser.isEof() {
		return nil
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
		return nil
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
