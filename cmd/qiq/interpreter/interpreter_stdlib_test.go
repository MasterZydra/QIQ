package interpreter

import (
	"QIQ/cmd/qiq/phpError"
	"testing"
)

// -------------------------------------- math -------------------------------------- MARK: math

func TestLibMath(t *testing.T) {
	// abs
	testInputOutput(t, `<?php var_dump(abs(-4.2));`, "float(4.2)\n")
	testInputOutput(t, `<?php var_dump(abs(5));`, "int(5)\n")
	testInputOutput(t, `<?php var_dump(abs(-5));`, "int(5)\n")

	// acos
	testInputOutput(t, `<?php var_dump(acos(1.0));`, "float(0)\n")
	testInputOutput(t, `<?php var_dump(acos(0.5)/M_PI*180);`, "float(60)\n")

	// acosh
	testInputOutput(t, `<?php var_dump(acosh(1.0));`, "float(0)\n")

	// asin
	testInputOutput(t, `<?php var_dump(asin(0.0));`, "float(0)\n")

	// asinh
	testInputOutput(t, `<?php var_dump(asinh(0.0));`, "float(0)\n")

	// pi
	testInputOutput(t, `<?php var_dump(M_PI === pi());`, "bool(true)\n")
}

// -------------------------------------- constant -------------------------------------- MARK: constant

func TestLibConstant(t *testing.T) {
	testInputOutput(t, `<?php var_dump(constant('E_ALL'));`, "int(32767)\n")
	testForError(t, `<?php constant('NOT_DEFINED_CONSTANT');`, phpError.NewError("Undefined constant \"NOT_DEFINED_CONSTANT\""))
	// TODO Add test cases for user defined constants
}

// -------------------------------------- defined -------------------------------------- MARK: defined

func TestLibDefined(t *testing.T) {
	testInputOutput(t, `<?php var_dump(defined('PHP_VERSION'));`, "bool(true)\n")
	testInputOutput(t, `<?php var_dump(defined('NOT_DEFINED_CONSTANT'));`, "bool(false)\n")
	// TODO Add test cases for user defined constants
}

// -------------------------------------- ob_ functions -------------------------------------- MARK: ob_ functions

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

// -------------------------------------- date -------------------------------------- MARK: date

func TestLibDate(t *testing.T) {
	// checkdate
	testInputOutput(t, `<?php var_dump(checkdate(12, 31, 2000));`, "bool(true)\n")
	testInputOutput(t, `<?php var_dump(checkdate(2, 29, 2001));`, "bool(false)\n")

	// date
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

	// getdate
	testInputOutput(t, `<?php print_r(getdate(1722707036));`,
		"Array\n(\n    [seconds] => 56\n    [minutes] => 43\n    [hours] => 17\n    [mday] => 3\n    [wday] => 6\n    [mon] => 8\n"+
			"    [year] => 2024\n    [yday] => 215\n    [weekday] => Saturday\n    [month] => August\n    [0] => 1722707036\n)\n",
	)
}

// -------------------------------------- strings -------------------------------------- MARK: strings

