package there

import (
	"errors"
	"log"
	"net/http"
)

type (

	//Middleware can be called before a specific Handler or before any Handler (GlobalMiddleware)
	//If a Middleware does not return a nil Response, then this Response will be returned and all
	//following Middlewares or Handlers will not be called.
	Middleware func(request Request) *Response

	//Handler the final handler of a request. Is at the bottom
	Handler func(request Request) Response
)

//HandlerContainer contains the Handler Func and the allowed methods. If the methods slice is empty, then all methods are allowed.
type HandlerContainer struct {
	path        string
	handler     Handler
	methods     []string
	middlewares []Middleware
}

//Handle registers a HandlerContainer with it's allowed methods. If no http method is provided, then all methods are allowed
func (router *Router) Handle(route string, handler Handler, methods ...string) *HandlerContainer {
	if router.handlers == nil {
		router.handlers = make([]*HandlerContainer, 0)
	}

	for _, method := range methods {
		// check if the method exists
		if !CheckArrayContains(AllMethods, method) {
			panic(errors.New("there: Method " + method + " does not exist"))
		}
	}

	handlerContainer := HandlerContainer{
		path:    route,
		handler: handler,
		methods: methods,
	}
	router.handlers = append(router.handlers, &handlerContainer)
	return &handlerContainer
}

//AddMiddleware allows you to register a Middleware on this Route
func (h *HandlerContainer) AddMiddleware(middleware ...Middleware) *HandlerContainer {
	h.middlewares = append(h.middlewares, middleware...)
	return h
}

//HandleGet register Get request
func (router *Router) HandleGet(route string, handler Handler) *HandlerContainer {
	return router.Handle(route, handler, MethodGet)
}

//HandleHead register Head request
func (router *Router) HandleHead(route string, handler Handler) *HandlerContainer {
	return router.Handle(route, handler, MethodHead)
}

//HandlePost register Post request
func (router *Router) HandlePost(route string, handler Handler) *HandlerContainer {
	return router.Handle(route, handler, MethodPost)
}

//HandlePut register Put request
func (router *Router) HandlePut(route string, handler Handler) *HandlerContainer {
	return router.Handle(route, handler, MethodPut)
}

//HandlePatch register Patch request
func (router *Router) HandlePatch(route string, handler Handler) *HandlerContainer {
	return router.Handle(route, handler, MethodPatch)
}

//HandleDelete register v request
func (router *Router) HandleDelete(route string, handler Handler) *HandlerContainer {
	return router.Handle(route, handler, MethodDelete)
}

//HandleConnect register Connect request
func (router *Router) HandleConnect(route string, handler Handler) *HandlerContainer {
	return router.Handle(route, handler, MethodConnect)
}

//HandleOptions register Options request
func (router *Router) HandleOptions(route string, handler Handler) *HandlerContainer {
	return router.Handle(route, handler, MethodOptions)
}

//HandleTrace register Trace request
func (router *Router) HandleTrace(route string, handler Handler) *HandlerContainer {
	return router.Handle(route, handler, MethodTrace)
}

//HandlerRouteNotFound 404
func HandlerRouteNotFound(request Request) Response {
	return ResponseData(StatusNotFound, "Route "+request.Method+" "+request.URL.Path+" does not exist")
}

type RouteHandler struct {
	router *Router
}

func (routeHandler *RouteHandler) ServeHTTP(responseWriter http.ResponseWriter, rawRequest *http.Request) {

	request := &Request{
		*rawRequest,
		&responseWriter,
	}
	router := routeHandler.router

	// setup headers
	for key, value := range router.SetupResponseHeaders {
		responseWriter.Header().Set(key, value)
	}

	// default Content-Type
	// can be overridden by a middleware
	responseWriter.Header().Set(ResponseHeaderContentType, ContentTypeApplicationJson)

	var response *Response

	for _, middleware := range router.globalMiddlewares {
		middlewareResponse := middleware(*request)
		if middlewareResponse != nil {
			response = middlewareResponse
			break
		}
	}

	// if the GlobalMiddleware did not produce a response, then we get the
	// handler of the route...
	if response == nil {
		handlerContainer := router.FindHandler(request)
		if handlerContainer == nil {
			// looks like we found no handler -> 404
			var handlerResponse Response
			if router.HandlerRouteNotFound == nil {
				handlerResponse = HandlerRouteNotFound(*request)
			} else {
				handlerResponse = router.HandlerRouteNotFound(*request)
			}
			response = &handlerResponse
		} else {
			// we got a handler, go through all middlewares
			for _, middleware := range handlerContainer.middlewares {
				middlewareResponse := middleware(*request)
				if middlewareResponse != nil {
					response = middlewareResponse
					break
				}
			}

			// if the middlewares did not produce a output
			// then finally call the handler
			if response == nil {
				handlerResponse := handlerContainer.handler(*request)
				response = &handlerResponse
			}

		}
	}

	if len(response.redirect) != 0 {
		http.Redirect(responseWriter, rawRequest, response.redirect, response.Status)
	} else {
		responseWriter.WriteHeader(response.Status)
	}
	if response.headers != nil {
		for key, value := range response.headers {
			responseWriter.Header().Set(key, value)
		}
	}

	_, err := responseWriter.Write(response.ToJson())
	if err != nil {
		return
	}

	router.logResponse(request, response)
}

func (router *Router) logResponse(request *Request, response *Response) {
	if router.LogRouteCalls {
		if len(response.redirect) != 0 {
			log.Println(request.Method, request.URL.Path, "resulted in", response.Status, "redirecting to", response.redirect)
		} else if (router.AlwaysLogErrors && response.IsError()) || router.LogResponseBodies {
			log.Println(request.Method, request.URL.Path, "resulted in", response.Status, "with body", string(response.ToJson()))
		} else {
			log.Println(request.Method, request.URL.Path, "resulted in", response.Status)
		}
	}

}
