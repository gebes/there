package middlewares

import (
	. "github.com/Gebes/there/v2"
	"net/http"
)

func Recoverer(request Request, next Response) Response {
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
