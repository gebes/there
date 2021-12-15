package main

import (
	. "github.com/Gebes/there/v2"
	"log"
)

func main() {
	router := NewRouter()

	user := router.Group("/user").
		Get("/", Handler). // /user
		Post("/", Handler).
		Patch("/", Handler)

	user.Group("/post").
		Get("/:id", Handler).
		Post("/", Handler)

	router.
		Get("/details", Handler).IgnoreCase() // don't differentiate between "details", "DETAILS" or a mix of these two

	err := router.Listen(8080)
	if err != nil {
		log.Fatalln("Could not listen to 8080", err)
	}
}

func Handler(request HttpRequest) HttpResponse {
	return Json(StatusOK, Map{
		"method":       request.Method,
		"path":         request.Request.URL.Path,
		"route_params": request.RouteParams,
	})
}
