package adapter

import (
	"github.com/nevalang/neva/internal/compiler/ir"
	"github.com/nevalang/neva/internal/runtime"
)

type Adapter struct{}

func (a Adapter) Adapt(irProg *ir.Program) (runtime.Program, error) {
	ports := a.getPorts(irProg)

	funcs, err := a.getFuncs(irProg, ports)
	if err != nil {
		return runtime.Program{}, err
	}

	startAddr := runtime.PortAddr{Path: "in", Port: "start"}
	start := runtime.NewSingleOutport(startAddr, ports[startAddr])

	stopAddr := runtime.PortAddr{Path: "out", Port: "stop"}
	stop := runtime.NewSingleInport(ports[stopAddr])

	return runtime.Program{
		Start: start,
		Stop:  stop,
		Funcs: funcs,
	}, nil
}

func (Adapter) getPorts(prog *ir.Program) map[runtime.PortAddr]chan runtime.IndexedMsg {
	result := make(
		map[runtime.PortAddr]chan runtime.IndexedMsg,
		len(prog.Ports),
	)

	for irAddr := range prog.Ports {
		runtimeAddr := irAddrToRuntime(irAddr)

		// channel for this port might be already created in previous iterations of this loop
		// in case this is receiver and corresponding sender was already processed
		if _, ok := result[runtimeAddr]; ok {
			continue
		}

		// it was not created yet, let's see if it's sender or receiver
		if _, isSender := prog.Connections[irAddr]; !isSender {
			// it's a receiver, so we just create a new channel for it and go to the next iteration
			result[runtimeAddr] = make(chan runtime.IndexedMsg)
			continue
		}

		// if it's a sender, so we need to find corresponding receiver channel to re-use it
		var receiverAddr ir.PortAddr
		for v := range prog.Connections[irAddr] { // pick first one
			receiverAddr = v
			break
		}

		// now we know address of the receiver channel, so we just check if it's already created
		// if it's there, we just re-use it, otherwise we create new channel and re-use it
		receiverAddrRuntime := irAddrToRuntime(receiverAddr)

		ch, ok := result[receiverAddrRuntime]
		if ok {
			result[runtimeAddr] = ch // assign receiver channel to sender so they are connected
			continue
		}

		// receiver channel was not created yet,
		// so we create new one and assign it to both sender and receiver
		ch = make(chan runtime.IndexedMsg)
		result[receiverAddrRuntime] = ch
		result[runtimeAddr] = ch
	}

	return result
}

func irAddrToRuntime(irAddr ir.PortAddr) runtime.PortAddr {
	runtimeAddr := runtime.PortAddr{
		Path: irAddr.Path,
		Port: irAddr.Port,
	}

	if irAddr.Idx != nil {
		runtimeAddr.Idx = *irAddr.Idx
		runtimeAddr.Arr = true
	}
	return runtimeAddr
}

func NewAdapter() Adapter {
	return Adapter{}
}
