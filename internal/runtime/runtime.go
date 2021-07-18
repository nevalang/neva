package runtime

import (
	"errors"
	"fmt"
)

type Runtime interface {
	Start(env map[string]Module, root string) (io NodeIO, err error)
}

func New() Runtime {
	return runtime{}
}

type runtime struct{}

func (r runtime) Start(env map[string]Module, root string) (NodeIO, error) {
	rootMod, ok := env[root]
	if !ok {
		return NodeIO{}, fmt.Errorf("%w: '%s'", ErrModNotFound, root)
	}
	return rootMod.SpawnWorker(env)
}

var (
	ErrModNotFound = errors.New("module not found in env")
)
