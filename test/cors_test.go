package test

import (
	"github.com/Gebes/there/examples"
	"testing"
)

func TestCors(t *testing.T) {


	err := unblock(func() error {
		return examples.CorsRouter()
	})
	if err != nil {
		panic(err)
	}

	get("http://localhost:8080/cors")


}


