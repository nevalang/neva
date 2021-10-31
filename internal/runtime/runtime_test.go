package runtime

// import (
// 	"errors"
// 	"reflect"
// 	"testing"
// 	"time"

// 	"github.com/emil14/respect/internal/runtime/program"
// )

// func TestRuntime_connect(t *testing.T) {
// 	r := New(map[string]Operator{})

// 	tests := []struct {
// 		name string
// 		r    Runtime
// 		c    pair
// 	}{
// 		{
// 			name: "1->3",
// 			r:    r,
// 			c: pair{
// 				from: make(chan Msg),
// 				to: []chan Msg{
// 					make(chan Msg),
// 					make(chan Msg),
// 					make(chan Msg),
// 				},
// 			},
// 		},
// 	}

// 	for i := range tests {
// 		tt := tests[i]

// 		t.Run(tt.name, func(t *testing.T) {
// 			go func() {
// 				tt.r.connect(tt.c)
// 			}()

// 			msg := 42

// 			select {
// 			case <-time.After(time.Second):
// 				t.Error("timeout")
// 				return
// 			case tt.c.from <- NewIntMsg(msg):
// 			}

// 			for i := range tt.c.to {
// 				select {
// 				case <-time.After(time.Second):
// 					t.Error("timeout")
// 					return
// 				case v := <-tt.c.to[i]:
// 					if v != NewIntMsg(msg) {
// 						t.Errorf("want %d, got %d", msg, v)
// 						return
// 					}
// 				}
// 			}
// 		})
// 	}
// }

// func TestRuntime_connectMany(t *testing.T) {
// 	t.Parallel()

// 	r := New(map[string]Operator{})

// 	tests := []struct {
// 		name string
// 		r    Runtime
// 		cc   []pair
// 	}{
// 		{
// 			name: "{1->2, 1->2}",
// 			r:    r,
// 			cc: []pair{
// 				{
// 					from: make(chan Msg),
// 					to: []chan Msg{
// 						make(chan Msg),
// 						make(chan Msg),
// 					},
// 				},
// 				{
// 					from: make(chan Msg),
// 					to: []chan Msg{
// 						make(chan Msg),
// 						make(chan Msg),
// 					},
// 				},
// 			},
// 		},
// 	}

// 	for i := range tests {
// 		tt := tests[i]
// 		t.Run(tt.name, func(t *testing.T) {
// 			t.Parallel()

// 			tt.r.connectMany(tt.cc)

// 			for _, c := range tt.cc {
// 				msg := 42

// 				select {
// 				case <-time.After(time.Second):
// 					t.Error("timeout")
// 					return
// 				case c.from <- NewIntMsg(msg):
// 				}

// 				for _, ch := range c.to {
// 					select {
// 					case <-time.After(time.Second):
// 						t.Error("timeout")
// 						return
// 					case v := <-ch:
// 						if v != NewIntMsg(msg) {
// 							t.Errorf("want %d, got %d", msg, v)
// 							return
// 						}
// 					}
// 				}
// 			}

// 		})
// 	}
// }

// func TestRuntime_connections(t *testing.T) {
// 	t.Parallel()

// 	r := New(map[string]Operator{})

// 	tests := []struct {
// 		name string
// 		f    func() (map[string]IO, []program.Connection, []pair)
// 	}{
// 		{
// 			name: "a.out.x -> b.in.x",
// 			f: func() (map[string]IO, []program.Connection, []pair) {
// 				aOut := map[PortAddr]chan Msg{
// 					{"x", 0}: make(chan Msg),
// 				}
// 				bIn := map[PortAddr]chan Msg{
// 					{"x", 0}: make(chan Msg),
// 				}
// 				nodesIO := map[string]IO{
// 					"a": {Out: aOut},
// 					"b": {In: bIn},
// 				}

// 				from := program.PortAddr{Node: "a", Port: "x", Idx: 0}
// 				to := []program.PortAddr{{Node: "b", Port: "x", Idx: 0}}
// 				net := []program.Connection{
// 					{From: from, To: to},
// 				}

// 				want := []pair{
// 					{
// 						from: aOut[PortAddr{"x", 0}],
// 						to: []chan Msg{
// 							bIn[PortAddr{"x", 0}],
// 						},
// 					},
// 				}

// 				return nodesIO, net, want
// 			},
// 		},
// 	}

// 	for _, tt := range tests {
// 		tt := tt

// 		t.Run(tt.name, func(t *testing.T) {
// 			t.Parallel()

