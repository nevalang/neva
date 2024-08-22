package adapter

import (
	"github.com/nevalang/neva/internal/compiler/ir"
	"github.com/nevalang/neva/internal/runtime"
)

type Adapter struct{}

func (a Adapter) Adapt(irProg *ir.Program, debug bool) (runtime.Program, error) {
	portToChan := a.getPorts(irProg)

	var interceptor runtime.Interceptor
	if debug {
		interceptor = debugInterceptor{}
	} else {
		interceptor = prodInterceptor{}
	}

	funcs, err := a.getFuncs(irProg, portToChan, interceptor)
	if err != nil {
		return runtime.Program{}, err
	}

	start := runtime.NewSingleOutport(
		runtime.PortAddr{
			Path: "in",
			Port: "start",
		},
		interceptor,
		portToChan[ir.PortAddr{
			Path: "in",
			Port: "start",
		}],
	)

	stop := runtime.NewSingleInport(
		portToChan[ir.PortAddr{
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

func (Adapter) getPorts(prog *ir.Program) map[ir.PortAddr]chan runtime.OrderedMsg {
	result := make(
		map[ir.PortAddr]chan runtime.OrderedMsg,
		len(prog.Ports),
	)

	for irAddr := range prog.Ports {
		// might be already created if it's receiver and its sender was processed
		if _, created := result[irAddr]; created {
			continue
		}

		// it was not created yet, let's see if it's sender or receiver
		if _, isSender := prog.Connections[irAddr]; !isSender {
			// it's a receiver, so we just create a new channel for it and go to the next iteration
			result[irAddr] = make(chan runtime.OrderedMsg)
			continue
		}

		// if it's a sender, so we need to find corresponding receiver channel to re-use it
		var receiverIrAddr ir.PortAddr
		for tmp := range prog.Connections[irAddr] { // pick first one
			receiverIrAddr = tmp
			break
		}

		// now we know address of the receiver channel, so we just check if it's already created
		// if it's there, we just re-use it, otherwise we create new channel and re-use it
		if receiverChan, ok := result[receiverIrAddr]; ok {
			result[irAddr] = receiverChan // assign receiver channel to sender so they are connected
			continue
		}

		// receiver channel was not created yet,
		// so we create new one and assign it to both sender and receiver
		sharedChan := make(chan runtime.OrderedMsg)
		result[irAddr] = sharedChan
		result[receiverIrAddr] = sharedChan
	}

	return result
}

func NewAdapter() Adapter {
	return Adapter{}
}
