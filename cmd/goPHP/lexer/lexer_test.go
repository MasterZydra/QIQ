package lexer

import (
	"testing"
)

func compareTokens(t1 *Token, t2 *Token) bool {
	return t1.Value == t2.Value && t1.TokenType == t2.TokenType
}

func TestText(t *testing.T) {
	tokens, _ := NewLexer().Tokenize("Hello world")
	if expected := NewToken(TextToken, "Hello world"); !compareTokens(expected, tokens[0]) {
		t.Errorf("Expected: \"%s\", Got \"%s\"", expected, tokens[0])
	}
}

func TestStartTag(t *testing.T) {
	tokens, _ := NewLexer().Tokenize("<?php")
	if expected := NewToken(StartTagToken, ""); !compareTokens(expected, tokens[0]) {
		t.Errorf("Expected: \"%s\", Got \"%s\"", expected, tokens[0])
	}

	tokens, _ = NewLexer().Tokenize("<?=")
	if expected := NewToken(StartTagToken, ""); !compareTokens(expected, tokens[0]) {
		t.Errorf("Expected: \"%s\", Got \"%s\"", expected, tokens[0])
	}
	if expected := NewToken(KeywordToken, "echo"); !compareTokens(expected, tokens[1]) {
		t.Errorf("Expected: \"%s\", Got \"%s\"", expected, tokens[1])
	}
}

func TestEndTag(t *testing.T) {
	tokens, _ := NewLexer().Tokenize("<?php ?>")
	if expected := NewToken(StartTagToken, ""); !compareTokens(expected, tokens[0]) {
		t.Errorf("Expected: \"%s\", Got \"%s\"", expected, tokens[0])
	}
	if expected := NewToken(OperatorOrPunctuatorToken, ";"); !compareTokens(expected, tokens[1]) {
		t.Errorf("Expected: \"%s\", Got \"%s\"", expected, tokens[1])
	}
	if expected := NewToken(EndTagToken, ""); !compareTokens(expected, tokens[2]) {
		t.Errorf("Expected: \"%s\", Got \"%s\"", expected, tokens[2])
	}

	tokens, _ = NewLexer().Tokenize("<?= ?>")
	if expected := NewToken(StartTagToken, ""); !compareTokens(expected, tokens[0]) {
		t.Errorf("Expected: \"%s\", Got \"%s\"", expected, tokens[0])
	}
	if expected := NewToken(OperatorOrPunctuatorToken, ";"); !compareTokens(expected, tokens[2]) {
		t.Errorf("Expected: \"%s\", Got \"%s\"", expected, tokens[2])
	}
	if expected := NewToken(EndTagToken, ""); !compareTokens(expected, tokens[3]) {
		t.Errorf("Expected: \"%s\", Got \"%s\"", expected, tokens[3])
	}

	tokens, _ = NewLexer().Tokenize("<?php ?> ?>")
	if expected := NewToken(StartTagToken, ""); !compareTokens(expected, tokens[0]) {
		t.Errorf("Expected: \"%s\", Got \"%s\"", expected, tokens[0])
	}
	if expected := NewToken(OperatorOrPunctuatorToken, ";"); !compareTokens(expected, tokens[1]) {
		t.Errorf("Expected: \"%s\", Got \"%s\"", expected, tokens[1])
	}
	if expected := NewToken(EndTagToken, ""); !compareTokens(expected, tokens[2]) {
		t.Errorf("Expected: \"%s\", Got \"%s\"", expected, tokens[2])
	}
	if expected := NewToken(TextToken, " ?>"); !compareTokens(expected, tokens[3]) {
		t.Errorf("Expected: \"%s\", Got \"%s\"", expected, tokens[3])
	}
}

