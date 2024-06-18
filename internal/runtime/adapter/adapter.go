package adapter

import (
	"strings"

	"github.com/nevalang/neva/internal/runtime"
	"github.com/nevalang/neva/internal/runtime/ir"
)

type Adapter struct{}

func (a Adapter) Adapt(irProg *ir.Program) (runtime.Program, error) {
	runtimePorts := a.getPorts(irProg)
	runtimeConnections := a.getConnections(irProg)
	runtimeFuncs, err := a.getFuncs(irProg, runtimePorts)
	if err != nil {
		return runtime.Program{}, err
	}
	return runtime.Program{
		QueueChan:   make(chan runtime.QueueItem),
		Ports:       runtimePorts,
		Connections: runtimeConnections,
		Funcs:       runtimeFuncs,
	}, nil
}

func (Adapter) getConnections(
	irProg *ir.Program,
) map[runtime.PortAddr]map[runtime.PortAddr][]runtime.PortAddr {
	runtimeConnections := make(
		map[runtime.PortAddr]map[runtime.PortAddr][]runtime.PortAddr,
		len(irProg.Connections),
	)

	for sender, receivers := range irProg.Connections {
		senderPortAddr := runtime.PortAddr{
			Path: sender.Path,
			Port: sender.Port,
			Idx:  sender.Idx,
		}

		receiverChans := make(map[runtime.PortAddr][]runtime.PortAddr, len(receivers))

		for rcvr := range receivers {
			receiverPortAddr := runtime.PortAddr{
				Path: rcvr.Path,
				Port: rcvr.Port,
				Idx:  rcvr.Idx,
			}
			intermediatePorts := []runtime.PortAddr{} // TODO intermediate ports
			receiverChans[receiverPortAddr] = intermediatePorts
		}

		runtimeConnections[senderPortAddr] = receiverChans
	}

	return runtimeConnections
}

func (Adapter) getPorts(irProg *ir.Program) map[runtime.PortAddr]chan runtime.Msg {
	runtimePorts := make(
		map[runtime.PortAddr]chan runtime.Msg,
		len(irProg.Ports),
	)

	for portInfo := range irProg.Ports {
		if strings.HasSuffix(portInfo.Path, "out") {
			continue // all outports use queue
		}

		addr := runtime.PortAddr{
			Path: portInfo.Path,
			Port: portInfo.Port,
			Idx:  portInfo.Idx,
		}

		runtimePorts[addr] = make(chan runtime.Msg)
	}

	// stop outport is special
	runtimePorts[runtime.PortAddr{
		Path: "out",
		Port: "stop",
	}] = make(chan runtime.Msg)

	return runtimePorts
}

func NewAdapter() Adapter {
	return Adapter{}
}
