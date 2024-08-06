package main

import (
	"GoPHP/cmd/goPHP/common"
	"GoPHP/cmd/goPHP/ini"
	"GoPHP/cmd/goPHP/interpreter"
	"GoPHP/cmd/goPhpTester/phpt"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

var succeeded int = 0
var failed int = 0

func main() {
	args := os.Args[1:]

	if len(args) == 0 {
		fmt.Println("Usage: goPhpTester [list of folders or files]")
		os.Exit(1)
	}

	failed = 0
	succeeded = 0

	for _, arg := range args {
		if err := process(arg); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}

	fmt.Printf("\n%d Tests succeeded. %d Tests failed\n", succeeded, failed)
}

func process(path string) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()

	info, err := file.Stat()
	if err != nil {
		return err
	}

	if info.IsDir() {
		return filepath.Walk(path, doTest)
	} else {
		if !strings.HasSuffix(strings.ToLower(file.Name()), ".phpt") {
			return fmt.Errorf("Test files must have the extension \"phpt\". Got: \"%s\"", file.Name())
		}
		return filepath.Walk(path, doTest)
	}
}

func doTest(path string, info os.FileInfo, err error) error {
	if info.IsDir() || !strings.HasSuffix(strings.ToLower(info.Name()), ".phpt") {
		return nil
	}

	reader, err := phpt.NewReader(path)
	if err != nil {
		fmt.Println("FAIL ", path)
		fmt.Println("     ", err)
		// return err
		failed++
		return nil
	}
	testFile, err := reader.GetTestFile()
	if err != nil {
		fmt.Println("FAIL ", path)
		fmt.Println("     ", err)
		// return err
		failed++
		return nil
	}

	request := interpreter.NewRequest()
	request.Env = testFile.Env
	request.GetParams = testFile.GetParams
	request.PostParams = testFile.PostParams

	result, phpError := interpreter.NewInterpreter(ini.NewIniFromArray(testFile.Ini), request, "").Process(testFile.File)
	if phpError != nil {
		fmt.Println("FAIL ", path)
		fmt.Println("     ", phpError)
		// return err
		failed++
		return nil
	}

	if runtime.GOOS == "windows" {
		testFile.Expect = strings.ReplaceAll(testFile.Expect, "\r\n", "\n")
		result = strings.ReplaceAll(result, "\r\n", "\n")
	}

	var equal bool
	switch testFile.ExpectType {
	case "--EXPECT--":
		equal = testFile.Expect == common.TrimTrailingLineBreaks(result)
	default:
		failed++
		fmt.Errorf("Unsupported expect section: %s", testFile.ExpectType)
		return nil
	}

	if equal {
		fmt.Println("OK   ", path)
		succeeded++
		return nil
	} else {
		fmt.Println("FAIL ", path)
		fmt.Println("--------------- Expected ---------------")
		fmt.Println(testFile.Expect)
		fmt.Println("---------------   Got    ---------------")
		fmt.Println(result)
		fmt.Println("----------------------------------------")
		fmt.Println("")
		// return fmt.Errorf("Test \"%s\" failed", path)
		failed++
		return nil
	}
}
