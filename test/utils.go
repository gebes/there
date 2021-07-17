package test

import (
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

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