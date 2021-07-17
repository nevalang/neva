package runtime

import "github.com/emil14/refactored-garbanzo/internal/types"

type Module interface {
	Run(in, out map[string]chan Msg)
	Interface() (InportsInterface, OutportsInterface)
}

type InportsInterface PortsInterface

type OutportsInterface PortsInterface

type PortsInterface map[string]types.Type
