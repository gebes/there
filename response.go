package there

// This file contains all responses there provides by default.

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"github.com/Gebes/there/v2/header"
	"github.com/Gebes/there/v2/status"
	"html/template"
	"io"
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

// Status sets the given status code to the http.ResponseWriter.
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
//				header.ResponseAccessControlAllowOrigin:  configuration.AccessControlAllowOrigin,
//				header.ResponseAccessControlAllowMethods: configuration.AccessControlAllowMethods,
//				header.ResponseAccessControlAllowHeaders: configuration.AccessControlAllowHeaders,
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
		if rw.Header().Get(key) == "" {
			rw.Header().Set(key, value)
		}
	}
	if h.response != nil {
		h.response.ServeHTTP(rw, r)
	}
}

// Gzip wraps around your current Response and compresses all the data
// written to it, if the client has specified 'gzip' in the Accept-Encoding
// header.
func Gzip(response Response) Response {
	r := &gzipMiddleware{response}
	return r
}

type gzipMiddleware struct {
	response Response
}

type gzipResponseWriter struct {
	io.Writer
	http.ResponseWriter
}

func (w gzipResponseWriter) Write(b []byte) (int, error) {
	return w.Writer.Write(b)
}

func (j gzipMiddleware) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	if !strings.Contains(r.Header.Get(header.RequestAcceptEncoding), "gzip") {
		j.response.ServeHTTP(rw, r)
		return
	}

	gz := gzip.NewWriter(rw)
	defer gz.Close()

	rw.Header().Set(header.ContentEncoding, "gzip")

	var responseWriter http.ResponseWriter = gzipResponseWriter{Writer: gz, ResponseWriter: rw}
	j.response.ServeHTTP(responseWriter, r)
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
		switch e[i] {
		case '"':
			b.WriteString("\\\"")
		case '\'':
			b.WriteString("\\'")
		case '\v':
			b.WriteString("\\v")
		case '\f':
			b.WriteString("\\f")
		case '\r':
			b.WriteString("\\r")
		case '\n':
			b.WriteString("\\n")
		case '\t':
			b.WriteString("\\t")
		case '\b':
			b.WriteString("\\b")
		case '\a':
			b.WriteString("\\a")
		default:
			b.WriteByte(e[i])
		}
	}
	b.Write(errorJsonClose)
	return &jsonResponse{code: code, data: b.Bytes()}
}

var (
	errorJsonOpen  = []byte("{\"error\":\"")
	errorJsonClose = []byte("\"}")
)

const errorJsonLength = 14 // the total length of errorJsonOpen + errorJsonClose

// Html takes a status code, the path to the html file and a map for the template parsing
func Html(code int, file string, template any) Response {
	content, err := parseTemplate(file, template)
	if err != nil {
		return Error(status.InternalServerError, fmt.Errorf("html: parseTemplate: %v", err))
	}
	return htmlResponse{code: code, data: []byte(*content)}
}

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

