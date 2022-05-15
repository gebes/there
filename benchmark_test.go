package there_test

/* Useful for testing. Commented for dependency cleanup
import (
	"encoding/json"
	. "github.com/Gebes/there/v2"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/mux"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func BenchmarkThere(b *testing.B) {
	router := NewRouter()
	router.Get("/", func(request Request) Response {
		return Json(StatusOK, map[string]interface{}{
			"text": "Hi",
		})
	})
	for i := 0; i < b.N; i++ {
		request := httptest.NewRequest(MethodGet, "/", nil)
		recorder := httptest.NewRecorder()

		router.ServeHTTP(recorder, request)

		result := recorder.Result()

		_, err := ioutil.ReadAll(result.Body)
		if err != nil {
			b.Fatalf("could not read body %v", err)
		}
		result.Body.Close()
	}
}

func BenchmarkMux(b *testing.B) {
	r := mux.NewRouter()
	r.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		writer.Headers().Set(ResponseHeaderContentType, ContentTypeApplicationJson)
		writer.WriteHeader(200)
		data, err := json.Marshal(map[string]interface{}{
			"text": "Hi",
		})
		if err != nil {
			panic(err)
		}
		writer.Write(data)
	})
	for i := 0; i < b.N; i++ {
		request := httptest.NewRequest(MethodGet, "/", nil)
		recorder := httptest.NewRecorder()

		r.ServeHTTP(recorder, request)

		result := recorder.Result()

		_, err := ioutil.ReadAll(result.Body)
		if err != nil {
			b.Fatalf("could not read body %v", err)
		}
		result.Body.Close()
	}
}

func BenchmarkGin(b *testing.B) {
	r := gin.New()
	gin.SetMode(gin.ReleaseMode)
	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"text": "Hi",
		})
	})

	for i := 0; i < b.N; i++ {
		request := httptest.NewRequest(MethodGet, "/", nil)
		recorder := httptest.NewRecorder()

		r.ServeHTTP(recorder, request)

		result := recorder.Result()

		_, err := ioutil.ReadAll(result.Body)
		if err != nil {
			b.Fatalf("could not read body %v", err)
		}
		result.Body.Close()
	}
}
*/
