package there

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"html/template"
	"net/http"
)


//HttpResponse is the base for every return you can make in an Endpoint.
//Necessary to render the Response by calling Execute and for the WithHeaders Builder.
type HttpResponse http.Handler

type HttpResponseFunc func(http.ResponseWriter, *http.Request)

// ServeHTTP calls f(w, r).
func (f HttpResponseFunc) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	f(w, r)
}

//Bytes takes a StatusCode and a series of bytes to render
func Bytes(code int, data []byte) HttpResponse {
	return StatusWithResponse(code, &bytesResponse{data: data})
}

type bytesResponse struct {
	data []byte
}

func (j bytesResponse) ServeHTTP(rw http.ResponseWriter, r *http.Request){
	_, err := rw.Write(j.data)
	if err != nil {
		panic(err)
	}
}

// Status takes a StatusCode and renders nothing
func Status(code int) HttpResponse {
	return &statusResponse{code: code}
}

// StatusWithResponse writes the StatusCode and renders the HttpResponse
func StatusWithResponse(code int, response HttpResponse) HttpResponse {
	return &statusResponse{code: code, response: response}
}

type statusResponse struct {
	code     int
	response HttpResponse
}

func (j statusResponse) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	rw.WriteHeader(j.code)
	if j.response != nil {
		j.response.ServeHTTP(rw, r)
	}
}

// WithHeaders writes the given headers and the Http Response
func WithHeaders(headers MapString, response HttpResponse) HttpResponse {
	return &headerResponse{headers: headers, response: response}
}

type headerResponse struct {
	headers  MapString
	response HttpResponse
}

func (h headerResponse) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	for key, value := range h.headers {
		// Only set the header, if it wasn't set yet
		if len(rw.Header().Get(key)) == 0 {
			rw.Header().Set(key, value)
		}
	}
	if h.response != nil {
		h.response.ServeHTTP(rw, r)
	}
}

//String takes a StatusCode and renders the plain string
func String(code int, data string) HttpResponse {
	return Bytes(code, []byte(data))
}

//Error takes a StatusCode and err which rendering is specified by the Serializers in the RouterConfiguration
func Error(code int, err interface{}) HttpResponse {
	return Json(code, MapString{
		"error": fmt.Sprint(err),
	})
}

//Html takes a status code, the path to the html file and a map for the template parsing
func Html(code int, file string, template interface{}) HttpResponse {
	content, err := parseTemplate(file, template)
	if err != nil {
		panic(err)
	}
	return WithHeaders(MapString{
		ResponseHeaderContentType: ContentTypeTextHtml,
	}, Bytes(code, []byte(*content)))
}


func parseTemplate(templateFileName string, data interface{}) (*string, error) {
	t, err := template.ParseFiles(templateFileName)
	if err != nil {
		return nil, err
	}
	buf := new(bytes.Buffer)
	if err = t.Execute(buf, data); err != nil {
		return nil, err
	}
	body := buf.String()
	return &body, nil
}

//Json takes a StatusCode and data which gets marshaled to Json
func Json(code int, data interface{}) HttpResponse {
	jsonData, err := json.Marshal(data)
	if err != nil {
		panic(err)
	}
	return WithHeaders(MapString{
		ResponseHeaderContentType: ContentTypeApplicationJson,
	}, Bytes(code, jsonData))
}

//Message takes StatusCode and a message which will be put into a JSON object
func Message(code int, message string) HttpResponse {
	return Json(code, map[string]interface{}{
		"message": message,
	})
}

//Redirect redirects to the specific URL
func Redirect(url string) HttpResponse {
	return &redirectResponse{url: url}
}

type redirectResponse struct {
	url string
}

func (j redirectResponse) ServeHTTP(rw http.ResponseWriter, r *http.Request)  {
	http.Redirect(rw, r, j.url, StatusMovedPermanently)
}

//Xml takes a StatusCode and data which gets marshaled to Xml
func Xml(code int, data interface{}) HttpResponse {
	xmlData, err := xml.Marshal(data)
	if err != nil {
		panic(err)
	}
	return WithHeaders(MapString{
		ResponseHeaderContentType: ContentTypeApplicationXml,
	}, Bytes(code, xmlData))
}


