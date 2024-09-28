package adapter

import (
	"os"

	"github.com/nevalang/neva/internal/compiler/ir"
	"github.com/nevalang/neva/internal/runtime"
)

type Adapter struct{}

func (a Adapter) Adapt(irProg *ir.Program, debug bool, debugLogFilePath string) (runtime.Program, error) {
	portToChan := a.getPorts(irProg)

	var interceptor runtime.Interceptor
	if debug {
		if err := a.clearDebugLogFile(debugLogFilePath); err != nil {
			return runtime.Program{}, err
		}
		interceptor = debugInterceptor{logger: fileLogger{debugLogFilePath}}
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
			Path: "out",
			Port: "stop",
		},
		interceptor,
	)

	return runtime.Program{
		Start:     start,
		Stop:      stop,
		FuncCalls: funcs,
	}, nil
}

func (a Adapter) clearDebugLogFile(filepath string) error {
	file, err := os.OpenFile(filepath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	return nil
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
			// if it's a receiver, we just create a new channel for it and go to the next iteration
			result[irAddr] = make(chan runtime.OrderedMsg)
			continue
		}

		// if it's a sender, we need to find corresponding receiver channel to re-use it
		receiverIrAddr := prog.Connections[irAddr]

		// now we know address of the receiver channel, so we just check if it's already created
		// if it's there, we just re-use it, otherwise we create new channel and re-use it
		if receiverChan, ok := result[receiverIrAddr]; ok {
			result[irAddr] = receiverChan // assign receiver channel to sender so they are connected
			continue
		}

		sharedChan := make(chan runtime.OrderedMsg)
		result[irAddr] = sharedChan
		result[receiverIrAddr] = sharedChan
	}

	return result
}

func NewAdapter() Adapter {
	return Adapter{}
}
