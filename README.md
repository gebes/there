# there
`there`, also called `GoThere`, aims to be a simple Go Library to reduce redundant code for REST APIs. The name `GoThere` because it tells the incoming requests to "go there".
With `there` it is a delight to create a REST API. You can focus on writing your API without writing boilerplate and duplicate code. Just start coding and save time!

## Table of contents

* [Installation](#install)
* [Basic Example](#basic-example)

## Install

```sh
go get -u github.com/Gebes/there
```

## Example

### Basic Example
`GET /user/:id`
`POST /user`
`DELETE /user/:id` 
```go
package main

import (
  "errors"
  . "github.com/Gebes/there"
  "log"
)

type User struct {
  Id   string `json:"id"`
  Name string `json:"name"`
}

var users = map[string]User{}

func main() {

  router := NewRouter().
    Get("/user/:id", GetUser).
    Post("/user", PostUser).
    Delete("/user/:id", DeleteUser)

  err := router.Listen(8080)

  if err != nil {
    log.Fatalln("Could not bind router to 8080", err)
  }

}

func GetUser(request HttpRequest) HttpResponse {
  queryId, _ := request.RouteParams.Get("id")
  user, ok := users[queryId]

  if !ok {
    return Error(StatusNotFound, errors.New("user not found"))
  }
  return Json(StatusOK, user)
}

func PostUser(request HttpRequest) HttpResponse {
  var userToAdd User
  err := request.Body.BindJson(&userToAdd)

  if err != nil {
    return Error(StatusBadRequest, err)
  }

  users[userToAdd.Id] = userToAdd
  return Empty(StatusOK)
}

func DeleteUser(request HttpRequest) HttpResponse {
  queryId, _ := request.RouteParams.Get("id")
  delete(users, queryId)
  return Empty(StatusOK)
}

```