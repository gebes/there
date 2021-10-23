package there

import (
	"errors"
	. "github.com/Gebes/there/there/http/middlewares"
	. "github.com/Gebes/there/there/http/request"
	. "github.com/Gebes/there/there/http/response"
	. "github.com/Gebes/there/there/utils"
)

type Endpoint func(request HttpRequest) HttpResponse

type EndpointContainer struct {
	Endpoint    Endpoint
	Methods     []string
	Path        string
	Middlewares []Middleware
}

func (router *Router) handle(route string, endpoint Endpoint, methods ...string) *RouterEndpoint {

	if router.Endpoints == nil {
		router.Endpoints = make([]*EndpointContainer, 0)
	}

	for _, container := range router.Endpoints {
		if container.Path != route {
			continue
		}
		for _, containerMethod := range container.Methods {
			for _, method := range methods {
				if containerMethod == method {
					panic(errors.New("route already defined: " + containerMethod + " " + route))
				}
			}
		}
	}

	container := &EndpointContainer{
		endpoint,
		methods,
		route,
		make([]Middleware, 0),
	}
	router.Endpoints = append(router.Endpoints, container)

	return &RouterEndpoint{
		router,
		container,
	}
}

type RouterEndpoint struct {
	*Router
	*EndpointContainer
}

func (router *Router) Get(route string, endpoint Endpoint) *RouterEndpoint {
	return router.handle(route, endpoint, MethodGet)
}

func (router *Router) Post(route string, endpoint Endpoint) *RouterEndpoint {
	return router.handle(route, endpoint, MethodPost)
}

func (router *Router) Patch(route string, endpoint Endpoint) *RouterEndpoint {
	return router.handle(route, endpoint, MethodPatch)
}

func (router *Router) Delete(route string, endpoint Endpoint) *RouterEndpoint {
	return router.handle(route, endpoint, MethodDelete)
}

func (router *Router) Connect(route string, endpoint Endpoint) *RouterEndpoint {
	return router.handle(route, endpoint, MethodConnect)
}

func (router *Router) Head(route string, endpoint Endpoint) *RouterEndpoint {
	return router.handle(route, endpoint, MethodHead)
}

func (router *Router) Trace(route string, endpoint Endpoint) *RouterEndpoint {
	return router.handle(route, endpoint, MethodTrace)
}

func (router *Router) Put(route string, endpoint Endpoint) *RouterEndpoint {
	return router.handle(route, endpoint, MethodPut)
}

func (router *Router) Options(route string, endpoint Endpoint) *RouterEndpoint {
	return router.handle(route, endpoint, MethodOptions)
}

func (e *RouterEndpoint) AddMiddleware(middleware Middleware) *RouterEndpoint {
	e.Middlewares = append(e.Middlewares, middleware)
	return e
}
