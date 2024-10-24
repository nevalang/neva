package desugarer

import (
	"testing"

	"github.com/nevalang/neva/internal/compiler"
	src "github.com/nevalang/neva/internal/compiler/sourcecode"
	"github.com/nevalang/neva/internal/compiler/sourcecode/core"
	ts "github.com/nevalang/neva/internal/compiler/sourcecode/typesystem"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Note: some cases are hard to test this way because desugarer depends on Scope object
// which is normally passed from top-level functions in this package.
func TestDesugarNetwork(t *testing.T) {
	d := Desugarer{}
	scope := src.Scope{}

	tests := []struct {
		name           string
		net            []src.Connection
		nodes          map[string]src.Node
		expectedResult handleNetworkResult
	}{
		{
			// node1:out -> node2:in
			name: "simple_1_to_1",
			net: []src.Connection{
				{
					Normal: &src.NormalConnection{
						SenderSide: []src.ConnectionSender{
							{
								PortAddr: &src.PortAddr{Node: "node1", Port: "out"},
							},
						},
						ReceiverSide: []src.ConnectionReceiver{
							{
								PortAddr: &src.PortAddr{Node: "node2", Port: "in"},
							},
						},
					},
				},
			},
			nodes: map[string]src.Node{
				"node1": {EntityRef: core.EntityRef{Pkg: "test", Name: "Node1"}},
				"node2": {EntityRef: core.EntityRef{Pkg: "test", Name: "Node2"}},
			},
			expectedResult: handleNetworkResult{
				desugaredConnections: []src.Connection{
					{
						Normal: &src.NormalConnection{
							SenderSide: []src.ConnectionSender{
								{
									PortAddr: &src.PortAddr{Node: "node1", Port: "out"},
								},
							},
							ReceiverSide: []src.ConnectionReceiver{
								{
									PortAddr: &src.PortAddr{Node: "node2", Port: "in"},
								},
							},
						},
					},
				},
				constsToInsert: map[string]src.Const{},
				nodesToInsert:  map[string]src.Node{},
			},
		},
		// [node1:out, node2:out] -> node3:in
		{
			name: "fan_in",
			net: []src.Connection{
				{
					Normal: &src.NormalConnection{
						SenderSide: []src.ConnectionSender{
							{PortAddr: &src.PortAddr{Node: "node1", Port: "out"}},
							{PortAddr: &src.PortAddr{Node: "node2", Port: "out"}},
						},
						ReceiverSide: []src.ConnectionReceiver{
							{
								PortAddr: &src.PortAddr{Node: "node3", Port: "in"},
							},
						},
					},
				},
			},
			nodes: map[string]src.Node{
				"node1": {EntityRef: core.EntityRef{Pkg: "test", Name: "Node1"}},
				"node2": {EntityRef: core.EntityRef{Pkg: "test", Name: "Node2"}},
				"node3": {EntityRef: core.EntityRef{Pkg: "test", Name: "Node3"}},
			},
			expectedResult: handleNetworkResult{
				desugaredConnections: []src.Connection{
					{
						Normal: &src.NormalConnection{
							SenderSide: []src.ConnectionSender{
								{
									PortAddr: &src.PortAddr{Node: "node1", Port: "out"},
								},
							},
							ReceiverSide: []src.ConnectionReceiver{
								{
									PortAddr: &src.PortAddr{Node: "__fanIn__1", Port: "data", Idx: compiler.Pointer(uint8(0))},
								},
							},
						},
					},
					{
						Normal: &src.NormalConnection{
							SenderSide: []src.ConnectionSender{
								{
									PortAddr: &src.PortAddr{Node: "node2", Port: "out"},
								},
							},
							ReceiverSide: []src.ConnectionReceiver{
								{
									PortAddr: &src.PortAddr{Node: "__fanIn__1", Port: "data", Idx: compiler.Pointer(uint8(1))},
								},
							},
						},
					},
					{
						Normal: &src.NormalConnection{
							SenderSide: []src.ConnectionSender{
								{
									PortAddr: &src.PortAddr{Node: "__fanIn__1", Port: "res"},
								},
							},
							ReceiverSide: []src.ConnectionReceiver{
								{
									PortAddr: &src.PortAddr{Node: "node3", Port: "in"},
								},
							},
						},
					},
				},
				constsToInsert: map[string]src.Const{},
				nodesToInsert: map[string]src.Node{
					"__fanIn__1": {
						EntityRef: core.EntityRef{Pkg: "builtin", Name: "FanIn"},
					},
				},
			},
		},
		// node1:foo -> node2:bar -> node3:baz
		{
			name: "chained",
			net: []src.Connection{
				{
					Normal: &src.NormalConnection{
						SenderSide: []src.ConnectionSender{
							{
								PortAddr: &src.PortAddr{Node: "node1", Port: "foo"},
							},
						},
						ReceiverSide: []src.ConnectionReceiver{
							{
								ChainedConnection: &src.Connection{
									Normal: &src.NormalConnection{
										SenderSide: []src.ConnectionSender{
											{
												PortAddr: &src.PortAddr{Node: "node2", Port: "bar"},
											},
										},
										ReceiverSide: []src.ConnectionReceiver{
											{
												PortAddr: &src.PortAddr{Node: "node3", Port: "baz"},
											},
										},
									},
								},
							},
						},
					},
				},
			},
			nodes: map[string]src.Node{
				"node1": {EntityRef: core.EntityRef{Pkg: "test", Name: "Node1"}},
				"node2": {EntityRef: core.EntityRef{Pkg: "test", Name: "Node2"}},
				"node3": {EntityRef: core.EntityRef{Pkg: "test", Name: "Node3"}},
			},
			expectedResult: handleNetworkResult{
				desugaredConnections: []src.Connection{
					{
						Normal: &src.NormalConnection{
							SenderSide: []src.ConnectionSender{
								{
									PortAddr: &src.PortAddr{Node: "node1", Port: "foo"},
								},
							},
							ReceiverSide: []src.ConnectionReceiver{
								{
									PortAddr: &src.PortAddr{Node: "node2", Port: "bar"},
								},
							},
						},
					},
					{
						Normal: &src.NormalConnection{
							SenderSide: []src.ConnectionSender{
								{
									PortAddr: &src.PortAddr{Node: "node2", Port: "bar"},
								},
							},
							ReceiverSide: []src.ConnectionReceiver{
								{
									PortAddr: &src.PortAddr{Node: "node3", Port: "baz"},
								},
							},
						},
					},
				},
				constsToInsert: map[string]src.Const{},
				nodesToInsert:  map[string]src.Node{},
			},
		},
		// foo:bar -> .a.b.c -> baz:bax
		{
			name: "struct_selector_chain",
			net: []src.Connection{
				{
					Normal: &src.NormalConnection{
						SenderSide: []src.ConnectionSender{
							{
								PortAddr: &src.PortAddr{Node: "foo", Port: "bar"},
							},
						},
						ReceiverSide: []src.ConnectionReceiver{
							{
								ChainedConnection: &src.Connection{
									Normal: &src.NormalConnection{
										SenderSide: []src.ConnectionSender{
											{StructSelector: []string{"a", "b", "c"}},
										},
										ReceiverSide: []src.ConnectionReceiver{
											{PortAddr: &src.PortAddr{Node: "baz", Port: "bax"}},
										},
									},
								},
							},
						},
					},
				},
			},
			nodes: map[string]src.Node{
				"foo": {EntityRef: core.EntityRef{Pkg: "test", Name: "Foo"}},
				"baz": {EntityRef: core.EntityRef{Pkg: "test", Name: "Baz"}},
			},
			expectedResult: handleNetworkResult{
				desugaredConnections: []src.Connection{
					{
						Normal: &src.NormalConnection{
							SenderSide: []src.ConnectionSender{
								{
									PortAddr: &src.PortAddr{Node: "foo", Port: "bar"},
								},
							},
							ReceiverSide: []src.ConnectionReceiver{
								{
									PortAddr: &src.PortAddr{Node: "__field__1", Port: "data"},
								},
							},
						},
					},
					{
						Normal: &src.NormalConnection{
							SenderSide: []src.ConnectionSender{
								{
									PortAddr: &src.PortAddr{Node: "__field__1", Port: "res"},
								},
							},
							ReceiverSide: []src.ConnectionReceiver{
								{
									PortAddr: &src.PortAddr{Node: "baz", Port: "bax"},
								},
							},
						},
					},
				},
				nodesToInsert: map[string]src.Node{
					"__field__1": {
						EntityRef: core.EntityRef{Pkg: "builtin", Name: "Field"},
						Directives: map[src.Directive][]string{
							compiler.BindDirective: {"__const__1"},
						},
					},
				},
				constsToInsert: map[string]src.Const{
					"__const__1": {
						TypeExpr: ts.Expr{
							Inst: &ts.InstExpr{
								Ref: core.EntityRef{Pkg: "builtin", Name: "list"},
								Args: []ts.Expr{
									{
										Inst: &ts.InstExpr{
											Ref: core.EntityRef{Pkg: "builtin", Name: "string"},
										},
									},
								},
							},
						},
						Value: src.ConstValue{
							Message: &src.MsgLiteral{
								List: []src.ConstValue{
									{Message: &src.MsgLiteral{Str: compiler.Pointer("a")}},
									{Message: &src.MsgLiteral{Str: compiler.Pointer("b")}},
									{Message: &src.MsgLiteral{Str: compiler.Pointer("c")}},
								},
							},
						},
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := d.handleNetwork(tt.net, tt.nodes, scope)

			require.Nil(t, err)
			assert.Equal(t, tt.expectedResult.desugaredConnections, result.desugaredConnections)
			assert.Equal(t, tt.expectedResult.constsToInsert, result.constsToInsert)
			assert.Equal(t, tt.expectedResult.nodesToInsert, result.nodesToInsert)
		})
	}
}
