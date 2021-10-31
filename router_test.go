package there

import "testing"

func TestPort_ToAddr(t *testing.T) {
	tests := []struct {
		name string
		p    Port
		want string
	}{
		{
			name: "Port to string",
			p:    8080,
			want: ":8080",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.p.ToAddr(); got != tt.want {
				t.Errorf("ToAddr() = %v, want %v", got, tt.want)
			}
		})
	}
}
