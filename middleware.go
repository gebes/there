package there

type Middleware func(request Request, next Response) Response
