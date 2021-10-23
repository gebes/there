package there



func Message(code int, message string) *JsonResponse {
	return Json(code, map[string]interface{}{
		"message": message,
	})
}
