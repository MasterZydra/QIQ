package common

import (
	"QIQ/cmd/qiq/common/os"
	"testing"
)

func TestGetAbsPathFromWorkingDir(t *testing.T) {
	if os.IS_WIN {
		got := GetAbsPathForWorkingDir("QIQ/cmd/qiq", "../../test.php")
		expected := "QIQ\\test.php"
		if got != expected {
			t.Errorf("\nExpected: \"%s\".\nGot: \"%s\"", expected, got)
		}

		got = GetAbsPathForWorkingDir("QIQ\\cmd\\qiq", "../../test.php")
		expected = "QIQ\\test.php"
		if got != expected {
			t.Errorf("\nExpected: \"%s\".\nGot: \"%s\"", expected, got)
		}
	} else {
		got := GetAbsPathForWorkingDir("QIQ/cmd/qiq", "../../test.php")
		expected := "QIQ/test.php"
		if got != expected {
			t.Errorf("\nExpected: \"%s\".\nGot: \"%s\"", expected, got)
		}
	}
}
