package main

import (
	"GoPHP/cmd/goPHP/interpreter"
	"bufio"
	"fmt"
	"os"
)

func main() {
	// Read input
	fileContent := ""

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		if fileContent != "" || scanner.Text() == "" {
			fileContent += "\n"
		}
		fileContent += scanner.Text()
	}
	if err := scanner.Err(); err != nil {
		fmt.Println("Error:", err)
	}

	// lexer := lexer.NewLexer()
	// tokens, err := lexer.Tokenize(fileContent)
	// if err != nil {
	// 	fmt.Println("Error:", err)
	// }
	// fmt.Printf("Tokens:   %s\n", tokens)

	// parser := parser.NewParser()
	// program, err := parser.ProduceAST(fileContent)
	// if err != nil {
	// 	fmt.Println("Error:", err)
	// }

	interpreter := interpreter.NewInterpreter(interpreter.NewDevConfig(), &interpreter.Request{}, "")
	result, err := interpreter.Process(fileContent)
	fmt.Println(result)
	if err != nil {
		fmt.Println("Error:", err)
	}
	os.Exit(interpreter.GetExitCode())
}
