package there

import (
	"context"
	"encoding/json"
	"encoding/xml"
	"errors"
	"io"
	"net/http"
)

var (
	ErrorParameterNotPresent = errors.New("parameter not present")
)

type Request struct {
	Request        *http.Request
	ResponseWriter http.ResponseWriter

	Method        string
	Body          *BodyReader
	Params        *MapReader
	Headers       *MapReader
	RouteParams   *RouteParamReader
	RemoteAddress string
	Host          string
	URI           string
}

func NewHttpRequest(responseWriter http.ResponseWriter, request *http.Request) Request {
	paramReader := MapReader(request.URL.Query())
	headerReader := MapReader(request.Header)
	return Request{
		Request:        request,
		ResponseWriter: responseWriter,
		Method:         request.Method,
		Body:           &BodyReader{request: request},
		Params:         &paramReader,
		Headers:        &headerReader,
		RouteParams:    &RouteParamReader{request},
		RemoteAddress:  request.RemoteAddr,
		URI:            request.RequestURI,
	}
}

func (r *Request) Context() context.Context {
	return r.Request.Context()
}

func (r *Request) WithContext(ctx context.Context) {
	*r.Request = *r.Request.WithContext(ctx)
}

// BodyReader reads the body and unmarshal it to the specified destination
type BodyReader struct {
	request *http.Request
}

func (read BodyReader) BindJson(dest any) error {
	return read.bind(dest, json.Unmarshal)
}

func (read BodyReader) BindXml(dest any) error {
	return read.bind(dest, xml.Unmarshal)
}

func (read BodyReader) bind(dest any, formatter func(data []byte, v any) error) error {
	body, err := read.ToBytes()
	if err != nil {
		return err
	}
	err = formatter(body, dest)
	return err
}

func (read BodyReader) ToString() (string, error) {
	data, err := read.ToBytes()
	if err != nil {
		return "", err
	}
	return string(data), err
}

func (read BodyReader) ToBytes() ([]byte, error) {
	data, err := io.ReadAll(read.request.Body)
	defer read.request.Body.Close()
	if err != nil {
		return nil, err
	}
	return data, nil
}

// MapReader reads http params
type MapReader map[string][]string

func (reader MapReader) Has(key string) bool {
	_, ok := reader.GetSlice(key)
	return ok
}

func (reader MapReader) GetDefault(key, defaultValue string) string {
	s, ok := reader.Get(key)
	if !ok {
		return defaultValue
	}
	return s
}

func (reader MapReader) Get(key string) (string, bool) {
	list, ok := reader.GetSlice(key)
	if !ok {
		return "", false
	}
	return list[0], true
}

func (reader MapReader) GetSlice(key string) ([]string, bool) {
	list, ok := reader[key]
	if !ok || len(list) == 0 {
		return nil, false
	}
	return list, true
}

type RouteParamReader struct {
	request *http.Request
}

func (reader RouteParamReader) Get(key string) string {
	return reader.request.PathValue(key)
}
