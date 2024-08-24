package request

import "context"

type Request struct {
	Method        string
	Path          string
	Header        map[string]string
	Body          []byte
	Proto         string
	ContentLength int
	Context       context.Context
}

func (r Request) WithContext(ctx context.Context) *Request {
	r.Context = ctx
	return &r
}
