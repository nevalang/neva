package funcs

import (
	"context"
	"errors"
	"fmt"

	"github.com/emil14/neva/internal/runtime"
	"github.com/emil14/neva/internal/runtime/core"
	"github.com/emil14/neva/internal/runtime/src"
	"github.com/emil14/neva/pkg/tools"
	"golang.org/x/sync/errgroup"
)

var (
	ErrRepo         = errors.New("repo")
	ErrOperatorFunc = errors.New("operator func")
)

type (
	Repo interface {
		Func(src.FuncRef) (Func, error)
	}
	Func func(context.Context, core.IO) error
)

type Effector struct {
	repo Repo
}

func (e Effector) Effect(ctx context.Context, effects []runtime.ComponentNode) error {
	g, gctx := errgroup.WithContext(ctx)

	for i := range effects {
		effect := effects[i]

		f, err := e.repo.Func(effect.Ref)
		if err != nil {
			return fmt.Errorf("%w: ref %v, err %v", ErrRepo, effect.Ref, err)
		}

		g.Go(func() error {
			if err := f(gctx, effect.IO); err != nil {
				return fmt.Errorf("%w: ref %v, err %v", ErrOperatorFunc, effect.Ref, err)
			}
			return nil
		})
	}

	return g.Wait()
}

func MustNewEffector(repo Repo) Effector {
	tools.NilPanic(repo)
	return Effector{repo}
}
