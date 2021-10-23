package there

import (
	. "github.com/Gebes/there/there/utils"
	"net/http"
)

type RedirectResponse struct {
	url string
	header *HeaderWrapper
}

func Redirect(url string) *RedirectResponse {
	r := &RedirectResponse{url: url}
	r.header = Header(r)
	return r
}

func (j RedirectResponse) Execute(r *http.Request, w *http.ResponseWriter) error {
	(*w).WriteHeader(StatusMovedPermanently)
	http.Redirect(*w, r, j.url, StatusMovedPermanently)
	return nil
}

func (j *RedirectResponse) Header() *HeaderWrapper {
	return j.header
}
