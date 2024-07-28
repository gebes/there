package there

import (
	"context"
	"encoding/json"
	"encoding/xml"
	"errors"
	"io"
	"net/http"
	"net/url"
)

var (
	ErrorParameterNotPresent = errors.New("parameter not present")
)

type Request struct {
	Request        *http.Request
	ResponseWriter http.ResponseWriter

	Method        string
	Body          *BodyReader
	Params        DefaultUrlValues
	Headers       DefaultHttpHeader
	RouteParams   *RouteParamReader
	RemoteAddress string
	Host          string
	URI           string
}

func NewHttpRequest(responseWriter http.ResponseWriter, request *http.Request) Request {
	return Request{
		Request:        request,
		ResponseWriter: responseWriter,
		Method:         request.Method,
		Body:           &BodyReader{request: request},
		Params:         DefaultUrlValues(request.URL.Query()),
		Headers:        DefaultHttpHeader(request.Header),
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

type DefaultHttpHeader http.Header

func (reader DefaultHttpHeader) GetDefault(key, defaultValue string) string {
	s := http.Header(reader).Get(key)
	if s == "" {
		return defaultValue
	}
	return s
}

type DefaultUrlValues url.Values

func (reader DefaultUrlValues) GetDefault(key, defaultValue string) string {
	s := url.Values(reader).Get(key)
	if s == "" {
		return defaultValue
	}
	return s
}

type RouteParamReader struct {
	request *http.Request
}

func (reader RouteParamReader) Get(key string) string {
	return reader.request.PathValue(key)
}
