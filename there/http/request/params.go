package there

import "net/http"

type ParamReader struct {
	request http.Request
}

func (paramReader ParamReader) Map() map[string][]string {
	return paramReader.request.URL.Query()
}
