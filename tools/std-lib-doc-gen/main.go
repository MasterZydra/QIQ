package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

var stdLibFunctions = map[string][]string{}

func main() {
	if err := searchDirectory("./cmd/goPHP/runtime/stdlib"); err != nil {
		fmt.Println(err)
		return
	}

	// Get and sort categories
	categories := []string{}
	for category := range stdLibFunctions {
		categories = append(categories, category)
	}
	sort.Slice(categories, func(i, j int) bool {
		return categories[i] < categories[j]
	})

	writeToFile("./doc/StdLib.md", generateMarkdown(categories))
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

func generateMarkdown(categories []string) string {
	var markdown strings.Builder

	markdown.WriteString("# StdLib Functions\n")

	for _, category := range categories {
		markdown.WriteString("\n## ")
		markdown.WriteString(category)
		markdown.WriteString("\n")

		functions := stdLibFunctions[category]
		sort.Slice(functions, func(i, j int) bool {
			return functions[i] < functions[j]
		})

		for _, function := range functions {
			markdown.WriteString("- ")
			markdown.WriteString(function)
			markdown.WriteString("\n")
		}
	}

	return markdown.String()
}

func searchDirectory(path string) error {
	files, err := os.ReadDir(path)
	if err != nil {
		return fmt.Errorf("failed reading dir %s: %s", path, err)
	}

	for _, file := range files {
		// Iterate sub directories
		if file.IsDir() {
			err := searchDirectory(filepath.Join(path, file.Name()))
			if err != nil {
				return err
			}
			continue
		}

		if !strings.HasSuffix(file.Name(), ".go") {
			continue
		}
		// Read all .go files
		err = readGoFile(filepath.Join(path, file.Name()))
		if err != nil {
			return err
		}
	}
	return nil
}

func readGoFile(path string) error {
	file, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("failed to open file %s: %s", path, err)
	}
	defer file.Close()

	category := ""
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}

		if strings.HasPrefix(line, "// Category: ") {
			category = strings.Replace(line, "// Category: ", "", 1)
			continue
		}

		if strings.HasPrefix(line, "environment.AddNativeFunction(") {
			if category == "" {
				return fmt.Errorf("category is not set before adding native functions in file %s", path)
			}

			if _, found := stdLibFunctions[category]; !found {
				stdLibFunctions[category] = []string{}
			}

			functionName := strings.Replace(line, "environment.AddNativeFunction(\"", "", 1)
			functionName = functionName[:strings.Index(functionName, "\"")]

			stdLibFunctions[category] = append(stdLibFunctions[category], functionName)
		}
	}
	return nil
}
