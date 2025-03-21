package main

import (
	"GoPHP/cmd/goPHP/common"
	"GoPHP/cmd/goPHP/ini"
	"GoPHP/cmd/goPHP/interpreter"
	"GoPHP/cmd/goPHP/request"
	"GoPHP/cmd/goPhpTester/phpt"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
)

var succeeded int = 0
var failed int = 0
var skipped int = 0

var verbosity1 bool
var verbosity2 bool
var onlyFailed bool

func main() {
	verbosity1Flag := flag.Bool("v1", false, "Verbosity level 1: Show all tests")
	verbosity2Flag := flag.Bool("v2", false, "Verbosity level 2: Show all tests and failure reason")
	onlyFailedFlag := flag.Bool("only-failed", false, "Show only failed tests")
	flag.Parse()
	verbosity1 = *verbosity1Flag
	verbosity2 = *verbosity2Flag
	onlyFailed = *onlyFailedFlag

	args := os.Args[1:]

	if len(args) == 0 {
		fmt.Println("Usage: goPhpTester [-v(1|2)] [-only-failed] [list of folders or files]")
		os.Exit(1)
	}

	failed = 0
	succeeded = 0
	skipped = 0

	if !verbosity1 && !verbosity2 {
		println("Running test...")
	}
	for _, arg := range args {
		if arg == "-v1" || arg == "-v2" || arg == "-only-failed" {
			continue
		}

		if err := process(arg); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}

	fmt.Printf("\n%d Tests succeeded.\n%d Tests failed.\n%d Tests skipped.\n", succeeded, failed, skipped)
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
		if verbosity1 || verbosity2 {
			fmt.Println("FAIL ", path)
		}
		if verbosity2 {
			fmt.Println("     ", err)
		}
		// return err
		failed++
		return nil
	}
	testFile, err := reader.GetTestFile()
	if err != nil {
		if verbosity1 || verbosity2 {
			fmt.Println("FAIL ", path)
		}
		if verbosity2 {
			fmt.Println("     ", err)
		}
		// return err
		failed++
		return nil
	}

	request := request.NewRequest()
	request.Env = testFile.Env
	request.Args = testFile.Args
	request.Cookie = testFile.Cookie
	request.QueryString = testFile.Get
	request.Post = testFile.Post

	result, phpError := interpreter.NewInterpreter(ini.NewDevIniFromArray(testFile.Ini), request, testFile.Filename).Process(testFile.File)
	if phpError != nil {
		result += phpError.GetMessage()
	}
	if err := common.DeleteFiles(request.UploadedFiles); err != nil {
		fmt.Printf("Cleanup failed: %s\n", err)
	}

	if runtime.GOOS == "windows" {
		testFile.Expect = strings.ReplaceAll(testFile.Expect, "\r", "")
		result = strings.ReplaceAll(result, "\r", "")
	}

	if strings.HasPrefix(result, "skip for") || strings.HasPrefix(result, "skip Run") ||
		strings.HasPrefix(result, "skip only") || strings.HasPrefix(result, "skip this") ||
		strings.HasPrefix(result, "skip.. ") || strings.HasPrefix(result, "skip ") {
		if !onlyFailed && (verbosity1 || verbosity2) {
			fmt.Println("SKIP ", path)
		}
		if !onlyFailed && (verbosity2) {
			reason := strings.TrimPrefix(result, "skip ")
			reason = strings.TrimPrefix(reason, "skip.. ")
			reason = strings.ToUpper(string(reason[0])) + reason[1:]
			fmt.Println("     ", reason)
		}
		skipped++
		return nil
	}

	var equal bool
	switch testFile.ExpectType {
	case "--EXPECT--":
		equal = testFile.Expect == common.TrimLineBreaks(result)

	case "--EXPECTF--", "--EXPECTREGEX--":
		pattern := testFile.Expect
		if testFile.ExpectType == "--EXPECTF--" {
			// Special case for goPHP:
			// The location of the error is given with line and column:
			// e.g. "... in tests/basic/025.phpt:2:18"
			// PHP only returns the line:
			// e.g. "... in tests/basic/025.phpt on line 2"
			if strings.Contains(testFile.Expect, "in %s on line %d") {
				testFile.Expect = strings.ReplaceAll(testFile.Expect, "in %s on line %d", "in %s")
			}
			pattern = replaceExpectfTags(testFile.Expect)
		}
		equal, err = regexp.MatchString(pattern, common.TrimLineBreaks(result))
		if err != nil {
			if verbosity1 || verbosity2 {
				fmt.Println("FAIL ", path)
			}
			if verbosity2 {
				fmt.Printf("      %s\n", err)
			}
			failed++
			return nil
		}

	default:
		if verbosity1 || verbosity2 {
			fmt.Println("FAIL ", path)
		}
		if verbosity2 {
			fmt.Printf("      Unsupported expect section: %s\n", testFile.ExpectType)
		}
		failed++
		return nil
	}

	if equal {
		if !onlyFailed && (verbosity1 || verbosity2) {
			fmt.Println("OK   ", path)
		}
		succeeded++
		return nil
	} else {
		if verbosity1 || verbosity2 {
			fmt.Println("FAIL ", path)
		}
		if verbosity2 {
			fmt.Println("--------------- Expected ---------------")
			fmt.Println(testFile.Expect)
			fmt.Println("---------------   Got    ---------------")
			fmt.Println(result)
			fmt.Println("----------------------------------------")
			fmt.Println("")
		}
		failed++
		return nil
	}
}

func replaceExpectfTags(value string) string {
	// Spec: https://qa.php.net/phpt_details.php#expectf_section

	replacements := map[string]string{
		`%e`: interpreter.DIR_SEP,           // Directory separator
		`%s`: `[^\n]+`,                      // One or more of anything except the end of line
		`%S`: `[^\n]*`,                      // Zero or more of anything except the end of line
		`%a`: `.+`,                          // One or more of anything including the end of line
		`%A`: `.*`,                          // Zero or more of anything including the end of line
		`%w`: `\s*`,                         // Zero or more white space characters
		`%i`: `[+-]?\d+`,                    // A signed integer value
		`%d`: `\d+`,                         // An unsigned integer value
		`%x`: `[0-9a-fA-F]+`,                // One or more hexadecimal characters
		`%f`: common.FloatingLiteralPattern, // A floating point number
		`%c`: `.`,                           // A single character of any sort
	}

	value = regexp.QuoteMeta(value)
	for key, replacement := range replacements {
		value = strings.ReplaceAll(value, key, replacement)
	}

	// Handle %r...%r for regular expressions
	re := regexp.MustCompile(`%r(.*?)%r`)
	value = re.ReplaceAllString(value, `($1)`)

	return value
}
