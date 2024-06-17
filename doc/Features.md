# Supported syntax and features

## Expressions and statements
- [x] short open tag: `<? 1 + 2;`
- [x] declare and access variable: `$var = "abc"; echo $abc;`
- [x] declare and access constant: `const PI = 3.141; echo PI;`
- [x] simple assignment: `$var = 42;`
- [x] compound assignment: `$var += 42;`
- [x] cast expression: `(int)$a;(string)$a;`
- [x] conditional expression: `$var ? $a : "b";`
- [x] coalesce expression: `$var ?? "b";`
- [x] equality expression: `$var === 42;`
- [x] relational expression: `$var >= 42;`
- [x] additive expression: `$var + 42; $var - 42; "a" . "b";`
- [x] multiplicative expression: `$var * 42; $var / 42; $var % 42;`
- [x] logical and expression: `$var && 8`;
- [x] logical inc or expression: `$var || 8`;
- [x] bitwise exc or expression: `$var ^ 8;`
- [x] bitwise inc or expression: `$var | 8;`
- [x] bitwise and expression: `$var & 8`;
- [x] shift expression: `$var << 8`;
- [x] exponentiation expression: `$var ** 42;`
- [x] unary expression: `-1; +1; ~1;`
- [x] prefix (in/de)crease expression: `++$var; --$var;`
- [x] postfix (in/de)crease expression: `$var++; $var--;`
- [x] logical not expression: `!$var;`
- [x] parenthesized expression: `(1 + 2) * 3;`
- [x] subscript expression: `$a[1];`
- [x] variable substitution: `echo "{$a}";`
- [x] if statement: `if (true) { ...} elseif (false) { ... } else { ... }`
- [x] function definition: `function func1($param1) {...}`
- [x] return statement: `return 42;`
- [x] require(_once), include(_once): `require 'lib.php';`

## Data types
- [x] array
- [x] bool
- [x] float
- [x] int
- [x] null
- [x] string

## Predefined variables
- [x] $_GET
- [x] $_POST 

## Standard library
- [x] array_key_exists
- [x] boolval
- [x] die
- [x] empty
- [x] error_reporting
- [x] exit
- [x] floatval
- [x] getenv
- [x] gettype
- [x] ini_get
- [x] intval
- [x] is_null
- [x] is_scalar
- [x] isset
- [x] key_exists
- [x] strlen
- [x] strval
- [x] unset
- [x] var_dump