package constants

import (
	"context"
	"sync"

	"github.com/emil14/neva/internal/runtime"
)

type Spawner struct{}

func (s Spawner) Spawn(ctx context.Context, consts []runtime.ConstEffect) error {
	wg := sync.WaitGroup{}
	wg.Add(len(consts))

	for i := range consts {
		c := consts[i]
		go func() {
			defer wg.Done()
			for {
				select {
				case <-ctx.Done():
					return
				default:
					c.OutPort <- c.Msg
				}
			}
		}()
	}

	wg.Wait()

	return ctx.Err()
}
