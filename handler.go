package there

import (
	"net/http"
	"path"
)

func (router *Router) ServeHTTP(rw http.ResponseWriter, request *http.Request) {
	_, pattern := router.serveMux.Handler(request)
	if len(pattern) == 0 { // no handler was found
		// not found with global middlewares applied
		wrappedNotFoundHandler := router.applyGlobalMiddlewares(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
			router.Configuration.RouteNotFoundHandler(NewHttpRequest(rw, req)).ServeHTTP(rw, req)
		}))
		wrappedNotFoundHandler.ServeHTTP(rw, request)
	} else {
		router.serveMux.ServeHTTP(rw, request)
	}
}

// muxHandler defines a struct that encapsulates a handler and its middleware.
type (
	muxHandler struct {
		router  *Router
		methods map[method]*muxHandlerEndpoint
	}
	muxHandlerEndpoint struct {
		endpoint    Endpoint
		middlewares []Middleware
	}
)

// newMuxHandler initializes and returns a new muxHandler.
func newMuxHandler(router *Router, endpoint Endpoint) *muxHandler {
	return &muxHandler{
		router:  router,
		methods: map[method]*muxHandlerEndpoint{},
	}
}

// AddMiddleware adds a new middleware to the handler's stack.
func (h *muxHandlerEndpoint) AddMiddleware(m Middleware) {
	h.middlewares = append(h.middlewares, m)
}

// ServeHTTP implements the http.Handler interface for muxHandler.
func (h *muxHandler) ServeHTTP(rw http.ResponseWriter, request *http.Request) {
	httpRequest := NewHttpRequest(rw, request)
	method := methodToInt(request.Method)

	sanitizedPath := request.URL.Path
	if h.router.Configuration.SanitizePaths {
		sanitizedPath = path.Clean(sanitizedPath)
	}

	muxHandlerEndpoint, ok := h.methods[method]
	if !ok {
		// not found with global middlewares applied
		h.router.applyGlobalMiddlewares(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
			h.router.Configuration.RouteNotFoundHandler(httpRequest).ServeHTTP(rw, req)
		})).ServeHTTP(rw, request)
		return
	}
	endpoint, middlewares := muxHandlerEndpoint.endpoint, muxHandlerEndpoint.middlewares

	var next Response = ResponseFunc(func(rw http.ResponseWriter, r *http.Request) {
		endpoint(httpRequest).ServeHTTP(rw, r)
	})

	// Apply endpoint-specific middleware in reverse order.
	for i := len(middlewares) - 1; i >= 0; i-- {
		next = middlewares[i](httpRequest, next)
	}

	// Apply global middlewares in reverse order.
	for i := len(h.router.globalMiddlewares) - 1; i >= 0; i-- {
		next = h.router.globalMiddlewares[i](httpRequest, next)
	}

	next.ServeHTTP(rw, request)
}

// applyGlobalMiddlewares wraps the given http.Handler with the router's global middlewares.
func (router *Router) applyGlobalMiddlewares(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, request *http.Request) {
		httpRequest := NewHttpRequest(rw, request)
		var next Response = ResponseFunc(func(rw http.ResponseWriter, r *http.Request) {
			handler.ServeHTTP(rw, r)
		})

		// Apply global middlewares in reverse order.
		for i := len(router.globalMiddlewares) - 1; i >= 0; i-- {
			next = router.globalMiddlewares[i](httpRequest, next)
		}

		next.ServeHTTP(rw, request)
	})
}
