package examples

import (
	"github.com/Gebes/there/there"
)

type User struct {
	Id   int    `json:"id,omitempty" validate:"required"`
	Name string `json:"name,omitempty" validate:"required"`
}

func CrudRouter() error{
	router := there.Router{
		Port:              8080,
		LogRouteCalls:     true,
		LogResponseBodies: true,
		AlwaysLogErrors:   true,
	}

	router.HandleGet("/user", func(request there.Request) there.Response {
		return there.ResponseData(there.StatusOK, User{
			1, "John",
		})
	})
	router.HandlePost("/user", func(request there.Request) there.Response {

		var user User
		err := request.ReadBody(&user)
		if err != nil {
			return there.ResponseData(there.StatusBadRequest, err)
		}

		// code

		return there.ResponseData(there.StatusOK, "Saved the user to the database")

	})

	return router.Listen()
}
