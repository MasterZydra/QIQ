package interpreter

import "testing"

// ------------------- MARK: checkdate -------------------

func TestLibCheckdate(t *testing.T) {
	testInputOutput(t, `<?php var_dump(checkdate(12, 31, 2000));`, "bool(true)\n")
	testInputOutput(t, `<?php var_dump(checkdate(2, 29, 2001));`, "bool(false)\n")
}

// ------------------- MARK: getdate -------------------

func TestLibGetdate(t *testing.T) {
	testInputOutput(t, `<?php print_r(getdate(1722707036));`,
		"Array\n(\n    [seconds] => 56\n    [minutes] => 43\n    [hours] => 17\n    [mday] => 3\n    [wday] => 6\n    [mon] => 8\n"+
			"    [year] => 2024\n    [yday] => 215\n    [weekday] => Saturday\n    [month] => August\n    [0] => 1722707036\n)",
	)
}
