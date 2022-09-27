package runtime

import (
	"context"
	"errors"
	"fmt"

	"github.com/emil14/neva/internal/core"
	"github.com/emil14/neva/internal/pkg/utils"
	"github.com/emil14/neva/internal/runtime/src"
)

type (
	Decoder interface {
		Decode([]byte) (src.Program, error)
	}

	PortGenerator interface {
		Ports([]src.Port) map[src.AbsolutePortAddr]chan core.Msg // sure that we need an interface?
	}

	ConstSpawner interface {
		Spawn(context.Context, []Const) error
	}
	Const struct {
		Port chan core.Msg
		Msg  src.Msg
	}

	OperatorSpawner interface {
		Spawn(context.Context, []Operator) error
	}
	Operator struct {
		Ref src.OperatorRef
		IO  core.IO
	}

	Connector interface {
		Connect(context.Context, []Connection) error
	}
	Connection struct {
		sender    Sender
		receivers []Receiver
	}
	Sender struct {
		addr src.AbsolutePortAddr
		port chan core.Msg
	}
	Receiver struct {
		point src.ReceiverConnectionPoint
		port  chan core.Msg
	}
)

type Runtime struct {
	decoder      Decoder
	portGen      PortGenerator
	opSpawner    OperatorSpawner
	constSpawner ConstSpawner
	connector    Connector
}

var (
	ErrDecoder           = errors.New("program decoder")
	ErrStartPortNotFound = errors.New("start port not found")
	ErrConnector         = errors.New("connector")
	ErrOpSpawner         = errors.New("operator-node spawner")
	ErrConstSpawner      = errors.New("const spawner")
	ErrStartPortBlocked  = errors.New("start port blocked")
)

func (r Runtime) Run(ctx context.Context, bb []byte) error {
	prog, err := r.decoder.Decode(bb)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrDecoder, err)
	}

	ports := r.portGen.Ports(prog.Ports)

	_, ok := ports[prog.StartPort]
	if !ok {
		return fmt.Errorf("%w: %v", ErrStartPortNotFound, prog.StartPort)
	}

	// TODO

	return nil
}

func MustNew(
	decoder Decoder,
	portGen PortGenerator,
	opSpawner OperatorSpawner,
	constSpawner ConstSpawner,
	connector Connector,
) Runtime {
	utils.NilPanic(decoder, portGen, opSpawner, constSpawner, connector)

	return Runtime{
		decoder:      decoder,
		portGen:      portGen,
		opSpawner:    opSpawner,
		constSpawner: constSpawner,
		connector:    connector,
	}
}
