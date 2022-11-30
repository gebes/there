package there_test

/*
All benchmarks have been commented out, to avoid having the dependencies
*/
/*

import (
	"encoding/json"
	"errors"
	"github.com/Gebes/there/v2"
	"github.com/Gebes/there/v2/status"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/mux"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

var invalidError = errors.New("no authorization header provided")

func BenchmarkErrorThere(b *testing.B) {
	r := there.NewRouter()
	r.Get("/user", func(request there.Request) there.Response {
		header, ok := request.Headers.Get("authorization")
		if !ok {
			return there.Error(status.BadRequest, invalidError)
		}
		return there.Json(status.OK, map[string]string{
			"text": header,
		})
	})
	err := r.HasError()
	if err != nil {
		b.Fatalf("%v", err)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		assertBodyResponse(b, r, "GET", "/user", "{\"error\":\"no authorization header provided\"}")
	}
}

func BenchmarkErrorGin(b *testing.B) {
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()

	r.GET("/user", func(c *gin.Context) {
		header := c.GetHeader("authorization")
		if len(header) == 0 {
			c.AbortWithStatusJSON(http.StatusBadRequest, map[string]string{
				"error": invalidError.Error(),
			})
			return
		}
		c.JSON(200, gin.H{
			"text": header,
		})
	})
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		assertBodyResponse(b, r, "GET", "/user", "{\"error\":\"no authorization header provided\"}")
	}
}

func BenchmarkErrorMux(b *testing.B) {
	r := mux.NewRouter()
	r.HandleFunc("/user", func(writer http.ResponseWriter, request *http.Request) {
		header := request.Header.Get("authorization")
		if header == "" {
			j, err := json.Marshal(map[string]string{
				"error": invalidError.Error(),
			})
			if err != nil {
				b.Fatalf("could not marshal error: %v", err)
			}
			writer.WriteHeader(http.StatusBadRequest)
			_, err = writer.Write(j)
			if err != nil {
				b.Fatalf("could not write to response writer: %v", err)
			}
			return
		}
		j, err := json.Marshal(gin.H{
			"text": header,
		})
		if err != nil {
			b.Fatalf("could not marshal body: %v", err)
		}
		writer.WriteHeader(http.StatusOK)
		_, err = writer.Write(j)
		if err != nil {
			b.Fatalf("could not write to response writer: %v", err)
		}
		return
	})

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		assertBodyResponse(b, r, "GET", "/user", "{\"error\":\"no authorization header provided\"}")
	}
}

var staticRoutingPaths = []string{
	"/user",
	"/user/create",
	"/user/delete",
	"/user/list",
	"/user/all",
	"/user/student/find",
	"/user/student/query",
	"/user/student",
	"/job",
	"/job/create",
	"/job/describe",
	"/settings",
	"/settings/delete",
	"/settings/list",
	"/settings/filter",
	"/pricing/fetch",
	"/pricing/something",
}

func BenchmarkStaticRoutingThere(b *testing.B) {
	r := there.NewRouter()
	handler := func(request there.Request) there.Response {
		return there.String(status.OK, request.Request.URL.Path)
	}
	for _, path := range staticRoutingPaths {
		r.Get(path, handler)
	}
	err := r.HasError()
	if err != nil {
		b.Fatalf("%v", err)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for _, path := range staticRoutingPaths {
			assertBodyResponse(b, r, "GET", path, path)
		}
	}
}

func BenchmarkStaticRoutingGin(b *testing.B) {
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()

	handler := func(c *gin.Context) {
		c.String(200, c.Request.URL.Path)
	}
	for _, path := range staticRoutingPaths {
		r.GET(path, handler)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for _, path := range staticRoutingPaths {
			assertBodyResponse(b, r, "GET", path, path)
		}
	}
}

func BenchmarkStaticRoutingMux(b *testing.B) {
	r := mux.NewRouter()
	for _, path := range staticRoutingPaths {
		r.HandleFunc(path, func(writer http.ResponseWriter, request *http.Request) {
			writer.WriteHeader(http.StatusOK)
			_, err := writer.Write([]byte(request.URL.Path))
			if err != nil {
				b.Fatalf("could not write to response writer: %v", err)
			}
			return
		})
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for _, path := range staticRoutingPaths {
			assertBodyResponse(b, r, "GET", path, path)
		}
	}
}

var dynamicRoutingPaths = []string{
	"/user",
	"/user/:id",
	"/user/:id/create",
	"/student/:id",
	"/teacher/:id",
	"/class/:id",
	"/:id",
	"/user/delete/:id",
	"/user/list/:id",
	"/user/:id/all",
	"/user/student/:id/find",
	"/user/student/query/:id",
	"/job/:id",
	"/job/create/:id",
	"/job/describe/:id",
	"/settings/:id",
	"/settings/:id/delete",
	"/settings/list/:id",
	"/settings/filter/:id",
	"/pricing/fetch/:id",
	"/pricing/something/:id",
}

func BenchmarkDynamicRoutingThere(b *testing.B) {
	r := there.NewRouter()
	handler := func(request there.Request) there.Response {
		return there.String(status.OK, request.Request.URL.Path)
	}
	for _, path := range dynamicRoutingPaths {
		r.Get(path, handler)
	}
	err := r.HasError()
	if err != nil {
		b.Fatalf("%v", err)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for _, path := range dynamicRoutingPaths {
			assertBodyResponse(b, r, "GET", path, path)
		}
	}
}

func BenchmarkDynamicRoutingGin(b *testing.B) {
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()

	handler := func(c *gin.Context) {
		c.String(200, c.Request.URL.Path)
	}
	for _, path := range dynamicRoutingPaths {
		r.GET(path, handler)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for _, path := range dynamicRoutingPaths {
			assertBodyResponse(b, r, "GET", path, path)
		}
	}
}

func BenchmarkDynamicRoutingMux(b *testing.B) {
	r := mux.NewRouter()
	for _, path := range dynamicRoutingPaths {
		path = strings.ReplaceAll(path, ":id", "{id}")
		r.HandleFunc(path, func(writer http.ResponseWriter, request *http.Request) {
			writer.WriteHeader(http.StatusOK)
			_, err := writer.Write([]byte(request.URL.Path))
			if err != nil {
				b.Fatalf("could not write to response writer: %v", err)
			}
			return
		})
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for _, path := range dynamicRoutingPaths {
			assertBodyResponse(b, r, "GET", path, path)
		}
	}
}

func assertBodyResponse(b *testing.B, r http.Handler, method, route string, expected string) {
	request := httptest.NewRequest(method, route, nil)
	recorder := httptest.NewRecorder()
	r.ServeHTTP(recorder, request)
	result := recorder.Result()
	re, err := io.ReadAll(result.Body)
	if err != nil {
		b.Fatalf("could not read body: %v", err)
	}
	err = result.Body.Close()
	if err != nil {
		b.Fatalf("could not close body: %v", err)
	}
	if string(re) != expected {
		b.Fatalf("invalid response. actual: %v, expected: %v", string(re), expected)
	}
}

*/
