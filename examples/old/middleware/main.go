package main

import (
	"context"
	"errors"
	. "github.com/Gebes/there/v2"
)

func main() {

	// Register Global Middleware
	router := NewRouter().Use(RandomMiddleware).Use(CorsMiddleware(AllowAllConfiguration()))

	router.
		// Registers Middleware only for the "/" route
		Get("/", GetAuthHeader).With(DataMiddleware)

	err := router.Listen(8080)
	if err != nil {
		panic(err)
	}
}

var count = 0

//RandomMiddleware returns an example error for every second request
func RandomMiddleware(request HttpRequest, next HttpResponse) HttpResponse {
	count++
	if count%2 == 0 {
		// If you do not return Next(), then the Invocation-Chain will be broken, and the Response will be returned
		return Error(StatusInternalServerError, errors.New("lost database connection"))
	}
	// Next() means, that either the next Middleware or Handler (if it is the last Middleware) should be executed
	return next
}

//DataMiddleware checks if the user provided an Authorization header. If so, then it will be passed on to the handler via Context
func DataMiddleware(request HttpRequest, next HttpResponse) HttpResponse {
	auth := request.Headers.GetDefault(RequestHeaderAuthorization, "")
	if len(auth) == 0 {
		return Error(StatusBadRequest, errors.New("no authorization header provider"))
	}
	// We wrap Next() with a Context, by using the WithContext Wrapper.
	// In the GetAuthHeader Handler, we can then use the current Context to read "auth"
	// WithContext can also be returned in a regular Handler, but it would make no sense. To where do you want to pass the context???
	// The WithContext() and Next() HttpResponse should only be used for Middlewares
	request.WithContext(context.WithValue(request.Context(), "auth", auth))
	return next
}

func GetAuthHeader(request HttpRequest) HttpResponse {
	// Read from the context
	data, ok := request.Context().Value("auth").(string)

	if !ok { // Could not read from the context... should not happen, except we forgot to add the DataMiddleware to the Route
		return Error(StatusUnprocessableEntity, errors.New("could not get auth from context"))
	}

	return String(StatusOK, "Auth: "+data)
}
