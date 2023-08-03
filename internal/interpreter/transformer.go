package interpreter

import (
	"context"

	"github.com/nevalang/neva/internal/runtime"
	"github.com/nevalang/neva/internal/shared"
)

type transformer struct{}

func (t transformer) Transform(ctx context.Context, ll shared.LowLvlProgram) (runtime.Program, error) {
	rPorts := make(runtime.Ports, len(ll.Ports))
	for addr, buf := range ll.Ports {
		rPorts[runtime.PortAddr{
			Path: addr.Path,
			Name: addr.Name,
			Idx:  addr.Idx,
		}] = make(chan runtime.Msg, buf)
	}

	rConns := make([]runtime.Connection, len(ll.Net))
	for _, conn := range ll.Net {
		senderAddr := runtime.PortAddr{
			Path: conn.SenderSide.Path,
			Name: conn.SenderSide.Name,
			Idx:  conn.SenderSide.Idx,
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
				Name: rcvr.PortAddr.Name,
				Idx:  rcvr.PortAddr.Idx,
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
		rIOIn := make(map[string][]chan runtime.Msg, len(f.IO.In))
		for _, addr := range f.IO.In {
			rPort := rPorts[runtime.PortAddr{
				Path: addr.Path,
				Name: addr.Name,
				Idx:  addr.Idx,
			}]
			rIOIn[addr.Name] = append(rIOIn[addr.Name], rPort)
		}

		rIOOut := make(map[string][]chan runtime.Msg, len(f.IO.Out))
		for _, addr := range f.IO.Out {
			rPort := rPorts[runtime.PortAddr{
				Path: addr.Path,
				Name: addr.Name,
				Idx:  addr.Idx,
			}]
			rIOOut[addr.Name] = append(rIOOut[addr.Name], rPort)
		}

		rMsg := t.msg(f.Msg)

		rFuncs = append(rFuncs, runtime.FuncRoutine{
			Ref: runtime.FuncRef{
				Pkg:  f.Ref.Pkg,
				Name: f.Ref.Name,
			},
			IO: runtime.FuncIO{
				In:  rIOIn,
				Out: rIOOut,
			},
			Msg: rMsg,
		})
	}

	return runtime.Program{
		Ports:       rPorts,
		Connections: rConns,
		Funcs:       rFuncs,
	}, nil
}

func (t transformer) msg(msg shared.LLMsg) runtime.Msg {
	var rMsg runtime.Msg

	switch msg.Type {
	case shared.LLBoolMsg:
		rMsg = runtime.NewBoolMsg(msg.Bool)
	case shared.LLIntMsg:
		rMsg = runtime.NewIntMsg(msg.Int)
	case shared.LLFloatMsg:
		rMsg = runtime.NewFloatMsg(msg.Float)
	case shared.LLStrMsg:
		rMsg = runtime.NewStrMsg(msg.Str)
	default:
		panic("unknown message type")
	}

	return rMsg
}

func MustNewTransformer() transformer {
	return transformer{}
}
