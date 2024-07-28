package middlewares

import (
	"github.com/gebes/there/v2/status"

	"github.com/gebes/there/v2"
)

// RequireHost is a middleware, that forces incoming requests to have a specific host header set.
// If this header is not set, the String response with StatusBadRequest and the message
// "Invalid host for access to resource" is returned.
func RequireHost(host string) func(request there.Request, next there.Response) there.Response {
	return func(request there.Request, next there.Response) there.Response {
		if request.Host != host {
			return there.String(status.BadRequest, "Invalid host for access to resource")
		}
		return next
	}
}

// RequireHeader is a middleware, that forces incoming requests to have a specific header.
// If this method is not used, only Status with status StatusBadRequest is returned.
func RequireHeader(key string, value string) func(request there.Request, next there.Response) there.Response {
	return func(request there.Request, next there.Response) there.Response {
		if header := request.Headers.Get(key); header != value {
			return there.Status(status.BadRequest)
		}
		return next
	}
}
