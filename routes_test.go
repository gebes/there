package there

import (
	"log"
	"reflect"
	"testing"
)

func TestNewRouteGroup(t *testing.T) {
	router := NewRouter()
	if router.RouteGroup.prefix != "/" {
		log.Fatalln("Route prefix is not /")
	}

	subGroup := router.Group("home")
	if subGroup.prefix != "/home/" {
		log.Fatalln("Route prefix is not /home/ but", subGroup.prefix)
	}

	handler := func(request HttpRequest) HttpResponse {
		return Status(StatusOK)
	}

	tests := map[string]func(route string, endpoint Endpoint) *RouteRouteGroupBuilder{
		MethodGet:     subGroup.Get,
		MethodPost:    subGroup.Post,
		MethodPatch:   subGroup.Patch,
		MethodDelete:  subGroup.Delete,
		MethodConnect: subGroup.Connect,
		MethodHead:    subGroup.Head,
		MethodTrace:   subGroup.Trace,
		MethodPut:     subGroup.Put,
		MethodOptions: subGroup.Options,
	}

	for method, methodFunc := range tests {
		if len(subGroup.routes) != 0 {
			log.Fatalln("Amount of routes should be 0")
		}

		h := methodFunc("/", handler)

		if subGroup.routes[0].Methods[0] != method {
			log.Fatalln("Method should be", method, "but is", subGroup.routes[0].Methods[0])
		}

		if len(subGroup.routes) != 1 {
			log.Fatalln("Amount of routes should be 1")
		}

		subGroup.routes.RemoveRoute(h.Route)

		if len(subGroup.routes) != 0 {
			log.Fatalln("Amount of routes should be 0")
		}
	}

}

func TestNewRouteGroup2(t *testing.T) {
	router := NewRouter()
	group := NewRouteGroup(router, "home")
	if group.prefix != "/home/" {
		log.Fatalln("Route prefix is not /home/ but", group.prefix)
	}
}

func TestNewRouteGroup3(t *testing.T) {

	handler := func(request HttpRequest) HttpResponse {
		return Status(StatusOK)
	}

	router := NewRouter()

	router.routes = nil // this should be covered and not throw any panics

	defer func() { recover() }()

	router.
		Get("/", handler).IgnoreCase().IgnoreCase().
		Get("/home/", handler).
		Get("home", handler) // duplicate route, MUST throw

	t.Errorf("did not panic")
}

func TestRouteRouteGroupBuilder_AddMiddleware2(t *testing.T) {
	handler := func(request HttpRequest) HttpResponse {
		return Status(StatusOK)
	}
	middleware := func(request HttpRequest, next HttpResponse) HttpResponse {
		return next
	}
	router := NewRouter()
	h := router.Get("/", handler).With(middleware).With(middleware)
	if len(h.Middlewares) != 2 {
		log.Fatalln("container did not have two middlewares")
	}
}

func TestRoute_OverlapsWith2(t *testing.T) {
	routeA := &Route{
		Endpoint:    nil,
		Methods:     []string{MethodGet},
		Path:        ConstructPath("/home", false),
		Middlewares: nil,
	}
	routeB := &Route{
		Endpoint:    nil,
		Methods:     []string{MethodGet, MethodPost},
		Path:        ConstructPath("/home", false),
		Middlewares: nil,
	}
	routeC := &Route{
		Endpoint:    nil,
		Methods:     []string{MethodGet, MethodPost},
		Path:        ConstructPath("/HOME", false),
		Middlewares: nil,
	}
	if !routeA.OverlapsWith(*routeB) {
		log.Fatalln("routes a and b must overlap!")
	}
	if routeA.OverlapsWith(*routeC) {
		log.Fatalln("routes a and c must NOT overlap!")
	}

}

