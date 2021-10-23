package there

import "net/http"

type Request interface {
	ReadBody() BodyReader
	ReadParams() ParamReader
}

type HttpRequest struct {
	request http.Request

	Method string
}

func NewHttpRequest(request http.Request) HttpRequest {
	return HttpRequest{
		request: request,
		Method:  request.Method,
	}
}

func (r HttpRequest) ReadBody() *BodyReader {
	return &BodyReader{request: r.request	}
}

func (r HttpRequest) ReadParams() *ParamReader {
	return &ParamReader{request: r.request}
}
