package there

import "testing"

func TestStatusText(t *testing.T) {
	type args struct {
		code int
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "200 = OK",
			args: args{code: StatusOK},
			want: "OK",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := StatusText(tt.args.code); got != tt.want {
				t.Errorf("StatusText() = %v, want %v", got, tt.want)
			}
		})
	}
}
