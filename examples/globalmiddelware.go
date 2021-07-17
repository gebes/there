package examples

import (
	. "github.com/Gebes/there/there"
	"log"
)


func GlobalMiddlewareRouter()error {
	router := Router{
		Port:              8080,
		LogRouteCalls:     true,
		LogResponseBodies: true,
		AlwaysLogErrors:   true,
	}

	router.AddGlobalMiddleware(func(request Request) *Response {
		log.Println("First middleware")
		return nil
	})

	router.AddGlobalMiddleware(func(request Request) *Response {
		log.Println("Second middleware")
		return ResponseDataP(StatusForbidden, "I won't let you in!")
	})

	router.HandleGet("/message", func(request Request) Response {
		return ResponseData(StatusOK, "Hello from the Handler")
	})


	return router.Listen()
}
