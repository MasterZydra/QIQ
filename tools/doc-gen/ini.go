package main

import (
	"bufio"
	"fmt"
	"os"
	"slices"
	"strings"
)

var iniDirectives = map[string]string{}

func docIniDirectives() {
	err := readGoFileIni("./cmd/qiq/ini/ini.go")
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

	for _, directive := range directives {
		markdown.WriteString("- ")
		markdown.WriteString(directive)
		markdown.WriteString(" (Default: \"")
		markdown.WriteString(iniDirectives[directive])
		markdown.WriteString("\")\n")
	}

	return markdown.String()
}

func readGoFileIni(path string) error {
	file, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("failed to open file %s: %s", path, err)
	}
	defer file.Close()

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

		// Only process lines that look like: '"key": "value",'
		if !strings.HasPrefix(line, "\"") {
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

		iniDirectives[directiveName] = defaultValue
	}
	return nil
}