func TestIntegerLiteral(t *testing.T) {
	// binary-literal

	tokens, _ := NewLexer().Tokenize("<?php 0b1010")
	if expected := NewToken(IntegerLiteralToken, "0b1010"); !compareTokens(expected, tokens[1]) {
		t.Errorf("Expected: \"%s\", Got \"%s\"", expected, tokens[1])
	}

	tokens, _ = NewLexer().Tokenize("<?php 0B1010")
	if expected := NewToken(IntegerLiteralToken, "0B1010"); !compareTokens(expected, tokens[1]) {
		t.Errorf("Expected: \"%s\", Got \"%s\"", expected, tokens[1])
	}

	// hexadecimal-literal

	tokens, _ = NewLexer().Tokenize("<?php 0x0123456789AbCdEf")
	if expected := NewToken(IntegerLiteralToken, "0x0123456789AbCdEf"); !compareTokens(expected, tokens[1]) {
		t.Errorf("Expected: \"%s\", Got \"%s\"", expected, tokens[1])
	}

	tokens, _ = NewLexer().Tokenize("<?php 0X0123456789AbCdEf")
	if expected := NewToken(IntegerLiteralToken, "0X0123456789AbCdEf"); !compareTokens(expected, tokens[1]) {
		t.Errorf("Expected: \"%s\", Got \"%s\"", expected, tokens[1])
	}

	// decimal-literal

	tokens, _ = NewLexer().Tokenize("<?php 124")
	if expected := NewToken(IntegerLiteralToken, "124"); !compareTokens(expected, tokens[1]) {
		t.Errorf("Expected: \"%s\", Got \"%s\"", expected, tokens[1])
	}

	// octal-literal

	tokens, _ = NewLexer().Tokenize("<?php 047")
	if expected := NewToken(IntegerLiteralToken, "047"); !compareTokens(expected, tokens[1]) {
		t.Errorf("Expected: \"%s\", Got \"%s\"", expected, tokens[1])
	}
}

func TestFloatingLiteral(t *testing.T) {
	tokens, _ := NewLexer().Tokenize("<?php .5")
	if expected := NewToken(FloatingLiteralToken, ".5"); !compareTokens(expected, tokens[1]) {
		t.Errorf("Expected: \"%s\", Got \"%s\"", expected, tokens[1])
	}

	tokens, _ = NewLexer().Tokenize("<?php 1.2")
	if expected := NewToken(FloatingLiteralToken, "1.2"); !compareTokens(expected, tokens[1]) {
		t.Errorf("Expected: \"%s\", Got \"%s\"", expected, tokens[1])
	}

	tokens, _ = NewLexer().Tokenize("<?php .5e-4")
	if expected := NewToken(FloatingLiteralToken, ".5e-4"); !compareTokens(expected, tokens[1]) {
		t.Errorf("Expected: \"%s\", Got \"%s\"", expected, tokens[1])
	}

	tokens, _ = NewLexer().Tokenize("<?php 2.5E+4")
	if expected := NewToken(FloatingLiteralToken, "2.5E+4"); !compareTokens(expected, tokens[1]) {
		t.Errorf("Expected: \"%s\", Got \"%s\"", expected, tokens[1])
	}

	tokens, _ = NewLexer().Tokenize("<?php 2e4")
	if expected := NewToken(FloatingLiteralToken, "2e4"); !compareTokens(expected, tokens[1]) {
		t.Errorf("Expected: \"%s\", Got \"%s\"", expected, tokens[1])
	}
}

func TestStringLiteralSingleQuote(t *testing.T) {
	tokens, _ := NewLexer().Tokenize("<?php b''")
	if expected := NewToken(StringLiteralToken, "b''"); !compareTokens(expected, tokens[1]) {
		t.Errorf("Expected: \"%s\", Got \"%s\"", expected, tokens[1])
	}

	tokens, _ = NewLexer().Tokenize("<?php B''")
	if expected := NewToken(StringLiteralToken, "B''"); !compareTokens(expected, tokens[1]) {
		t.Errorf("Expected: \"%s\", Got \"%s\"", expected, tokens[1])
	}

	tokens, _ = NewLexer().Tokenize("<?php ''")
	if expected := NewToken(StringLiteralToken, "''"); !compareTokens(expected, tokens[1]) {
		t.Errorf("Expected: \"%s\", Got \"%s\"", expected, tokens[1])
	}

	tokens, _ = NewLexer().Tokenize("<?php 'abc'")
	if expected := NewToken(StringLiteralToken, "'abc'"); !compareTokens(expected, tokens[1]) {
		t.Errorf("Expected: \"%s\", Got \"%s\"", expected, tokens[1])
	}

	tokens, _ = NewLexer().Tokenize(`<?php '\'abc\\'`)
	if expected := NewToken(StringLiteralToken, `'\'abc\\'`); !compareTokens(expected, tokens[1]) {
		t.Errorf("Expected: \"%s\", Got \"%s\"", expected, tokens[1])
	}
}

