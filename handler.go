package there

import (
	"net/http"
)

func (router *Router) ServeHTTP(rw http.ResponseWriter, request *http.Request) {
	writer := &rw

	write := func(response HttpResponse) {
		writeHeader(writer, response)
		_ = response.Execute(router, request, writer)
	}

	defer func() {
		err := recover()
		if err == nil {
			return
		}
		write(Error(StatusInternalServerError, err))
	}()

	httpRequest := NewHttpRequest(request, writer)

	for _, middleware := range router.GlobalMiddlewares {
		response := middleware(httpRequest)
		write(response)
		if !isNextResponse(response) {
			return
		}
	}

	var endpoint Endpoint = nil
	var middlewares = make([]Middleware, 0)

	for _, route := range router.routes {
		routeParams, ok := route.Path.Parse(request.URL.Path)
		if ok && CheckArrayContains(route.Methods, request.Method) {
			endpoint = route.Endpoint
			middlewares = route.Middlewares
			routeParamReader := RouteParamReader(routeParams)
			httpRequest.RouteParams = &routeParamReader
			break
		}
	}

	for _, middleware := range middlewares {
		response := middleware(httpRequest)
		write(response)
		if !isNextResponse(response) {
			return
		}
	}

	if endpoint == nil {
		endpoint = func(_ HttpRequest) HttpResponse {
			return Error(StatusNotFound, router.RouterConfiguration.RouteNotFound(request))
		}
	}

	httpResponse := endpoint(httpRequest)
	writeHeader(writer, httpResponse)

	err := httpResponse.Execute(router, request, writer)
	if err != nil {
		write(Error(StatusInternalServerError, err))
	}

}

func writeHeader(w *http.ResponseWriter, httpResponse HttpResponse) {
	for key, values := range httpResponse.Header().Values {
		(*w).Header().Del(key)
		for _, value := range values {
			(*w).Header().Add(key, value)
		}
	}
}

func isNextResponse(response HttpResponse) bool {
	switch v := response.(type) {
	case *nextMiddleware:
		return true
	case *HeaderWrapper:
		switch v.HttpResponse.(type) {
		case *nextMiddleware:
			return true
		default:
			return false
		}
	case *contextResponse:
		for {
			switch res := v.response.(type) {
			case *contextResponse:
				v = res
			case *nextMiddleware:
				return true
			default:
				return false
			}
		}
	default:
		return false
	}
}
