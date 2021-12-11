# ⭐️ Straightforward routing

Routing with **There** is easy! Simply create a new router and add `GET`, `POST` or different handlers. Use groups to define multiple handlers simultaneously and protect your handlers with middlewares.
Define route variables with `:` and you have all the things you need for routing.

## Why three options?

With **There** you have multiple "design patterns" to register groups. However, all options produce the same outcome.  
I personally like option one the most because it gets properly formatted by Go. Option three would make more sense, but it adds more complexity and isn't as easy to read as option one.

Go with the option you like the most, but stick to one!


### Option 1
```go
	router.Group("/user").
		Get("/", Handler). // /user
		Post("/", Handler).
		Patch("/", Handler)

	router.Group("/user/post").
		Get("/:id", Handler).
		Post("/", Handler)
```

### Option 2
```go
	router.Group("/user").
		Get("/", Handler). // /user
		Post("/", Handler).
		Patch("/", Handler).
		Group("/post").
			Get("/:id", Handler).
			Post("/", Handler)
```

### Option 3
```go
	user := router.Group("/user").
		Get("/", Handler). // /user
		Post("/", Handler).
		Patch("/", Handler)

	user.Group("/post").
		Get("/:id", Handler).
		Post("/", Handler)
```