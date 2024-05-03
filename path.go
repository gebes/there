package there

type method uint8

const (
	methodGet method = iota
	methodHead
	methodPost
	methodPut
	methodPatch
	methodDelete
	methodConnect
	methodOptions
	methodTrace
	methods
)

func methodToInt(method string) method {
	switch method {
	case MethodHead:
		return methodHead
	case MethodPost:
		return methodPost
	case MethodPut:
		return methodPut
	case MethodPatch:
		return methodPatch
	case MethodDelete:
		return methodDelete
	case MethodConnect:
		return methodConnect
	case MethodOptions:
		return methodOptions
	case MethodTrace:
		return methodTrace
	default:
		return methodGet
	}
}

func methodToString(method method) string {
	switch method {
	case methodHead:
		return MethodHead
	case methodPost:
		return MethodPost
	case methodPut:
		return MethodPut
	case methodPatch:
		return MethodPatch
	case methodDelete:
		return MethodDelete
	case methodConnect:
		return MethodConnect
	case methodOptions:
		return MethodOptions
	case methodTrace:
		return MethodTrace
	default:
		return MethodGet
	}
}
