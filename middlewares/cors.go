package middlewares

import . "github.com/Gebes/there/v2"

func Cors(configuration CorsConfiguration) Middleware {
	return func(request Request, next Response) Response {
		headers := map[string]string{
			ResponseHeaderAccessControlAllowOrigin:  configuration.AccessControlAllowOrigin,
			ResponseHeaderAccessControlAllowMethods: configuration.AccessControlAllowMethods,
			ResponseHeaderAccessControlAllowHeaders: configuration.AccessControlAllowHeaders,
		}
		if request.Method == MethodOptions {
			return Headers(headers, Status(StatusOK))
		}
		return Headers(headers, next)
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
