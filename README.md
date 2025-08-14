[![Build and Test](https://github.com/MasterZydra/QIQ/actions/workflows/go-build-and-test.yml/badge.svg)](https://github.com/MasterZydra/QIQ/actions/workflows/go-build-and-test.yml)
# QIQ

QIQ (**Q**uick **I**nterpreter for **Q**uasi-PHP)
GoPHP is an implementation of the [PHP language specification](https://phplang.org/) written in the Go programming language.

The goals of the project are:
- Deep dive into the PHP language syntax and the internals of its mode of operation
- Gain more experience in writing lexers, parser and interpreter
- *Very long-term goal*: Implement as many parts of the standard library and language features as needed to run a simple Laravel application :sweat:

**More documentation:**
- [AST](doc/AST.md)
- [Features](doc/Features.md)
- [Internal workings](doc/Internal%20workings.md)
- [Packages](doc/Packages.md)

## Usage
```
Usage of ./qiq:
  -h            Show help
  -dev          Run in developer mode.
  -stats        Show statistics.
  -S string     Run with built-in web server. <addr>:<port>
  -t string     Specify document root <docroot> for built-in web server.
  -f string     Parse and execute <file>.
```

**Parse file:**  
`cat index.php | ./qiq` or `./qiq -f index.php`

**Run web server:**  
`./qiq -S localhost:8080` - Document root is current working directory  
`./qiq -S localhost:8080 -dev` - Web server in developer mode  
`./qiq -S localhost:8080 -t /srv/www/html` - Document root is `/srv/www/html`

## Development

**Compile and run**  
`go run ./...`

**Build executable**  
`go build -o . ./...`

**Run all tests**  
`go test -v ./...`

**See test coverage**  
`go test -coverprofile=coverage.out ./...`  
`go tool cover -html=coverage.out`

## Run official PHP `phpt` test cases
There are a lot of test cases in the source repository for PHP under the folder [tests](https://github.com/php/php-src/tree/master/tests).  
In order to test the QIQ implementation against this cases the binary `qiqTester` can be used.

**Usage:**  
`./qiqTester [-v(1|2)] [-only-failed] <list of directory or phpt-file>`

**Examples:**  
`./qiqTester php-src/tests`  
`./qiqTester -v2 php-src/tests/basic/001.phpt`

## Used resources
For some part of this project, the following resources were used as a guide, inspiration, or concept:
- [PHP Language Specification](https://phplang.org/)
- YouTube playlist [Build a Custom Scripting Language In Typescript](https://www.youtube.com/playlist?list=PL_2VhOvlMk4UHGqYCLWc6GO8FaPl8fQTh) by [tylerlaceby](https://www.youtube.com/@tylerlaceby)
- Book [Crafting Interpreters](https://craftinginterpreters.com/) by Robert Nystorm
- Book [Writing An Interpreter In Go](https://interpreterbook.com/) by Thorsten Ball
