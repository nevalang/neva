package desugarer

import (
	src "github.com/nevalang/neva/internal/compiler/sourcecode"
	"github.com/nevalang/neva/internal/compiler/sourcecode/core"
)

type unusedOutportsResult struct {
	voidNodeName       string
	voidNode           src.Node
	virtualConnections []src.Connection
}

func (Desugarer) handleUnusedOutports(unusedOutports nodeOutportsUsed) unusedOutportsResult {
	destructorNodeName := "__del__"

	result := unusedOutportsResult{
		voidNodeName: destructorNodeName,
		voidNode: src.Node{
			EntityRef: core.EntityRef{
				Pkg:  "builtin",
				Name: "Del",
			},
		},
		virtualConnections: make([]src.Connection, 0, len(unusedOutports.m)),
	}

	receiverSides := []src.ConnectionReceiver{
		{
			PortAddr: &src.PortAddr{
				Node: destructorNodeName,
				Port: "msg",
			},
		},
	}

	voidConns := make([]src.Connection, 0, len(unusedOutports.m))
	for nodeName, ports := range unusedOutports.m {
		for portName := range ports {
			voidConns = append(voidConns, src.Connection{
				Normal: &src.NormalConnection{
					SenderSide: []src.ConnectionSender{
						{
							PortAddr: &src.PortAddr{
								Node: nodeName,
								Port: portName,
							},
						},
					},
					ReceiverSide: receiverSides,
				},
				Meta: core.Meta{},
			})
		}
	}

	result.virtualConnections = voidConns

	return result
}

func (Desugarer) findUnusedOutports(
	component src.Component,
	scope src.Scope,
	usedNodePorts nodeOutportsUsed,
) nodeOutportsUsed {
	unusedOutports := newNodePortsMap()

	for nodeName, node := range component.Nodes {
		entity, _, err := scope.Entity(node.EntityRef)
		if err != nil {
			continue
		}
		if entity.Kind != src.InterfaceEntity && entity.Kind != src.ComponentEntity {
			continue
		}

		var io src.IO
		if entity.Kind == src.InterfaceEntity {
			io = entity.Interface.IO
		} else {
			io = entity.Component.Interface.IO
		}

		for outportName := range io.Out {
			ok := usedNodePorts.get(nodeName, outportName)
			if !ok {
				unusedOutports.set(nodeName, outportName)
			}
		}
	}

	if unusedOutports.len() == 0 {
		return nodeOutportsUsed{}
	}

	return unusedOutports
}

type nodeOutportsUsed struct {
	m map[string]map[string]struct{}
}

func (n nodeOutportsUsed) set(node string, outport string) {
	if n.m[node] == nil {
		n.m[node] = map[string]struct{}{}
	}
	n.m[node][outport] = struct{}{}
}

func (n nodeOutportsUsed) get(node, port string) bool {
	ports, ok := n.m[node]
	if !ok {
		return false
	}
	_, ok = ports[port]
	return ok
}

func (n nodeOutportsUsed) len() int {
	return len(n.m)
}

func newNodePortsMap() nodeOutportsUsed {
	return nodeOutportsUsed{
		m: map[string]map[string]struct{}{},
	}
}
