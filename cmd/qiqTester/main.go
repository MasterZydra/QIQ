package main

import (
	"QIQ/cmd/qiq/common"
	"QIQ/cmd/qiq/common/os"
	"QIQ/cmd/qiq/ini"
	"QIQ/cmd/qiq/interpreter"
	"QIQ/cmd/qiq/phpError"
	"QIQ/cmd/qiq/request"
	"QIQ/cmd/qiq/runtime"
	"QIQ/cmd/qiqTester/phpt"
	replacejson "QIQ/cmd/qiqTester/replaceJson"
	"encoding/json"
	"flag"
	"fmt"
	goOs "os"
	"path/filepath"
	"regexp"
	goRuntime "runtime"
	"strings"
)

var succeeded int = 0
var failed int = 0
var skipped int = 0

var verbosity1 bool
var verbosity2 bool
var onlyFailed bool
var noColor bool
var replaceJson string

func main() {
	verbosity1Flag := flag.Bool("v1", false, "Verbosity level 1: Show all tests")
	verbosity2Flag := flag.Bool("v2", false, "Verbosity level 2: Show all tests and failure reason")
	onlyFailedFlag := flag.Bool("only-failed", false, "Show only failed tests")
	noColorFlag := flag.Bool("no-color", false, "Do not use color in output")
	replaceJsonFlag := flag.String("replace-json", "", "JSON to replace parts of a test")

	flag.Parse()
	verbosity1 = *verbosity1Flag
	verbosity2 = *verbosity2Flag
	onlyFailed = *onlyFailedFlag
	noColor = *noColorFlag
	replaceJson = *replaceJsonFlag

	args := goOs.Args[1:]

	if len(args) == 0 {
		fmt.Println("Usage: qiqTester [-v(1|2)] [-only-failed] [-no-color] [-replace-json <file-name>] [list of folders or files]")
		goOs.Exit(1)
	}

	if replaceJson != "" {
		if err := loadReplaceJson(); err != nil {
			fmt.Println(err)
			goOs.Exit(1)
		}
	}

	failed = 0
	succeeded = 0
	skipped = 0

	if !verbosity1 && !verbosity2 {
		println("Running test...")
	}
	continueNext := false
	for _, arg := range args {
		if continueNext {
			continueNext = false
			continue
		}
		if arg == "-replace-json" {
			continueNext = true
			continue
		}
		if arg == "-v1" || arg == "-v2" || arg == "-only-failed" || arg == "-no-color" {
			continue
		}

		if err := process(arg); err != nil {
			fmt.Println(err)
			goOs.Exit(1)
		}
	}

	if noColor {
		fmt.Printf("%5d Tests succeeded.\n%5d Tests failed.\n%5d Tests skipped.\n", succeeded, failed, skipped)
	} else {
		fmt.Printf("\n\033[32m%5d\033[0m Tests succeeded.\n\033[31m%5d\033[0m Tests failed.\n\033[33m%5d\033[0m Tests skipped.\n", succeeded, failed, skipped)
	}
}

func process(path string) error {
	file, err := goOs.Open(path)
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
			return fmt.Errorf(`Test files must have the extension "phpt". Got: "%s"`, file.Name())
		}
		return filepath.Walk(path, doTest)
	}
}

