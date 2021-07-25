package core

import (
	"errors"

	"github.com/emil14/refactored-garbanzo/internal/types"
)

var (
	ErrModNotFound = errors.New("module not found in env")
)

type Module interface {
	Interface() Interface
}

type Interface struct {
	In  Inport
	Out Outports
}

type Deps map[string]Interface

type Inport PortsInterface

type Outports PortsInterface

type PortsInterface map[string]struct {
	t   types.Type
	arr bool
}
