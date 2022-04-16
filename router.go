package there

import (
	"errors"
	"fmt"
	"net/http"
)

type Router struct {
	*RouteGroup
	Configuration *RouterConfiguration

	*http.Server

	globalMiddlewares []Middleware

	//routes is a list of Routes which checks for duplicate entries
	//on insert.
	routes RouteManager
}

func NewRouter() *Router {
	r := &Router{
		globalMiddlewares: make([]Middleware, 0),
		routes:            make([]*Route, 0),
		Server:            &http.Server{},
		Configuration: &RouterConfiguration{
			RouteNotFoundHandler: func(request HttpRequest) HttpResponse {
				return Error(StatusNotFound, errors.New("could not find route "+request.Method+" "+request.Request.URL.Path))
			},
		},
	}
	r.Server.Handler = r
	r.RouteGroup = NewRouteGroup(r, "/")
	return r
}

type Port uint16

func (p Port) ToAddr() string {
	return fmt.Sprintf(":%d", p)
}

func (router *Router) Listen(port Port) error {
	router.Server.Addr = port.ToAddr()
	return router.Server.ListenAndServe()
}

func (router *Router) ListenToTLS(port Port, certFile, keyFile string) error {
	router.Server.Addr = port.ToAddr()
	return router.Server.ListenAndServeTLS(certFile, keyFile)
}

//Use registers a Middleware
func (router *Router) Use(middleware Middleware) *Router {
	router.globalMiddlewares = append(router.globalMiddlewares, middleware)
	return router
}

//RouterConfiguration is a straightforward place to override default behavior of the router
type RouterConfiguration struct {
	//RouteNotFoundHandler gets invoked, when the specified URL and method have no handlers
	RouteNotFoundHandler Endpoint
}
