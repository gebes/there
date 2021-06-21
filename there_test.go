package there

import (
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"testing"
	"time"
)

func TestGoThere(t *testing.T) {
	router := Router{Port: 8080}
	router.Handle("/monkeys", func(request Request) Response {
		return ResponseData(http.StatusOK, "I like monkeys too")
	}, http.MethodGet)
	router.Handle("/monkey", func(request Request) Response {
		return ResponseData(http.StatusOK, map[string]string{
			"name":           "Alfred",
			"age":            "Unkown",
			"favourite_meal": "Affensuppe",
		})
	}, http.MethodGet)
	router.Handle("/user", func(request Request) Response {
		return Response{Status: http.StatusOK, Data: map[string]string{
			"name": "Henrick",
		}}
	})

	err := unblock(func() error {
		return router.Listen()
	})
	if err != nil {
		panic(err)
	}

	get("http://localhost:8080/monkeys?test=123")
	get("http://localhost:8080/monkeys")
	get("http://localhost:8080/monkey")
}

func get(route string) {
	get, err := http.Get(route)
	if err != nil {
		log.Fatalln("Could not make RawRequest", err)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			panic(err)
		}
	}(get.Body)

	bodyBytes, _ := ioutil.ReadAll(get.Body)
	body := string(bodyBytes)
	log.Println("Body", body)
}

func unblock(h func() error) error {
	w := make(chan error)
	go func() {
		w <- h()
	}()
	select {
	case err := <-w:
		return err
	case <-time.After(time.Millisecond * 50):
		return nil
	}
}
