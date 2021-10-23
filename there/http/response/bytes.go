package there

import (
	"net/http"
)

type BytesResponse struct {
	code   int
	data   []byte
	header *HeaderWrapper
}

func Bytes(code int, data []byte) *BytesResponse {
	r := &BytesResponse{data: data, code: code}
	r.header = Header(r)
	return r
}

func (j BytesResponse) Execute(r *http.Request, w *http.ResponseWriter) error {
	(*w).WriteHeader(j.code)
	_, err := (*w).Write(j.data)
	return err
}

func (j *BytesResponse) Header() *HeaderWrapper {
	return j.header
}
