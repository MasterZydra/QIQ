package lexer

import (
	"GoPHP/cmd/goPHP/ini"
	"GoPHP/cmd/goPHP/position"
	"fmt"
	"testing"
)

var testFile = "test.php"

func testTokenize(t *testing.T, php string, expected []*Token) {
	compareTokens := func(t1 *Token, t2 *Token) bool {
		if t1.Position != nil && t2.Position != nil {
			t1.Position.Filename = t2.Position.Filename
		}
		return t1.String() == t2.String()
	}

	tokens, err := NewLexer(ini.NewDevIni()).Tokenize(php, testFile)
	if err != nil {
		fmt.Println("    Code:", php)
		t.Errorf("Unexpected error: \"%s\"", err)
		return
	}
	for index, token := range expected {
		if !compareTokens(token, tokens[index]) {
			fmt.Println("    Code:", php)
			t.Errorf("\nExpected: \"%s\"\nGot:      \"%s\"", token, tokens[index])
			return
		}
	}
}

func TestText(t *testing.T) {
	testTokenize(t, "Hello world", []*Token{NewToken(TextToken, "Hello world", nil)})
}

func TestStartTag(t *testing.T) {
	testTokenize(t, "<?php", []*Token{NewToken(StartTagToken, "", position.NewPosition(testFile, 1, 1))})
	testTokenize(t, "<?=", []*Token{
		NewToken(StartTagToken, "", position.NewPosition(testFile, 1, 1)),
		NewToken(KeywordToken, "echo", position.NewPosition(testFile, 1, 1)),
	})
}

func TestEndTag(t *testing.T) {
	testTokenize(t,
		"<?php ?>",
		[]*Token{
			NewToken(StartTagToken, "", position.NewPosition(testFile, 1, 1)),
			NewToken(OpOrPuncToken, ";", position.NewPosition(testFile, 1, 7)),
			NewToken(EndTagToken, "", position.NewPosition(testFile, 1, 7)),
		},
	)

	testTokenize(t,
		"<?= ?>",
		[]*Token{
			NewToken(StartTagToken, "", position.NewPosition(testFile, 1, 1)),
			NewToken(KeywordToken, "echo", position.NewPosition(testFile, 1, 1)),
			NewToken(OpOrPuncToken, ";", position.NewPosition(testFile, 1, 5)),
			NewToken(EndTagToken, "", position.NewPosition(testFile, 1, 5)),
		},
	)

	testTokenize(t,
		"<?php ?> ?>",
		[]*Token{
			NewToken(StartTagToken, "", position.NewPosition(testFile, 1, 1)),
			NewToken(OpOrPuncToken, ";", position.NewPosition(testFile, 1, 7)),
			NewToken(EndTagToken, "", position.NewPosition(testFile, 1, 7)),
			NewToken(TextToken, " ?>", nil),
		},
	)
}

func TestIntegerLiteral(t *testing.T) {
	// binary-literal
	testTokenize(t, "<?php 0b1010", []*Token{
		NewToken(StartTagToken, "", position.NewPosition(testFile, 1, 1)),
		NewToken(IntegerLiteralToken, "0b1010", position.NewPosition(testFile, 1, 7)),
	})
	testTokenize(t, "<?php 0B1010", []*Token{
		NewToken(StartTagToken, "", position.NewPosition(testFile, 1, 1)),
		NewToken(IntegerLiteralToken, "0B1010", position.NewPosition(testFile, 1, 7)),
	})

	// hexadecimal-literal
	testTokenize(t, "<?php 0x0123456789AbCdEf", []*Token{
		NewToken(StartTagToken, "", position.NewPosition(testFile, 1, 1)),
		NewToken(IntegerLiteralToken, "0x0123456789AbCdEf", position.NewPosition(testFile, 1, 7)),
	})
	testTokenize(t, "<?php 0X0123456789AbCdEf", []*Token{
		NewToken(StartTagToken, "", position.NewPosition(testFile, 1, 1)),
		NewToken(IntegerLiteralToken, "0X0123456789AbCdEf", position.NewPosition(testFile, 1, 7)),
	})

	// decimal-literal
	testTokenize(t, "<?php 124", []*Token{
		NewToken(StartTagToken, "", position.NewPosition(testFile, 1, 1)),
		NewToken(IntegerLiteralToken, "124", position.NewPosition(testFile, 1, 7)),
	})

	// octal-literal
	testTokenize(t, "<?php 047", []*Token{
		NewToken(StartTagToken, "", position.NewPosition(testFile, 1, 1)),
		NewToken(IntegerLiteralToken, "047", position.NewPosition(testFile, 1, 7)),
	})
}

