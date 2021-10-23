package there

import (
	. "github.com/Gebes/there/there/http/middlewares"
	. "github.com/Gebes/there/there/http/request"
	. "github.com/Gebes/there/there/http/response"
	"github.com/Gebes/there/there/http/router/handlers"
	. "github.com/Gebes/there/there/utils"
	"net/http"
)

func (router *Router) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	route := request.URL.Path
	method := request.Method

	httpRequest := NewHttpRequest(*request)
	var httpResponse HttpResponse

	errorOut := func(err error) {
		httpResponse = Error(StatusInternalServerError, err)
		writeHeader(&writer, httpResponse)
		err = httpResponse.Execute(request, &writer)
	}

	for _, middleware := range router.GlobalMiddlewares {
		httpResponse = middleware(httpRequest)
		writeHeader(&writer, httpResponse)
		if !isNextMiddleware(httpResponse) {
			break
		}
	}

	var endpoint Endpoint = handlers.RouteNotFound
	var middlewares = make([]Middleware, 0)

	for _, container := range router.Endpoints {
		if container.Path != route {
			continue
		}
		if CheckArrayContains(container.Methods, method) {
			endpoint = container.Endpoint
			middlewares = container.Middlewares
			break
		}
	}

	for _, middleware := range middlewares {
		httpResponse = middleware(httpRequest)
		writeHeader(&writer, httpResponse)
		if !isNextMiddleware(httpResponse) {
			break
		}
	}

	if httpResponse == nil || isNextMiddleware(httpResponse) {
		httpResponse = endpoint(httpRequest)
		writeHeader(&writer, httpResponse)
	}

	err := httpResponse.Execute(request, &writer)
	if err != nil {
		errorOut(err)
		return
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

func isNextMiddleware(response HttpResponse) bool {

	switch v := response.(type) {
	case *NextMiddleware:
		return true
	case *HeaderWrapper:
		switch v.HttpResponse.(type){
		case *NextMiddleware:
			return true
		default:
			return false
		}
	default:
		return false
	}
}
