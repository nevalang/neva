package runtime

import "github.com/emil14/refactored-garbanzo/internal/types"

type Module interface {
	Interface() ModuleInterface
	SpawnWorker(env map[string]Module) (NodeIO, error)
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
