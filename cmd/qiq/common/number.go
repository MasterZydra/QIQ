package common

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

func addRegexEOL(str string, add bool) string {
	if add {
		return str + `$`
	}
	return str
}

// -------------------------------------- Integer -------------------------------------- MARK: Integer

func IsDigit(char string) bool {
	// Spec: https://phplang.org/spec/19-grammar.html#grammar-digit

	// digit:: one of
	//    0   1   2   3   4   5   6   7   8   9

	match, _ := regexp.MatchString(`^\d$`, char)
	return match
}

const decimalLiteralPattern = `^[1-9]\d*(_\d+)*`

func IsDecimalLiteral(str string, leadingNumeric bool) bool {
	// Spec: https://phplang.org/spec/09-lexical-structure.html#grammar-decimal-literal

	// decimal-literal::
	//    nonzero-digit
	//    decimal-literal   digit

	// nonzero-digit:: one of
	//    1   2   3   4   5   6   7   8   9

	// digit:: one of
	//    0   1   2   3   4   5   6   7   8   9

	match, _ := regexp.MatchString(addRegexEOL(decimalLiteralPattern, !leadingNumeric), str)
	return match
}

func DecimalLiteralToInt64(str string, leadingNumeric bool) int64 {
	r, _ := regexp.Compile(addRegexEOL(decimalLiteralPattern, !leadingNumeric))
	integer, _ := strconv.ParseInt(ReplaceUnderscores(r.FindString(str)), 10, 64)
	return integer
}

const octalLiteralPattern = `^[0-7]+(_[0-7]+)*`

func IsOctalLiteral(str string, leadingNumeric bool) bool {
	// Spec: https://phplang.org/spec/09-lexical-structure.html#grammar-octal-literal

	// octal-literal::
	//    0
	//    octal-literal   octal-digit

	// octal-digit:: one of
	//    0   1   2   3   4   5   6   7

	match, _ := regexp.MatchString(addRegexEOL(octalLiteralPattern, !leadingNumeric), str)
	return match
}

