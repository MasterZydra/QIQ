package lexer

import (
	"testing"
)

func testTokenize(t *testing.T, php string, expected []*Token) {
	compareTokens := func(t1 *Token, t2 *Token) bool {
		return t1.Value == t2.Value && t1.TokenType == t2.TokenType
	}

	tokens, err := NewLexer().Tokenize(php)
	if err != nil {
		t.Errorf("Unexpected error: \"%s\"", err)
		return
	}
	for index, token := range expected {
		if !compareTokens(token, tokens[index]) {
			t.Errorf("Expected: \"%s\", Got \"%s\"", token, tokens[index])
			return
		}
	}
}

func TestText(t *testing.T) {
	testTokenize(t, "Hello world", []*Token{NewToken(TextToken, "Hello world")})
}

func TestStartTag(t *testing.T) {
	testTokenize(t, "<?php", []*Token{NewToken(StartTagToken, "")})
	testTokenize(t, "<?=", []*Token{NewToken(StartTagToken, ""), NewToken(KeywordToken, "echo")})
}

func TestEndTag(t *testing.T) {
	testTokenize(t,
		"<?php ?>",
		[]*Token{NewToken(StartTagToken, ""), NewToken(OperatorOrPunctuatorToken, ";"), NewToken(EndTagToken, "")},
	)

	testTokenize(t,
		"<?= ?>",
		[]*Token{NewToken(StartTagToken, ""),
			NewToken(KeywordToken, "echo"),
			NewToken(OperatorOrPunctuatorToken, ";"),
			NewToken(EndTagToken, "")},
	)

	testTokenize(t,
		"<?php ?> ?>",
		[]*Token{NewToken(StartTagToken, ""),
			NewToken(OperatorOrPunctuatorToken, ";"),
			NewToken(EndTagToken, ""),
			NewToken(TextToken, " ?>")},
	)
}

func TestIntegerLiteral(t *testing.T) {
	// binary-literal
	testTokenize(t, "<?php 0b1010", []*Token{NewToken(StartTagToken, ""), NewToken(IntegerLiteralToken, "0b1010")})
	testTokenize(t, "<?php 0B1010", []*Token{NewToken(StartTagToken, ""), NewToken(IntegerLiteralToken, "0B1010")})

	// hexadecimal-literal
	testTokenize(t, "<?php 0x0123456789AbCdEf",
		[]*Token{NewToken(StartTagToken, ""), NewToken(IntegerLiteralToken, "0x0123456789AbCdEf")})
	testTokenize(t, "<?php 0X0123456789AbCdEf",
		[]*Token{NewToken(StartTagToken, ""), NewToken(IntegerLiteralToken, "0X0123456789AbCdEf")})

	// decimal-literal
	testTokenize(t, "<?php 124", []*Token{NewToken(StartTagToken, ""), NewToken(IntegerLiteralToken, "124")})

	// octal-literal
	testTokenize(t, "<?php 047", []*Token{NewToken(StartTagToken, ""), NewToken(IntegerLiteralToken, "047")})
}

func TestFloatingLiteral(t *testing.T) {
	testTokenize(t, "<?php .5", []*Token{NewToken(StartTagToken, ""), NewToken(FloatingLiteralToken, ".5")})
	testTokenize(t, "<?php 1.2", []*Token{NewToken(StartTagToken, ""), NewToken(FloatingLiteralToken, "1.2")})
	testTokenize(t, "<?php .5e-4", []*Token{NewToken(StartTagToken, ""), NewToken(FloatingLiteralToken, ".5e-4")})
	testTokenize(t, "<?php 2.5e-4", []*Token{NewToken(StartTagToken, ""), NewToken(FloatingLiteralToken, "2.5e-4")})
	testTokenize(t, "<?php 2e4", []*Token{NewToken(StartTagToken, ""), NewToken(FloatingLiteralToken, "2e4")})
}

