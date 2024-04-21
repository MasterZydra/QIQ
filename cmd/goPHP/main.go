package main

import (
	"GoPHP/cmd/goPHP/lexer"
	"bufio"
	"fmt"
	"os"
)

func main() {
	// Read input
	fileContent := ""

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		if fileContent != "" {
			fileContent += "\n"
		}
		fileContent += scanner.Text()
	}
	if err := scanner.Err(); err != nil {
		fmt.Println("Error:", err)
	}

	lexer := lexer.NewLexer()
	tokens, err := lexer.Tokenize(fileContent)
	if err != nil {
		fmt.Println("Error:", err)
	}
	fmt.Printf("Tokens:   %s\n", tokens)
}
