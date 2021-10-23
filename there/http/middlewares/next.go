package there

import (
	"errors"
	. "github.com/Gebes/there/there/http/response"
	"net/http"
)

type NextMiddleware struct {
	header *HeaderWrapper
}

func Next() *NextMiddleware {
	r := &NextMiddleware{}
	r.header = &HeaderWrapper{Values: map[string][]string{}, HttpResponse: r}
	return r
}

func (j NextMiddleware) Execute(r *http.Request, w *http.ResponseWriter) error {
	return errors.New("cannot execute next middleware")
}

func (j *NextMiddleware) Header() *HeaderWrapper {
	return j.header
}
