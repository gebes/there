package there_test

import (
	"context"
	"errors"
	"io/ioutil"
	"net/http/httptest"
	"testing"

	"github.com/Gebes/there/v2"
	"github.com/Gebes/there/v2/middlewares"
)

var (
	users = map[string]simpleUser{
		"1": {"Person"},
		"2": {"Hannes"},
	}
)

func createRouter() *there.Router {
	router := there.NewRouter()
	router.Use(middlewares.Recoverer)
	router.
		Use(func(request there.Request, next there.Response) there.Response {
			authorization := request.Headers.GetDefault(there.RequestHeaderAuthorization, "")

			user, ok := users[authorization]
			if !ok {
				return there.Error(there.StatusUnauthorized, errors.New("not authorized"))
			}

			request.WithContext(context.WithValue(request.Context(), "user", user))
			return next
		})
	router.
		Get("/user", func(request there.Request) there.Response {
			user, ok := request.Context().Value("user").(simpleUser)

			if !ok {
				return there.Error(there.StatusUnprocessableEntity, errors.New("could not get user from context"))
			}

			return there.Json(there.StatusOK, user)
		}).With(func(request there.Request, next there.Response) there.Response {

		request.WithContext(context.WithValue(request.Context(), "world", "hello"))
		request.WithContext(context.WithValue(request.Context(), "hello", "world"))
		return next
	})

	router.Get("/user/test2", func(request there.Request) there.Response {
		return there.Status(there.StatusOK)
	}).With(func(request there.Request, next there.Response) there.Response {
		request.WithContext(context.WithValue(request.Context(), "hello", "world"))
		return there.String(there.StatusBadRequest, "there.Error")
	})

	return router
}

func TestContextMiddleware1(t *testing.T) {

	router := createRouter()

	request := httptest.NewRequest(there.MethodGet, "/user", nil)
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

	request := httptest.NewRequest(there.MethodGet, "/user/test2", nil)
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
	want := "there.Error"
	if got != want {
		t.Errorf("%v != %v", want, got)
	}
}
