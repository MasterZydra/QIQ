package interpreter

import (
	"GoPHP/cmd/goPHP/phpError"
	"GoPHP/cmd/goPHP/runtime"
	"GoPHP/cmd/goPHP/runtime/values"
	"encoding/hex"
	"strings"
)

func registerNativeStringsFunctions(environment *Environment) {
	environment.nativeFunctions["bin2hex"] = nativeFn_bin2hex
	environment.nativeFunctions["chr"] = nativeFn_chr
	environment.nativeFunctions["lcfirst"] = nativeFn_lcfirst
	environment.nativeFunctions["quotemeta"] = nativeFn_quotemeta
	environment.nativeFunctions["str_contains"] = nativeFn_str_contains
	environment.nativeFunctions["str_ends_with"] = nativeFn_str_ends_with
	environment.nativeFunctions["str_repeat"] = nativeFn_str_repeat
	environment.nativeFunctions["str_starts_with"] = nativeFn_str_starts_with
	environment.nativeFunctions["strlen"] = nativeFn_strlen
	environment.nativeFunctions["strtolower"] = nativeFn_strtolower
	environment.nativeFunctions["strtoupper"] = nativeFn_strtoupper
	environment.nativeFunctions["ucfirst"] = nativeFn_ucfirst
}

// ------------------- MARK: bin2hex -------------------

func nativeFn_bin2hex(args []values.RuntimeValue, _ runtime.Context) (values.RuntimeValue, phpError.Error) {
	// Spec: https://www.php.net/manual/en/function.bin2hex.php

	args, err := NewFuncParamValidator("bin2hex").addParam("$string", []string{"string"}, nil).validate(args)
	if err != nil {
		return values.NewVoid(), err
	}

	input := args[0].(*values.Str).Value

	var output strings.Builder
	for i := 0; i < len(input); i++ {
		output.WriteString(strings.ToLower(hex.EncodeToString([]byte{input[i]})))
	}

	return values.NewStr(output.String()), nil
}

// ------------------- MARK: chr -------------------

func nativeFn_chr(args []values.RuntimeValue, _ runtime.Context) (values.RuntimeValue, phpError.Error) {
	// Spec: https://www.php.net/manual/en/function.chr.php

	args, err := NewFuncParamValidator("chr").addParam("$codepoint", []string{"int"}, nil).validate(args)
	if err != nil {
		return values.NewVoid(), err
	}

	codepoint := args[0].(*values.Int).Value

	// Spec: https://www.php.net/manual/en/function.chr.php
	// Values outside the valid range (0..255) will be bitwise and'ed with 255
	for codepoint < 0 {
		codepoint += 256
	}
	codepoint %= 256

	return values.NewStr(string(rune(codepoint))), nil
}

// ------------------- MARK: lcfirst -------------------

func nativeFn_lcfirst(args []values.RuntimeValue, _ runtime.Context) (values.RuntimeValue, phpError.Error) {
	// Spec: https://www.php.net/manual/en/function.lcfirst.php

	args, err := NewFuncParamValidator("lcfirst").addParam("$string", []string{"string"}, nil).validate(args)
	if err != nil {
		return values.NewVoid(), err
	}

	input := args[0].(*values.Str).Value

	if len(input) == 0 {
		return values.NewStr(""), nil
	}

	// Spec: https://www.php.net/manual/en/function.lcfirst.php
	// Returns a string with the first character of string lowercased
	// if that character is an ASCII character in the range "A" (0x41) to "Z" (0x5a).
	if input[0] >= 'A' && input[0] <= 'Z' {
		input = string(input[0]+32) + input[1:]
	}

	return values.NewStr(input), nil
}

// ------------------- MARK: quotemeta -------------------

func nativeFn_quotemeta(args []values.RuntimeValue, _ runtime.Context) (values.RuntimeValue, phpError.Error) {
	// Spec: https://www.php.net/manual/en/function.quotemeta.php

	args, err := NewFuncParamValidator("quotemeta").addParam("$string", []string{"string"}, nil).validate(args)
	if err != nil {
		return values.NewVoid(), err
	}

	input := args[0].(*values.Str).Value

	if input == "" {
		return values.NewBool(false), nil
	}

	var output strings.Builder
	for i := 0; i < len(input); i++ {
		if strings.ContainsAny(string(input[i]), `.\+*?[^($)`) {
			output.WriteByte('\\')
		}
		output.WriteByte(input[i])
	}

	return values.NewStr(output.String()), nil
}

