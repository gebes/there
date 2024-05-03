package main

import (
	"errors"
	"github.com/Gebes/there/v2"
	"github.com/Gebes/there/v2/middlewares"
	"github.com/Gebes/there/v2/status"
	"log"
	"math/rand"
)

type User struct {
	Id    int    `json:"id,omitempty"`
	Name  string `json:"name,omitempty"`
	Admin bool   `json:"admin,omitempty"`
}

func main() {
	router := there.NewRouter()

	router.Use(
		middlewares.Recoverer,
		middlewares.Sanitizer,
		middlewares.Gzip,
		middlewares.Logger(),
		middlewares.Cors(middlewares.CorsAllowAllConfiguration()),
	)

	router.Get("/", Get)
	router.Get("/user", Get)
	router.Get("/user/{name}", Get)

	err := router.Listen(8080)
	if err != nil {
		log.Fatalln("Could not listen to 8080", err)
	}
}

func Get(request there.Request) there.Response {
	if rand.Int()%3 == 0 {
		return there.Error(status.InternalServerError, errors.New("something went wrong"))
	}
	return there.Auto(status.OK, User{
		Id:    5,
		Name:  "Chris",
		Admin: true,
	})
}
