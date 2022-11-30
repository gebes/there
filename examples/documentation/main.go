package main

import (
	"errors"
	"fmt"
	"github.com/Gebes/there/v2"
	"github.com/Gebes/there/v2/status"
)

func main() {
	router := there.NewRouter()

	router.Group("/example").
		Get("/json", ExampleJsonGet).
		Get("/jsonerror", ExampleJsonGet).
		Get("/xml", ExampleXmlGet).
		Get("/error", ExampleErrorGet).
		Get("/message", ExampleMessageGet).
		Get("/status", ExampleStatusGet).
		Get("/string", ExampleStringGet)

	err := router.Listen(8080)
	if err != nil {
		fmt.Printf("err: %v\n", err)
	}
}

func ExampleJsonGet(request there.Request) there.Response {
	user := map[string]string{
		"firstname": "John",
		"surname":   "Smith",
	}
	return there.Json(status.OK, user)
}

func ExampleJsonErrorGet(request there.Request) there.Response {
	user := map[string]string{
		"firstname": "John",
		"surname":   "Smith",
	}
	resp, err := there.JsonError(status.OK, user)
	if err != nil {
		return there.Error(status.InternalServerError, fmt.Errorf("something went wrong: %v", err))
	}
	return resp
}

type User struct {
	Firstname string `xml:"firstname"`
	Surname   string `xml:"surname"`
}

func ExampleXmlErrorGet(request there.Request) there.Response {
	user := User{"John", "Smith"}
	resp, err := there.XmlError(status.OK, user)
	if err != nil {
		return there.Error(status.InternalServerError, fmt.Errorf("something went wrong: %v", err))
	}
	return resp
}

func ExampleXmlGet(request there.Request) there.Response {
	user := User{"John", "Smith"}
	return there.Xml(status.OK, user)
}

func ExampleErrorGet(request there.Request) there.Response {
	if 1 != 2 {
		return there.Error(status.InternalServerError, errors.New("something went wrong"))
	}
	return there.Status(status.OK)
}

func ExampleMessageGet(request there.Request) there.Response {
	return there.Message(status.OK, "Hello there")
}

func ExampleStatusGet(request there.Request) there.Response {
	return there.Status(status.OK)
}

func ExampleStringGet(request there.Request) there.Response {
	return there.String(status.OK, "Hello there")
}

func ExampleBytesGet(request there.Request) there.Response {
	data := []byte("Hello there")
	return there.Bytes(status.OK, data)
}
