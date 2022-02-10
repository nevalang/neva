package opspawner

import (
	"errors"
	"fmt"

	"github.com/emil14/neva/internal/new/core"
	"github.com/emil14/neva/internal/new/runtime"
)

var (
	ErrRepo  = errors.New("operator not loaded from repo")
	ErrSpawn = errors.New("operator was spawned with errors")
)

type Repo interface {
	Operator(ref runtime.OperatorRef) (func(core.IO) error, error)
}

type Spawner struct {
	repo Repo
}

func (s Spawner) Spawn(ref runtime.OperatorRef, io core.IO) error {
	op, err := s.repo.Operator(ref)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrRepo, err)
	}

	if err := op(io); err != nil {
		return fmt.Errorf("%w: %v", ErrSpawn, err)
	}

	return nil
}

func New(repo Repo) Spawner {
	return Spawner{
		repo: repo,
	}
}
