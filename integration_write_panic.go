package there_test

import (
	"encoding/json"
	. "github.com/Gebes/there/v2"
	"io"
	"io/ioutil"
	"log"
	"net/http/httptest"
	"testing"
)

func TestWriteError(t *testing.T) {
	router := NewRouter()
	router.Get("/", func(request HttpRequest) HttpResponse {
		// No body with 1xx status
		return Json(StatusContinue, "not writeable")
	})
	var data interface{}
	readJsonBody(router, t, MethodGet, "/", nil, &data)
	log.Println(data)
}

func readBody(router *Router, t *testing.T, method, route string, body io.Reader) []byte {

	request := httptest.NewRequest(method, route, body)
	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, request)

	result := recorder.Result()

	defer result.Body.Close()
	data, err := ioutil.ReadAll(result.Body)
	if err != nil {
		t.Fatalf("could not read body %v", err)
	}
	return data
}

func readAndUnmarshal(router *Router, t *testing.T, method, route string, body io.Reader, unmarshal func(data []byte, v interface{}) error, res interface{}) {
	data := readBody(router, t, method, route, body)
	err := unmarshal(data, res)

	if err != nil {
		t.Fatalf("could not parse body %v, %v", err, string(data))
	}

}

func readJsonBody(router *Router, t *testing.T, method, route string, body io.Reader, res interface{}) {
	readAndUnmarshal(router, t, method, route, body, json.Unmarshal, res)
}
