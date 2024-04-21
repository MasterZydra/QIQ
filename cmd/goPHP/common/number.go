package common

import (
	"regexp"
	"strconv"
)

func IsDigit(char string) bool {
	// Spec: https://phplang.org/spec/19-grammar.html#grammar-digit

	// digit:: one of
	//    0   1   2   3   4   5   6   7   8   9

	match, _ := regexp.MatchString("^[0-9]$", char)
	return match
}

func IsDecimalLiteral(str string) bool {
	// Spec: https://phplang.org/spec/09-lexical-structure.html#grammar-decimal-literal

	// decimal-literal::
	//    nonzero-digit
	//    decimal-literal   digit

	// nonzero-digit:: one of
	//    1   2   3   4   5   6   7   8   9

	// digit:: one of
	//    0   1   2   3   4   5   6   7   8   9

	match, _ := regexp.MatchString("^[1-9][0-9]*$", str)
	return match
}

func DecimalLiteralToInt64(str string) int64 {
	integer, _ := strconv.ParseInt(str, 10, 64)
	return integer
}

func IsOctalLiteral(str string) bool {
	// Spec: https://phplang.org/spec/09-lexical-structure.html#grammar-octal-literal

	// octal-literal::
	//    0
	//    octal-literal   octal-digit

	// octal-digit:: one of
	//    0   1   2   3   4   5   6   7

	match, _ := regexp.MatchString("^[0-7]+$", str)
	return match
}

func OctalLiteralToInt64(str string) int64 {
	integer, _ := strconv.ParseInt(str, 8, 64)
	return integer
}

func IsHexadecimalDigit(char string) bool {
	// Spec: https://phplang.org/spec/19-grammar.html#grammar-hexadecimal-digit

	// hexadecimal-digit:: one of
	//    0   1   2   3   4   5   6   7   8   9
	//    a   b   c   d   e   f
	//    A   B   C   D   E   F

	match, _ := regexp.MatchString("^[0-9a-fA-F]$", char)
	return match
}

func IsHexadecimalLiteral(str string) bool {
	// Spec: https://phplang.org/spec/09-lexical-structure.html#grammar-hexadecimal-literal

	// hexadecimal-literal::
	//    hexadecimal-prefix   hexadecimal-digit
	//    hexadecimal-literal   hexadecimal-digit

	// hexadecimal-prefix:: one of
	//    0x   0X

	// hexadecimal-digit:: one of
	//    0   1   2   3   4   5   6   7   8   9
	//    a   b   c   d   e   f
	//    A   B   C   D   E   F

	match, _ := regexp.MatchString("^0[xX][0-9a-fA-F]+$", str)
	return match
}

func HexadecimalLiteralToInt64(str string) int64 {
	integer, _ := strconv.ParseInt(str[2:], 16, 64)
	return integer
}

func IsBinaryLiteral(str string) bool {
	// Spec: https://phplang.org/spec/09-lexical-structure.html#grammar-binary-literal

	// binary-literal::
	//    binary-prefix   binary-digit
	//    binary-literal   binary-digit

	// binary-prefix:: one of
	//    0b   0B

	// binary-digit:: one of
	//    0   1

	match, _ := regexp.MatchString("^0[bB][01]+$", str)
	return match
}

func BinaryLiteralToInt64(str string) int64 {
	integer, _ := strconv.ParseInt(str[2:], 2, 64)
	return integer
}

func IsFloatingLiteral(str string) bool {
	// Spec: https://phplang.org/spec/09-lexical-structure.html#grammar-floating-literal

	// floating-literal::
	//    fractional-literal   exponent-part(opt)
	//    digit-sequence   exponent-part

	// fractional-literal::
	//    digit-sequence(opt)   .   digit-sequence
	//    digit-sequence   .

	// exponent-part::
	//    e   sign(opt)   digit-sequence
	//    E   sign(opt)   digit-sequence

	// sign:: one of
	//    +   -

	// digit-sequence::
	//    digit
	//    digit-sequence   digit

	// fractional-literal   exponent-part(opt)
	match, _ := regexp.MatchString(`^([0-9]*\.[0-9]+|[0-9]+\.)([eE][+-]?[0-9]+)?$`, str)
	if match {
		return true
	}

	// digit-sequence   exponent-part
	match, _ = regexp.MatchString("^[0-9]+[eE][+-]?[0-9]+$", str)
	return match
}

func FloatingLiteralToFloat64(str string) float64 {
	float, _ := strconv.ParseFloat(str, 64)
	return float
}
