package main

import (
	"github.com/Gebes/there/v2"
	"log"
)

func main() {
	router := there.NewRouter()

	router.Get("/", func(request there.Request) there.Response {
		return there.File("./examples/file-serving/main.go", there.ContentTypeTextPlain)
	})

	err := router.Listen(8080)
	if err != nil {
		log.Fatalln("Could not listen to 8080", err)
	}

	select {}

}
