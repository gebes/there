package main

import (
	"errors"
	"github.com/Gebes/there/v2/header"
	"github.com/Gebes/there/v2/status"
	"log"

	"github.com/Gebes/there/v2"
	"github.com/Gebes/there/v2/middlewares"
)

func main() {
	router := there.NewRouter()

	// Register global middlewares
	router.Use(middlewares.Recoverer)
	router.Use(middlewares.RequireHost("localhost:8080"))
	router.Use(middlewares.Cors(middlewares.CorsAllowAllConfiguration()))
	router.Use(GlobalMiddleware)

	router.
		Get("/", Get).With(RouteSpecificMiddleware) // Register route with middleware

	err := router.Listen(8080)
	if err != nil {
		log.Fatalln("Could not listen to 8080:", err)
	}
}

func GlobalMiddleware(request there.Request, next there.Response) there.Response {
	// Check the request content-type
	if request.Headers.GetDefault(header.ContentType, "") != there.ContentTypeApplicationJson {
		return there.Error(status.UnsupportedMediaType, errors.New("header "+header.ContentType+" is not "+there.ContentTypeApplicationJson))
	}

	return next // Everything is fine until here, continue
}

func RouteSpecificMiddleware(request there.Request, next there.Response) there.Response {
	return there.Headers(map[string]string{
		header.ResponseContentLanguage: "en",
	}, next) // Set the content-language header by wrapping next with Headers
}

func Get(request there.Request) there.Response {
	return there.Json(status.OK, map[string]string{
		"Hello": "World",
		"How":   "are you?",
	})
}