// 			nodesIO, net, want := tt.f()
// 			if got := r.connections(nodesIO, net); !reflect.DeepEqual(got, want) {
// 				t.Errorf("want %v, got %v", want, got)
// 				return
// 			}
// 		})
// 	}
// }

// func TestRuntime_connectOperator(t *testing.T) {
// 	t.Parallel()

// 	tests := []struct {
// 		name     string
// 		ops      map[string]Operator
// 		operator string
// 		nodeMeta program.NodeMeta
// 		wantMsg  Msg
// 		wantErr  bool
// 	}{
// 		{
// 			name:     "operator error",
// 			operator: "err",
// 			ops: map[string]Operator{
// 				"err": func(IO) error { return errors.New("") },
// 			},
// 			nodeMeta: program.NodeMeta{},
// 			wantMsg:  nil,
// 			wantErr:  true,
// 		},
// 		{
// 			name:     "undefined operator",
// 			operator: "undefined",
// 			ops:      map[string]Operator{},
// 			nodeMeta: program.NodeMeta{},
// 			wantMsg:  nil,
// 			wantErr:  true,
// 		},
// 		{
// 			name:     "bypass success",
// 			operator: "bypass",
// 			ops: map[string]Operator{
// 				"bypass": func(io IO) error {
// 					in := io.In[PortAddr{"x", 0}]
// 					out := io.Out[PortAddr{"x", 0}]
// 					go func() {
// 						v := <-in
// 						out <- v
// 					}()
// 					return nil
// 				},
// 			},
// 			nodeMeta: program.NodeMeta{
// 				In:  map[string]uint8{"x": 0},
// 				Out: map[string]uint8{"x": 0},
// 			},
// 			wantMsg: NewIntMsg(42),
// 			wantErr: false,
// 		},
// 	}

// 	for _, tt := range tests {
// 		tt := tt

// 		t.Run(tt.name, func(t *testing.T) {
// 			t.Parallel()

// 			r := New(tt.ops)
// 			io := r.nodeIO(tt.nodeMeta)

// 			err := r.connectOperator(tt.operator, io)
// 			if err != nil {
// 				if !tt.wantErr {
// 					t.Errorf("Runtime.connectOperator() error = %v, wantErr %v", err, tt.wantErr)
// 				}
// 				return
// 			}

// 			for _, in := range io.In {
// 				in <- tt.wantMsg
// 			}
// 			for _, out := range io.Out {
// 				gotMsg := <-out
// 				if gotMsg != tt.wantMsg {
// 					t.Errorf("Runtime.connectOperator() want %v, got %v", tt.wantMsg, gotMsg)
// 					return
// 				}
// 			}
// 		})
// 	}
// }

// func TestRuntime_nodeIO(t *testing.T) {
// 	t.Parallel()

// 	r := New(nil)

// 	tests := []struct {
// 		name    string
// 		node    program.NodeMeta
// 		wantIn  Ports
// 		wantOut Ports
// 	}{
// 		{
// 			name: "",
// 			node: program.NodeMeta{
// 				In: map[string]uint8{
// 					"x": 0,
// 					"y": 1,
// 					"z": 2,
// 				},
// 				Out: map[string]uint8{
// 					"x": 0,
// 					"y": 1,
// 					"z": 2,
// 				},
// 			},
// 			wantIn: map[PortAddr]chan Msg{
// 				{"x", 0}: make(chan Msg),
// 				{"y", 0}: make(chan Msg),
// 				{"z", 0}: make(chan Msg),
// 				{"z", 1}: make(chan Msg),
// 			},
// 			wantOut: map[PortAddr]chan Msg{
// 				{"x", 0}: make(chan Msg),
// 				{"y", 0}: make(chan Msg),
// 				{"z", 0}: make(chan Msg),
// 				{"z", 1}: make(chan Msg),
// 			},
// 		},
// 	}

// 	for _, tt := range tests {
// 		tt := tt
// 		t.Run(tt.name, func(t *testing.T) {
// 			t.Parallel()

// 			io := r.nodeIO(tt.node)
// 			if !comparePorts(io.In, tt.wantIn) {
// 				t.Error("r.nodeIO: got != want")
// 				return
// 			}

// 			if !comparePorts(io.Out, tt.wantOut) {
// 				t.Error("r.nodeIO: got != want")
// 				return
// 			}
// 		})
// 	}
// }

// func comparePorts(got, want Ports) bool {
// 	if len(got) != len(want) {
// 		return false
// 	}
// 	for k := range want {
// 		if got[k] == nil {
// 			return false
// 		}
// 	}
// 	return true
// }