func TestStringLiteralDoubleQuote(t *testing.T) {
	tokens, _ := NewLexer().Tokenize(`<?php b""`)
	if expected := NewToken(StringLiteralToken, `b""`); !compareTokens(expected, tokens[1]) {
		t.Errorf("Expected: \"%s\", Got \"%s\"", expected, tokens[1])
	}

	tokens, _ = NewLexer().Tokenize(`<?php B""`)
	if expected := NewToken(StringLiteralToken, `B""`); !compareTokens(expected, tokens[1]) {
		t.Errorf("Expected: \"%s\", Got \"%s\"", expected, tokens[1])
	}

	tokens, _ = NewLexer().Tokenize(`<?php ""`)
	if expected := NewToken(StringLiteralToken, `""`); !compareTokens(expected, tokens[1]) {
		t.Errorf("Expected: \"%s\", Got \"%s\"", expected, tokens[1])
	}

	tokens, _ = NewLexer().Tokenize(`<?php "abc"`)
	if expected := NewToken(StringLiteralToken, `"abc"`); !compareTokens(expected, tokens[1]) {
		t.Errorf("Expected: \"%s\", Got \"%s\"", expected, tokens[1])
	}

	tokens, _ = NewLexer().Tokenize(`<?php "\"abc\\\n\$"`)
	if expected := NewToken(StringLiteralToken, `"\"abc\\\n\$"`); !compareTokens(expected, tokens[1]) {
		t.Errorf("Expected: \"%s\", Got \"%s\"", expected, tokens[1])
	}
}

func TestOperatorOrPunctuator(t *testing.T) {
	tokens, _ := NewLexer().Tokenize("<?php a === 1")
	if expected := NewToken(OperatorOrPunctuatorToken, "==="); !compareTokens(expected, tokens[2]) {
		t.Errorf("Expected: \"%s\", Got \"%s\"", expected, tokens[2])
	}

	tokens, _ = NewLexer().Tokenize("<?php a == 1")
	if expected := NewToken(OperatorOrPunctuatorToken, "=="); !compareTokens(expected, tokens[2]) {
		t.Errorf("Expected: \"%s\", Got \"%s\"", expected, tokens[2])
	}

	tokens, _ = NewLexer().Tokenize("<?php a = 1")
	if expected := NewToken(OperatorOrPunctuatorToken, "="); !compareTokens(expected, tokens[2]) {
		t.Errorf("Expected: \"%s\", Got \"%s\"", expected, tokens[2])
	}
}

func TestVariableVarname(t *testing.T) {
	tokens, _ := NewLexer().Tokenize(`<?php $$var = "someValue";`)
	tokenIndex := 0
	if expected := NewToken(StartTagToken, ""); !compareTokens(expected, tokens[tokenIndex]) {
		t.Errorf("Expected: \"%s\", Got \"%s\"", expected, tokens[tokenIndex])
	}
	tokenIndex += 1
	if expected := NewToken(OperatorOrPunctuatorToken, "$"); !compareTokens(expected, tokens[tokenIndex]) {
		t.Errorf("Expected: \"%s\", Got \"%s\"", expected, tokens[tokenIndex])
	}
	tokenIndex += 1
	if expected := NewToken(VariableNameToken, "$var"); !compareTokens(expected, tokens[tokenIndex]) {
		t.Errorf("Expected: \"%s\", Got \"%s\"", expected, tokens[tokenIndex])
	}
	tokenIndex += 1
	if expected := NewToken(OperatorOrPunctuatorToken, "="); !compareTokens(expected, tokens[tokenIndex]) {
		t.Errorf("Expected: \"%s\", Got \"%s\"", expected, tokens[tokenIndex])
	}
	tokenIndex += 1
	if expected := NewToken(StringLiteralToken, `"someValue"`); !compareTokens(expected, tokens[tokenIndex]) {
		t.Errorf("Expected: \"%s\", Got \"%s\"", expected, tokens[tokenIndex])
	}
	tokenIndex += 1
	if expected := NewToken(OperatorOrPunctuatorToken, ";"); !compareTokens(expected, tokens[tokenIndex]) {
		t.Errorf("Expected: \"%s\", Got \"%s\"", expected, tokens[tokenIndex])
	}

	tokens, _ = NewLexer().Tokenize(`<?php echo 12, $var;`)
	tokenIndex = 0
	if expected := NewToken(StartTagToken, ""); !compareTokens(expected, tokens[tokenIndex]) {
		t.Errorf("Expected: \"%s\", Got \"%s\"", expected, tokens[tokenIndex])
	}
	tokenIndex += 1
	if expected := NewToken(KeywordToken, "echo"); !compareTokens(expected, tokens[tokenIndex]) {
		t.Errorf("Expected: \"%s\", Got \"%s\"", expected, tokens[tokenIndex])
	}
	tokenIndex += 1
	if expected := NewToken(IntegerLiteralToken, "12"); !compareTokens(expected, tokens[tokenIndex]) {
		t.Errorf("Expected: \"%s\", Got \"%s\"", expected, tokens[tokenIndex])
	}
	tokenIndex += 1
	if expected := NewToken(OperatorOrPunctuatorToken, ","); !compareTokens(expected, tokens[tokenIndex]) {
		t.Errorf("Expected: \"%s\", Got \"%s\"", expected, tokens[tokenIndex])
	}
	tokenIndex += 1
	if expected := NewToken(VariableNameToken, "$var"); !compareTokens(expected, tokens[tokenIndex]) {
		t.Errorf("Expected: \"%s\", Got \"%s\"", expected, tokens[tokenIndex])
	}
}

