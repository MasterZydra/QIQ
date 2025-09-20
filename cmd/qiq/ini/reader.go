package ini

import (
	"QIQ/cmd/qiq/config"
	"QIQ/cmd/qiq/phpError"
	"bufio"
	"os"
	"strings"
)

type Reader struct {
	ini *Ini
}

func NewReader() *Reader {
	reader := &Reader{}
	if config.IsDevMode {
		reader.ini = NewDevIni()
	} else {
		reader.ini = NewDefaultIni()
	}
	return reader
}

func (reader *Reader) GetIni() *Ini { return reader.ini }

func (reader *Reader) Read(filename string) phpError.Error {
	file, err := os.Open(filename)
	if err != nil {
		return phpError.NewError("%s", err.Error())
	}
	defer file.Close()

	section := ""
	scanner := bufio.NewScanner(file)
	lineNumber := 0
	for scanner.Scan() {
		lineNumber++

		line := strings.TrimSpace(scanner.Text())

		if line == "" || strings.HasPrefix(line, ";") {
			continue
		}

		// Read section
		if strings.HasPrefix(line, "[") && strings.Contains(line, "]") {
			section = strings.TrimSpace(line[strings.Index(line, "[")+1 : strings.Index(line, "]")])
			continue
		}

		// PHP section
		if section != "PHP" {
			continue
		}

		if !strings.Contains(line, "=") {
			return phpError.NewError("Cannot parse ini file line %d", lineNumber)
		}

		parts := strings.SplitN(line, "=", 2)
		directive, value := strings.TrimSpace(parts[0]), strings.TrimSpace(parts[1])
		value = strings.TrimPrefix(value, `"`)
		value = strings.TrimSuffix(value, `"`)
		reader.ini.Set(directive, value, INI_ALL)

		// TODO Ini Reader: Add support for environment variables
		// Spec: https://www.php.net/manual/en/configuration.file.php
		/*
			; PHP_MEMORY_LIMIT is taken from environment
			memory_limit = ${PHP_MEMORY_LIMIT}
			; If PHP_MAX_EXECUTION_TIME is not defined, it will fall back to 30
			max_execution_time = ${PHP_MAX_EXECUTION_TIME:-30}
		*/
		/*
			; any text on a line after an unquoted semicolon (;) is ignored
			[php] ; section markers (text within square brackets) are also ignored
			; Boolean values can be set to either:
			;    true, on, yes
			; or false, off, no, none
			register_globals = off
			track_errors = yes

			; you can enclose strings in double-quotes
			include_path = ".:/usr/local/lib/php"

			; backslashes are treated the same as any other character
			include_path = ".;c:\php\lib"
		*/
	}

	return nil
}
