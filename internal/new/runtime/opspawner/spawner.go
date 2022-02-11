package opspawner

import (
	"errors"
	"fmt"

	"github.com/emil14/neva/internal/new/core"
	"github.com/emil14/neva/internal/new/runtime"
)

var (
	ErrRepo      = errors.New("repo")
	ErrOper      = errors.New("operator")
	ErrCollector = errors.New("collector")
)

type (
	Repo interface {
		Operator(ref runtime.OpRef) (func(core.IO) error, error)
	}
	Collector interface {
		Collect(runtime.OperatorIO, map[runtime.FullPortAddr]chan core.Msg) (core.IO, error)
	}
)

type Spawner struct {
	repo      Repo
	collector Collector
}

func (s Spawner) Spawn(ops []runtime.Operator, ports map[runtime.FullPortAddr]chan core.Msg) error {
	for i := range ops {
		op, err := s.repo.Operator(ops[i].Ref)
		if err != nil {
			return fmt.Errorf("%w: %v", ErrRepo, err)
		}

		io, err := s.collector.Collect(ops[i].IO, ports)
		if err != nil {
			return fmt.Errorf("%w: %v", ErrCollector, err)
		}

		if err := op(io); err != nil {
			return fmt.Errorf("%w: %v", ErrOper, err)
		}
	}

	return nil
}

func New(repo Repo) Spawner {
	return Spawner{
		repo: repo,
	}
}
