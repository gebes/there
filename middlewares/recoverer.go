package middlewares

import (
	"net/http"

	"github.com/Gebes/there/v2"
)

func Recoverer(request there.Request, next there.Response) there.Response {
	fn := func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if rvr := recover(); rvr != nil && rvr != http.ErrAbortHandler {
				there.Error(there.StatusInternalServerError, rvr).ServeHTTP(w, r)
			}
		}()
		next.ServeHTTP(w, r)
	}
	return there.HttpResponseFunc(fn)
}