func TestFloatingLiteral(t *testing.T) {
	testTokenize(t, "<?php .5", []*Token{
		NewToken(StartTagToken, "", position.NewPosition(testFile, 1, 1)),
		NewToken(FloatingLiteralToken, ".5", position.NewPosition(testFile, 1, 7)),
	})
	testTokenize(t, "<?php 1.2", []*Token{
		NewToken(StartTagToken, "", position.NewPosition(testFile, 1, 1)),
		NewToken(FloatingLiteralToken, "1.2", position.NewPosition(testFile, 1, 7)),
	})
	testTokenize(t, "<?php .5e-4", []*Token{
		NewToken(StartTagToken, "", position.NewPosition(testFile, 1, 1)),
		NewToken(FloatingLiteralToken, ".5e-4", position.NewPosition(testFile, 1, 7)),
	})
	testTokenize(t, "<?php 2.5e-4", []*Token{
		NewToken(StartTagToken, "", position.NewPosition(testFile, 1, 1)),
		NewToken(FloatingLiteralToken, "2.5e-4", position.NewPosition(testFile, 1, 7)),
	})
	testTokenize(t, "<?php 2e4", []*Token{
		NewToken(StartTagToken, "", position.NewPosition(testFile, 1, 1)),
		NewToken(FloatingLiteralToken, "2e4", position.NewPosition(testFile, 1, 7)),
	})
}

func TestStringLiteral(t *testing.T) {
	// Single quote
	testTokenize(t, "<?php b''", []*Token{
		NewToken(StartTagToken, "", position.NewPosition(testFile, 1, 1)),
		NewToken(StringLiteralToken, "b''", position.NewPosition(testFile, 1, 7)),
	})
	testTokenize(t, "<?php B''", []*Token{
		NewToken(StartTagToken, "", position.NewPosition(testFile, 1, 1)),
		NewToken(StringLiteralToken, "B''", position.NewPosition(testFile, 1, 7)),
	})
	testTokenize(t, "<?php ''", []*Token{
		NewToken(StartTagToken, "", position.NewPosition(testFile, 1, 1)),
		NewToken(StringLiteralToken, "''", position.NewPosition(testFile, 1, 7)),
	})
	testTokenize(t, "<?php 'abc'", []*Token{
		NewToken(StartTagToken, "", position.NewPosition(testFile, 1, 1)),
		NewToken(StringLiteralToken, "'abc'", position.NewPosition(testFile, 1, 7)),
	})
	testTokenize(t, `<?php '\'abc\\'`, []*Token{
		NewToken(StartTagToken, "", position.NewPosition(testFile, 1, 1)),
		NewToken(StringLiteralToken, `'\'abc\\'`, position.NewPosition(testFile, 1, 7)),
	})

	// Double quote
	testTokenize(t, `<?php b""`, []*Token{
		NewToken(StartTagToken, "", position.NewPosition(testFile, 1, 1)),
		NewToken(StringLiteralToken, `b""`, position.NewPosition(testFile, 1, 7)),
	})
	testTokenize(t, `<?php B""`, []*Token{
		NewToken(StartTagToken, "", position.NewPosition(testFile, 1, 1)),
		NewToken(StringLiteralToken, `B""`, position.NewPosition(testFile, 1, 7)),
	})
	testTokenize(t, `<?php ""`, []*Token{
		NewToken(StartTagToken, "", position.NewPosition(testFile, 1, 1)),
		NewToken(StringLiteralToken, `""`, position.NewPosition(testFile, 1, 7)),
	})
	testTokenize(t, `<?php "abc"`, []*Token{
		NewToken(StartTagToken, "", position.NewPosition(testFile, 1, 1)),
		NewToken(StringLiteralToken, `"abc"`, position.NewPosition(testFile, 1, 7)),
	})
	testTokenize(t, `<?php "\"abc\\\n\$"`, []*Token{
		NewToken(StartTagToken, "", position.NewPosition(testFile, 1, 1)),
		NewToken(StringLiteralToken, `"\"abc\\\n\$"`, position.NewPosition(testFile, 1, 7)),
	})
}

