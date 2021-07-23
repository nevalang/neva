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
	Deps() Deps
}

type Interface struct {
	In  InportsInterface
	Out OutportsInterface
}

type Deps map[string]Interface

type InportsInterface PortsInterface

type OutportsInterface PortsInterface

type PortsInterface map[string]types.Type