func TestLibStrings(t *testing.T) {
	// bin2hex
	testInputOutput(t, `<?php var_dump(bin2hex('Hello world!'));`, "string(24) \"48656c6c6f20776f726c6421\"\n")
	testInputOutput(t, `<?php var_dump(bin2hex('Äàßê'));`, "string(16) \"c384c3a0c39fc3aa\"\n")
	testInputOutput(t, `<?php var_dump(bin2hex(''));`, "string(0) \"\"\n")

	// chr
	testInputOutput(t, `<?php var_dump(chr(60));`, "string(1) \"<\"\n")
	testInputOutput(t, `<?php var_dump(chr(60-256));`, "string(1) \"<\"\n")
	testInputOutput(t, `<?php var_dump(chr(60+256));`, "string(1) \"<\"\n")

	// lcfirst
	testInputOutput(t, `<?php var_dump(lcfirst('ABC'));`, "string(3) \"aBC\"\n")
	testInputOutput(t, `<?php var_dump(lcfirst('Abc'));`, "string(3) \"abc\"\n")
	testInputOutput(t, `<?php var_dump(lcfirst('abc'));`, "string(3) \"abc\"\n")
	testInputOutput(t, `<?php var_dump(lcfirst(''));`, "string(0) \"\"\n")

	// quotemeta
	testInputOutput(t, `<?php var_dump(quotemeta('. \ + * ? [ ^ ] ( $ )'));`, `string(31) "\. \\ \+ \* \? \[ \^ ] \( \$ \)"`+"\n")
	testInputOutput(t, `<?php var_dump(quotemeta('Hello. (can you hear me?)'));`, `string(29) "Hello\. \(can you hear me\?\)"`+"\n")
	testInputOutput(t, `<?php var_dump(quotemeta(''));`, "bool(false)\n")

	// str_contains
	testInputOutput(t, `<?php var_dump(str_contains('abc', ''));`, "bool(true)\n")
	testInputOutput(t, `<?php var_dump(str_contains('The lazy fox', 'lazy'));`, "bool(true)\n")
	testInputOutput(t, `<?php var_dump(str_contains('The lazy fox', 'Lazy'));`, "bool(false)\n")

	// str_ends_with
	testInputOutput(t, `<?php var_dump(str_ends_with('abc', ''));`, "bool(true)\n")
	testInputOutput(t, `<?php var_dump(str_ends_with('The lazy fox', 'fox'));`, "bool(true)\n")
	testInputOutput(t, `<?php var_dump(str_ends_with('The lazy fox', 'Fox'));`, "bool(false)\n")

	// str_repeat
	testInputOutput(t, `<?php var_dump(str_repeat('abc', 0));`, "string(0) \"\"\n")
	testInputOutput(t, `<?php var_dump(str_repeat('abc', 1));`, "string(3) \"abc\"\n")
	testInputOutput(t, `<?php var_dump(str_repeat('abc', 2));`, "string(6) \"abcabc\"\n")
	testForError(t, `<?php var_dump(str_repeat('abc', -1));`, phpError.NewError("Uncaught ValueError: str_repeat(): Argument #2 ($times) must be greater than or equal to 0"))

	// str_starts_with
	testInputOutput(t, `<?php var_dump(str_starts_with('abc', ''));`, "bool(true)\n")
	testInputOutput(t, `<?php var_dump(str_starts_with('The lazy fox', 'The'));`, "bool(true)\n")
	testInputOutput(t, `<?php var_dump(str_starts_with('The lazy fox', 'the'));`, "bool(false)\n")

	// strlen
	testInputOutput(t, `<?php var_dump(strlen('abcdef'));`, "int(6)\n")
	testInputOutput(t, `<?php var_dump(strlen(' ab cd '));`, "int(7)\n")
	testInputOutput(t, `<?php var_dump(strlen(' äb ćd '));`, "int(9)\n")

	// strtolower
	testInputOutput(t, `<?php var_dump(strtolower('Mary Had A Little Lamb and She LOVED It So'));`, "string(42) \"mary had a little lamb and she loved it so\"\n")
	testInputOutput(t, `<?php var_dump(strtolower(''));`, "string(0) \"\"\n")
	testInputOutput(t, `<?php var_dump(strtolower('AÄOÖUÜSß'));`, "string(12) \"aÄoÖuÜsß\"\n")

	// strtoupper
	testInputOutput(t, `<?php var_dump(strtoupper('Mary Had A Little Lamb and She LOVED It So'));`, "string(42) \"MARY HAD A LITTLE LAMB AND SHE LOVED IT SO\"\n")
	testInputOutput(t, `<?php var_dump(strtoupper(''));`, "string(0) \"\"\n")
	testInputOutput(t, `<?php var_dump(strtoupper('aäoöuüsß'));`, "string(12) \"AäOöUüSß\"\n")

	// ucfirst
	testInputOutput(t, `<?php var_dump(ucfirst('ABC'));`, "string(3) \"ABC\"\n")
	testInputOutput(t, `<?php var_dump(ucfirst('abc'));`, "string(3) \"Abc\"\n")
	testInputOutput(t, `<?php var_dump(ucfirst('Abc'));`, "string(3) \"Abc\"\n")
	testInputOutput(t, `<?php var_dump(ucfirst(''));`, "string(0) \"\"\n")
}

