package runtime

import (
	"errors"
	"fmt"

	"github.com/emil14/neva/internal/new/core"
)

type (
	ProgramDecoder interface {
		Decode([]byte) (Program, error)
	}
	PortGenerator interface {
		Ports(NodeIO) core.IO
	}
	ConstSpawner interface {
		Spawn(map[string]Msg, map[core.PortAddr]chan core.Msg) error
	}
	OperatorSpawner interface {
		Spawn(OperatorRef, core.IO) error
	}
	NetworkConnector interface {
		Connect([]Connection, map[string]core.IO, chan<- error)
	}
)

var (
	ErrProgDecoder  = errors.New("program decoder")
	ErrOpSpawner    = errors.New("operator-node spawner")
	ErrConstSpawner = errors.New("const-node spawner")
	ErrNetConnector = errors.New("network connector")
)

type Runtime struct {
	decoder      ProgramDecoder
	portGen      PortGenerator
	opSpawner    OperatorSpawner
	constSpawner ConstSpawner
	connector    NetworkConnector
}

func (r Runtime) Run(raw []byte) error {
	prog, err := r.decoder.Decode(raw)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrProgDecoder, err)
	}

	nodesIO := make(map[string]core.IO, len(prog.Nodes))
	for nodeName, node := range prog.Nodes {
		nodesIO[nodeName] = r.portGen.Ports(node.IO)

		switch node.Type {
		case OperatorNode:
			if err := r.opSpawner.Spawn(node.OperatorRef, nodesIO[nodeName]); err != nil {
				return fmt.Errorf("%w: %v", ErrOpSpawner, err)
			}
		case ConstNode:
			if err := r.constSpawner.Spawn(node.Const, nodesIO[nodeName].Out); err != nil {
				return fmt.Errorf("%w: %v", ErrConstSpawner, err)
			}
		}
	}

	stop := make(chan error)
	r.connector.Connect(prog.Connections, nodesIO, stop)
	err = <-stop

	return fmt.Errorf("%w: %v", ErrNetConnector, err)
}
