package there

import (
	"encoding/json"
	. "github.com/Gebes/there/there/utils"
	"net/http"
)

type MsgpackResponse struct {
	data   interface{}
	code   int
	header *HeaderWrapper
}

func Msgpack(code int, data interface{}) *JsonResponse {
	r := &JsonResponse{data: data, code: code}
	r.header = Header(r).Set(ResponseHeaderContentType, "application/msgpack")
	return r
}

func (j *MsgpackResponse) Execute(r *http.Request, w *http.ResponseWriter) error {
	b, err := json.Marshal(j.data)
	if err != nil {
		return err
	}
	(*w).WriteHeader(j.code)
	_, err = (*w).Write(b)
	return err
}

func (j *MsgpackResponse) Header() *HeaderWrapper {
	return j.header
}
