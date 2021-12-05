package there

import "testing"

func Test_isNextMiddleware(t *testing.T) {
	type args struct {
		response HttpResponse
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "anything else",
			args: args{response: Empty(202)},
			want: false,
		},
		{
			name: "next middleware",
			args: args{response: Next()},
			want: true,
		},
		{
			name: "next middleware wrapped",
			args: args{response: Next().Header()},
			want: true,
		},
		{
			name: "anything else wrapped",
			args: args{response: Empty(202).Header()},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isNextResponse(tt.args.response); got != tt.want {
				t.Errorf("isNextResponse() = %v, want %v", got, tt.want)
			}
		})
	}
}
