package middlewares

import (
	"github.com/Gebes/there/v2"
	"github.com/Gebes/there/v2/status"
	"path"
)

func Sanitizer(request there.Request, next there.Response) there.Response {
	cleaned := path.Clean(request.Request.URL.Path)
	if cleaned != request.Request.URL.Path {
		return there.Redirect(status.TemporaryRedirect, cleaned)
	}
	return next
}
