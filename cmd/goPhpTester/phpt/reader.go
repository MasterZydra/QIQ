package phpt

import (
	"bufio"
	"fmt"
	"os"
	"slices"
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
	getError := func(section string) error {
		return fmt.Errorf("Expected \"%s\". Got: \"%s\"", section, reader.at())
	}

	if reader.at() == "--TEST--" {
		reader.eat()
		reader.testFile.Title = reader.eat()
		reader.sections = append(reader.sections, "--TEST--")
	} else {
		return nil, getError("--TEST--")
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
	}
	if reader.at() == "--FILE--" {
		reader.eat()
		file := ""
		for !reader.isEof() && !reader.isSection(reader.at()) {
			file += reader.eat() + "\n"
		}
		reader.testFile.File = file
		reader.sections = append(reader.sections, "--FILE--")
	} else {
		return nil, getError("--FILE--")
	}
	if reader.at() == "--EXPECT--" {
		reader.eat()
		expect := ""
		for !reader.isEof() && !reader.isSection(reader.at()) {
			expect += reader.eat() + "\n"
		}
		reader.testFile.Expect = expect
		reader.sections = append(reader.sections, "--EXPECT--")
	} else {
		return nil, getError("--EXPECT--")
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
