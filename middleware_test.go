package there

import (
	"log"
	"reflect"
	"testing"
)

func TestNextMiddleware(t *testing.T) {
	nextMiddlewareResponse := Next().
		Header().Set("a", "b")

	if !reflect.DeepEqual(nextMiddlewareResponse.Header().Values["a"], []string{"b"}) {
		log.Fatalln("header did not set properly")
	}

	if nextMiddlewareResponse.Execute(nil, nil, nil) == nil {
		log.Fatalln("execute did not return error")
	}

}

func TestCorsMiddleware(t *testing.T) {
	middleware := CorsMiddleware(AllowAllConfiguration())

	response := middleware(HttpRequest{
		Method: MethodGet,
	})

	if !isNextResponse(response) {
		log.Fatalln("response is not next middleware")
	}

	if !reflect.DeepEqual(response.Header().Values[ResponseHeaderAccessControlAllowOrigin], []string{"*"}) ||
		!reflect.DeepEqual(response.Header().Values[ResponseHeaderAccessControlAllowMethods], []string{AllMethodsString}) ||
		!reflect.DeepEqual(response.Header().Values[ResponseHeaderAccessControlAllowHeaders], []string{"Accept, Content-Type, Content-Length, Authorization"}) {
		log.Fatalln("headers did not match allow all configuration", response.Header().Values)
	}

	response = middleware(HttpRequest{
		Method: MethodOptions,
	})

	switch v := response.(type) {

	case *HeaderWrapper:
		switch v.HttpResponse.(type) {
		case *emptyResponse:
			// everything is fine
		default:
			log.Fatalln("not empty response")
		}
	default:
		log.Fatalln("not empty response")
	}

}
