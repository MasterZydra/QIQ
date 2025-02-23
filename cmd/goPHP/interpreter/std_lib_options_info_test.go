package interpreter

import "testing"

// ------------------- MARK: ini_get -------------------

func TestLibIniGet(t *testing.T) {
	testInputOutput(t, `<?php var_dump(ini_get('none_existing'));`, "bool(false)\n")
	testInputOutput(t, `<?php var_dump(ini_get('variables_order'));`, "string(5) \"EGPCS\"\n")
	testInputOutput(t, `<?php var_dump(ini_get('error_reporting'));`, "string(5) \"32767\"\n")
}

// ------------------- MARK: ini_set -------------------

func TestLibIniSet(t *testing.T) {
	testInputOutput(t, `<?php var_dump(ini_set('none_existing', true));`, "bool(false)\n")
	testInputOutput(t, `<?php var_dump(ini_set('variables_order', true));`, "bool(false)\n")
	testInputOutput(t, `<?php var_dump(ini_set('error_reporting', E_ERROR)); var_dump(ini_set('error_reporting', E_ERROR));`, "string(5) \"32767\"\nstring(1) \"1\"\n")
}
