package main

import (
	"errors"
	"log"

	. "github.com/Gebes/there/v2"
	"github.com/Gebes/there/v2/middlewares"
)

func main() {
	router := NewRouter()

	// Register global middlewares
	router.Use(middlewares.Recoverer)
	router.Use(middlewares.Cors(middlewares.AllowAllConfiguration()))
	router.Use(GlobalMiddleware)

	router.
		Get("/", Get).With(RouteSpecificMiddleware) // Register route with middleware

	err := router.Listen(8080)
	if err != nil {
		log.Fatalln("Could not listen to 8080", err)
	}
}

func GlobalMiddleware(request HttpRequest, next HttpResponse) HttpResponse {
	// Check the request content-type
	if request.Headers.GetDefault(RequestHeaderContentType, "") != ContentTypeApplicationJson {
		return Error(StatusUnsupportedMediaType, errors.New("Header "+RequestHeaderContentType+" is not "+ContentTypeApplicationJson))
	}

	return next // Everything is fine until here, continue
}

func RouteSpecificMiddleware(request HttpRequest, next HttpResponse) HttpResponse {
	return WithHeaders(MapString{
		ResponseHeaderContentLanguage: "en",
	}, next) // Set the content-language header by wrapping next with WithHeaders
}

func Get(request HttpRequest) HttpResponse {
	return Json(StatusOK, map[string]string{
		"Hello": "World",
		"How":   "are you?",
	})
}
