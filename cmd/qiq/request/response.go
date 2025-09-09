package request

type Response struct {
	ExitCode int
}

func NewResponse() *Response { return &Response{ExitCode: 0} }
