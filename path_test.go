package there

import (
	"reflect"
	"testing"
)

func TestConstructPath(t *testing.T) {
	type args struct {
		pathString string
		ignoreCase bool
	}
	tests := []struct {
		name string
		args args
		want Path
	}{
		{
			name: "/home",
			args: args{
				pathString: "/home",
				ignoreCase: false,
			},
			want: Path{
				parts: []pathPart{
					{value: "home", variable: false},
				},
				ignoreCase: false,
			},
		},
		{
			name: "/user/:id",
			args: args{
				pathString: "/user/:id",
				ignoreCase: false,
			},
			want: Path{
				parts: []pathPart{
					{value: "user", variable: false},
					{value: "id", variable: true},
				},
				ignoreCase: false,
			},
		},
		{
			name: "/home",
			args: args{
				pathString: "/home",
				ignoreCase: true,
			},
			want: Path{
				parts: []pathPart{
					{value: "home", variable: false},
				},
				ignoreCase: true,
			},
		},
		{
			name: "/user/:id",
			args: args{
				pathString: "/user/:id",
				ignoreCase: true,
			},
			want: Path{
				parts: []pathPart{
					{value: "user", variable: false},
					{value: "id", variable: true},
				},
				ignoreCase: true,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ConstructPath(tt.args.pathString, tt.args.ignoreCase); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ConstructPath() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestConstructPathPanic(t *testing.T) {
	defer func() { recover() }()

	//should panic because id is defined twice
	ConstructPath(":id/:id", false)
	ConstructPath(":id/:Id", false)

	t.Errorf("did not panic")
}

func TestPath_Equals(t *testing.T) {
	type args struct {
		toCompare Path
	}
	tests := []struct {
		name string
		path Path
		args args
		want bool
	}{
		{
			name: "/",
			path: ConstructPath("/", true),
			args: args{ConstructPath("/", true)},
			want: true,
		},
		{
			name: "/",
			path: ConstructPath("/", true),
			args: args{ConstructPath("/", false)},
			want: false,
		},
		{
			name: "/",
			path: ConstructPath("/", false),
			args: args{ConstructPath("/", true)},
			want: false,
		},
		{
			name: "/",
			path: ConstructPath("/", false),
			args: args{ConstructPath("/", false)},
			want: true,
		},

		{
			name: "/home/:id == /home/:uid",
			path: ConstructPath("/home/:id", false),
			args: args{ConstructPath("/home/:uid", false)},
			want: true,
		},
		{
			name: "/home/:id != /Home/:uid",
			path: ConstructPath("/home/:id", false),
			args: args{ConstructPath("/Home/:uid", false)},
			want: false,
		},

		{
			name: "/home/:id != /home/about",
			path: ConstructPath("/home/:id", false),
			args: args{ConstructPath("/home/about", false)},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.path.Equals(tt.args.toCompare); got != tt.want {
				t.Errorf("Equals() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPath_Parse(t *testing.T) {

	type args struct {
		route string
	}
	tests := []struct {
		name  string
		path  Path
		args  args
		want  map[string]string
		want1 bool
	}{
		{
			name:  "/",
			path:  ConstructPath("/", false),
			args:  args{route: "/"},
			want:  map[string]string{},
			want1: true,
		},
		{
			name: "/:id",
			path: ConstructPath("/:id", false),
			args: args{route: "/101"},
			want: map[string]string{
				"id": "101",
			},
			want1: true,
		},
		{
			name: "/user/:id",
			path: ConstructPath("/user/:id", false),
			args: args{route: "/user/101"},
			want: map[string]string{
				"id": "101",
			},
			want1: true,
		},
		{
			name: "/user/:id",
			path: ConstructPath("/user/:id", true),
			args: args{route: "/USER/101"},
			want: map[string]string{
				"id": "101",
			},
			want1: true,
		},
		{
			name: "/user/:id",
			path: ConstructPath("/user/:id", true),
			args: args{route: "/USER/101"},
			want: map[string]string{
				"id": "101",
			},
			want1: true,
		},
		{
			name: "/USER/:id",
			path: ConstructPath("/USER/:id", true),
			args: args{route: "/useR/101"},
			want: map[string]string{
				"id": "101",
			},
			want1: true,
		},
		{
			name:  "/USER/:id",
			path:  ConstructPath("/USER/:id", false),
			args:  args{route: "/useR/101"},
			want:  nil,
			want1: false,
		},
		{
			name:  "/user/:id",
			path:  ConstructPath("/user/:id", false),
			args:  args{route: "/"},
			want:  nil,
			want1: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := tt.path.Parse(tt.args.route)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Parse() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("Parse() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestPath_ToString(t *testing.T) {
	type fields struct {
		parts      []pathPart
		ignoreCase bool
	}
	type test struct {
		name   string
		fields fields
		want   string
	}
	tests := []test{}

	add := func(constructWith string, expect string) {
		tests = append(tests, test{
			name: "\"" + constructWith + "\" -> \"" + expect + "\"",
			fields: fields{
				parts:      ConstructPath(constructWith, false).parts,
				ignoreCase: false,
			},
			want: expect,
		})
	}
	add("", "/")
	add("/", "/")
	add("//////", "/")
	add("//", "/")
	add("/home", "/home")
	add("home/", "/home")
	add("/home/", "/home")
	add("home/user", "/home/user")
	add("/home/user", "/home/user")
	add("/home/user/", "/home/user")
	add("//home/user//", "/home/user")
	add("///home//user///", "/home/user")
	add("/user/:id", "/user/:id")
	add("/user/:id/", "/user/:id")
	add("user/:id/", "/user/:id")
	add("user/:id", "/user/:id")

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := Path{
				parts:      tt.fields.parts,
				ignoreCase: tt.fields.ignoreCase,
			}
			if got := p.ToString(); got != tt.want {
				t.Errorf("ToString() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_splitUrl(t *testing.T) {
	type args struct {
		route string
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{
			name: "/url",
			args: args{"/url"},
			want: []string{"url"},
		},
		{
			name: "/",
			args: args{"/"},
			want: []string{},
		},
		{
			name: "/url///home",
			args: args{"/url///home"},
			want: []string{"url", "home"},
		},
		{
			name: "url///home",
			args: args{"url///home"},
			want: []string{"url", "home"},
		},
		{
			name: "url///home/",
			args: args{"url///home/"},
			want: []string{"url", "home"},
		},
		{
			name: "url/home",
			args: args{"url/home"},
			want: []string{"url", "home"},
		},
		{
			name: "",
			args: args{""},
			want: []string{},
		},
		{
			name: "///////",
			args: args{"///////"},
			want: []string{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := splitUrl(tt.args.route); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("splitUrl() = %v, want %v", got, tt.want)
			}
		})
	}
}