// ------------------- MARK: str_contains -------------------

func nativeFn_str_contains(args []values.RuntimeValue, _ runtime.Context) (values.RuntimeValue, phpError.Error) {
	// Spec: https://www.php.net/manual/en/function.str-contains.php

	args, err := NewFuncParamValidator("str_contains").
		addParam("$haystack", []string{"string"}, nil).
		addParam("$needle", []string{"string"}, nil).
		validate(args)
	if err != nil {
		return values.NewVoid(), err
	}

	haystack := args[0].(*values.Str).Value
	needle := args[1].(*values.Str).Value

	return values.NewBool(strings.Contains(haystack, needle)), nil
}

// ------------------- MARK: str_ends_with -------------------

func nativeFn_str_ends_with(args []values.RuntimeValue, _ runtime.Context) (values.RuntimeValue, phpError.Error) {
	// Spec: https://www.php.net/manual/en/function.str-ends-with.php

	args, err := NewFuncParamValidator("str_ends_with").
		addParam("$haystack", []string{"string"}, nil).
		addParam("$needle", []string{"string"}, nil).
		validate(args)
	if err != nil {
		return values.NewVoid(), err
	}

	haystack := args[0].(*values.Str).Value
	needle := args[1].(*values.Str).Value

	return values.NewBool(strings.HasSuffix(haystack, needle)), nil
}

// ------------------- MARK: str_repeat -------------------

func nativeFn_str_repeat(args []values.RuntimeValue, _ runtime.Context) (values.RuntimeValue, phpError.Error) {
	// Spec: https://www.php.net/manual/en/function.str-repeat.php

	args, err := NewFuncParamValidator("str_repeat").
		addParam("$string", []string{"string"}, nil).
		addParam("$times", []string{"int"}, nil).
		validate(args)
	if err != nil {
		return values.NewVoid(), err
	}

	input := args[0].(*values.Str).Value
	times := args[1].(*values.Int).Value

	// Spec: https://www.php.net/manual/en/function.str-repeat.php
	// times has to be greater than or equal to 0.
	// If the times is set to 0, the function will return an empty string.
	if times < 0 {
		return values.NewVoid(), phpError.NewError("Uncaught ValueError: str_repeat(): Argument #2 ($times) must be greater than or equal to 0")
	}

	return values.NewStr(strings.Repeat(input, int(times))), nil
}

// ------------------- MARK: str_starts_with -------------------

func nativeFn_str_starts_with(args []values.RuntimeValue, _ runtime.Context) (values.RuntimeValue, phpError.Error) {
	// Spec: https://www.php.net/manual/en/function.str-starts-with.php

	args, err := NewFuncParamValidator("str_starts_with").
		addParam("$haystack", []string{"string"}, nil).
		addParam("$needle", []string{"string"}, nil).
		validate(args)
	if err != nil {
		return values.NewVoid(), err
	}

	haystack := args[0].(*values.Str).Value
	needle := args[1].(*values.Str).Value

	return values.NewBool(strings.HasPrefix(haystack, needle)), nil
}

// ------------------- MARK: strlen -------------------

func nativeFn_strlen(args []values.RuntimeValue, _ runtime.Context) (values.RuntimeValue, phpError.Error) {
	// Spec: https://www.php.net/manual/en/function.strlen.php

	args, err := NewFuncParamValidator("strlen").addParam("$string", []string{"string"}, nil).validate(args)
	if err != nil {
		return values.NewVoid(), err
	}

	return values.NewInt(int64(len(args[0].(*values.Str).Value))), nil
}

// ------------------- MARK: strtolower -------------------

func nativeFn_strtolower(args []values.RuntimeValue, _ runtime.Context) (values.RuntimeValue, phpError.Error) {
	// Spec: https://www.php.net/manual/en/function.strtolower.php

	args, err := NewFuncParamValidator("strtolower").addParam("$string", []string{"string"}, nil).validate(args)
	if err != nil {
		return values.NewVoid(), err
	}

	// Spec: https://www.php.net/manual/en/function.strtolower.php
	// Bytes in the range "A" (0x41) to "Z" (0x5a) will be converted to the corresponding lowercase letter by adding 32 to each byte value.

	input := args[0].(*values.Str).Value

	for i := 0; i < len(input); i++ {
		if input[i] >= 'A' && input[i] <= 'Z' {
			input = input[:i] + string(input[i]+32) + input[i+1:]
		}
	}

	return values.NewStr(input), nil
}

