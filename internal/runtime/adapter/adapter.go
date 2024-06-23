package adapter

import (
	"github.com/nevalang/neva/internal/runtime"
	"github.com/nevalang/neva/internal/runtime/ir"
)

type Adapter struct{}

func (a Adapter) Adapt(irProg *ir.Program) (runtime.Program, error) {
	fanIn := irProg.FanIn()

	ports := a.getPorts(irProg.Ports, fanIn, irProg.Connections)

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

	for receiverAddr, senderAddrs := range fanIn {
		runtimeReceiverAddr := runtime.PortAddr{
			Path: receiverAddr.Path,
			Port: receiverAddr.Port,
			Idx:  receiverAddr.Idx,
		}

		receiverChan := ports[runtimeReceiverAddr]

		receiver := runtime.Receiver{
			Addr: runtimeReceiverAddr,
			Port: receiverChan,
		}

		senders := make([]runtime.Sender, 0, len(senderAddrs))

		for senderAddr := range senderAddrs {
			senderRuntimeAddr := runtime.PortAddr{
				Path: senderAddr.Path,
				Port: senderAddr.Port,
				Idx:  senderAddr.Idx,
			}
			senders = append(senders, runtime.Sender{
				Addr: senderRuntimeAddr,
				Port: ports[senderRuntimeAddr],
			})
		}

		runtimeConnections[receiver] = senders
	}

	return runtimeConnections
}

func (Adapter) getPorts(
	ports map[ir.PortAddr]struct{},
	fanIn map[ir.PortAddr]map[ir.PortAddr]struct{},
	fanOut map[ir.PortAddr]map[ir.PortAddr]struct{},
) map[runtime.PortAddr]chan runtime.IndexedMsg {
	runtimePorts := make(
		map[runtime.PortAddr]chan runtime.IndexedMsg,
		len(ports),
	)

	for irAddr := range ports {
		runtimeAddr := runtime.PortAddr{
			Path: irAddr.Path,
			Port: irAddr.Port,
			Idx:  irAddr.Idx,
		}

		var bufSize int
		if receivers, ok := fanOut[irAddr]; ok {
			bufSize = len(receivers)
		} else if senders, ok := fanIn[irAddr]; ok {
			bufSize = len(senders)
		}

		runtimePorts[runtimeAddr] = make(chan runtime.IndexedMsg, bufSize)
	}

	return runtimePorts
}

func NewAdapter() Adapter {
	return Adapter{}
}