func OctalLiteralToInt64(str string, leadingNumeric bool) int64 {
	r, _ := regexp.Compile(addRegexEOL(octalLiteralPattern, !leadingNumeric))
	integer, _ := strconv.ParseInt(ReplaceUnderscores(r.FindString(str)), 8, 64)
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

const hexadecimalLiteralPattern = `^0[xX][\da-fA-F]+(_[\da-fA-F]+)*`

func IsHexadecimalLiteral(str string, leadingNumeric bool) bool {
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

	match, _ := regexp.MatchString(addRegexEOL(hexadecimalLiteralPattern, !leadingNumeric), str)
	return match
}

func HexadecimalLiteralToInt64(str string, leadingNumeric bool) int64 {
	r, _ := regexp.Compile(addRegexEOL(hexadecimalLiteralPattern, !leadingNumeric))
	integer, _ := strconv.ParseInt(ReplaceUnderscores(r.FindString(str)[2:]), 16, 64)
	return integer
}

const binaryLiteralPattern = `^0[bB][01]+(_[01]+)*`

func IsBinaryLiteral(str string, leadingNumeric bool) bool {
	// Spec: https://phplang.org/spec/09-lexical-structure.html#grammar-binary-literal

	// binary-literal::
	//    binary-prefix   binary-digit
	//    binary-literal   binary-digit

	// binary-prefix:: one of
	//    0b   0B

	// binary-digit:: one of
	//    0   1

	match, _ := regexp.MatchString(addRegexEOL(binaryLiteralPattern, !leadingNumeric), str)
	return match
}

func BinaryLiteralToInt64(str string, leadingNumeric bool) int64 {
	r, _ := regexp.Compile(addRegexEOL(binaryLiteralPattern, !leadingNumeric))
	integer, _ := strconv.ParseInt(ReplaceUnderscores(r.FindString(str)[2:]), 2, 64)
	return integer
}

func IsIntegerLiteral(str string, leadingNumeric bool) bool {
	return IsDecimalLiteral(str, leadingNumeric) ||
		IsOctalLiteral(str, leadingNumeric) ||
		IsHexadecimalLiteral(str, leadingNumeric) ||
		IsBinaryLiteral(str, leadingNumeric)
}

func IsIntegerLiteralWithSign(str string, leadingNumeric bool) bool {
	if strings.HasPrefix(str, "-") || strings.HasPrefix(str, "+") {
		str = str[1:]
	}
	return IsIntegerLiteral(str, leadingNumeric)
}

func IntegerLiteralToInt64WithSign(str string, leadingNumeric bool) (int64, error) {
	isNegative := strings.HasPrefix(str, "-")
	if strings.HasPrefix(str, "+") || strings.HasPrefix(str, "-") {
		str = str[1:]
	}

	result, err := IntegerLiteralToInt64(str, leadingNumeric)
	if err != nil {
		return 0, err
	}
	if isNegative {
		result = -result
	}
	return result, nil
}

func IntegerLiteralToInt64(str string, leadingNumeric bool) (int64, error) {
	// decimal-literal
	if IsDecimalLiteral(str, leadingNumeric) {
		return DecimalLiteralToInt64(str, leadingNumeric), nil
	}

	// octal-literal
	if IsOctalLiteral(str, leadingNumeric) {
		return OctalLiteralToInt64(str, leadingNumeric), nil
	}

	// hexadecimal-literal
	if IsHexadecimalLiteral(str, leadingNumeric) {
		return HexadecimalLiteralToInt64(str, leadingNumeric), nil
	}

	// binary-literal
	if IsBinaryLiteral(str, leadingNumeric) {
		return BinaryLiteralToInt64(str, leadingNumeric), nil
	}

	return 0, fmt.Errorf("Given string is not an integer literal")
}

// -------------------------------------- Float -------------------------------------- MARK: Float

const floatingFractLiteralPattern = `^(\d*(_\d+)*\.\d+(_\d+)*|\d+(_\d+)*\.)([eE][+-]?\d+(_\d+)*)?`

const floatingDigitLiteralPattern = `^\d+(_\d+)*[eE][+-]?\d+(_\d+)*`

const FloatingLiteralPattern = `(([+-]?\d*(_\d+)*\.\d+(_\d+)*|\d+(_\d+)*\.)([eE][+-]?\d+(_\d+)*)?)|([+-]?\d+(_\d+)*[eE][+-]?\d+(_\d+)*)`

func IsFloatingLiteral(str string, leadingNumeric bool) bool {
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
	match, _ := regexp.MatchString(addRegexEOL(floatingFractLiteralPattern, !leadingNumeric), str)
	if match {
		return true
	}

	// digit-sequence   exponent-part
	match, _ = regexp.MatchString(addRegexEOL(floatingDigitLiteralPattern, !leadingNumeric), str)
	return match
}

func IsFloatingLiteralWithSign(str string, leadingNumeric bool) bool {
	if strings.HasPrefix(str, "-") || strings.HasPrefix(str, "+") {
		str = str[1:]
	}
	return IsFloatingLiteral(str, leadingNumeric)
}

func FloatingLiteralToFloat64(str string, leadingNumeric bool) float64 {
	// fractional-literal   exponent-part(opt)
	match, _ := regexp.MatchString(addRegexEOL(floatingFractLiteralPattern, !leadingNumeric), str)
	if match {
		r, _ := regexp.Compile(addRegexEOL(floatingFractLiteralPattern, !leadingNumeric))
		str = r.FindString(str)
	} else {
		r, _ := regexp.Compile(addRegexEOL(floatingDigitLiteralPattern, !leadingNumeric))
		str = r.FindString(str)
	}

	float, _ := strconv.ParseFloat(str, 64)
	return float
}

func FloatingLiteralToFloat64WithSign(str string, leadingNumeric bool) float64 {
	isNegative := strings.HasPrefix(str, "-")
	if strings.HasPrefix(str, "+") || strings.HasPrefix(str, "-") {
		str = str[1:]
	}

	result := FloatingLiteralToFloat64(str, leadingNumeric)
	if isNegative {
		result = -result
	}
	return result
}
