package there

import (
	"encoding/json"
	. "github.com/Gebes/there/there/utils"
	"net/http"
)

type JsonResponse struct {
	data   interface{}
	code   int
	header *HeaderWrapper
}

func Json(code int, data interface{}) *JsonResponse {
	r := &JsonResponse{data: data, code: code}
	r.header = Header(r).Set(ResponseHeaderContentType, ContentTypeApplicationJson)
	return r
}

func (j *JsonResponse) Execute(r *http.Request, w *http.ResponseWriter) error {
	b, err := json.Marshal(j.data)
	if err != nil {
		return err
	}
	(*w).WriteHeader(j.code)
	_, err = (*w).Write(b)
	return err
}

func (j *JsonResponse) Header() *HeaderWrapper {
	return j.header
}