func doTest(path string, info goOs.FileInfo, _ error) error {
	if info.IsDir() || !strings.HasSuffix(strings.ToLower(info.Name()), ".phpt") {
		return nil
	}

	replaceEntry, hasReplaceEntry := replaceData.GetEntry(path)
	reader, err := phpt.NewReader(path, replaceEntry)
	if err != nil {
		if verbosity1 || verbosity2 {
			printFail(path)
		}
		if verbosity2 {
			fmt.Println("     ", err)
		}
		failed++
		return nil
	}
	testFile, err := reader.GetTestFile()
	if err != nil {
		if verbosity1 || verbosity2 {
			printFail(path)
		}
		if verbosity2 {
			fmt.Println("     ", err)
		}
		failed++
		return nil
	}

	request := request.NewRequest()
	request.Env = testFile.Env
	// TODO Add path to TEST_PHP_CGI_EXECUTABLE
	// request.Env["TEST_PHP_CGI_EXECUTABLE"] = "1"
	request.Args = testFile.Args
	request.Cookie = testFile.Cookie
	request.QueryString = testFile.Get
	request.Post = testFile.Post

	result := ""
	devIni, phpErr := ini.NewDevIniFromArray(testFile.Ini)
	if phpErr != nil {
		result = phpErr.GetMessage()
	} else {
		interpr, phpErr := interpreter.NewInterpreter(runtime.NewExecutionContext(), devIni, request, testFile.Filename[:len(testFile.Filename)-1])
		if phpErr != nil {
			result = phpErr.GetMessage()
		} else {
			result, phpErr = interpr.Process(testFile.File)
			if phpErr != nil {
				result += "\n" + phpErr.GetMessage()
			}
		}
		if err := common.DeleteFiles(request.UploadedFiles); err != nil {
			fmt.Printf("Cleanup failed: %s\n", err)
		}
	}

	if goRuntime.GOOS == "windows" {
		testFile.Expect = strings.ReplaceAll(testFile.Expect, "\r", "")
		result = strings.ReplaceAll(result, "\r", "")
	}

	if strings.HasPrefix(result, "skip for") || strings.HasPrefix(result, "skip Run") ||
		strings.HasPrefix(result, "skip only") || strings.HasPrefix(result, "skip this") ||
		strings.HasPrefix(result, "skip.. ") || strings.HasPrefix(result, "skip ") {
		if !onlyFailed && (verbosity1 || verbosity2) {
			if noColor {
				fmt.Println("SKIP ", path)
			} else {
				fmt.Println("\033[33mSKIP\033[0m ", path)
			}
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
			if !hasReplaceEntry {
				// Special case for QIQ:
				// The location of the error is given with line and column:
				// e.g. "... in tests/basic/025.phpt:2:18"
				// PHP only returns the line:
				// e.g. "... in tests/basic/025.phpt on line 2"
				if strings.Contains(testFile.Expect, " on line %d") {
					testFile.Expect = strings.ReplaceAll(testFile.Expect, " on line %d", ":%d:%d")
				}
				if matched, _ := regexp.MatchString(`in %s on line \d+`, testFile.Expect); matched {
					testFile.Expect = regexp.MustCompile(`in %s on line \d+`).ReplaceAllString(testFile.Expect, "in %s")
				}
			}
			pattern = replaceExpectfTags(testFile.Expect)
		}
		equal, err = regexp.MatchString(common.TrimLineBreaks(pattern), common.TrimLineBreaks(result))
		if err != nil {
			if verbosity1 || verbosity2 {
				printFail(path)
			}
			if verbosity2 {
				fmt.Printf("      %s\n", err)
			}
			failed++
			return nil
		}

	default:
		if verbosity1 || verbosity2 {
			printFail(path)
		}
		if verbosity2 {
			fmt.Printf("      Unsupported expect section: %s\n", testFile.ExpectType)
		}
		failed++
		return nil
	}

	if equal {
		if !onlyFailed && (verbosity1 || verbosity2) {
			if noColor {
				fmt.Println("OK   ", path)
			} else {
				fmt.Println("\033[32mOK\033[0m   ", path)
			}
		}
		succeeded++
		return nil
	} else {
		if verbosity1 || verbosity2 {
			printFail(path)
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
		`%e`: os.DIR_SEP,                    // Directory separator
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

func printFail(path string) {
	if noColor {
		fmt.Println("FAIL ", path)
	} else {
		fmt.Println("\033[31mFAIL\033[0m ", path)
	}
}

var replaceData replacejson.ReplaceJson

func loadReplaceJson() phpError.Error {
	info, err := goOs.Stat(replaceJson)
	if err != nil {
		return phpError.NewError("Replace JSON file not found: %s", err)
	}
	if info.IsDir() {
		return phpError.NewError("Replace JSON path is a directory, not a file")
	}
	data, err := goOs.ReadFile(replaceJson)
	if err != nil {
		return phpError.NewError("Failed to read replace JSON file: %s", err)
	}
	err = json.Unmarshal(data, &replaceData)
	if err != nil {
		return phpError.NewError("Failed to parse replace JSON: %s", err)
	}
	return nil
}
