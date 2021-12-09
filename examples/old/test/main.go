package main

import . "github.com/Gebes/there/v2"

func main() {
	router := NewRouter().
		Get("/", Get)

	router.Use(Recoverer)

	err := router.Listen(8080)
	if err != nil {
		panic(err)
	}
}


func Get(request HttpRequest) HttpResponse {
	return Message(StatusOK, "Hello!")
}