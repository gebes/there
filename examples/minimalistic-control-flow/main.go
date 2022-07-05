package main

import (
	"errors"
	"log"

	. "github.com/Gebes/there/v2"
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
	router := NewRouter()

	router.Group("/post").
		Get("/", GetPosts).
		Get("/:id", GetPostById).
		Post("/", CreatePost)

	err := router.Listen(8080)
	if err != nil {
		log.Fatalln("Could not listen to 8080", err)
	}
}

func GetPosts(request HttpRequest) HttpResponse {
	return Json(StatusOK, posts)
}

func GetPostById(request HttpRequest) HttpResponse {
	id := request.RouteParams.GetDefault("id", "")

	post := postById(id)
	if post == nil {
		return Error(StatusNotFound, errors.New("Could not find post"))
	}

	return Json(StatusOK, post)
}

func CreatePost(request HttpRequest) HttpResponse {
	var body Post
	err := request.Body.BindJson(&body) // Decode body
	if err != nil {                     // If body was not valid json, return bad request error
		return Error(StatusBadRequest, errors.New("Could not parse body: "+err.Error()))
	}

	post := postById(body.Id)
	if post != nil { // if the post already exists, return conflict error
		return Error(StatusConflict, errors.New("Post with this ID already exists"))
	}

	posts = append(posts, body)      // create post
	return Json(StatusCreated, body) // return created post as json
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
