package examples

import (
	. "github.com/Gebes/there/there"
	"log"
)

func MiddlewareRouter() error {
	router := Router{
		Port:              8080,
		LogRouteCalls:     true,
		LogResponseBodies: true,
		AlwaysLogErrors:   true,
	}

	router.HandleGet("/message", messageHandler).AddMiddleware(messageMiddleware, authorizationMiddleware)
	router.HandleGet("/message/without/middleware", messageHandler)

	return router.Listen()
}

func messageMiddleware(request Request) *Response {
	log.Println("First Middleware")
	return nil
}

func authorizationMiddleware(request Request) *Response{
	params, err := request.ReadParams("authorization")
	if err != nil {
		return ResponseDataP(StatusBadRequest, err)
	}
	if params[0] != "password" {
		return ResponseDataP(StatusUnauthorized, err)
	}
	return nil
}

func messageHandler(request Request) Response {
	return ResponseData(StatusOK, "Hello from the Handler")
}
