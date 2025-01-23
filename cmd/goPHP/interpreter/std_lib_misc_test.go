package interpreter

import (
	"GoPHP/cmd/goPHP/phpError"
	"testing"
)

// ------------------- MARK: constant -------------------

func TestLibConstant(t *testing.T) {
	testInputOutput(t, `<?php var_dump(constant('E_ALL'));`, "int(32767)\n")
	testForError(t, `<?php constant('NOT_DEFINED_CONSTANT');`, phpError.NewError("Undefined constant \"NOT_DEFINED_CONSTANT\""))
	// TODO Add test cases for user defined constants
}

// ------------------- MARK: defined -------------------

func TestLibDefined(t *testing.T) {
	testInputOutput(t, `<?php var_dump(defined('PHP_VERSION'));`, "bool(true)\n")
	testInputOutput(t, `<?php var_dump(defined('NOT_DEFINED_CONSTANT'));`, "bool(false)\n")
	// TODO Add test cases for user defined constants
}
