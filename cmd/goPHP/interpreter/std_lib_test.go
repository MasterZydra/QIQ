package interpreter

import (
	"GoPHP/cmd/goPHP/runtime/values"
	"testing"
)

func TestArrayKeyExists(t *testing.T) {
	array := values.NewArray()
	array.SetElement(nil, values.NewInt(42))
	if actual, _ := lib_array_key_exists(values.NewInt(0), array); !actual {
		t.Errorf("Expected: \"%t\", Got \"%t\"", true, actual)
	}
	if actual, _ := lib_array_key_exists(values.NewInt(1), array); actual {
		t.Errorf("Expected: \"%t\", Got \"%t\"", false, actual)
	}
}

func TestLibErrorReporting(t *testing.T) {
	// Spec: https://www.php.net/manual/en/function.error-reporting.php - Example #1

	// Turn off all error reporting
	testInputOutput(t, `<?php error_reporting(0); echo error_reporting();`, "0")
	// Report simple running errors
	testInputOutput(t, `<?php error_reporting(E_ERROR | E_WARNING | E_PARSE); echo error_reporting();`, "7")
	// Reporting E_NOTICE can be good too (to report uninitialized variables or catch variable name misspellings ...)
	testInputOutput(t, `<?php error_reporting(E_ERROR | E_WARNING | E_PARSE | E_NOTICE); echo error_reporting();`, "15")
	// Report all errors except E_NOTICE
	testInputOutput(t, `<?php error_reporting(E_ALL & ~E_NOTICE); echo error_reporting();`, "32759")
	// Report all PHP errors
	testInputOutput(t, `<?php error_reporting(E_ALL); echo error_reporting();`, "32767")
	// Report all PHP errors
	testInputOutput(t, `<?php error_reporting(-1); echo error_reporting();`, "32767")
}
