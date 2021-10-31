package there

import (
	"reflect"
	"testing"
)

var (
	exampleReader = BasicReader{
		"id":    []string{"100", "101"},
		"name":  []string{},
		"query": []string{"all"},
	}
)

func TestParamReader_Get(t *testing.T) {
	type args struct {
		key string
	}
	tests := []struct {
		name   string
		reader BasicReader
		args   args
		want   string
		want1  bool
	}{
		{
			name:   "Not Existing Param",
			reader: exampleReader,
			args:   args{key: "something"},
			want:   "",
			want1:  false,
		},
		{
			name:   "Empty (Not Existing) Param",
			reader: exampleReader,
			args:   args{key: "name"},
			want:   "",
			want1:  false,
		},
		{
			name:   "One Param From Polluted Param",
			reader: exampleReader,
			args:   args{key: "id"},
			want:   "100",
			want1:  true,
		},
		{
			name:   "Successful From Not Polluted",
			reader: exampleReader,
			args:   args{key: "query"},
			want:   "all",
			want1:  true,
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

func TestParamReader_GetDefault(t *testing.T) {
	type args struct {
		key          string
		defaultValue string
	}
	tests := []struct {
		name   string
		reader BasicReader
		args   args
		want   string
	}{
		{
			name:   "Not Existing Param",
			reader: exampleReader,
			args:   args{key: "something", defaultValue: "abc"},
			want:   "abc",
		},
		{
			name:   "Empty (Not Existing) Param",
			reader: exampleReader,
			args:   args{key: "name", defaultValue: "abc"},
			want:   "abc",
		},
		{
			name:   "One Param From Polluted Param",
			reader: exampleReader,
			args:   args{key: "id", defaultValue: "abc"},
			want:   "100",
		},
		{
			name:   "Successful From Not Polluted",
			reader: exampleReader,
			args:   args{key: "query", defaultValue: "abc"},
			want:   "all",
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

func TestParamReader_GetSlice(t *testing.T) {
	type args struct {
		key string
	}
	tests := []struct {
		name   string
		reader BasicReader
		args   args
		want   []string
		want1  bool
	}{
		{
			name:   "Not Existing Param",
			reader: exampleReader,
			args:   args{key: "something"},
			want:   nil,
			want1:  false,
		},
		{
			name:   "Empty (Not Existing) Param",
			reader: exampleReader,
			args:   args{key: "name"},
			want:   nil,
			want1:  false,
		},
		{
			name:   "One Param From Polluted Param",
			reader: exampleReader,
			args:   args{key: "id"},
			want:   []string{"100", "101"},
			want1:  true,
		},
		{
			name:   "Successful From Not Polluted",
			reader: exampleReader,
			args:   args{key: "query"},
			want:   []string{"all"},
			want1:  true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := tt.reader.GetSlice(tt.args.key)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetSlice() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("GetSlice() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestParamReader_Has(t *testing.T) {
	type args struct {
		key string
	}
	tests := []struct {
		name   string
		reader BasicReader
		args   args
		want   bool
	}{
		{
			name:   "Not Existing Param",
			reader: exampleReader,
			args:   args{key: "something"},
			want:   false,
		},
		{
			name:   "Empty (Not Existing) Param",
			reader: exampleReader,
			args:   args{key: "name"},
			want:   false,
		},
		{
			name:   "One Param From Polluted Param",
			reader: exampleReader,
			args:   args{key: "id"},
			want:   true,
		},
		{
			name:   "Successful From Not Polluted",
			reader: exampleReader,
			args:   args{key: "query"},
			want:   true,
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
