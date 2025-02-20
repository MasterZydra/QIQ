package interpreter

import "testing"

// ------------------- MARK: checkdate -------------------

func TestLibCheckdate(t *testing.T) {
	testInputOutput(t, `<?php var_dump(checkdate(12, 31, 2000));`, "bool(true)\n")
	testInputOutput(t, `<?php var_dump(checkdate(2, 29, 2001));`, "bool(false)\n")
}

// ------------------- MARK: date -------------------

func TestLibDate(t *testing.T) {
	// Day
	testInputOutput(t, `<?php var_dump(date('d', mktime(12, 13, 14, 05, 04, 2024)));`, "string(2) \"04\"\n")
	testInputOutput(t, `<?php var_dump(date('j', mktime(12, 13, 14, 05, 04, 2024)));`, "string(1) \"4\"\n")
	testInputOutput(t, `<?php var_dump(date('z', mktime(12, 13, 14, 05, 04, 2024)));`, "string(3) \"124\"\n")
	testInputOutput(t, `<?php var_dump(date('w', mktime(12, 13, 14, 05, 04, 2024)));`, "string(1) \"6\"\n")
	testInputOutput(t, `<?php var_dump(date('w', mktime(12, 13, 14, 05, 05, 2024)));`, "string(1) \"0\"\n")
	testInputOutput(t, `<?php var_dump(date('N', mktime(12, 13, 14, 05, 04, 2024)));`, "string(1) \"6\"\n")
	testInputOutput(t, `<?php var_dump(date('N', mktime(12, 13, 14, 05, 05, 2024)));`, "string(1) \"7\"\n")
	// Week
	testInputOutput(t, `<?php var_dump(date('W', mktime(12, 13, 14, 05, 04, 2024)));`, "string(2) \"18\"\n")
	// Month
	testInputOutput(t, `<?php var_dump(date('m', mktime(12, 13, 14, 05, 04, 2024)));`, "string(2) \"05\"\n")
	testInputOutput(t, `<?php var_dump(date('n', mktime(12, 13, 14, 05, 04, 2024)));`, "string(1) \"5\"\n")
	testInputOutput(t, `<?php var_dump(date('t', mktime(12, 13, 14, 05, 04, 2024)));`, "string(2) \"31\"\n")
	// Year
	testInputOutput(t, `<?php var_dump(date('Y', mktime(12, 13, 14, 05, 04, 2024)));`, "string(4) \"2024\"\n")
	testInputOutput(t, `<?php var_dump(date('y', mktime(12, 13, 14, 05, 04, 2024)));`, "string(2) \"24\"\n")
	testInputOutput(t, `<?php var_dump(date('L', mktime(12, 13, 14, 05, 04, 2024)));`, "string(1) \"1\"\n")
	// Time
	testInputOutput(t, `<?php var_dump(date('i', mktime(12, 13, 14, 05, 04, 2024)));`, "string(2) \"13\"\n")
	testInputOutput(t, `<?php var_dump(date('i', mktime(12, 00, 14, 05, 04, 2024)));`, "string(2) \"00\"\n")
	testInputOutput(t, `<?php var_dump(date('s', mktime(12, 13, 14, 05, 04, 2024)));`, "string(2) \"14\"\n")
	testInputOutput(t, `<?php var_dump(date('s', mktime(12, 13, 00, 05, 04, 2024)));`, "string(2) \"00\"\n")
	testInputOutput(t, `<?php var_dump(date('G', mktime(12, 13, 14, 05, 04, 2024)));`, "string(2) \"12\"\n")
	testInputOutput(t, `<?php var_dump(date('G', mktime(20, 13, 14, 05, 04, 2024)));`, "string(2) \"20\"\n")
	testInputOutput(t, `<?php var_dump(date('H', mktime(12, 13, 14, 05, 04, 2024)));`, "string(2) \"12\"\n")
	testInputOutput(t, `<?php var_dump(date('H', mktime(03, 13, 14, 05, 04, 2024)));`, "string(2) \"03\"\n")
	testInputOutput(t, `<?php var_dump(date('H', mktime(20, 13, 14, 05, 04, 2024)));`, "string(2) \"20\"\n")
	testInputOutput(t, `<?php var_dump(date('g', mktime(00, 13, 14, 05, 04, 2024)));`, "string(2) \"12\"\n")
	testInputOutput(t, `<?php var_dump(date('g', mktime(12, 13, 14, 05, 04, 2024)));`, "string(2) \"12\"\n")
	testInputOutput(t, `<?php var_dump(date('g', mktime(14, 13, 14, 05, 04, 2024)));`, "string(1) \"2\"\n")
	testInputOutput(t, `<?php var_dump(date('h', mktime(00, 13, 14, 05, 04, 2024)));`, "string(2) \"12\"\n")
	testInputOutput(t, `<?php var_dump(date('h', mktime(12, 13, 14, 05, 04, 2024)));`, "string(2) \"12\"\n")
	testInputOutput(t, `<?php var_dump(date('h', mktime(14, 13, 14, 05, 04, 2024)));`, "string(2) \"02\"\n")
}

// ------------------- MARK: getdate -------------------

func TestLibGetdate(t *testing.T) {
	testInputOutput(t, `<?php print_r(getdate(1722707036));`,
		"Array\n(\n    [seconds] => 56\n    [minutes] => 43\n    [hours] => 17\n    [mday] => 3\n    [wday] => 6\n    [mon] => 8\n"+
			"    [year] => 2024\n    [yday] => 215\n    [weekday] => Saturday\n    [month] => August\n    [0] => 1722707036\n)\n",
	)
}
