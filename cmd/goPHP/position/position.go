package position

import "fmt"

type File struct {
	Filename     string
	IsStrictType bool
}

func NewFile(filename string) *File {
	// TODO position - Set IsStrictType default value to true/false depending on future ini setting for a strict mode
	return &File{Filename: filename, IsStrictType: false}
}

type Position struct {
	File   *File
	Line   int
	Column int
}

func NewPosition(file *File, line int, column int) *Position {
	return &Position{File: file, Line: line, Column: column}
}

func (pos *Position) ToPosString() string {
	if pos.File == nil {
		pos.File = NewFile("")
	}
	return fmt.Sprintf("%s:%d:%d", pos.File.Filename, pos.Line, pos.Column)
}

func (pos *Position) String() string {
	return fmt.Sprintf("{Position - file: \"%s\", ln: %d, col: %d}", pos.File.Filename, pos.Line, pos.Column)
}
