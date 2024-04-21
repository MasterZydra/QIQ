# GoPHP

GoPHP is an implementation of the [PHP language specification](https://phplang.org/) written in the Go programming language.

The goals of the project are:
- Deep dive into the PHP language syntax and the internals of its mode of operation
- Gain more experience in writing lexers, parser and interpreter

Usage: `cat index.php | ./goPHP`

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
