package middlewares

import (
	. "github.com/Gebes/there/v2"
	"net/http"
)

func Recoverer(request HttpRequest, next HttpResponse) HttpResponse {
	fn := func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if rvr := recover(); rvr != nil && rvr != http.ErrAbortHandler {
				Error(StatusInternalServerError, rvr).ServeHTTP(w, r)
			}
		}()
		next.ServeHTTP(w, r)
	}
	return HttpResponseFunc(fn)
}
