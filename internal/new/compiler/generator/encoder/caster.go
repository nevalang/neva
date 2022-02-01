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
		StartPort: &runtimesdk.AbsPortAddr{
			Node: prog.StartPort.Node,
			Port: prog.StartPort.Port,
			Idx:  uint32(prog.StartPort.Idx),
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
				Pkg:  node.OpRef.Pkg,
				Name: node.OpRef.Name,
			},
			Const: c.castConst(node.ConstOuts),
		}
	}

	return sdkNodes
}

func (c caster) castPorts(ports map[runtime.RelPortAddr]runtime.Port) []*runtimesdk.Port {
	sdkPorts := make([]*runtimesdk.Port, 0, len(ports))

	for addr, port := range ports {
		sdkPorts = append(sdkPorts, &runtimesdk.Port{
			Name: addr.Port,
			Idx:  uint32(addr.Idx),
			Buf:  uint32(port.Buf),
		})
	}

	return nil
}

func (c caster) castConst(messages map[runtime.RelPortAddr]runtime.ConstMsg) map[string]*runtimesdk.ConstValue {
	sdkConst := make(map[runtime.RelPortAddr]*runtimesdk.ConstValue, len(messages))

	for addr, value := range messages {
		sdkConst[addr] = &runtimesdk.ConstValue{
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
			From: c.castAbsPortAddr(connection.From),
			To:   c.castAbsPortAddrs(connection.To),
		})
	}

	return sdkConnections
}

func (c caster) castAbsPortAddr(addr runtime.AbsPortAddr) *runtimesdk.AbsPortAddr {
	return &runtimesdk.AbsPortAddr{
		Node: addr.Node,
		Port: addr.Port,
		Idx:  uint32(addr.Idx),
	}
}

func (c caster) castAbsPortAddrs(addrs []runtime.AbsPortAddr) []*runtimesdk.AbsPortAddr {
	sdkAddrs := make([]*runtimesdk.AbsPortAddr, 0, len(addrs))

	for _, addr := range addrs {
		sdkAddrs = append(sdkAddrs, c.castAbsPortAddr(addr))
	}

	return sdkAddrs
}

func NewCaster() caster {
	return caster{}
}
