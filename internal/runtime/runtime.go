package runtime

import (
	"context"
	"fmt"
	"sync"
	"sync/atomic"
)

var counter atomic.Uint64

type Runtime struct {
	funcRunner FuncRunner
}

func (p *Runtime) Run(ctx context.Context, prog Program) error {
	debugValidation(prog)

	ctx, cancel := context.WithCancel(ctx)
	go func() {
		prog.Stop.Receive(ctx)
		cancel() // this is how we normally stop the program
	}()

	// here we create runtime function instances to call them later
	funcrun, err := p.funcRunner.Run(prog.FuncCalls)
	if err != nil {
		return err
	}

	// we gonna use wg to make sure we don't terminate until either
	// program is successfully executed or context is closed by some runtime function
	wg := sync.WaitGroup{}
	wg.Add(1)

	go func() {
		ctx = context.WithValue(ctx, "cancel", cancel) //nolint:staticcheck // SA1029
		funcrun(ctx)                                   // funcrun blocks until context is closed
		wg.Done()
	}()

	// block until the program starts
	// by some node actually receiving from start port
	prog.Start.Send(ctx, &baseMsg{})

	// basically wait for the context to be closed
	// by either writing to the stop port
	// or calling cancel func by some runtime function like `panic`
	wg.Wait()

	return nil
}

func debugValidation(prog Program) {
	type info struct {
		PortAddr
		Ref  string
		Chan any
	}

	receivers := map[string]info{}
	senders := map[string]info{}
	for _, call := range prog.FuncCalls {
		for _, inport := range call.IO.In.ports {
			if inport.single != nil {
				k := fmt.Sprint(inport.single.ch)
				receivers[k] = info{inport.single.addr, call.Ref, inport.single.ch}
			} else if inport.array != nil {
				for _, ch := range inport.array.chans {
					k := fmt.Sprint(ch)
					receivers[k] = info{inport.array.addr, call.Ref, ch}
				}
			} else {
				panic("empty func call!")
			}
		}
		for _, outport := range call.IO.Out.ports {
			if outport.single != nil {
				k := fmt.Sprint(outport.single.ch)
				senders[k] = info{outport.single.addr, call.Ref, outport.single.ch}
			} else if outport.array != nil {
				for _, ch := range outport.array.slots {
					k := fmt.Sprint(ch)
					senders[k] = info{outport.array.addr, call.Ref, ch}
				}
			} else {
				panic("empty func call!")
			}
		}
	}

	if len(senders) != len(receivers) {
		fmt.Printf(
			"===\nWARNING: len(senders)!=len(receivers), senders=%d, receivers=%d\n===\n\n",
			len(senders),
			len(receivers),
		)
	}

	for sChStr, sInfo := range senders {
		if _, ok := receivers[sChStr]; !ok {
			fmt.Printf("%v:%v -> ???\n", sInfo.PortAddr.Path, sInfo.PortAddr.Port)
			continue
		}
		fmt.Printf(
			"%v:%v -> %v:%v\n",
			senders[sChStr].PortAddr.Path,
			receivers[sChStr].PortAddr.Port,
			receivers[sChStr].PortAddr.Path,
			receivers[sChStr].PortAddr.Port,
		)
	}

	for rChStr, rInfo := range receivers {
		if _, ok := senders[rChStr]; !ok {
			fmt.Printf("??? -> %v:%v\n", rInfo.PortAddr.Path, rInfo.PortAddr.Port)
			continue
		}
	}
}

func New(funcRunner FuncRunner) Runtime {
	return Runtime{
		funcRunner: funcRunner,
	}
}
