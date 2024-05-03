<p align="center">
<a href="https://there.gebes.io">
<img src="https://github.com/Gebes/there/blob/main/.github/logo/logo-compressed.png?raw=true" alt="Gopher There" height="256">
</a>
</p>
<p align="center">
<font size="4px">
Effortless Go Routing: Build Powerful Services with Minimal Overhead
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
3. Install `there` with the `go get` command

```sh
go get -u github.com/Gebes/there/v2
```

4. Create a `main.go` file

```go
package main

import (
	"github.com/Gebes/there/v2"
	"github.com/Gebes/there/v2/status"
	"log"
)

func main() {
	router := there.NewRouter() 

	router.Get("/", func(request there.Request) there.Response {
		return there.Json(status.OK, map[string]string{
			"message": "Hello World!",
		})
	})

	err := router.Listen(8080)
	if err != nil {
		log.Fatalln(err)
	}
}
```

[📚 Check the Documentation for more info](https://there.gebes.io)

# 🚀 Introducing `there` - Simplify Your Go Routing

Welcome to `there`, a Go routing library designed to streamline your project's API development. By wrapping the default http mux, `there` enhances control flow without stripping away the flexibility and power of standard HTTP capabilities.

## Why Choose `there`?

Developing robust APIs doesn't have to be complicated or time-consuming. With `there`, you can manage response handling intuitively and efficiently, avoiding the common pitfalls and verbosity typical of many routing frameworks.

- **Error Handling Simplified**: Encounter an error? Respond with `return Error(status, err)`.
- **Effortless Data Response**: Need to send data? Use `return Json(status, data)`.
- **Optimize Large Data Transfers**: Have large data? Compress easily with `return Gzip(Json(status, data))`.

This approach ensures your code is cleaner, more readable, and maintains full HTTP functionality.


## ⭐️ Key Features

### Minimalistic Control Flow
Enjoy simple and clear control flow in your routes with out-of-the-box handlers that make coding a breeze:
```go
func CreatePost(request there.Request) there.Response {
    var body Post
    err := request.Body.BindJson(&body)
    if err != nil {
        return there.Error(status.BadRequest, "Could not parse body: " + err.Error())
    }
    if postById(body.Id) != nil {
        return there.Error(status.Conflict, "Post with this ID already exists")
    }
    posts = append(posts, body)
    return there.Json(status.Created, body)
}
```
[🧑‍💻 View full code example](https://github.com/Gebes/there/tree/main/examples/minimalistic-control-flow)

### Straightforward Routing
`there` makes routing straightforward. Define routes and methods easily with groupings and middleware support:
```go
router := there.NewRouter()

router.Group("/user")
    .Get("/", Handler)
    .Post("/", Handler)
    .Patch("/", Handler)

router.Group("/user/post")
    .Get("/{id}", Handler)
    .Post("/", Handler)

router.Get("/details", Handler)
```
[🧑‍💻 View full code example and best practices](https://github.com/Gebes/there/tree/main/examples/straightforward-routing)



### Expandable - Customize Your Control Flow
Easily extend `there` by adding your own responses like a Msgpack response, without the need for third-party dependencies:
```go
import (
    "github.com/Gebes/there/v2"
    "github.com/vmihailenco/msgpack/v5"
)

func Msgpack(code int, data interface{}) there.Response {
    msgpackData, err := msgpack.Marshal(data)
    if err != nil {
        panic(err)
    }
    return there.WithHeaders(map[string]string{
        there.ResponseHeaderContentType: "application/x-msgpack",
    }, there.Bytes(code, msgpackData))
}
```
[🧑‍💻 View full code example](https://github.com/Gebes/there/tree/main/examples/custom-http-response)

## 🌟 Additional Features for the Skeptical Gophers

### Complete Middleware Support
Utilize existing middleware with minimal changes or create specific middleware for individual routes, enhancing flexibility and control over the request/response lifecycle.

### Lightweight - No Dependencies
`there` is designed to be dependency-free, ensuring a lightweight integration into your projects without bloating your application.

### Robust and Tested
With nearly complete test coverage, `there` is proven in production environments, offering reliability and stability for your critical applications.

### Easy Integration with Existing Code
Thanks to `there`'s compatibility and flexible architecture, integrating with existing Go codebases and middleware is seamless and straightforward.

### Performance

`there` is as efficient as it gets—built on the default http mux, it matches routes without any performance overhead. Just focus on building your application; `there` handles the rest.

Check out our benchmarks for more details:
- [Benchmark Analysis](https://pastebin.com/iP7NhtZH)
- [Performance Results](https://pastebin.com/ETWF8cqt)

# 👨‍💻 Contributions
Feel free to contribute to this project in any way! May it be a simple issue, idea or a finished pull request. Every helpful hand is welcomed.

<a href="https://gitHub.com/Gebes/there/graphs/commit-activity">
    <img src="https://img.shields.io/badge/Maintained%3F-yes-green.svg" alt="Maintained?">
</a>
<a href="https://github.com/Gebes">
    <img src="https://img.shields.io/badge/Maintainer-Gebes-blue" alt="Maintainer">
</a>
</p>