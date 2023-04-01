package main

// import (
// 	"context"
// 	"sync"
// )

// // https://github.com/nevalang/nevalang/issues/149#issuecomment-1368855353
// func Recorder(ctx context.Context, io runtime.FuncIO) error {
// 	outport, err := io.Out.SinglePort("v")
// 	if err != nil {
// 		return err
// 	}

// 	wg := sync.WaitGroup{}
// 	mu := sync.Mutex{}

// 	for {
// 		wg.Add(len(io.In))
// 		rec := make(map[string]runtime.Msg, len(io.In))

// 		for addr := range io.In {
// 			wg.Add(1)
// 			go func(field string, port chan runtime.Msg) {
// 				v := <-port
// 				mu.Lock()
// 				rec[field] = v
// 				mu.Unlock()
// 				wg.Done()
// 			}(addr.Port, io.In[addr])
// 		}

// 		wg.Wait()
// 		outport <- runtime.NewDictMsg(rec)
// 	}
// }

// // Unpacker takes message that can be nil and checks it.
// // If it's not, it sends that message to `some` outport.
// // Otherwise it sends empty map (signal) to `none` outport.
// func Unpacker(ctx context.Context, io runtime.FuncIO) error {
// 	option, err := io.In.SinglePort("option")
// 	if err != nil {
// 		return err
// 	}

// 	some, err := io.Out.SinglePort("some")
// 	if err != nil {
// 		return err
// 	}

// 	none, err := io.Out.SinglePort("none")
// 	if err != nil {
// 		return err
// 	}

// 	for opt := range option {
// 		if opt != nil {
// 			some <- opt
// 			continue
// 		}
// 		none <- runtime.NewDictMsg(map[string]runtime.Msg{})
// 	}

// 	return nil
// }
