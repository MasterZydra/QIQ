package interpreter

import "testing"

func TestObFunctions(t *testing.T) {
	// Implicit flush
	testInputOutput(t, `<?php ob_start(); echo '123';`, "123")
	// ob_clean
	testInputOutput(t, `<?php ob_start(); echo '123'; ob_clean(); echo '456';`, "456")
	// ob_flush
	testInputOutput(t, `<?php ob_start(); echo '123'; ob_flush(); echo '456'; ob_end_clean();`, "123")
	// ob_end_clean
	testInputOutput(t, `<?php ob_start(); echo '123'; ob_end_clean(); echo '456';`, "456")
	// ob_end_flush
	testInputOutput(t, `<?php ob_start(); echo '123'; ob_end_flush(); echo '456';`, "123456")
	// ob_get_clean
	testInputOutput(t, `<?php ob_start(); echo '123'; $ob = ob_get_clean(); echo '456' . $ob;`, "456123")
	// ob_get_flush
	testInputOutput(t, `<?php ob_start(); echo '123'; $ob = ob_get_flush(); echo '456' . $ob;`, "123456123")
	// ob_get_contents
	testInputOutput(t, `<?php ob_start(); echo '123'; $ob = ob_get_contents(); ob_end_clean(); echo '456' . $ob;`, "456123")
	// ob_get_level
	testInputOutput(t, `<?php ob_start(); echo 'A' . ob_get_level(); ob_start(); echo 'B' . ob_get_level();`, "A1B2")
	// Stacked output buffers
	testInputOutput(t,
		`<?php
            echo 0;
                ob_start();
                    ob_start();
                        ob_start();
                            ob_start();
                                echo 1;
                            ob_end_flush();
                            echo 2;
                        $ob = ob_get_clean();
                    echo 3;
                    ob_flush();
                    ob_end_clean();
                echo 4;
                ob_end_flush();
            echo '-' . $ob;
        ?>`,
		"034-12")
}
