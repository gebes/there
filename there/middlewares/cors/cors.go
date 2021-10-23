package cors

type Configuration struct {
	AccessControlAllowOrigin  string
	AccessControlAllowMethods string
	AccessControlAllowHeaders string
}
/*
func DefaultConfiguration() Configuration {
	return Configuration{
		AccessControlAllowOrigin:  "*",
		AccessControlAllowMethods: there.AllMethodsString,
		AccessControlAllowHeaders: "Accept, Content-Type, Content-Length, Authorization",
	}
}

func MiddlewareCors(configuration Configuration) there.Middleware {
	return func(request there.Request) *there.Response {
		// we need to unpack the ResponseWriter first
		rw := *request.ResponseWriter
		rw.Header().Set(there.ResponseHeaderAccessControlAllowOrigin, configuration.AccessControlAllowOrigin)
		rw.Header().Set(there.ResponseHeaderAccessControlAllowMethods, configuration.AccessControlAllowMethods)
		rw.Header().Set(there.ResponseHeaderAccessControlAllowHeaders, configuration.AccessControlAllowHeaders)
		if request.Method == there.MethodOptions {
			return there.ResponseStatusP(there.StatusOK)
		}
		return nil
	}
}
*/