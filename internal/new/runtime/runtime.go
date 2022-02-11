package runtime

import (
	"errors"
	"fmt"

	"github.com/emil14/neva/internal/new/core"
	"github.com/emil14/neva/internal/pkg/utils"
)

type (
	Decoder interface {
		Decode([]byte) (Program, error)
	}
	PortGen interface {
		Ports(map[FullPortAddr]Port) map[FullPortAddr]chan core.Msg
	}
	ConstSpawner interface {
		Spawn(map[FullPortAddr]ConstMsg, map[FullPortAddr]chan core.Msg) error
	}
	OperatorSpawner interface {
		Spawn([]Operator, map[FullPortAddr]chan core.Msg) error
	}
	Connector interface {
		Connect(map[FullPortAddr]chan core.Msg, []Connection) error
	}
)

var (
	ErrDecoder           = errors.New("program decoder")
	ErrOpSpawner         = errors.New("operator-node spawner")
	ErrConstSpawner      = errors.New("const-node spawner")
	ErrConnector         = errors.New("network connector")
	ErrStartPortNotFound = errors.New("start port not found")
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

	ports := r.portGen.Ports(prog.Ports)

	if err := r.connector.Connect(ports, prog.Connections); err != nil {
		return fmt.Errorf("%w: %v", ErrConnector, err)
	}

	if err := r.opSpawner.Spawn(prog.Effects.Ops, ports); err != nil {
		return fmt.Errorf("%w: %v", ErrOpSpawner, err)
	}

	if err := r.constSpawner.Spawn(prog.Effects.Const, ports); err != nil {
		return fmt.Errorf("%w: %v", ErrConstSpawner, err)
	}

	start, ok := ports[prog.StartPort]
	if !ok {
		return fmt.Errorf("%w: %v", ErrStartPortNotFound, prog.StartPort)
	}
	start <- core.NewSigMsg()

	return nil // block?
}

func MustNew(
	decoder Decoder,
	portGen PortGen,
	opSpawner OperatorSpawner,
	constSpawner ConstSpawner,
	connector Connector,
) Runtime {
	utils.NilArgsFatal(decoder, portGen, opSpawner, constSpawner, connector)

	return Runtime{
		decoder:      decoder,
		portGen:      portGen,
		opSpawner:    opSpawner,
		constSpawner: constSpawner,
		connector:    connector,
	}
}