func TestStringLiteral(t *testing.T) {
	// Single quote
	testTokenize(t, "<?php b''", []*Token{NewToken(StartTagToken, ""), NewToken(StringLiteralToken, "b''")})
	testTokenize(t, "<?php B''", []*Token{NewToken(StartTagToken, ""), NewToken(StringLiteralToken, "B''")})
	testTokenize(t, "<?php ''", []*Token{NewToken(StartTagToken, ""), NewToken(StringLiteralToken, "''")})
	testTokenize(t, "<?php 'abc'", []*Token{NewToken(StartTagToken, ""), NewToken(StringLiteralToken, "'abc'")})
	testTokenize(t, `<?php '\'abc\\'`, []*Token{NewToken(StartTagToken, ""), NewToken(StringLiteralToken, `'\'abc\\'`)})

	// Double quote
	testTokenize(t, `<?php b""`, []*Token{NewToken(StartTagToken, ""), NewToken(StringLiteralToken, `b""`)})
	testTokenize(t, `<?php B""`, []*Token{NewToken(StartTagToken, ""), NewToken(StringLiteralToken, `B""`)})
	testTokenize(t, `<?php ""`, []*Token{NewToken(StartTagToken, ""), NewToken(StringLiteralToken, `""`)})
	testTokenize(t, `<?php "abc"`, []*Token{NewToken(StartTagToken, ""), NewToken(StringLiteralToken, `"abc"`)})
	testTokenize(t, `<?php "\"abc\\\n\$"`, []*Token{NewToken(StartTagToken, ""), NewToken(StringLiteralToken, `"\"abc\\\n\$"`)})
}

func TestOperatorOrPunctuator(t *testing.T) {
	testTokenize(t, "<?php a === 1",
		[]*Token{NewToken(StartTagToken, ""), NewToken(NameToken, "a"), NewToken(OperatorOrPunctuatorToken, "===")})

	testTokenize(t, "<?php a == 1",
		[]*Token{NewToken(StartTagToken, ""), NewToken(NameToken, "a"), NewToken(OperatorOrPunctuatorToken, "==")})

	testTokenize(t, "<?php a = 1",
		[]*Token{NewToken(StartTagToken, ""), NewToken(NameToken, "a"), NewToken(OperatorOrPunctuatorToken, "=")})

	testTokenize(t, "<?php a ** 1",
		[]*Token{NewToken(StartTagToken, ""), NewToken(NameToken, "a"), NewToken(OperatorOrPunctuatorToken, "**")})
}

func TestVariableVarname(t *testing.T) {
	testTokenize(t, `<?php $$var = "someValue";`,
		[]*Token{
			NewToken(StartTagToken, ""), NewToken(OperatorOrPunctuatorToken, "$"), NewToken(VariableNameToken, "$var"),
			NewToken(OperatorOrPunctuatorToken, "="), NewToken(StringLiteralToken, `"someValue"`),
			NewToken(OperatorOrPunctuatorToken, ";"),
		})

	testTokenize(t, `<?php echo 12, $var;`,
		[]*Token{
			NewToken(StartTagToken, ""), NewToken(KeywordToken, "echo"), NewToken(IntegerLiteralToken, "12"),
			NewToken(OperatorOrPunctuatorToken, ","), NewToken(VariableNameToken, "$var"),
		})
}

func TestHtmlAndPhp(t *testing.T) {
	testTokenize(t, `<body><?php $heading = "My Heading"; ?><h1><?= $heading ?></h1>`,
		[]*Token{
			NewToken(TextToken, "<body>"), NewToken(StartTagToken, ""), NewToken(VariableNameToken, "$heading"),
			NewToken(OperatorOrPunctuatorToken, "="), NewToken(StringLiteralToken, `"My Heading"`),
			NewToken(OperatorOrPunctuatorToken, ";"), NewToken(EndTagToken, ""), NewToken(TextToken, "<h1>"),
			NewToken(StartTagToken, ""), NewToken(KeywordToken, "echo"), NewToken(VariableNameToken, "$heading"),
			NewToken(OperatorOrPunctuatorToken, ";"), NewToken(EndTagToken, ""), NewToken(TextToken, "</h1>"),
		})
}
