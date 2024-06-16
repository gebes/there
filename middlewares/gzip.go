package middlewares

import (
	"github.com/gebes/there/v2"
)

func Gzip(request there.Request, next there.Response) there.Response {
	return there.Gzip(next)
}
