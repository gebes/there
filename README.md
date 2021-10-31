# there
`there`, also called `GoThere`, aims to be a simple Go Library to reduce redundant code for REST APIs. The name `GoThere` because it tells the incoming requests to "go there".
With `there` it is a delight to create a REST API. You can focus on writing your API without writing boilerplate and duplicate code. Just start coding and save time!

## Table of contents

* [Installation](#install)
* [Basic Example](#basic-example)
    * [Validating Body](#validating-body)
* [Cors Example](#cors-example)
* [Redirect Example](#redirect-example)
* [Middlewares](#middlewares)
    * [Global Middlewares](#global-middlewares)
    * [Route Specific Middlewares](#route-specific-middlewares)


## Install

```sh
go get -u github.com/Gebes/there
```

## Example

### Basic Example
`GET /user` to read a user from the database  
`POST /user` to store a user in the database
```go
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
```
#### Validating Body
`ReadBody` will return an error, if either `Id` or `Name` are not present, because of the `validate:"required"` tag. 
You can find more information [here](https://gopkg.in/go-playground/validator.v9)
### Cors Example
Automatically adds all the Access-Control headers. If you are up to it, you can use your config only to allow specific origins, methods, or headers. By default, all origins and headers are allowed 
```go
package examples

import (
	. "github.com/Gebes/there/there"
	"github.com/Gebes/there/there/middlewares/cors"
)

func CorsRouter() error {
	router := Router{
		Port:              8080,
		LogRouteCalls:     true,
		LogResponseBodies: true,
		AlwaysLogErrors:   true,
	}

	router.AddGlobalMiddleware(cors.MiddlewareCors(cors.DefaultConfiguration()))

	router.HandleGet("/cors", func(request Request) Response {
		return ResponseData(StatusOK, "Hello there")
	})


	return router.Listen()
}


```

### Redirect Example
Redirects the user to google.com if you call `localhost:8080/google`
```go
package examples

import (
	. "github.com/Gebes/there/there"
	"github.com/Gebes/there/there/middlewares/cors"
)

func RedirectRouter() error {
	router := Router{
		Port:              8080,
		LogRouteCalls:     true,
		LogResponseBodies: true,
		AlwaysLogErrors:   true,
	}

	router.AddGlobalMiddleware(cors.MiddlewareCors(cors.DefaultConfiguration()))

	router.HandleGet("/google", func(request Request) Response {
		return ResponseRedirect("https://google.com")
	})

	return router.Listen()
}


```

### Middlewares
GoThere supports route-specific and global middlewares
#### Global Middlewares
```go
package examples

import (
	. "github.com/Gebes/there/there"
	"log"
)


func GlobalMiddlewareRouter()error {
	router := Router{
		Port:              8080,
		LogRouteCalls:     true,
		LogResponseBodies: true,
		AlwaysLogErrors:   true,
	}

	router.AddGlobalMiddleware(func(request Request) *Response {
		log.Println("First middleware")
		return nil
	})

	router.AddGlobalMiddleware(func(request Request) *Response {
		log.Println("Second middleware")
		return ResponseDataP(StatusForbidden, "I won't let you in!")
	})

	router.HandleGet("/message", func(request Request) Response {
		return ResponseData(StatusOK, "Hello from the Handler")
	})


	return router.Listen()
}
```
Calling the `/message` route returns in
```
2021/07/17 10:45:16 First middleware
2021/07/17 10:45:16 Second middleware
2021/07/17 10:45:16 GET /message resulted in 403 with body {"status":403,"message":"I won't let you in!"}
2021/07/17 10:45:16 Body {"status":403,"message":"I won't let you in!"}
```
#### Route Specific Middlewares
```go
package examples

import (
	. "github.com/Gebes/there/there"
	"log"
)

func MiddlewareRouter() error {
	router := Router{
		Port:              8080,
		LogRouteCalls:     true,
		LogResponseBodies: true,
		AlwaysLogErrors:   true,
	}

	router.HandleGet("/message", messageHandler).AddMiddleware(messageMiddleware, authorizationMiddleware)
	router.HandleGet("/message/without/middleware", messageHandler)

	return router.Listen()
}

func messageMiddleware(request Request) *Response {
	log.Println("First Middleware")
	return nil
}

func authorizationMiddleware(request Request) *Response{
	params, err := request.ReadParams("authorization")
	if err != nil {
		return ResponseDataP(StatusBadRequest, err)
	}
	if params[0] != "password" {
		return ResponseDataP(StatusUnauthorized, err)
	}
	return nil
}

func messageHandler(request Request) Response {
	return ResponseData(StatusOK, "Hello from the Handler")
}
```
The result if you call `/message` without any parameters
```
2021/07/17 10:47:21 First Middleware
2021/07/17 10:47:21 GET /message resulted in 400 with body {"status":400,"message":"required parameter not existing: authorization"}
2021/07/17 10:47:21 Body {"status":400,"message":"required parameter not existing: authorization"}
```
The result if you call `/message/without/middleware`
```
2021/07/17 10:47:21 GET /message/without/middleware resulted in 200 with body {"status":200,"message":"Hello from the Handler"}
2021/07/17 10:47:21 Body {"status":200,"message":"Hello from the Handler"}
```