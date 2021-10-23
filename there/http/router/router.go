package there

import (
	. "github.com/Gebes/there/there/http/middlewares"
	"net/http"
)

type Router struct {
	RunningPort *string

	Endpoints         []*EndpointContainer
	GlobalMiddlewares []Middleware
}

func NewRouter() *Router {
	return &Router{}
}

func (router Router) Listen(address string) error {
	router.RunningPort = &address
	return http.ListenAndServe(address, &router)
}

func (router *Router) AddGlobalMiddleware(middleware Middleware) *Router {
	router.GlobalMiddlewares = append(router.GlobalMiddlewares, middleware)
	return router
}
