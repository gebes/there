# ⭐️ Middlewares

To keep things simple, you can use already existing middlewares with little to no change, and you can use the simple control flow from **There** in your middlewares.

Here is an example:
```go

func main() {
   router := NewRouter()

   router.Use(middlewares.Recoverer)
   router.Use(middlewares.Cors(middlewares.AllowAllConfiguration()))
   router.Use(GlobalMiddleware)

   router.Get("/", Get).With(RouteSpecificMiddleware)

   err := router.Listen(8080)
   if err != nil {
      log.Fatalln("Could not listen to 8080", err)
   }
}

func GlobalMiddleware(request HttpRequest, next HttpResponse) HttpResponse {

   if request.Headers.GetDefault(RequestHeaderContentType, "") != ContentTypeApplicationJson {
      return Error(StatusUnsupportedMediaType, "Header " + RequestHeaderContentType + " is not " + ContentTypeApplicationJson)
   }

   return next
}

func RouteSpecificMiddleware(request HttpRequest, next HttpResponse) HttpResponse {
   return WithHeaders(MapString{
      ResponseHeaderContentLanguage: "en",
   }, next)
}
```

With the `.Use` method, you can add a global middleware. No matter on which group you call it, it will be **global**.  
On the other side, if you use the `.With` method you can only add a middleware to **one handler**! Not to a whole group.

The `GlobalMiddleware` in this code checks if the request has `application/json` as content-type. If not, the request will fail with an error.
Compared to the `GlobalMiddleware`, the `RouteSpecificMiddleware` does not change the control flow but adds data to the response.

Be careful in this example. Global middlewares will always be called first, so if the global middleware returns an error, the content-language header won't be set by the `RouteSpecificMiddleware` middleware.

## Using already existing middlewares
If you have other middlewares, which you created using other routers, then there is a high chance that you can use it in **There** without changing much.

As an example, let us have a look at the Recoverer middleware.

```go
func Recoverer(request HttpRequest, next HttpResponse) HttpResponse {
   fn := func(w http.ResponseWriter, r *http.Request) {
      defer func() {
         if rvr := recover(); rvr != nil && rvr != http.ErrAbortHandler {
            Error(StatusInternalServerError, rvr).ServeHTTP(w, r)
         }
      }()
      next.ServeHTTP(w, r)
   }
   return HttpResponseFunc(fn)
}
```
It is a trivial Recovere