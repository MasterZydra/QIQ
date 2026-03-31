package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strings"
)

var statements = []string{}
var expressions = []string{}
var intrinsics = []string{}

func docStatements() {
	if err := searchDirectoryStatements("./cmd/qiq"); err != nil {
		fmt.Println(err)
		return
	}

	writeToFile("./doc/StatementsAndExpressions.md", generateMarkdownStatements())
}

func generateMarkdownStatements() string {
	var markdown strings.Builder

	slices.Sort(statements)
	markdown.WriteString("# Statements\n")
	for _, statement := range statements {
		markdown.WriteString("- ")
		markdown.WriteString(statement)
		markdown.WriteString("\n")
	}
	markdown.WriteString("\n")

	slices.Sort(expressions)
	markdown.WriteString("# Expressions\n")
	for _, expression := range expressions {
		markdown.WriteString("- ")
		markdown.WriteString(expression)
		markdown.WriteString("\n")
	}
	markdown.WriteString("\n")

	slices.Sort(intrinsics)
	markdown.WriteString("# Intrinsics\n")
	for _, intrinsic := range intrinsics {
		markdown.WriteString("- ")
		markdown.WriteString(intrinsic)
		markdown.WriteString("\n")
	}

	return markdown.String()
}

func searchDirectoryStatements(path string) error {
	files, err := os.ReadDir(path)
	if err != nil {
		return fmt.Errorf("failed reading dir %s: %s", path, err)
	}

	for _, file := range files {
		// Iterate sub directories
		if file.IsDir() {
			err := searchDirectoryStatements(filepath.Join(path, file.Name()))
			if err != nil {
				return err
			}
			continue
		}

		if !strings.HasSuffix(file.Name(), ".go") {
			continue
		}
		// Read all .go files
		err = readGoFileStatements(filepath.Join(path, file.Name()))
		if err != nil {
			return err
		}
	}
	return nil
}

func readGoFileStatements(path string) error {
	file, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("failed to open file %s: %s", path, err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		if strings.HasPrefix(line, "// Supported expression: ") {
			expressions = append(expressions, strings.Replace(line, "// Supported expression: ", "", 1))
			continue
		}

		if strings.HasPrefix(line, "// Supported statement: ") {
			statements = append(statements, strings.Replace(line, "// Supported statement: ", "", 1))
			continue
		}

		if strings.HasPrefix(line, "// Supported intrinsic: ") {
			intrinsics = append(intrinsics, strings.Replace(line, "// Supported intrinsic: ", "", 1))
			continue
		}
	}
	return nil
}
