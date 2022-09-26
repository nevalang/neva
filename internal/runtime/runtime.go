package runtime

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/emil14/neva/internal/core"
	"github.com/emil14/neva/internal/pkg/utils"
	"github.com/emil14/neva/internal/runtime/src"
	"golang.org/x/sync/errgroup"
)

type (
	Decoder interface {
		Decode([]byte) (src.Program, error)
	}

	PortGenerator interface {
		Ports([]src.Port) map[src.AbsolutePortAddr]chan core.Msg
	}

	ConstSpawner interface {
		Spawn(context.Context, []Const) error
	}
	Const struct {
		port chan core.Msg
		msg  src.Msg
	}

	OperatorSpawner interface {
		Spawn(context.Context, []Operator) error
	}
	Operator struct {
		Ref src.OperatorRef
		IO  OperatorIO
	}
	OperatorIO struct {
		Inports, Outports map[core.RelativePortAddr]chan core.Msg
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

	startPort, ok := ports[prog.StartPort]
	if !ok {
		return fmt.Errorf("%w: %v", ErrStartPortNotFound, prog.StartPort)
	}

	g, ctx := errgroup.WithContext(ctx)

	g.Go(func() error {
		if err := r.connector.Connect(ctx, ports, prog.Connections); err != nil {
			return fmt.Errorf("%w: %v", ErrConnector, err)
		}
		return nil
	})

	g.Go(func() error {
		if err := r.opSpawner.Spawn(ctx, prog.Effects.Operators, ports); err != nil {
			fmt.Println(err) // FIXME get rid
			return fmt.Errorf("%w: %v", ErrOpSpawner, err)
		}
		return nil
	})

	g.Go(func() error {
		if err := r.constSpawner.Spawn(ctx, prog.Effects.Constants, ports); err != nil {
			return fmt.Errorf("%w: %v", ErrConstSpawner, err)
		}
		return nil
	})

	g.Go(func() error {
		select {
		case startPort <- core.NewDictMsg(nil):
			return nil
		case <-time.After(time.Second): // FIXME deadline for all 3 jobs +  independent from program size and hardware
			return ErrStartPortBlocked
		}
	})

	if err := g.Wait(); err != nil { // FIXME all goroutines must respect context cancelation
		return fmt.Errorf("wait group: %w", err)
	}

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
