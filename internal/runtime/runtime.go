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
	// debugValidation(prog)

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
		PortSlotAddr
		FuncRef string
		Chan    any
	}

	receivers := map[string]info{}
	senders := map[string]info{}
	for _, call := range prog.FuncCalls {
		for _, inport := range call.IO.In.ports {
			if inport.single != nil {
				k := fmt.Sprint(inport.single.ch)
				receivers[k] = info{
					PortSlotAddr: PortSlotAddr{inport.single.addr, nil},
					FuncRef:      call.Ref,
					Chan:         inport.single.ch,
				}
			} else if inport.array != nil {
				for i, ch := range inport.array.chans {
					k := fmt.Sprint(ch)
					idx := uint8(i)
					receivers[k] = info{
						PortSlotAddr: PortSlotAddr{inport.array.addr, &idx},
						FuncRef:      call.Ref,
						Chan:         ch,
					}
				}
			} else {
				panic("empty func call!")
			}
		}
		for _, outport := range call.IO.Out.ports {
			if outport.single != nil {
				k := fmt.Sprint(outport.single.ch)
				senders[k] = info{
					PortSlotAddr: PortSlotAddr{outport.single.addr, nil},
					FuncRef:      call.Ref,
					Chan:         outport.single.ch,
				}
			} else if outport.array != nil {
				for i, ch := range outport.array.slots {
					k := fmt.Sprint(ch)
					idx := uint8(i)
					senders[k] = info{
						PortSlotAddr: PortSlotAddr{outport.array.addr, &idx},
						FuncRef:      call.Ref,
						Chan:         ch,
					}
				}
			} else {
				panic("empty func call!")
			}
		}
	}

	senders[fmt.Sprint(prog.Start.ch)] = info{
		PortSlotAddr: PortSlotAddr{PortAddr{Path: "prog", Port: "Start"}, nil},
		FuncRef:      "Program",
		Chan:         prog.Start,
	}

	receivers[fmt.Sprint(prog.Stop.ch)] = info{
		PortSlotAddr: PortSlotAddr{PortAddr{Path: "prog", Port: "Stop"}, nil},
		FuncRef:      "Program",
		Chan:         prog.Stop,
	}

	if len(senders) != len(receivers) {
		fmt.Printf(
			"===\nWARNING: len(senders)!=len(receivers), senders=%d, receivers=%d\n===\n\n",
			len(senders),
			len(receivers),
		)
	}

	formatSlotIndex := func(idx *uint8) string {
		if idx != nil {
			return fmt.Sprintf("[%d]", *idx)
		}
		return ""
	}

	for senderChanString, sInfo := range senders {
		rInfo, ok := receivers[senderChanString]
		if !ok {
			fmt.Printf(
				"%v | %v:%v%s -> ???\n",
				senderChanString,
				sInfo.PortSlotAddr.Path,
				sInfo.PortSlotAddr.Port,
				formatSlotIndex(sInfo.PortSlotAddr.Index),
			)
			continue
		}
		fmt.Printf(
			"%v | %v:%v%s -> %v:%v%s\n",
			senderChanString,
			sInfo.PortSlotAddr.Path,
			sInfo.PortSlotAddr.Port,
			formatSlotIndex(sInfo.PortSlotAddr.Index),
			rInfo.PortSlotAddr.Path,
			rInfo.PortSlotAddr.Port,
			formatSlotIndex(rInfo.PortSlotAddr.Index),
		)
	}

	for rChStr, rInfo := range receivers {
		if _, ok := senders[rChStr]; !ok {
			fmt.Printf(
				"%v | ??? -> %v:%v%s\n",
				rChStr,
				rInfo.PortSlotAddr.Path,
				rInfo.PortSlotAddr.Port,
				formatSlotIndex(rInfo.PortSlotAddr.Index),
			)
			continue
		}
	}
}

func New(funcRunner FuncRunner) Runtime {
	return Runtime{
		funcRunner: funcRunner,
	}
}
