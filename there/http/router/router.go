package there

import (
	"fmt"
	. "github.com/Gebes/there/there/http/middlewares"
	"net/http"
)

type Router struct {
	RunningPort uint16

	Endpoints         []*EndpointContainer
	GlobalMiddlewares []Middleware
}

func NewRouter() *Router {
	return &Router{}
}

func (router Router) Listen(port uint16) error {
	router.RunningPort = port
	return http.ListenAndServe(fmt.Sprintf(":%d", port), &router)
}

func (router *Router) AddGlobalMiddleware(middleware Middleware) *Router {
	router.GlobalMiddlewares = append(router.GlobalMiddlewares, middleware)
	return router
}
