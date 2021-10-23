package there

import (
	"net/http"
)

type HttpResponse interface {
	Header() *HeaderWrapper
	Execute(r *http.Request, w *http.ResponseWriter) error
}



