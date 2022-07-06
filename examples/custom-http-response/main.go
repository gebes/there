package main

// This example is commented out, to prevent the third party dependency from being added.

/*
import (
	"github.com/Gebes/there/v2"
	"github.com/vmihailenco/msgpack/v5"
	"log"
)

//Msgpack takes a StatusCode and data which gets marshaled to Msgpack
func Msgpack(code int, data any) there.Response {
   msgpackData, err := msgpack.Marshal(data) // marshal the data
   if err != nil {
      panic(err) // panic if the data was invalid. can be caught by Recoverer
   }
   return there.WithHeaders(MapString{ // set proper content-type
      there.ResponseHeaderContentType: "application/x-msgpack",
   }, there.Bytes(code, msgpackData))
}

func main() {
	router := there.NewRouter()

	router.
		Get("/", Get)

	err := router.Listen(8080)
	if err != nil {
		log.Fatalln("Could not listen to 8080", err)
	}
}

func Get(request there.Request) there.Response {
   return Msgpack(there.StatusOK, map[string]string{ // now use the created response
      "Hello": "World",
      "How":   "are you?",
   })
}
*/
