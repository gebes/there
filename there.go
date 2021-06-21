package there

import (
	"net/http"
	"strconv"
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

	// LogErrorsAlways defines if errors in a there.Response should be logged
	// or not, independent of the status code. Overrides Router.LogResponseBodies
	// and Router.LogRouteCalls
	LogErrorsAlways bool

	server   *http.Server
	handlers []HandlerContainer
}

//Listen start listening on the provided port blocking
func (router *Router) Listen() error {
	router.server = &http.Server{Addr: ":" + strconv.Itoa(router.Port), Handler: &GlobalHandler{router: router}}
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
