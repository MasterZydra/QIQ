package main

import (
	"GoPHP/cmd/goPHP/interpreter"
	"GoPHP/cmd/goPhpTester/phpt"
	"fmt"
	"os"
	"path/filepath"
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
	if len(testFile.PostParams) > 0 {
		for key, value := range testFile.PostParams {
			request.PostParams[interpreter.NewStringRuntimeValue(key)] =
				interpreter.NewStringRuntimeValue(value[0])
		}
	}
	if len(testFile.GetParams) > 0 {
		for key, value := range testFile.GetParams {
			request.GetParams[interpreter.NewStringRuntimeValue(key)] =
				interpreter.NewStringRuntimeValue(value[0])
		}
	}

	result, err := interpreter.NewInterpreter(request).Process(testFile.File)
	if err != nil {
		fmt.Println("FAIL ", path)
		fmt.Println("     ", err)
		// return err
		failed++
		return nil
	}

	if testFile.Expect == result {
		fmt.Println("OK   ", path)
		succeeded++
		return nil
	} else {
		fmt.Println("FAIL ", path)
		fmt.Println("--------------- Expected ---------------")
		fmt.Print(testFile.Expect)
		fmt.Println("---------------   Got    ---------------")
		fmt.Print(result)
		fmt.Println("----------------------------------------")
		fmt.Println("")
		// return fmt.Errorf("Test \"%s\" failed", path)
		failed++
		return nil
	}
}
