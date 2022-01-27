package encoder

import (
	"github.com/emil14/neva/internal/new/runtime"
	"github.com/emil14/neva/pkg/runtimesdk"
)

type caster struct{}

func (c caster) Cast(prog runtime.Program) (runtimesdk.Program, error) {
	return runtimesdk.Program{
		Nodes:       c.castNodes(prog.Nodes),
		Connections: c.castConnections(prog.Connections),
		StartPort: &runtimesdk.PortAddr{
			Node: prog.StartPort.Node,
			Port: prog.StartPort.Port,
			Slot: uint32(prog.StartPort.Idx),
		},
	}, nil
}

func (c caster) castNodes(nodes map[string]runtime.Node) map[string]*runtimesdk.Node {
	sdkNodes := make(map[string]*runtimesdk.Node, len(nodes))

	for name, node := range nodes {
		sdkNodes[name] = &runtimesdk.Node{
			Io: &runtimesdk.NodeIO{
				In:  c.castPorts(node.IO.In),
				Out: c.castPorts(node.IO.Out),
			},
			Type: runtimesdk.NodeType(node.Type),
			OpRef: &runtimesdk.OpRef{
				Pkg:  node.OperatorRef.Pkg,
				Name: node.OperatorRef.Name,
			},
			Const: c.castConst(node.Const),
		}
	}

	return sdkNodes
}

func (c caster) castPorts(ports map[string]runtime.PortMeta) map[string]*runtimesdk.PortMeta {
	sdkPorts := make(map[string]*runtimesdk.PortMeta, len(ports))

	for name, port := range ports {
		sdkPorts[name] = &runtimesdk.PortMeta{
			Slots: uint32(port.Slots),
			Buf:   uint32(port.Buf),
		}
	}

	return nil
}

func (c caster) castConst(cnst map[string]runtime.ConstValue) map[string]*runtimesdk.ConstValue {
	sdkConst := make(map[string]*runtimesdk.ConstValue, len(cnst))

	for name, value := range cnst {
		sdkConst[name] = &runtimesdk.ConstValue{
			Type:      runtimesdk.ValueType(value.Type),
			IntValue:  int64(value.Int),
			StrValue:  value.Str,
			BoolValue: value.Bool,
		}
	}

	return nil
}

func (c caster) castConnections(connections []runtime.Connection) []*runtimesdk.Connection {
	sdkConnections := make([]*runtimesdk.Connection, 0, len(connections))

	for _, connection := range connections {
		sdkConnections = append(sdkConnections, &runtimesdk.Connection{
			From: c.castPortAddr(connection.From),
			To:   c.castPortAddrs(connection.To),
		})
	}

	return sdkConnections
}

func (c caster) castPortAddr(addr runtime.PortAddr) *runtimesdk.PortAddr {
	return &runtimesdk.PortAddr{
		Node: addr.Node,
		Port: addr.Port,
		Slot: uint32(addr.Idx),
	}
}

func (c caster) castPortAddrs(addrs []runtime.PortAddr) []*runtimesdk.PortAddr {
	sdkAddrs := make([]*runtimesdk.PortAddr, 0, len(addrs))

	for _, addr := range addrs {
		sdkAddrs = append(sdkAddrs, c.castPortAddr(addr))
	}

	return sdkAddrs
}

func NewCaster() caster {
	return caster{}
}
