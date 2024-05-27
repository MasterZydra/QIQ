package interpreter

type Request struct {
	GetParams  [][]string
	PostParams [][]string
}

func NewRequest() *Request {
	return &Request{GetParams: [][]string{}, PostParams: [][]string{}}
}
