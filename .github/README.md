<p align="center">
<a href="https://there.gebes.io">
<img src="https://github.com/Gebes/there/blob/main/.github/logo/logo-compressed.png?raw=true" alt="Gopher There" height="256">
</a>
</p>
<p align="center">
<font size="4px">
A robust and minimalistic Web Framework that helps you to build fabulous Go services in no time!
</font>
</p>
<p align="center">
<a href="http://golang.org">
    <img src="https://img.shields.io/badge/Made%20with-Go-1f425f.svg" alt="Made with Go">
</a>
<img src="https://img.shields.io/github/go-mod/go-version/Gebes/there.svg" alt="Go Version">
<a href="https://there.gebes.io/">
    <img src="https://img.shields.io/badge/Documentation-there-blue.svg" alt="Documentation">
</a>
<a href="https://pkg.go.dev/github.com/Gebes/there/v2">
    <img src="https://img.shields.io/badge/godoc-reference-blue.svg" alt="GoDoc">
</a>
<a href="https://goreportcard.com/report/github.com/Gebes/there">
    <img src="https://goreportcard.com/badge/github.com/Gebes/there" alt="GoReportCard">
</a>
<a href="https://gocover.io/github.com/Gebes/there">
    <img src="https://gocover.io/_badge/github.com/Gebes/there" alt="CodeCoverage">
</a>
<a href="https://github.com/Gebes/there/blob/master/LICENSE">
    <img src="https://img.shields.io/github/license/Gebes/there.svg" alt="License">
</a>
<a href="https://GitHub.com/Gebes/there/releases/">
    <img src="https://img.shields.io/github/release/Gebes/there" alt="Latest release">
</a>


## 🔥 Get Started

