package middlewares

import (
	"errors"
	"fmt"
	"net/http"

	. "github.com/Gebes/there/v2"
)

func Recoverer(request there.Request, next there.Response) there.Response {
	fn := func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if rvr := recover(); rvr != nil && rvr != http.ErrAbortHandler {
				err := fmt.Sprintf("%v", rvr)
				Error(StatusInternalServerError, errors.New(err)).ServeHTTP(w, r)
			}
		}()
		next.ServeHTTP(w, r)
	}
	return there.ResponseFunc(fn)
}
