package main

import (
	. "github.com/Gebes/there/v2"
	"github.com/Gebes/there/v2/middlewares"
	"log"
)

func main() {
	router := NewRouter()

	router.Use(middlewares.Recoverer)
	router.Use(middlewares.Cors(middlewares.AllowAllConfiguration()))
	router.Use(GlobalMiddleware)

	router.Get("/", Get).With(RouteSpecificMiddleware)

	err := router.Listen(8080)
	if err != nil {
		log.Fatalln("Could not listen to 8080", err)
	}
}

func GlobalMiddleware(request HttpRequest, next HttpResponse) HttpResponse {

	if request.Headers.GetDefault(RequestHeaderContentType, "") != ContentTypeApplicationJson {
		return Error(StatusUnsupportedMediaType, "Header "+RequestHeaderContentType+" is not "+ContentTypeApplicationJson)
	}

	return next
}

func RouteSpecificMiddleware(request HttpRequest, next HttpResponse) HttpResponse {
	return WithHeaders(MapString{
		ResponseHeaderContentLanguage: "en",
	}, next)
}

func Get(request HttpRequest) HttpResponse {
	return Json(StatusOK, map[string]string{
		"Hello": "World",
		"How":   "are you?",
	})
}
