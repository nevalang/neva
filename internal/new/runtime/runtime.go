package runtime

import (
	"errors"
	"fmt"

	"github.com/emil14/neva/internal/new/core"
)

type (
	Decoder interface {
		Decode([]byte) (Program, error)
	}
	PortGen interface {
		Ports(NodeIO) core.IO
	}
	ConstSpawner interface {
		Spawn(map[string]Msg, map[core.PortAddr]chan core.Msg) error
	}
	OperatorSpawner interface {
		Spawn(OpRef, core.IO) error
	}
	Connector interface {
		Connect([]Connection, map[string]core.IO, chan<- error)
	}
)

var (
	ErrDecoder      = errors.New("decoder")
	ErrOpSpawner    = errors.New("operator spawner")
	ErrConstSpawner = errors.New("const spawner")
	ErrConnector    = errors.New("connector")
)

type Runtime struct {
	decoder      Decoder
	portGen      PortGen
	opSpawner    OperatorSpawner
	constSpawner ConstSpawner
	connector    Connector
}

func (r Runtime) Run(raw []byte) error {
	prog, err := r.decoder.Decode(raw)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrDecoder, err)
	}

	nodesIO := make(map[string]core.IO, len(prog.Nodes))
	for nodeName, node := range prog.Nodes {
		nodesIO[nodeName] = r.portGen.Ports(node.IO)

		switch node.Type {
		case OperatorNode:
			if err := r.opSpawner.Spawn(node.OpRef, nodesIO[nodeName]); err != nil {
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

	return fmt.Errorf("%w: %v", ErrConnector, err)
}