// ------------------- MARK: strtoupper -------------------

func nativeFn_strtoupper(args []values.RuntimeValue, _ runtime.Context) (values.RuntimeValue, phpError.Error) {
	// Spec: https://www.php.net/manual/en/function.strtoupper.php

	args, err := NewFuncParamValidator("strtoupper").addParam("$string", []string{"string"}, nil).validate(args)
	if err != nil {
		return values.NewVoid(), err
	}

	// Spec: https://www.php.net/manual/en/function.strtoupper.php
	// Bytes in the range "a" (0x61) to "z" (0x7a) will be converted to the corresponding uppercase letter by subtracting 32 from each byte value.

	input := args[0].(*values.Str).Value

	for i := 0; i < len(input); i++ {
		if input[i] >= 'a' && input[i] <= 'z' {
			input = input[:i] + string(input[i]-32) + input[i+1:]
		}
	}

	return values.NewStr(input), nil
}

// ------------------- MARK: ucfirst -------------------

func nativeFn_ucfirst(args []values.RuntimeValue, _ runtime.Context) (values.RuntimeValue, phpError.Error) {
	// Spec: https://www.php.net/manual/en/function.ucfirst.php

	args, err := NewFuncParamValidator("ucfirst").addParam("$string", []string{"string"}, nil).validate(args)
	if err != nil {
		return values.NewVoid(), err
	}

	input := args[0].(*values.Str).Value

	if len(input) == 0 {
		return values.NewStr(""), nil
	}

	// Spec: https://www.php.net/manual/en/function.ucfirst.php
	// Returns a string with the first character of string capitalized,
	// if that character is an ASCII character in the range from "a" (0x61) to "z" (0x7a).
	if input[0] >= 'a' && input[0] <= 'z' {
		input = string(input[0]-32) + input[1:]
	}

	return values.NewStr(input), nil
}

// TODO addcslashes
// TODO addslashes
// TODO chop
// TODO chunk_split
// TODO convert_uudecode
// TODO convert_uuencode
// TODO count_hars
// TODO crc32
// TODO crypt
// TODO explode
// TODO fprintf
// TODO get_html_translation_table
// TODO hebrev
// TODO hex2bin
// TODO html_entity_decode
// TODO htmlentities
// TODO htmlspecialchars
// TODO htmlspecialchars_decode
// TODO implode
// TODO join
// TODO levenshtein
// TODO localeconv
// TODO ltrim
// TODO md5
// TODO md5_file
// TODO metaphone
// TODO money_format
// TODO nl_langinfo
// TODO nl2br
// TODO number_format
// TODO ord
// TODO parse_str
// TODO printf
// TODO quoted_printable_decode
// TODO quoted_printable_encode
// TODO rtrim
// TODO setlocale
// TODO sha1
// TODO sha1_file
// TODO similar_text
// TODO soundex
// TODO sprintf
// TODO sscanf
// TODO str_decrement
// TODO str_getcsv
// TODO str_increment
// TODO str_ireplace
// TODO str_pad
// TODO str_replace
// TODO str_rot13
// TODO str_shuffle
// TODO str_split
// TODO str_word_count
// TODO strcasecmp
// TODO strchr
// TODO strcmp
// TODO strcoll
// TODO strcspn
// TODO strip_tags
// TODO stripcslashes
// TODO stripos
// TODO stripslashes
// TODO stristr
// TODO strnatcasecmp
// TODO strnatcmp
// TODO strncasecmp
// TODO strncmp
// TODO strpbrk
// TODO strpos
// TODO strrchr
// TODO strrev
// TODO strripos
// TODO strrpos
// TODO strspn
// TODO strstr
// TODO strtok
// TODO strtr
// TODO substr
// TODO substr_compare
// TODO substr_count
// TODO substr_replace
// TODO trim
// TODO ucwords
// TODO vfprintf
// TODO vprintf
// TODO vsprintf
// TODO wordwrap
// Deprecated:
// TODO convert_cyr_string
// TODO hebrevc
// TODO utf8_decode
// TODO utf8_encode
