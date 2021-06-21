package there

import (
	"encoding/json"
	"errors"
	"gopkg.in/go-playground/validator.v9"
	"io"
	"net/http"
	"strings"
)

type Request struct {
	RawRequest *http.Request
	RawWriter  *http.ResponseWriter
}

type Response struct {
	Status  int         `json:"status,omitempty"`
	Data    interface{} `json:"message,omitempty"`
	headers map[string]string
}

//ResponseStatus creates a simple response only with a status code and the text from the status code
func ResponseStatus(statusCode int) Response {
	return Response{Status: statusCode, Data: http.StatusText(statusCode)}
}

//ResponseData is a response with a status code and data
func ResponseData(statusCode int, data interface{}) Response {
	return Response{Status: statusCode, Data: data}
}

//ReadBody reads the body of an http.Request as a json to the provided interface{}.
//
func (request *Request) ReadBody(body interface{}) error {
	err := json.NewDecoder(request.RawRequest.Body).Decode(body)
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			panic(err)
		}
	}(request.RawRequest.Body)
	if err != nil {
		return err
	}

	v := validator.New()
	err = v.Struct(body)
	if err != nil {
		return err
	}
	return nil
}

func (request *Request) Method() string {
	return request.RawRequest.Method
}

func (request *Request) Path() string {
	return request.RawRequest.URL.Path
}

//ReadParams generates a slice of parameters in the order you parse them. Returns an error if at least one parameter is missing
func (request *Request) ReadParams(requiredParameters ...string) ([]string, error) {
	parameters := make([]string, len(requiredParameters))

	var missing []string

	for i, parameter := range requiredParameters {
		fetchedParams, ok := request.RawRequest.URL.Query()[parameter]
		if ok && len(fetchedParams) == 1 {
			parameters[i] = fetchedParams[0]
		} else {
			missing = append(missing, parameter)
		}
	}

	if len(missing) != 0 {
		return nil, errors.New("required parameter not existing: " + strings.Join(missing, " "))
	}

	return parameters, nil

}
