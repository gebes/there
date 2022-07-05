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
	"strings"
)

//Response is the base for every return you can make in an Endpoint.
//Necessary to render the Response by calling Execute and for the Headers Builder.
type Response http.Handler

//ResponseFunc is the type for a http.Handler
type ResponseFunc func(http.ResponseWriter, *http.Request)

// ServeHTTP calls f(w, r).
func (f ResponseFunc) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	f(w, r)
}

//Bytes takes a status code and a byte array to write to the ResponseWriter
func Bytes(code int, data []byte) Response {
	return &bytesResponse{status: code, data: data}
}

type bytesResponse struct {
	status int
	data   []byte
}

func (j bytesResponse) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	_, err := rw.Write(j.data)
	if err != nil {
		log.Printf("bytesResponse: ServeHttp write failed: %v", err)
	}
}

// Status takes a status code and renders nothing
func Status(code int) Response {
	return &statusResponse{code: code}
}

// StatusWithResponse writes the status code and renders the Response
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

// Headers adds the given map of headers to the current http.ResponseWriter and Response
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
		log.Printf("stringResponse: ServeHttp write failed: %v", err)
	}
}

//String takes a status code and renders the plain string
func String(code int, data string) Response {
	return stringResponse{code: code, data: []byte(data)}
}

//Error takes a status code and err which rendering is specified by the Serializers in the RouterConfiguration
func Error(code int, err error) Response {
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
	const (
		jsonOpen  = "{\"error\":\""
		jsonClose = "\"}"
	)
	return jsonResponse{code: code, data: []byte(jsonOpen + b.String() + jsonClose)}
}

type htmlResponse struct {
	code int
	data []byte
}

func (h htmlResponse) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	rw.WriteHeader(h.code)
	rw.Header().Set(ResponseHeaderContentType, ContentTypeTextHtml)
	_, err := rw.Write(h.data)
	if err != nil {
		log.Printf("htmlResponse: ServeHttp write failed: %v", err)
	}
}

//Html takes a status code, the path to the html file and a map for the template parsing
func Html(code int, file string, template any) Response {
	content, err := parseTemplate(file, template)
	if err != nil {
		return Error(StatusInternalServerError, fmt.Errorf("html: parseTemplate: %v", err))
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
		log.Printf("jsonResponse: ServeHttp write failed: %v", err)
	}
}

//Json takes a status code and data which gets marshaled to Json
func Json(code int, data interface{}) Response {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return Error(StatusInternalServerError, fmt.Errorf("json: json.Marshal: %v", err))
	}
	return jsonResponse{code: code, data: jsonData}
}

//Message takes status code and a message which will be put into a JSON object
func Message(code int, message string) Response {
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
		log.Printf("xmlResponse: ServeHttp write failed: %v", err)
	}
}

//Xml takes a status code and data which gets marshaled to Xml
func Xml(code int, data interface{}) Response {
	xmlData, err := xml.Marshal(data)
	if err != nil {
		return Error(StatusInternalServerError, fmt.Errorf("xml: xml.Marshal: %v", err))
	}
	return xmlResponse{code: code, data: xmlData}
}

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
		log.Printf("fileResponse: ServeHttp write failed: %v", err)
	}
}

// File returns the contents of the file provided by the path. The content type gets automatically guessed by the file extension.
// If the extension is unknown, then the fallback ContentType is ContentTypeTextPlain. Additionally, a custom ContentType can be set,
// by providing a second argument.
func File(path string, contentType ...string) Response {
	data, err := os.ReadFile(path)
	if err != nil {
		return Error(StatusNotFound, err)
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
