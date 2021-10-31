package there

import (
	"errors"
	"net/http"
)

type Middleware func(request HttpRequest) HttpResponse

//Next if returned will continue to the next middleware or the response
func Next() *nextMiddleware {
	r := &nextMiddleware{}
	r.header = &HeaderWrapper{Values: map[string][]string{}, HttpResponse: r}
	return r
}

type nextMiddleware struct {
	header *HeaderWrapper
}

func (j nextMiddleware) Execute(*Router, *http.Request, *http.ResponseWriter) error {
	return errors.New("cannot execute next middleware")
}

func (j *nextMiddleware) Header() *HeaderWrapper {
	return j.header
}

func CorsMiddleware(configuration CorsConfiguration) Middleware {
	return func(request HttpRequest) HttpResponse {
		// we need to unpack the ResponseWriter first
		headers := map[string][]string{
			ResponseHeaderAccessControlAllowOrigin:  {configuration.AccessControlAllowOrigin},
			ResponseHeaderAccessControlAllowMethods: {configuration.AccessControlAllowMethods},
			ResponseHeaderAccessControlAllowHeaders: {configuration.AccessControlAllowHeaders},
		}

		if request.Method == MethodOptions {
			return Empty(StatusOK).
				Header().SetAll(headers)
		}
		return Next().
			Header().SetAll(headers)
	}
}

type CorsConfiguration struct {
	AccessControlAllowOrigin  string
	AccessControlAllowMethods string
	AccessControlAllowHeaders string
}

func AllowAllConfiguration() CorsConfiguration {
	return CorsConfiguration{
		AccessControlAllowOrigin:  "*",
		AccessControlAllowMethods: AllMethodsString,
		AccessControlAllowHeaders: "Accept, Content-Type, Content-Length, Authorization",
	}
}
