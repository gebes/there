package main

import (
	"errors"
	. "github.com/Gebes/there/v2"
)

func main() {

	router := NewRouter().SetProductionMode()

	router.Use(CorsMiddleware(AllowAllConfiguration()))

	router.Get("/", GetWelcome).
		With(RandomMiddleware)

	router.Group("home").
		Get("/", GetPage).
		Get("/params", GetParams).IgnoreCase().
		Get("/user", GetUser).IgnoreCase().
		Get("/user/:name", GetUserByName).IgnoreCase()

	err := router.Listen(8080)

	if err != nil {
		panic(err)
	}
}

type User struct {
	Name string
}

var count = 0

func RandomMiddleware(request HttpRequest, next HttpResponse) HttpResponse {
	count++
	if count%2 == 0 {
		// If you do not return Next(), then the Invocation-Chain will be broken, and the Response will be returned
		return Error(StatusInternalServerError, errors.New("lost database connection"))
	}
	// Next() means, that either the next Middleware or Handler (if it is the last Middleware) should be executed
	return next
}


func GetWelcome(HttpRequest) HttpResponse {
	return String(StatusOK, "Hello")
}

func GetUser(HttpRequest) HttpResponse {
	return Json(StatusOK, User{
		Name: "Hannes",
	})
}

func GetUserByName(request HttpRequest) HttpResponse {
	name := request.RouteParams.GetDefault("name", "Hannes")
	return String(StatusOK, "Hallo "+name)
}

func GetParams(request HttpRequest) HttpResponse {
	return Json(StatusOK, map[string]interface{}{
		"params":      request.Params,
		"routeParams": request.RouteParams,
	})
}

func GetPage(request HttpRequest) HttpResponse {
	user := request.Params.GetDefault("user", "Gebes")
	return Html(StatusOK, "./examples/index.html", map[string]string{
		"user": user,
	})
}
