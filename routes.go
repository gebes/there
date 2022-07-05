package there

import (
	"fmt"
	"strings"
)

type RouteGroup struct {
	*Router
	prefix string
}

func (group RouteGroup) Group(prefix string) *RouteGroup {

	prefix = strings.TrimPrefix(prefix, "/")

	Assert(len(prefix) > 1, "route group needs to have at least one symbol")

	if !strings.HasSuffix(prefix, "/") {
		prefix += "/"
	}

	return &RouteGroup{
		Router: group.Router,
		prefix: group.prefix + prefix,
	}
}

func NewRouteGroup(router *Router, route string) *RouteGroup {

	Assert(route != "", "route must not be empty")

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

//Route adds attributes to an Endpoint func
type Route struct {
	Endpoint    Endpoint
	Methods     []string
	Path        Path
	Middlewares []Middleware
}

//OverlapsWith checks if an Route somehow overlaps with another container. For this to be true, the path and at least one method must equal
func (e Route) OverlapsWith(toCompare Route) bool {
	if !e.Path.Equals(toCompare.Path) {
		return false
	}
	return CheckArraysOverlap(e.Methods, toCompare.Methods)
}

func (e Route) ToString() string {
	r := fmt.Sprint(e.Methods, " ", e.Path.ToString())
	if e.Path.ignoreCase {
		r += " *IgnoreCase"
	}
	return r
}

func (group *RouteGroup) Handle(path string, endpoint Endpoint, methods ...string) *RouteRouteGroupBuilder {
	if path == "" {
		path = "/"
	}
	if path[0] != '/' {
		path = "/" + path
	}

	path = group.prefix + path

	if group.Router.routes == nil {
		group.Router.routes = make([]*Route, 0)
	}

	route := &Route{
		endpoint,
		methods,
		ConstructPath(path, false),
		make([]Middleware, 0),
	}
	group.routes.AddRoute(route)

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

//With adds a middleware to the handler the method is called on
func (group *RouteRouteGroupBuilder) With(middleware Middleware) *RouteRouteGroupBuilder {
	group.Middlewares = append(group.Middlewares, middleware)
	return group
}

func (group *RouteRouteGroupBuilder) IgnoreCase() *RouteRouteGroupBuilder {
	// cancel if already ignore case
	if group.Route.Path.ignoreCase {
		return group
	}

	// delete route
	group.routes.RemoveRoute(group.Route)

	group.Route.Path.ignoreCase = true
	group.routes.AddRoute(group.Route)

	return group
}

type RouteManager []*Route

func (r *RouteManager) FindOverlappingRoute(routeToCheck *Route) *Route {
	for _, toCompare := range *r {
		if toCompare.OverlapsWith(*routeToCheck) {
			return toCompare
		}
	}
	return nil
}

func (r *RouteManager) AddRoute(routeToAdd *Route) *Route {

	if overlapsWith := r.FindOverlappingRoute(routeToAdd); overlapsWith != nil {
		panic("The route \"" + routeToAdd.ToString() + "\" overlaps with the existing route \"" + overlapsWith.ToString() + "\"")
	}

	*r = append(*r, routeToAdd)

	return routeToAdd
}

func (r *RouteManager) RemoveRoute(toRemove *Route) {

	for i, container := range *r {
		if container.Path.Equals(toRemove.Path) {
			*r = append((*r)[:i], (*r)[i+1:]...)
		}
	}

}
