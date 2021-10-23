package there

import (
	"encoding/json"
	. "github.com/Gebes/there/there/utils"
	"net/http"
)

type ErrorResponse struct {
	code   int
	err    error
	header *HeaderWrapper
}

func Error(code int, err error) *ErrorResponse {
	r := &ErrorResponse{err: err, code: code}
	r.header = Header(r).Set(ResponseHeaderContentType, ContentTypeApplicationJson)
	return r
}

var ErrorSerializer = func(err error) ([]byte, error) {
	data, err := json.Marshal(map[string]interface{}{
		"error": err.Error(),
	})
	return data, err
}

func (j ErrorResponse) Execute(r *http.Request, w *http.ResponseWriter) error {
	data, err := ErrorSerializer(j.err)
	if err != nil {
		return err
	}
	(*w).WriteHeader(j.code)
	_, err = (*w).Write(data)
	return err
}
func (j *ErrorResponse) Header() *HeaderWrapper {
	return j.header
}
