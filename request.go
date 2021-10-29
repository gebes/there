package there

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"net/http"
)

type HttpRequest struct {
	request http.Request

	Method      string
	Body        *BodyReader
	Params      *ParamReader
	RouteParams *RouteParamReader
}

func NewHttpRequest(request http.Request) HttpRequest {
	paramReader := ParamReader(request.URL.Query())
	return HttpRequest{
		request:     request,
		Method:      request.Method,
		Body:        &BodyReader{request: request},
		Params:      &paramReader,
		RouteParams: nil, // inject routeParams in Handler
	}
}

//BodyReader reads the body and unmarshal it to the specified destination
type BodyReader struct {
	request http.Request
}

func (read BodyReader) AsJson(dest interface{}) error {
	return read.format(&dest, json.Unmarshal)
}

func (read BodyReader) AsXml(dest interface{}) error {
	return read.format(&dest, xml.Unmarshal)
}

func (read BodyReader) AsYaml(dest interface{}) error {
	fmt.Println(ContentTypeApplicationEdiDashX12)
	return read.format(&dest, yaml.Unmarshal)
}

func (read BodyReader) format(dest *interface{}, formatter func(data []byte, v interface{}) error) error {
	body, err := read.AsBytes()
	if err != nil {
		return err
	}
	err = formatter(body, dest)
	return err
}

func (read BodyReader) AsString() (string, error) {
	data, err := read.AsBytes()
	if err != nil {
		return "", err
	}
	return string(data), err
}

func (read BodyReader) AsBytes() ([]byte, error) {
	data, err := ioutil.ReadAll(read.request.Body)
	defer read.request.Body.Close()
	if err != nil {
		return nil, err
	}
	return data, nil
}

//ParamReader reads http params
type ParamReader map[string][]string

func (reader ParamReader) Has(key string) bool {
	_, ok := reader.GetSlice(key)
	return ok
}

func (reader ParamReader) GetDefault(key, defaultValue string) string {
	s, ok := reader.Get(key)
	if !ok {
		return defaultValue
	}
	return s
}

func (reader ParamReader) Get(key string) (string, bool) {
	list, ok := reader.GetSlice(key)
	if !ok {
		return "", false
	}
	return list[0], true
}

func (reader ParamReader) GetSlice(key string) ([]string, bool) {
	list, ok := reader[key]
	if !ok || len(list) == 0 {
		return nil, false
	}
	return list, true
}

//RouteParamReader reads dynamic route params
type RouteParamReader map[string]string

func (reader RouteParamReader) Has(key string) bool {
	_, ok := reader[key]
	return ok
}

func (reader RouteParamReader) GetDefault(key, defaultValue string) string {
	s, ok := reader.Get(key)
	if !ok {
		return defaultValue
	}
	return s
}

func (reader RouteParamReader) Get(key string) (string, bool) {
	v, ok := reader[key]
	return v, ok
}
