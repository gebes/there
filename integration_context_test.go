package there_test

import (
	"context"
	"errors"
	. "github.com/Gebes/there/v2"
	"github.com/Gebes/there/v2/middlewares"
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
	router.Use(middlewares.Recoverer)
	router.
		Use(func(request HttpRequest, next HttpResponse) HttpResponse {
			authorization := request.Headers.GetDefault(RequestHeaderAuthorization, "")

			user, ok := users[authorization]
			if !ok {
				return Error(StatusUnauthorized, errors.New("not authorized"))
			}

			request.WithContext(context.WithValue(request.Context(), "user", user))
			return next
		})
	router.
		Get("/user", func(request HttpRequest) HttpResponse {
			user, ok := request.Context().Value("user").(simpleUser)

			if !ok {
				return Error(StatusUnprocessableEntity, errors.New("could not get user from context"))
			}

			return Json(StatusOK, user)
		}).With(func(request HttpRequest, next HttpResponse) HttpResponse {

		request.WithContext(context.WithValue(request.Context(), "world", "hello"))
		request.WithContext(context.WithValue(request.Context(), "hello", "world"))
		return next
	})

	router.Get("/user/test2", func(request HttpRequest) HttpResponse {
		return Status(StatusOK)
	}).With(func(request HttpRequest, next HttpResponse) HttpResponse {
		request.WithContext(context.WithValue(request.Context(), "hello", "world"))
		return String(StatusBadRequest, "Error")
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
	got := string(data)
	want := "{\"Name\":\"Person\"}"
	if got != want {
		t.Errorf("%v != %v", want, got)
	}
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
	got := string(data)
	want := "Error"
	if got != want {
		t.Errorf("%v != %v", want, got)
	}
}
