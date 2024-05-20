package interpreter

type Request struct {
	GetParams  map[IRuntimeValue]IRuntimeValue
	PostParams map[IRuntimeValue]IRuntimeValue
}

func NewRequest() *Request {
	return &Request{GetParams: map[IRuntimeValue]IRuntimeValue{}, PostParams: map[IRuntimeValue]IRuntimeValue{}}
}
