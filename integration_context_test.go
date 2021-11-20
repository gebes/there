package there

import (
	"context"
	"errors"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http/httptest"
	"testing"
)

var (
	users = map[string]simpleUser{
		"1": {"Person"},
		"2": {"Hannes"},
	}
)

func createRouter() *Router {
	router := NewRouter()
	router.
		Use(func(request HttpRequest) HttpResponse {
			authorization := request.Headers.GetDefault(RequestHeaderAuthorization, "")

			user, ok := users[authorization]
			if !ok {
				return Error(StatusUnauthorized, errors.New("not authorized"))
			}

			return WithContext(context.WithValue(request.Context(), "user", user), Next())
		})
	router.
		Get("/user", func(request HttpRequest) HttpResponse {
			user, ok := request.Context().Value("user").(simpleUser)

			if !ok {
				return Error(StatusUnprocessableEntity, errors.New("could not get user from context"))
			}

			return Json(StatusOK, user)
		}).With(func(request HttpRequest) HttpResponse {
		// Nested WithContext Test
		ctx := context.WithValue(request.Context(), "world", "hello")
		r := WithContext(ctx, Next())
		return WithContext(context.WithValue(ctx, "hello", "world"), r)
	})

	router.Get("/user/test2", func(request HttpRequest) HttpResponse {
		return Empty(StatusOK)
	}).With(func(request HttpRequest) HttpResponse {
		return WithContext(context.WithValue(request.Context(), "hello", "world"), String(StatusBadRequest, "Error"))
	})

	return router
}

func TestContextMiddleware1(t *testing.T) {

	router := createRouter()

	request := httptest.NewRequest(MethodGet, "/user", nil)
	request.Header.Set("Authorization", "1")
	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, request)

	result := recorder.Result()

	defer result.Body.Close()
	data, err := ioutil.ReadAll(result.Body)
	if err != nil {
		t.Fatalf("could not read body %v", err)
	}
	dat := string(data)
	assert.Equal(t, "{\"Name\":\"Person\"}", dat)
}

func TestContextMiddleware2(t *testing.T) {

	router := createRouter()

	request := httptest.NewRequest(MethodGet, "/user/test2", nil)
	request.Header.Set("Authorization", "1")
	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, request)

	result := recorder.Result()

	defer result.Body.Close()
	data, err := ioutil.ReadAll(result.Body)
	if err != nil {
		t.Fatalf("could not read body %v", err)
	}
	dat := string(data)
	assert.Equal(t, "Error", dat)
}
