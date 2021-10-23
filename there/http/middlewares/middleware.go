package there

import (
	. "github.com/Gebes/there/there/http/request"
	. "github.com/Gebes/there/there/http/response"
)

type Middleware func(request HttpRequest) HttpResponse


