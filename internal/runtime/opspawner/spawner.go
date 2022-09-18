package opspawner

import (
	"context"
	"errors"
	"fmt"

	"github.com/emil14/neva/internal/core"
	"github.com/emil14/neva/internal/pkg/utils"
	"github.com/emil14/neva/internal/runtime"
)

var (
	ErrRepo         = errors.New("repo")
	ErrOper         = errors.New("operator")
	ErrPortSearcher = errors.New("port searcher")
)

type (
	Repo interface {
		Operator(ref runtime.OperatorRef) (func(core.IO) error, error)
	}
	PortSearcher interface {
		SearchPorts(runtime.OperatorPortAddrs, map[runtime.AbsolutePortAddr]chan core.Msg) (core.IO, error)
	}
)

type Spawner struct {
	repo         Repo
	portSearcher PortSearcher
}

func (s Spawner) Spawn(
	ctx context.Context,
	ops []runtime.Operator,
	ports map[runtime.AbsolutePortAddr]chan core.Msg,
) error {
	for i := range ops {
		op, err := s.repo.Operator(ops[i].Ref)
		if err != nil {
			return fmt.Errorf("%w: %v", ErrRepo, err)
		}

		io, err := s.portSearcher.SearchPorts(ops[i].PortAddrs, ports)
		if err != nil {
			return fmt.Errorf("%w: %v", ErrPortSearcher, err)
		}

		if err := op(io); err != nil {
			return fmt.Errorf("%w: %v", ErrOper, err)
		}
	}

	return nil
}

func MustNew(repo Repo, portSearcher PortSearcher) Spawner {
	utils.PanicOnNil(repo, portSearcher)
	return Spawner{
		repo:         repo,
		portSearcher: portSearcher,
	}
}