func parseTemplate(templateFileName string, data any) (*string, error) {
	t, err := template.ParseFiles(templateFileName)
	if err != nil {
		return nil, err
	}
	buf := new(bytes.Buffer)
	if err := t.Execute(buf, data); err != nil {
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
		switch message[i] {
		case '"':
			b.WriteString("\\\"")
		case '\'':
			b.WriteString("\\'")
		case '\v':
			b.WriteString("\\v")
		case '\f':
			b.WriteString("\\f")
		case '\r':
			b.WriteString("\\r")
		case '\n':
			b.WriteString("\\n")
		case '\t':
			b.WriteString("\\t")
		case '\b':
			b.WriteString("\\b")
		case '\a':
			b.WriteString("\\a")
		default:
			b.WriteByte(message[i])
		}
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

// XmlError marshalls the given data parameter with the xml.Marshal function and
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

var AutoHandlers = map[string]func(code int, data any) Response{
	"fallback":                 Json,
	ContentTypeApplicationJson: Json,
	ContentTypeApplicationXml:  Xml,
}

func Auto(code int, data any) Response {
	return autoResponse{code, data}
}

type autoResponse struct {
	code int
	data any
}

func (a autoResponse) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	contentTypes := make([]string, 0, len(AutoHandlers))
	for s := range AutoHandlers {
		contentTypes = append(contentTypes, s)
	}

	contentType := NegotiateContentType(r.Header[header.RequestAccept], contentTypes, "fallback")
	handler, ok := AutoHandlers[contentType]
	if !ok {
		Error(status.BadRequest, errors.New("no suitable content-type provided")).ServeHTTP(rw, r)
	} else {
		handler(a.code, a.data).ServeHTTP(rw, r)
	}
}

// NegotiateContentType returns the best offered content type for the request's
// Accept header. If two offers match with equal weight, then the more specific
// offer is preferred.  For example, text/* trumps */*. If two offers match
// with equal weight and specificity, then the offer earlier in the list is
// preferred. If no offers match, then defaultOffer is returned.
func NegotiateContentType(headerValue, offers []string, defaultOffer string) string {
	bestOffer := defaultOffer
	bestQ := -1.0
	bestWild := 3
	specs := ParseAccept(headerValue)
	for _, offer := range offers {
		for _, spec := range specs {
			switch {
			case spec.Q == 0.0:
				// ignore
			case spec.Q < bestQ:
				// better match found
			case spec.Value == "*/*":
				if spec.Q > bestQ || bestWild > 2 {
					bestQ = spec.Q
					bestWild = 2
					bestOffer = offer
				}
			case strings.HasSuffix(spec.Value, "/*"):
				if strings.HasPrefix(offer, spec.Value[:len(spec.Value)-1]) &&
					(spec.Q > bestQ || bestWild > 1) {
					bestQ = spec.Q
					bestWild = 1
					bestOffer = offer
				}
			default:
				if spec.Value == offer &&
					(spec.Q > bestQ || bestWild > 0) {
					bestQ = spec.Q
					bestWild = 0
					bestOffer = offer
				}
			}
		}
	}
	return bestOffer
}

// AcceptSpec describes an Accept* header.
type AcceptSpec struct {
	Value string
	Q     float64
}

// ParseAccept parses Accept* headers.
func ParseAccept(header []string) (specs []AcceptSpec) {
loop:
	for _, s := range header {
		for {
			var spec AcceptSpec
			spec.Value, s = expectTokenSlash(s)
			if spec.Value == "" {
				continue loop
			}
			spec.Q = 1.0
			s = skipSpace(s)
			if strings.HasPrefix(s, ";") {
				s = skipSpace(s[1:])
				if !strings.HasPrefix(s, "q=") {
					continue loop
				}
				spec.Q, s = expectQuality(s[2:])
				if spec.Q < 0.0 {
					continue loop
				}
			}
			specs = append(specs, spec)
			s = skipSpace(s)
			if !strings.HasPrefix(s, ",") {
				continue loop
			}
			s = skipSpace(s[1:])
		}
	}
	return
}

func skipSpace(s string) (rest string) {
	i := 0
	for ; i < len(s); i++ {
		if octetTypes[s[i]]&isSpace == 0 {
			break
		}
	}
	return s[i:]
}

func expectTokenSlash(s string) (token, rest string) {
	i := 0
	for ; i < len(s); i++ {
		b := s[i]
		if (octetTypes[b]&isToken == 0) && b != '/' {
			break
		}
	}
	return s[:i], s[i:]
}

func expectQuality(s string) (quality float64, rest string) {
	switch {
	case s == "":
		return -1, ""
	case s[0] == '0':
		quality = 0
	case s[0] == '1':
		quality = 1
	default:
		return -1, ""
	}
	s = s[1:]
	if !strings.HasPrefix(s, ".") {
		return quality, s
	}
	s = s[1:]
	i := 0
	n := 0
	d := 1
	for ; i < len(s); i++ {
		b := s[i]
		if b < '0' || b > '9' {
			break
		}
		n = n*10 + int(b) - '0'
		d *= 10
	}
	return quality + float64(n)/float64(d), s[i:]
}

// Octet types from RFC 2616.
var octetTypes [256]octetType

type octetType byte

const (
	isToken octetType = 1 << iota
	isSpace
)

func init() {
	// OCTET      = <any 8-bit sequence of data>
	// CHAR       = <any US-ASCII character (octets 0 - 127)>
	// CTL        = <any US-ASCII control character (octets 0 - 31) and DEL (127)>
	// CR         = <US-ASCII CR, carriage return (13)>
	// LF         = <US-ASCII LF, linefeed (10)>
	// SP         = <US-ASCII SP, space (32)>
	// HT         = <US-ASCII HT, horizontal-tab (9)>
	// <">        = <US-ASCII double-quote mark (34)>
	// CRLF       = CR LF
	// LWS        = [CRLF] 1*( SP | HT )
	// TEXT       = <any OCTET except CTLs, but including LWS>
	// separators = "(" | ")" | "<" | ">" | "@" | "," | ";" | ":" | "\" | <">
	//              | "/" | "[" | "]" | "?" | "=" | "{" | "}" | SP | HT
	// token      = 1*<any CHAR except CTLs or separators>
	// qdtext     = <any TEXT except <">>

	for c := 0; c < 256; c++ {
		var t octetType
		isCtl := c <= 31 || c == 127
		isChar := 0 <= c && c <= 127
		isSeparator := strings.ContainsRune(" \t\"(),/:;<=>?@[]\\{}", rune(c))
		if strings.ContainsRune(" \t\r\n", rune(c)) {
			t |= isSpace
		}
		if isChar && !isCtl && !isSeparator {
			t |= isToken
		}
		octetTypes[c] = t
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
		if header == "" {
			header = ContentTypeTextPlain
		}
	}
	return fileResponse{
		code:   status.OK,
		header: header,
		data:   data,
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
