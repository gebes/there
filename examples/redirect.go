package examples

import (
	. "github.com/Gebes/there/there"
	"github.com/Gebes/there/there/middlewares/cors"
)

func RedirectRouter() error {
	router := Router{
		Port:              8080,
		LogRouteCalls:     true,
		LogResponseBodies: true,
		AlwaysLogErrors:   true,
	}

	router.AddGlobalMiddleware(cors.MiddlewareCors(cors.DefaultConfiguration()))

	router.HandleGet("/google", func(request Request) Response {
		return ResponseRedirect("https://google.com")
	})

	return router.Listen()
}


