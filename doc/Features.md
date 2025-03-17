# Supported syntax and features

[StdLib Functions](./StdLib.md)

## Ini directives
- arg_separator.input
- arg_separator.output
- default_charset
- error_reporting
- input_encoding
- internal_encoding
- max_input_nesting_level
- open_basedir
- output_encoding
- register_argc_argv
- short_open_tag
- variables_order

## Expressions and statements
- echo statement: `echo "abc", 123, true;`
- print statement: `print "abc";`
- short open tag: `<? 1 + 2;`
- short echo statement: `<?= "123";`
- declare and access variable: `$var = "abc"; echo $abc;`
- declare and access constant: `const PI = 3.141; echo PI;`
- simple assignment: `$var = 42;`
- compound assignment: `$var += 42;`
- cast expression: `(int)$a;(string)$a;`
- conditional expression: `$var ? $a : "b";`
- coalesce expression: `$var ?? "b";`
- equality expression: `$var === 42;`
- relational expression: `$var >= 42;`
- additive expression: `$var + 42; $var - 42; "a" . "b";`
- multiplicative expression: `$var * 42; $var / 42; $var % 42;`
- logical and expression: `$var && 8;`
- logical and expression 2: `$var and 8;`
- logical exc or expression: `$var xor 8;`
- logical inc or expression: `$var || 8;`
- logical inc or expression 2: `$var or 8;`
- bitwise exc or expression: `$var ^ 8;`
- bitwise inc or expression: `$var | 8;`
- bitwise and expression: `$var & 8;`
- shift expression: `$var << 8;`
- exponentiation expression: `$var ** 42;`
- unary expression: `-1; +1; ~1;`
- prefix (in/de)crease expression: `++$var; --$var;`
- postfix (in/de)crease expression: `$var++; $var--;`
- logical not expression: `!$var;`
- parenthesized expression: `(1 + 2) * 3;`
- subscript expression: `$a[1];`
- variable substitution: `echo "{$a}";`
- if statement: `if (true) { ... } elseif (false) { ... } else { ... }`
- for statement: `for (...; ...; ...) { ... }`
- while statement: `while (true) { ... }`
- do statement: `do { ... } while (true);`
- function definition: `function func1($param1) { ... }`
- break statement: `break 1;`
- continue statement: `continue (2);`
- return statement: `return 42;`
- require(_once), include(_once): `require 'lib.php';`
- error control expression: `@fn();`
- global declaration: `global $var;`

## Data types
- array
- bool
- float (including numeric literal separator)
- int  (including numeric literal separator)
- null
- string

## Predefined variables
- $_ENV
- $_COOKIE
- $_FILES
- $_GET
- $_POST
- $_REQUEST
- $_SERVER

## Predefined constants
- DIRECTORY_SEPARATOR
- FALSE
- NULL
- PHP_EOL
- PHP_EXTRA_VERSION
- PHP_INT_MAX
- PHP_INT_MIN
- PHP_INT_SIZE
- PHP_MAJOR_VERSION
- PHP_MINOR_VERSION
- PHP_OS
- PHP_OS_FAMILY
- PHP_RELEASE_VERSION
- PHP_VERSION
- PHP_VERSION_ID
- TRUE

### Error Handling
- E_ALL
- E_COMPILE_ERROR
- E_COMPILE_WARNING
- E_CORE_ERROR
- E_CORE_WARNING
- E_DEPRECATED
- E_ERROR
- E_NOTICE
- E_PARSE
- E_RECOVERABLE_ERROR
- E_STRICT
- E_USER_DEPRECATED
- E_USER_ERROR
- E_USER_NOTICE
- E_USER_WARNING
- E_WARNING

### INI mode constants
- INI_ALL
- INI_PERDIR
- INI_SYSTEM
- INI_USER

### Mathematical Constants
- M_1_PI
- M_2_PI
- M_2_SQRTPI
- M_E
- M_EULER
- M_LN10
- M_LN2
- M_LNPI
- M_LOG10E
- M_LOG2E
- M_PI
- M_PI_2
- M_PI_4
- M_SQRT1_2
- M_SQRT2
- M_SQRT3
- M_SQRTPI

### Math - Rounding constants
- PHP_ROUND_HALF_DOWN
- PHP_ROUND_HALF_EVEN
- PHP_ROUND_HALF_ODD
- PHP_ROUND_HALF_UP

## Magic constants
- \_\_DIR\_\_
- \_\_FILE\_\_
- \_\_FUNCTION\_\_
- \_\_LINE\_\_

## Intrinsics
- die
- empty
- eval
- exit
- isset
- unset
