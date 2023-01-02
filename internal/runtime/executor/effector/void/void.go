package void

import (
	"context"

	"github.com/emil14/neva/internal/runtime/core"
	"golang.org/x/sync/errgroup"
)

type Effector struct{}

func (e Effector) Effect(ctx context.Context, ports []chan core.Msg) error {
	g, gctx := errgroup.WithContext(ctx)

	for i := range ports {
		p := ports[i]

		g.Go(func() error {
			for {
				select {
				case <-gctx.Done():
					return ctx.Err()
				case <-p:
				}
			}
		})
	}

	return g.Wait()
}
