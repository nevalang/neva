package constants

import (
	"context"
	"sync"

	"github.com/emil14/neva/internal/runtime"
)

type Spawner struct{}

func (s Spawner) Spawn(ctx context.Context,triggers []runtime.TriggerEffect) error {
	wg := sync.WaitGroup{}
	wg.Add(len(triggers))

	for i := range triggers {
		c := triggers[i]
		go func() {
			defer wg.Done()
			for {
				select {
				case <-ctx.Done():
					return
				default:
                                        <-c.InPort
					c.OutPort <- c.Msg
				}
			}
		}()
	}

	wg.Wait()

	return ctx.Err()
}
