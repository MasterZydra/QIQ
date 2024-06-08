package common

import "testing"

func TestReplaceControlChars(t *testing.T) {
	doTest := func(t *testing.T, input string, output string) {
		if got := ReplaceControlChars(input); got != output {
			t.Errorf("\nExpected: \"%s\".\nGot: \"%s\"", output, got)
		}
	}

	doTest(t, `\r`, "\r")
	doTest(t, `\n`, "\n")
	doTest(t, `\t`, "\t")
	doTest(t, `\\`, "\\")

	doTest(t, `\\n`, `\n`)
	doTest(t, `\\\n`, `\`+"\n")
	doTest(t, `\\hi\\\n`, `\hi\`+"\n")
}
