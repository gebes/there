package there

import (
	"context"
	"encoding/json"
	"encoding/xml"
	"errors"
	"io"
	"io/ioutil"
	"log"
	"net/http/httptest"
	"reflect"
	"strconv"
	"testing"
	"time"
)

// constants.go tests

func TestStatusText(t *testing.T) {
	type args struct {
		code int
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "200 = OK",
			args: args{code: StatusOK},
			want: "OK",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := StatusText(tt.args.code); got != tt.want {
				t.Errorf("StatusText() = %v, want %v", got, tt.want)
			}
		})
	}
}

// Integration tests

func TestWriteError(t *testing.T) {
	router := NewRouter()
	router.Get("/", func(request Request) Response {
		// No body with 1xx status
		return Json(StatusContinue, "not writeable")
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
		return Json(StatusOK, sampleData)
	}
	// Samle User
	xml := func(request Request) Response {
		return Xml(StatusOK, sampleUser)
	}

	router := NewRouter()

	data := router.Group("/data")

	data.Handle("/json", json, MethodGet, MethodPost, MethodPut, MethodDelete)
	data.Get("/xml", xml)
	data.Get("/empty", func(request Request) Response {
		return Status(StatusAccepted)
	})
	data.Get("/message", func(request Request) Response {
		return Message(StatusOK, "Hello there")
	})
	data.Get("/string", func(request Request) Response {
		return String(StatusOK, "Hello there")
	})
	data.Get("/redirect", func(request Request) Response {
		return Redirect(StatusMovedPermanently, "https://google.com")
	})
	data.Get("/html", func(request Request) Response {
		return Html(StatusOK, "./test/index.html", map[string]string{
			"user": "Hannes",
		})
	})
	data.Get("/bytes", func(request Request) Response {
		return Bytes(StatusOK, []byte{'a', 'b'})
	})

	data.Post("/return/json", func(request Request) Response {
		var user simpleUser
		err := request.Body.BindJson(&user)
		if err != nil {
			log.Fatalln("Could not bind", err)
		}
		return String(StatusOK, user.Name)
	})
	data.Post("/return/xml", func(request Request) Response {
		var user simpleUser
		err := request.Body.BindXml(&user)
		if err != nil {
			log.Fatalln("Could not bind", err)
		}
		return String(StatusOK, user.Name)
	})
	data.Post("/return/string", func(request Request) Response {
		body, err := request.Body.ToString()
		if err != nil {
			log.Fatalln("Could not bind", err)
		}
		return String(StatusOK, body)
	})

	return router
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
		Get("error/data/global", func(request Request) Response {
			return Status(StatusOK)
		}).
		Use(func(request Request, next Response) Response {
			return Error(StatusInternalServerError, errors.New("errored out"))
		})
	testErrorResponse(router, t, "data/global")
}

func TestPanicErrorResponse(t *testing.T) {
	router := NewRouter().
		Get("error/data/panic", func(request Request) Response {
			panic("oh no panic")
		}).
		Use(func(request Request, next Response) Response {
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
				return Error(StatusInternalServerError, errors.New("not every bind threw an error: "+strconv.Itoa(did)+"/"+strconv.Itoa(tests)))
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
	// assert.Equal(t, "https://google.com", result.Headers.Get("Location"))

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
		err := router.ListenToTLS(8081, "./test/server.crt", "./test/server.key")
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
			authorization := request.Headers.GetDefault(RequestHeaderAuthorization, "")

			user, ok := users[authorization]
			if !ok {
				return Error(StatusUnauthorized, errors.New("not authorized"))
			}

			request.WithContext(context.WithValue(request.Context(), "user", user))
			return next
		})
	router.
		Get("/user", func(request Request) Response {
			user, ok := request.Context().Value("user").(simpleUser)

			if !ok {
				return Error(StatusUnprocessableEntity, errors.New("could not get user from context"))
			}

			return Json(StatusOK, user)
		}).With(func(request Request, next Response) Response {

		request.WithContext(context.WithValue(request.Context(), "world", "hello"))
		request.WithContext(context.WithValue(request.Context(), "hello", "world"))
		return next
	})

	router.Get("/user/test2", func(request Request) Response {
		return Status(StatusOK)
	}).With(func(request Request, next Response) Response {
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

// Tests for path.go

func TestConstructPath(t *testing.T) {
	type args struct {
		pathString string
		ignoreCase bool
	}
	tests := []struct {
		name string
		args args
		want Path
	}{
		{
			name: "/home",
			args: args{
				pathString: "/home",
				ignoreCase: false,
			},
			want: Path{
				parts: []pathPart{
					{value: "home", variable: false},
				},
				ignoreCase: false,
			},
		},
		{
			name: "/user/:id",
			args: args{
				pathString: "/user/:id",
				ignoreCase: false,
			},
			want: Path{
				parts: []pathPart{
					{value: "user", variable: false},
					{value: "id", variable: true},
				},
				ignoreCase: false,
			},
		},
		{
			name: "/home",
			args: args{
				pathString: "/home",
				ignoreCase: true,
			},
			want: Path{
				parts: []pathPart{
					{value: "home", variable: false},
				},
				ignoreCase: true,
			},
		},
		{
			name: "/user/:id",
			args: args{
				pathString: "/user/:id",
				ignoreCase: true,
			},
			want: Path{
				parts: []pathPart{
					{value: "user", variable: false},
					{value: "id", variable: true},
				},
				ignoreCase: true,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ConstructPath(tt.args.pathString, tt.args.ignoreCase); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ConstructPath() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestConstructPathPanic(t *testing.T) {
	defer func() { recover() }()

	//should panic because id is defined twice
	ConstructPath(":id/:id", false)
	ConstructPath(":id/:Id", false)

	t.Errorf("did not panic")
}

func TestPath_Equals(t *testing.T) {
	type args struct {
		toCompare Path
	}
	tests := []struct {
		name string
		path Path
		args args
		want bool
	}{
		{
			name: "/",
			path: ConstructPath("/", true),
			args: args{ConstructPath("/", true)},
			want: true,
		},
		{
			name: "/",
			path: ConstructPath("/", true),
			args: args{ConstructPath("/", false)},
			want: false,
		},
		{
			name: "/",
			path: ConstructPath("/", false),
			args: args{ConstructPath("/", true)},
			want: false,
		},
		{
			name: "/",
			path: ConstructPath("/", false),
			args: args{ConstructPath("/", false)},
			want: true,
		},

		{
			name: "/home/:id == /home/:uid",
			path: ConstructPath("/home/:id", false),
			args: args{ConstructPath("/home/:uid", false)},
			want: true,
		},
		{
			name: "/home/:id != /Home/:uid",
			path: ConstructPath("/home/:id", false),
			args: args{ConstructPath("/Home/:uid", false)},
			want: false,
		},

		{
			name: "/home/:id != /home/about",
			path: ConstructPath("/home/:id", false),
			args: args{ConstructPath("/home/about", false)},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.path.Equals(tt.args.toCompare); got != tt.want {
				t.Errorf("Equals() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPath_Parse(t *testing.T) {

	type args struct {
		route string
	}
	tests := []struct {
		name  string
		path  Path
		args  args
		want  map[string]string
		want1 bool
	}{
		{
			name:  "/",
			path:  ConstructPath("/", false),
			args:  args{route: "/"},
			want:  map[string]string{},
			want1: true,
		},
		{
			name: "/:id",
			path: ConstructPath("/:id", false),
			args: args{route: "/101"},
			want: map[string]string{
				"id": "101",
			},
			want1: true,
		},
		{
			name: "/user/:id",
			path: ConstructPath("/user/:id", false),
			args: args{route: "/user/101"},
			want: map[string]string{
				"id": "101",
			},
			want1: true,
		},
		{
			name: "/user/:id",
			path: ConstructPath("/user/:id", true),
			args: args{route: "/USER/101"},
			want: map[string]string{
				"id": "101",
			},
			want1: true,
		},
		{
			name: "/user/:id",
			path: ConstructPath("/user/:id", true),
			args: args{route: "/USER/101"},
			want: map[string]string{
				"id": "101",
			},
			want1: true,
		},
		{
			name: "/USER/:id",
			path: ConstructPath("/USER/:id", true),
			args: args{route: "/useR/101"},
			want: map[string]string{
				"id": "101",
			},
			want1: true,
		},
		{
			name:  "/USER/:id",
			path:  ConstructPath("/USER/:id", false),
			args:  args{route: "/useR/101"},
			want:  nil,
			want1: false,
		},
		{
			name:  "/user/:id",
			path:  ConstructPath("/user/:id", false),
			args:  args{route: "/"},
			want:  nil,
			want1: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := tt.path.Parse(tt.args.route)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Parse() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("Parse() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestPath_ToString(t *testing.T) {
	type fields struct {
		parts      []pathPart
		ignoreCase bool
	}
	type test struct {
		name   string
		fields fields
		want   string
	}
	tests := []test{}

	add := func(constructWith string, expect string) {
		tests = append(tests, test{
			name: "\"" + constructWith + "\" -> \"" + expect + "\"",
			fields: fields{
				parts:      ConstructPath(constructWith, false).parts,
				ignoreCase: false,
			},
			want: expect,
		})
	}
	add("", "/")
	add("/", "/")
	add("//////", "/")
	add("//", "/")
	add("/home", "/home")
	add("home/", "/home")
	add("/home/", "/home")
	add("home/user", "/home/user")
	add("/home/user", "/home/user")
	add("/home/user/", "/home/user")
	add("//home/user//", "/home/user")
	add("///home//user///", "/home/user")
	add("/user/:id", "/user/:id")
	add("/user/:id/", "/user/:id")
	add("user/:id/", "/user/:id")
	add("user/:id", "/user/:id")

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := Path{
				parts:      tt.fields.parts,
				ignoreCase: tt.fields.ignoreCase,
			}
			if got := p.ToString(); got != tt.want {
				t.Errorf("ToString() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_splitUrl(t *testing.T) {
	type args struct {
		route string
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{
			name: "/url",
			args: args{"/url"},
			want: []string{"url"},
		},
		{
			name: "/",
			args: args{"/"},
			want: []string{},
		},
		{
			name: "/url///home",
			args: args{"/url///home"},
			want: []string{"url", "home"},
		},
		{
			name: "url///home",
			args: args{"url///home"},
			want: []string{"url", "home"},
		},
		{
			name: "url///home/",
			args: args{"url///home/"},
			want: []string{"url", "home"},
		},
		{
			name: "url/home",
			args: args{"url/home"},
			want: []string{"url", "home"},
		},
		{
			name: "",
			args: args{""},
			want: []string{},
		},
		{
			name: "///////",
			args: args{"///////"},
			want: []string{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := splitUrl(tt.args.route); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("splitUrl() = %v, want %v", got, tt.want)
			}
		})
	}
}

// Tests for request.go

var (
	exampleRouteParamReader = RouteParamReader{
		"id":   "101",
		"name": "Max",
	}
)

func TestRouteParamReader_Get(t *testing.T) {
	type args struct {
		key string
	}
	tests := []struct {
		name   string
		reader RouteParamReader
		args   args
		want   string
		want1  bool
	}{
		{
			name:   "Existing Param",
			reader: exampleRouteParamReader,
			args:   args{key: "id"},
			want:   "101",
			want1:  true,
		},
		{
			name:   "Existing Param",
			reader: exampleRouteParamReader,
			args:   args{key: "name"},
			want:   "Max",
			want1:  true,
		},
		{
			name:   "Not Existing Param",
			reader: exampleRouteParamReader,
			args:   args{key: "something"},
			want:   "",
			want1:  false,
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

func TestRouteParamReader_GetDefault(t *testing.T) {
	type args struct {
		key          string
		defaultValue string
	}
	tests := []struct {
		name   string
		reader RouteParamReader
		args   args
		want   string
	}{
		{
			name:   "Existing Param",
			reader: exampleRouteParamReader,
			args:   args{key: "id", defaultValue: "abc"},
			want:   "101",
		},
		{
			name:   "Existing Param",
			reader: exampleRouteParamReader,
			args:   args{key: "name", defaultValue: "abc"},
			want:   "Max",
		},
		{
			name:   "Not Existing Param",
			reader: exampleRouteParamReader,
			args:   args{key: "something", defaultValue: "abc"},
			want:   "abc",
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

func TestRouteParamReader_Has(t *testing.T) {
	type args struct {
		key string
	}
	tests := []struct {
		name   string
		reader RouteParamReader
		args   args
		want   bool
	}{
		{
			name:   "Existing Param",
			reader: exampleRouteParamReader,
			args:   args{key: "id"},
			want:   true,
		},
		{
			name:   "Existing Param",
			reader: exampleRouteParamReader,
			args:   args{key: "name"},
			want:   true,
		},
		{
			name:   "Not Existing Param",
			reader: exampleRouteParamReader,
			args:   args{key: "something"},
			want:   false,
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

var (
	exampleReader = BasicReader{
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
		reader BasicReader
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
		reader BasicReader
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
		reader BasicReader
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
		reader BasicReader
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

func TestPort_ToAddr(t *testing.T) {
	tests := []struct {
		name string
		p    Port
		want string
	}{
		{
			name: "Port to string",
			p:    8080,
			want: ":8080",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.p.ToAddr(); got != tt.want {
				t.Errorf("ToAddr() = %v, want %v", got, tt.want)
			}
		})
	}
}

// Tests for routes.go

func TestNewRouteGroup(t *testing.T) {
	router := NewRouter()
	if router.RouteGroup.prefix != "/" {
		log.Fatalln("Route prefix is not /")
	}

	subGroup := router.Group("home")
	if subGroup.prefix != "/home/" {
		log.Fatalln("Route prefix is not /home/ but", subGroup.prefix)
	}

	handler := func(request Request) Response {
		return Status(StatusOK)
	}

	tests := map[string]func(route string, endpoint Endpoint) *RouteRouteGroupBuilder{
		MethodGet:     subGroup.Get,
		MethodPost:    subGroup.Post,
		MethodPatch:   subGroup.Patch,
		MethodDelete:  subGroup.Delete,
		MethodConnect: subGroup.Connect,
		MethodHead:    subGroup.Head,
		MethodTrace:   subGroup.Trace,
		MethodPut:     subGroup.Put,
		MethodOptions: subGroup.Options,
	}

	for method, methodFunc := range tests {
		if len(subGroup.routes) != 0 {
			log.Fatalln("Amount of routes should be 0")
		}

		h := methodFunc("/", handler)

		if subGroup.routes[0].Methods[0] != method {
			log.Fatalln("Method should be", method, "but is", subGroup.routes[0].Methods[0])
		}

		if len(subGroup.routes) != 1 {
			log.Fatalln("Amount of routes should be 1")
		}

		subGroup.routes.RemoveRoute(h.Route)

		if len(subGroup.routes) != 0 {
			log.Fatalln("Amount of routes should be 0")
		}
	}

}

func TestNewRouteGroup2(t *testing.T) {
	router := NewRouter()
	group := NewRouteGroup(router, "home")
	if group.prefix != "/home/" {
		log.Fatalln("Route prefix is not /home/ but", group.prefix)
	}
}

func TestNewRouteGroup3(t *testing.T) {

	handler := func(request Request) Response {
		return Status(StatusOK)
	}

	router := NewRouter()

	router.routes = nil // this should be covered and not throw any panics

	defer func() { recover() }()

	router.
		Get("/", handler).IgnoreCase().IgnoreCase().
		Get("/home/", handler).
		Get("home", handler) // duplicate route, MUST throw

	t.Errorf("did not panic")
}

func TestRouteRouteGroupBuilder_AddMiddleware2(t *testing.T) {
	handler := func(request Request) Response {
		return Status(StatusOK)
	}
	middleware := func(request Request, next Response) Response {
		return next
	}
	router := NewRouter()
	h := router.Get("/", handler).With(middleware).With(middleware)
	if len(h.Middlewares) != 2 {
		log.Fatalln("container did not have two middlewares")
	}
}

func TestRoute_OverlapsWith2(t *testing.T) {
	routeA := &Route{
		Endpoint:    nil,
		Methods:     []string{MethodGet},
		Path:        ConstructPath("/home", false),
		Middlewares: nil,
	}
	routeB := &Route{
		Endpoint:    nil,
		Methods:     []string{MethodGet, MethodPost},
		Path:        ConstructPath("/home", false),
		Middlewares: nil,
	}
	routeC := &Route{
		Endpoint:    nil,
		Methods:     []string{MethodGet, MethodPost},
		Path:        ConstructPath("/HOME", false),
		Middlewares: nil,
	}
	if !routeA.OverlapsWith(*routeB) {
		log.Fatalln("routes a and b must overlap!")
	}
	if routeA.OverlapsWith(*routeC) {
		log.Fatalln("routes a and c must NOT overlap!")
	}

}

func TestRouteGroup_Connect(t *testing.T) {
	type fields struct {
		Router *Router
		prefix string
	}
	type args struct {
		route    string
		endpoint Endpoint
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *RouteRouteGroupBuilder
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			group := &RouteGroup{
				Router: tt.fields.Router,
				prefix: tt.fields.prefix,
			}
			if got := group.Connect(tt.args.route, tt.args.endpoint); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Connect() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRouteGroup_Delete(t *testing.T) {
	type fields struct {
		Router *Router
		prefix string
	}
	type args struct {
		route    string
		endpoint Endpoint
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *RouteRouteGroupBuilder
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			group := &RouteGroup{
				Router: tt.fields.Router,
				prefix: tt.fields.prefix,
			}
			if got := group.Delete(tt.args.route, tt.args.endpoint); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Delete() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRouteGroup_Get(t *testing.T) {
	type fields struct {
		Router *Router
		prefix string
	}
	type args struct {
		route    string
		endpoint Endpoint
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *RouteRouteGroupBuilder
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			group := &RouteGroup{
				Router: tt.fields.Router,
				prefix: tt.fields.prefix,
			}
			if got := group.Get(tt.args.route, tt.args.endpoint); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Get() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRouteGroup_Group(t *testing.T) {
	type fields struct {
		Router *Router
		prefix string
	}
	type args struct {
		prefix string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *RouteGroup
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			group := RouteGroup{
				Router: tt.fields.Router,
				prefix: tt.fields.prefix,
			}
			if got := group.Group(tt.args.prefix); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Group() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRouteGroup_Handle(t *testing.T) {
	type fields struct {
		Router *Router
		prefix string
	}
	type args struct {
		path     string
		endpoint Endpoint
		methods  []string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *RouteRouteGroupBuilder
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			group := &RouteGroup{
				Router: tt.fields.Router,
				prefix: tt.fields.prefix,
			}
			if got := group.Handle(tt.args.path, tt.args.endpoint, tt.args.methods...); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Handle() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRouteGroup_Head(t *testing.T) {
	type fields struct {
		Router *Router
		prefix string
	}
	type args struct {
		route    string
		endpoint Endpoint
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *RouteRouteGroupBuilder
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			group := &RouteGroup{
				Router: tt.fields.Router,
				prefix: tt.fields.prefix,
			}
			if got := group.Head(tt.args.route, tt.args.endpoint); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Head() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRouteGroup_Options(t *testing.T) {
	type fields struct {
		Router *Router
		prefix string
	}
	type args struct {
		route    string
		endpoint Endpoint
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *RouteRouteGroupBuilder
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			group := &RouteGroup{
				Router: tt.fields.Router,
				prefix: tt.fields.prefix,
			}
			if got := group.Options(tt.args.route, tt.args.endpoint); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Options() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRouteGroup_Patch(t *testing.T) {
	type fields struct {
		Router *Router
		prefix string
	}
	type args struct {
		route    string
		endpoint Endpoint
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *RouteRouteGroupBuilder
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			group := &RouteGroup{
				Router: tt.fields.Router,
				prefix: tt.fields.prefix,
			}
			if got := group.Patch(tt.args.route, tt.args.endpoint); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Patch() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRouteGroup_Post(t *testing.T) {
	type fields struct {
		Router *Router
		prefix string
	}
	type args struct {
		route    string
		endpoint Endpoint
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *RouteRouteGroupBuilder
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			group := &RouteGroup{
				Router: tt.fields.Router,
				prefix: tt.fields.prefix,
			}
			if got := group.Post(tt.args.route, tt.args.endpoint); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Post() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRouteGroup_Put(t *testing.T) {
	type fields struct {
		Router *Router
		prefix string
	}
	type args struct {
		route    string
		endpoint Endpoint
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *RouteRouteGroupBuilder
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			group := &RouteGroup{
				Router: tt.fields.Router,
				prefix: tt.fields.prefix,
			}
			if got := group.Put(tt.args.route, tt.args.endpoint); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Put() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRouteGroup_Trace(t *testing.T) {
	type fields struct {
		Router *Router
		prefix string
	}
	type args struct {
		route    string
		endpoint Endpoint
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *RouteRouteGroupBuilder
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			group := &RouteGroup{
				Router: tt.fields.Router,
				prefix: tt.fields.prefix,
			}
			if got := group.Trace(tt.args.route, tt.args.endpoint); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Trace() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRouteManager_FindOverlappingRoute(t *testing.T) {
	type args struct {
		routeToCheck *Route
	}
	tests := []struct {
		name string
		r    routeManager
		args args
		want *Route
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.r.FindOverlappingRoute(tt.args.routeToCheck); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("FindOverlappingRoute() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRouteManager_RemoveRoute(t *testing.T) {
	type args struct {
		toRemove *Route
	}
	tests := []struct {
		name string
		r    routeManager
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
		})
	}
}

func TestRouteRouteGroupBuilder_AddMiddleware(t *testing.T) {
	type fields struct {
		Route      *Route
		RouteGroup *RouteGroup
	}
	type args struct {
		middleware Middleware
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *RouteRouteGroupBuilder
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			group := &RouteRouteGroupBuilder{
				Route:      tt.fields.Route,
				RouteGroup: tt.fields.RouteGroup,
			}
			if got := group.With(tt.args.middleware); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("With() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRouteRouteGroupBuilder_IgnoreCase(t *testing.T) {
	type fields struct {
		Route      *Route
		RouteGroup *RouteGroup
	}
	tests := []struct {
		name   string
		fields fields
		want   *RouteRouteGroupBuilder
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			group := &RouteRouteGroupBuilder{
				Route:      tt.fields.Route,
				RouteGroup: tt.fields.RouteGroup,
			}
			if got := group.IgnoreCase(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("IgnoreCase() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRoute_OverlapsWith(t *testing.T) {
	type fields struct {
		Endpoint    Endpoint
		Methods     []string
		Path        Path
		Middlewares []Middleware
	}
	type args struct {
		toCompare Route
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := Route{
				Endpoint:    tt.fields.Endpoint,
				Methods:     tt.fields.Methods,
				Path:        tt.fields.Path,
				Middlewares: tt.fields.Middlewares,
			}
			if got := e.OverlapsWith(tt.args.toCompare); got != tt.want {
				t.Errorf("OverlapsWith() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRoute_ToString(t *testing.T) {
	type fields struct {
		Endpoint    Endpoint
		Methods     []string
		Path        Path
		Middlewares []Middleware
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "ToString Uppercase",
			fields: fields{
				Endpoint: nil,
				Methods:  []string{MethodGet},
				Path: Path{
					parts: []pathPart{
						{
							value:    "",
							variable: false,
						},
					},
					ignoreCase: true,
				},
				Middlewares: nil,
			},
			want: "[GET] / *IgnoreCase",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := Route{
				Endpoint:    tt.fields.Endpoint,
				Methods:     tt.fields.Methods,
				Path:        tt.fields.Path,
				Middlewares: tt.fields.Middlewares,
			}
			if got := e.ToString(); got != tt.want {
				t.Errorf("ToString() = %v, want %v", got, tt.want)
			}
		})
	}
}

// Tests for array_utils.go

func TestCheckArrayContains(t *testing.T) {
	type args struct {
		slice    []string
		toSearch string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "a b c d includes a",
			args: args{
				slice:    []string{"a", "b", "c", "d"},
				toSearch: "a",
			},
			want: true,
		},
		{
			name: "a b c d includes b",
			args: args{
				slice:    []string{"a", "b", "c", "d"},
				toSearch: "b",
			},
			want: true,
		},
		{
			name: "a b c d includes c",
			args: args{
				slice:    []string{"a", "b", "c", "d"},
				toSearch: "c",
			},
			want: true,
		},
		{
			name: "a b c d includes d",
			args: args{
				slice:    []string{"a", "b", "c", "d"},
				toSearch: "d",
			},
			want: true,
		},
		{
			name: "a b c d does not include e",
			args: args{
				slice:    []string{"a", "b", "c", "d"},
				toSearch: "e",
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := CheckArrayContains(tt.args.slice, tt.args.toSearch); got != tt.want {
				t.Errorf("CheckArrayContains() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCheckArraysOverlap(t *testing.T) {
	type args struct {
		a []string
		b []string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "a b c overlaps with c d e",
			args: args{
				a: []string{"a", "b", "c"},
				b: []string{"c", "d", "e"},
			},
			want: true,
		},
		{
			name: "a b c overlaps with d e f",
			args: args{
				a: []string{"a", "b", "c"},
				b: []string{"d", "e", "f"},
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := CheckArraysOverlap(tt.args.a, tt.args.b); got != tt.want {
				t.Errorf("CheckArraysOverlap() = %v, want %v", got, tt.want)
			}
		})
	}
}
