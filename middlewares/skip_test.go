package middlewares

import (
	"errors"
	"github.com/gebes/there/v2/status"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/gebes/there/v2"
)

func TestSkip(t *testing.T) {
	router := there.NewRouter()

	router.Use(Skip(dummyMiddleware, dummyValidatorFalse))

	router.Get("/", dummyEndpointHandler)

	request := httptest.NewRequest(string(there.MethodGet), "/", nil)
	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, request)
	result := recorder.Result()

	reflect.DeepEqual(result.StatusCode, status.InternalServerError)
}

func TestSkipFalse(t *testing.T) {
	router := there.NewRouter()

	router.Use(Skip(dummyMiddleware, dummyValidatorTrue))

	router.Get("/", dummyEndpointHandler)

	request := httptest.NewRequest(string(there.MethodGet), "/", nil)
	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, request)
	result := recorder.Result()

	reflect.DeepEqual(result.Body, dummyData)
}

func TestSkipNil(t *testing.T) {
	router := there.NewRouter()

	router.Use(Skip(dummyMiddleware, nil))

	router.Get("/", dummyEndpointHandler)

	request := httptest.NewRequest(string(there.MethodGet), "/", nil)
	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, request)
	result := recorder.Result()

	reflect.DeepEqual(result.StatusCode, status.InternalServerError)
}

var dummyData = map[string]interface{}{
	"Hello": "There",
}

func dummyMiddleware(request there.Request, next there.Response) there.Response {
	return there.Error(status.InternalServerError, errors.New("i'm just a dummy"))
}

func dummyEndpointHandler(request there.Request) there.Response {
	return there.Json(status.OK, dummyData)
}

func dummyValidatorFalse(request there.Request) bool {
	return false
}

func dummyValidatorTrue(request there.Request) bool {
	return true
}
