<p align="center">
<font size="7px">
<b>
there
</b>
</font>
</p>
<p align="center">
<font size="4px">
⚡️ Robust Web Framework to build Go Services with easier control flow ⚡️
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

The general problem with many routers is the way you handle responses. Most frameworks make it to complex or do not offer the proper abstraction, to get the same result in a short amount of time.  
The goal of **There** is to give developers the right tool, to create robust API's in no time.  

We solve this problem, by providing simple interfaces to control the flow of your API.  
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

The goal of **There** is to safe time and provide all the functionality a backend developer needs on daily basis. **There** should always keep it simple and only add complexity, if there is no good workaround.  
New Go Developers often struggle with chores they are not used to and **There** should be a help for them, to make life easier.

## ⭐️ Features
* [Straightforward **routing**](#straightforward-routing)
* Minimalistic **control flow**
* **Expandable** - add your own control flow
* Complete **middlewares** support
* **Lightweight** - no dependencies
* Robust - 99,6% coverage

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
