# ⭐️ Custom Http Response

Simply create your own HttpResponse to save time. However, if you need some inspiration, look into
the [response.go](https://github.com/gebes/there/blob/main/response.go) file.

For example, let us create a Msgpack response. By default, there does not provide a Msgpack response, because this would
require a third-party dependency. But it is not much work to create your own Msgpack HttpResponse:

```go
import (
"github.com/gebes/there/v2"
"github.com/vmihailenco/msgpack/v5"
)

//Msgpack takes a StatusCode and data which gets marshaled to Msgpack
func Msgpack(code int, data interface{}) there.Response {
msgpackData, err := msgpack.Marshal(data)
if err != nil {
panic(err)
}
return there.WithHeaders(MapString{
there.ResponseHeaderContentType: "application/x-msgpack",
}, there.Bytes(code, msgpackData))
}

func Get(request there.Request) there.Response {
return Msgpack(StatusOK, map[string]string{
"Hello": "World",
"How":   "are you?",
})
}
```

**There** provides enough lower-level HttpResponses to build another one on top of it. At the bottom, we have a "Bytes"
response, which writes the given bytes and the status code.  
Wrapped around the "Bytes" response, you can find a "WithHeaders" response, adding the ContentType header.

As you see, it is only a few lines of code to have a custom HttpResponse.