package there

type Middleware func(request HttpRequest, next HttpResponse) HttpResponse