// -------------------------------------- get_debug_type -------------------------------------- MARK: get_debug_type

func TestLibGetDebugType(t *testing.T) {
	testInputOutput(t, `<?php echo get_debug_type(false);`, "bool")
	testInputOutput(t, `<?php echo get_debug_type(true);`, "bool")
	testInputOutput(t, `<?php echo get_debug_type(0);`, "int")
	testInputOutput(t, `<?php echo get_debug_type(-1);`, "int")
	testInputOutput(t, `<?php echo get_debug_type(42);`, "int")
	testInputOutput(t, `<?php echo get_debug_type(0.0);`, "float")
	testInputOutput(t, `<?php echo get_debug_type(-1.5);`, "float")
	testInputOutput(t, `<?php echo get_debug_type(42.5);`, "float")
	testInputOutput(t, `<?php echo get_debug_type("");`, "string")
	testInputOutput(t, `<?php echo get_debug_type("abc");`, "string")
	testInputOutput(t, `<?php echo get_debug_type([]);`, "array")
	testInputOutput(t, `<?php echo get_debug_type([42]);`, "array")
	testInputOutput(t, `<?php echo get_debug_type(null);`, "null")
}

// -------------------------------------- gettype -------------------------------------- MARK: gettype

func TestLibGettype(t *testing.T) {
	testInputOutput(t, `<?php echo gettype(false);`, "boolean")
	testInputOutput(t, `<?php echo gettype(true);`, "boolean")
	testInputOutput(t, `<?php echo gettype(0);`, "integer")
	testInputOutput(t, `<?php echo gettype(-1);`, "integer")
	testInputOutput(t, `<?php echo gettype(42);`, "integer")
	testInputOutput(t, `<?php echo gettype(0.0);`, "double")
	testInputOutput(t, `<?php echo gettype(-1.5);`, "double")
	testInputOutput(t, `<?php echo gettype(42.5);`, "double")
	testInputOutput(t, `<?php echo gettype("");`, "string")
	testInputOutput(t, `<?php echo gettype("abc");`, "string")
	testInputOutput(t, `<?php echo gettype([]);`, "array")
	testInputOutput(t, `<?php echo gettype([42]);`, "array")
	testInputOutput(t, `<?php echo gettype(null);`, "NULL")
}

// -------------------------------------- is_XXX -------------------------------------- MARK: is_XXX

func TestLibIsType(t *testing.T) {
	// is_array
	testInputOutput(t, `<?php $a = [true]; var_dump(is_array($a));`, "bool(true)\n")
	testInputOutput(t, `<?php $a = 0; var_dump(is_array($a));`, "bool(false)\n")

	// is_bool
	testInputOutput(t, `<?php $a = true; var_dump(is_bool($a));`, "bool(true)\n")
	testInputOutput(t, `<?php $a = 0; var_dump(is_bool($a));`, "bool(false)\n")

	// is_float
	testInputOutput(t, `<?php $a = 42.0; var_dump(is_float($a));`, "bool(true)\n")
	testInputOutput(t, `<?php $a = 0; var_dump(is_float($a));`, "bool(false)\n")

	// is_int
	testInputOutput(t, `<?php $a = 42; var_dump(is_int($a));`, "bool(true)\n")
	testInputOutput(t, `<?php $a = "42"; var_dump(is_int($a));`, "bool(false)\n")

	// is_null
	testInputOutput(t, `<?php $a = null; var_dump(is_null($a));`, "bool(true)\n")
	testInputOutput(t, `<?php $a = 42; var_dump(is_null($a));`, "bool(false)\n")

	// is_scalar
	testInputOutput(t, `<?php $a = true; var_dump(is_scalar($a));`, "bool(true)\n")
	testInputOutput(t, `<?php $a = false; var_dump(is_scalar($a));`, "bool(true)\n")
	testInputOutput(t, `<?php $a = 42; var_dump(is_scalar($a));`, "bool(true)\n")
	testInputOutput(t, `<?php $a = 3.5; var_dump(is_scalar($a));`, "bool(true)\n")
	testInputOutput(t, `<?php $a = "abc"; var_dump(is_scalar($a));`, "bool(true)\n")
	testInputOutput(t, `<?php $a = null; var_dump(is_scalar($a));`, "bool(false)\n")
	testInputOutput(t, `<?php $a = []; var_dump(is_scalar($a));`, "bool(false)\n")

	// is_string
	testInputOutput(t, `<?php $a = " "; var_dump(is_string($a));`, "bool(true)\n")
	testInputOutput(t, `<?php $a = 42; var_dump(is_string($a));`, "bool(false)\n")

	// is_object
	testInputOutput(t, `<?php $a = null; var_dump(is_object($a));`, "bool(false)\n")
	testInputOutput(t, `<?php $a = []; var_dump(is_object($a));`, "bool(false)\n")
	testInputOutput(t, `<?php class C {} $a = new C; var_dump(is_object($a));`, "bool(true)\n")
}