func TestRouteGroup_Connect(t *testing.T) {
	type fields struct {
		Router *Router
		prefix string
	}
	type args struct {
		route    string
		endpoint Endpoint
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *RouteRouteGroupBuilder
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			group := &RouteGroup{
				Router: tt.fields.Router,
				prefix: tt.fields.prefix,
			}
			if got := group.Connect(tt.args.route, tt.args.endpoint); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Connect() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRouteGroup_Delete(t *testing.T) {
	type fields struct {
		Router *Router
		prefix string
	}
	type args struct {
		route    string
		endpoint Endpoint
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *RouteRouteGroupBuilder
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			group := &RouteGroup{
				Router: tt.fields.Router,
				prefix: tt.fields.prefix,
			}
			if got := group.Delete(tt.args.route, tt.args.endpoint); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Delete() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRouteGroup_Get(t *testing.T) {
	type fields struct {
		Router *Router
		prefix string
	}
	type args struct {
		route    string
		endpoint Endpoint
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *RouteRouteGroupBuilder
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			group := &RouteGroup{
				Router: tt.fields.Router,
				prefix: tt.fields.prefix,
			}
			if got := group.Get(tt.args.route, tt.args.endpoint); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Get() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRouteGroup_Group(t *testing.T) {
	type fields struct {
		Router *Router
		prefix string
	}
	type args struct {
		prefix string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *RouteGroup
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			group := RouteGroup{
				Router: tt.fields.Router,
				prefix: tt.fields.prefix,
			}
			if got := group.Group(tt.args.prefix); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Group() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRouteGroup_Handle(t *testing.T) {
	type fields struct {
		Router *Router
		prefix string
	}
	type args struct {
		path     string
		endpoint Endpoint
		methods  []string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *RouteRouteGroupBuilder
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			group := &RouteGroup{
				Router: tt.fields.Router,
				prefix: tt.fields.prefix,
			}
			if got := group.Handle(tt.args.path, tt.args.endpoint, tt.args.methods...); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Handle() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRouteGroup_Head(t *testing.T) {
	type fields struct {
		Router *Router
		prefix string
	}
	type args struct {
		route    string
		endpoint Endpoint
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *RouteRouteGroupBuilder
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			group := &RouteGroup{
				Router: tt.fields.Router,
				prefix: tt.fields.prefix,
			}
			if got := group.Head(tt.args.route, tt.args.endpoint); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Head() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRouteGroup_Options(t *testing.T) {
	type fields struct {
		Router *Router
		prefix string
	}
	type args struct {
		route    string
		endpoint Endpoint
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *RouteRouteGroupBuilder
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			group := &RouteGroup{
				Router: tt.fields.Router,
				prefix: tt.fields.prefix,
			}
			if got := group.Options(tt.args.route, tt.args.endpoint); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Options() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRouteGroup_Patch(t *testing.T) {
	type fields struct {
		Router *Router
		prefix string
	}
	type args struct {
		route    string
		endpoint Endpoint
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *RouteRouteGroupBuilder
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			group := &RouteGroup{
				Router: tt.fields.Router,
				prefix: tt.fields.prefix,
			}
			if got := group.Patch(tt.args.route, tt.args.endpoint); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Patch() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRouteGroup_Post(t *testing.T) {
	type fields struct {
		Router *Router
		prefix string
	}
	type args struct {
		route    string
		endpoint Endpoint
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *RouteRouteGroupBuilder
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			group := &RouteGroup{
				Router: tt.fields.Router,
				prefix: tt.fields.prefix,
			}
			if got := group.Post(tt.args.route, tt.args.endpoint); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Post() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRouteGroup_Put(t *testing.T) {
	type fields struct {
		Router *Router
		prefix string
	}
	type args struct {
		route    string
		endpoint Endpoint
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *RouteRouteGroupBuilder
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			group := &RouteGroup{
				Router: tt.fields.Router,
				prefix: tt.fields.prefix,
			}
			if got := group.Put(tt.args.route, tt.args.endpoint); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Put() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRouteGroup_Trace(t *testing.T) {
	type fields struct {
		Router *Router
		prefix string
	}
	type args struct {
		route    string
		endpoint Endpoint
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *RouteRouteGroupBuilder
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			group := &RouteGroup{
				Router: tt.fields.Router,
				prefix: tt.fields.prefix,
			}
			if got := group.Trace(tt.args.route, tt.args.endpoint); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Trace() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRouteManager_AddRoute(t *testing.T) {
	type args struct {
		routeToAdd *Route
	}
	tests := []struct {
		name string
		r    RouteManager
		args args
		want *Route
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.r.AddRoute(tt.args.routeToAdd); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("AddRoute() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRouteManager_FindOverlappingRoute(t *testing.T) {
	type args struct {
		routeToCheck *Route
	}
	tests := []struct {
		name string
		r    RouteManager
		args args
		want *Route
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.r.FindOverlappingRoute(tt.args.routeToCheck); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("FindOverlappingRoute() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRouteManager_RemoveRoute(t *testing.T) {
	type args struct {
		toRemove *Route
	}
	tests := []struct {
		name string
		r    RouteManager
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
		})
	}
}

func TestRouteRouteGroupBuilder_AddMiddleware(t *testing.T) {
	type fields struct {
		Route      *Route
		RouteGroup *RouteGroup
	}
	type args struct {
		middleware Middleware
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *RouteRouteGroupBuilder
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			group := &RouteRouteGroupBuilder{
				Route:      tt.fields.Route,
				RouteGroup: tt.fields.RouteGroup,
			}
			if got := group.With(tt.args.middleware); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("With() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRouteRouteGroupBuilder_IgnoreCase(t *testing.T) {
	type fields struct {
		Route      *Route
		RouteGroup *RouteGroup
	}
	tests := []struct {
		name   string
		fields fields
		want   *RouteRouteGroupBuilder
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			group := &RouteRouteGroupBuilder{
				Route:      tt.fields.Route,
				RouteGroup: tt.fields.RouteGroup,
			}
			if got := group.IgnoreCase(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("IgnoreCase() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRoute_OverlapsWith(t *testing.T) {
	type fields struct {
		Endpoint    Endpoint
		Methods     []string
		Path        Path
		Middlewares []Middleware
	}
	type args struct {
		toCompare Route
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := Route{
				Endpoint:    tt.fields.Endpoint,
				Methods:     tt.fields.Methods,
				Path:        tt.fields.Path,
				Middlewares: tt.fields.Middlewares,
			}
			if got := e.OverlapsWith(tt.args.toCompare); got != tt.want {
				t.Errorf("OverlapsWith() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRoute_ToString(t *testing.T) {
	type fields struct {
		Endpoint    Endpoint
		Methods     []string
		Path        Path
		Middlewares []Middleware
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "ToString Uppercase",
			fields: fields{
				Endpoint: nil,
				Methods:  []string{MethodGet},
				Path: Path{
					parts: []PathPart{
						{
							value:    "",
							variable: false,
						},
					},
					ignoreCase: true,
				},
				Middlewares: nil,
			},
			want: "[GET] / *IgnoreCase",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := Route{
				Endpoint:    tt.fields.Endpoint,
				Methods:     tt.fields.Methods,
				Path:        tt.fields.Path,
				Middlewares: tt.fields.Middlewares,
			}
			if got := e.ToString(); got != tt.want {
				t.Errorf("ToString() = %v, want %v", got, tt.want)
			}
		})
	}
}
