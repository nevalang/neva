package main

import (
	"context"
	"sync"

	"github.com/emil14/neva/internal/runtime/core"
)

func Builder(ctx context.Context, io core.IO) error {
	outport, err := io.Out.Port("v")
	if err != nil {
		return err
	}

	wg := sync.WaitGroup{}
	mu := sync.Mutex{}

	for {
		wg.Add(len(io.In))
		rec := make(map[string]core.Msg, len(io.In))

		for addr := range io.In {
			wg.Add(1)
			go func(field string, port chan core.Msg) {
				v := <-port
				mu.Lock()
				rec[field] = v
				mu.Unlock()
				wg.Done()
			}(addr.Port, io.In[addr])
		}

		wg.Wait()
		outport <- core.NewDictMsg(rec)
	}
}
