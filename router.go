package there

import (
	"errors"
	"fmt"
	"github.com/gebes/there/v2/status"
	"net/http"
	"sync"
)

type Router struct {
	*RouteGroup
	Configuration *RouterConfiguration
	Server        *http.Server

	assertionErrors

	globalMiddlewares []Middleware

	// routes is a list of Routes which checks for duplicate entries
	// on insert.
	serveMux      *http.ServeMux
	handlerKeeper map[string]*muxHandler
	mutex         sync.Mutex
}

func NewRouter() *Router {
	r := &Router{
		globalMiddlewares: make([]Middleware, 0),
		RouteGroup:        &RouteGroup{prefix: "/"},
		Server:            &http.Server{},
		Configuration: &RouterConfiguration{
			RouteNotFoundHandler: func(request Request) Response {
				type Error struct {
					Error  string `json:"error,omitempty" xml:"Error" yaml:"error" bson:"error"`
					Path   string `json:"path,omitempty" xml:"Path" yaml:"path" bson:"path"`
					Method string `json:"method,omitempty" xml:"Method" yaml:"method" bson:"method"`
				}
				return Auto(status.NotFound, Error{
					Error:  "could not find specified path",
					Path:   request.Request.URL.Path,
					Method: request.Request.Method,
				})
			},
			SanitizePaths: true,
		},
		serveMux:      http.NewServeMux(),
		handlerKeeper: map[string]*muxHandler{},
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
	err := router.HasError()
	if err != nil {
		return err
	}
	router.Server.Addr = port.ToAddr()
	return router.Server.ListenAndServe()
}

func (router *Router) ListenToTLS(port Port, certFile, keyFile string) error {
	err := router.HasError()
	if err != nil {
		return err
	}
	router.Server.Addr = port.ToAddr()
	return router.Server.ListenAndServeTLS(certFile, keyFile)
}

// Use registers a Middleware
func (router *Router) Use(middleware ...Middleware) *Router {
	router.globalMiddlewares = append(router.globalMiddlewares, middleware...)
	return router
}

// RouterConfiguration is a straightforward place to override default behavior of the router
type RouterConfiguration struct {
	// RouteNotFoundHandler gets invoked, when the specified URL and method have no handlers
	RouteNotFoundHandler Endpoint
	SanitizePaths        bool
}

type assertionErrors []error

func (a *assertionErrors) HasError() error {
	if len(*a) == 0 {
		return nil
	}
	var err error
	err, *a = (*a)[0], (*a)[1:]
	return err
}

func (a *assertionErrors) assert(condition bool, errorString string) {
	if !condition {
		*a = append(*a, errors.New(errorString))
	}
}
