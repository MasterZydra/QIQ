package phpt

import (
	"GoPHP/cmd/goPHP/common"
	"bufio"
	"fmt"
	"os"
	"slices"
	"strings"
)

type Reader struct {
	filename    string
	sections    []string
	lines       []string
	currentLine int
	testFile    *TestFile
}

func NewReader(filename string) (*Reader, error) {
	lines := []string{}

	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return &Reader{filename: filename, sections: []string{}, lines: lines, currentLine: 0, testFile: NewTestFile()}, nil
}

func (reader *Reader) GetTestFile() (*TestFile, error) {
	for !reader.isEof() {
		if reader.at() == "--TEST--" {
			reader.eat()
			reader.testFile.Title = reader.eat()
			for !reader.isEof() && !reader.isSection(reader.at()) {
				reader.testFile.Title += reader.eat()
			}
			reader.sections = append(reader.sections, "--TEST--")
			continue
		}

		if reader.at() == "--CREDITS--" || reader.at() == "--EXTENSIONS--" {
			reader.eat()
			for !reader.isEof() && !reader.isSection(reader.at()) {
				reader.eat()
			}
			continue
		}

		if reader.at() == "--SKIPIF--" {
			reader.eat()
			file := ""
			for !reader.isEof() && !reader.isSection(reader.at()) {
				if file != "" {
					file += "\n"
				}
				file += reader.eat()
			}
			reader.testFile.File = file
			reader.sections = append(reader.sections, "--SKIPIF--")
			continue
		}

		if reader.at() == "--POST--" {
			reader.eat()
			params := ""
			for !reader.isEof() && !reader.isSection(reader.at()) {
				params += reader.eat()
			}
			paramsMap, err := parseQuery(params)
			if err != nil {
				return nil, fmt.Errorf("--POST--\n%s", err)
			}
			reader.testFile.PostParams = paramsMap
			reader.sections = append(reader.sections, "--POST--")
			continue
		}

		if reader.at() == "--GET--" {
			reader.eat()
			params := ""
			for !reader.isEof() && !reader.isSection(reader.at()) {
				params += reader.eat()
			}
			paramsMap, err := parseQuery(params)
			if err != nil {
				return nil, fmt.Errorf("--GET--\n%s", err)
			}
			reader.testFile.GetParams = paramsMap
			reader.sections = append(reader.sections, "--GET--")
			continue
		}

		if reader.at() == "--INI--" {
			ini := []string{}
			reader.eat()
			for !reader.isEof() && !reader.isSection(reader.at()) {
				ini = append(ini, reader.eat())
			}
			reader.testFile.Ini = ini
			continue
		}

		if reader.at() == "--ENV--" {
			reader.eat()
			env := map[string]string{}
			for !reader.isEof() && !reader.isSection(reader.at()) {
				line := reader.eat()
				separator := strings.Index(line, "=")
				env[line[:separator]] = line[separator+1:]
			}
			reader.testFile.Env = env
			reader.sections = append(reader.sections, "--ENV--")
			continue
		}

		if reader.at() == "--FILE--" {
			reader.eat()
			file := reader.testFile.File
			for !reader.isEof() && !reader.isSection(reader.at()) {
				file += reader.eat() + "\n"
			}
			reader.testFile.File = file
			reader.sections = append(reader.sections, "--FILE--")
			continue
		}

		if reader.at() == "--EXPECT--" || reader.at() == "--EXPECTF--" {
			section := reader.eat()
			expect := ""
			for !reader.isEof() && !reader.isSection(reader.at()) {
				expect += reader.eat() + "\n"
			}
			reader.testFile.Expect = common.TrimTrailingLineBreaks(expect)
			reader.testFile.ExpectType = section
			reader.sections = append(reader.sections, section)
			continue
		}

		if reader.at() == "--CLEAN--" {
			reader.eat()
			file := reader.testFile.File
			for !reader.isEof() && !reader.isSection(reader.at()) {
				file += reader.eat() + "\n"
			}
			reader.testFile.File = file
			continue
		}

		return reader.testFile, fmt.Errorf("Unsupported section \"%s\"", reader.at())
	}

	if err := reader.isValid(); err != nil {
		return nil, err
	}

	return reader.testFile, nil
}

func (reader *Reader) isValid() error {
	if !reader.hasSection("--TEST--") {
		return fmt.Errorf("Section \"--TEST--\" is missing")
	}

	if !reader.hasSection("--FILE--") && !reader.hasSection("--FILEEOF--") &&
		!reader.hasSection("--FILE_EXTERNAL--") && !reader.hasSection("--REDIRECTTEST--") {
		return fmt.Errorf("Section \"--FILE-- | --FILEEOF-- | --FILE_EXTERNAL-- | --REDIRECTTEST--\" is missing")
	}

	if !reader.hasSection("--EXPECT--") && !reader.hasSection("--EXPECTF--") && !reader.hasSection("--EXPECTREGEX--") &&
		!reader.hasSection("--EXPECT_EXTERNAL--") && !reader.hasSection("--EXPECTF_EXTERNAL--") &&
		!reader.hasSection("--EXPECTREGEX_EXTERNAL--") {
		return fmt.Errorf("Section \"--EXPECT-- | --EXPECTF-- | --EXPECTREGEX-- | --EXPECT_EXTERNAL-- | --EXPECTF_EXTERNAL-- | --EXPECTREGEX_EXTERNAL--\" is missing")
	}
	return nil
}

func (reader *Reader) hasSection(section string) bool {
	return slices.Contains(reader.sections, section)
}

var sections = []string{
	"--TEST--", "--DESCRIPTION--", "--CREDITS--", "--SKIPIF--", "--CONFLICTS--", "--WHITESPACE_SENSITIVE--", "--CAPTURE_STDIO--",
	"--EXTENSIONS--", "--POST--", "--PUT--", "--POST_RAW--", "--GZIP_POST--", "--DEFLATE_POST--", "--GET--", "--COOKIE--",
	"--STDIN--", "--INI--", "--ARGS--", "--ENV--", "--PHPDBG--", "--FILE--", "--FILEEOF--", "--FILE_EXTERNAL--",
	"--REDIRECTTEST--", "--CGI--", "--XFAIL--", "--EXPECTHEADERS--", "--EXPECT--", "--EXPECTF--", "--EXPECTREGEX--",
	"--EXPECT_EXTERNAL--", "--EXPECTF_EXTERNAL--", "--EXPECTREGEX_EXTERNAL--", "--CLEAN--",
}

func (reader *Reader) isSection(section string) bool {
	return slices.Contains(sections, section)
}

/*
Spec: https://qa.php.net/phpt_details.php
[] indicates optional sections.

--TEST--
[--DESCRIPTION--]
[--CREDITS--]
[--SKIPIF--]
[--CONFLICTS--]
[--WHITESPACE_SENSITIVE--]
[--CAPTURE_STDIO--]
[--EXTENSIONS--]
[--POST-- | --PUT-- | --POST_RAW-- | --GZIP_POST-- | --DEFLATE_POST-- | --GET--]
[--COOKIE--]
[--STDIN--]
[--INI--]
[--ARGS--]
[--ENV--]
[--PHPDBG--]
--FILE-- | --FILEEOF-- | --FILE_EXTERNAL-- | --REDIRECTTEST--
[--CGI--]
[--XFAIL--]
[--EXPECTHEADERS--]
--EXPECT-- | --EXPECTF-- | --EXPECTREGEX-- | --EXPECT_EXTERNAL-- | --EXPECTF_EXTERNAL-- | --EXPECTREGEX_EXTERNAL--
[--CLEAN--]
*/
