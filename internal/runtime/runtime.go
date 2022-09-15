package runtime

import (
	"errors"
	"fmt"

	"github.com/emil14/neva/internal/core"
	"github.com/emil14/neva/internal/pkg/utils"
	"golang.org/x/sync/errgroup"
)

type (
	Decoder interface {
		Decode([]byte) (Program, error)
	}
	PortGen interface {
		Ports([]AbsolutePortAddr) map[AbsolutePortAddr]chan core.Msg
	}
	// Effector interface { // instead of const and ops
	// 	MakeEffects(Effects) error
	// }
	ConstSpawner interface {
		Spawn(map[AbsolutePortAddr]Msg, map[AbsolutePortAddr]chan core.Msg) error
	}
	OperatorSpawner interface {
		Spawn([]Operator, map[AbsolutePortAddr]chan core.Msg) error
	}
	Connector interface {
		Connect(map[AbsolutePortAddr]chan core.Msg, []Connection) error
	}
)

var (
	ErrDecoder           = errors.New("program decoder")
	ErrOpSpawner         = errors.New("operator-node spawner")
	ErrConstSpawner      = errors.New("const spawner")
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
	startPort, ok := ports[prog.StartPort]
	if !ok {
		return fmt.Errorf("%w: %v", ErrStartPortNotFound, prog.StartPort)
	}

	g := errgroup.Group{}

	g.Go(func() error {
		if err := r.connector.Connect(ports, prog.Connections); err != nil {
			return fmt.Errorf("%w: %v", ErrConnector, err)
		}
		return nil
	})
	g.Go(func() error {
		if err := r.opSpawner.Spawn(prog.Effects.Operators, ports); err != nil {
			return fmt.Errorf("%w: %v", ErrOpSpawner, err)
		}
		return nil
	})
	g.Go(func() error {
		if err := r.constSpawner.Spawn(prog.Effects.Constants, ports); err != nil {
			return fmt.Errorf("%w: %v", ErrConstSpawner, err)
		}
		return nil
	})

	if err := g.Wait(); err != nil {
		return fmt.Errorf("wait group: %w", err)
	}

	startPort <- core.NewStructMsg(nil)

	return nil
}

func MustNew(
	decoder Decoder,
	portGen PortGen,
	opSpawner OperatorSpawner,
	constSpawner ConstSpawner,
	connector Connector,
) Runtime {
	utils.PanicOnNil(decoder, portGen, opSpawner, constSpawner, connector)

	return Runtime{
		decoder:      decoder,
		portGen:      portGen,
		opSpawner:    opSpawner,
		constSpawner: constSpawner,
		connector:    connector,
	}
}
