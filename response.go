package there

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"html/template"
	"net/http"
	"github.com/vmihailenco/msgpack/v5"
)

//HttpResponse is the base for every return you can make in an Endpoint.
//Necessary to render the Response by calling Execute and for the Header Builder.
type HttpResponse interface {
	Header() *HeaderWrapper
	Execute(router *Router, r *http.Request, w *http.ResponseWriter) error
}

//HeaderWrapper for the fluent Header Builder
type HeaderWrapper struct {
	Values map[string][]string
	HttpResponse
}

func Header(httpResponse HttpResponse) *HeaderWrapper {
	return &HeaderWrapper{Values: map[string][]string{}, HttpResponse: httpResponse}
}

func (h *HeaderWrapper) Set(key string, values ...string) *HeaderWrapper {
	h.Values[key] = values
	return h
}

func (h *HeaderWrapper) SetAll(values map[string][]string) *HeaderWrapper {
	for key, values := range values {
		h.Set(key, values...)
	}
	return h
}

//Bytes takes a StatusCode and a series of bytes to render
func Bytes(code int, data []byte) *bytesResponse {
	r := &bytesResponse{data: data, code: code}
	r.header = Header(r)
	return r
}

type bytesResponse struct {
	code   int
	data   []byte
	header *HeaderWrapper
}

func (j bytesResponse) Execute(_ *Router, _ *http.Request, w *http.ResponseWriter) error {
	(*w).WriteHeader(j.code)
	_, err := (*w).Write(j.data)
	return err
}

func (j *bytesResponse) Header() *HeaderWrapper {
	return j.header
}

//Empty takes a StatusCode and renders nothing
func Empty(code int) *emptyResponse {
	r := &emptyResponse{code: code}
	r.header = Header(r)
	return r
}

type emptyResponse struct {
	code   int
	header *HeaderWrapper
}

func (j emptyResponse) Execute(_ *Router, _ *http.Request, w *http.ResponseWriter) error {
	(*w).WriteHeader(j.code)
	_, err := (*w).Write([]byte(""))
	return err
}

func (j *emptyResponse) Header() *HeaderWrapper {
	return j.header
}

//String takes a StatusCode and renders the plain string
func String(code int, data string) *stringResponse {
	r := &stringResponse{data: data, code: code}
	r.header = Header(r).Set(ResponseHeaderContentType, ContentTypeTextPlain)
	return r
}

type stringResponse struct {
	code   int
	data   string
	header *HeaderWrapper
}

func (j stringResponse) Execute(_ *Router, _ *http.Request, w *http.ResponseWriter) error {
	(*w).WriteHeader(j.code)
	_, err := (*w).Write([]byte(j.data))
	if err != nil {
		return err
	}
	return nil
}

func (j stringResponse) Header() *HeaderWrapper {
	return j.header
}

//Error takes a StatusCode and err which rendering is specified by the Serializers in the Configuration
func Error(code int, err error) *errorResponse {
	r := &errorResponse{err: err, code: code}
	r.header = Header(r)
	return r
}

type errorResponse struct {
	code   int
	err    error
	header *HeaderWrapper
}

func (j errorResponse) Execute(router *Router, _ *http.Request, w *http.ResponseWriter) error {
	data, err := router.Configuration.ErrorToBytes(j.err)
	if err != nil {
		return err
	}
	(*w).Header().Set(ResponseHeaderContentType, router.Configuration.ErrorToBytesContentType())
	(*w).WriteHeader(j.code)
	_, err = (*w).Write(data)
	return err
}
func (j *errorResponse) Header() *HeaderWrapper {
	return j.header
}

//Html takes a status code, the path to the html file and a map for the template parsing
func Html(code int, file string, template map[string]string) *htmlResponse {
	r := &htmlResponse{file: file, template: template, code: code}
	r.header = Header(r).Set(ResponseHeaderContentType, ContentTypeTextHtml)
	return r
}

type htmlResponse struct {
	file     string
	template map[string]string
	code     int
	header   *HeaderWrapper
}

