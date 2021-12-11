<p align="center">
<a href="https://there.gebes.io">
<img src="https://github.com/Gebes/there/blob/v2.1/.github/logo/logo-compressed.png?raw=true" alt="Gopher There" height="256">
</a>
</p>
<p align="center">
<font size="4px">
⚡️ Robust Web Framework to build Go Services with ease ⚡️
</font>
</p>

## 🔥 Get Started

### 🔨 Create a Project
1. Ensure you have [Go](https://go.dev/dl/) installed.
2. Create a project with `go mod init github.com/user/repository`
3. Install **There** with the `go get` command

```sh
go get -u github.com/Gebes/there/v2
```

4. Create a `main.go` file

```go
package main

import . "github.com/Gebes/there/v2"

func main() {
	router := NewRouter().
		Get("/", Get)

	err := router.Listen(8080)

	if err != nil {
		panic(err)
	}
}

func Get(request HttpRequest) HttpResponse {
	return Json(StatusOK, Map{
		"message": "Hello World!",
    })
}
```




## 🤔 Why There?


The general problem with many routers is the way you handle responses. Most frameworks make it too complex or do not offer the proper abstraction to get the same result in a short amount of time.  
The goal of **There** is to give developers the right tool to create robust apis in no time.

We solve this problem by providing simple interfaces to control the flow of your API.  
Got an error while fetching the user? Just `return Error(status, err)`. Want to return some data? Just `return Json(status, data)`. Is the data too large? Compress it `return Gzip(Json(status, data))`.  
This type of control flow is way easier to read, and it doesn't take away any freedom!


### Imports
If you create an API with **There** you do not need to import `net/http` even once! Simply import
```go
import . "github.com/Gebes/there/v2"
```
and **There** provides you with all the handlers, constants and interfaces you need to create a router, middleware or anything else!  
**There** provides enough constants for you! In total there are 140 of them.
* Method (`MethodGet`, `MethodPost`)
* Status (`StatusOK`, `StatusInternalServerError`)
* RequestHeader/ResponseHeader (`RequestHeaderContentType`, `RequestHeaderAcceptEncoding`, `ResponseHeaderLocation`)
* ContentType (`ContentTypeApplicationJson`, `ContentTypeApplicationXml`)

## 🧠 Philosophy
> Focus on your project, not on the framework.

The goal of **There** is to save time and provide all the functionality a backend developer needs daily. Therefore, **There** should always keep it simple and only add complexity if there is no reasonable workaround.  
New Go Developers often struggle with chores they are not used to, and **There** should help them make life easier.

So **There** should be something for everyone.

## ⭐️ Features
* [Straightforward **routing**](#straightforward-routing)
* [Minimalistic **control flow**](#minimalistic-control-flow)
* [**Expandable** - add your own control flow](#expandable---add-your-own-control-flow)
* [Complete **middleware** support](#complete-middleware-support)
* [**Lightweight** - no dependencies](#lightweight---no-dependencies)
* [Robust - 99,6% coverage](#robust---996-coverage)

### Straightforward routing

Routing with **There** is easy! Simply create a new router and add `GET`, `POST` or different handlers. Use groups to define multiple handlers simultaneously and protect your handlers with middlewares.
Define route variables with `:` and you have all the things you need for routing.

```go
	router := NewRouter()

	router.Group("/user").
		Get("/", Handler). // /user
		Post("/", Handler).
		Patch("/", Handler)

	router.Group("/user/post").
		Get("/:id", Handler).
		Post("/", Handler)

	router.
		Get("/details", Handler).IgnoreCase()
```

[🧑‍💻 View more code examples](https://github.com/Gebes/there/tree/main/examples/straightforward-routing)

### Minimalistic control flow


Controlling your route's flow with **There** is a delight! It is easy to understand and fast to write.  
A HttpResponse is basically a `http.handler`. **There** provides several handlers out of the box!

```go
func CreatePost(request HttpRequest) HttpResponse {
	var body Post
	err := request.Body.BindJson(&body)
	if err != nil {
		return Error(StatusBadRequest, "Could not parse body: "+err.Error())
	}

	post := postById(body.Id)
	if post != nil {
		return Error(StatusConflict, "Post with this ID already exists")
	}

	posts = append(posts, body)
	return Json(StatusCreated, body)
}
```

This handler uses `Error` and `Json`. By returning a HttpResponse, the handler chain will break and **There** will render the given response.

[🧑‍💻 View full code example](https://github.com/Gebes/there/tree/main/examples/minimalistic-control-flow)

### Expandable - add your own control flow

Simply create your own HttpResponse to save time. However, if you need some inspiration, look into the [response.go](https://github.com/Gebes/there/blob/main/response.go) file.

For example, let us create a Msgpack response. By default, there does not provide a Msgpack response, because this would require a third-party dependency. But it is not much work to create your own Msgpack HttpResponse:
```go
import (
    . "github.com/Gebes/there/v2"
    "github.com/vmihailenco/msgpack/v5"
)

//Msgpack takes a StatusCode and data which gets marshaled to Msgpack
func Msgpack(code int, data interface{}) HttpResponse {
   msgpackData, err := msgpack.Marshal(data)
   if err != nil {
      panic(err)
   }
   return WithHeaders(MapString{
      ResponseHeaderContentType: "application/x-msgpack",
   }, Bytes(code, msgpackData))
}

func Get(request HttpRequest) HttpResponse {
   return Msgpack(StatusOK, map[string]string{
      "Hello": "World",
      "How":   "are you?",
   })
}
```

**There** provides enough lower-level HttpResponses to build another one on top of it. At the bottom, we have a "Bytes" response, which writes the given bytes and the status code.  
Wrapped around the "Bytes" response, you can find a "WithHeaders" response, adding the ContentType header.

As you see, it is only a few lines of code to have a custom HttpResponse.

[🧑‍💻 View full code example](https://github.com/Gebes/there/tree/main/examples/custom-http-response)

### Complete middleware support

To keep things simple, you can use already existing middlewares with little to no change, and you can use the simple control flow from **There** in your middlewares.

Here is an example:
```go

func main() {
   router := NewRouter()

   router.Use(middlewares.Recoverer)
   router.Use(middlewares.Cors(middlewares.AllowAllConfiguration()))
   router.Use(GlobalMiddleware)

   router.Get("/", Get).With(RouteSpecificMiddleware)

   err := router.Listen(8080)
   if err != nil {
      log.Fatalln("Could not listen to 8080", err)
   }
}

func GlobalMiddleware(request HttpRequest, next HttpResponse) HttpResponse {

   if request.Headers.GetDefault(RequestHeaderContentType, "") != ContentTypeApplicationJson {
      return Error(StatusUnsupportedMediaType, "Header " + RequestHeaderContentType + " is not " + ContentTypeApplicationJson)
   }

   return next
}

func RouteSpecificMiddleware(request HttpRequest, next HttpResponse) HttpResponse {
   return WithHeaders(MapString{
      ResponseHeaderContentLanguage: "en",
   }, next)
}
```

With the `.Use` method, you can add a global middleware. No matter on which group you call it, it will be **global**.  
On the other side, if you use the `.With` method you can only add a middleware to **one handler**! Not to a whole group.

The `GlobalMiddleware` in this code checks if the request has `application/json` as content-type. If not, the request will fail with an error.
Compared to the `GlobalMiddleware`, the `RouteSpecificMiddleware` does not change the control flow but adds data to the response.

Be careful in this example. Global middlewares will always be called first, so if the global middleware returns an error, the content-language header won't be set by the `RouteSpecificMiddleware` middleware.

#### Using already existing middlewares
If you have other middlewares, which you created using other routers, then there is a high chance that you can use it in **There** without changing much.

As an example, let us have a look at the Recoverer middleware.

```go
func Recoverer(request HttpRequest, next HttpResponse) HttpResponse {
   fn := func(w http.ResponseWriter, r *http.Request) {
      defer func() {
         if rvr := recover(); rvr != nil && rvr != http.ErrAbortHandler {
            Error(StatusInternalServerError, rvr).ServeHTTP(w, r)
         }
      }()
      next.ServeHTTP(w, r)
   }
   return HttpResponseFunc(fn)
}
```
It is a trivial Recoverer. The only things you need to change are the types and the parameters. **There** provides you all the types required, so that you don't need to import "net/http".

[🧑‍💻 View full code example](https://github.com/Gebes/there/tree/main/examples/middlewares)

### Lightweight - no dependencies

**There** was built to be lightweight and robust, so if you use **There**, you do not need to worry about having many more external dependencies because **There** has none!

### Robust - 99,6% coverage

Almost everything in **There** has a corresponding test. We tested the framework well in production and wrote enough test cases, so in the end almost everything that makes sense to test was tested.
