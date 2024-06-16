package main

import (
	"errors"
	"fmt"
	"github.com/gebes/there/v2/status"
	"log"

	"github.com/gebes/there/v2"
)

type (
	Post struct {
		Id      string `json:"id,omitempty"`
		Title   string `json:"title,omitempty"`
		Content string `json:"content,omitempty"`
	}
	Posts []Post
)

var posts = make(Posts, 0)

func main() {
	router := there.NewRouter()

	router.Group("/post").
		Get("/", GetPosts).
		Get("/{id}", GetPostById).
		Post("/", CreatePost)

	err := router.Listen(8080)
	if err != nil {
		log.Fatalln("Could not listen to 8080:", err)
	}
}

func GetPosts(request there.Request) there.Response {
	return there.Json(status.OK, posts)
}

func GetPostById(request there.Request) there.Response {
	id := request.RouteParams.Get("id")

	post := postById(id)
	if post == nil {
		return there.Error(status.NotFound, errors.New("could not find post"))
	}

	return there.Json(status.OK, post)
}

func CreatePost(request there.Request) there.Response {
	var body Post
	err := request.Body.BindJson(&body) // Decode body
	if err != nil {                     // If body was not valid json, return bad request error
		return there.Error(status.BadRequest, fmt.Errorf("could not parse body: %w", err))
	}

	post := postById(body.Id)
	if post != nil { // if the post already exists, return conflict error
		return there.Error(status.Conflict, errors.New("post with this ID already exists"))
	}

	posts = append(posts, body) // create post

	return there.Json(status.Created, body) // return created post as json
}

func postById(id string) *Post {
	var post *Post
	for _, current := range posts {
		if current.Id == id {
			post = &current
			break
		}
	}
	return post
}

func ExampleGet(request there.Request) there.Response {
	user := map[string]string{
		"firstname": "John",
		"surname":   "Smith",
	}
	return there.Json(status.OK, user)
}
