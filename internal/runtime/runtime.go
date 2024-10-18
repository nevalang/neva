package runtime

import (
	"context"
	"fmt"
	"sync"
	"sync/atomic"
)

var counter atomic.Uint64

type FuncCreator interface {
	Create(IO, Msg) (func(context.Context), error)
}

func Run(ctx context.Context, prog Program, registry map[string]FuncCreator) error {
	// debugValidation(prog)

	ctx, cancel := context.WithCancel(ctx)
	go func() {
		prog.Stop.Receive(ctx)
		cancel()
	}()

	runFuncs, err := deferFuncCalls(prog.FuncCalls, registry)
	if err != nil {
		return err
	}

	funcsFinished := make(chan struct{})

	go func() {
		runFuncs(context.WithValue(ctx, "cancel", cancel)) //nolint:staticcheck // SA1029
		close(funcsFinished)
	}()

	prog.Start.Send(
		ctx,
		NewStructMsg(nil, nil),
	)

	<-funcsFinished

	return nil
}

func deferFuncCalls(
	funcCalls []FuncCall,
	registry map[string]FuncCreator,
) (func(ctx context.Context), error) {
	handlers, err := createHandlers(funcCalls, registry)
	if err != nil {
		return nil, err
	}

	return func(ctx context.Context) {
		wg := sync.WaitGroup{}
		wg.Add(len(handlers))
		for i := range handlers {
			routine := handlers[i]
			go func() {
				routine(ctx)
				wg.Done()
			}()
		}
		wg.Wait()
	}, nil
}

func createHandlers(funcCalls []FuncCall, registry map[string]FuncCreator) ([]func(context.Context), error) {
	funcs := make([]func(context.Context), len(funcCalls))

	for i, call := range funcCalls {
		creator, ok := registry[call.Ref]
		if !ok {
			return nil, fmt.Errorf("func creator not found: %v", call.Ref)
		}

		handler, err := creator.Create(call.IO, call.Config)
		if err != nil {
			return nil, fmt.Errorf("%v: %w", call.Ref, err)
		}

		funcs[i] = handler
	}

	return funcs, nil
}

// func debugValidation(prog Program) {
// 	type info struct {
// 		PortSlotAddr
// 		FuncRef string
// 		Chan    any
// 	}

// 	receivers := map[string]info{}
// 	senders := map[string]info{}

// 	for _, call := range prog.FuncCalls {
// 		for _, inport := range call.IO.In.ports {
// 			if inport.single != nil {
// 				k := fmt.Sprint(inport.single.ch)
// 				receivers[k] = info{
// 					PortSlotAddr: PortSlotAddr{inport.single.addr, nil},
// 					FuncRef:      call.Ref,
// 					Chan:         inport.single.ch,
// 				}
// 			} else if inport.array != nil {
// 				for i, ch := range inport.array.chans {
// 					k := fmt.Sprint(ch)
// 					idx := uint8(i)
// 					receivers[k] = info{
// 						PortSlotAddr: PortSlotAddr{inport.array.addr, &idx},
// 						FuncRef:      call.Ref,
// 						Chan:         ch,
// 					}
// 				}
// 			} else {
// 				panic("empty func call!")
// 			}
// 		}

// 		for _, outport := range call.IO.Out.ports {
// 			if outport.single != nil {
// 				k := fmt.Sprint(outport.single.ch)
// 				senders[k] = info{
// 					PortSlotAddr: PortSlotAddr{outport.single.addr, nil},
// 					FuncRef:      call.Ref,
// 					Chan:         outport.single.ch,
// 				}
// 			} else if outport.array != nil {
// 				for i, ch := range outport.array.slots {
// 					k := fmt.Sprint(ch)
// 					idx := uint8(i)
// 					senders[k] = info{
// 						PortSlotAddr: PortSlotAddr{outport.array.addr, &idx},
// 						FuncRef:      call.Ref,
// 						Chan:         ch,
// 					}
// 				}
// 			} else {
// 				panic("empty func call!")
// 			}
// 		}
// 	}

// 	senders[fmt.Sprint(prog.Start.ch)] = info{
// 		PortSlotAddr: PortSlotAddr{PortAddr{Path: "prog", Port: "Start"}, nil},
// 		FuncRef:      "Program",
// 		Chan:         prog.Start,
// 	}

// 	receivers[fmt.Sprint(prog.Stop.ch)] = info{
// 		PortSlotAddr: PortSlotAddr{PortAddr{Path: "prog", Port: "Stop"}, nil},
// 		FuncRef:      "Program",
// 		Chan:         prog.Stop,
// 	}

// 	if len(senders) != len(receivers) {
// 		fmt.Printf(
// 			"[DEBUG] ===\nWARNING: len(senders)!=len(receivers), senders=%d, receivers=%d\n===\n\n",
// 			len(senders),
// 			len(receivers),
// 		)
// 	}

// 	formatSlotIndex := func(idx *uint8) string {
// 		if idx != nil {
// 			return fmt.Sprintf("[%d]", *idx)
// 		}
// 		return ""
// 	}

// 	for senderChanString, sInfo := range senders {
// 		if _, ok := receivers[senderChanString]; !ok {
// 			fmt.Printf(
// 				"[DEBUG] Unconnected Sender: %v | %v:%v%s -> ???\n",
// 				senderChanString,
// 				sInfo.PortSlotAddr.Path,
// 				sInfo.PortSlotAddr.Port,
// 				formatSlotIndex(sInfo.PortSlotAddr.Index),
// 			)
// 		}
// 	}

// 	for rChStr, rInfo := range receivers {
// 		if _, ok := senders[rChStr]; !ok {
// 			fmt.Printf(
// 				"[DEBUG] Unconnected Receiver: %v | ??? -> %v:%v%s\n",
// 				rChStr,
// 				rInfo.PortSlotAddr.Path,
// 				rInfo.PortSlotAddr.Port,
// 				formatSlotIndex(rInfo.PortSlotAddr.Index),
// 			)
// 		}
// 	}

// }
