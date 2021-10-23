package handlers

import (
	. "github.com/Gebes/there/there/http/request"
	. "github.com/Gebes/there/there/http/response"
	. "github.com/Gebes/there/there/utils"
)

var (
	RouteNotFound = func (request HttpRequest) HttpResponse {
		return Json(StatusNotFound, map[string]interface{}{
			"error": "Method not found",
		})
	}
)
