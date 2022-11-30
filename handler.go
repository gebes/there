package there

import (
	"net/http"
)

func (router *Router) ServeHTTP(rw http.ResponseWriter, request *http.Request) {
	httpRequest := NewHttpRequest(rw, request)
	method := methodToInt(request.Method)

	node, params := router.matcher.findNode(request.URL.Path)
	*httpRequest.RouteParams = params

	var endpoint Endpoint
	var middlewares []Middleware

	if node != nil {
		endpoint = node.handler[method]
		middlewares = node.middlewares[method]
	}
	if endpoint == nil {
		endpoint = router.Configuration.RouteNotFoundHandler
	}

	var next Response = ResponseFunc(func(rw http.ResponseWriter, r *http.Request) {
		endpoint(httpRequest).ServeHTTP(rw, r)
	})

	for i := len(middlewares) - 1; i >= 0; i-- {
		middleware := middlewares[i]
		next = middleware(httpRequest, next)
	}
	for i := len(router.globalMiddlewares) - 1; i >= 0; i-- {
		middleware := router.globalMiddlewares[i]
		next = middleware(httpRequest, next)
	}
	next.ServeHTTP(rw, request)
}
