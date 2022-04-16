package there

import "testing"

func AssertEquals(t *testing.T, got, want string) {
	if got != want {
		t.Errorf("%v != %v", want, got)
	}
}

func TestAssert(t *testing.T) {
	defer func() { recover() }()

	Assert(false, "assertion failure")

	t.Errorf("did not panic")
}

func TestAssert2(t *testing.T) {
	Assert(true, "assertion failure")
}

func TestCheckArrayContains(t *testing.T) {
	type args struct {
		slice    []string
		toSearch string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "a b c d includes a",
			args: args{
				slice:    []string{"a", "b", "c", "d"},
				toSearch: "a",
			},
			want: true,
		},
		{
			name: "a b c d includes b",
			args: args{
				slice:    []string{"a", "b", "c", "d"},
				toSearch: "b",
			},
			want: true,
		},
		{
			name: "a b c d includes c",
			args: args{
				slice:    []string{"a", "b", "c", "d"},
				toSearch: "c",
			},
			want: true,
		},
		{
			name: "a b c d includes d",
			args: args{
				slice:    []string{"a", "b", "c", "d"},
				toSearch: "d",
			},
			want: true,
		},
		{
			name: "a b c d does not include e",
			args: args{
				slice:    []string{"a", "b", "c", "d"},
				toSearch: "e",
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := CheckArrayContains(tt.args.slice, tt.args.toSearch); got != tt.want {
				t.Errorf("CheckArrayContains() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCheckArraysOverlap(t *testing.T) {
	type args struct {
		a []string
		b []string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "a b c overlaps with c d e",
			args: args{
				a: []string{"a", "b", "c"},
				b: []string{"c", "d", "e"},
			},
			want: true,
		},
		{
			name: "a b c overlaps with d e f",
			args: args{
				a: []string{"a", "b", "c"},
				b: []string{"d", "e", "f"},
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := CheckArraysOverlap(tt.args.a, tt.args.b); got != tt.want {
				t.Errorf("CheckArraysOverlap() = %v, want %v", got, tt.want)
			}
		})
	}
}
