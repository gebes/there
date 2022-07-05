package there

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

//Response is the base for every return you can make in an Endpoint.
//Necessary to render the Response by calling Execute and for the Headers Builder.
type Response http.Handler

type ResponseFunc func(http.ResponseWriter, *http.Request)

// ServeHTTP calls f(w, r).
func (f ResponseFunc) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	f(w, r)
}

//Bytes takes a StatusCode and a series of bytes to render
func Bytes(code int, data []byte) Response {
	return StatusWithResponse(code, &bytesResponse{data: data})
}

type bytesResponse struct {
	data []byte
}

func (j bytesResponse) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	_, err := rw.Write(j.data)
	if err != nil {
		log.Printf("bytesResponse: write failed: %v", err)
	}
}

// Status takes a StatusCode and renders nothing
func Status(code int) Response {
	return &statusResponse{code: code}
}

// StatusWithResponse writes the StatusCode and renders the Response
func StatusWithResponse(code int, response Response) Response {
	return &statusResponse{code: code, response: response}
}

type statusResponse struct {
	code     int
	response Response
}

func (j statusResponse) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	rw.WriteHeader(j.code)
	if j.response != nil {
		j.response.ServeHTTP(rw, r)
	}
}

// Headers writes the given headers and the Http Response
func Headers(headers map[string]string, response Response) Response {
	return &headerResponse{headers: headers, response: response}
}

type headerResponse struct {
	headers  map[string]string
	response Response
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

type stringResponse struct {
	code int
	data []byte
}

func (s stringResponse) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	rw.WriteHeader(s.code)
	rw.Header().Set(ResponseHeaderContentType, ContentTypeTextPlain)
	_, err := rw.Write(s.data)
	if err != nil {
		log.Printf("stringResponse: write failed: %v", err)
	}
}

//String takes a StatusCode and renders the plain string
func String(code int, data string) Response {
	return stringResponse{code: code, data: []byte(data)}
}

//Error takes a StatusCode and err which rendering is specified by the Serializers in the RouterConfiguration
func Error(code int, err error) HttpResponse {
	e := err.Error()
	var b strings.Builder
	b.Grow(len(e))
	for i := range e {
		if e[i] == '"' {
			b.WriteString("\"")
			continue
		}
		b.WriteByte(e[i])
	}
	return jsonResponse{code: code, data: []byte(jsonLeft + b.String() + jsonRight)}
}

type htmlResponse struct {
	code int
	data []byte
}

func (h htmlResponse) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	rw.WriteHeader(h.code)
	rw.Header().Set(ResponseHeaderContentType, ContentTypeTextHtml)
	rw.Write(h.data)
}

//Html takes a status code, the path to the html file and a map for the template parsing
func Html(code int, file string, template any) HttpResponse {
	content, err := parseTemplate(file, template)
	if err != nil {
		panic(err)
	}
	return htmlResponse{code: code, data: []byte(*content)}
}

func parseTemplate(templateFileName string, data any) (*string, error) {
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

type jsonResponse struct {
	code int
	data []byte
}

func (j jsonResponse) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	rw.WriteHeader(j.code)
	rw.Header().Set(ResponseHeaderContentType, ContentTypeApplicationJson)
	_, err := rw.Write(j.data)
	if err != nil {
		panic(err)
	}
}

//Json takes a StatusCode and data which gets marshaled to Json
func Json(code int, data interface{}) Response {
	jsonData, err := json.Marshal(data)
	if err != nil {
		panic(err)
	}
	return jsonResponse{code: code, data: jsonData}
}

//Message takes StatusCode and a message which will be put into a JSON object
func Message(code int, message string) HttpResponse {
	return Json(code, map[string]any{
		"message": message,
	})
}

//Redirect redirects to the specific URL
func Redirect(code int, url string) Response {
	return &redirectResponse{code: code, url: url}
}

type redirectResponse struct {
	code int
	url  string
}

func (j redirectResponse) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	http.Redirect(rw, r, j.url, j.code)
}

type xmlResponse struct {
	code int
	data []byte
}

func (x xmlResponse) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	rw.Header().Set(ResponseHeaderContentType, ContentTypeApplicationXml)
	rw.WriteHeader(x.code)
	_, err := rw.Write(x.data)
	if err != nil {
		panic(err)
	}
}

//Xml takes a StatusCode and data which gets marshaled to Xml
func Xml(code int, data interface{}) Response {
	xmlData, err := xml.Marshal(data)
	if err != nil {
		panic(err)
	}
	return xmlResponse{code: code, data: xmlData}
}

// File takes the path to a file-serving, and sets the response equal to the bytes of it.
// It also selects an appropriate content type header, depending on the file-serving extension.
// Additionally, a fallbackContentType can be passed, if the content type corresponding to
// the file-serving extension was not found.

type fileResponse struct {
	code   int
	header string
	data   []byte
}

func (f fileResponse) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	rw.Header().Set(ResponseHeaderContentType, f.header)
	rw.WriteHeader(f.code)
	_, err := rw.Write(f.data)
	if err != nil {
		panic(err)
	}
}

func File(path string, contentType ...string) Response {
	data, err := os.ReadFile(path)
	if err != nil {
		return Error(StatusNotFound, err.Error())
	}
	var header string
	if len(contentType) >= 1 {
		header = contentType[0]
	} else {
		extension := filepath.Ext(path)
		if extension != "" {
			extension = extension[1:]
		}

		header = ContentType(extension)
		if len(header) == 0 {
			header = ContentTypeTextPlain
		}
	}
	return fileResponse{
		code:   StatusOK,
		header: header,
		data:   data,
	}
}

// File takes the path to a file-serving, and sets the response equal to the bytes of it.
// It also selects an appropriate content type header, depending on the file-serving extension.
// Additionally, a fallbackContentType can be passed, if the content type corresponding to
// the file-serving extension was not found.

type fileResponse struct {
	code   int
	header string
	data   []byte
}

func (f fileResponse) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	rw.Header().Set(ResponseHeaderContentType, f.header)
	rw.WriteHeader(f.code)
	_, err := rw.Write(f.data)
	if err != nil {
		panic(err)
	}
}

func File(path string, contentType ...string) Response {
	data, err := os.ReadFile(path)
	if err != nil {
		return Error(StatusNotFound, err.Error())
	}
	var header string
	if len(contentType) >= 1 {
		header = contentType[0]
	} else {
		extension := filepath.Ext(path)
		if extension != "" {
			extension = extension[1:]
		}

		header = ContentType(extension)
		if len(header) == 0 {
			header = ContentTypeTextPlain
		}
	}
	return fileResponse{
		code:   StatusOK,
		header: header,
		data:   data,
	}
}
