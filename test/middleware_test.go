package test

import (
	"github.com/Gebes/there/examples"
	"testing"
)

func TestMiddleware(t *testing.T) {


	err := unblock(func() error {
		return examples.MiddlewareRouter()
	})
	if err != nil {
		panic(err)
	}

	get("http://localhost:8080/message")
	get("http://localhost:8080/message/without/middleware")

}


