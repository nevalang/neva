package irgen

import (
	"testing"

	"github.com/nevalang/neva/internal/runtime/ir"
	"github.com/stretchr/testify/require"
)

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

func Test_sortPortAddrs(t *testing.T) {
	tests := []struct {
		name  string
		addrs []ir.PortAddr
		want  []ir.PortAddr
	}{
		{
			name: "messed up order",
			addrs: []ir.PortAddr{
				{Path: "b", Port: "A", Idx: 1},
				{Path: "b", Port: "A", Idx: 0},
				{Path: "a", Port: "B", Idx: 0},
				{Path: "a", Port: "B", Idx: 1},
				{Path: "a", Port: "A", Idx: 2},
				{Path: "a", Port: "A", Idx: 1},
				{Path: "a", Port: "A", Idx: 0},
			},
			want: []ir.PortAddr{
				{Path: "a", Port: "A", Idx: 0},
				{Path: "a", Port: "A", Idx: 1},
				{Path: "a", Port: "A", Idx: 2},
				{Path: "a", Port: "B", Idx: 0},
				{Path: "a", Port: "B", Idx: 1},
				{Path: "b", Port: "A", Idx: 0},
				{Path: "b", Port: "A", Idx: 1},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sortPortAddrs(tt.addrs)
			require.Equal(t, tt.want, tt.addrs)
		})
	}
}
