package there

import (
	"context"
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"github.com/vmihailenco/msgpack/v5"
	"gopkg.in/yaml.v3"
	"net/http"
)

type Router struct {
	GlobalMiddlewares []Middleware

	//routes is a list of Routes which checks for duplicate entries
	//on insert.
	routes RouteManager
	mode   Mode

	*RouteGroup

	*RouterConfiguration

	HttpEngine
}

func NewRouter() *Router {
	r := &Router{
		GlobalMiddlewares:   make([]Middleware, 0),
		routes:              make([]*Route, 0),
		mode:                DebugMode,
		HttpEngine:          newDefaultHttpEngine(),
		RouterConfiguration: defaultConfiguration(),
	}
	r.RouteGroup = NewRouteGroup(r, "/")
	return r
}

func (router *Router) IsProductionMode() bool {
	return router.mode.IsProduction()
}

func (router *Router) IsDebugMode() bool {
	return router.mode.IsDebug()
}

func (router *Router) SetProductionMode() *Router {
	router.mode.SetProduction()
	return router
}

func (router *Router) SetDebugMode() *Router {
	router.mode.SetDebug()
	return router
}

type Port uint16

func (p Port) ToAddr() string {
	return fmt.Sprintf(":%d", p)
}

type HttpEngine interface {
	listenAndServe(addr string, handler http.Handler) error
	listenAndServeTLS(addr, certFile, keyFile string, handler http.Handler) error
	shutdown(ctx context.Context) error
}

func (router Router) Listen(port Port) error {
	return router.listenAndServe(port.ToAddr(), &router)
}

func (router Router) ListenToTLS(port Port, certFile, keyFile string) error {
	return router.listenAndServeTLS(port.ToAddr(), certFile, keyFile, &router)
}

//Use registers a Middleware
func (router *Router) Use(middleware Middleware) *Router {
	router.GlobalMiddlewares = append(router.GlobalMiddlewares, middleware)
	return router
}

//RouterConfiguration is a simple place for the user to override the behavior of the router
type RouterConfiguration struct {
	*Serializers
	Handlers
}

type Serializers struct {
	JsonMarshal    func(interface{}) ([]byte, error)
	XmlMarshal     func(interface{}) ([]byte, error)
	YamlMarshal    func(interface{}) ([]byte, error)
	MsgpackMarshal func(interface{}) ([]byte, error)

	ErrorMarshal            func(interface{}) []byte
	ErrorMarshalContentType string
}

func defaultConfiguration() *RouterConfiguration {
	c := &RouterConfiguration{Serializers: &Serializers{
		JsonMarshal:             json.Marshal,
		XmlMarshal:              xml.Marshal,
		YamlMarshal:             yaml.Marshal,
		MsgpackMarshal:          msgpack.Marshal,
		ErrorMarshalContentType: ContentTypeApplicationJson,
	}, Handlers: &defaultHandlers{}}
	c.ErrorMarshal = func(err interface{}) []byte {
		var errorString string
		switch v := err.(type) {
		case error:
			errorString = v.Error()
		default:
			errorString = fmt.Sprint(v)
		}
		data, err := c.JsonMarshal(map[string]interface{}{
			"error": errorString,
		})
		if err != nil {
			data = []byte(errorString)
		}
		return data
	}
	return c
}

type Handlers interface {
	RouteNotFound(request *http.Request) error
}

type defaultHandlers struct{}

func (d defaultHandlers) RouteNotFound(request *http.Request) error {
	return errors.New("Could not find route " + request.Method + " " + request.URL.Path)
}

// Default Http Engine

type defaultHttpEngine struct {
	*http.Server
}

func newDefaultHttpEngine() *defaultHttpEngine {
	return &defaultHttpEngine{&http.Server{}}
}

func (d defaultHttpEngine) prepare(addr string, handler http.Handler) {
	d.Addr = addr
	d.Handler = handler
}

func (d defaultHttpEngine) listenAndServe(addr string, handler http.Handler) error {
	d.prepare(addr, handler)
	return http.ListenAndServe(addr, handler)
}

func (d defaultHttpEngine) listenAndServeTLS(addr, certFile, keyFile string, handler http.Handler) error {
	d.prepare(addr, handler)
	return http.ListenAndServeTLS(addr, certFile, keyFile, handler)
}

func (d defaultHttpEngine) shutdown(ctx context.Context) error {
	return d.Shutdown(ctx)
}
