package core

import (
	"errors"

	"github.com/emil14/refactored-garbanzo/internal/types"
)

var (
	ErrModNotFound = errors.New("module not found in env")
)

type Module interface {
	Interface() ModuleInterface
	SpawnWorker(scope map[string]Module) (NodeIO, error)
}

type ModuleInterface struct {
	In  InportsInterface
	Out OutportsInterface
}

type InportsInterface PortsInterface

type OutportsInterface PortsInterface

type PortsInterface map[string]types.Type

type NodeIO struct {
	In  NodeInports
	Out NodeOutports
}

type NodeInports map[string]chan Msg

type NodeOutports map[string]chan Msg

type Msg struct {
	Str  string
	Int  int
	Bool bool
}
