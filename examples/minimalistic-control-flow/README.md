# ⭐️ Minimalistic control flow


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

[📚 Documentation](https://there.gebes.io/responses/status)  
