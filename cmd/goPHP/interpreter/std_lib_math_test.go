package interpreter

import "testing"

// ------------------- MARK: abs -------------------

func TestLibAbs(t *testing.T) {
	testInputOutput(t, `<?php var_dump(abs(-4.2));`, "float(4.2)\n")
	testInputOutput(t, `<?php var_dump(abs(5));`, "int(5)\n")
	testInputOutput(t, `<?php var_dump(abs(-5));`, "int(5)\n")
}

// ------------------- MARK: acos -------------------

func TestLibAcos(t *testing.T) {
	testInputOutput(t, `<?php var_dump(acos(1.0));`, "float(0)\n")
	testInputOutput(t, `<?php var_dump(acos(0.5)/M_PI*180);`, "float(60)\n")
}

// ------------------- MARK: acosh -------------------

func TestLibAcosh(t *testing.T) {
	testInputOutput(t, `<?php var_dump(acosh(1.0));`, "float(0)\n")
}

// ------------------- MARK: asin -------------------

func TestLibAsin(t *testing.T) {
	testInputOutput(t, `<?php var_dump(asin(0.0));`, "float(0)\n")
}

// ------------------- MARK: asinh -------------------

func TestLibAsinh(t *testing.T) {
	testInputOutput(t, `<?php var_dump(asinh(0.0));`, "float(0)\n")
}

// ------------------- MARK: pi -------------------

func TestLibPi(t *testing.T) {
	testInputOutput(t, `<?php var_dump(M_PI === pi());`, "bool(true)\n")
}
