package there

import (
	. "github.com/Gebes/there/there/http/request"
	. "github.com/Gebes/there/there/http/response"
	. "github.com/Gebes/there/there/utils"
)

type Configuration struct {
	AccessControlAllowOrigin  string
	AccessControlAllowMethods string
	AccessControlAllowHeaders string
}

func AllowAllConfiguration() Configuration {
	return Configuration{
		AccessControlAllowOrigin:  "*",
		AccessControlAllowMethods: AllMethodsString,
		AccessControlAllowHeaders: "Accept, Content-Type, Content-Length, Authorization",
	}
}

func CorsMiddleware(configuration Configuration) Middleware {
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
