# Development

## Compile and run
```bash
go run ./...
```

## Build executable
```bash
go build -o . ./...
```

## Run all tests
```bash
go test -v ./...
```

**See test coverage**
```bash
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

## Build docker image
```bash
docker build -t qiq:dev .
```
