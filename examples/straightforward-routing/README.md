# ⭐️ Straightforward routing

Routing with **There** is easy! Simply create a new router and add `GET`, `POST` or different handlers. Use groups to define multiple handlers simultaneously and protect your handlers with middlewares.
Define route variables with `:` and you have all the things you need for routing.

## 👆 Style Best Practices

With **There** you have multiple "design patterns" to register groups. However, all options produce the same outcome.  
I personally like the first option the most because it gets properly formatted by Go. Option three would make more sense, but it adds more complexity and isn't as easy to read as option one.

Go with the option you like the most, but be aware of readability problems!

### ✅ Good
This example is the easiest to understand
```go
router.Group("/user").
    Get("/all", Handler).
    Get("/:id", Handler).
    Post("/", Handler).
    Patch("/:id", Handler).
    Delete("/:id", Handler)

router.Group("/user/post").
    Get("/:id", Handler).
    Post("/", Handler).
    Patch("/:id", Handler).
    Delete("/:id", Handler)

router.Get("/details", Handler)
```

### ❌ Bad
Most linters will break the manual intends you placed.
```go
router.Group("/user").
    Get("/all", Handler).
    Get("/:id", Handler).
    Post("/", Handler).
    Patch("/:id", Handler).
    Delete("/:id", Handler).
    Group("/post"). 
        Get("/:id", Handler).
        Post("/", Handler).
        Patch("/:id", Handler).
        Delete("/:id", Handler)

router.Get("/details", Handler)
```


Requires more time to read and can lead more easily to evil bugs.
```go
user := router.Group("/user").
    Get("/all", Handler).
    Get("/:id", Handler).
    Post("/", Handler).
	Patch("/:id", Handler).
    Delete("/:id", Handler)

user.Group("/post").
    Get("/:id", Handler).
    Post("/", Handler).
    Patch("/:id", Handler).
    Delete("/:id", Handler)

router.Get("/details", Handler)
```


Looks stuffed and typos can happen more easily.
```go
router.Get("/user/all", Handler).
router.Get("/user/:id", Handler).
router.Post("/user", Handler).
router.ch("/user/:id", Handler).
router.Delete("/user/:id", Handler)

router.Get("/user/post/:id", Handler)
router.Post("/user/post", Handler)
router.Patch("/user/post/:id", Handler)
router.Delete("/user/post/:id", Handler)

router.Get("/details", Handler)
```


## 👆 Route Variables Best Practices

You can have as many path variables as you want. However, you need to consider one thing for naming.

### ✅ Good
```go
router.Group("/user").
    Get("/:id", Handler).
    Patch("/:id", Handler).
    Delete("/:id", Handler)

router.Group("/student").
    Get("/:name", Handler).
    Patch("/:name", Handler)
```

### ❌ Bad
This example won't work. 
```go
router.Group("/user").
    Get("/:id", Handler).
    Patch("/:name", Handler).
    Delete("/:serial", Handler)

router.Group("/student").
    Get("/:name", Handler).
    Patch("/:name", Handler)
```
It will compile, but the router will return an error.
```
path variable "name" for path "/user/:name" needs to equal "id" as in all other routes
```
If you have a parent path `/user/:id` you can't define a sub path with a different parameter name `/user/:name/test`. However, it is possible to have multiple route parameters `/user/:userId/post/:postId`.
