package irgen

import "testing"

func Test_joinNodePath(t *testing.T) {
	type args struct {
		nodePath []string
		nodeName string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "simple_join",
			args: args{
				nodePath: []string{"foo", "bar"},
				nodeName: "baz",
			},
			want: "foo/bar/baz",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := joinNodePath(tt.args.nodePath, tt.args.nodeName); got != tt.want {
				t.Errorf("joinNodePath() = %v, want %v", got, tt.want)
			}
		})
	}
}
