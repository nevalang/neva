package runtime

import (
	"fmt"

	"github.com/emil14/neva/internal/core"
	"github.com/emil14/neva/internal/pkg/utils"
)

type (
	Decoder interface {
		Decode([]byte) (Program, error)
	}
	PortGen interface {
		Ports([]PortAddr) map[PortAddr]chan core.Msg
	}
	// Effector interface { // instead of const and ops
	// 	MakeEffects(Effects) error
	// }
	ConstSpawner interface {
		Spawn(map[PortAddr]ConstMsg, map[PortAddr]chan core.Msg) error
	}
	OperatorSpawner interface {
		Spawn([]Operator, map[PortAddr]chan core.Msg) error
	}
	Connector interface {
		Connect(map[PortAddr]chan core.Msg, []Relation) error
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
