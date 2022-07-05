package there_test

/* Useful for testing. Commented for dependency cleanup
import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Gebes/there/v2"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/mux"
)

func BenchmarkThere(b *testing.B) {
<<<<<<< HEAD
	router := NewRouter()
	router.Get("/", func(request HttpRequest) HttpResponse {
		return Json(StatusOK, map[string]any{
=======
	router := there.NewRouter()
	router.Get("/", func(request there.Request) there.Response {
		return there.Json(there.StatusOK, map[string]interface{}{
>>>>>>> 3792dd04d46aad11797cc77f064a12953664b24f
			"text": "Hi",
		})
	})
	for i := 0; i < b.N; i++ {
		request := httptest.NewRequest(there.MethodGet, "/", nil)
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
		writer.Headers().Set(there.ResponseHeaderContentType, there.ContentTypeApplicationJson)
		writer.WriteHeader(200)
		data, err := json.Marshal(map[string]any{
			"text": "Hi",
		})
		if err != nil {
			panic(err)
		}
		writer.Write(data)
	})
	for i := 0; i < b.N; i++ {
		request := httptest.NewRequest(there.MethodGet, "/", nil)
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
		request := httptest.NewRequest(there.MethodGet, "/", nil)
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
