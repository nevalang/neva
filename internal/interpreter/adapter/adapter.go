package adapter

import (
	"github.com/nevalang/neva/internal/compiler/ir"
	"github.com/nevalang/neva/internal/runtime"
)

type Adapter struct{}

func (a Adapter) Adapt(irProg *ir.Program, debug bool) (runtime.Program, error) {
	ports := a.getPorts(irProg)

	var interceptor runtime.Interceptor
	if debug {
		interceptor = debugInterceptor{}
	} else {
		interceptor = prodInterceptor{}
	}

	funcs, err := a.getFuncs(irProg, ports, interceptor)
	if err != nil {
		return runtime.Program{}, err
	}

	start := runtime.NewSingleOutport(
		runtime.PortAddr{
			Path: "in",
			Port: "start",
		},
		interceptor,
		ports[ir.PortAddr{
			Path: "in",
			Port: "start",
		}],
	)

	stop := runtime.NewSingleInport(
		ports[ir.PortAddr{
			Path: "out",
			Port: "stop",
		}],
		runtime.PortAddr{
			Path: "in",
			Port: "start",
		},
		interceptor,
	)

	return runtime.Program{
		Start:     start,
		Stop:      stop,
		FuncCalls: funcs,
	}, nil
}

func (Adapter) getPorts(prog *ir.Program) map[ir.PortAddr]chan runtime.IndexedMsg {
	result := make(
		map[ir.PortAddr]chan runtime.IndexedMsg,
		len(prog.Ports),
	)

	for senderIrAddr := range prog.Ports {
		// channel for this port might be already created in previous iterations of this loop
		// in case this is receiver and corresponding sender was already processed
		if _, ok := result[senderIrAddr]; ok {
			continue
		}

		// it was not created yet, let's see if it's sender or receiver
		if _, isSender := prog.Connections[senderIrAddr]; !isSender {
			// it's a receiver, so we just create a new channel for it and go to the next iteration
			result[senderIrAddr] = make(chan runtime.IndexedMsg)
			continue
		}

		// if it's a sender, so we need to find corresponding receiver channel to re-use it
		var receiverIrAddr ir.PortAddr
		for v := range prog.Connections[senderIrAddr] { // pick first one
			receiverIrAddr = v
			break
		}

		// now we know address of the receiver channel, so we just check if it's already created
		// if it's there, we just re-use it, otherwise we create new channel and re-use it
		if receiverChan, ok := result[receiverIrAddr]; ok {
			result[senderIrAddr] = receiverChan // assign receiver channel to sender so they are connected
			continue
		}

		// receiver channel was not created yet,
		// so we create new one and assign it to both sender and receiver
		sharedChan := make(chan runtime.IndexedMsg)
		result[senderIrAddr] = sharedChan
		result[receiverIrAddr] = sharedChan
	}

	return result
}

func NewAdapter() Adapter {
	return Adapter{}
}
