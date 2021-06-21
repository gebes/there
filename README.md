# there
"there" also called "GoThere" aims to be a simple Go Library to reduce redundant code for REST APIs. 

Despite the existence of the other libraries, "there"s existence is justified by the minor amount of code you have to write to make API calls go there where you want to.

## Example
Basic GET and POST routing without database implementation
```go
package main

import (
	. "there"
)

type User struct {
	Id   int    `json:"id,omitempty" validate:"required"`
	Name string `json:"name,omitempty" validate:"required"`
}

func main() {
	router := Router{
		Port:              8080,
		LogRouteCalls:     true,
		LogResponseBodies: true,
		AlwaysLogErrors:   true,
	}

	router.HandleGet("/user", func(request Request) Response {
		return ResponseData(StatusOK, User{
			1, "John",
		})
	})
	router.HandlePost("/user", func(request Request) Response {

		var user User
		err := request.ReadBody(&user)
		if err != nil {
			return ResponseData(StatusBadRequest, err)
		}

		// code

		return ResponseData(StatusOK, "Saved the user to the database")

	})

	err := router.Listen()
	if err != nil {
		panic(err)
	}
}

```