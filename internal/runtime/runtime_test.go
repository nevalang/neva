package runtime

import (
	"reflect"
	"testing"
	"time"

	"github.com/emil14/neva/internal/runtime/program"
)

func TestRuntime_connect(t *testing.T) {
	r := New(map[string]Operator{})

	tests := []struct {
		name string
		r    Runtime
		c    connection
	}{
		{
			name: "1->3",
			r:    r,
			c: connection{
				from: make(chan Msg),
				to: []chan Msg{
					make(chan Msg),
					make(chan Msg),
					make(chan Msg),
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			go func() {
				tt.r.connect(tt.c)
			}()

			msg := 42

			select {
			case <-time.After(time.Second):
				t.Error("timeout")
			case tt.c.from <- NewIntMsg(msg):
			}

			for i := range tt.c.to {
				select {
				case <-time.After(time.Second):
					t.Error("timeout")
				case v := <-tt.c.to[i]:
					if v != NewIntMsg(msg) {
						t.Errorf("want %d, got %d", msg, v)
					}
				}
			}
		})
	}
}

func TestRuntime_connectMany(t *testing.T) {
	r := New(map[string]Operator{})

	tests := []struct {
		name string
		r    Runtime
		cc   []connection
	}{
		{
			name: "{1->2, 1->2}",
			r:    r,
			cc: []connection{
				{
					from: make(chan Msg),
					to: []chan Msg{
						make(chan Msg),
						make(chan Msg),
					},
				},
				{
					from: make(chan Msg),
					to: []chan Msg{
						make(chan Msg),
						make(chan Msg),
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.r.connectMany(tt.cc)

			for _, c := range tt.cc {
				msg := 42

				select {
				case <-time.After(time.Second):
					t.Error("timeout")
				case c.from <- NewIntMsg(msg):
				}

				for _, ch := range c.to {
					select {
					case <-time.After(time.Second):
						t.Error("timeout")
					case v := <-ch:
						if v != NewIntMsg(msg) {
							t.Errorf("want %d, got %d", msg, v)
						}
					}
				}
			}

		})
	}
}

func TestRuntime_connections(t *testing.T) {
	r := New(map[string]Operator{})

	tests := []struct {
		name string
		f    func() (map[string]IO, []program.Connection, []connection)
	}{
		{
			name: "a.out.x -> b.in.x",
			f: func() (map[string]IO, []program.Connection, []connection) {
				aOut := map[PortAddr]chan Msg{
					{"x", 0}: make(chan Msg),
				}
				bIn := map[PortAddr]chan Msg{
					{"x", 0}: make(chan Msg),
				}
				nodesIO := map[string]IO{
					"a": {Out: aOut},
					"b": {In: bIn},
				}

				from := program.PortAddr{"a", "x", 0}
				to := []program.PortAddr{{"b", "x", 0}}
				net := []program.Connection{
					{From: from, To: to},
				}

				want := []connection{
					{
						from: aOut[PortAddr{"x", 0}],
						to: []chan Msg{
							bIn[PortAddr{"x", 0}],
						},
					},
				}

				return nodesIO, net, want
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			nodesIO, net, want := tt.f()
			if got := r.connections(nodesIO, net); !reflect.DeepEqual(got, want) {
				t.Errorf("want %v, got %v", want, got)
			}
		})
	}
}
