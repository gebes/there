package test

import (
	"github.com/Gebes/there/examples"
	"testing"
)

func TestRedirect(t *testing.T) {

	// open http://localhost:8080/google

	err := examples.RedirectRouter()
	if err != nil {
		panic(err)
	}



}


