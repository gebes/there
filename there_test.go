package there

import (
	"context"
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"github.com/Gebes/there/v2/header"
	"github.com/Gebes/there/v2/status"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strconv"
	"testing"
	"time"
)

// Empty handler
func handler(request Request) Response {
	return Status(status.OK)
}

// Empty middleware
func middleware(request Request, next Response) Response {
	return next
}

func assertBodyResponse(t *testing.T, r http.Handler, method, route string, expected string) {
	request := httptest.NewRequest(method, route, nil)
	recorder := httptest.NewRecorder()
	r.ServeHTTP(recorder, request)
	result := recorder.Result()
	re, err := io.ReadAll(result.Body)
	if err != nil {
		t.Fatalf("could not read body: %v", err)
	}
	err = result.Body.Close()
	if err != nil {
		t.Fatalf("could not close body: %v", err)
	}
	if string(re) != expected {
		t.Fatalf("invalid response. actual: %v, expected: %v", string(re), expected)
	}
}

// constants.go tests

// Integration tests

func TestWriteError(t *testing.T) {
	router := NewRouter()
	router.Get("/", func(request Request) Response {
		// No body with 1xx status
		return Json(status.Continue, "not writeable")
	})
	var data any
	readJsonBody(router, t, MethodGet, "/", nil, &data)
	log.Println(data)
}

var (
	sampleData = map[string]any{
		"Hello": "There",
	}
	sampleUser       = user{"Hannes", "A cool user"}
	sampleSimpleUser = simpleUser{"Hannes"}
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
	json := func(request Request) Response {
		return Json(status.OK, sampleData)
	}
	// Samle User
	xml := func(request Request) Response {
		return Xml(status.OK, sampleUser)
	}

	router := NewRouter()

	data := router.Group("/data")

	data.Handle("/json", json, MethodGet, MethodPost, MethodPut, MethodDelete)
	data.Get("/xml", xml)
	data.Get("/empty", func(request Request) Response {
		return Status(status.Accepted)
	})
	data.Get("/message", func(request Request) Response {
		return Message(status.OK, "Hello there")
	})
	data.Get("/string", func(request Request) Response {
		return String(status.OK, "Hello there")
	})
	data.Get("/redirect", func(request Request) Response {
		return Redirect(status.MovedPermanently, "https://google.com")
	})
	/*	data.Get("/html", func(request Request) Response {
		return Html(status.OK, "./test/index.html", map[string]string{
			"user": "Hannes",
		})
	})*/
	data.Get("/bytes", func(request Request) Response {
		return Bytes(status.OK, []byte{'a', 'b'})
	})

	data.Post("/return/json", func(request Request) Response {
		var user simpleUser
		err := request.Body.BindJson(&user)
		if err != nil {
			log.Fatalln("Could not bind", err)
		}
		return String(status.OK, user.Name)
	})
	data.Post("/return/xml", func(request Request) Response {
		var user simpleUser
		err := request.Body.BindXml(&user)
		if err != nil {
			log.Fatalln("Could not bind", err)
		}
		return String(status.OK, user.Name)
	})
	data.Post("/return/string", func(request Request) Response {
		body, err := request.Body.ToString()
		if err != nil {
			log.Fatalln("Could not bind", err)
		}
		return String(status.OK, body)
	})

	return router
}

func readBody(router *Router, t *testing.T, method Method, route string, body io.Reader) []byte {

	request := httptest.NewRequest(string(method), route, body)
	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, request)

	result := recorder.Result()

	defer result.Body.Close()
	data, err := io.ReadAll(result.Body)
	if err != nil {
		t.Fatalf("could not read body %v", err)
	}
	return data
}

func readAndUnmarshal(router *Router, t *testing.T, method Method, route string, body io.Reader, unmarshal func(data []byte, v any) error, res any) {
	data := readBody(router, t, method, route, body)
	err := unmarshal(data, res)

	if err != nil {
		t.Fatalf("could not parse body %v, %v", err, string(data))
	}

}

func readJsonBody(router *Router, t *testing.T, method Method, route string, body io.Reader, res any) {
	readAndUnmarshal(router, t, method, route, body, json.Unmarshal, res)
}
func readXmlBody(router *Router, t *testing.T, method Method, route string, body io.Reader, res any) {
	readAndUnmarshal(router, t, method, route, body, xml.Unmarshal, res)
}

func TestJson(t *testing.T) {
	router := CreateRouter()

	methods := []Method{MethodGet, MethodPost, MethodPut, MethodDelete}

	for _, method := range methods {
		var res map[string]any
		readJsonBody(router, t, method, "/data/json", nil, &res)

		if !reflect.DeepEqual(res, sampleData) {
			t.Fatal(res, "does not equal", sampleData)
		}
	}

}

