package decoder

import (
	"github.com/emil14/neva/internal/runtime"
	"github.com/emil14/neva/pkg/runtimesdk"
)

type caster struct{}

func (c caster) Cast(in *runtimesdk.Program) (runtime.Program, error) {
	ports := c.castPorts(in)
	connections := c.castConnections(in)
	operators := c.castOperators(in)
	constants := c.castConstants(in)

	return runtime.Program{
		Ports:       ports,
		Connections: connections,
		Effects: runtime.Effects{
			Operators: operators,
			Constants: constants,
		},
		StartPort: runtime.AbsolutePortAddr{
			Path: in.StartPort.Path,
			Port: in.StartPort.Port,
			Idx:  uint8(in.StartPort.Idx),
		},
	}, nil
}

func (c caster) castConstants(in *runtimesdk.Program) map[runtime.AbsolutePortAddr]runtime.Msg {
	constants := make(map[runtime.AbsolutePortAddr]runtime.Msg, len(in.Constants))
	for _, constant := range in.Constants {
		addr := runtime.AbsolutePortAddr{
			Path: constant.OutPortAddr.Path,
			Port: constant.OutPortAddr.Port,
			Idx:  uint8(constant.OutPortAddr.Idx),
		}
		constants[addr] = c.castMsg(constant.Msg)
	}
	return constants
}

func (caster) castOperators(in *runtimesdk.Program) []runtime.Operator {
	operators := make([]runtime.Operator, 0, len(in.Operators))
	for _, operator := range in.Operators {
		inAddrs := make([]runtime.AbsolutePortAddr, 0, len(operator.InPortAddrs))
		for _, addr := range operator.InPortAddrs {
			inAddrs = append(inAddrs, runtime.AbsolutePortAddr{
				Path: addr.Path,
				Port: addr.Port,
				Idx:  uint8(addr.Idx),
			})
		}

		outAddrs := make([]runtime.AbsolutePortAddr, 0, len(operator.OutPortAddrs))
		for _, addr := range operator.OutPortAddrs {
			outAddrs = append(outAddrs, runtime.AbsolutePortAddr{
				Path: addr.Path,
				Port: addr.Port,
				Idx:  uint8(addr.Idx),
			})
		}

		operators = append(operators, runtime.Operator{
			Ref: runtime.OperatorRef{
				Pkg:  operator.Ref.Pkg,
				Name: operator.Ref.Name,
			},
			PortAddrs: runtime.OperatorPortAddrs{
				In:  inAddrs,
				Out: outAddrs,
			},
		})
	}
	return operators
}

func (caster) castConnections(in *runtimesdk.Program) []runtime.Connection {
	connections := make([]runtime.Connection, 0, len(in.Connections))
	for _, connection := range in.Connections {
		receivers := make([]runtime.ReceiverConnectionPoint, 0, len(connection.ReceiverConnectionPoints))
		for _, receiver := range connection.ReceiverConnectionPoints {
			receivers = append(receivers, runtime.ReceiverConnectionPoint{
				PortAddr: runtime.AbsolutePortAddr{
					Path: receiver.InPortAddr.Path,
					Port: receiver.InPortAddr.Port,
					Idx:  uint8(receiver.InPortAddr.Idx),
				},
				Type:            runtime.ConnectionPointType(receiver.Type),
				StructFieldPath: receiver.StructFieldPath,
			})
		}
		connections = append(connections, runtime.Connection{
			SenderPortAddr: runtime.AbsolutePortAddr{
				Path: connection.SenderOutPortAddr.Path,
				Port: connection.SenderOutPortAddr.Port,
				Idx:  uint8(connection.SenderOutPortAddr.Idx),
			},
			ReceiversConnectionPoints: receivers,
		})
	}
	return connections
}

func (caster) castPorts(in *runtimesdk.Program) []runtime.AbsolutePortAddr {
	ports := make([]runtime.AbsolutePortAddr, 0, len(in.Ports))
	for _, port := range in.Ports {
		ports = append(ports, runtime.AbsolutePortAddr{
			Path: port.Path,
			Port: port.Port,
			Idx:  uint8(port.Idx),
		})
	}
	return ports
}

func (c caster) castMsg(in *runtimesdk.Msg) runtime.Msg {
	msg := runtime.Msg{}

	switch in.Type {
	case runtimesdk.MsgType_VALUE_TYPE_BOOL:
		msg = runtime.Msg{
			Type: runtime.BoolMsg,
			Bool: in.Bool,
		}
	case runtimesdk.MsgType_VALUE_TYPE_INT:
		msg = runtime.Msg{
			Type: runtime.IntMsg,
			Int:  int(in.Int),
		}
	case runtimesdk.MsgType_VALUE_TYPE_STR:
		msg = runtime.Msg{
			Type: runtime.StrMsg,
			Str:  in.Str,
		}
	case runtimesdk.MsgType_VALUE_TYPE_STRUCT:
		structMsg := make(map[string]runtime.Msg, len(in.Struct))
		for k, v := range in.Struct {
			structMsg[k] = c.castMsg(v)
		}
		msg = runtime.Msg{
			Type:   runtime.StrMsg,
			Struct: structMsg,
		}
	}

	return msg
}

func NewCaster() caster {
	return caster{}
}
