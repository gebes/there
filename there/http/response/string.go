package there

import (
	"net/http"
)

type StringResponse struct {
	code   int
	data   string
	header *HeaderWrapper
}

func String(code int, data string) *StringResponse {
	r := &StringResponse{data: data, code: code}
	r.header = Header(r)
	return r
}

func (j StringResponse) Execute(r *http.Request, w *http.ResponseWriter) error {
	(*w).WriteHeader(j.code)
	(*w).Write([]byte(j.data))
	return nil
}

func (j StringResponse) Header() *HeaderWrapper {
	return j.header
}
