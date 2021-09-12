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
	http.Request
	ResponseWriter *http.ResponseWriter
}

type Response struct {
	Status  int         `json:"status,omitempty"`
	Data    interface{} `json:"message,omitempty"`
	headers map[string]string
	redirect string
}

//ResponseStatus creates a simple response only with a status code and the text from the status code
func ResponseStatus(statusCode int) Response {
	return Response{Status: statusCode, Data: http.StatusText(statusCode)}
}

//ResponseData is a response with a status code and data
func ResponseData(statusCode int, data interface{}) Response {
	return Response{Status: statusCode, Data: data}
}

//ResponseRedirect is a response with a status code and the url to redirect to
func ResponseRedirect(redirectUrl string) Response {
	return Response{Status: StatusSeeOther, Data: "Redirecting to " + redirectUrl, redirect: redirectUrl}
}


//ResponseStatusP creates a simple response pointer only with a status code and the text from the status code
func ResponseStatusP(statusCode int) *Response {
	r := ResponseStatus(statusCode)
	return &r
}

//ResponseDataP is a response pointer with a status code and data
func ResponseDataP(statusCode int, data interface{}) *Response {
	r := ResponseData(statusCode, data)
	return &r
}

//ResponseRedirectP is a response pointer with a status code and the url to redirect to
func ResponseRedirectP(redirectUrl string) *Response {
	r := ResponseRedirect(redirectUrl)
	return &r
}


//ReadBody reads the body of an http.Request as a json to the provided interface{}.
func (request *Request) ReadBody(body interface{}) error {
	err := json.NewDecoder(request.Body).Decode(body)
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			panic(err)
		}
	}(request.Body)
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


//ReadParams generates a slice of parameters in the order you parse them. Returns an error if at least one parameter is missing
func (request *Request) ReadParams(requiredParameters ...string) ([]string, error) {
	parameters := make([]string, len(requiredParameters))

	var missing []string

	for i, parameter := range requiredParameters {
		fetchedParams, ok := request.URL.Query()[parameter]
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

func (response *Response) ToJson() []byte {
	var j []byte
	switch response.Data.(type) {
	case string, *string:
		j, _ = json.Marshal(response)
	case error, *error:
		response.Data = response.Data.(error).Error()
		j, _ = json.Marshal(response)
	default:
		j, _ = json.Marshal(response.Data)
	}
	return j
}

//IsError returns true if the Data is an error or pointer to an error
func (response *Response) IsError() bool {
	switch response.Data.(type) {
	case error, *error:
		return true
	default:
		return false
	}
}