func testSampleUserRoutes(t *testing.T, route string, handler func(router *Router, t *testing.T, method Method, route string, body io.Reader, res any), res, expected any) {
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
		Get("error/data/global", func(request Request) Response {
			return Status(status.OK)
		}).
		Use(func(request Request, next Response) Response {
			return Error(status.InternalServerError, errors.New("errored out"))
		})
	testErrorResponse(router, t, "data/global")
}

func TestPanicErrorResponse(t *testing.T) {
	router := NewRouter().
		Get("error/data/panic", func(request Request) Response {
			panic("oh no panic")
		}).
		Use(func(request Request, next Response) Response {
			return Error(status.InternalServerError, errors.New("errored out"))
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
		Post("/test", func(request Request) Response {

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
				return Error(status.InternalServerError, errors.New("not every bind threw an error: "+strconv.Itoa(did)+"/"+strconv.Itoa(tests)))
			}

			return Status(status.OK)
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

	request := httptest.NewRequest(string(MethodGet), "/data/redirect", nil)
	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, request)

	// result := recorder.Result()

	// TODO FIX ASSERT
	// assert.Equal(t, "https://google.com", result.Headers.Get("Location"))

}

func readStringBody(router *Router, t *testing.T, method Method, route string, body io.Reader) string {

	request := httptest.NewRequest(string(method), route, body)
	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, request)

	result := recorder.Result()

	defer result.Body.Close()
	data, err := io.ReadAll(result.Body)
	if err != nil {
		t.Fatalf("could not read body %v", err)
	}
	return string(data)
}

func TestServer_Start(t *testing.T) {
	router := NewRouter()
	go func() {
		err := router.Listen(8080)
		if err.Error() != "http: Server closed" {
			t.Error("unexpected error:", err)
		}
	}()
	time.Sleep(time.Millisecond * 10)
	err := router.Server.Shutdown(context.Background())
	if err != nil {
		t.Error("unexpected error:", err)
	}
	time.Sleep(time.Millisecond * 50)
}

func TestServerTsl_Start(t *testing.T) {
	router := NewRouter()
	go func() {
		err := router.ListenToTLS(8081, "./test/server.crt", "./test2/server.key")
		if err.Error() != "open ./test/server.crt: no such file or directory" {
			t.Error("unexpected error:", err)
		}
	}()
	time.Sleep(time.Millisecond * 10)
	err := router.Server.Shutdown(context.Background())
	if err != nil {
		t.Error("unexpected error:", err)
	}
	time.Sleep(time.Millisecond * 50)
}

var (
	users = map[string]simpleUser{
		"1": {"Person"},
		"2": {"Hannes"},
	}
)

func createRouter() *Router {
	router := NewRouter()
	router.
		Use(func(request Request, next Response) Response {
			authorization := request.Headers.GetDefault(header.RequestAuthorization, "")

			user, ok := users[authorization]
			if !ok {
				return Error(status.Unauthorized, errors.New("not authorized"))
			}

			request.WithContext(context.WithValue(request.Context(), "user", user))
			return next
		})
	router.
		Get("/user", func(request Request) Response {
			user, ok := request.Context().Value("user").(simpleUser)

			if !ok {
				return Error(status.UnprocessableEntity, errors.New("could not get user from context"))
			}

			return Json(status.OK, user)
		}).With(func(request Request, next Response) Response {

		request.WithContext(context.WithValue(request.Context(), "world", "hello"))
		request.WithContext(context.WithValue(request.Context(), "hello", "world"))
		return next
	})

	router.Get("/user/test2", func(request Request) Response {
		return Status(status.OK)
	}).With(func(request Request, next Response) Response {
		request.WithContext(context.WithValue(request.Context(), "hello", "world"))
		return String(status.BadRequest, "Error")
	})

	return router
}

func TestContextMiddleware1(t *testing.T) {

	router := createRouter()

	request := httptest.NewRequest(string(MethodGet), "/user", nil)
	request.Header.Set("Authorization", "1")
	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, request)

	result := recorder.Result()

	defer result.Body.Close()
	data, err := io.ReadAll(result.Body)
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

	request := httptest.NewRequest(string(MethodGet), "/user/test2", nil)
	request.Header.Set("Authorization", "1")
	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, request)

	result := recorder.Result()

	defer result.Body.Close()
	data, err := io.ReadAll(result.Body)
	if err != nil {
		t.Fatalf("could not read body %v", err)
	}
	got := string(data)
	want := "Error"
	if got != want {
		t.Errorf("%v != %v", want, got)
	}
}

