package operators

import (
	"context"
	"errors"
	"fmt"

	"github.com/emil14/neva/internal/core"
	"github.com/emil14/neva/internal/pkg/utils"
	"github.com/emil14/neva/internal/runtime"
	"github.com/emil14/neva/internal/runtime/src"
)

var (
	ErrRepo     = errors.New("repo")
	ErrOperator = errors.New("operator")
)

type (
	Repo interface {
		Operator(ref src.OperatorRef) (func(core.IO) error, error)
	}
	PortSearcher interface {
		SearchPorts(src.OperatorPortAddrs, map[src.AbsolutePortAddr]chan core.Msg) (core.IO, error)
	}
)

type Spawner struct {
	repo Repo
}

func (s Spawner) Spawn(ctx context.Context, ops []runtime.OperatorEffect) error {
	for i := range ops {
		op, err := s.repo.Operator(ops[i].Ref) // FIXME no err on not existing operator?
		if err != nil {
			return fmt.Errorf("%w: %v", ErrRepo, err)
		}

		if err := op(ops[i].IO); err != nil {
			return fmt.Errorf("%w: ref %v, err %v", ErrOperator, ops[i].Ref, err)
		}
	}

	return nil
}

func MustNew(repo Repo) Spawner {
	utils.NilPanic(repo)
	return Spawner{repo}
}