// -------------------------------------- print_r -------------------------------------- MARK: print_r

func TestLibPrintR(t *testing.T) {
	testInputOutput(t, `<?php print_r(3.5);`, "3.5")
	testInputOutput(t, `<?php print_r(42);`, "42")
	testInputOutput(t, `<?php print_r("abc");`, "abc")
	testInputOutput(t, `<?php print_r(true);`, "1")
	testInputOutput(t, `<?php print_r(false);`, "")
	testInputOutput(t, `<?php print_r(null);`, "")
	testInputOutput(t, `<?php print_r([]);`, "Array\n(\n)\n")
	testInputOutput(t, `<?php print_r([1,2]);`, "Array\n(\n    [0] => 1\n    [1] => 2\n)\n")
	testInputOutput(t, `<?php print_r([1, [1]]);`, "Array\n(\n    [0] => 1\n    [1] => Array\n        (\n            [0] => 1\n        )\n\n)\n")
}

// -------------------------------------- var_dump -------------------------------------- MARK: var_dump

func TestLibVarDump(t *testing.T) {
	testInputOutput(t, `<?php var_dump(3.5);`, "float(3.5)\n")
	testInputOutput(t, `<?php var_dump(3.5, 42, true, false, null);`, "float(3.5)\nint(42)\nbool(true)\nbool(false)\nNULL\n")
	testInputOutput(t, `<?php var_dump([]);`, "array(0) {\n}\n")
	testInputOutput(t, `<?php var_dump([1,2]);`, "array(2) {\n  [0]=>\n  int(1)\n  [1]=>\n  int(2)\n}\n")
	testInputOutput(t, `<?php var_dump([1, [1]]);`, "array(2) {\n  [0]=>\n  int(1)\n  [1]=>\n  array(1) {\n    [0]=>\n    int(1)\n  }\n}\n")
}

// -------------------------------------- var_export -------------------------------------- MARK: var_export

func TestLibVarExport(t *testing.T) {
	testInputOutput(t, `<?php var_export(3.5);`, "3.5")
	testInputOutput(t, `<?php var_export(42);`, "42")
	testInputOutput(t, `<?php var_export("abc");`, "'abc'")
	testInputOutput(t, `<?php var_export(true);`, "true")
	testInputOutput(t, `<?php var_export(false);`, "false")
	testInputOutput(t, `<?php var_export(null);`, "NULL")
	testInputOutput(t, `<?php var_export([]);`, "array (\n)")
	testInputOutput(t, `<?php var_export([1,2]);`, "array (\n  0 => 1,\n  1 => 2,\n)")
	testInputOutput(t, `<?php var_export([1, [1]]);`, "array (\n  0 => 1,\n  1 => \n  array (\n    0 => 1,\n  ),\n)")
}

// -------------------------------------- option_info -------------------------------------- MARK: option_info

