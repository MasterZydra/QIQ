package interpreter

type Request struct {
	Env        map[string]string
	GetParams  [][]string
	PostParams [][]string
}

func NewRequest() *Request {
	return &Request{GetParams: [][]string{}, PostParams: [][]string{}}
}
