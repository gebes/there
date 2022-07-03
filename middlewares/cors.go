package middlewares

import (
	"github.com/Gebes/there/v2"
)

func Cors(configuration CorsConfiguration) there.Middleware {
	return func(request there.Request, next there.Response) there.Response {
		headers := map[string]string{
			there.ResponseHeaderAccessControlAllowOrigin:  configuration.AccessControlAllowOrigin,
			there.ResponseHeaderAccessControlAllowMethods: configuration.AccessControlAllowMethods,
			there.ResponseHeaderAccessControlAllowHeaders: configuration.AccessControlAllowHeaders,
		}
		if request.Method == there.MethodOptions {
			return there.Headers(headers, there.Status(there.StatusOK))
		}
		return there.Headers(headers, next)
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
		AccessControlAllowMethods: there.AllMethodsString,
		AccessControlAllowHeaders: "Accept, Content-Type, Content-Length, Authorization",
	}
}
