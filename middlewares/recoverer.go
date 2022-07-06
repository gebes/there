package middlewares

import (
	"fmt"
	"net/http"

	. "github.com/Gebes/there/v2"
)

func Recoverer(request HttpRequest, next HttpResponse) HttpResponse {
	fn := func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if rvr := recover(); rvr != nil && rvr != http.ErrAbortHandler {
				Error(StatusInternalServerError, fmt.Errorf("%v", rvr)).ServeHTTP(w, r)
			}
		}()
		next.ServeHTTP(w, r)
	}
	return HttpResponseFunc(fn)
}
