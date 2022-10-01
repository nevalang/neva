package operators

import (
	"context"
	"errors"
	"fmt"

	"github.com/emil14/neva/internal/core"
	"github.com/emil14/neva/internal/pkg/utils"
	"github.com/emil14/neva/internal/runtime"
	"github.com/emil14/neva/internal/runtime/src"
	"golang.org/x/sync/errgroup"
)

var (
	ErrRepo     = errors.New("repo")
	ErrOperator = errors.New("operator")
)

type (
	Repo interface {
		Operator(ref src.OperatorRef) (func(context.Context, core.IO) error, error)
	}
)

type Spawner struct {
	repo Repo
}

func (s Spawner) Spawn(ctx context.Context, opEffects []runtime.OperatorEffect) error {
	g, gctx := errgroup.WithContext(ctx)

	for i := range opEffects {
		opEffect := opEffects[i]

		g.Go(func() error {
			opFunc, err := s.repo.Operator(opEffect.Ref) // FIXME no err on not existing operator?
			if err != nil {
				return fmt.Errorf("%w: %v", ErrRepo, err)
			}

			if err := opFunc(gctx, opEffect.IO); err != nil {
				return fmt.Errorf("%w: ref %v, err %v", ErrOperator, opEffect.Ref, err)
			}

			return nil
		})
	}

	return g.Wait()
}

func MustNew(repo Repo) Spawner {
	utils.NilPanic(repo)
	return Spawner{repo}
}
