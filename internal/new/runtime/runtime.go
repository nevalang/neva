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
		Ports(IO) core.IO
	}
	ConstSpawner interface {
		Spawn(map[RelPortAddr]ConstMsg, map[core.PortAddr]chan core.Msg) error
	}
	OperatorSpawner interface {
		Spawn(OperatorRef, core.IO) error
	}
	NetworkConnector interface {
		Connect([]Connection, map[string]core.IO, chan<- error)
	}
)

var (
	ErrProgDecoder       = errors.New("program decoder")
	ErrOpSpawner         = errors.New("operator-node spawner")
	ErrConstSpawner      = errors.New("const-node spawner")
	ErrNetConnector      = errors.New("network connector")
	ErrStartNodeNotFound = errors.New("start node not found")
	ErrStartPortNotFound = errors.New("start port not found")
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
	for name, node := range prog.Nodes {
		nodesIO[name] = r.portGen.Ports(node.IO)
		switch node.Type {
		case OperatorNode:
			if err := r.opSpawner.Spawn(node.OpRef, nodesIO[name]); err != nil {
				return fmt.Errorf("%w: %v", ErrOpSpawner, err)
			}
		case ConstNode:
			if err := r.constSpawner.Spawn(node.ConstOuts, nodesIO[name].Out); err != nil {
				return fmt.Errorf("%w: %v", ErrConstSpawner, err)
			}
		}
	}

	startNode, ok := nodesIO[prog.StartPort.Node]
	if !ok {
		return fmt.Errorf("%w: %v", ErrStartNodeNotFound, err)
	}

	startPort, err := startNode.In.Port("start")
	if err != nil {
		return fmt.Errorf("%w: %v", ErrStartPortNotFound, err)
	}

	stopChan := make(chan error)
	r.connector.Connect(prog.Connections, nodesIO, stopChan)
	go func() { startPort <- core.NewSigMsg() }()

	return fmt.Errorf("%w: %v", ErrNetConnector, <-stopChan)
}
