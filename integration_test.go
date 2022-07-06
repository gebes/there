package there_test

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http/httptest"
	"reflect"
	"strconv"
	"testing"

	. "github.com/Gebes/there/v2"
	"github.com/Gebes/there/v2/middlewares"
)

var (
	sampleData = map[string]any{
		"Hello": "There",
	}
	sampleUser       = user{"Hannes", "A cool user"}
	sampleSimpleUser = simpleUser{"Hannes"}
	errorData        = map[any]any{}
	errorMarshal     = func(i any) ([]byte, error) {
		return nil, errors.New("test")
	}
)

type user struct {
	Name        string
	Description string
}

type simpleUser struct {
	Name string `yaml:"Name"`
}

func CreateRouter() *Router {
	// Sample Data
	json := func(request HttpRequest) HttpResponse {
		return Json(StatusOK, sampleData)
	}
	// Samle User
	xml := func(request HttpRequest) HttpResponse {
		return Xml(StatusOK, sampleUser)
	}

	router := NewRouter().
		Use(middlewares.Cors(middlewares.AllowAllConfiguration()))

	router.Use(middlewares.Recoverer)

	data := router.Group("/data")

	data.Handle("/json", json, MethodGet, MethodPost, MethodPut, MethodDelete)
	data.Get("/xml", xml)
	data.Get("/empty", func(request HttpRequest) HttpResponse {
		return Status(StatusAccepted)
	})
	data.Get("/message", func(request HttpRequest) HttpResponse {
		return Message(StatusOK, "Hello there")
	})
	data.Get("/string", func(request HttpRequest) HttpResponse {
		return String(StatusOK, "Hello there")
	})
	data.Get("/redirect", func(request HttpRequest) HttpResponse {
		return Redirect(StatusMovedPermanently, "https://google.com")
	})
	data.Get("/html", func(request HttpRequest) HttpResponse {
		return Html(StatusOK, "./test/index.html", map[string]string{
			"user": "Hannes",
		})
	})
	data.Get("/bytes", func(request HttpRequest) HttpResponse {
		return Bytes(StatusOK, []byte{'a', 'b'})
	})

	errorGroup := router.Group("/error")
	errorGroup.Get("/json", func(request HttpRequest) HttpResponse {
		return Json(StatusOK, errorData)
	})
	errorGroup.Get("/xml", func(request HttpRequest) HttpResponse {
		return Xml(StatusOK, errorData)
	})

	errorGroup.Get("/error/1", func(request HttpRequest) HttpResponse {
		return Error(StatusOK, errors.New("test2"))
	})
	errorGroup.Get("/error/2", func(request HttpRequest) HttpResponse {
		return Error(StatusOK, errors.New("test3"))
	})
	errorGroup.Get("/html/1", func(request HttpRequest) HttpResponse {
		return Html(StatusOK, "./non/existing/folder/for/the/test", map[string]string{
			"user": "Hannes",
		})
	})
	errorGroup.Get("/html/2", func(request HttpRequest) HttpResponse {
		return Html(StatusOK, "./examples/index.html", "A string cannot be used as a template, hence this will fail")
	})
	errorGroup.Get("/data", func(request HttpRequest) HttpResponse {
		return Status(StatusOK)
	}).With(func(request HttpRequest, next HttpResponse) HttpResponse {
		return Error(StatusInternalServerError, errors.New("lol"))
	})

	data.Post("/return/json", func(request HttpRequest) HttpResponse {
		var user simpleUser
		err := request.Body.BindJson(&user)
		if err != nil {
			log.Fatalln("Could not bind", err)
		}
		return String(StatusOK, user.Name)
	})
	data.Post("/return/xml", func(request HttpRequest) HttpResponse {
		var user simpleUser
		err := request.Body.BindXml(&user)
		if err != nil {
			log.Fatalln("Could not bind", err)
		}
		return String(StatusOK, user.Name)
	})
	data.Post("/return/string", func(request HttpRequest) HttpResponse {
		body, err := request.Body.ToString()
		if err != nil {
			log.Fatalln("Could not bind", err)
		}
		return String(StatusOK, body)
	})

	return router.Router
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

func readAndUnmarshal(router *Router, t *testing.T, method, route string, body io.Reader, unmarshal func(data []byte, v any) error, res any) {
	data := readBody(router, t, method, route, body)
	err := unmarshal(data, res)

	if err != nil {
		t.Fatalf("could not parse body %v, %v", err, string(data))
	}

}

func readJsonBody(router *Router, t *testing.T, method, route string, body io.Reader, res any) {
	readAndUnmarshal(router, t, method, route, body, json.Unmarshal, res)
}
func readXmlBody(router *Router, t *testing.T, method, route string, body io.Reader, res any) {
	readAndUnmarshal(router, t, method, route, body, xml.Unmarshal, res)
}

func TestJson(t *testing.T) {
	router := CreateRouter()

	methods := []string{MethodGet, MethodPost, MethodPut, MethodDelete}

	for _, method := range methods {
		var res map[string]any
		readJsonBody(router, t, method, "/data/json", nil, &res)

		if !reflect.DeepEqual(res, sampleData) {
			t.Fatal(res, "does not equal", sampleData)
		}
	}

}

func testSampleUserRoutes(t *testing.T, route string, handler func(router *Router, t *testing.T, method, route string, body io.Reader, res any), res, expected any) {
	router := CreateRouter()
	handler(router, t, MethodGet, "/data/"+route, nil, res)

	if !reflect.DeepEqual(res, expected) {
		t.Fatal(res, "does not equal", expected)
	}

}

func TestXml(t *testing.T) {
	var res user
	testSampleUserRoutes(t, "xml", readXmlBody, &res, &sampleUser)
}

func testErrorResponse(router *Router, t *testing.T, route string) {
	var res map[string]any
	readJsonBody(router, t, MethodGet, "/error/"+route, nil, &res)

	_, ok := res["error"]
	if !ok {
		log.Fatalln("no error key provided", res)
	}
}

func TestNotFoundError(t *testing.T) {
	router := CreateRouter()
	testErrorResponse(router, t, "not_existing_route")
}

func TestJsonErrorResponse(t *testing.T) {
	router := CreateRouter()
	testErrorResponse(router, t, "json")
}

func TestXmlErrorResponse(t *testing.T) {
	router := CreateRouter()
	testErrorResponse(router, t, "xml")
}

func TestHtmlErrorResponse(t *testing.T) {
	router := CreateRouter()
	testErrorResponse(router, t, "html/1")
}
func TestHtml2ErrorResponse(t *testing.T) {
	router := CreateRouter()
	testErrorResponse(router, t, "html/2")
}

func TestErrorErrorResponse2(t *testing.T) {
	router := CreateRouter()
	testErrorResponse(router, t, "error/2")
}

func TestMiddlewareErrorResponse(t *testing.T) {
	router := CreateRouter()
	testErrorResponse(router, t, "data")
}

func TestGlobalMiddlewareErrorResponse(t *testing.T) {
	router := NewRouter().
		Get("error/data/global", func(request HttpRequest) HttpResponse {
			return Status(StatusOK)
		}).
		Use(func(request HttpRequest, next HttpResponse) HttpResponse {
			return Error(StatusInternalServerError, errors.New("errored out"))
		})
	testErrorResponse(router, t, "data/global")
}

func TestPanicErrorResponse(t *testing.T) {
	router := NewRouter().
		Get("error/data/panic", func(request HttpRequest) HttpResponse {
			panic("oh no panic")
		}).
		Use(func(request HttpRequest, next HttpResponse) HttpResponse {
			return Error(StatusInternalServerError, errors.New("errored out"))
		})
	testErrorResponse(router, t, "data/panic")
}

type errReader int

func (errReader) Read(p []byte) (n int, err error) {
	return 0, errors.New("test error")
}

func TestBodyToStringError(t *testing.T) {
	router := NewRouter()
	router.
		Post("/test", func(request HttpRequest) HttpResponse {

			tests := 3
			did := 0

			var s any

			_, err := request.Body.ToString()
			if err != nil {
				did++
			}
			err = request.Body.BindJson(&s)
			if err != nil {
				did++
			}
			err = request.Body.BindXml(&s)
			if err != nil {
				did++
			}

			if tests != did {
				err := fmt.Errorf("not every bind threw an error: \"%v\"/\"%v", strconv.Itoa(did), strconv.Itoa(tests))
				return Error(StatusInternalServerError, err)
			}

			return Status(StatusOK)
		})

	res := readStringBody(router, t, MethodPost, "/test", errReader(0))
	if len(res) != 0 {
		log.Fatalln("res was empty but shouldn't be", res)
	}
}

func TestStringResponse(t *testing.T) {
	router := CreateRouter()
	r := readBody(router, t, MethodGet, "/data/string", nil)
	res := string(r)
	shouldBe := "Hello there"
	if res != shouldBe {
		log.Fatalln("res was", res, "and not", shouldBe)
	}
}

func TestJsonResponse(t *testing.T) {

	router := CreateRouter()
	var jsonBody map[string]any
	readJsonBody(router, t, MethodGet, "/data/message", nil, &jsonBody)

	v, ok := jsonBody["message"]
	if !ok {
		log.Fatalln("key message not present")
	}
	shouldBe := "Hello there"

	if v != shouldBe {
		log.Fatalln("value was", v, "and not", shouldBe)
	}

}

func TestEmptyResponse(t *testing.T) {
	router := CreateRouter()
	r := readBody(router, t, MethodGet, "/data/empty", nil)
	res := string(r)
	if len(res) != 0 {
		log.Fatalln("res was", res, "and not empty \"\"")
	}
}

func TestRedirectResponse(t *testing.T) {
	router := CreateRouter()

	request := httptest.NewRequest(MethodGet, "/data/redirect", nil)
	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, request)

	// result := recorder.Result()

	// TODO FIX ASSERT
	// assert.Equal(t, "https://google.com", result.WithHeaders.Get("Location"))

}

