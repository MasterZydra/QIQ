package interpreter

type Request struct {
	Env        map[string]string
	Args       [][]string
	GetParams  [][]string
	PostParams [][]string
}

func NewRequest() *Request {
	return &Request{Args: [][]string{}, GetParams: [][]string{}, PostParams: [][]string{}}
}
