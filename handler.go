package there

import (
	"net/http"
	"path"
)

func (router *Router) ServeHTTP(rw http.ResponseWriter, request *http.Request) {
	handler, pattern := router.serveMux.Handler(request)
	if len(pattern) == 0 { // no handler was found
		router.Configuration.RouteNotFoundHandler(NewHttpRequest(rw, request)).ServeHTTP(rw, request)
	} else {
		handler.ServeHTTP(rw, request)
	}
}

// muxHandler defines a struct that encapsulates a handler and its middleware.
type muxHandler struct {
	router      *Router
	endpoint    Endpoint
	middlewares []Middleware
}

// newMuxHandler initializes and returns a new muxHandler.
func newMuxHandler(router *Router, endpoint Endpoint) *muxHandler {
	return &muxHandler{
		router:      router,
		endpoint:    endpoint,
		middlewares: []Middleware{},
	}
}

// AddMiddleware adds a new middleware to the handler's stack.
func (h *muxHandler) AddMiddleware(m Middleware) {
	h.middlewares = append(h.middlewares, m)
}

// ServeHTTP implements the http.Handler interface for muxHandler.
func (h *muxHandler) ServeHTTP(rw http.ResponseWriter, request *http.Request) {
	httpRequest := NewHttpRequest(rw, request)

	sanitizedPath := request.URL.Path
	if h.router.Configuration.SanitizePaths {
		sanitizedPath = path.Clean(sanitizedPath)
	}

	var next Response = ResponseFunc(func(rw http.ResponseWriter, r *http.Request) {
		h.endpoint(httpRequest).ServeHTTP(rw, r)
	})

	// Apply endpoint-specific middleware in reverse order.
	for i := len(h.middlewares) - 1; i >= 0; i-- {
		next = h.middlewares[i](httpRequest, next)
	}

	// Apply global middlewares in reverse order.
	for i := len(h.router.globalMiddlewares) - 1; i >= 0; i-- {
		next = h.router.globalMiddlewares[i](httpRequest, next)
	}

	next.ServeHTTP(rw, request)
}
