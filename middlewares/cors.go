package middlewares

import . "github.com/Gebes/there/v2"

func Cors(configuration CorsConfiguration) Middleware {
	return func(request HttpRequest, next HttpResponse) HttpResponse {
		headers := MapString{
			ResponseHeaderAccessControlAllowOrigin:  configuration.AccessControlAllowOrigin,
			ResponseHeaderAccessControlAllowMethods: configuration.AccessControlAllowMethods,
			ResponseHeaderAccessControlAllowHeaders: configuration.AccessControlAllowHeaders,
		}
		if request.Method == MethodOptions {
			return WithHeaders(headers, Status(StatusOK))
		}
		return WithHeaders(headers, next)
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

