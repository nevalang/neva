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
		Ports(map[string]Node) (map[PortAddr]chan core.Msg, error)
	}
	OpsSpawner interface {
		Spawn(map[string]Node, map[PortAddr]chan core.Msg) error
	}
	Connector interface {
		Connect([]Connection, map[PortAddr]chan core.Msg, chan<- error)
	}
)

var (
	ErrDecoder    = errors.New("decoder")
	ErrPortGen    = errors.New("port gen")
	ErrOpsSpawner = errors.New("ops spawner")
	ErrConnector  = errors.New("connector")
)

type Runtime struct {
	decoder    Decoder
	portGen    PortGen
	opsSpawner OpsSpawner
	connector  Connector
}

func (r Runtime) Run(raw []byte) error {
	prog, err := r.decoder.Decode(raw)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrDecoder, err)
	}

	ports, err := r.portGen.Ports(prog.Nodes)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrPortGen, err)
	}

	if err := r.opsSpawner.Spawn(prog.Nodes, ports); err != nil {
		return fmt.Errorf("%w: %v", ErrOpsSpawner, err)
	}

	stop := make(chan error)
	r.connector.Connect(prog.Connections, ports, stop)
	return fmt.Errorf("%w: %v", ErrConnector, <-stop)
}
