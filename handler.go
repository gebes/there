package there

import (
	"errors"
	"log"
	"net/http"
	"strings"
)

type (
	Handler func(request Request) Response
)

//HandlerContainer contains the Handler Func and the allowed methods. If the methods slice is empty, then all methods are allowed.
type HandlerContainer struct {
	path    string
	handler Handler
	methods []string
}

//Handle registers a HandlerContainer with it's allowed methods. If no http method is provided, then all methods are allowed
func (router *Router) Handle(route string, handler Handler, methods ...string) {
	if router.handlers == nil {
		router.handlers = make([]HandlerContainer, 1)
	}

	for _, method := range methods {
		// check if the method exists
		if !CheckArrayContains(AllMethods, method) {
			panic(errors.New("there: Method " + method + " does not exist"))
		}
	}

	router.handlers = append(router.handlers, HandlerContainer{
		path:    route,
		handler: handler,
		methods: methods,
	})
}

//HandleGet register Get request
func (router *Router) HandleGet(route string, handler Handler) {
	router.Handle(route, handler, MethodGet)
}

//HandleHead register Head request
func (router *Router) HandleHead(route string, handler Handler) {
	router.Handle(route, handler, MethodHead)
}

//HandlePost register Post request
func (router *Router) HandlePost(route string, handler Handler) {
	router.Handle(route, handler, MethodPost)
}

//HandlePut register Put request
func (router *Router) HandlePut(route string, handler Handler) {
	router.Handle(route, handler, MethodPut)
}

//HandlePatch register Patch request
func (router *Router) HandlePatch(route string, handler Handler) {
	router.Handle(route, handler, MethodPatch)
}

//HandleDelete register v request
func (router *Router) HandleDelete(route string, handler Handler) {
	router.Handle(route, handler, MethodDelete)
}

//HandleConnect register Connect request
func (router *Router) HandleConnect(route string, handler Handler) {
	router.Handle(route, handler, MethodConnect)
}

//HandleOptions register Options request
func (router *Router) HandleOptions(route string, handler Handler) {
	router.Handle(route, handler, MethodOptions)
}

//HandleTrace register Trace request
func (router *Router) HandleTrace(route string, handler Handler) {
	router.Handle(route, handler, MethodTrace)
}

type RouteHandler struct {
	router *Router
}

func (globalHandler *RouteHandler) ServeHTTP(rawWriter http.ResponseWriter, rawRequest *http.Request) {
	var response *Response

	request := Request{
		RawRequest: rawRequest,
		RawWriter:  &rawWriter,
	}
	router := globalHandler.router

	for _, container := range router.handlers {
		if !(container.path == request.Path() && (len(container.methods) == 0 || strings.Contains(strings.Join(container.methods, " "), rawRequest.Method))) {
			continue
		}

		temp := container.handler(request)
		response = &temp
		break
	}

	if response == nil {
		response = &Response{Status: http.StatusNotFound}
	}

	// default Content-Type
	rawWriter.Header().Set("Content-Type", "application/json")
	rawWriter.WriteHeader(response.Status)

	if response.headers != nil {
		for key, value := range response.headers {
			rawWriter.Header().Set(key, value)
		}
	}

	_, err := rawWriter.Write(response.ToJson())
	if err != nil {
		return
	}

	router.logResponse(&request, response)
}

func (router *Router) logResponse(request *Request, response *Response) {
	if router.LogRouteCalls {
		if (router.AlwaysLogErrors && response.IsError()) || router.LogResponseBodies {
			log.Println(request.Method(), request.Path(), "resulted in", response.Status, "with body", string(response.ToJson()))
		} else {
			log.Println(request.Method(), request.Path(), "resulted in", response.Status)
		}
	}

}
