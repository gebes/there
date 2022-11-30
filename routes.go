package there

import (
	"path/filepath"
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
	node    *node
	methods []method
}

func (group *RouteGroup) Handle(path string, endpoint Endpoint, methodsString ...Method) *RouteRouteGroupBuilder {
	var methods []method
	for _, m := range methodsString {
		methods = append(methods, methodToInt(m))
	}

	path = filepath.Clean(path)
	if path == "" {
		path = "/"
	}
	if path[0] == '/' {
		path = "/" + path
	}

	path = group.prefix + path

	for strings.Contains(path, "//") {
		path = strings.ReplaceAll(path, "//", "/")
	}

	node, err := group.Router.ensureNodeExists(path)
	if err != nil {
		group.Router.assertionErrors = append(group.Router.assertionErrors, err)
	}
	if node != nil {
		for _, method := range methods {
			group.Router.assert(node.handler[method] == nil, string(methodToString(method))+" "+path+" already defined")
			node.handler[method] = endpoint
		}

	}

	route := &Route{
		node,
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
		group.node.middlewares[method] = append(group.node.middlewares[method], middleware)
	}
	return group
}
