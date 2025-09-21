package common

import "testing"

func TestReplaceSingleQuoteControlChars(t *testing.T) {
	doTest := func(t *testing.T, input string, output string) {
		if got := ReplaceSingleQuoteControlChars(input); got != output {
			t.Errorf("\nExpected: \"%s\".\nGot: \"%s\"", output, got)
		}
	}

	doTest(t, `\r`, `\r`)
	doTest(t, `\n`, `\n`)
	doTest(t, `\t`, `\t`)
	doTest(t, `\\`, `\`)
	doTest(t, `\'`, "'")

	doTest(t, `\\n`, `\n`)
	doTest(t, `\\\n`, `\\n`)
	doTest(t, `\\hi\\\n`, `\hi\\n`)
	doTest(t, `\n\\\'a\\\b\\`, `\n\'a\\b\`)
}

func TestReplaceDoubleQuoteControlChars(t *testing.T) {
	doTest := func(t *testing.T, input string, output string) {
		if got := ReplaceDoubleQuoteControlChars(input); got != output {
			t.Errorf("\nExpected: \"%s\".\nGot: \"%s\"", output, got)
		}
	}

	doTest(t, `\r`, "\r")
	doTest(t, `\n`, "\n")
	doTest(t, `\t`, "\t")
	doTest(t, `\\`, `\`)

	doTest(t, `\\n`, `\n`)
	doTest(t, `\\\n`, `\`+"\n")
	doTest(t, `\\hi\\\n`, `\hi\`+"\n")
}

func TestReplaceAtPos(t *testing.T) {
	doTest := func(t *testing.T, input string, new string, pos int, output string) {
		if got := ReplaceAtPos(input, new, pos); got != output {
			t.Errorf("\nExpected: \"%s\".\nGot: \"%s\"", output, got)
		}
	}

	doTest(t, "abc", "*", 1, "a*c")
	doTest(t, "abc", "**", 1, "a*c")
	doTest(t, "abc", "**", 0, "*bc")
}

func TestExtendWithSpaces(t *testing.T) {
	doTest := func(t *testing.T, input string, length int, output string) {
		if got := ExtendWithSpaces(input, length); got != output {
			t.Errorf("\nExpected: \"%s\".\nGot: \"%s\"", output, got)
		}
	}

	doTest(t, "abc", 0, "abc")
	doTest(t, "abc", 1, "abc")
	doTest(t, "abc", 4, "abc ")
	doTest(t, "abc", 10, "abc       ")
}

func TestIsName(t *testing.T) {
	doTest := func(t *testing.T, input string, output bool) {
		if got := IsName(input); got != output {
			t.Errorf("\nExpected: \"%v\" Got: \"%v\"", output, got)
		}
	}

	doTest(t, "abc", true)
	doTest(t, "abc123", true)
	doTest(t, "_abc", true)
	doTest(t, "__abc", true)
	doTest(t, "_a1", true)
	doTest(t, "_1", true)

	doTest(t, " ", false)
	doTest(t, "1", false)
	doTest(t, "1abc", false)
	doTest(t, "$abc", false)
}

func TestIsQualifiedName(t *testing.T) {
	doTest := func(t *testing.T, input string, output bool) {
		if got := IsQualifiedName(input); got != output {
			t.Errorf("\nExpected: \"%v\" Got: \"%v\"", output, got)
		}
	}

	doTest(t, "abc", true)
	doTest(t, `a\b\c\d`, true)
	doTest(t, `\a\b\c`, true)
	doTest(t, `\Exception`, true)

	doTest(t, `a\b\`, false)
	doTest(t, `a/b`, false)
}

func TestIsVariableName(t *testing.T) {
	doTest := func(t *testing.T, input string, output bool) {
		if got := IsVariableName(input); got != output {
			t.Errorf("\nExpected: \"%v\" Got: \"%v\"", output, got)
		}
	}

	doTest(t, "$_", true)
	doTest(t, "$var", true)
	doTest(t, "$var123", true)

	doTest(t, "$", false)
	doTest(t, "$0", false)
	doTest(t, "var", false)
}
