package middlewares

import (
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/Gebes/there/v2"
)

func TestCorsMiddleware(t *testing.T) {

	router := there.NewRouter()
	router.Use(Cors(AllowAllConfiguration()))
	router.Get("/", func(request there.Request) there.Response {
		return there.Status(there.StatusOK)
	})

	request := httptest.NewRequest(there.MethodGet, "/", nil)
	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, request)

	result := recorder.Result()
	checkHeaders(t, result)

	request = httptest.NewRequest(there.MethodOptions, "/", nil)
	recorder = httptest.NewRecorder()

	router.ServeHTTP(recorder, request)

	result = recorder.Result()
	checkHeaders(t, result)

}

func checkHeaders(t *testing.T, result *http.Response) {
	if !reflect.DeepEqual(result.Header.Get(ResponseHeaderAccessControlAllowOrigin), "*") ||
		!reflect.DeepEqual(result.Header.Get(ResponseHeaderAccessControlAllowMethods), AllMethodsString) ||
		!reflect.DeepEqual(result.Header.Get(ResponseHeaderAccessControlAllowHeaders), "Accept, Content-Type, Content-Length, Authorization") {
		t.Fatal("headers did not match allow all configuration", result.Header)
	}
}
