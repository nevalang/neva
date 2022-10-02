package effector

import (
	"context"
	"errors"
	"fmt"

	"github.com/emil14/neva/internal/pkg/initutils"
	"github.com/emil14/neva/internal/runtime"
	"golang.org/x/sync/errgroup"
)

type (
	ConstantEffector interface {
		Effect(context.Context, []runtime.ConstantEffect) error
	}
	OperatorEffector interface {
		Effect(context.Context, []runtime.OperatorEffect) error
	}
	TriggerEffector interface {
		Effect(context.Context, []runtime.TriggerEffect) error
	}
)

type Effector struct {
	constant ConstantEffector
	operator OperatorEffector
	trigger  TriggerEffector
}

var (
	ErrOperatorEffector = errors.New("operator effector")
	ErrConstantEffector = errors.New("constant effector")
	ErrTriggerEffector  = errors.New("trigger effector")
)

func (e Effector) Effect(ctx context.Context, effects runtime.Effects) error {
	g, gctx := errgroup.WithContext(ctx)

	g.Go(func() error {
		if err := e.constant.Effect(gctx, effects.Constants); err != nil {
			return fmt.Errorf("%w: %v", ErrConstantEffector, err)
		}
		return nil
	})

	g.Go(func() error {
		if err := e.operator.Effect(gctx, effects.Operators); err != nil {
			return fmt.Errorf("%w: %v", ErrOperatorEffector, err)
		}
		return nil
	})

	g.Go(func() error {
		if err := e.trigger.Effect(gctx, effects.Triggers); err != nil {
			return fmt.Errorf("%w: %v", ErrTriggerEffector, err)
		}
		return nil
	})

	return g.Wait()
}

func MustNew(c ConstantEffector, o OperatorEffector, t TriggerEffector) Effector {
	initutils.NilPanic(c, o, t)
	return Effector{c, o, t}
}