func TestHtmlAndPhp(t *testing.T) {
	tokenIndex := 0
	tokens, _ := NewLexer().Tokenize(`<body><?php $heading = "My Heading"; ?><h1><?= $heading ?></h1>`)
	if expected := NewToken(TextToken, "<body>"); !compareTokens(expected, tokens[tokenIndex]) {
		t.Errorf("Expected: \"%s\", Got \"%s\"", expected, tokens[tokenIndex])
	}
	tokenIndex += 1
	if expected := NewToken(StartTagToken, ""); !compareTokens(expected, tokens[tokenIndex]) {
		t.Errorf("Expected: \"%s\", Got \"%s\"", expected, tokens[tokenIndex])
	}
	tokenIndex += 1
	if expected := NewToken(VariableNameToken, "$heading"); !compareTokens(expected, tokens[tokenIndex]) {
		t.Errorf("Expected: \"%s\", Got \"%s\"", expected, tokens[tokenIndex])
	}
	tokenIndex += 1
	if expected := NewToken(OperatorOrPunctuatorToken, "="); !compareTokens(expected, tokens[tokenIndex]) {
		t.Errorf("Expected: \"%s\", Got \"%s\"", expected, tokens[tokenIndex])
	}
	tokenIndex += 1
	if expected := NewToken(StringLiteralToken, `"My Heading"`); !compareTokens(expected, tokens[tokenIndex]) {
		t.Errorf("Expected: \"%s\", Got \"%s\"", expected, tokens[tokenIndex])
	}
	tokenIndex += 1
	if expected := NewToken(OperatorOrPunctuatorToken, ";"); !compareTokens(expected, tokens[tokenIndex]) {
		t.Errorf("Expected: \"%s\", Got \"%s\"", expected, tokens[tokenIndex])
	}
	tokenIndex += 1
	if expected := NewToken(EndTagToken, ""); !compareTokens(expected, tokens[tokenIndex]) {
		t.Errorf("Expected: \"%s\", Got \"%s\"", expected, tokens[tokenIndex])
	}
	tokenIndex += 1
	if expected := NewToken(TextToken, "<h1>"); !compareTokens(expected, tokens[tokenIndex]) {
		t.Errorf("Expected: \"%s\", Got \"%s\"", expected, tokens[tokenIndex])
	}
	tokenIndex += 1
	if expected := NewToken(StartTagToken, ""); !compareTokens(expected, tokens[tokenIndex]) {
		t.Errorf("Expected: \"%s\", Got \"%s\"", expected, tokens[tokenIndex])
	}
	tokenIndex += 1
	if expected := NewToken(KeywordToken, "echo"); !compareTokens(expected, tokens[tokenIndex]) {
		t.Errorf("Expected: \"%s\", Got \"%s\"", expected, tokens[tokenIndex])
	}
	tokenIndex += 1
	if expected := NewToken(VariableNameToken, "$heading"); !compareTokens(expected, tokens[tokenIndex]) {
		t.Errorf("Expected: \"%s\", Got \"%s\"", expected, tokens[tokenIndex])
	}
	tokenIndex += 1
	if expected := NewToken(OperatorOrPunctuatorToken, ";"); !compareTokens(expected, tokens[tokenIndex]) {
		t.Errorf("Expected: \"%s\", Got \"%s\"", expected, tokens[tokenIndex])
	}
	tokenIndex += 1
	if expected := NewToken(EndTagToken, ""); !compareTokens(expected, tokens[tokenIndex]) {
		t.Errorf("Expected: \"%s\", Got \"%s\"", expected, tokens[tokenIndex])
	}
	tokenIndex += 1
	if expected := NewToken(TextToken, "</h1>"); !compareTokens(expected, tokens[tokenIndex]) {
		t.Errorf("Expected: \"%s\", Got \"%s\"", expected, tokens[tokenIndex])
	}
}
