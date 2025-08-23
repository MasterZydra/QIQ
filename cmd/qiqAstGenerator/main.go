package main

import (
	"QIQ/cmd/qiqAstGenerator/astGenerator"
	"fmt"
	"os"
)

func main() {
	args := os.Args[1:]

	if len(args) != 1 {
		fmt.Println("Usage: qiqAstGenerator [path to php file]")
		os.Exit(1)
	}

	filename := args[0]
	content, err := os.ReadFile(filename)
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}

	output, err := astGenerator.NewAstGenerator().Process(string(content), filename)
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}

	ouputFileName := filename + ".ast"
	outputFile, err := os.Create(ouputFileName)
	if err != nil {
		fmt.Printf("failed to create file %s: %s\n", ouputFileName, err)
		os.Exit(1)
	}

	_, err = outputFile.WriteString(output)
	if err != nil {
		fmt.Printf("failed to write to file %s: %s\n", ouputFileName, err)
		os.Exit(1)
	}
}
