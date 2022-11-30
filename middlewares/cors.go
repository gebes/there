package middlewares

import (
	"github.com/Gebes/there/v2"
	"github.com/Gebes/there/v2/header"
	"github.com/Gebes/there/v2/status"
)

func Cors(configuration CorsConfiguration) there.Middleware {
	return func(request there.Request, next there.Response) there.Response {
		headers := map[string]string{
			header.ResponseAccessControlAllowOrigin:  configuration.AccessControlAllowOrigin,
			header.ResponseAccessControlAllowMethods: configuration.AccessControlAllowMethods,
			header.ResponseAccessControlAllowHeaders: configuration.AccessControlAllowHeaders,
		}
		if request.Method == there.MethodOptions {
			return there.Headers(headers, there.Status(status.OK))
		}
		return there.Headers(headers, next)
	}
}

type CorsConfiguration struct {
	AccessControlAllowOrigin  string
	AccessControlAllowMethods string
	AccessControlAllowHeaders string
}

func CorsAllowAllConfiguration() CorsConfiguration {
	return CorsConfiguration{
		AccessControlAllowOrigin:  "*",
		AccessControlAllowMethods: there.AllMethodsJoined,
		AccessControlAllowHeaders: "Accept, Content-Type, Content-Length, Authorization",
	}
}
