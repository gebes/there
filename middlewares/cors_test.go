package middlewares

import (
	"github.com/Gebes/there/v2/header"
	"github.com/Gebes/there/v2/status"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/Gebes/there/v2"
)

func TestCorsMiddleware(t *testing.T) {

	router := there.NewRouter()
	router.Use(Cors(CorsAllowAllConfiguration()))
	router.Get("/", func(request there.Request) there.Response {
		return there.Status(status.OK)
	})

	request := httptest.NewRequest(string(there.MethodGet), "/", nil)
	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, request)

	result := recorder.Result()
	checkHeaders(t, result)

	request = httptest.NewRequest(string(there.MethodOptions), "/", nil)
	recorder = httptest.NewRecorder()

	router.ServeHTTP(recorder, request)

	result = recorder.Result()
	checkHeaders(t, result)

}

func checkHeaders(t *testing.T, result *http.Response) {
	if !reflect.DeepEqual(result.Header.Get(header.ResponseAccessControlAllowOrigin), "*") ||
		!reflect.DeepEqual(result.Header.Get(header.ResponseAccessControlAllowMethods), there.AllMethodsJoined) ||
		!reflect.DeepEqual(result.Header.Get(header.ResponseAccessControlAllowHeaders), "Accept, Content-Type, Content-Length, Authorization") {
		t.Fatal("headers did not match allow all configuration", result.Header)
	}
}
