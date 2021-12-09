package main

import (
	. "github.com/Gebes/there/v2"
	"log"
)

type User struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

func NewUser(name string, description string) *User {
	return &User{Name: name, Description: description}
}

var (
	users = []*User{
		NewUser("Steve Jobs", "Apple Founder"),
		NewUser("Elon Musk", "Cool guy"),
		NewUser("Bill Gates", "Microsoft Founder"),
		NewUser("Tim Cook", "Current Apple Ceo"),
	}
)

func main() {
	router := NewRouter().
		Get("/users", GetUsers)

	err := router.Listen(8080)
	if err != nil {
		log.Fatalln("Could not start listening on port 8080", err)
	}

}

func GetUsers(request HttpRequest) HttpResponse {
	// there automatically sets the Content-Type header
	return Json(StatusOK, users) // return all the users as JSON
	//	return Xml(StatusOK, users)
}

