package there

import (
	"net/http"
)

type EmptyResponse struct {
	code int
	header *HeaderWrapper
}

func Empty(code int) *EmptyResponse {
	r := &EmptyResponse{code: code}
	r.header = Header(r)
	return r
}

func (j EmptyResponse) Execute(r *http.Request, w *http.ResponseWriter) error {
	(*w).WriteHeader(j.code)
	_, err := (*w).Write([]byte(""))
	return err
}

func (j *EmptyResponse) Header() *HeaderWrapper {
	return j.header
}
