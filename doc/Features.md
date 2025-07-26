# Supported syntax and features

[StdLib Functions](./StdLib.md)  
[Constants](./Constants.md)  
[Ini directives](./IniDirectives.md)

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
- class declaration: `class MyClass extends ParentC implements I, J {}`
- object creation expression: `new myClass;`
- member access expression: `$obj->member`

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

## Magic constants
- \_\_CLASS\_\_
- \_\_DIR\_\_
- \_\_FILE\_\_
- \_\_FUNCTION\_\_
- \_\_LINE\_\_
- \_\_METHOD\_\_

## Intrinsics
- die
- empty
- eval
- exit
- isset
- unset
