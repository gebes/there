package there

type HeaderWrapper struct {
	Values map[string][]string
	HttpResponse
}

func Header(httpResponse HttpResponse) *HeaderWrapper {
	return &HeaderWrapper{Values: map[string][]string{}, HttpResponse: httpResponse}
}

func (h *HeaderWrapper) Set(key string, values ...string) *HeaderWrapper {
	h.Values[key] = values
	return h
}

func (h *HeaderWrapper) SetAll(values map[string][]string) *HeaderWrapper {
	for key, values := range values {
		h.Set(key, values...)
	}
	return h
}

