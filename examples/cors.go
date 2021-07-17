package examples

import (
	. "github.com/Gebes/there/there"
	"github.com/Gebes/there/there/middlewares/cors"
)

func CorsRouter() error {
	router := Router{
		Port:              8080,
		LogRouteCalls:     true,
		LogResponseBodies: true,
		AlwaysLogErrors:   true,
	}

	router.AddGlobalMiddleware(cors.MiddlewareCors(cors.DefaultConfiguration()))

	router.HandleGet("/cors", func(request Request) Response {
		return ResponseData(StatusOK, "Hello there")
	})


	return router.Listen()
}


