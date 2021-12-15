package middlewares

import (
	. "github.com/Gebes/there/v2"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func TestCorsMiddleware(t *testing.T) {

	router := NewRouter()
	router.Use(Cors(AllowAllConfiguration()))
	router.Get("/", func(request HttpRequest) HttpResponse {
		return Status(StatusOK)
	})

	request := httptest.NewRequest(MethodGet, "/", nil)
	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, request)

	result := recorder.Result()
	checkHeaders(t, result)

	request = httptest.NewRequest(MethodOptions, "/", nil)
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
