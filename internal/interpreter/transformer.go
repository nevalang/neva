package interpreter

import (
	"context"

	"github.com/nevalang/neva/internal/runtime"
	ir "github.com/nevalang/neva/pkg/ir/api"
)

type transformer struct{}

func (t transformer) Transform(ctx context.Context, ll *ir.LLProgram) (runtime.Program, error) {
	rPorts := make(runtime.Ports, len(ll.Ports))
	for _, portInfo := range ll.Ports {
		rPorts[runtime.PortAddr{
			Path: portInfo.PortAddr.Path,
			Port: portInfo.PortAddr.Port,
			Idx:  uint8(portInfo.PortAddr.Idx),
		}] = make(chan runtime.Msg, portInfo.BufSize)
	}

	rConns := make([]runtime.Connection, len(ll.Net))
	for _, conn := range ll.Net {
		senderAddr := runtime.PortAddr{
			Path: conn.SenderSide.Path,
			Port: conn.SenderSide.Port,
			Idx:  uint8(conn.SenderSide.Idx),
		}

		senderPort, ok := rPorts[senderAddr]
		if !ok {
			panic("!ok")
		}

		rSenderConnSide := runtime.SenderConnectionSide{
			Port: senderPort,
			Meta: runtime.SenderConnectionSideMeta{
				PortAddr: senderAddr,
			},
		}

		rReceivers := make([]runtime.ReceiverConnectionSide, len(conn.ReceiverSides))
		for _, rcvr := range conn.ReceiverSides {
			receiverPortAddr := runtime.PortAddr{
				Path: rcvr.PortAddr.Path,
				Port: rcvr.PortAddr.Port,
				Idx:  uint8(rcvr.PortAddr.Idx),
			}

			receiverPort, ok := rPorts[receiverPortAddr]
			if !ok {
				panic("!ok")
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

	rFuncs := make([]runtime.FuncRoutine, len(ll.Funcs))
	for _, f := range ll.Funcs {
		rIOIn := make(map[string][]chan runtime.Msg, len(f.Io.In))
		for _, addr := range f.Io.In {
			rPort := rPorts[runtime.PortAddr{
				Path: addr.Path,
				Port: addr.Port,
				Idx:  uint8(addr.Idx),
			}]
			rIOIn[addr.Port] = append(rIOIn[addr.Port], rPort)
		}

		rIOOut := make(map[string][]chan runtime.Msg, len(f.Io.Out))
		for _, addr := range f.Io.Out {
			rPort := rPorts[runtime.PortAddr{
				Path: addr.Path,
				Port: addr.Port,
				Idx:  uint8(addr.Idx),
			}]
			rIOOut[addr.Port] = append(rIOOut[addr.Port], rPort)
		}

		rMsg := t.msg(f.Params)

		rFuncs = append(rFuncs, runtime.FuncRoutine{
			Ref: runtime.FuncRef{
				Pkg:  f.Ref.Pkg,
				Name: f.Ref.Name,
			},
			IO: runtime.FuncIO{
				In:  rIOIn,
				Out: rIOOut,
			},
			MetaMsg: rMsg,
		})
	}

	return runtime.Program{
		Ports:       rPorts,
		Connections: rConns,
		Funcs:       rFuncs,
	}, nil
}

func (t transformer) msg(msg *ir.LLMsg) runtime.Msg {
	var rMsg runtime.Msg

	switch msg.Type {
	case ir.LLMsgType_LLBoolMsg:
		rMsg = runtime.NewBoolMsg(msg.Bool)
	case ir.LLMsgType_LLIntMsg:
		rMsg = runtime.NewIntMsg(msg.Int)
	case ir.LLMsgType_LLFloatMsg:
		rMsg = runtime.NewFloatMsg(msg.Float)
	case ir.LLMsgType_LLStrMsg:
		rMsg = runtime.NewStrMsg(msg.Str)
	default:
		panic("unknown message type")
	}

	return rMsg
}

func MustNewTransformer() transformer {
	return transformer{}
}
