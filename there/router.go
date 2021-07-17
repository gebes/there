package there

import (
	"net/http"
	"strconv"
	"strings"
)

type Router struct {

	// Port to listen
	Port int

	// LogRouteCalls defines if all accesses to a route should be logged
	//	2021/06/20 21:07:40 GET /user resulted in 200
	LogRouteCalls bool

	// LogResponseBodies defines if all the different responses from
	// all routes should be logged or not, provided that Router.LogRouteCalls
	// is enabled
	LogResponseBodies bool

	// AlwaysLogErrors defines if errors in a there.Response should be logged
	// or not, independent of the status code. Overrides Router.LogResponseBodies
	// and Router.LogRouteCalls
	AlwaysLogErrors bool

	// SetupResponseHeaders are automatically set, before a request gets forwarded
	// to its handler
	SetupResponseHeaders map[string]string


	// HandlerRouteNotFound gets called when a route gets called, which method is
	// not mapped
	HandlerRouteNotFound Handler

	server   *http.Server
	handlers []*HandlerContainer
	globalMiddlewares []Middleware
}



//Listen start listening on the provided port blocking
func (router *Router) Listen() error {
	router.server = &http.Server{Addr: ":" + strconv.Itoa(router.Port), Handler: &RouteHandler{router: router}}
	return router.server.ListenAndServe()
}

//IsRunning returns if the server is running, as long as the router.server object is not nil
func (router *Router) IsRunning() bool {
	return router.server != nil
}

//EnsureRunning panics if start hasn't been called, because the router cannot work if the http.Server is nil
func (router *Router) EnsureRunning() {
	if !router.IsRunning() {
		panic(ErrorNotRunning)
	}
}

func (router*Router) FindHandler(request*Request) *HandlerContainer{
	// find specific handler
	for _, container := range router.handlers {
		if container.path == request.URL.Path && (len(container.methods) == 0 || strings.Contains(strings.Join(container.methods, " "), request.Method)) {
			return container
		}
	}
	return nil
}

func (router*Router) AddGlobalMiddleware(middleware... Middleware){
	router.globalMiddlewares = append(router.globalMiddlewares, middleware...)
}