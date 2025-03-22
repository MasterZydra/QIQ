package main

import (
	"fmt"
	"os"
)

func main() {
	docStdLib()
	docConst()
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
