package position

import "fmt"

type Position struct {
	Filename string
	Line     int
	Column   int
}

func NewPosition(filename string, line int, column int) *Position {
	return &Position{Filename: filename, Line: line, Column: column}
}

func (pos *Position) ToPosString() string {
	return fmt.Sprintf("%s:%d:%d", pos.Filename, pos.Line, pos.Column)
}

func (pos *Position) String() string {
	return fmt.Sprintf("{Position - file: \"%s\", ln: %d, col: %d}", pos.Filename, pos.Line, pos.Column)
}