func (j *htmlResponse) Execute(_ *Router, _ *http.Request, w *http.ResponseWriter) error {
	content, err := parseTemplate(j.file, j.template)
	if err != nil {
		return err
	}
	(*w).WriteHeader(j.code)
	_, err = (*w).Write([]byte(*content))
	return err
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

func (j *htmlResponse) Header() *HeaderWrapper {
	return j.header
}

//Json takes a StatusCode and data which gets marshaled to Json
func Json(code int, data interface{}) *jsonResponse {
	r := &jsonResponse{data: data, code: code}
	r.header = Header(r).Set(ResponseHeaderContentType, ContentTypeApplicationJson)
	return r
}

type jsonResponse struct {
	data   interface{}
	code   int
	header *HeaderWrapper
}

func (j *jsonResponse) Execute(_ *Router, _ *http.Request, w *http.ResponseWriter) error {
	b, err := json.Marshal(j.data)
	if err != nil {
		return err
	}
	(*w).WriteHeader(j.code)
	_, err = (*w).Write(b)
	return err
}

func (j *jsonResponse) Header() *HeaderWrapper {
	return j.header
}

//Message takes StatusCode and a message which will be put into a JSON object
func Message(code int, message string) *jsonResponse {
	return Json(code, map[string]interface{}{
		"message": message,
	})
}

//Msgpack takes a StatusCode and data which gets marshaled to Msgpack
func Msgpack(code int, data interface{}) *msgpackResponse {
	r := &msgpackResponse{data: data, code: code}
	r.header = Header(r).Set(ResponseHeaderContentType, "application/msgpack")
	return r
}

type msgpackResponse struct {
	data   interface{}
	code   int
	header *HeaderWrapper
}

func (j *msgpackResponse) Execute(_ *Router, _ *http.Request, w *http.ResponseWriter) error {
	b, err := msgpack.Marshal(j.data)
	if err != nil {
		return err
	}
	(*w).WriteHeader(j.code)
	_, err = (*w).Write(b)
	return err
}

func (j *msgpackResponse) Header() *HeaderWrapper {
	return j.header
}

//Redirect redirects to the specific URL
func Redirect(url string) *redirectResponse {
	r := &redirectResponse{url: url}
	r.header = Header(r)
	return r
}

type redirectResponse struct {
	url    string
	header *HeaderWrapper
}

func (j redirectResponse) Execute(_ *Router, r *http.Request, w *http.ResponseWriter) error {
	(*w).WriteHeader(StatusMovedPermanently)
	http.Redirect(*w, r, j.url, StatusMovedPermanently)
	return nil
}

func (j *redirectResponse) Header() *HeaderWrapper {
	return j.header
}

//Xml takes a StatusCode and data which gets marshaled to Xml
func Xml(code int, data interface{}) *xmlResponse {
	r := &xmlResponse{data: data, code: code}
	r.header = Header(r).Set(ResponseHeaderContentType, ContentTypeApplicationXml)
	return r
}

type xmlResponse struct {
	data   interface{}
	code   int
	header *HeaderWrapper
}

func (j *xmlResponse) Execute(_ *Router, _ *http.Request, w *http.ResponseWriter) error {
	b, err := xml.Marshal(j.data)
	if err != nil {
		return err
	}
	(*w).WriteHeader(j.code)
	_, err = (*w).Write(b)
	return err
}

func (j *xmlResponse) Header() *HeaderWrapper {
	return j.header
}

//Yaml takes a StatusCode and data which gets marshaled to Yaml
func Yaml(code int, data interface{}) *yamlResponse {
	r := &yamlResponse{data: data, code: code}
	r.header = Header(r).Set(ResponseHeaderContentType, "text/x-yaml")
	return r
}

type yamlResponse struct {
	data   interface{}
	code   int
	header *HeaderWrapper
}

func (j *yamlResponse) Execute(_ *Router, _ *http.Request, w *http.ResponseWriter) error {
	b, err := xml.Marshal(j.data)
	if err != nil {
		return err
	}
	(*w).WriteHeader(j.code)
	_, err = (*w).Write(b)
	return err
}

func (j *yamlResponse) Header() *HeaderWrapper {
	return j.header
}
