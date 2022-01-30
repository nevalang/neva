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
				Pkg:  node.OperatorRef.Pkg,
				Name: node.OperatorRef.Name,
			},
			Const: c.castConst(node.ConstOut),
		}
	}

	return sdkNodes
}

func (c caster) castPorts(ports map[runtime.RelPortAddr]runtime.Port) []*runtimesdk.Port {
	sdkPorts := make([]*runtimesdk.Port, 0, len(ports))

	for addr, port := range ports {
		sdkPorts = append(sdkPorts, &runtimesdk.Port{
			Addr: &runtimesdk.RelPortAddr{
				Port: addr.Port,
				Idx:  uint32(addr.Idx),
			},
			Meta: &runtimesdk.PortMeta{
				ArrSize: uint32(port.ArrSize),
				Buf:     uint32(port.Buf),
			},
		})
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
