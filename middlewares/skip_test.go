package middlewares

import (
	"errors"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/Gebes/there/v2"
)

func TestSkip(t *testing.T) {
	router := there.NewRouter()

	router.Use(Skip(dummyMiddleware, func(request there.Request) bool {
		return false
	}))

	router.Get("/", dummyEndpointHandler)

	request := httptest.NewRequest(there.MethodGet, "/", nil)
	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, request)
	result := recorder.Result()

	reflect.DeepEqual(result.StatusCode, there.StatusInternalServerError)
}

func TestSkipFalse(t *testing.T) {
	router := there.NewRouter()

	router.Use(Skip(dummyMiddleware, func(request there.Request) bool {
		return true
	}))

	router.Get("/", dummyEndpointHandler)

	request := httptest.NewRequest(there.MethodGet, "/", nil)
	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, request)
	result := recorder.Result()

	reflect.DeepEqual(result.StatusCode, there.StatusOK)
}

func TestSkipNil(t *testing.T) {
	router := there.NewRouter()

	router.Use(Skip(dummyMiddleware, nil))

	router.Get("/", dummyEndpointHandler)

	request := httptest.NewRequest(there.MethodGet, "/", nil)
	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, request)
	result := recorder.Result()

	reflect.DeepEqual(result.StatusCode, there.StatusInternalServerError)
}

var dummyData  = map[string]interface{}{
	"Hello": "There",
}

func dummyMiddleware(request there.Request, next there.Response) there.Response {
	return there.Error(there.StatusInternalServerError, errors.New("I'm just a dummy"))
}

func dummyEndpointHandler(request there.Request) there.Response {
	return there.Json(there.StatusOK, "Skip")
}
