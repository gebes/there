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
   msgpackData, err := msgpack.Marshal(data) // marshal the data
   if err != nil {
      panic(err) // panic if the data was invalid. can be caught by Recoverer
   }
   return WithHeaders(MapString{ // set proper content-type
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
   return Msgpack(StatusOK, map[string]string{ // now use the created response
      "Hello": "World",
      "How":   "are you?",
   })
}
*/