func TestLibOptionInfo(t *testing.T) {
	// ini_get
	testInputOutput(t, `<?php var_dump(ini_get('none_existing'));`, "bool(false)\n")
	testInputOutput(t, `<?php var_dump(ini_get('variables_order'));`, "string(5) \"EGPCS\"\n")
	testInputOutput(t, `<?php var_dump(ini_get('error_reporting'));`, "string(5) \"32767\"\n")

	// ini_set
	testInputOutput(t, `<?php var_dump(ini_set('none_existing', true));`, "bool(false)\n")
	testInputOutput(t, `<?php var_dump(ini_set('variables_order', true));`, "bool(false)\n")
	testInputOutput(t, `<?php var_dump(ini_set('error_reporting', E_ERROR)); var_dump(ini_set('error_reporting', E_ERROR));`, "string(5) \"32767\"\nstring(1) \"1\"\n")
}

// -------------------------------------- error_reporting -------------------------------------- MARK: error_reporting

func TestLibErrorReporting(t *testing.T) {
	// Spec: https://www.php.net/manual/en/function.error-reporting.php - Example #1

	// Turn off all error reporting
	testInputOutput(t, `<?php error_reporting(0); echo error_reporting();`, "0")
	// Report simple running errors
	testInputOutput(t, `<?php error_reporting(E_ERROR | E_WARNING | E_PARSE); echo error_reporting();`, "7")
	// Reporting E_NOTICE can be good too (to report uninitialized variables or catch variable name misspellings ...)
	testInputOutput(t, `<?php error_reporting(E_ERROR | E_WARNING | E_PARSE | E_NOTICE); echo error_reporting();`, "15")
	// Report all errors except E_NOTICE
	testInputOutput(t, `<?php error_reporting(E_ALL & ~E_NOTICE); echo error_reporting();`, "32759")
	// Report all PHP errors
	testInputOutput(t, `<?php error_reporting(E_ALL); echo error_reporting();`, "32767")
	// Report all PHP errors
	testInputOutput(t, `<?php error_reporting(-1); echo error_reporting();`, "32767")
}

// -------------------------------------- classes_object -------------------------------------- MARK: classesobject

func TestLibClassesObject(t *testing.T) {
	// class_exists
	testInputOutput(t, `<?php class C {} var_dump(class_exists('C'));`, "bool(true)\n")
	testInputOutput(t, `<?php class C {} var_dump(class_exists('c'));`, "bool(true)\n")
	testInputOutput(t, `<?php var_dump(class_exists('c'));`, "bool(false)\n")

	// get_class
	testInputOutput(t, `<?php class Foo {} $bar = new Foo; var_dump(get_class($bar));`, "string(3) \"Foo\"\n")
	testInputOutput(t, `<?php namespace Space; class Foo {} $bar = new Foo; var_dump(get_class($bar));`, "string(9) \"Space\\Foo\"\n")

	// get_parent_class
	testInputOutput(t, `<?php class Dad {} class Child extends Dad {} $c = new Child; var_dump(get_parent_class($c));`, "string(3) \"Dad\"\n")
	testInputOutput(t, `<?php namespace Space; class Dad {} class Child extends Dad {} $c = new Child; var_dump(get_parent_class($c));`, "string(9) \"Space\\Dad\"\n")

	// is_subclass_of
	testInputOutput(t, `<?php class Dad {} class Child extends Dad {} $c = new Child; var_dump(is_subclass_of($c, 'Dad'));`, "bool(true)\n")
	testInputOutput(t, `<?php class Dad {} class Child extends Dad {} $c = new Child; var_dump(is_subclass_of($c, 'dad'));`, "bool(true)\n")
	testInputOutput(t, `<?php class Dad {} class Child extends Dad {} var_dump(is_subclass_of('Child', 'dad'));`, "bool(true)\n")
	testInputOutput(t, `<?php class Dad {} class Child extends Dad {} var_dump(is_subclass_of('someChild', 'Dad'));`, "bool(false)\n")
	testInputOutput(t, `<?php namespace Space; class Dad {} class Child extends Dad {} var_dump(is_subclass_of('someChild', 'Dad'));`, "bool(false)\n")
	testInputOutput(t, `<?php namespace Space; class Dad {} class Child extends Dad {} $c = new Child; var_dump(is_subclass_of('Space\Child', 'space\dad'));`, "bool(true)\n")
}
