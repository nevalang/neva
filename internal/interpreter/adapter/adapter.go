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

	return runtime.Program{
		Start: ports[runtime.PortAddr{
			Path: "in",
			Port: "start",
			Idx:  0,
			Arr:  false,
		}],
		Stop: ports[runtime.PortAddr{
			Path: "out",
			Port: "stop",
			Idx:  0,
			Arr:  false,
		}],
		Funcs: funcs,
	}, nil
}

func (Adapter) getConnections(
	fanOut map[ir.PortAddr]map[ir.PortAddr]struct{},
	ports map[runtime.PortAddr]chan runtime.IndexedMsg,
) map[runtime.Sender]runtime.Receiver {
	runtimeConnections := make(
		map[runtime.Sender]runtime.Receiver,
		len(fanOut),
	)

	for irSenderAddr, receiverAddrs := range fanOut {
		runtimeSenderAddr := runtime.PortAddr{
			Path: irSenderAddr.Path,
			Port: irSenderAddr.Port,
		}

		if irSenderAddr.Idx != nil {
			runtimeSenderAddr.Idx = *irSenderAddr.Idx
			runtimeSenderAddr.Arr = true
		}

		senderChan, ok := ports[runtimeSenderAddr]
		if !ok {
			panic("sender chan not found")
		}

		sender := runtime.Sender{
			Addr: runtimeSenderAddr,
			Port: senderChan,
		}

		receivers := make([]runtime.Receiver, 0, len(receiverAddrs))

		for receiverAddr := range receiverAddrs {
			receiverRuntimeAddr := runtime.PortAddr{
				Path: receiverAddr.Path,
				Port: receiverAddr.Port,
			}

			if receiverAddr.Idx != nil {
				receiverRuntimeAddr.Idx = *receiverAddr.Idx
				receiverRuntimeAddr.Arr = true
			}

			receiverChan, ok := ports[receiverRuntimeAddr]
			if !ok {
				panic("receiver chan not found: " + receiverRuntimeAddr.String())
			}

			receivers = append(receivers, runtime.Receiver{
				Addr: receiverRuntimeAddr,
				Port: receiverChan,
			})
		}

		runtimeConnections[sender] = receivers[0] // TODO use single receiver
	}

	return runtimeConnections
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
