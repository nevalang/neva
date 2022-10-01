package effector

import (
	"context"
	"errors"
	"fmt"

	"github.com/emil14/neva/internal/pkg/utils"
	"github.com/emil14/neva/internal/runtime"
	"golang.org/x/sync/errgroup"
)

type (
	ConstSpawner interface {
		Spawn(context.Context, []runtime.ConstEffect) error
	}
	OperatorSpawner interface {
		Spawn(context.Context, []runtime.OperatorEffect) error
	}
)

type Effector struct {
	constants ConstSpawner
	operators OperatorSpawner
}

var (
	ErrOpSpawner    = errors.New("operator-node spawner")
	ErrConstSpawner = errors.New("const spawner")
)

func (e Effector) MakeEffects(ctx context.Context, effects runtime.Effects) error {
	g, gctx := errgroup.WithContext(ctx)

	g.Go(func() error {
		if err := e.constants.Spawn(gctx, effects.Consts); err != nil {
			return fmt.Errorf("%w: %v", ErrConstSpawner, err)
		}
		return nil
	})

	g.Go(func() error {
		if err := e.operators.Spawn(gctx, effects.Operators); err != nil {
			return fmt.Errorf("%w: %v", ErrOpSpawner, err)
		}
		return nil
	})

	if err := g.Wait(); err != nil {
		return fmt.Errorf("wait: %w", err)
	}

	return nil
}

func MustNew(c ConstSpawner, o OperatorSpawner) Effector {
	utils.NilPanic(c, o)
	return Effector{c, o}
}
