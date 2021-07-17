package test

import (
	"github.com/Gebes/there/examples"
	"testing"
)

func TestGeneral(t *testing.T) {


	err := examples.CorsRouter()
	if err != nil {
		panic(err)
	}



}


