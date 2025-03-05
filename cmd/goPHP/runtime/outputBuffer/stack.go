package outputBuffer

type Stack struct {
	buffers []*Buffer
}

func NewStack() *Stack {
	return &Stack{buffers: []*Buffer{}}
}

func (stack *Stack) Push() {
	stack.buffers = append(stack.buffers, NewBuffer())
}

func (stack *Stack) Pop() {
	stack.buffers = stack.buffers[:stack.Len()-1]
}

func (stack *Stack) Get(index int) *Buffer {
	return stack.buffers[index]
}

func (stack *Stack) GetLast() *Buffer {
	return stack.buffers[stack.Len()-1]
}

func (stack *Stack) Len() int {
	return len(stack.buffers)
}
