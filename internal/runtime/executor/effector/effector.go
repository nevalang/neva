package effector

import (
	"context"
	"errors"
	"fmt"

	"github.com/emil14/neva/internal/core"
	"github.com/emil14/neva/internal/runtime"
	"github.com/emil14/neva/pkg/tools"
	"golang.org/x/sync/errgroup"
)

type (
	ConstantEffector interface {
		Effect(context.Context, []runtime.ConstNode) error
	}
	FuncEffector interface {
		Effect(context.Context, []runtime.FuncEffect) error
	}
	TriggerEffector interface {
		Effect(context.Context, []runtime.TriggerNode) error
	}
	VoidEffector interface {
		Effect(context.Context, []chan core.Msg) error
	}
)

type Effector struct {
	constant ConstantEffector
	operator FuncEffector
	trigger  TriggerEffector
	void     VoidEffector
}

var (
	ErrOperatorEffector = errors.New("operator effector")
	ErrConstantEffector = errors.New("constant effector")
	ErrTriggerEffector  = errors.New("trigger effector")
)

func (e Effector) Effect(ctx context.Context, effects runtime.Nodes) error {
	g, gctx := errgroup.WithContext(ctx)

	g.Go(func() error {
		if err := e.constant.Effect(gctx, effects.Const); err != nil {
			return fmt.Errorf("%w: %v", ErrConstantEffector, err)
		}
		return nil
	})

	g.Go(func() error {
		if err := e.operator.Effect(gctx, effects.Func); err != nil {
			return fmt.Errorf("%w: %v", ErrOperatorEffector, err)
		}
		return nil
	})

	g.Go(func() error {
		if err := e.trigger.Effect(gctx, effects.Trigger); err != nil {
			return fmt.Errorf("%w: %v", ErrTriggerEffector, err)
		}
		return nil
	})

	g.Go(func() error {
		return e.void.Effect(gctx, effects.Void)
	})

	return g.Wait()
}

func MustNew(c ConstantEffector, o FuncEffector, t TriggerEffector, v VoidEffector) Effector {
	tools.PanicWithNil(c, o, t, v)
	return Effector{c, o, t, v}
}
