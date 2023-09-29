package proto

import (
	"errors"

	"github.com/nevalang/neva/internal/runtime"
	"github.com/nevalang/neva/pkg/ir"
)

type Adapter struct{}

func (a Adapter) Adapt(irProg *ir.Program) (runtime.Program, error) { //nolint:funlen
	rPorts := make(runtime.Ports, len(irProg.Ports))
	for _, portInfo := range irProg.Ports {
		rPorts[runtime.PortAddr{
			Path: portInfo.PortAddr.Path,
			Port: portInfo.PortAddr.Port,
			Idx:  uint8(portInfo.PortAddr.Idx),
		}] = make(chan runtime.Msg, portInfo.BufSize)
	}

	rConns := make([]runtime.Connection, len(irProg.Connections))
	for _, conn := range irProg.Connections {
		senderPortAddr := runtime.PortAddr{ // reference
			Path: conn.SenderSide.Path,
			Port: conn.SenderSide.Port,
			Idx:  uint8(conn.SenderSide.Idx),
		}

		senderPort, ok := rPorts[senderPortAddr]
		if !ok {
			return runtime.Program{}, errors.New("sender port not found")
		}

		rSenderConnSide := runtime.SenderConnectionSide{
			Port: senderPort,
			Meta: runtime.SenderConnectionSideMeta{
				PortAddr: senderPortAddr,
			},
		}

		rReceivers := make([]runtime.ReceiverConnectionSide, 0, len(conn.ReceiverSides))
		for _, rcvr := range conn.ReceiverSides {
			receiverPortAddr := runtime.PortAddr{
				Path: rcvr.PortAddr.Path,
				Port: rcvr.PortAddr.Port,
				Idx:  uint8(rcvr.PortAddr.Idx),
			}

			receiverPort, ok := rPorts[receiverPortAddr]
			if !ok {
				return runtime.Program{}, errors.New("receiver port not found")
			}

			rReceivers = append(rReceivers, runtime.ReceiverConnectionSide{
				Port: receiverPort,
				Meta: runtime.ReceiverConnectionSideMeta{
					PortAddr:  receiverPortAddr,
					Selectors: rcvr.Selectors,
				},
			})
		}

		rConns = append(rConns, runtime.Connection{
			Sender:    rSenderConnSide,
			Receivers: rReceivers,
		})
	}

	rFuncs := make([]runtime.FuncRoutine, 0, len(irProg.Funcs))
	for _, f := range irProg.Funcs {
		rIOIn := make(map[string][]chan runtime.Msg, len(f.Io.Inports))
		for _, addr := range f.Io.Inports {
			rPort := rPorts[runtime.PortAddr{
				Path: addr.Path,
				Port: addr.Port,
				Idx:  uint8(addr.Idx),
			}]
			rIOIn[addr.Port] = append(rIOIn[addr.Port], rPort)
		}

		rIOOut := make(map[string][]chan runtime.Msg, len(f.Io.Outports))
		for _, addr := range f.Io.Outports {
			rPort := rPorts[runtime.PortAddr{
				Path: addr.Path,
				Port: addr.Port,
				Idx:  uint8(addr.Idx),
			}]
			rIOOut[addr.Port] = append(rIOOut[addr.Port], rPort)
		}

		rFunc := runtime.FuncRoutine{
			Ref: f.Ref,
			IO: runtime.FuncIO{
				In:  rIOIn,
				Out: rIOOut,
			},
		}

		if f.Params != nil {
			rMsg, err := a.msg(f.Params)
			if err != nil {
				return runtime.Program{}, err
			}
			rFunc.MetaMsg = rMsg
		}

		rFuncs = append(rFuncs, rFunc)
	}

	return runtime.Program{
		Ports:       rPorts,
		Connections: rConns,
		Funcs:       rFuncs,
	}, nil
}

func (a Adapter) msg(msg *ir.Msg) (runtime.Msg, error) {
	var rMsg runtime.Msg

	//nolint:nosnakecase
	switch msg.Type {
	case ir.MsgType_MSG_TYPE_BOOL:
		rMsg = runtime.NewBoolMsg(msg.Bool)
	case ir.MsgType_MSG_TYPE_INT:
		rMsg = runtime.NewIntMsg(msg.Int)
	case ir.MsgType_MSG_TYPE_FLOAT:
		rMsg = runtime.NewFloatMsg(msg.Float)
	case ir.MsgType_MSG_TYPE_STR:
		rMsg = runtime.NewStrMsg(msg.Str)
	default:
		return nil, errors.New("unknown message type")
	}

	return rMsg, nil
}