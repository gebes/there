package main

import (
	"log"

	"github.com/Gebes/there/v2"
)

func main() {
	router := there.NewRouter()

	router.Group("/user").
		Get("/", Handler). // /user
		Post("/", Handler).
		Patch("/", Handler)

	router.Group("/user/post").
		Get("/:id", Handler).
		Post("/", Handler)

	router.
		Get("/details", Handler).IgnoreCase() // don't differentiate between "details", "DETAILS" or a mix of these two

	err := router.Listen(8080)
	if err != nil {
		log.Fatalln("Could not listen to 8080", err)
	}
}

func Handler(request there.Request) there.Response {
	return there.Json(there.StatusOK, map[string]any{
		"method":       request.Method,
		"path":         request.Request.URL.Path,
		"route_params": request.RouteParams,
	})
}
