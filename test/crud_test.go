package test

import (
	"github.com/Gebes/there/there"
	"net/http"
	"testing"
)

func TestCrud(t *testing.T) {
	router := there.Router{Port: 8080}
	router.Handle("/monkeys", func(request there.Request) there.Response {
		return there.ResponseData(http.StatusOK, "I like monkeys too")
	}, http.MethodGet)
	router.Handle("/monkey", func(request there.Request) there.Response {
		return there.ResponseData(http.StatusOK, map[string]string{
			"name":           "Alfred",
			"age":            "Unkown",
			"favourite_meal": "Affensuppe",
		})
	}, http.MethodGet)
	router.Handle("/user", func(request there.Request) there.Response {
		return there.Response{Status: http.StatusOK, Data: map[string]string{
			"name": "Henrick",
		}}
	})

	err := unblock(func() error {
		return router.Listen()
	})
	if err != nil {
		panic(err)
	}

	get("http://localhost:8080/monkeys?test=123")
	get("http://localhost:8080/monkeys")
	get("http://localhost:8080/monkey")
}


