package there

import (
	"encoding/json"
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
		if !CheckArrayContains(AllMethods, method){
			panic(errors.New("there: Method " + method + " does not exist"))
		}
	}

	router.handlers = append(router.handlers, HandlerContainer{
		path: route,
		handler: handler,
		methods: methods,
	})
}

type GlobalHandler struct {
	router *Router
}

func (globalHandler *GlobalHandler) ServeHTTP(rawWriter http.ResponseWriter, rawRequest *http.Request) {
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

	router.logResponse(&request, response)


	// TODO extract into own method
	switch response.Data.(type) {
	case string, *string:
		_ = json.NewEncoder(rawWriter).Encode(response)
	case error, *error:
		response.Data = response.Data.(error).Error()
		_ = json.NewEncoder(rawWriter).Encode(response)
	default:
		_ = json.NewEncoder(rawWriter).Encode(response.Data)
	}
}

func (router*Router)logResponse(request*Request, response*Response)  {
	if router.LogRouteCalls {
		if router.LogResponseBodies {
			// FIXME
			log.Println(request.Method(), request.Path(), "resulted in", response.Status)
		}else{
			log.Println(request.Method(), request.Path(), "resulted in", response.Status)
		}
	}

}