package trigger

import (
	"context"
	"sync"

	"github.com/emil14/neva/internal/runtime"
)

type Effector struct{}

func (e Effector) Effect(ctx context.Context, effects []runtime.TriggerEffect) error {
	wg := sync.WaitGroup{}
	wg.Add(len(effects))

	for i := range effects {
		effect := effects[i]

		go func() {
			defer wg.Done()
			for {
				select {
				case <-ctx.Done():
					return
				case <-effect.InPort:
					select {
					case <-ctx.Done():
						return
					case effect.OutPort <- effect.Msg:
						continue
					}
				}
			}
		}()
	}

	wg.Wait()

	return ctx.Err()
}
