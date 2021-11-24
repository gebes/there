# there

`there`, also called `GoThere`, aims to be a simple Go Library to reduce redundant code for REST APIs. The
name `GoThere` because it tells the incoming requests to "go there". With `there` it is a delight to create a REST API.
You can focus on writing your API without writing boilerplate and duplicate code. Just start coding and save time!

## Table of contents

If it states `Complete example`, you can just copy the whole example and run it directly without changing something.

* [Installation](#install)
* [Examples](#example)
  * [Create a router](#create-a-router)
  * [Listen to a port](#listen-to-a-port)
  * [Complete example: returning Json, Xml, Yaml or Msgpack](#complete-example-returning-json-xml-yaml-or-msgpack)
  * [All HttpResponses](#all-httpresponses)
  * [Complete example: Middlewares](#middlewares)

## Install 

```sh
go get -u github.com/Gebes/there/v2
```

## Examples

Let's go through some basic examples, which make you understand the library in less than 10 minutes. Feel free to play around!  

I recommend you import `there` in every file you use it like the following:

```go
import (
	. "github.com/Gebes/there/v2"
)
```

If you are not familiar with this syntax, this allows you to use `there` without the `there.` prefix.  

### Create a router

```go
func main(){
	router := NewRouter().
		Get("/user/:id", GetUser).
		Post("/user", PostUser).
		Patch("/user/:id", PatchUser).
		Delete("/user/:id", DeleteUser)
}
```

Just create a new router instance and register some routes. `there` provides simple builder patterns, so you don't need to
write `router.` for every route or middleware.

### Listen to a port
```go
func main(){
	router := NewRouter()
	// ...
	err := router.Listen(8080)
	if err != nil {
		log.Fatalln("Could not start listening on port 8080", err)
	}
}
```

The `Listen` method binds the router blocking to this port. A possible error could be that the port is already in use by a different program.

### Handle a request

We already know how to define routes. But how do we handle them?  

We need to pass a handler function, as the following

```go
router := NewRouter().
		Get("/route", func(request HttpRequest) HttpResponse {
		
		})
```

or

```go
func main() {
	router := NewRouter().
		Get("/route", RouteHandler)
}

func RouteHandler(request HttpRequest) HttpResponse {

}
```

### Complete example: returning Json, Xml, Yaml or Msgpack

This example provides a /users route, which returns a list of users in the JSON Format.  

If you want to, you can also return the data in Xml, Yaml, or even Msgpack.

```go
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
//	return Yaml(StatusOK, users)
//	return Msgpack(StatusOK, users) // Msgpack is supported out of the box
}


```

If you run this example and open [localhost:8080/users](http://localhost:8080/users) in your browser, then you get the following result:

```json
[{"name":"Steve Jobs","description":"Apple Founder"},{"name":"Elon Musk","description":"Cool guy"},{"name":"Bill Gates","description":"Microsoft Founder"},{"name":"Tim Cook","description":"Current Apple Ceo"}]
```

Here is the result formatted:

```json
[
  {
    "name": "Steve Jobs",
    "description": "Apple Founder"
  },
  {
    "name": "Elon Musk",
    "description": "Cool guy"
  },
  {
    "name": "Bill Gates",
    "description": "Microsoft Founder"
  },
  {
    "name": "Tim Cook",
    "description": "Current Apple Ceo"
  }
]
```

### All HttpResponses

Here is a list of all the valid HttpResponse Returns you can make:

```go
func RouteHandler(request HttpRequest) HttpResponse {
	return Empty(StatusOK)                      // Returns nothing

	return Bytes(StatusOK, []byte("A message")) // Same result, but the input is a byte array
	return String(StatusOK, "A message")        // Just return a plain string

	return Redirect("https://www.google.com")   // Redirect the request to another page

	return Error(StatusInternalServerError, errors.New("parse an error")) // Will be formatted accordingly to router.RouterConfiguration.ErrorMarshal
	// Default is JSON: {"error": "parse an error"}
	
	return Json(StatusOK, users)
	return Xml(StatusOK, users)
	return Yaml(StatusOK, users)
	return Msgpack(StatusOK, users)                                  // Msgpack is supported out of the box
	return Html(StatusOK, "./files/index.html", map[string]string{}) // Reads the HTML File and uses it as a template. Variables from the map will be replaced into the response
}
```

If you are up to it, you can also create your own Response by creating a struct that implements the `HttpResponse` interface.

### Middlewares

Of course, `there` has middleware support. You can either have global middlewares by using `router.Use(middleware)` or route-specific middlewares by using `router.Get("/route", handler).With(middleware)`.  
You cannot do `router.With(middleware)` because the `.With(middleware)` method requires a route to be defined before. It will add the middleware always to the last added route.

```go
package main

import (
	"context"
	"errors"
	. "github.com/Gebes/there/v2"
)

func main() {

	// Register global middleware 
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
func RandomMiddleware(HttpRequest) HttpResponse {
	count++
	if count%2 == 0 {
		// If you do not return Next(), then the Invocation-Chain will be broken, and the Response will be returned
		return Error(StatusInternalServerError, errors.New("lost database connection"))
	}
	// Next() means, that either the next middleware or Handler (if it is the last middleware) should be executed
	return Next()
}

//DataMiddleware checks if the user provided an Authorization header. If so, then it will be passed on to the handler via Context
func DataMiddleware(request HttpRequest) HttpResponse {
	auth := request.Headers.GetDefault(RequestHeaderAuthorization, "")
	if len(auth) == 0 {
		return Error(StatusBadRequest, errors.New("no authorization header provider"))
	}
	// We wrap Next() with a Context, by using the WithContext Wrapper.
	// In the GetAuthHeader Handler, we can then use the current Context to read "auth"
	// WithContext can also be returned in a regular Handler, but it would make no sense. To where do you want to pass the context???
	// The WithContext() and Next() HttpResponse should only be used for middlewares
	return WithContext(context.WithValue(request.Context(), "auth", auth), Next())
}

func GetAuthHeader(request HttpRequest) HttpResponse {
	// Read from the context
	data, ok := request.Context().Value("auth").(string)

	if !ok { // Could not read from the context... should not happen, except we forgot to add the DataMiddleware to the Route
		return Error(StatusUnprocessableEntity, errors.New("could not get auth from context"))
	}

	return String(StatusOK, "Auth: "+data)
}

```
The example seems a bit too big, but it shows everything that you can do with middlewares. We added two global middlewares. One which we defined on our own and one cors middleware, which allows everything.  
Our RandomMiddleware is now used globally, which means it will be called before any Route Handler. As a result, every second call to our API will fail with the defined error.   
Our DataMiddleware is only used for the GetAuthHeader Route. Therefore, it gets the "Authorization" Header. If the Header is empty, then it will return an error. If not, it will pass the Authorization Header via Context to the next middleware or final Route Handler.
In this case, we do not have any extra middlewares, so it will call the GetAuthHeader handler, read from the Context, and return it as a String.
