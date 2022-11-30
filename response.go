// This file contains all responses there provides by default.
package there

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"github.com/Gebes/there/v2/header"
	"github.com/Gebes/there/v2/status"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

// Response is the base for every return you can make in an Endpoint.
// Necessary to render the Response by calling Execute and for the Headers Builder.
type Response http.Handler

// ResponseFunc is the type for a http.Handler
type ResponseFunc func(http.ResponseWriter, *http.Request)

// ServeHTTP calls f(w, r).
func (f ResponseFunc) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	f(w, r)
}

// Bytes writes the data parameter with the given status code
// to the http.ResponseWriter
//
// The Content-Type header is set to nothing at all.
//
//	func ExampleStringGet(request there.Request) there.Response {
//		return there.String(status.OK, "Hello there")
//	}
//
// When this handler gets called, the final rendered result will be
//
//	Hello there
func Bytes(code int, data []byte) Response {
	return &bytesResponse{code: code, data: data}
}

type bytesResponse struct {
	code int
	data []byte
}

func (j bytesResponse) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	_, err := rw.Write(j.data)
	if err != nil {
		log.Printf("bytesResponse: ServeHttp write failed: %v", err)
	}
}

// String sets the given status code to the http.ResponseWriter.
//
// No Content-Type header is set at all.
//
//	func ExampleStatusGet(request there.Request) there.Response {
//		return there.Status(status.OK)
//	}
//
// When this handler gets called, the final rendered result will be
// an empty body with the status 200
func Status(code int) Response {
	return &statusResponse{code: code}
}

type statusResponse struct {
	code int
}

func (j statusResponse) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	rw.WriteHeader(j.code)
}

// Headers wraps around your current Response and sets all the headers
// parsed in the headers parameter provided they are not set. If a header
// was already set by a previous Response, then it will be skipped.
//
//	func Cors(configuration CorsConfiguration) there.Middleware {
//		return func(request there.Request, next there.Response) there.Response {
//			headers := map[string]string{
//				there.ResponseHeaderAccessControlAllowOrigin:  configuration.AccessControlAllowOrigin,
//				there.ResponseHeaderAccessControlAllowMethods: configuration.AccessControlAllowMethods,
//				there.ResponseHeaderAccessControlAllowHeaders: configuration.AccessControlAllowHeaders,
//			}
//			if request.Method == there.MethodOptions {
//				return there.Headers(headers, there.Status(status.OK))
//			}
//			return there.Headers(headers, next)
//		}
//	}
//
// When this middleware gets called, all the Cors Headers will be set.
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

// String writes the data parameter with the given status code
// to the http.ResponseWriter
//
// The Content-Type header is set accordingly to text/plain
//
//	func ExampleStringGet(request there.Request) there.Response {
//		return there.String(status.OK, "Hello there")
//	}
//
// When this handler gets called, the final rendered result will be
//
//	Hello there
func String(code int, data string) Response {
	return stringResponse{code: code, data: []byte(data)}
}

type stringResponse struct {
	code int
	data []byte
}

func (s stringResponse) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	rw.WriteHeader(s.code)
	rw.Header().Set(header.ContentType, ContentTypeTextPlain)
	_, err := rw.Write(s.data)
	if err != nil {
		log.Printf("stringResponse: ServeHttp write failed: %v", err)
	}
}

// Error builds a json response using the err parameter and writes
// the result with the given status code to the http.ResponseWriter
//
// The Content-Type header is set accordingly to application/json
//
//	func ExampleErrorGet(request there.Request) there.Response {
//		if 1 != 2 {
//			return there.Error(status.InternalServerError, errors.New("something went wrong"))
//		}
//		return there.Status(status.OK)
//	}
//
// When this handler gets called, the final rendered result will be
//
//	{"error":"something went wrong"}
//
// For optimal performance the use of json.Marshal is avoided and the response
// body is built directly. With this way, there is no error that could occur.
func Error(code int, err error) Response {
	e := err.Error()
	var b bytes.Buffer
	b.Grow(len(e) + errorJsonLength)
	b.Write(errorJsonOpen)
	for i := range e {
		if e[i] == '"' {
			b.WriteString("\\\"")
			continue
		}
		b.WriteByte(e[i])
	}
	b.Write(errorJsonClose)
	return &jsonResponse{code: code, data: b.Bytes()}
}

var (
	errorJsonOpen  = []byte("{\"error\":\"")
	errorJsonClose = []byte("\"}")
)

const errorJsonLength = 14 // the total length of errorJsonOpen + errorJsonClose

type htmlResponse struct {
	code int
	data []byte
}

func (h htmlResponse) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	rw.WriteHeader(h.code)
	rw.Header().Set(header.ContentType, ContentTypeTextHtml)
	_, err := rw.Write(h.data)
	if err != nil {
		log.Printf("htmlResponse: ServeHttp write failed: %v", err)
	}
}

