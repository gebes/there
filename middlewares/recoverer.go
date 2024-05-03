package middlewares

import (
	"errors"
	"fmt"
	"github.com/Gebes/there/v2/status"
	"net/http"

	"github.com/Gebes/there/v2"
)

func Recoverer(request there.Request, next there.Response) there.Response {
	fn := func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if rvr := recover(); rvr != nil && rvr != http.ErrAbortHandler {
				switch t := rvr.(type) {
				case error:
					there.Error(status.InternalServerError, t).ServeHTTP(w, r)
				default:
					there.Error(status.InternalServerError, errors.New(fmt.Sprint(t))).ServeHTTP(w, r)
				}
			}
		}()
		next.ServeHTTP(w, r)
	}
	return there.ResponseFunc(fn)
}
