package test

import (
	"github.com/Gebes/there/examples"
	"testing"
)

func TestGlobalMiddleware(t *testing.T) {


	err := unblock(func() error {
		return examples.GlobalMiddlewareRouter()
	})
	if err != nil {
		panic(err)
	}

	get("http://localhost:8080/message")
}