// Html takes a status code, the path to the html file and a map for the template parsing
func Html(code int, file string, template any) Response {
	content, err := parseTemplate(file, template)
	if err != nil {
		return Error(status.InternalServerError, fmt.Errorf("html: parseTemplate: %v", err))
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

// Json marshalls the given data parameter with the json.Marshal function
// and writes the result with the given status code to the http.ResponseWriter
//
// The Content-Type header is set accordingly to application/json
//
//	func ExampleGet(request there.Request) there.Response {
//		user := map[string]string{
//			"firstname": "John",
//			"surname": "Smith",
//		}
//		return there.Json(status.OK, user)
//	}
//
// When this handler gets called, the final rendered result will be
//
//	{"firstname":"John","surname":"Smith"}
//
// If the json.Marshal fails with an error, then an Error with StatusInternalServerError will be returned, with the error format "json: json.Marshal: %v"
func Json(code int, data any) Response {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return Error(status.InternalServerError, fmt.Errorf("json: json.Marshal: %v", err))
	}
	return jsonResponse{code: code, data: jsonData}
}

// JsonError marshalls the given data parameter with the json.Marshal function
// and writes the result with the given status code to the http.ResponseWriter
//
// The Content-Type header is set accordingly to application/json
//
//	func ExampleJsonErrorGet(request there.Request) there.Response{
//		user := map[string]string{
//			"firstname": "John",
//			"surname": "Smith",
//		}
//		resp, err := there.JsonError(status.OK, user)
//		if err != nil {
//			return there.Error(status.InternalServerError, fmt.Errorf("something went wrong: %v", err))
//		}
//		return resp
//	}
//
// When this handler gets called, the final rendered result will be
//
//	{"firstname":"John","surname":"Smith"}
//
// If the json.Marshal fails with an error, then a nil response with a non-nil error will be returned to handle.
func JsonError(code int, data any) (Response, error) {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}
	return jsonResponse{code: code, data: jsonData}, nil
}

type jsonResponse struct {
	code int
	data []byte
}

func (j jsonResponse) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	rw.WriteHeader(j.code)
	rw.Header().Set(header.ContentType, ContentTypeApplicationJson)
	_, err := rw.Write(j.data)
	if err != nil {
		log.Printf("jsonResponse: ServeHttp write failed: %v", err)
	}
}

// Message builds a json response using the message parameter and writes
// the result with the given status code to the http.ResponseWriter
//
// The Content-Type header is set accordingly to application/json
//
//	func ExampleMessageGet(request there.Request) there.Response {
//		return there.Message(status.OK, "Hello there")
//	}
//
// When this handler gets called, the final rendered result will be
//
//	{"message":"Hello there"}
//
// For optimal performance the use of json.Marshal is avoided and the response
// body is built directly. With this way, there is no error that could occur.
func Message(code int, message string) Response {
	var b strings.Builder
	b.Grow(len(message))
	for i := range message {
		if message[i] == '"' {
			b.WriteString("\\\"")
			continue
		}
		b.WriteByte(message[i])
	}
	const (
		jsonOpen  = "{\"message\":\""
		jsonClose = "\"}"
	)
	return jsonResponse{code: code, data: []byte(jsonOpen + b.String() + jsonClose)}
}

// Redirect redirects to the specific URL
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

// Xml marshalls the given data parameter with the xml.Marshal function and
// writes the result with the given status code to the http.ResponseWriter
//
// The Content-Type header is set accordingly to application/xml
//
//	type User struct {
//		Firstname string `xml:"firstname"`
//		Surname string  `xml:"surname"`
//	}
//
//	func ExampleXmlGet(request there.Request) there.Response{
//		user := User{"John", "Smith"}
//		return there.Xml(status.OK, user)
//	}
//
// When this handler gets called, the final rendered result will be
//
//	<User><firstname>John</firstname><surname>Smith</surname></User>
//
// If the xml.Marshal fails with an error, then an Error with StatusInternalServerError will be returned, with the error format "xml: xml.Marshal: %v"
func Xml(code int, data any) Response {
	xmlData, err := xml.Marshal(data)
	if err != nil {
		return Error(status.InternalServerError, fmt.Errorf("xml: xml.Marshal: %v", err))
	}
	return xmlResponse{code: code, data: xmlData}
}

// Xml marshalls the given data parameter with the xml.Marshal function and
// writes the result with the given status code to the http.ResponseWriter
//
// The Content-Type header is set accordingly to application/xml
//
//	type User struct {
//		Firstname string `xml:"firstname"`
//		Surname   string `xml:"surname"`
//	}
//
//	func ExampleXmlErrorGet(request there.Request) there.Response {
//		user := User{"John", "Smith"}
//		resp, err := there.XmlError(status.OK, user)
//		if err != nil {
//			return there.Error(status.InternalServerError, fmt.Errorf("something went wrong: %v", err))
//		}
//		return resp
//	}
//
// When this handler gets called, the final rendered result will be
//
//	<User><firstname>John</firstname><surname>Smith</surname></User>
//
// If the xml.Marshal fails with an error, then a nil response with a non-nil error will be returned to handle.
func XmlError(code int, data any) (Response, error) {
	xmlData, err := xml.Marshal(data)
	if err != nil {
		return nil, err
	}
	return xmlResponse{code: code, data: xmlData}, nil
}

type xmlResponse struct {
	code int
	data []byte
}

func (x xmlResponse) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	rw.Header().Set(header.ContentType, ContentTypeApplicationXml)
	rw.WriteHeader(x.code)
	_, err := rw.Write(x.data)
	if err != nil {
		log.Printf("xmlResponse: ServeHttp write failed: %v", err)
	}
}

type fileResponse struct {
	code   int
	header string
	data   []byte
}

func (f fileResponse) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	rw.Header().Set(header.ContentType, f.header)
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
		return Error(status.NotFound, err)
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
		code:   status.OK,
		header: header,
		data:   data,
	}
}
