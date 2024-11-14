package common

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

// ------------------- MARK: Integer -------------------

func IsDigit(char string) bool {
	// Spec: https://phplang.org/spec/19-grammar.html#grammar-digit

	// digit:: one of
	//    0   1   2   3   4   5   6   7   8   9

	match, _ := regexp.MatchString(`^\d$`, char)
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

	match, _ := regexp.MatchString(`^[1-9]\d*(_\d+)*$`, str)
	return match
}

func DecimalLiteralToInt64(str string) int64 {
	integer, _ := strconv.ParseInt(ReplaceUnderscores(str), 10, 64)
	return integer
}

func IsOctalLiteral(str string) bool {
	// Spec: https://phplang.org/spec/09-lexical-structure.html#grammar-octal-literal

	// octal-literal::
	//    0
	//    octal-literal   octal-digit

	// octal-digit:: one of
	//    0   1   2   3   4   5   6   7

	match, _ := regexp.MatchString("^[0-7]+(_[0-7]+)*$", str)
	return match
}

func OctalLiteralToInt64(str string) int64 {
	integer, _ := strconv.ParseInt(ReplaceUnderscores(str), 8, 64)
	return integer
}

func IsHexadecimalDigit(char string) bool {
	// Spec: https://phplang.org/spec/19-grammar.html#grammar-hexadecimal-digit

	// hexadecimal-digit:: one of
	//    0   1   2   3   4   5   6   7   8   9
	//    a   b   c   d   e   f
	//    A   B   C   D   E   F

	match, _ := regexp.MatchString(`^[\da-fA-F]$`, char)
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

	match, _ := regexp.MatchString(`^0[xX][\da-fA-F]+(_[\da-fA-F]+)*$`, str)
	return match
}

func HexadecimalLiteralToInt64(str string) int64 {
	integer, _ := strconv.ParseInt(ReplaceUnderscores(str[2:]), 16, 64)
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

	match, _ := regexp.MatchString("^0[bB][01]+(_[01]+)*$", str)
	return match
}

func BinaryLiteralToInt64(str string) int64 {
	integer, _ := strconv.ParseInt(ReplaceUnderscores(str[2:]), 2, 64)
	return integer
}

func IsIntegerLiteral(str string) bool {
	return IsDecimalLiteral(str) || IsOctalLiteral(str) || IsHexadecimalLiteral(str) || IsBinaryLiteral(str)
}

func IsIntegerLiteralWithSign(str string) bool {
	if strings.HasPrefix(str, "-") || strings.HasPrefix(str, "+") {
		str = str[1:]
	}
	return IsIntegerLiteral(str)
}

func IntegerLiteralToInt64WithSign(str string) (int64, error) {
	isNegative := strings.HasPrefix(str, "-")
	if strings.HasPrefix(str, "+") || strings.HasPrefix(str, "-") {
		str = str[1:]
	}

	result, err := IntegerLiteralToInt64(str)
	if err != nil {
		return 0, err
	}
	if isNegative {
		result = -result
	}
	return result, nil
}

func IntegerLiteralToInt64(str string) (int64, error) {
	// decimal-literal
	if IsDecimalLiteral(str) {
		return DecimalLiteralToInt64(str), nil
	}

	// octal-literal
	if IsOctalLiteral(str) {
		return OctalLiteralToInt64(str), nil
	}

	// hexadecimal-literal
	if IsHexadecimalLiteral(str) {
		return HexadecimalLiteralToInt64(str), nil
	}

	// binary-literal
	if IsBinaryLiteral(str) {
		return BinaryLiteralToInt64(str), nil
	}

	return 0, fmt.Errorf("Given string is not an integer literal")
}

// ------------------- MARK: Float -------------------

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
	match, _ := regexp.MatchString(`^(\d*(_\d+)*\.\d+(_\d+)*|\d+(_\d+)*\.)([eE][+-]?\d+(_\d+)*)?$`, str)
	if match {
		return true
	}

	// digit-sequence   exponent-part
	match, _ = regexp.MatchString(`^\d+(_\d+)*[eE][+-]?\d+(_\d+)*$`, str)
	return match
}

func IsFloatingLiteralWithSign(str string) bool {
	if strings.HasPrefix(str, "-") || strings.HasPrefix(str, "+") {
		str = str[1:]
	}
	return IsFloatingLiteral(str)
}

func FloatingLiteralToFloat64(str string) float64 {
	float, _ := strconv.ParseFloat(str, 64)
	return float
}

func FloatingLiteralToFloat64WithSign(str string) float64 {
	isNegative := strings.HasPrefix(str, "-")
	if strings.HasPrefix(str, "+") || strings.HasPrefix(str, "-") {
		str = str[1:]
	}

	result := FloatingLiteralToFloat64(str)
	if isNegative {
		result = -result
	}
	return result
}
