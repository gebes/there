package middlewares

import (
	"fmt"
	. "github.com/Gebes/there/v2"
)

// RequireHost is a middleware, that forces incoming requests to have a specific host header set.
// If this header is not set, the String response with StatusBadRequest and the message
// "Invalid host for access to resource" is returned.
func RequireHost(host string) func(request Request, next Response) Response {
	return func(request Request, next Response) Response {
		if request.Host != host {
			return String(StatusBadRequest, fmt.Sprintf("Invalid host for access to resource"))
		}
		return next
	}
}

// RequireHeader is a middleware, that forces incoming requests to have a specific header.
// If this method is not used, only Status with status StatusBadRequest is returned.
func RequireHeader(key string, value string) func(request Request, next Response) Response {
	return func(request Request, next Response) Response {
		if header, exists := request.Headers.Get(key); !exists || header != value {
			return Status(StatusBadRequest)
		}
		return next
	}
}
