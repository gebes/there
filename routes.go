package there

import (
	path2 "path"
	"strings"
)

type RouteGroup struct {
	*Router
	prefix string
}

func (group RouteGroup) Group(prefix string) *RouteGroup {

	prefix = strings.TrimPrefix(prefix, "/")

	group.assert(len(prefix) > 1, "route group needs to have at least one symbol")

	if !strings.HasSuffix(prefix, "/") {
		prefix += "/"
	}

	group.prefix += prefix
	return &group
}

func NewRouteGroup(router *Router, route string) *RouteGroup {

	router.assert(route != "", "route \""+route+"\" must not be empty")

	if !strings.HasPrefix(route, "/") {
		route = "/" + route
	}
	if !strings.HasSuffix(route, "/") {
		route += "/"
	}

	return &RouteGroup{
		Router: router,
		prefix: route,
	}
}

type Endpoint func(request Request) Response

// Route adds attributes to an Endpoint func
type Route struct {
	muxHandler *muxHandler
	methods    []method
}

func (group *RouteGroup) Handle(path string, endpoint Endpoint, methodsString ...string) *RouteRouteGroupBuilder {
	group.Router.mutex.Lock()
	defer group.Router.mutex.Unlock()

	var methods []method
	for _, m := range methodsString {
		methods = append(methods, methodToInt(m))
	}

	if path == "" {
		path = "/"
	}
	if path[0] == '/' {
		path = "/" + path
	}

	path = group.prefix + path
	path = path2.Clean(path)

	var ok bool
	var muxHandler *muxHandler
	muxHandler, ok = group.Router.handlerKeeper[path]
	if !ok {
		muxHandler = newMuxHandler(group.Router, endpoint)
		group.serveMux.Handle(path, muxHandler)
		group.Router.handlerKeeper[path] = muxHandler
	}

	for _, m := range methods {
		muxHandler.methods[m] = &muxHandlerEndpoint{
			endpoint: endpoint,
		}
	}

	route := &Route{
		muxHandler,
		methods,
	}

	return &RouteRouteGroupBuilder{
		route,
		group,
	}
}

type RouteRouteGroupBuilder struct {
	*Route
	*RouteGroup
}

func (group *RouteGroup) Get(route string, endpoint Endpoint) *RouteRouteGroupBuilder {
	return group.Handle(route, endpoint, MethodGet)
}

func (group *RouteGroup) Post(route string, endpoint Endpoint) *RouteRouteGroupBuilder {
	return group.Handle(route, endpoint, MethodPost)
}

func (group *RouteGroup) Patch(route string, endpoint Endpoint) *RouteRouteGroupBuilder {
	return group.Handle(route, endpoint, MethodPatch)
}

func (group *RouteGroup) Delete(route string, endpoint Endpoint) *RouteRouteGroupBuilder {
	return group.Handle(route, endpoint, MethodDelete)
}

func (group *RouteGroup) Connect(route string, endpoint Endpoint) *RouteRouteGroupBuilder {
	return group.Handle(route, endpoint, MethodConnect)
}

func (group *RouteGroup) Head(route string, endpoint Endpoint) *RouteRouteGroupBuilder {
	return group.Handle(route, endpoint, MethodHead)
}

func (group *RouteGroup) Trace(route string, endpoint Endpoint) *RouteRouteGroupBuilder {
	return group.Handle(route, endpoint, MethodTrace)
}

func (group *RouteGroup) Put(route string, endpoint Endpoint) *RouteRouteGroupBuilder {
	return group.Handle(route, endpoint, MethodPut)
}

func (group *RouteGroup) Options(route string, endpoint Endpoint) *RouteRouteGroupBuilder {
	return group.Handle(route, endpoint, MethodOptions)
}

// With adds a middleware to the handler the method is called on
func (group *RouteRouteGroupBuilder) With(middleware Middleware) *RouteRouteGroupBuilder {
	for _, method := range group.methods {
		endpoint := group.muxHandler.methods[method]
		endpoint.AddMiddleware(middleware)
	}
	return group
}
