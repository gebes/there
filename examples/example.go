package main

import (
	"errors"
	. "github.com/Gebes/there/there/http/middlewares"
	. "github.com/Gebes/there/there/http/request"
	. "github.com/Gebes/there/there/http/response"
	. "github.com/Gebes/there/there/http/router"
	. "github.com/Gebes/there/there/utils"
)

func main() {

	err := NewRouter().
		Get("/", GetWelcome).AddMiddleware(RandomMiddleware).
		Get("/user", GetUser).
		Get("/test", GetTest).
		AddGlobalMiddleware(CorsMiddleware(AllowAllConfiguration())).
		Listen(":8080")

	if err != nil {
		panic(err)
	}
}

type User struct {
	Name string
}

var count = 0

func RandomMiddleware(request HttpRequest) HttpResponse {
	count++

	if count%2 == 0 {
		return Error(StatusInternalServerError, errors.New("lost database connection")).
					Header().Set(ResponseHeaderContentLanguage, "English")
	}

	return Next()
}

func GetWelcome(request HttpRequest) HttpResponse {
	return String(StatusOK, "Hello")
}

func GetUser(request HttpRequest) HttpResponse {
	return Json(StatusOK, User{
		Name: "Hannes",
	})
}

func GetTest(request HttpRequest) HttpResponse {
	body, err := request.ReadBody().String()

	if err != nil {
		return Error(StatusBadRequest, err)
	}
	return Json(StatusOK, map[string]interface{}{
		"body":   body,
		"params": request.ReadParams().Map(),
	})
}
