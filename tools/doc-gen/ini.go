package main

import (
	"bufio"
	"fmt"
	"os"
	"slices"
	"strings"
)

var iniDirectives = map[string][]string{}

func docIniDirectives() {
	err := readGoFileIni("./cmd/qiq/ini/ini_directives.go")
	if err != nil {
		fmt.Println(err)
		return
	}

	// Get and sort directives
	directives := []string{}
	for directive := range iniDirectives {
		directives = append(directives, directive)
	}
	slices.Sort(directives)

	writeToFile("./doc/IniDirectives.md", generateMarkdownIni(directives))
}

func generateMarkdownIni(directives []string) string {
	var markdown strings.Builder

	markdown.WriteString("# Ini directives\n")

	// Get and sort categories
	categories := []string{}
	for category := range iniDirectives {
		categories = append(categories, category)
	}
	slices.Sort(categories)
	for _, category := range categories {
		if category != "" {
			markdown.WriteString("\n## ")
			markdown.WriteString(category)
			markdown.WriteString("\n")
		}

		directives := iniDirectives[category]
		slices.Sort(directives)

		for _, directive := range directives {
			markdown.WriteString("- ")
			markdown.WriteString(directive)
			markdown.WriteString("\n")
		}
	}

	return markdown.String()
}

func readGoFileIni(path string) error {
	file, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("failed to open file %s: %s", path, err)
	}
	defer file.Close()

	category := ""
	scanner := bufio.NewScanner(file)
	isIniSection := false
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}

		if strings.HasPrefix(line, "// Ini Directives:") {
			isIniSection = true
			continue
		}

		if !isIniSection {
			continue
		}

		if strings.HasPrefix(line, "// Category: ") {
			category = strings.Replace(line, "// Category: ", "", 1)
			continue
		}

		// Only process lines that look like: '"key": "value",'
		if !strings.HasPrefix(line, "\"") {
			continue
		}

		if line == "}" {
			break
		}

		parts := strings.SplitN(line, ":", 2)
		if len(parts) < 2 {
			continue
		}

		// Extract directive name
		directiveName := strings.Trim(parts[0], "\" ")

		// Extract default value
		defaultValue := strings.TrimSpace(parts[1])
		defaultValue = strings.TrimSuffix(defaultValue, ",")
		defaultValue = strings.Trim(defaultValue, "\"")

		if _, found := iniDirectives[category]; !found {
			iniDirectives[category] = []string{}
		}

		iniDirectives[category] = append(iniDirectives[category], directiveName)
	}
	return nil
}
