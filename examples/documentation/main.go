package main

import (
	"errors"
	"fmt"
	"github.com/Gebes/there/v2"
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
	return there.Json(there.StatusOK, user)
}

func ExampleJsonErrorGet(request there.Request) there.Response {
	user := map[string]string{
		"firstname": "John",
		"surname":   "Smith",
	}
	resp, err := there.JsonError(there.StatusOK, user)
	if err != nil {
		return there.Error(there.StatusInternalServerError, fmt.Errorf("something went wrong: %v", err))
	}
	return resp
}

type User struct {
	Firstname string `xml:"firstname"`
	Surname   string `xml:"surname"`
}

func ExampleXmlErrorGet(request there.Request) there.Response {
	user := User{"John", "Smith"}
	resp, err := there.XmlError(there.StatusOK, user)
	if err != nil {
		return there.Error(there.StatusInternalServerError, fmt.Errorf("something went wrong: %v", err))
	}
	return resp
}

func ExampleXmlGet(request there.Request) there.Response {
	user := User{"John", "Smith"}
	return there.Xml(there.StatusOK, user)
}

func ExampleErrorGet(request there.Request) there.Response {
	if 1 != 2 {
		return there.Error(there.StatusInternalServerError, errors.New("something went wrong"))
	}
	return there.Status(there.StatusOK)
}

func ExampleMessageGet(request there.Request) there.Response {
	return there.Message(there.StatusOK, "Hello there")
}

func ExampleStatusGet(request there.Request) there.Response {
	return there.Status(there.StatusOK)
}

func ExampleStringGet(request there.Request) there.Response {
	return there.String(there.StatusOK, "Hello there")
}


func ExampleBytesGet(request there.Request) there.Response {
	data := []byte("Hello there")
	return there.Bytes(there.StatusOK, data)
}