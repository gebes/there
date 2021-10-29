package there

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
)

type Router struct {
	RunningPort Port

	GlobalMiddlewares []Middleware

	//routes is a list of Routes which checks for duplicate entries
	//on insert.
	routes RouteManager
	mode   Mode

	*RouteGroup

	*Configuration

	HttpEngine
}

func NewRouter() *Router {
	r := &Router{
		RunningPort:       0,
		GlobalMiddlewares: make([]Middleware, 0),
		routes:            make([]*Route, 0),
		mode:              DebugMode,
		HttpEngine:        newDefaultHttpEngine(),
		Configuration:     defaultConfiguration(),
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
}

type defaultHttpEngine struct{}

func newDefaultHttpEngine() *defaultHttpEngine {
	return &defaultHttpEngine{}
}

func (d defaultHttpEngine) listenAndServe(addr string, handler http.Handler) error {
	return http.ListenAndServe(addr, handler)
}

func (d defaultHttpEngine) listenAndServeTLS(addr, certFile, keyFile string, handler http.Handler) error {
	return http.ListenAndServeTLS(addr, certFile, keyFile, handler)
}

func (router Router) Listen(port Port) error {
	router.RunningPort = port
	return router.listenAndServe(port.ToAddr(), &router)
}

func (router Router) ListenToTLS(port Port, certFile, keyFile string) error {
	router.RunningPort = port
	return router.listenAndServeTLS(port.ToAddr(), certFile, keyFile, &router)
}

//Use registers a Middleware
func (router *Router) Use(middleware Middleware) *Router {
	router.GlobalMiddlewares = append(router.GlobalMiddlewares, middleware)
	return router
}

//Configuration is a simple place for the user to override the behavior of the router
type Configuration struct {
	Serializers
	Handlers
}

type Serializers interface {
	ErrorToBytes(err error) ([]byte, error)
	ErrorToBytesContentType() string
}

func defaultConfiguration() *Configuration {
	return &Configuration{Serializers: &defaultSerializers{}, Handlers: &defaultHandlers{}}
}

type defaultSerializers struct{}

func (d defaultSerializers) ErrorToBytes(err error) ([]byte, error) {
	data, err := json.Marshal(map[string]interface{}{
		"error": err.Error(),
	})
	return data, err
}

func (d defaultSerializers) ErrorToBytesContentType() string {
	return ContentTypeApplicationJson
}

type Handlers interface {
	RouteNotFound(request *http.Request) error
}

type defaultHandlers struct{}

func (d defaultHandlers) RouteNotFound(request *http.Request) error {
	return errors.New("Could not find route " + request.Method + " " + request.URL.Path)
}

var mode = DebugMode

func IsDebug() bool {
	return mode == DebugMode
}

func IsProduction() bool {
	return mode == ProductionMode
}

func DebugPrintln(v ...interface{}) {
	if IsDebug() {
		log.Println(v...)
	}
}

func init() {
	SetMode(os.Getenv(envThereMode))
}

func SetMode(newMode string) {
	if newMode == "" {
		mode = DebugMode
		return
	}

	newMode = strings.ToLower(newMode)
	Assert(newMode == DebugMode || newMode == ProductionMode, "unknown mode: "+newMode)
	mode = newMode
}
