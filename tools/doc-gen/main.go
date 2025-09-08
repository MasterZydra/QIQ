package main

import (
	"fmt"
	"os"
)

func main() {
	println("Generating documentation for ...")

	println("- std lib")
	docStdLib()

	println("- constants")
	docConst()

	println("- packages")
	docPackage()

	println("- ini directives")
	docIniDirectives()

	println("- statements and expressions")
	docStatements()
}

func writeToFile(path, content string) error {
	file, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("failed to create file %s: %s", path, err)
	}

	_, err = file.WriteString(content)
	if err != nil {
		return fmt.Errorf("failed to write to file %s: %s", path, err)
	}
	return nil
}
