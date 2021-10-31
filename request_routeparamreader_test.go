package there

import (
	"testing"
)

var (
	exampleRouteParamReader = RouteParamReader{
		"id":   "101",
		"name": "Max",
	}
)

func TestRouteParamReader_Get(t *testing.T) {
	type args struct {
		key string
	}
	tests := []struct {
		name   string
		reader RouteParamReader
		args   args
		want   string
		want1  bool
	}{
		{
			name:   "Existing Param",
			reader: exampleRouteParamReader,
			args:   args{key: "id"},
			want:   "101",
			want1:  true,
		},
		{
			name:   "Existing Param",
			reader: exampleRouteParamReader,
			args:   args{key: "name"},
			want:   "Max",
			want1:  true,
		},
		{
			name:   "Not Existing Param",
			reader: exampleRouteParamReader,
			args:   args{key: "something"},
			want:   "",
			want1:  false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := tt.reader.Get(tt.args.key)
			if got != tt.want {
				t.Errorf("Get() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("Get() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestRouteParamReader_GetDefault(t *testing.T) {
	type args struct {
		key          string
		defaultValue string
	}
	tests := []struct {
		name   string
		reader RouteParamReader
		args   args
		want   string
	}{
		{
			name:   "Existing Param",
			reader: exampleRouteParamReader,
			args:   args{key: "id", defaultValue: "abc"},
			want:   "101",
		},
		{
			name:   "Existing Param",
			reader: exampleRouteParamReader,
			args:   args{key: "name", defaultValue: "abc"},
			want:   "Max",
		},
		{
			name:   "Not Existing Param",
			reader: exampleRouteParamReader,
			args:   args{key: "something", defaultValue: "abc"},
			want:   "abc",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.reader.GetDefault(tt.args.key, tt.args.defaultValue); got != tt.want {
				t.Errorf("GetDefault() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRouteParamReader_Has(t *testing.T) {
	type args struct {
		key string
	}
	tests := []struct {
		name   string
		reader RouteParamReader
		args   args
		want   bool
	}{
		{
			name:   "Existing Param",
			reader: exampleRouteParamReader,
			args:   args{key: "id"},
			want:   true,
		},
		{
			name:   "Existing Param",
			reader: exampleRouteParamReader,
			args:   args{key: "name"},
			want:   true,
		},
		{
			name:   "Not Existing Param",
			reader: exampleRouteParamReader,
			args:   args{key: "something"},
			want:   false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.reader.Has(tt.args.key); got != tt.want {
				t.Errorf("Has() = %v, want %v", got, tt.want)
			}
		})
	}
}
