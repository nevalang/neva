package desugarer

import (
	"testing"

	src "github.com/nevalang/neva/internal/compiler/sourcecode"
	"github.com/nevalang/neva/internal/compiler/sourcecode/core"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDesugarNetwork(t *testing.T) {
	d := Desugarer{}
	scope := src.Scope{}

	tests := []struct {
		name           string
		net            []src.Connection
		nodes          map[string]src.Node
		expectedResult handleNetResult
	}{
		{
			// node1:out -> node2:in
			name: "one_simple_connection",
			net: []src.Connection{
				{
					Normal: &src.NormalConnection{
						SenderSide: src.ConnectionSenderSide{
							PortAddr: &src.PortAddr{Node: "node1", Port: "out"},
						},
						ReceiverSide: src.ConnectionReceiverSide{
							Receivers: []src.ConnectionReceiver{
								{PortAddr: src.PortAddr{Node: "node2", Port: "in"}},
							},
						},
					},
				},
			},
			nodes: map[string]src.Node{
				"node1": {EntityRef: core.EntityRef{Pkg: "test", Name: "Node1"}},
				"node2": {EntityRef: core.EntityRef{Pkg: "test", Name: "Node2"}},
			},
			expectedResult: handleNetResult{
				desugaredConnections: []src.Connection{
					{
						Normal: &src.NormalConnection{
							SenderSide: src.ConnectionSenderSide{
								PortAddr: &src.PortAddr{Node: "node1", Port: "out"},
							},
							ReceiverSide: src.ConnectionReceiverSide{
								Receivers: []src.ConnectionReceiver{
									{PortAddr: src.PortAddr{Node: "node2", Port: "in"}},
								},
							},
						},
					},
				},
				constsToInsert: map[string]src.Const{},
				nodesToInsert:  map[string]src.Node{},
			},
		},
		// [node1:out, node2:out] ->
		{
			name: "fan_in_connection",
			net: []src.Connection{
				{
					Normal: &src.NormalConnection{
						SenderSide: src.ConnectionSenderSide{
							PortAddr: &src.PortAddr{Node: "node1", Port: "out"},
						},
						ReceiverSide: src.ConnectionReceiverSide{
							Receivers: []src.ConnectionReceiver{
								{PortAddr: src.PortAddr{Node: "node3", Port: "in"}},
							},
						},
					},
				},
				{
					Normal: &src.NormalConnection{
						SenderSide: src.ConnectionSenderSide{
							PortAddr: &src.PortAddr{Node: "node2", Port: "out"},
						},
						ReceiverSide: src.ConnectionReceiverSide{
							Receivers: []src.ConnectionReceiver{
								{PortAddr: src.PortAddr{Node: "node3", Port: "in"}},
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
			expectedResult: handleNetResult{
				desugaredConnections: []src.Connection{
					{
						Normal: &src.NormalConnection{
							SenderSide: src.ConnectionSenderSide{
								PortAddr: &src.PortAddr{Node: "node1", Port: "out"},
							},
							ReceiverSide: src.ConnectionReceiverSide{
								Receivers: []src.ConnectionReceiver{
									{PortAddr: src.PortAddr{Node: "__fanIn__1", Port: "data", Idx: uintPtr(0)}},
								},
							},
						},
					},
					{
						Normal: &src.NormalConnection{
							SenderSide: src.ConnectionSenderSide{
								PortAddr: &src.PortAddr{Node: "node2", Port: "out"},
							},
							ReceiverSide: src.ConnectionReceiverSide{
								Receivers: []src.ConnectionReceiver{
									{PortAddr: src.PortAddr{Node: "__fanIn__1", Port: "data", Idx: uintPtr(1)}},
								},
							},
						},
					},
					{
						Normal: &src.NormalConnection{
							SenderSide: src.ConnectionSenderSide{
								PortAddr: &src.PortAddr{Node: "__fanIn__1", Port: "res"},
							},
							ReceiverSide: src.ConnectionReceiverSide{
								Receivers: []src.ConnectionReceiver{
									{PortAddr: src.PortAddr{Node: "node3", Port: "in"}},
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

func uintPtr(i uint8) *uint8 {
	return &i
}
