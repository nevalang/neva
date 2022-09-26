package opspawner

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
	ErrRepo         = errors.New("repo")
	ErrOperator     = errors.New("operator")
	ErrPortSearcher = errors.New("port searcher")
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
	repo         Repo
	portSearcher PortSearcher
}

func (s Spawner) Spawn(
	ctx context.Context,
	ops []runtime.Operator,
	ports map[src.AbsolutePortAddr]chan core.Msg,
) error {
	for i := range ops {
		op, err := s.repo.Operator(ops[i].Ref) // FIXME no err on not existing operator?
		if err != nil {
			return fmt.Errorf("%w: %v", ErrRepo, err)
		}

		io, err := s.portSearcher.SearchPorts(ops[i].PortAddrs, ports)
		if err != nil {
			return fmt.Errorf("%w: %v", ErrPortSearcher, err)
		}

		if err := op(io); err != nil {
			return fmt.Errorf("%w: ref %v, err %v", ErrOperator, ops[i].Ref, err)
		}
	}

	return nil
}

func MustNew(repo Repo, portSearcher PortSearcher) Spawner {
	utils.NilPanic(repo, portSearcher)
	return Spawner{
		repo:         repo,
		portSearcher: portSearcher,
	}
}
