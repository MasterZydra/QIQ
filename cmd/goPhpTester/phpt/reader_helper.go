package phpt

func (reader *Reader) isEof() bool {
	return reader.currentLine > len(reader.lines)-1
}

func (reader *Reader) at() string {
	if reader.isEof() {
		return ""
	}
	return reader.lines[reader.currentLine]
}

func (reader *Reader) eat() string {
	if reader.isEof() {
		return ""
	}

	result := reader.at()
	reader.currentLine++
	return result
}
