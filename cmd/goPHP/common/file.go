package common

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"
)

func GetAbsPath(path string) (string, error) {
	if !PathExists(path) {
		return path, fmt.Errorf("Could not open file: %s", path)
	}
	return ToAbsPath(path), nil
}

func GetAbsPathForWorkingDir(workingDir string, subPath string) string {
	if IsAbsPath(subPath) {
		return subPath
	}
	return path.Join(workingDir, subPath)
}

func ToAbsPath(path string) string {
	absFilePath, _ := filepath.Abs(path)
	return absFilePath
}

func IsAbsPath(path string) bool {
	return strings.HasPrefix(path, "/")
}

func PathExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

func IsDir(path string) bool {
	fileInfo, err := os.Stat(path)
	if err != nil {
		return false
	}
	return fileInfo.IsDir()
}

func ExtractPath(path string) string {
	return filepath.Dir(path)
}

func MkDir(path string) error {
	return os.MkdirAll(path, os.ModePerm)
}

func WriteFile(filename, content string) error {
	path := ExtractPath(filename)
	if !PathExists(path) {
		if err := MkDir(path); err != nil {
			return err
		}
	}

	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("failed to create file %s: %s", path, err)
	}

	_, err = file.WriteString(content)
	if err != nil {
		return fmt.Errorf("failed to write to file %s: %s", path, err)
	}
	return nil
}

func ReadFile(filename string) (string, error) {
	content, err := os.ReadFile(filename)
	return string(content), err
}