func TestOperatorOrPunctuator(t *testing.T) {
	testTokenize(t, "<?php a === 1", []*Token{
		NewToken(StartTagToken, "", position.NewPosition(testFile, 1, 1)),
		NewToken(NameToken, "a", position.NewPosition(testFile, 1, 7)),
		NewToken(OpOrPuncToken, "===", position.NewPosition(testFile, 1, 9)),
	})

	testTokenize(t, "<?php a == 1", []*Token{
		NewToken(StartTagToken, "", position.NewPosition(testFile, 1, 1)),
		NewToken(NameToken, "a", position.NewPosition(testFile, 1, 7)),
		NewToken(OpOrPuncToken, "==", position.NewPosition(testFile, 1, 9)),
	})

	testTokenize(t, "<?php a = 1", []*Token{
		NewToken(StartTagToken, "", position.NewPosition(testFile, 1, 1)),
		NewToken(NameToken, "a", position.NewPosition(testFile, 1, 7)),
		NewToken(OpOrPuncToken, "=", position.NewPosition(testFile, 1, 9)),
	})

	testTokenize(t, "<?php a ** 1", []*Token{
		NewToken(StartTagToken, "", position.NewPosition(testFile, 1, 1)),
		NewToken(NameToken, "a", position.NewPosition(testFile, 1, 7)),
		NewToken(OpOrPuncToken, "**", position.NewPosition(testFile, 1, 9)),
	})
}

func TestVariableVarname(t *testing.T) {
	testTokenize(t, `<?php $$var = "someValue";`,
		[]*Token{
			NewToken(StartTagToken, "", position.NewPosition(testFile, 1, 1)),
			NewToken(OpOrPuncToken, "$", position.NewPosition(testFile, 1, 7)),
			NewToken(VariableNameToken, "$var", position.NewPosition(testFile, 1, 8)),
			NewToken(OpOrPuncToken, "=", position.NewPosition(testFile, 1, 13)),
			NewToken(StringLiteralToken, `"someValue"`, position.NewPosition(testFile, 1, 15)),
			NewToken(OpOrPuncToken, ";", position.NewPosition(testFile, 1, 26)),
		})

	testTokenize(t, `<?php echo 12, $var;`,
		[]*Token{
			NewToken(StartTagToken, "", position.NewPosition(testFile, 1, 1)),
			NewToken(KeywordToken, "echo", position.NewPosition(testFile, 1, 7)),
			NewToken(IntegerLiteralToken, "12", position.NewPosition(testFile, 1, 12)),
			NewToken(OpOrPuncToken, ",", position.NewPosition(testFile, 1, 14)),
			NewToken(VariableNameToken, "$var", position.NewPosition(testFile, 1, 16)),
		})
}

func TestHtmlAndPhp(t *testing.T) {
	testTokenize(t, "<body>\n"+`    <?php $heading = "My Heading"; ?>`+"\n    <h1><?= $heading ?></h1>",
		[]*Token{
			NewToken(TextToken, "<body>\n    ", nil),
			NewToken(StartTagToken, "", position.NewPosition(testFile, 2, 5)),
			NewToken(VariableNameToken, "$heading", position.NewPosition(testFile, 2, 11)),
			NewToken(OpOrPuncToken, "=", position.NewPosition(testFile, 2, 20)),
			NewToken(StringLiteralToken, `"My Heading"`, position.NewPosition(testFile, 2, 22)),
			NewToken(OpOrPuncToken, ";", position.NewPosition(testFile, 2, 34)),
			NewToken(EndTagToken, "", position.NewPosition(testFile, 2, 36)),
			NewToken(TextToken, "    <h1>", nil),
			NewToken(StartTagToken, "", position.NewPosition(testFile, 3, 9)),
			NewToken(KeywordToken, "echo", position.NewPosition(testFile, 3, 9)),
			NewToken(VariableNameToken, "$heading", position.NewPosition(testFile, 3, 13)),
			NewToken(OpOrPuncToken, ";", position.NewPosition(testFile, 3, 22)),
			NewToken(EndTagToken, "", position.NewPosition(testFile, 3, 22)),
			NewToken(TextToken, "</h1>", nil),
		})
}
