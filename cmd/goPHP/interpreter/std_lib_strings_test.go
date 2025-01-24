package interpreter

import (
	"GoPHP/cmd/goPHP/phpError"
	"testing"
)

// ------------------- MARK: bin2hex -------------------

func TestLibBin2Hex(t *testing.T) {
	testInputOutput(t, `<?php var_dump(bin2hex('Hello world!'));`, "string(24) \"48656c6c6f20776f726c6421\"\n")
	testInputOutput(t, `<?php var_dump(bin2hex('Äàßê'));`, "string(16) \"c384c3a0c39fc3aa\"\n")
	testInputOutput(t, `<?php var_dump(bin2hex(''));`, "string(0) \"\"\n")
}

// ------------------- MARK: chr -------------------

func TestLibChr(t *testing.T) {
	testInputOutput(t, `<?php var_dump(chr(60));`, "string(1) \"<\"\n")
	testInputOutput(t, `<?php var_dump(chr(60-256));`, "string(1) \"<\"\n")
	testInputOutput(t, `<?php var_dump(chr(60+256));`, "string(1) \"<\"\n")
}

// ------------------- MARK: lcfirst -------------------

func TestLibLcFirst(t *testing.T) {
	testInputOutput(t, `<?php var_dump(lcfirst('ABC'));`, "string(3) \"aBC\"\n")
	testInputOutput(t, `<?php var_dump(lcfirst('Abc'));`, "string(3) \"abc\"\n")
	testInputOutput(t, `<?php var_dump(lcfirst('abc'));`, "string(3) \"abc\"\n")
	testInputOutput(t, `<?php var_dump(lcfirst(''));`, "string(0) \"\"\n")
}

// ------------------- MARK: quotemeta -------------------

func TestLibQuoteMeta(t *testing.T) {
	testInputOutput(t, `<?php var_dump(quotemeta('. \ + * ? [ ^ ] ( $ )'));`, `string(31) "\. \\ \+ \* \? \[ \^ ] \( \$ \)"`+"\n")
	testInputOutput(t, `<?php var_dump(quotemeta('Hello. (can you hear me?)'));`, `string(29) "Hello\. \(can you hear me\?\)"`+"\n")
	testInputOutput(t, `<?php var_dump(quotemeta(''));`, "bool(false)\n")
}

// ------------------- MARK: str_contains -------------------

func TestLibStrContains(t *testing.T) {
	testInputOutput(t, `<?php var_dump(str_contains('abc', ''));`, "bool(true)\n")
	testInputOutput(t, `<?php var_dump(str_contains('The lazy fox', 'lazy'));`, "bool(true)\n")
	testInputOutput(t, `<?php var_dump(str_contains('The lazy fox', 'Lazy'));`, "bool(false)\n")
}

// ------------------- MARK: str_ends_with -------------------

func TestLibStrEndsWith(t *testing.T) {
	testInputOutput(t, `<?php var_dump(str_ends_with('abc', ''));`, "bool(true)\n")
	testInputOutput(t, `<?php var_dump(str_ends_with('The lazy fox', 'fox'));`, "bool(true)\n")
	testInputOutput(t, `<?php var_dump(str_ends_with('The lazy fox', 'Fox'));`, "bool(false)\n")
}

// ------------------- MARK: str_repeat -------------------

func TestLibStrRepeat(t *testing.T) {
	testInputOutput(t, `<?php var_dump(str_repeat('abc', 0));`, "string(0) \"\"\n")
	testInputOutput(t, `<?php var_dump(str_repeat('abc', 1));`, "string(3) \"abc\"\n")
	testInputOutput(t, `<?php var_dump(str_repeat('abc', 2));`, "string(6) \"abcabc\"\n")
	testForError(t, `<?php var_dump(str_repeat('abc', -1));`, phpError.NewError("Uncaught ValueError: str_repeat(): Argument #2 ($times) must be greater than or equal to 0"))
}

// ------------------- MARK: str_starts_with -------------------

func TestLibStrStartsWith(t *testing.T) {
	testInputOutput(t, `<?php var_dump(str_starts_with('abc', ''));`, "bool(true)\n")
	testInputOutput(t, `<?php var_dump(str_starts_with('The lazy fox', 'The'));`, "bool(true)\n")
	testInputOutput(t, `<?php var_dump(str_starts_with('The lazy fox', 'the'));`, "bool(false)\n")
}

// ------------------- MARK: strlen -------------------

func TestLibStrlen(t *testing.T) {
	testInputOutput(t, `<?php var_dump(strlen('abcdef'));`, "int(6)\n")
	testInputOutput(t, `<?php var_dump(strlen(' ab cd '));`, "int(7)\n")
	testInputOutput(t, `<?php var_dump(strlen(' äb ćd '));`, "int(9)\n")
}

// ------------------- MARK: strtolower -------------------

func TestLibStrToLower(t *testing.T) {
	testInputOutput(t, `<?php var_dump(strtolower('Mary Had A Little Lamb and She LOVED It So'));`, "string(42) \"mary had a little lamb and she loved it so\"\n")
	testInputOutput(t, `<?php var_dump(strtolower(''));`, "string(0) \"\"\n")
	testInputOutput(t, `<?php var_dump(strtolower('AÄOÖUÜSß'));`, "string(12) \"aÄoÖuÜsß\"\n")
}

// ------------------- MARK: strtoupper -------------------

func TestLibStrToUpper(t *testing.T) {
	testInputOutput(t, `<?php var_dump(strtoupper('Mary Had A Little Lamb and She LOVED It So'));`, "string(42) \"MARY HAD A LITTLE LAMB AND SHE LOVED IT SO\"\n")
	testInputOutput(t, `<?php var_dump(strtoupper(''));`, "string(0) \"\"\n")
	testInputOutput(t, `<?php var_dump(strtoupper('aäoöuüsß'));`, "string(12) \"AäOöUüSß\"\n")
}

// ------------------- MARK: ucfirst -------------------

func TestLibUcFirst(t *testing.T) {
	testInputOutput(t, `<?php var_dump(ucfirst('ABC'));`, "string(3) \"ABC\"\n")
	testInputOutput(t, `<?php var_dump(ucfirst('abc'));`, "string(3) \"Abc\"\n")
	testInputOutput(t, `<?php var_dump(ucfirst('Abc'));`, "string(3) \"Abc\"\n")
	testInputOutput(t, `<?php var_dump(ucfirst(''));`, "string(0) \"\"\n")
}
