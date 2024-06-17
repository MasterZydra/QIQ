[![Go Build and Test](https://github.com/MasterZydra/GoPHP/actions/workflows/go-build-and-test.yml/badge.svg)](https://github.com/MasterZydra/GoPHP/actions/workflows/go-build-and-test.yml)
# GoPHP

GoPHP is an implementation of the [PHP language specification](https://phplang.org/) written in the Go programming language.

The goals of the project are:
- Deep dive into the PHP language syntax and the internals of its mode of operation
- Gain more experience in writing lexers, parser and interpreter

## Usage
```
Usage of ./goPHP:
  -h            Show help
  -dev          Run in developer mode.
  -S string     Run with built-in web server. <addr>:<port>
  -t string     Specify document root <docroot> for built-in web server.
  -f string     Parse and execute <file>.
```

**Parse file:**  
`cat index.php | ./goPHP` or `./goPHP -f index.php`

**Run web server:**  
`./goPHP -S localhost:8080` - Document root is current working directory  
`./goPHP -S localhost:8080 -dev` - Web server in developer mode  
`./goPHP -S localhost:8080 -t /srv/www/html` - Document root is `/srv/www/html`

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