### 🔨 Create a Project
1. Ensure you have [Go](https://go.dev/dl/) installed.
2. Create a project with `go mod init github.com/{user}/{repository}`
3. Install **There** with the `go get` command

```sh
go get -u github.com/Gebes/there/v2
```

4. Create a `main.go` file

```go
package main

import (
	"github.com/Gebes/there/v2"
	"github.com/Gebes/there/v2/status"
)

func main() {
	router := there.NewRouter() // Create a new router

	// Register GET route /
	router.Get("/", func(request there.Request) there.Response {
		return there.Json(status.OK, map[string]string{
			"message": "Hello World!",
		})
	})

	err := router.Listen(8080) // Start listening on 8080

	if err != nil {
		panic(err)
	}
}
```

[📚 Check the Documentation for more info](https://there.gebes.io)

## 🤔 Why There?


The general problem with many routers is the way you handle responses. Most frameworks make it too complex or do not offer the proper abstraction to get the required result in a short amount of time.  
The goal of **There** is to give developers the right tool to create robust apis in a shorter amount of time.

We solve this problem by providing simple interfaces to control the flow of your API.  
Got an error while fetching the user? Just `return Error(status, err)`. Want to return some data? Just `return Json(status, data)`. Is the data too large? Compress it `return Gzip(Json(status, data))`.  
This type of control flow is way easier to read, and it doesn't take away any freedom!

### ⚡️ Speed
Speed is critical, even though your Go router will never be a bottleneck. As a comparison, **There** is faster than Gin and Mux ([Benchmark](https://pastebin.com/iP7NhtZH), [Result](https://pastebin.com/RrKi8B3J)).

### 📤 Imports
If you create an API with **There** you do not need to import `net/http` even once! Simply import
```go
import "github.com/Gebes/there/v2"
```
and **There** provides you with all the handlers, constants and interfaces you need to create a router, middleware or anything else!  
**There** provides enough constants for you! In total there are 140 of them.
* Method (`there.MethodGet`, `there.MethodPost`)
* Status (`status.OK`, `statusInternalServerError`)
* Header/Request only Header/Response only Header (`header.ContentType`, `header.RequestAcceptEncoding`, `header.ResponseLocation`)
* ContentType (`there.ContentTypeApplicationJson`, `there.ContentTypeApplicationXml`)

## 🧠 Philosophy
> Focus on your project, not on the framework.

The goal of **There** is to save time and provide all the functionality a backend developer needs daily. Therefore, **There** should always keep it simple and only add complexity if there is no reasonable workaround.  
New Go Developers often struggle with chores they are not used to, and **There** should help them make life easier.

From beginner to expert, **There** should be something for everyone.

## ⭐️ Features
* [Straightforward **routing**](#straightforward-routing)
* [Minimalistic **control flow**](#minimalistic-control-flow)
* [**Expandable** - add your own control flow](#expandable---add-your-own-control-flow)
* [Complete **middleware** support](#complete-middleware-support)
* [**Lightweight** - no dependencies](#lightweight---no-dependencies)
* [Robust](#robust---996-coverage)

### Straightforward routing

Routing with **There** is easy! Simply create a new router and add `GET`, `POST` or different handlers. Use groups to define multiple handlers simultaneously and protect your handlers with middlewares.
Define route variables with `:` and you have all the things you need for routing.

```go
	router := there.NewRouter()

	router.Group("/user").
		Get("/", Handler).
		Post("/", Handler).
		Patch("/", Handler)

	router.Group("/user/post").
		Get("/:id", Handler).
		Post("/", Handler)

	router.Get("/details", Handler)
```

[🧑‍💻 View full code example and best practices](https://github.com/Gebes/there/tree/main/examples/straightforward-routing)

### Minimalistic control flow


Controlling your route's flow with **There** is a delight! It is easy to understand and fast to write.  
A HttpResponse is basically a `http.Handler`. **There** provides several handlers out of the box!

```go
func CreatePost(request there.Request) there.Response {
	var body Post
	err := request.Body.BindJson(&body) // Decode body
	if err != nil { // If body was not valid json, return bad request error
		return there.Error(status.BadRequest, "Could not parse body: "+err.Error())
	}

	post := postById(body.Id)
	if post != nil { // if the post already exists, return conflict error
		return there.Error(status.Conflict, "Post with this ID already exists")
	}

	posts = append(posts, body) // create post
	return there.Json(status.Created, body) // return created post as json
}
```

This handler uses `Error` and `Json`. By returning a HttpResponse, the handler chain will break and **There** will render the given response.

[📚 Documentation](https://there.gebes.io/responses/status)  
[🧑‍💻 View full code example](https://github.com/Gebes/there/tree/main/examples/minimalistic-control-flow)

### Expandable - add your own control flow

Simply create your own HttpResponse to save time. However, if you need some inspiration, look into the [response.go](https://github.com/Gebes/there/blob/main/response.go) file.

For example, let us create a Msgpack response. By default, there does not provide a Msgpack response, because this would require a third-party dependency. But it is not much work to create your own Msgpack HttpResponse:
```go
import (
    "github.com/Gebes/there/v2"
    "github.com/vmihailenco/msgpack/v5"
)

//Msgpack takes a StatusCode and data which gets marshaled to Msgpack
func Msgpack(code int, data interface{}) there.Response {
   msgpackData, err := msgpack.Marshal(data) // marshal the data
   if err != nil {
      panic(err) // panic if the data was invalid. can be caught by Recoverer
   }
   return there.WithHeaders(map[string]string{ // set proper content-type
      there.ResponseHeaderContentType: "application/x-msgpack",
   }, there.Bytes(code, msgpackData))
}

func Get(request there.Request) there.Response {
   return Msgpack(status.OK, map[string]string{ // now use the created response
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
   router := there.NewRouter()

   // Register global middlewares
   router.Use(middlewares.Recoverer)
   router.Use(middlewares.Cors(middlewares.AllowAllConfiguration()))
   router.Use(GlobalMiddleware)

   router.
   	Get("/", Get).With(RouteSpecificMiddleware) // Register route with middleware

   err := router.Listen(8080)
   if err != nil {
      log.Fatalln("Could not listen to 8080", err)
   }
}

func GlobalMiddleware(request there.Request, next there.Response) there.Response {
   // Check the request content-type
   if request.Headers.GetDefault(there.RequestHeaderContentType, "") != there.ContentTypeApplicationJson {
      return there.Error(status.UnsupportedMediaType, "Header " + there.RequestHeaderContentType + " is not " + there.ContentTypeApplicationJson)
   }
   
   return next // Everything is fine until here, continue
}

func RouteSpecificMiddleware(request there.Request, next there.Response) there.Response {
   return there.WithHeaders(MapString{
      there.ResponseHeaderContentLanguage: "en",
   }, next) // Set the content-language header by wrapping next with WithHeaders
}
```

With the `.Use` method, you can add a global middleware. No matter on which group you call it, it will be **global**.  
On the other side, if you use the `.With` method you can only add a middleware to **one handler**! **Not to a whole group.**

The `GlobalMiddleware` in this code checks if the request has `application/json` as content-type. If not, the request will fail with an error.
Compared to the `GlobalMiddleware`, the `RouteSpecificMiddleware` does not change the control flow but adds data to the response.

Be careful in this example. Global middlewares will always be called first, so if the global middleware returns an error, the content-language header won't be set by the `RouteSpecificMiddleware` middleware.

#### Using already existing middlewares
If you have other middlewares, which you created using other routers, then there is a high chance that you can use it in **There** without changing much.

As an example, let us have a look at the Recoverer middleware.

```go
func Recoverer(request there.Request, next there.Response) there.Response {
   fn := func(w http.ResponseWriter, r *http.Request) {
      defer func() {
         if rvr := recover(); rvr != nil && rvr != http.ErrAbortHandler {
            there.Error(status.InternalServerError, rvr).ServeHTTP(w, r)
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

### Robust

Almost everything in **There** has a corresponding test. We tested the framework well in production and wrote enough test cases, so in the end almost everything that makes sense to test was tested.


# 👨‍💻 Contributions
Feel free to contribute to this project in any way! May it be a simple issue, idea or a finished pull request. Every helpful hand is welcomed.

<a href="https://gitHub.com/Gebes/there/graphs/commit-activity">
    <img src="https://img.shields.io/badge/Maintained%3F-yes-green.svg" alt="Maintained?">
</a>
<a href="https://github.com/Gebes">
    <img src="https://img.shields.io/badge/Maintainer-Gebes-blue" alt="Maintainer">
</a>
</p>