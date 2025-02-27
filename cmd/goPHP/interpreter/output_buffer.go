package interpreter

type OutputBuffer struct {
	Content string
}

func NewOutputBuffer() *OutputBuffer {
	return &OutputBuffer{}
}
