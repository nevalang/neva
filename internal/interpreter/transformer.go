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
			SenderConnectionSideMeta: runtime.SenderConnectionSideMeta{
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

			selectors := make([]runtime.Selector, len(rcvr.Selectors))
			for _, sel := range rcvr.Selectors {
				selectors = append(selectors, runtime.Selector{
					RecField: sel.RecField,
					ArrIdx:   sel.ArrIdx,
				})
			}

			rReceivers = append(rReceivers, runtime.ReceiverConnectionSide{
				Port: receiverPort,
				Meta: runtime.ReceiverConnectionSideMeta{
					PortAddr:  receiverPortAddr,
					Selectors: selectors,
				},
			})
		}

		rConns = append(rConns, runtime.Connection{
			Sender:    rSenderConnSide,
			Receivers: rReceivers,
		})
	}

	return runtime.Program{
		Ports:       rPorts,
		Connections: rConns,
		Funcs:       []runtime.FuncRoutine{},
	}, nil
}
