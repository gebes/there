package there

import (
	"net/http"
)

func (router *Router) ServeHTTP(rw http.ResponseWriter, request *http.Request) {

	httpRequest := NewHttpRequest(rw, request)
	var middlewares = make([]Middleware, 0)
	middlewares = append(middlewares, router.globalMiddlewares...)

	var endpoint Endpoint = nil

	for _, current := range router.routes {
		routeParams, ok := current.Path.Parse(request.URL.Path)
		if ok && CheckArrayContains(current.Methods, request.Method) {
			endpoint = current.Endpoint
			middlewares = append(middlewares, current.Middlewares...)
			routeParamReader := RouteParamReader(routeParams)
			httpRequest.RouteParams = &routeParamReader
			break
		}
	}

	if endpoint == nil {
		endpoint = router.Configuration.RouteNotFoundHandler
	}

	var next HttpResponse = HttpResponseFunc(func(rw http.ResponseWriter, r *http.Request) {
		endpoint(httpRequest).ServeHTTP(rw, r)
	})
	for i := len(middlewares) - 1; i >= 0; i-- {
		middleware := middlewares[i]
		next = middleware(httpRequest, next)
	}
	next.ServeHTTP(rw, request)
}
