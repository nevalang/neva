package desugarer

import src "github.com/nevalang/neva/pkg/sourcecode"

type voidResult struct {
	voidNodeName string
	voidNode     src.Node
	voidConns    []src.Connection
}

func (Desugarer) getVoidNodeAndConns(unusedOutports nodePortsMap) voidResult {
	voidNodeName := "__void__"

	result := voidResult{
		voidNodeName: voidNodeName,
		voidNode: src.Node{
			EntityRef: src.EntityRef{
				Pkg:  "builtin",
				Name: "Destructor",
			},
		},
		voidConns: make([]src.Connection, 0, len(unusedOutports.m)),
	}

	receiverSides := []src.ConnectionReceiver{
		{
			PortAddr: src.PortAddr{
				Node: voidNodeName,
				Port: "v",
			},
		},
	}

	voidConns := make([]src.Connection, 0, len(unusedOutports.m))
	for nodeName, ports := range unusedOutports.m {
		for portName := range ports {
			voidConns = append(voidConns, src.Connection{
				SenderSide: src.ConnectionSenderSide{
					PortAddr: &src.PortAddr{
						Node: nodeName,
						Port: portName,
					},
				},
				ReceiverSide: src.ConnectionReceiverSide{
					Receivers: receiverSides,
				},
				Meta: src.Meta{},
			})
		}
	}

	result.voidConns = voidConns

	return result
}

func (Desugarer) findUnusedOutports(
	component src.Component,
	scope src.Scope,
	usedNodePorts nodePortsMap,
	desugaredNodes map[string]src.Node,
	desugaredNet []src.Connection,
) nodePortsMap {
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
		return nodePortsMap{}
	}

	return unusedOutports
}

type nodePortsMap struct {
	m map[string]map[string]struct{}
}

func (n nodePortsMap) set(node string, outport string) {
	if n.m[node] == nil {
		n.m[node] = map[string]struct{}{}
	}
	n.m[node][outport] = struct{}{}
}

func (n nodePortsMap) get(node, port string) bool {
	ports, ok := n.m[node]
	if !ok {
		return false
	}
	_, ok = ports[port]
	return ok
}

func (n nodePortsMap) len() int {
	return len(n.m)
}

func newNodePortsMap() nodePortsMap {
	return nodePortsMap{
		m: map[string]map[string]struct{}{},
	}
}