var (
	exampleReader = MapReader{
		"id":    []string{"100", "101"},
		"name":  []string{},
		"query": []string{"all"},
	}
)

func TestParamReader_Get(t *testing.T) {
	type args struct {
		key string
	}
	tests := []struct {
		name   string
		reader MapReader
		args   args
		want   string
		want1  bool
	}{
		{
			name:   "Not Existing Param",
			reader: exampleReader,
			args:   args{key: "something"},
			want:   "",
			want1:  false,
		},
		{
			name:   "Status (Not Existing) Param",
			reader: exampleReader,
			args:   args{key: "name"},
			want:   "",
			want1:  false,
		},
		{
			name:   "One Param From Polluted Param",
			reader: exampleReader,
			args:   args{key: "id"},
			want:   "100",
			want1:  true,
		},
		{
			name:   "Successful From Not Polluted",
			reader: exampleReader,
			args:   args{key: "query"},
			want:   "all",
			want1:  true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := tt.reader.Get(tt.args.key)
			if got != tt.want {
				t.Errorf("Get() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("Get() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestParamReader_GetDefault(t *testing.T) {
	type args struct {
		key          string
		defaultValue string
	}
	tests := []struct {
		name   string
		reader MapReader
		args   args
		want   string
	}{
		{
			name:   "Not Existing Param",
			reader: exampleReader,
			args:   args{key: "something", defaultValue: "abc"},
			want:   "abc",
		},
		{
			name:   "Status (Not Existing) Param",
			reader: exampleReader,
			args:   args{key: "name", defaultValue: "abc"},
			want:   "abc",
		},
		{
			name:   "One Param From Polluted Param",
			reader: exampleReader,
			args:   args{key: "id", defaultValue: "abc"},
			want:   "100",
		},
		{
			name:   "Successful From Not Polluted",
			reader: exampleReader,
			args:   args{key: "query", defaultValue: "abc"},
			want:   "all",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.reader.GetDefault(tt.args.key, tt.args.defaultValue); got != tt.want {
				t.Errorf("GetDefault() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParamReader_GetSlice(t *testing.T) {
	type args struct {
		key string
	}
	tests := []struct {
		name   string
		reader MapReader
		args   args
		want   []string
		want1  bool
	}{
		{
			name:   "Not Existing Param",
			reader: exampleReader,
			args:   args{key: "something"},
			want:   nil,
			want1:  false,
		},
		{
			name:   "Status (Not Existing) Param",
			reader: exampleReader,
			args:   args{key: "name"},
			want:   nil,
			want1:  false,
		},
		{
			name:   "One Param From Polluted Param",
			reader: exampleReader,
			args:   args{key: "id"},
			want:   []string{"100", "101"},
			want1:  true,
		},
		{
			name:   "Successful From Not Polluted",
			reader: exampleReader,
			args:   args{key: "query"},
			want:   []string{"all"},
			want1:  true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := tt.reader.GetSlice(tt.args.key)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetSlice() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("GetSlice() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestParamReader_Has(t *testing.T) {
	type args struct {
		key string
	}
	tests := []struct {
		name   string
		reader MapReader
		args   args
		want   bool
	}{
		{
			name:   "Not Existing Param",
			reader: exampleReader,
			args:   args{key: "something"},
			want:   false,
		},
		{
			name:   "Status (Not Existing) Param",
			reader: exampleReader,
			args:   args{key: "name"},
			want:   false,
		},
		{
			name:   "One Param From Polluted Param",
			reader: exampleReader,
			args:   args{key: "id"},
			want:   true,
		},
		{
			name:   "Successful From Not Polluted",
			reader: exampleReader,
			args:   args{key: "query"},
			want:   true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.reader.Has(tt.args.key); got != tt.want {
				t.Errorf("Has() = %v, want %v", got, tt.want)
			}
		})
	}
}

// Tests for router.go

func TestPort(t *testing.T) {
	tests := []struct {
		name string
		p    Port
		want string
	}{
		{name: "80", p: 80, want: ":80"},
		{name: "443", p: 443, want: ":443"},
		{name: "3000", p: 3000, want: ":3000"},
		{name: "8080", p: 8080, want: ":8080"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.p.ToAddr(); got != tt.want {
				t.Errorf("ToAddr() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRouteBuilder(t *testing.T) {
	t.Run("two middlewares", func(t *testing.T) {
		router := NewRouter()
		h := router.Get("/", handler).With(middleware).With(middleware)
		if len(h.node.middlewares[methodGet]) != 2 {
			t.Fatalf("node did not have two middlewares")
		}
	})
	t.Run("wrong parameter name", func(t *testing.T) {
		router := NewRouter()

		router.
			Patch("/user/:id/update", handler).
			Get("/user/:name/create", handler) // cant name parameter :name here, because :id was previously defined

		err := router.HasError()
		fmt.Println(err)
		if err == nil {
			t.Errorf("did not collect any error")
		}
	})
	t.Run("group prefix", func(t *testing.T) {
		router := NewRouter()
		group := router.Group("/home")
		if group.prefix != "/home/" {
			t.Error("Route prefix is not /home/ but", group.prefix)
		}
	})
}

func TestMethods(t *testing.T) {
	for _, tt := range AllMethods {
		t.Run(string(tt), func(t *testing.T) {
			r := NewRouter()
			r.Handle("/user", handler, tt)
			assertBodyResponse(t, r, string(tt), "/user", "")
		})
	}
}

func TestStatusText(t *testing.T) {
	statusText := map[int]string{
		status.Continue:           "Continue",
		status.SwitchingProtocols: "Switching Protocols",
		status.Processing:         "Processing",
		status.EarlyHints:         "Early Hints",

		status.OK:                   "OK",
		status.Created:              "Created",
		status.Accepted:             "Accepted",
		status.NonAuthoritativeInfo: "Non-Authoritative Information",
		status.NoContent:            "No Content",
		status.ResetContent:         "Reset Content",
		status.PartialContent:       "Partial Content",
		status.MultiStatus:          "Multi-Status",
		status.AlreadyReported:      "Already Reported",
		status.IMUsed:               "IM Used",

		status.MultipleChoices:   "Multiple Choices",
		status.MovedPermanently:  "Moved Permanently",
		status.Found:             "Found",
		status.SeeOther:          "See Other",
		status.NotModified:       "Not Modified",
		status.UseProxy:          "Use Proxy",
		status.TemporaryRedirect: "Temporary Redirect",
		status.PermanentRedirect: "Permanent Redirect",

		status.BadRequest:                   "Bad Request",
		status.Unauthorized:                 "Unauthorized",
		status.PaymentRequired:              "Payment Required",
		status.Forbidden:                    "Forbidden",
		status.NotFound:                     "Not Found",
		status.MethodNotAllowed:             "Method Not Allowed",
		status.NotAcceptable:                "Not Acceptable",
		status.ProxyAuthRequired:            "Proxy Authentication Required",
		status.RequestTimeout:               "Request Timeout",
		status.Conflict:                     "Conflict",
		status.Gone:                         "Gone",
		status.LengthRequired:               "Length Required",
		status.PreconditionFailed:           "Precondition Failed",
		status.RequestEntityTooLarge:        "Request Entity Too Large",
		status.RequestURITooLong:            "Request URI Too Long",
		status.UnsupportedMediaType:         "Unsupported Media Type",
		status.RequestedRangeNotSatisfiable: "Requested Range Not Satisfiable",
		status.ExpectationFailed:            "Expectation Failed",
		status.Teapot:                       "I'm a teapot",
		status.MisdirectedRequest:           "Misdirected Request",
		status.UnprocessableEntity:          "Unprocessable Entity",
		status.Locked:                       "Locked",
		status.FailedDependency:             "Failed Dependency",
		status.TooEarly:                     "Too Early",
		status.UpgradeRequired:              "Upgrade Required",
		status.PreconditionRequired:         "Precondition Required",
		status.TooManyRequests:              "Too Many Requests",
		status.RequestHeaderFieldsTooLarge:  "Request Headers Fields Too Large",
		status.UnavailableForLegalReasons:   "Unavailable For Legal Reasons",

		status.InternalServerError:           "Internal Server Error",
		status.NotImplemented:                "Not Implemented",
		status.BadGateway:                    "Bad Gateway",
		status.ServiceUnavailable:            "Service Unavailable",
		status.GatewayTimeout:                "Gateway Timeout",
		status.HTTPVersionNotSupported:       "HTTP Version Not Supported",
		status.VariantAlsoNegotiates:         "Variant Also Negotiates",
		status.InsufficientStorage:           "Insufficient Storage",
		status.LoopDetected:                  "Loop Detected",
		status.NotExtended:                   "Not Extended",
		status.NetworkAuthenticationRequired: "Network Authentication Required",
	}
	for code, text := range statusText {
		t.Run(strconv.Itoa(code)+" = "+text, func(t *testing.T) {
			if got := status.Text(code); got != text {
				t.Errorf("StatusText() = %v, want %v", got, text)
			}
		})
	}
}

func TestContentType(t *testing.T) {
	for file, value := range fileContentTypes {
		t.Run(file+" = "+value, func(t *testing.T) {
			got := ContentType(file)
			if got != value {
				t.Errorf("ContentType() = %v, want %v", got, value)
			}
		})
	}
}
