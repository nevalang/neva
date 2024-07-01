package adapter

import (
	"github.com/nevalang/neva/internal/compiler/ir"
	"github.com/nevalang/neva/internal/runtime"
)

type Adapter struct{}

func (a Adapter) Adapt(irProg *ir.Program) (runtime.Program, error) {
	fanIn := irProg.FanIn()

	ports := a.getPorts(irProg.Ports, fanIn)

	connections := a.getConnections(fanIn, ports)

	funcs, err := a.getFuncs(irProg, ports)
	if err != nil {
		return runtime.Program{}, err
	}

	start := ports[runtime.PortAddr{
		Path: "in",
		Port: "start",
	}]

	stop := ports[runtime.PortAddr{
		Path: "out",
		Port: "stop",
	}]

	return runtime.Program{
		Start:       start,
		Stop:        stop,
		Connections: connections,
		Funcs:       funcs,
	}, nil
}

func (Adapter) getConnections(
	fanIn map[ir.PortAddr]map[ir.PortAddr]struct{},
	ports map[runtime.PortAddr]chan runtime.IndexedMsg,
) map[runtime.Receiver][]runtime.Sender {
	runtimeConnections := make(
		map[runtime.Receiver][]runtime.Sender,
		len(fanIn),
	)

	for irReceiverAddr, senderAddrs := range fanIn {
		runtimeReceiverAddr := runtime.PortAddr{
			Path: irReceiverAddr.Path,
			Port: irReceiverAddr.Port,
		}

		if irReceiverAddr.Idx != nil {
			runtimeReceiverAddr.Idx = *irReceiverAddr.Idx
			runtimeReceiverAddr.Arr = true
		}

		receiverChan, ok := ports[runtimeReceiverAddr]
		if !ok {
			panic("receiver chan not found")
		}

		receiver := runtime.Receiver{
			Addr: runtimeReceiverAddr,
			Port: receiverChan,
		}

		senders := make([]runtime.Sender, 0, len(senderAddrs))

		for senderAddr := range senderAddrs {
			senderRuntimeAddr := runtime.PortAddr{
				Path: senderAddr.Path,
				Port: senderAddr.Port,
			}

			if senderAddr.Idx != nil {
				senderRuntimeAddr.Idx = *senderAddr.Idx
				senderRuntimeAddr.Arr = true
			}

			senderChan, ok := ports[senderRuntimeAddr]
			if !ok {
				panic("sender chan not found: " + senderRuntimeAddr.String())
			}

			senders = append(senders, runtime.Sender{
				Addr: senderRuntimeAddr,
				Port: senderChan,
			})
		}

		runtimeConnections[receiver] = senders
	}

	return runtimeConnections
}

func (Adapter) getPorts(
	ports map[ir.PortAddr]struct{},
	fanIn map[ir.PortAddr]map[ir.PortAddr]struct{},
) map[runtime.PortAddr]chan runtime.IndexedMsg {
	runtimePorts := make(
		map[runtime.PortAddr]chan runtime.IndexedMsg,
		len(ports),
	)

	for irAddr := range ports {
		runtimeAddr := runtime.PortAddr{
			Path: irAddr.Path,
			Port: irAddr.Port,
		}

		if irAddr.Idx != nil {
			runtimeAddr.Idx = *irAddr.Idx
			runtimeAddr.Arr = true
		}

		// TODO figure out how to set buf for senders (we don't have fan-out)
		// IDEA we can do it in src lvl

		bufSize := 0
		// if senders, ok := fanIn[irAddr]; ok {
		// 	bufSize = len(senders)
		// }

		runtimePorts[runtimeAddr] = make(chan runtime.IndexedMsg, bufSize)
	}

	return runtimePorts
}

func NewAdapter() Adapter {
	return Adapter{}
}
