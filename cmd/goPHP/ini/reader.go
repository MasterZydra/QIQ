package ini

import (
	"GoPHP/cmd/goPHP/config"
	"GoPHP/cmd/goPHP/phpError"
	"bufio"
	"os"
	"strings"
)

type IniReader struct {
	ini *Ini
}

func NewIniReader() *IniReader {
	reader := &IniReader{}
	if config.IsDevMode {
		reader.ini = NewDevIni()
	} else {
		reader.ini = NewDefaultIni()
	}
	return reader
}

func (reader *IniReader) Read(filename string) phpError.Error {
	file, err := os.Open(filename)
	if err != nil {
		return phpError.NewError(err.Error())
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	lineNumber := 0
	for scanner.Scan() {
		lineNumber++

		line := strings.TrimSpace(scanner.Text())

		if line == "" || strings.HasPrefix(line, ";") {
			continue
		}
		if !strings.Contains(line, "=") {
			return phpError.NewError("Cannot parse ini file line %d", lineNumber)
		}

		println(strings.SplitN(line, "=", 2))
		// Trim directive, trim value
		// set directive in ini struct
	}

	return nil
}
