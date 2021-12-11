package main

// This example is commented out, to prevent the third party dependency from being added.

/*
import (
	. "github.com/Gebes/there/v2"
	"github.com/vmihailenco/msgpack/v5"
	"log"
)

//Msgpack takes a StatusCode and data which gets marshaled to Msgpack
func Msgpack(code int, data interface{}) HttpResponse {
	msgpackData, err := msgpack.Marshal(data)
	if err != nil {
		panic(err)
	}
	return WithHeaders(MapString{
		ResponseHeaderContentType: "application/x-msgpack",
	}, Bytes(code, msgpackData))
}

func main() {
	router := NewRouter()

	router.
		Get("/", Get)

	err := router.Listen(8080)
	if err != nil {
		log.Fatalln("Could not listen to 8080", err)
	}
}

func Get(request HttpRequest) HttpResponse {
	return Msgpack(StatusOK, map[string]string{
		"Hello": "World",
		"How":   "are you?",
	})
}
*/