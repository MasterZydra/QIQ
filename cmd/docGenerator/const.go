package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strings"
)

var constants = map[string][]string{}

func docConst() {
	if err := searchDirectoryConstants("./cmd/qiq"); err != nil {
		fmt.Println(err)
		return
	}

	// Get and sort categories
	categories := []string{}
	for category := range constants {
		categories = append(categories, category)
	}
	slices.Sort(categories)

	writeToFile("./doc/Constants.md", generateMarkdownConstants(categories))
}

func generateMarkdownConstants(categories []string) string {
	var markdown strings.Builder

	markdown.WriteString("# Constants\n")

	for _, category := range categories {
		markdown.WriteString("\n## ")
		markdown.WriteString(category)
		markdown.WriteString("\n")

		functions := constants[category]
		slices.Sort(functions)

		for _, function := range functions {
			markdown.WriteString("- ")
			markdown.WriteString(function)
			markdown.WriteString("\n")
		}
	}

	return markdown.String()
}

func searchDirectoryConstants(path string) error {
	files, err := os.ReadDir(path)
	if err != nil {
		return fmt.Errorf("failed reading dir %s: %s", path, err)
	}

	for _, file := range files {
		// Iterate sub directories
		if file.IsDir() {
			err := searchDirectoryConstants(filepath.Join(path, file.Name()))
			if err != nil {
				return err
			}
			continue
		}

		if !strings.HasSuffix(file.Name(), ".go") {
			continue
		}
		// Read all .go files
		err = readGoFileConstants(filepath.Join(path, file.Name()))
		if err != nil {
			return err
		}
	}
	return nil
}

func readGoFileConstants(path string) error {
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
		if strings.HasPrefix(line, "// Const Category: ") {
			category = strings.Replace(line, "// Const Category: ", "", 1)
			continue
		}

		if strings.HasPrefix(line, "environment.AddPredefinedConstant(") {
			if category == "" {
				return fmt.Errorf("category is not set before adding predefined constants in file %s", path)
			}

			if _, found := constants[category]; !found {
				constants[category] = []string{}
			}

			functionName := strings.Replace(line, `environment.AddPredefinedConstant("`, "", 1)
			functionName = functionName[:strings.Index(functionName, `"`)]

			constants[category] = append(constants[category], functionName)
		}
	}
	return nil
}
