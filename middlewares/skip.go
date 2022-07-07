package middlewares

import "github.com/Gebes/there/v2"



func Skip(handler there.Middleware, exclude func(request there.Request) bool) there.Middleware {
	if exclude == nil {
		return handler
	}

	fn := func(request there.Request, next there.Response) there.Response {
		if exclude(request) {
			return next
		}

		return handler(request, next)
	}

	return fn
}
