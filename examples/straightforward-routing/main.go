package main

import (
	"github.com/Gebes/there/v2/status"
	"log"

	"github.com/Gebes/there/v2"
)

func main() {
	router := there.NewRouter()

	router.Group("/user").
		Get("/all", Handler).
		Get("/{id}", Handler).
		Post("/", Handler).
		Patch("/{id}", Handler).
		Delete("/{id}", Handler)

	router.Group("/user/post").
		Get("/{id}", Handler).
		Post("/", Handler).
		Patch("/{id}", Handler).
		Delete("/{id}", Handler)

	router.Get("/details", Handler)

	err := router.Listen(8080)
	if err != nil {
		log.Fatalf("Could not listen to 8080: %v", err)
	}
}

func Handler(request there.Request) there.Response {
	return there.Json(status.OK, map[string]any{
		"path":         request.Request.URL.Path,
		"method":       request.Method,
		"route_params": request.RouteParams,
	})
}