func TestHtmlResponse(t *testing.T) {
	router := CreateRouter()
	r := readBody(router, t, MethodGet, "/data/html", nil)
	res := string(r)
	AssertEquals(t, "Hello Hannes", res)
}
func TestBytesResponse(t *testing.T) {
	router := CreateRouter()
	r := readBody(router, t, MethodGet, "/data/bytes", nil)
	res := string(r)
	AssertEquals(t, "ab", res)
}

func readStringBody(router *Router, t *testing.T, method, route string, body io.Reader) string {

	request := httptest.NewRequest(method, route, body)
	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, request)

	result := recorder.Result()

	defer result.Body.Close()
	data, err := ioutil.ReadAll(result.Body)
	if err != nil {
		t.Fatalf("could not read body %v", err)
	}
	return string(data)
}

func testBind(t *testing.T, marshaller func(v any) ([]byte, error), route string) {
	router := CreateRouter()
	data, err := marshaller(sampleSimpleUser)
	if err != nil {
		t.Fatal(err)
	}
	AssertEquals(t, readStringBody(router, t, MethodPost, "/data/return/"+route, bytes.NewReader(data)), sampleSimpleUser.Name)
}

func TestJsonBodyBind(t *testing.T) {
	testBind(t, json.Marshal, "json")
}

func TestXmlBodyBind(t *testing.T) {
	testBind(t, xml.Marshal, "xml")
}

func TestStringBodyBind(t *testing.T) {
	router := CreateRouter()
	s := "Hello !!!"
	AssertEquals(t, readStringBody(router, t, MethodPost, "/data/return/string", bytes.NewReader([]byte(s))), s)
}
